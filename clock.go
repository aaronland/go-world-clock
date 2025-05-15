package clock

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aaronland/go-world-clock/timezones"
)

type TimeResults struct {
	Label         string `json:"label"`
	TimeZone      string `json:"timezone"`
	DayOfWeek     string `json:"day_of_week"`
	Date          string `json:"date"`
	Time          string `json:"time"`
	UnixTimestamp int64  `json:"unix_timestamp"`
}

func TimeFromStrings(ctx context.Context, date string, tz string, locations ...string) ([]*TimeResults, error) {

	var source time.Time

	if date != "" {

		if tz == "" {
			return nil, fmt.Errorf("Missing timezone")
		}

		loc, err := time.LoadLocation(tz)

		if err != nil {
			return nil, fmt.Errorf("Failed to load timezone, %w", err)
		}

		t, err := time.ParseInLocation("2006-01-02 15:04", date, loc)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse date, %w", err)
		}

		source = t

	} else {
		now := time.Now()
		source = now.Local()
	}

	filters := &Filters{
		Timezones: locations,
	}

	clock_results, err := Time(ctx, source, filters)

	if err != nil {
		return nil, fmt.Errorf("Failed to query clock, %w", err)
	}

	// Filter and sort

	results := make([]*TimeResults, 0)

	source_zn, source_offset := source.Zone()
	seen := new(sync.Map)

	day_fmt := "Monday"
	date_fmt := "2006-01-02"
	time_fmt := "15:04"

	for _, r := range clock_results {

		r_zn, r_offset := r.Time.Zone()

		if r_zn == source_zn && r_offset == source_offset {

			if r.Timezone != tz {
				continue
			}
		}

		_, exists := seen.LoadOrStore(r.Timezone, true)

		if exists {
			continue
		}

		label_parts := strings.Split(r.Timezone, "/")
		label := fmt.Sprintf("%s (%s)", label_parts[1], label_parts[0])

		wasm_r := &TimeResults{
			Label:         label,
			TimeZone:      r.Timezone,
			DayOfWeek:     r.Time.Format(day_fmt),
			Date:          r.Time.Format(date_fmt),
			Time:          r.Time.Format(time_fmt),
			UnixTimestamp: r.Time.Unix(),
		}

		results = append(results, wasm_r)
	}

	// Sort in descending order (future to past)

	sort.Slice(results, func(i, j int) bool {
		return results[i].UnixTimestamp > results[j].UnixTimestamp
	})

	return results, nil
}

// Time will return zero or more Location records in other timezones for the time defined in source.
// Results may be limited by passing in a Filter instance with zero or more limits.
func Time(ctx context.Context, source time.Time, f *Filters) ([]*Location, error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	source_zn, source_offset := source.Zone()

	slog.Info("GET TIME", "source", source_zn, "offset", source_offset)

	tz_fs := timezones.FS
	tz_fh, err := tz_fs.Open("timezones.csv")

	if err != nil {
		return nil, fmt.Errorf("Failed to load timezones, %v", err)
	}

	csv_r := csv.NewReader(tz_fh)

	unsorted := make(map[int][]*Location)

	candidates_ch := make(chan *Location)
	done_ch := make(chan bool)

	wg := new(sync.WaitGroup)

	go func() {

		for {
			select {
			case <-ctx.Done():
				return
			case <-done_ch:
				return
			case r := <-candidates_ch:

				_, offset := r.Time.Zone()

				responses, exists := unsorted[offset]

				if !exists {
					responses = make([]*Location, 0)
				}

				responses = append(responses, r)
				unsorted[offset] = responses
			}
		}
	}()

	for {

		select {
		case <-ctx.Done():
			break
		default:
			// pass
		}

		row, err := csv_r.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("Failed to read row, %w", err)
		}

		wg.Add(1)

		go func(row []string) {

			defer wg.Done()

			select {
			case <-ctx.Done():
				return
			default:
				// pass
			}

			row_id := row[0]
			row_tz := row[1]
			row_zn := row[2]
			row_offset, _ := strconv.Atoi(row[3])

			if row_zn == source_zn || row_offset == source_offset {

				slog.Info("WTF", "row", row_zn, "source", source_zn)

				l := &Location{
					Time:     source,
					Id:       row_id,
					Timezone: row_tz,
				}

				candidates_ch <- l
				return
			}

			if f != nil {

				if len(f.Timezones) > 0 {

					ok := false

					for _, label := range f.Timezones {

						if strings.Contains(row_tz, label) {
							// slog.Info("OK", "row", row_tz, "label", label)
							ok = true
							break
						}
					}

					if !ok {
						return
					}
				}
			}

			loc, err := time.LoadLocation(row_tz)

			if err != nil {
				// err_ch <- fmt.Errorf("Failed to load '%s', %w", row_tz, err)
				return
			}

			there := source.In(loc)

			l := &Location{
				Time:     there,
				Id:       row_id,
				Timezone: row_tz,
			}

			candidates_ch <- l

		}(row)
	}

	wg.Wait()

	done_ch <- true

	//

	offsets := make([]int, 0)

	for i, _ := range unsorted {
		offsets = append(offsets, i)
	}

	sort.Ints(offsets)

	sorted := make([]*Location, 0)

	for _, i := range offsets {

		for _, r := range unsorted[i] {
			sorted = append(sorted, r)
		}
	}

	return sorted, nil
}

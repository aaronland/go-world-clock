package clock

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	_ "log"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/aaronland/go-world-clock/timezones"
)

// Time will return zero or more Location records in other timezones for the time defined in source.
// Results may be limited by passing in a Filter instance with zero or more limits.
func Time(ctx context.Context, source time.Time, f *Filters) ([]*Location, error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	source_zn, _ := source.Zone()

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

			if row_zn == source_zn {

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

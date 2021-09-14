package main

import (
	"encoding/csv"
	_ "fmt"
	"github.com/aaronland/go-world-clock/timezones"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
	"io"
	"log"
	"sort"
	"strings"
	"sync"
	"time"
)

func main() {

	fs := flagset.NewFlagSet("time")

	var labels multi.MultiString
	fs.Var(&labels, "in", "...")

	flagset.Parse(fs)

	tz_fs := timezones.FS
	fh, err := tz_fs.Open("timezones.csv")

	if err != nil {
		log.Fatalf("Failed to open timezones, %v")
	}

	defer fh.Close()

	now := time.Now()
	here := now.Local()
	zn, offset := here.Zone()

	csv_r := csv.NewReader(fh)

	wg := new(sync.WaitGroup)

	tmp := make(map[int][]time.Time)

	tmp[offset] = []time.Time{
		here,
	}

	candidates_ch := make(chan time.Time)
	done_ch := make(chan bool)

	go func() {

		for {
			select {
			case <-done_ch:
				return
			case t := <-candidates_ch:

				_, offset := t.Zone()

				offset_times, exists := tmp[offset]

				if !exists {
					offset_times = make([]time.Time, 0)
				}

				offset_times = append(offset_times, t)
				tmp[offset] = offset_times

			}
		}
	}()

	for {

		row, err := csv_r.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		wg.Add(1)

		go func(row []string) {

			defer wg.Done()

			row_tz := row[1]
			row_zn := row[2]

			if row_zn == zn {
				return
			}

			if len(labels) > 0 {

				ok := false

				for _, label := range labels {

					if strings.Contains(row_tz, label) {
						ok = true
						break
					}
				}

				if !ok {
					return
				}
			}

			loc, _ := time.LoadLocation(row_tz)

			there := here.In(loc)
			candidates_ch <- there

		}(row)
	}

	wg.Wait()

	done_ch <- true

	//

	offsets := make([]int, 0)

	for i, _ := range tmp {
		offsets = append(offsets, i)
	}

	sort.Ints(offsets)

	for _, i := range offsets {

		for _, t := range tmp[i] {
			log.Println(t)
		}
	}
}

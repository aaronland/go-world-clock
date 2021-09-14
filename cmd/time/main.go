package main

import (
	"encoding/csv"
	"fmt"
	"github.com/aaronland/go-world-clock/timezones"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
	"io"
	"log"
	"strings"
	"sync"
	"time"
)

func main() {

	fs := flagset.NewFlagSet("time")

	var wof_ids multi.MultiString
	fs.Var(&wof_ids, "wof-id", "...")

	var labels multi.MultiString
	fs.Var(&labels, "label", "...")

	flagset.Parse(fs)

	tz_fs := timezones.FS
	fh, err := tz_fs.Open("timezones.csv")

	if err != nil {
		log.Fatalf("Failed to open timezones, %v")
	}

	defer fh.Close()

	now := time.Now()
	here := now.Local()
	zn, _ := here.Zone()

	csv_r := csv.NewReader(fh)

	wg := new(sync.WaitGroup)

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

			if len(wof_ids) > 0 {

				row_id := row[0]
				ok := false

				for _, id := range wof_ids {

					if id == row_id {
						ok = true
						break
					}
				}

				if !ok {
					return
				}
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
			fmt.Printf("%s %v\n", row_tz, there)
		}(row)
	}

	wg.Wait()
}

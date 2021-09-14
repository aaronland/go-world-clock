package main

import (
	"encoding/csv"
	"fmt"
	"github.com/aaronland/go-world-clock/timezones"
	"io"
	"log"
	"sync"
	"time"
)

func main() {

	fs := timezones.FS
	fh, err := fs.Open("timezones.csv")

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

			loc, _ := time.LoadLocation(row_tz)

			there := here.In(loc)
			fmt.Printf("%s %v\n", row_tz, there)
		}(row)
	}

	wg.Wait()
}

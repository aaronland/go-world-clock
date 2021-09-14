// time will print the local time as well as the time in one or more timezones.
package main

import (
	"context"
	"fmt"
	"github.com/aaronland/go-world-clock"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
	"log"
	"os"
	"time"
)

func main() {

	fs := flagset.NewFlagSet("time")

	var in_timezones multi.MultiString
	fs.Var(&in_timezones, "in", "Zero or more strings to test whether they are contained by a given timezone's longform (major/minor) label.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Print the local time as well as the time in one or more timezones.\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n")
		fs.PrintDefaults()
	}

	flagset.Parse(fs)

	ctx := context.Background()

	filters := &clock.Filters{
		Timezones: in_timezones,
	}

	now := time.Now()
	here := now.Local()

	err := process(ctx, here, filters)

	if err != nil {
		log.Fatalf("Failed to process '%v', %v", here, err)
	}
}

func process(ctx context.Context, source time.Time, filters *clock.Filters) error {

	results, err := clock.Time(ctx, source, filters)

	if err != nil {
		return fmt.Errorf("Failed to determine time, %v", err)
	}

	zn, _ := source.Zone()
	seen := false

	d_fmt := "Monday"
	t_fmt := "2006-01-02 15:04"

	for _, r := range results {

		r_zn, _ := r.Time.Zone()

		if r_zn == zn {

			if !seen {
				label := "Local"
				label = padding(label, 24)

				d := r.Time.Format(d_fmt)
				d = padding(d, 10)

				t := r.Time.Format(t_fmt)

				str_t := fmt.Sprintf("%s %s", d, t)

				fmt.Printf("%s %s\t%s 👈\n", label, r_zn, str_t)
				seen = true
			}

			continue
		}

		label := r.Timezone
		label = padding(label, 24)

		d := r.Time.Format(d_fmt)
		d = padding(d, 10)

		t := r.Time.Format(t_fmt)

		str_t := fmt.Sprintf("%s %s", d, t)

		fmt.Printf("%s %s\t%s\n", label, r_zn, str_t)
	}

	return nil
}

func padding(input string, final int) string {

	input_len := len(input)
	padding := ""

	for i := 0; i < final-input_len; i++ {
		padding = padding + " "
	}

	return input + padding
}

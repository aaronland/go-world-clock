// time will print the local time as well as the time in one or more timezones.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"unicode/utf8"

	"github.com/aaronland/go-world-clock"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
)

func main() {

	fs := flagset.NewFlagSet("time")

	var in_timezones multi.MultiString
	fs.Var(&in_timezones, "in", "Zero or more strings to test whether they are contained by a given timezone's longform (major/minor) label.")

	date := fs.String("date", "", "YYYY-MM-dd HH:mm. If empty the current time in the computer's locale will be used.")
	timezone := fs.String("timezone", "", "A valid major/minor timezone location. Required if -date is not empty.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Print the local time as well as the time in one or more timezones.\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n")
		fs.PrintDefaults()
	}

	flagset.Parse(fs)

	ctx := context.Background()

	results, err := clock.TimeFromStrings(ctx, *date, *timezone, in_timezones...)

	if err != nil {
		log.Fatal(err)
	}

	for _, r := range results {

		label := padding(r.Label, 23)
		date := padding(r.Date, 12)
		dow := padding(r.DayOfWeek, 8)
		time := r.Time

		fmt.Printf("%s %s %s %s\n", label, date, dow, time)
	}

}

func padding(input string, final int) string {

	// input_len := len(input)
	input_len := utf8.RuneCountInString(input)

	// log.Println(input, input_len)
	padding := ""

	for i := 0; i < final-input_len; i++ {
		padding = padding + " "
	}

	//log.Println(input, len(padding))
	return input + padding
}

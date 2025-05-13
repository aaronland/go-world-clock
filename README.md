# go-world-clock

There are many world clocks. This one is mine.

## Documentation

Documentation is in progress and incomplete at this time.

[![Go Reference](https://pkg.go.dev/badge/github.com/aaronland/go-world-clock.svg)](https://pkg.go.dev/github.com/aaronland/go-world-clock)

## Usage

```
import (
	"log"
        "time"

	"github.com/aaronland/go-world-clock"
)

func main(){

	now := time.Now()
	here := now.Local()

	filters := &clock.Filters{
		Timezones: []string{ "Montreal" },
	}

	results, _ := clock.Time(ctx, here, filters)

	for _, r := range results {
		log.Printf("%s %s\n", r.Timezone, r.Time.Format(time.RFC3339))
	}
}
```

## Tools

### time

Print the local time as well as the time in one or more timezones.

```
$> ./bin/time -h
Print the local time as well as the time in one or more timezones.

Usage:
	 ./bin/time [options]

Valid options are:
  -date string
    	YYYY-MM-dd HH:mm. If empty the current time in the computer's locale will be used.
  -in value
    	Zero or more strings to test whether they are contained by a given timezone's longform (major/minor) label.
  -timezone string
    	A valid major/minor timezone location. Required if -date is not empty.
```

#### Example:

```
$> ./bin/time  -in Melbourne -in Cairo -in Honolulu -in China -in Montreal -in London
Pacific/Honolulu         HST	Tuesday    2021-09-14 13:38
👉                       PDT	Tuesday    2021-09-14 16:38 👈
America/Montreal         EDT	Tuesday    2021-09-14 19:38
Europe/London            BST	Wednesday  2021-09-15 00:38
Africa/Cairo             EET	Wednesday  2021-09-15 01:38
Australia/Melbourne      AEST	Wednesday  2021-09-15 09:38
```

Or specifying a custom date in a timezone:

```
$> ./bin/time  -date '2021-09-15 03:00' -timezone 'Asia/Singapore' -in London -in 'Los_Angeles'
America/Los_Angeles      PDT	Tuesday    2021-09-14 12:00
Europe/London            BST	Tuesday    2021-09-14 20:00
👉                       +08	Wednesday  2021-09-15 03:00 👈
```
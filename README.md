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
$> ./bin/time -date '2023-12-31 20:32' -timezone America/Montreal -in Europe/London -in Australia/Melbourne
Montreal (America)      2023-12-31   Sunday   20:32
London (Europe)         2024-01-01   Monday   01:32
Melbourne (Australia)   2024-01-01   Monday   12:32
```
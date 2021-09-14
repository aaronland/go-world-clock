# go-world-clock

There are many world clocks. This one is mine.

## Documentation

Documentation is in progress and incomplete at this time.

[![Go Reference](https://pkg.go.dev/badge/github.com/aaronland/go-world-clock.svg)](https://pkg.go.dev/github.com/aaronland/go-world-clock)

## Usage

```
import (
	"github.com/aaronland/go-world-clock"
	"log"
        "time"
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
> ./bin/time -h
Print the local time as well as the time in one or more timezones.

Usage:
	 ./bin/time [options]

Valid options are:
  -in value
    	Zero or more strings to test whether they are contained by a given timezone's longform (major/minor) label.
```

#### Example:

```
$> ./bin/time -in Melbourne -in Cairo -in Honolulu -in China -in Montreal -in London
Pacific/Honolulu         HST	Tuesday    2021-09-14 06:01
Local                    PDT	Tuesday    2021-09-14 09:01 👈
America/Montreal         EDT	Tuesday    2021-09-14 12:01
Europe/London            BST	Tuesday    2021-09-14 17:01
Africa/Cairo             EET	Tuesday    2021-09-14 18:01
Australia/Melbourne      AEST	Wednesday  2021-09-15 02:01
```
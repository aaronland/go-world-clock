//go:build wasmjs
package wasm

// See the way we are importing time/tzdata? That is important in order for
// time.LoadLocation to work in WASM-land. This fact isn't really documented
// anywhere but here: https://github.com/golang/go/issues/44408#issuecomment-1062548031

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"syscall/js"
	"time"
	_ "time/tzdata" 
	"strings"
	"sort"
	"sync"
	
	"github.com/aaronland/go-world-clock"	
)

type TimeFuncResults struct {
	Label string `json:"label"`
	TimeZone string `json:"timezone"`
	DayOfWeek string `json:"day_of_week"`
	Date string `json:"date"`
	Time string `json:"time"`	
	UnixTimestamp int64 `json:"unix_timestamp"`
}

func TimeFunc() js.Func {

	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		date := args[0].String()
		tz := args[1].String()
		str_locs := args[2].String()

		locations := strings.Split(str_locs, ",")
		
		logger := slog.Default()
		logger = logger.With("date", date)
		logger = logger.With("tz", tz)
		logger = logger.With("locations", str_locs)
		
		filters := &clock.Filters{
			Timezones: locations,
		}

		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {

			resolve := args[0]
			reject := args[1]

			ctx := context.Background()

			logger.Info("Lookup times")
			
			var source time.Time
			
			if date != "" {

				if tz == "" {
					logger.Error("Missing timezone")
					reject.Invoke(fmt.Printf("Missing timezone"))
					return nil					
				}

				// Failed to load timezone, %!w(syscall.Errno=38)2025/05/13 12:46:41 ERROR Failed load timezone date=2025-06-07T13:25 tz=Europe/London error="not implemented on js"
				loc, err := time.LoadLocation(tz)
				
				if err != nil {
					logger.Error("Failed load timezone", "error", err)
					reject.Invoke(fmt.Printf("Failed to load timezone, %w", err))
					return nil										
				}
				
				t, err := time.ParseInLocation("2006-01-02 15:04", date, loc)
				
				if err != nil {
					logger.Error("Failed to parse date", "error", err)
					reject.Invoke(fmt.Printf("Failed to parse date, %w", err))
					return nil										
				}
				
				source = t
				
			} else {
				now := time.Now()
				source = now.Local()
			}

			clock_results, err := clock.Time(ctx, source, filters)

			if err != nil {
				logger.Error("Failed to query clock", "error", err)
				reject.Invoke(fmt.Printf("Failed to query clock, %w", err))
				return nil										
			}

			// START OF put me in a function
			
			results := make([]*TimeFuncResults, 0)

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
				
				wasm_r := &TimeFuncResults{
					Label: label,
					TimeZone: r.Timezone,
					DayOfWeek: r.Time.Format(day_fmt),
					Date: r.Time.Format(date_fmt),
					Time: r.Time.Format(time_fmt),					
					UnixTimestamp: r.Time.Unix(),
				}
				
				results = append(results, wasm_r)
			}
			
			// Sort in descending order (future to past)
			
			sort.Slice(results, func(i, j int) bool {
				return results[i].UnixTimestamp > results[j].UnixTimestamp
			})

			// END OF put me in a function
			
			enc_results, err := json.Marshal(results)

			if err != nil {
				logger.Error("Failed to marshal results", "error", err)
				reject.Invoke(fmt.Printf("Failed to marshal results, %w", err))
				return nil										
			}

			resolve.Invoke(string(enc_results))
			return nil
		})

		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	})
}

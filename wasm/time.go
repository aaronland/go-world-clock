//go:build wasmjs
package wasm

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"syscall/js"
	"time"
	"strings"
	
	"github.com/aaronland/go-world-clock"	
)

type TimeFuncResults struct {
	Label string `json:"label"`
	TimeZone string `json:"timezone"`
	DateTime string `json:"datetime"`
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

		filters := &clock.Filters{
			Timezones: locations,
		}
		

		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {

			resolve := args[0]
			reject := args[1]

			ctx := context.Background()
			
			var source time.Time
			
			if date != "" {

				if tz == "" {
					logger.Error("Missing timezone")
					reject.Invoke(fmt.Printf("Missing timezone"))
					return nil					
				}
				
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

			results := make([]*TimeFuncResults, 0)

			// zn, _ := source.Zone()
			// seen := false
			
			d_fmt := "Monday"
			t_fmt := "2006-01-02 15:04"
			
			for _, r := range clock_results {

				// TBD...
				
				r_zn, _ := r.Time.Zone()
				
				d := r.Time.Format(d_fmt)
				
				t := r.Time.Format(t_fmt)
				str_t := fmt.Sprintf("%s %s", d, t)
				
				wasm_r := &TimeFuncResults{
					// Label: label,
					TimeZone: r_zn,
					DateTime: str_t,
				}
				
				results = append(results, wasm_r)
			}

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

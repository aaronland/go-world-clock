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
	_ "time/tzdata" 
	"strings"
	
	"github.com/aaronland/go-world-clock"	
)

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
		
		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {

			resolve := args[0]
			reject := args[1]

			ctx := context.Background()

			logger.Info("Lookup times")

			results, err := clock.TimeFromStrings(ctx, date, tz, locations...)

			if err != nil {
				logger.Error("Failed to determine times", "error", err)
				reject.Invoke(fmt.Printf("Failed to determine times, %w", err))
				return nil										
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

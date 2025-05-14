//go:build wasmjs
package wasm

// See the way we are importing time/tzdata? That is important in order for
// time.LoadLocation to work in WASM-land. This fact isn't really documented
// anywhere but here: https://github.com/golang/go/issues/44408#issuecomment-1062548031

import (
	"encoding/json"
	"encoding/csv"	
	"fmt"
	"log/slog"
	"io"
	"syscall/js"
	"strconv"
	
	"github.com/aaronland/go-world-clock/timezones"
)

type TimeZone struct {
	Name string `json:"name"`
	WhosOnFirstId int64 `json:"wof:id"`
}

func TimeZonesFunc() js.Func {

	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		logger := slog.Default()
		
		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {

			resolve := args[0]
			reject := args[1]
			
			tz_results := make([]*TimeZone, 0)
			
			tz_r, err := timezones.FS.Open("timezones.csv")

			if err != nil {
				logger.Error("Failed to open timezones", "error", err)
				reject.Invoke(fmt.Printf("File to open timezones, %w", err))
				return nil					
			}

			defer tz_r.Close()

			csv_r := csv.NewReader(tz_r)

			for {
				
				row, err := csv_r.Read()

				if err == io.EOF {
					break
				}
				
				if err != nil {
					logger.Error("Failed to iterate CSV row", "error", err)
					reject.Invoke(fmt.Printf("File to iterate CSV row, %w", err))
					return nil					
				}

				str_id := row[0]
				name := row[1]

				id, err := strconv.ParseInt(str_id, 10, 64)

				if err != nil {
					logger.Error("Failed to parse string ID, skipping", "id", str_id, "error", err)
					continue
				}

				tz := &TimeZone{
					Name: name,
					WhosOnFirstId: id,
				}

				tz_results = append(tz_results, tz)
			}
			
			enc_results, err := json.Marshal(tz_results)

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

//go:build wasmjs
package main

import (
	"log/slog"
	"syscall/js"
	
	"github.com/aaronland/go-world-clock/wasm"
)

func main() {

	time_func := wasm.TimeFunc()
	defer time_func.Release()

	js.Global().Set("world_clock_time", time_func)

	c := make(chan struct{}, 0)

	slog.Info("WASM world_clock_time function initialized")
	<-c
	
}

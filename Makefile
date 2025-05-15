GOROOT=$(shell go env GOROOT)
GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")

cli:
	go build -mod vendor -o bin/time cmd/time/main.go

tz:
	go run -mod vendor cmd/timezones/main.go > timezones/timezones.csv

wasmjs:
	GOOS=js GOARCH=wasm \
		go build -mod $(GOMOD) -ldflags="-s -w" -tags wasmjs \
		-o www/wasm/world_clock_time.wasm \
		cmd/time-wasm/main.go

server:
	fileserver -root www

debug:
	@make wasmjs
	@make server

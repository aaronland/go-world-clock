cli:
	go build -mod vendor -o bin/time cmd/time/main.go

tz:
	go run -mod vendor cmd/timezones/main.go > timezones/timezones.csv

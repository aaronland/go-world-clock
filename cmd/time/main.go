package main

import (
	"context"
	"fmt"
	"github.com/aaronland/go-world-clock"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
	"log"
	"time"
)

func main() {

	fs := flagset.NewFlagSet("time")

	var labels multi.MultiString
	fs.Var(&labels, "in", "...")

	flagset.Parse(fs)

	ctx := context.Background()

	now := time.Now()
	here := now.Local()

	results, err := clock.Time(ctx, here, labels...)

	if err != nil {
		log.Fatalf("Failed to determine time, %v", err)
	}

	for _, r := range results {
		fmt.Println(r)
	}
}

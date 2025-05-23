package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-iterate/v2/iterator"
)

func main() {

	iterator_uri := flag.String("iterator-uri", "repo://?include=properties.wof:placetype=timezone", "A valid whosonfirst/go-whosonfirst-iterate URI string.")
	iterator_source := flag.String("iterator-source", "/usr/local/data/whosonfirst/whosonfirst-data-admin-xy", "A valid whosonfirst/go-whosonfirst-iterate data source.")

	flag.Parse()

	ctx := context.Background()

	writers := []io.Writer{
		os.Stdout,
	}

	mw := io.MultiWriter(writers...)

	csv_wr := csv.NewWriter(mw)

	mu := new(sync.RWMutex)

	iter_cb := func(ctx context.Context, path string, r io.ReadSeeker, args ...interface{}) error {

		body, err := io.ReadAll(r)

		if err != nil {
			return fmt.Errorf("Failed to read body, %v", err)
		}

		id_rsp := gjson.GetBytes(body, "properties.wof:id")

		if !id_rsp.Exists() {
			return fmt.Errorf("Record is missing properties.wof:id")
		}

		name_rsp := gjson.GetBytes(body, "properties.wof:name")

		if !name_rsp.Exists() {
			return fmt.Errorf("Record is missing properties.wof:name")
		}

		id := id_rsp.Int()
		tz := name_rsp.String()

		loc, err := time.LoadLocation(tz)

		if err != nil {
			return fmt.Errorf("Failed to load location for '%s', %v", tz, err)
		}

		now := time.Now()
		rows := new(sync.Map)

		for i := 0; i < 12; i++ {

			ymd := fmt.Sprintf("%04d-%02d-01", now.Year(), i+1)
			dt, err := time.Parse("2006-01-02", ymd)

			if err != nil {
				return fmt.Errorf("Failed to parse %s, %w", ymd, err)
			}

			here := dt.In(loc)
			zn, offset := here.Zone()

			key := fmt.Sprintf("%s %s %d", tz, zn, offset)

			row := []string{
				strconv.FormatInt(id, 10),
				tz,
				zn,
				strconv.Itoa(offset),
			}

			rows.LoadOrStore(key, row)
		}

		mu.Lock()
		defer mu.Unlock()

		rows.Range(func(k interface{}, v interface{}) bool {

			key := k.(string)
			out := v.([]string)

			slog.Debug("Write timezone", "key", key)
			csv_wr.Write(out)
			csv_wr.Flush()
			return true
		})

		return nil
	}

	iter, err := iterator.NewIterator(ctx, *iterator_uri, iter_cb)

	if err != nil {
		log.Fatalf("Failed to create iterator, %v", err)
	}

	err = iter.IterateURIs(ctx, *iterator_source)

	if err != nil {
		log.Fatalf("Failed to iterate URIs, %v", err)
	}
}

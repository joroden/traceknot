package dashboard

import (
	"errors"
	"fmt"
	"time"
)

type Request struct {
	Range       string
	StartUnixMs int64
	EndUnixMs   int64
}

type Range struct {
	Key         string
	StartMs     int64
	EndMs       int64
	PrevStartMs int64
	PrevEndMs   int64
}

func Resolve(now time.Time, request Request) (Range, error) {
	switch request.Range {
	case "", "week":
		start := weekStart(now)
		return Range{
			Key:         "week",
			StartMs:     start.UnixMilli(),
			EndMs:       now.UnixMilli(),
			PrevStartMs: start.AddDate(0, 0, -7).UnixMilli(),
			PrevEndMs:   start.UnixMilli(),
		}, nil
	case "today":
		local := now.Local()
		day := time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, local.Location())
		start := day.UnixMilli()
		end := now.UnixMilli()
		return Range{
			Key:         "today",
			StartMs:     start,
			EndMs:       end,
			PrevStartMs: start - (end - start),
			PrevEndMs:   start,
		}, nil
	case "month":
		local := now.Local()
		monthStart := time.Date(local.Year(), local.Month(), 1, 0, 0, 0, 0, local.Location())
		return Range{
			Key:         "month",
			StartMs:     monthStart.UnixMilli(),
			EndMs:       now.UnixMilli(),
			PrevStartMs: monthStart.AddDate(0, -1, 0).UnixMilli(),
			PrevEndMs:   monthStart.UnixMilli(),
		}, nil
	case "all":
		return Range{
			Key:     "all",
			StartMs: 0,
			EndMs:   now.UnixMilli(),
		}, nil
	case "custom":
		if request.StartUnixMs <= 0 || request.EndUnixMs <= request.StartUnixMs {
			return Range{}, errors.New("custom range requires start_unix_ms < end_unix_ms")
		}
		length := request.EndUnixMs - request.StartUnixMs
		return Range{
			Key:         "custom",
			StartMs:     request.StartUnixMs,
			EndMs:       request.EndUnixMs,
			PrevStartMs: request.StartUnixMs - length,
			PrevEndMs:   request.StartUnixMs,
		}, nil
	default:
		return Range{}, fmt.Errorf("unknown range %q", request.Range)
	}
}

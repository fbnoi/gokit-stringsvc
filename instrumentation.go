package main

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
)

type instrumentingMiddleWare struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	countResult    metrics.Histogram
	next           IStringService
}

func instrumentingMiddleware(requestCount metrics.Counter, requestLatency metrics.Histogram, countResult metrics.Histogram) ServiceMiddleware {
	return func(next IStringService) IStringService {
		return instrumentingMiddleWare{requestCount, requestLatency, countResult, next}
	}
}

func (mw instrumentingMiddleWare) UpperCase(str string) (out string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "UpperCase", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	out, err = mw.next.UpperCase(str)
	return
}

func (mw instrumentingMiddleWare) Count(str string) (out int) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Count", "error", "false"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
		mw.countResult.Observe(float64(out))
	}(time.Now())
	out = mw.next.Count(str)
	return
}

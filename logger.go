package main

import (
	"time"

	"github.com/go-kit/kit/log"
)

// LoggingMiddleware LoggingMiddleware
type LoggingMiddleware struct {
	logger log.Logger
	next   IStringService
}

// UpperCase UpperCase
func (mw LoggingMiddleware) UpperCase(s string) (output string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "uppercase",
			"input", s,
			"output", output,
			"err", err,
			"took", time.Since(begin).Microseconds(),
		)
	}(time.Now())

	output, err = mw.next.UpperCase(s)
	return
}

// Count Count
func (mw LoggingMiddleware) Count(s string) (n int) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "count",
			"input", s,
			"n", n,
			"took", time.Since(begin).Microseconds(),
		)
	}(time.Now())

	n = mw.next.Count(s)
	return
}

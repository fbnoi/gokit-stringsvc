package main

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/sony/gobreaker"

	"golang.org/x/time/rate"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	"github.com/go-kit/kit/transport/http"
)

func proxyingmw(ctx context.Context, instances string, logger log.Logger) ServiceMiddleware {
	if instances == "" {
		logger.Log("proxy_to", "none")
		return func(next IStringService) IStringService {
			return next
		}
	}

	var (
		qps         = 100
		maxAttempts = 3
		maxTime     = 250 * time.Millisecond
	)

	var (
		instanceList = split(instances)
		ep           sd.FixedEndpointer
	)

	logger.Log("proxy_to", fmt.Sprint(instanceList))

	for _, instance := range instanceList {
		var e endpoint.Endpoint
		e = makeUpperCaseProxy(ctx, instance)
		e = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(e)
		e = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), qps))(e)
		ep = append(ep, e)
	}

	// Now, build a single, retrying, load-balancing endpoint out of all of
	// those individual endpoints.
	balancer := lb.NewRoundRobin(ep)
	retry := lb.Retry(maxAttempts, maxTime, balancer)

	// And finally, return the ServiceMiddleware, implemented by proxymw.
	return func(next IStringService) IStringService {
		return proxymw{ctx: ctx, next: next, proxy: retry}
	}
}

type proxymw struct {
	next  IStringService
	proxy endpoint.Endpoint //proxy upper case service to other service
	ctx   context.Context
}

func (mw proxymw) UpperCase(str string) (string, error) {
	response, err := mw.proxy(mw.ctx, upperCaseRequest{Str: str})
	if err != nil {
		return "", err
	}
	res := response.(upperCaseResponse)
	return res.Str, nil
}

func (mw proxymw) Count(str string) int {
	return mw.next.Count(str)
}

func makeUpperCaseProxy(ctx context.Context, instance string) endpoint.Endpoint {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}

	u, err := url.Parse(instance)
	if err != nil {
		panic(err)
	}

	if u.Path == "" {
		u.Path += "/uppercase"
	}
	return http.NewClient("GET", u, encodeRequest, decodeUpperCaseResponse).Endpoint()
}

func split(s string) []string {
	a := strings.Split(s, ",")
	for i := range a {
		a[i] = strings.TrimSpace(a[i])
	}
	return a
}

package transport

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/go-hclog"

	"github.com/okta/terraform-provider-okta/okta/internal/apimutex"
)

const (
	X_RATE_LIMIT_LIMIT     = "x-rate-limit-limit"
	X_RATE_LIMIT_REMAINING = "x-rate-limit-remaining"
	X_RATE_LIMIT_RESET     = "x-rate-limit-reset"
)

type GovernedTransport struct {
	base     http.RoundTripper
	apiMutex *apimutex.APIMutex
	logger   hclog.Logger
}

// NewGovernedTransport returns a governed transport that relies on pre- and post-
// requests from the http round tripper. The pre request consults the api mutex
// to determine if sleeping for the Okta API one minute bucket is called for.
// The post request updates the information it is holding about the current api
// rate limits.
func NewGovernedTransport(base http.RoundTripper, apiMutex *apimutex.APIMutex, logger hclog.Logger) *GovernedTransport {
	return &GovernedTransport{
		base:     base,
		apiMutex: apiMutex,
		logger:   logger,
	}
}

// RoundTrip returns the final http response after it has managed the api rate
// limit accounting in the pre and post request hooks.
func (t *GovernedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	path := req.URL.Path
	if err := t.preRequestHook(req.Context(), req.Method, path); err != nil {
		return nil, err
	}

	resp, err := t.base.RoundTrip(req)
	// always attempt to save x-headers
	t.postRequestHook(req.Method, path, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (t *GovernedTransport) preRequestHook(ctx context.Context, method, path string) error {
	if t.apiMutex.HasCapacity(method, path) {
		return nil
	}

	status := t.apiMutex.Status(method, path)
	now := time.Now().Unix()
	timeToSleep := status.Reset() - now

	line := fmt.Sprintf("Throttling API requests; sleeping for %d seconds until rate limit reset (path class %q: %d remaining of %d total); current request \"%s %s\"",
		timeToSleep,
		t.apiMutex.Class(method, path),
		status.Remaining(),
		status.Limit(),
		method,
		path,
	)
	t.logger.Info(line)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.NewTimer(time.Second * time.Duration(timeToSleep)).C:
		return nil
	}
}

func (t *GovernedTransport) postRequestHook(method, path string, resp *http.Response) {
	if resp == nil {
		return
	}
	reset, err := strconv.ParseInt(resp.Header.Get(X_RATE_LIMIT_RESET), 10, 64)
	if err != nil {
		t.logger.Warn(fmt.Sprintf("%q response header is missing or invalid, skipping postRequestHook: %+v", X_RATE_LIMIT_RESET, err))
		return
	}
	limit, err := strconv.Atoi(resp.Header.Get(X_RATE_LIMIT_LIMIT))
	if err != nil {
		t.logger.Warn(fmt.Sprintf("%q response header is missing or invalid, skipping postRequestHook: %+v", X_RATE_LIMIT_LIMIT, err))
		return
	}
	remaining, err := strconv.Atoi(resp.Header.Get(X_RATE_LIMIT_REMAINING))
	if err != nil {
		t.logger.Warn(fmt.Sprintf("%q response header is missing or invalid, skipping postRequestHook: %+v", X_RATE_LIMIT_REMAINING, err))
		return
	}

	t.apiMutex.Update(method, path, limit, remaining, reset)
}

package api

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/require"
)

func newTestRetryableClient(c *OktaAPIConfig) *http.Client {
	retryableClient := retryablehttp.NewClient()
	retryableClient.RetryWaitMin = 0
	retryableClient.RetryWaitMax = 0
	retryableClient.RetryMax = 3
	retryableClient.Logger = hclog.NewNullLogger()
	setRetryPolicy(retryableClient, c)
	return retryableClient.StandardClient()
}

func TestRetryPolicyDPoP(t *testing.T) {
	privateKeyConfig := &OktaAPIConfig{PrivateKey: "test-key"}
	apiTokenConfig := &OktaAPIConfig{ApiToken: "test-token"}
	cases := []struct {
		name         string
		config       *OktaAPIConfig
		status       int
		dpop         bool
		wantAttempts int32
	}{
		{name: "dpop 429 is not retried", config: privateKeyConfig, status: http.StatusTooManyRequests, dpop: true, wantAttempts: 1},
		{name: "dpop 503 is not retried", config: privateKeyConfig, status: http.StatusServiceUnavailable, dpop: true, wantAttempts: 1},
		{name: "private key without dpop 503 is retried", config: privateKeyConfig, status: http.StatusServiceUnavailable, wantAttempts: 4},
		{name: "api token 503 is retried", config: apiTokenConfig, status: http.StatusServiceUnavailable, wantAttempts: 4},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var attempts atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				attempts.Add(1)
				w.WriteHeader(tc.status)
			}))
			defer server.Close()

			req, err := http.NewRequest(http.MethodGet, server.URL, nil)
			require.NoError(t, err)
			if tc.dpop {
				req.Header.Set("Dpop", "proof")
			}
			resp, err := newTestRetryableClient(tc.config).Do(req)
			if tc.dpop {
				// The response is passed through for the SDK to handle.
				require.NoError(t, err)
				require.Equal(t, tc.status, resp.StatusCode)
				_ = resp.Body.Close()
			} else {
				require.Error(t, err)
			}
			require.Equal(t, tc.wantAttempts, attempts.Load())
		})
	}
}

func TestRetryPolicyDPoPConnectionError(t *testing.T) {
	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts.Add(1)
		hijacker, ok := w.(http.Hijacker)
		require.True(t, ok)
		conn, _, err := hijacker.Hijack()
		require.NoError(t, err)
		_ = conn.Close()
	}))
	defer server.Close()

	for _, dpop := range []bool{true, false} {
		attempts.Store(0)
		req, err := http.NewRequest(http.MethodGet, server.URL, nil)
		require.NoError(t, err)
		if dpop {
			req.Header.Set("Dpop", "proof")
		}
		resp, err := newTestRetryableClient(&OktaAPIConfig{PrivateKey: "test-key"}).Do(req)
		if resp != nil {
			_ = resp.Body.Close()
		}
		require.Error(t, err)
		if dpop {
			require.Equal(t, int32(1), attempts.Load())
		} else {
			require.Equal(t, int32(4), attempts.Load())
		}
	}
}

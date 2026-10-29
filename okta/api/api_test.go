package api

import (
	"net/http"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-retryablehttp"
)

func TestPrivateKeyAuthSkipsRetryableHTTP(t *testing.T) {
	tests := []struct {
		name            string
		backoff         bool
		privateKey      string
		expectRetryable bool
	}{
		{"backoff without private key uses retryablehttp", true, "", true},
		{"backoff with private key skips retryablehttp", true, "fake-private-key", false},
		{"no backoff without private key skips retryablehttp", false, "", false},
		{"no backoff with private key skips retryablehttp", false, "fake-private-key", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &OktaAPIConfig{
				Backoff:    tt.backoff,
				PrivateKey: tt.privateKey,
				OrgName:    "test",
				Domain:     "okta.com",
				MinWait:    30,
				MaxWait:    300,
				RetryCount: 5,
				Logger:     hclog.NewNullLogger(),
			}

			config, _, err := GetV3ClientConfig(c)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			assertTransport(t, config.HTTPClient, tt.expectRetryable)
		})
	}
}

func assertTransport(t *testing.T, client *http.Client, expectRetryable bool) {
	t.Helper()
	_, isRetryable := client.Transport.(*retryablehttp.RoundTripper)
	if expectRetryable && !isRetryable {
		t.Error("expected retryablehttp transport, got standard transport")
	}
	if !expectRetryable && isRetryable {
		t.Error("expected standard transport, got retryablehttp transport")
	}
}

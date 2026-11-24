package idaas

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	v6okta "github.com/okta/okta-sdk-golang/v6/okta"
)

// v6ErrorFromServer returns the error a real v6 SDK client produces when the
// Okta API responds with a 400 and the given content type and body. The SDK's
// GenericOpenAPIError has unexported fields, so a fake server is the only way
// to build a realistic one in a test.
func v6ErrorFromServer(t *testing.T, contentType, body string) error {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, body)
	}))
	t.Cleanup(server.Close)

	// The v6 SDK drops the port from the org URL, so route every connection
	// to the fake server regardless of the address the SDK dials.
	dialer := &net.Dialer{}
	httpClient := &http.Client{Transport: &http.Transport{
		DialContext: func(ctx context.Context, network, _ string) (net.Conn, error) {
			return dialer.DialContext(ctx, network, server.Listener.Addr().String())
		},
	}}

	config, err := v6okta.NewConfiguration(
		v6okta.WithOrgUrl(server.URL),
		v6okta.WithToken("test-token"),
		v6okta.WithAuthorizationMode("SSWS"),
		v6okta.WithTestingDisableHttpsCheck(true),
		v6okta.WithCache(false),
		v6okta.WithRateLimitMaxRetries(0),
		v6okta.WithHttpClientPtr(httpClient),
	)
	if err != nil {
		t.Fatalf("building v6 config: %v", err)
	}
	client := v6okta.NewAPIClient(config)

	app := v6okta.OpenIdConnectApplicationAsListApplications200ResponseInner(v6okta.NewOpenIdConnectApplicationWithDefaults())
	_, _, err = client.ApplicationAPI.ReplaceApplication(context.Background(), "0oa123").Application(app).Execute()
	if err == nil {
		t.Fatal("expected an error from the fake server, got nil")
	}
	return err
}

func TestOktaV6ErrorDetail(t *testing.T) {
	const withCauses = `{
		"errorCode": "E0000001",
		"errorSummary": "Api validation failed: redirect_uris",
		"errorId": "oae123",
		"errorCauses": [
			{"errorSummary": "'redirect_uris' is required when 'grant_types' includes 'authorization_code'."},
			{"errorSummary": "Duplicate redirect URI."}
		]
	}`
	const noCauses = `{"errorCode": "E0000001", "errorSummary": "Api validation failed: label"}`

	tests := []struct {
		name string
		err  func(t *testing.T) error
		want string
	}{
		{
			name: "okta error with causes",
			err: func(t *testing.T) error {
				return v6ErrorFromServer(t, "application/json", withCauses)
			},
			want: "Api validation failed: redirect_uris. Causes: 'redirect_uris' is required when 'grant_types' includes 'authorization_code'., Duplicate redirect URI.",
		},
		{
			name: "okta error without causes",
			err: func(t *testing.T) error {
				return v6ErrorFromServer(t, "application/json", noCauses)
			},
			want: "Api validation failed: label",
		},
		{
			name: "non-json body falls back to raw body",
			err: func(t *testing.T) error {
				return v6ErrorFromServer(t, "text/plain", "  upstream proxy rejected request\n")
			},
			want: "upstream proxy rejected request",
		},
		{
			name: "wrapped okta error is still unwrapped",
			err: func(t *testing.T) error {
				return fmt.Errorf("context: %w", v6ErrorFromServer(t, "application/json", noCauses))
			},
			want: "Api validation failed: label",
		},
		{
			name: "non-okta error returns empty string",
			err:  func(t *testing.T) error { return errors.New("dial tcp: connection refused") },
			want: "",
		},
		{
			name: "nil error returns empty string",
			err:  func(t *testing.T) error { return nil },
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := oktaV6ErrorDetail(tt.err(t))
			if got != tt.want {
				t.Errorf("oktaV6ErrorDetail()\n got: %q\nwant: %q", got, tt.want)
			}
		})
	}
}

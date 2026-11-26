package utils

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	v7okta "github.com/okta/okta-sdk-golang/v7/okta"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// GH-3043: an Okta 400 must surface its errorSummary and errorCauses, not only "400 Bad Request".
func TestAPIErrorDetail_GH3043(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"errorCode":"E0000001","errorSummary":"Api validation failed: screenLockType","errorId":"oae123","errorCauses":[{"errorSummary":"screenLockType: Invalid value for macOS: [PASSCODE]"}]}`))
	}))
	defer ts.Close()

	cfg, err := v7okta.NewConfiguration(
		v7okta.WithOrgUrl(ts.URL),
		v7okta.WithToken("test-token"),
		v7okta.WithTestingDisableHttpsCheck(true),
		v7okta.WithCache(false),
		v7okta.WithRateLimitMaxRetries(0),
	)
	require.NoError(t, err)
	cfg.Host = strings.TrimPrefix(ts.URL, "http://") // NewConfiguration keeps only the hostname, not the port
	client := v7okta.NewAPIClient(cfg)

	body := v7okta.NewDeviceAssuranceMacOSPlatformWithDefaults()
	body.SetName("example")
	body.SetPlatform("MACOS")
	_, _, apiErr := client.DeviceAssuranceAPI.CreateDeviceAssurancePolicy(context.Background()).
		DeviceAssurance(v7okta.ListDeviceAssurancePolicies200ResponseInner{DeviceAssuranceMacOSPlatform: body}).
		Execute()
	require.Error(t, apiErr)

	// The bug: the SDK error's message is only the HTTP status line.
	assert.Equal(t, "400 Bad Request", apiErr.Error())

	got := APIErrorDetail(apiErr)
	assert.Contains(t, got, "400 Bad Request")
	assert.Contains(t, got, "Api validation failed: screenLockType")
	assert.Contains(t, got, "E0000001")
	assert.Contains(t, got, "screenLockType: Invalid value for macOS: [PASSCODE]")
}

// A body that is not Okta's JSON error (e.g. a proxy page) is appended raw.
func TestAPIErrorDetail_RawBody(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte("blocked by upstream proxy"))
	}))
	defer ts.Close()

	cfg, err := v7okta.NewConfiguration(
		v7okta.WithOrgUrl(ts.URL),
		v7okta.WithToken("test-token"),
		v7okta.WithTestingDisableHttpsCheck(true),
		v7okta.WithCache(false),
		v7okta.WithRateLimitMaxRetries(0),
	)
	require.NoError(t, err)
	cfg.Host = strings.TrimPrefix(ts.URL, "http://")
	client := v7okta.NewAPIClient(cfg)

	_, _, apiErr := client.DeviceAssuranceAPI.GetDeviceAssurancePolicy(context.Background(), "dae123").Execute()
	require.Error(t, apiErr)
	assert.Contains(t, APIErrorDetail(apiErr), "blocked by upstream proxy")
}

func TestAPIErrorDetail_PlainError(t *testing.T) {
	assert.Equal(t, "boom", APIErrorDetail(errors.New("boom")))
	assert.Equal(t, "", APIErrorDetail(nil))
}

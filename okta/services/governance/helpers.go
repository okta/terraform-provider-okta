package governance

import (
	"errors"
	"strings"

	sdkgov "github.com/okta/okta-governance-sdk-golang/governance"
)

// APIErrorMessage returns a detailed error message from the governance SDK.
// When the SDK successfully decoded a ModelError it uses the structured fields
// (errorCode, errorSummary, errorCauses); otherwise it falls back to the raw
// response body so the caller always gets the most useful information available.
func APIErrorMessage(err error) string {
	if err == nil {
		return ""
	}

	var apiErr *sdkgov.GenericOpenAPIError
	if !errors.As(err, &apiErr) {
		return err.Error()
	}

	var parts []string
	parts = append(parts, err.Error()) // HTTP status, e.g. "400 Bad Request"

	if model := apiErr.Model(); model != nil {
		if modelErr, ok := model.(sdkgov.ModelError); ok {
			// Prefer structured fields over raw body.
			if code := strings.TrimSpace(modelErr.GetErrorCode()); code != "" {
				parts = append(parts, code)
			}
			if summary := strings.TrimSpace(modelErr.GetErrorSummary()); summary != "" {
				parts = append(parts, summary)
			}
			for _, cause := range modelErr.GetErrorCauses() {
				if cs := strings.TrimSpace(cause.GetErrorSummary()); cs != "" {
					parts = append(parts, cs)
				}
			}
			return strings.Join(parts, ": ")
		}
	}

	// Model was not decoded — fall back to the raw response body.
	if body := strings.TrimSpace(string(apiErr.Body())); body != "" {
		parts = append(parts, body)
	}

	return strings.Join(parts, ": ")
}

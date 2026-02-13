package governance

import (
	"errors"
	"strings"

	sdkgov "github.com/okta/okta-governance-sdk-golang/governance"
)

// APIErrorMessage returns a detailed error message from the governance SDK,
// including the API response body and decoded ModelError summary when available.
// Use this when reporting SDK errors in diagnostics so users see the API's message.
func APIErrorMessage(err error) string {
	if err == nil {
		return ""
	}
	msg := err.Error()
	var apiErr *sdkgov.GenericOpenAPIError
	if errors.As(err, &apiErr) {
		if body := apiErr.Body(); len(body) > 0 {
			bodyStr := strings.TrimSpace(string(body))
			if bodyStr != "" && !strings.Contains(msg, bodyStr) {
				msg += ". Response: " + bodyStr
			}
		}
		if model := apiErr.Model(); model != nil {
			if modelErr, ok := model.(sdkgov.ModelError); ok {
				if s := modelErr.GetErrorSummary(); s != "" && !strings.Contains(msg, s) {
					msg += ". " + s
				}
			}
		}
	}
	return msg
}

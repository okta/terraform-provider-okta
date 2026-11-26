package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// openAPIError is implemented by GenericOpenAPIError in every okta-sdk-golang major
// version. Its Error() is only the HTTP status line; Okta's explanation is in the
// decoded Model() or the raw Body().
type openAPIError interface {
	error
	Body() []byte
	Model() interface{}
}

type oktaAPIErrorPayload struct {
	ErrorCode    string `json:"errorCode"`
	ErrorSummary string `json:"errorSummary"`
	ErrorID      string `json:"errorId"`
	ErrorCauses  []struct {
		ErrorSummary string `json:"errorSummary"`
	} `json:"errorCauses"`
}

// maxRawErrorBody caps an undecodable response body appended to the message.
const maxRawErrorBody = 1024

// APIErrorDetail returns err's message, followed by Okta's errorSummary, errorCode,
// errorId and errorCauses when err comes from an okta-sdk-golang API call. Use it as
// the detail of a diagnostic instead of err.Error(), which for SDK errors is only the
// HTTP status line (e.g. "400 Bad Request").
func APIErrorDetail(err error) string {
	if err == nil {
		return ""
	}
	var apiErr openAPIError
	if !errors.As(err, &apiErr) {
		return err.Error()
	}

	var payload oktaAPIErrorPayload
	decoded := false
	if model := apiErr.Model(); model != nil {
		if raw, mErr := json.Marshal(model); mErr == nil && json.Unmarshal(raw, &payload) == nil && payload.ErrorSummary != "" {
			decoded = true
		}
	}
	if !decoded && len(apiErr.Body()) > 0 {
		if json.Unmarshal(apiErr.Body(), &payload) == nil && payload.ErrorSummary != "" {
			decoded = true
		}
	}

	var b strings.Builder
	b.WriteString(err.Error())
	if !decoded {
		body := strings.TrimSpace(string(apiErr.Body()))
		if body != "" {
			if len(body) > maxRawErrorBody {
				body = body[:maxRawErrorBody] + "..."
			}
			fmt.Fprintf(&b, ": %s", body)
		}
		return b.String()
	}

	fmt.Fprintf(&b, ": %s", payload.ErrorSummary)
	if payload.ErrorCode != "" {
		fmt.Fprintf(&b, " (%s)", payload.ErrorCode)
	}
	for _, cause := range payload.ErrorCauses {
		if cause.ErrorSummary != "" {
			fmt.Fprintf(&b, "\n  - %s", cause.ErrorSummary)
		}
	}
	if payload.ErrorID != "" {
		fmt.Fprintf(&b, "\nOkta error ID: %s", payload.ErrorID)
	}
	return b.String()
}

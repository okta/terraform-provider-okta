package sdk

import (
	"fmt"
	"strings"
)

type Error struct {
	ErrorMessage     string                   `json:"error"`
	ErrorDescription string                   `json:"error_description"`
	ErrorCode        string                   `json:"errorCode,omitempty"`
	ErrorSummary     string                   `json:"errorSummary,omitempty" toml:"error_description"`
	ErrorLink        string                   `json:"errorLink,omitempty"`
	ErrorId          string                   `json:"errorId,omitempty"`
	ErrorCauses      []map[string]interface{} `json:"errorCauses,omitempty"`
}

func (e *Error) Error() string {
	formattedErr := "the API returned an unknown error"
	if e.ErrorDescription != "" {
		formattedErr = fmt.Sprintf("the API returned an error: %s", e.ErrorDescription)
	} else if e.ErrorSummary != "" {
		formattedErr = fmt.Sprintf("the API returned an error: %s", e.ErrorSummary)
	}
	if len(e.ErrorCauses) > 0 {
		var causes []string
		for _, cause := range e.ErrorCauses {
			for key, val := range cause {
				causes = append(causes, fmt.Sprintf("%s: %v", key, val))
			}
		}
		formattedErr = fmt.Sprintf("%s. Causes: %s", formattedErr, strings.Join(causes, ", "))
	}
	return formattedErr
}

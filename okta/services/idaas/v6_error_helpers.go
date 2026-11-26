package idaas

import (
	"errors"
	"strings"

	v6okta "github.com/okta/okta-sdk-golang/v6/okta"
)

// oktaV6ErrorDetail returns the Okta API errorSummary and errorCauses
// from a v6 SDK error, or "" if none are available.
func oktaV6ErrorDetail(err error) string {
	var oae *v6okta.GenericOpenAPIError
	if !errors.As(err, &oae) {
		return ""
	}

	if oe, ok := oae.Model().(v6okta.Error); ok && oe.GetErrorSummary() != "" {
		var causes []string
		for _, cause := range oe.GetErrorCauses() {
			if s := cause.GetErrorSummary(); s != "" {
				causes = append(causes, s)
			}
		}
		if len(causes) == 0 {
			return oe.GetErrorSummary()
		}
		return oe.GetErrorSummary() + ". Causes: " + strings.Join(causes, ", ")
	}

	return strings.TrimSpace(string(oae.Body()))
}

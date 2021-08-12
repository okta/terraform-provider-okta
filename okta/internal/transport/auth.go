package transport

import (
	"net/http"
	"strings"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

type AuthTransport struct {
	base http.RoundTripper
}

// NewAuthTransport stops the provider execution in case Okta API returns either 401 or 403 error codes.
func NewAuthTransport(base http.RoundTripper) *AuthTransport {
	return &AuthTransport{
		base: base,
	}
}

// RoundTrip read the code based on the response from the API and terminates any further
// execution in case of unauthenticated or unauthorized requests.
// This can be invalid API token or a private key, no permissions to execute the request,
// no correct scopes were granted,etc.
func (t *AuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.base.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusBadRequest {
		oktaErr := okta.CheckResponseForError(resp).Error()
		if strings.Contains(oktaErr, "You are not allowed any of the requested scopes") ||
			strings.Contains(oktaErr, "Invalid value for 'client_id' parameter") {
			// panic here because hlog doesn't have Fatal method
			panic(oktaErr)
		}
	}
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		// panic here because hlog doesn't have Fatal method
		panic(okta.CheckResponseForError(resp))
	}
	return resp, nil
}

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-jose/go-jose/v4"
	"github.com/hashicorp/go-hclog"
	"github.com/okta/terraform-provider-okta/sdk"
)

// v3ClientAssertionTransport fixes private key (OAuth 2.0 service app) auth for
// the okta-sdk-golang v4 client that backs OktaSDKClientV3() when the org or
// app requires DPoP.
//
// In v4.1.2 the SDK sends the first token request without a DPoP proof. When
// Okta answers with invalid_dpop_proof (and later use_dpop_nonce), the SDK
// re-sends the same *http.Request with a DPoP header added, which means the
// same client_assertion JWT is presented again. Okta rejects the replayed
// assertion, and the SDK drops that error on the floor, so callers only see
// "Empty access token" (see issue #2658). Newer SDK majors (v5+) and the
// provider's own sdk package already sign a fresh assertion for the DPoP
// retry; v4 is frozen, so the provider does it here instead.
//
// For every token request that carries a DPoP proof and a client_assertion,
// the transport swaps in a freshly signed assertion (new jti, iat and exp).
// It also logs the token endpoint's error response so failures are not
// reduced to "Empty access token".
type v3ClientAssertionTransport struct {
	base     http.RoundTripper
	orgURL   string
	tokenURL *url.URL
	clientID string
	signer   jose.Signer
	logger   hclog.Logger
}

func newV3ClientAssertionTransport(base http.RoundTripper, orgURL, clientID, privateKey, privateKeyID string, logger hclog.Logger) (*v3ClientAssertionTransport, error) {
	if base == nil {
		base = http.DefaultTransport
	}
	tokenURL, err := url.Parse(orgURL + "/oauth2/v1/token")
	if err != nil {
		return nil, err
	}
	// Like the SDK's WithPrivateKey, accept either PEM content or a file path.
	if content, readErr := os.ReadFile(filepath.Clean(privateKey)); readErr == nil {
		privateKey = string(content)
	}
	signer, err := sdk.CreateKeySigner(privateKey, privateKeyID)
	if err != nil {
		return nil, err
	}
	return &v3ClientAssertionTransport{
		base:     base,
		orgURL:   orgURL,
		tokenURL: tokenURL,
		clientID: clientID,
		signer:   signer,
		logger:   logger,
	}, nil
}

func (t *v3ClientAssertionTransport) isTokenRequest(req *http.Request) bool {
	return req.Method == http.MethodPost &&
		req.URL != nil &&
		strings.EqualFold(req.URL.Host, t.tokenURL.Host) &&
		req.URL.Path == t.tokenURL.Path
}

func (t *v3ClientAssertionTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if !t.isTokenRequest(req) {
		return t.base.RoundTrip(req)
	}

	out := req
	query := req.URL.Query()
	hasProof := req.Header.Get("DPoP") != ""
	if hasProof && query.Get("client_assertion") != "" {
		assertion, err := sdk.CreateClientAssertion(t.orgURL, t.clientID, t.signer)
		if err != nil {
			return nil, fmt.Errorf("failed to sign client assertion for DPoP token request: %w", err)
		}
		query.Set("client_assertion", assertion)
		// RoundTrippers must not modify the caller's request.
		out = req.Clone(req.Context())
		u := *req.URL
		u.RawQuery = query.Encode()
		out.URL = &u
	}

	resp, err := t.base.RoundTrip(out)
	if err != nil || resp == nil || resp.StatusCode < 300 {
		return resp, err
	}
	t.logTokenError(resp, hasProof)
	return resp, nil
}

// logTokenError records Okta's OAuth error (error / error_description) for a
// failed token request. The body is restored so the SDK can still read it.
func (t *v3ClientAssertionTransport) logTokenError(resp *http.Response, hadProof bool) {
	if t.logger == nil || resp.Body == nil {
		return
	}
	body, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewReader(body))
	if err != nil {
		return
	}
	var oauthErr struct {
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
	}
	_ = json.Unmarshal(body, &oauthErr)
	// A missing proof on the first attempt and a nonce challenge are expected
	// steps of the DPoP handshake, not failures.
	expected := oauthErr.Error == "use_dpop_nonce" || (oauthErr.Error == "invalid_dpop_proof" && !hadProof)
	if expected {
		t.logger.Debug("okta token endpoint requested DPoP", "status", resp.StatusCode, "error", oauthErr.Error)
	} else {
		t.logger.Error("okta token endpoint rejected the access token request", "status", resp.StatusCode, "error", oauthErr.Error, "error_description", oauthErr.ErrorDescription)
	}
}

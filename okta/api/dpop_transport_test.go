package api

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/hashicorp/go-hclog"
	"github.com/okta/okta-sdk-golang/v4/okta"
)

const (
	testOrgURL   = "https://dpop-test.okta.com"
	testClientID = "0oaTESTCLIENT"
	testKeyID    = "test-kid"
)

func testPrivateKeyPEM(t *testing.T) (*rsa.PrivateKey, string) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	pemBytes := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	return key, string(pemBytes)
}

// fakeOktaOrg emulates the parts of an Okta org that matter for private key
// auth: the /oauth2/v1/token endpoint (optionally enforcing DPoP, RFC 9449)
// and one management API endpoint.
type fakeOktaOrg struct {
	t           *testing.T
	pub         *rsa.PublicKey
	requireDPoP bool

	mu             sync.Mutex
	seenAssertions map[string]bool
	tokenRequests  []tokenRequestRecord
	apiAuthHeaders []string
}

type tokenRequestRecord struct {
	assertion string
	hasProof  bool
	nonce     string
}

func newFakeOktaOrg(t *testing.T, pub *rsa.PublicKey, requireDPoP bool) *fakeOktaOrg {
	return &fakeOktaOrg{t: t, pub: pub, requireDPoP: requireDPoP, seenAssertions: map[string]bool{}}
}

func writeOAuthError(w http.ResponseWriter, status int, code, desc string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": code, "error_description": desc})
}

func (f *fakeOktaOrg) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && r.URL.Path == "/oauth2/v1/token":
		f.serveToken(w, r)
	case r.Method == http.MethodGet && r.URL.Path == "/api/v1/users/me":
		f.mu.Lock()
		f.apiAuthHeaders = append(f.apiAuthHeaders, r.Header.Get("Authorization"))
		f.mu.Unlock()
		want := "Bearer at-bearer"
		if f.requireDPoP {
			want = "DPoP at-dpop"
			if r.Header.Get("Dpop") == "" {
				writeOAuthError(w, http.StatusUnauthorized, "invalid_dpop_proof", "The DPoP proof JWT header is missing.")
				return
			}
		}
		if r.Header.Get("Authorization") != want {
			writeOAuthError(w, http.StatusUnauthorized, "invalid_token", "unexpected authorization header")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"00uME","status":"ACTIVE","profile":{"login":"admin@example.com"}}`))
	default:
		http.NotFound(w, r)
	}
}

func (f *fakeOktaOrg) serveToken(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm() // the v4 SDK sends parameters in the query string
	assertion := r.Form.Get("client_assertion")
	proof := r.Header.Get("Dpop")
	nonce := ""
	if proof != "" {
		parsed, err := jwt.ParseSigned(proof, []jose.SignatureAlgorithm{jose.RS256})
		if err != nil {
			f.t.Errorf("unparseable DPoP proof: %v", err)
		} else {
			var claims map[string]interface{}
			_ = parsed.UnsafeClaimsWithoutVerification(&claims)
			if htu, _ := claims["htu"].(string); htu != testOrgURL+"/oauth2/v1/token" {
				f.t.Errorf("DPoP proof htu = %q", htu)
			}
			nonce, _ = claims["nonce"].(string)
		}
	}

	f.mu.Lock()
	f.tokenRequests = append(f.tokenRequests, tokenRequestRecord{assertion: assertion, hasProof: proof != "", nonce: nonce})
	reused := f.seenAssertions[assertion]
	f.seenAssertions[assertion] = true
	f.mu.Unlock()

	// Okta authenticates the client before looking at DPoP, and a client
	// assertion is accepted only once.
	if reused {
		writeOAuthError(w, http.StatusUnauthorized, "invalid_client", "The client_assertion token has already been used.")
		return
	}
	if err := f.verifyAssertion(assertion); err != nil {
		writeOAuthError(w, http.StatusUnauthorized, "invalid_client", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if !f.requireDPoP {
		_, _ = w.Write([]byte(`{"token_type":"Bearer","access_token":"at-bearer","expires_in":3600,"scope":"okta.users.read"}`))
		return
	}
	if proof == "" {
		writeOAuthError(w, http.StatusBadRequest, "invalid_dpop_proof", "The DPoP proof JWT header is missing.")
		return
	}
	if nonce == "" {
		w.Header().Set("DPoP-Nonce", "server-nonce-1")
		writeOAuthError(w, http.StatusBadRequest, "use_dpop_nonce", "Authorization server requires nonce in DPoP proof.")
		return
	}
	_, _ = w.Write([]byte(`{"token_type":"DPoP","access_token":"at-dpop","expires_in":3600,"scope":"okta.users.read"}`))
}

func (f *fakeOktaOrg) verifyAssertion(assertion string) error {
	parsed, err := jwt.ParseSigned(assertion, []jose.SignatureAlgorithm{jose.RS256})
	if err != nil {
		return err
	}
	var claims jwt.Claims
	if err := parsed.Claims(f.pub, &claims); err != nil {
		return err
	}
	return claims.Validate(jwt.Expected{
		Issuer:      testClientID,
		Subject:     testClientID,
		AnyAudience: jwt.Audience{testOrgURL + "/oauth2/v1/token"},
	})
}

// redirectTransport sends requests addressed to testOrgURL to a local test
// server without changing what the SDK (and the transport under test) sees.
type redirectTransport struct {
	target *url.URL
}

func (rt redirectTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	out := req.Clone(req.Context())
	u := *req.URL
	u.Scheme = rt.target.Scheme
	u.Host = rt.target.Host
	out.URL = &u
	out.Host = ""
	return http.DefaultTransport.RoundTrip(out)
}

func testAPIConfig(privateKeyPEM string, logger hclog.Logger) *OktaAPIConfig {
	return &OktaAPIConfig{
		OrgName:        "dpop-test",
		Domain:         "okta.com",
		ClientID:       testClientID,
		PrivateKey:     privateKeyPEM,
		PrivateKeyId:   testKeyID,
		Scopes:         []string{"okta.users.read"},
		Logger:         logger,
		RequestTimeout: 30,
		RetryCount:     1,
		MaxWait:        1,
		MinWait:        1,
	}
}

// newTestV3Client builds the v3 client exactly like the provider does and
// points the innermost transport at the fake org.
func newTestV3Client(t *testing.T, cfg *OktaAPIConfig, server *httptest.Server) *okta.APIClient {
	t.Helper()
	config, _, err := GetV3ClientConfig(cfg)
	if err != nil {
		t.Fatal(err)
	}
	target, _ := url.Parse(server.URL)
	rt, ok := config.HTTPClient.Transport.(*v3ClientAssertionTransport)
	if !ok {
		t.Fatalf("v3 client transport is %T, want *v3ClientAssertionTransport", config.HTTPClient.Transport)
	}
	rt.base = redirectTransport{target: target}
	return okta.NewAPIClient(config)
}

func TestV3ClientPrivateKeyAuthAgainstDPoPEnforcedOrg(t *testing.T) {
	key, keyPEM := testPrivateKeyPEM(t)
	org := newFakeOktaOrg(t, &key.PublicKey, true)
	server := httptest.NewServer(org)
	defer server.Close()

	client := newTestV3Client(t, testAPIConfig(keyPEM, hclog.NewNullLogger()), server)
	user, _, err := client.UserAPI.GetUser(context.Background(), "me").Execute()
	if err != nil {
		t.Fatalf("GetUser against a DPoP-enforced org: %v", err)
	}
	if user.GetId() != "00uME" {
		t.Fatalf("unexpected user id %q", user.GetId())
	}

	org.mu.Lock()
	defer org.mu.Unlock()
	// proofless attempt -> invalid_dpop_proof, proof without nonce ->
	// use_dpop_nonce, proof with nonce -> DPoP-bound token.
	if len(org.tokenRequests) != 3 {
		t.Fatalf("token requests = %d, want 3: %+v", len(org.tokenRequests), org.tokenRequests)
	}
	if org.tokenRequests[0].hasProof || !org.tokenRequests[1].hasProof || org.tokenRequests[2].nonce != "server-nonce-1" {
		t.Fatalf("unexpected DPoP handshake: %+v", org.tokenRequests)
	}
	seen := map[string]bool{}
	for _, r := range org.tokenRequests {
		if seen[r.assertion] {
			t.Fatal("a client assertion was sent more than once")
		}
		seen[r.assertion] = true
	}
	if len(org.apiAuthHeaders) != 1 || org.apiAuthHeaders[0] != "DPoP at-dpop" {
		t.Fatalf("API authorization headers = %v", org.apiAuthHeaders)
	}
}

// Documents the underlying SDK behavior the transport works around: without
// it, okta-sdk-golang v4 replays its assertion and reports "Empty access token".
func TestV3ClientWithoutAssertionRefreshFailsAgainstDPoPEnforcedOrg(t *testing.T) {
	key, keyPEM := testPrivateKeyPEM(t)
	org := newFakeOktaOrg(t, &key.PublicKey, true)
	server := httptest.NewServer(org)
	defer server.Close()
	target, _ := url.Parse(server.URL)

	config, err := okta.NewConfiguration(
		okta.WithOrgUrl(testOrgURL),
		okta.WithCache(false),
		okta.WithHttpClientPtr(&http.Client{Transport: redirectTransport{target: target}}),
		okta.WithPrivateKey(keyPEM), okta.WithPrivateKeyId(testKeyID), okta.WithScopes([]string{"okta.users.read"}), okta.WithClientId(testClientID),
		okta.WithAuthorizationMode("PrivateKey"),
		okta.WithRateLimitMaxRetries(1),
	)
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = okta.NewAPIClient(config).UserAPI.GetUser(context.Background(), "me").Execute()
	if err == nil || !strings.Contains(err.Error(), "Empty access token") {
		t.Fatalf("expected the SDK's \"Empty access token\" error, got %v", err)
	}
}

func TestV3ClientPrivateKeyAuthAgainstNonDPoPOrg(t *testing.T) {
	key, keyPEM := testPrivateKeyPEM(t)
	org := newFakeOktaOrg(t, &key.PublicKey, false)
	server := httptest.NewServer(org)
	defer server.Close()

	client := newTestV3Client(t, testAPIConfig(keyPEM, hclog.NewNullLogger()), server)
	if _, _, err := client.UserAPI.GetUser(context.Background(), "me").Execute(); err != nil {
		t.Fatalf("GetUser against a non-DPoP org: %v", err)
	}
	org.mu.Lock()
	defer org.mu.Unlock()
	if len(org.tokenRequests) != 1 || org.tokenRequests[0].hasProof {
		t.Fatalf("token requests = %+v, want one proofless request", org.tokenRequests)
	}
	if len(org.apiAuthHeaders) != 1 || org.apiAuthHeaders[0] != "Bearer at-bearer" {
		t.Fatalf("API authorization headers = %v", org.apiAuthHeaders)
	}
}

func TestV3ClientLogsTokenEndpointErrors(t *testing.T) {
	_, keyPEM := testPrivateKeyPEM(t)
	otherKey, _ := testPrivateKeyPEM(t) // org trusts a different key
	org := newFakeOktaOrg(t, &otherKey.PublicKey, true)
	server := httptest.NewServer(org)
	defer server.Close()

	var logs bytes.Buffer
	logger := hclog.New(&hclog.LoggerOptions{Output: &logs, Level: hclog.Debug})
	client := newTestV3Client(t, testAPIConfig(keyPEM, logger), server)
	if _, _, err := client.UserAPI.GetUser(context.Background(), "me").Execute(); err == nil {
		t.Fatal("expected an error when the org rejects the client assertion")
	}
	if !strings.Contains(logs.String(), "invalid_client") {
		t.Fatalf("token endpoint error was not logged; logs:\n%s", logs.String())
	}
}

func TestV3ClientAssertionTransportScope(t *testing.T) {
	_, keyPEM := testPrivateKeyPEM(t)
	var got []*http.Request
	base := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		got = append(got, r)
		return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody, Request: r}, nil
	})
	rt, err := newV3ClientAssertionTransport(base, testOrgURL, testClientID, keyPEM, testKeyID, hclog.NewNullLogger())
	if err != nil {
		t.Fatal(err)
	}

	tokenURL := testOrgURL + "/oauth2/v1/token?client_assertion=original&grant_type=client_credentials"
	cases := []struct {
		name       string
		method     string
		url        string
		dpop       bool
		wantResign bool
	}{
		{"token request without proof", http.MethodPost, tokenURL, false, false},
		{"token request with proof", http.MethodPost, tokenURL, true, true},
		{"management API request with proof", http.MethodGet, testOrgURL + "/api/v1/users/me?client_assertion=original", true, false},
		{"token request to another host", http.MethodPost, "https://other.okta.com/oauth2/v1/token?client_assertion=original", true, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got = nil
			req, _ := http.NewRequest(tc.method, tc.url, nil)
			if tc.dpop {
				req.Header.Set("DPoP", "proof")
			}
			resp, err := rt.RoundTrip(req)
			if err != nil {
				t.Fatal(err)
			}
			_ = resp.Body.Close()
			if req.URL.Query().Get("client_assertion") != "original" {
				t.Fatal("caller's request was modified")
			}
			sent := got[0].URL.Query().Get("client_assertion")
			if resigned := sent != "original"; resigned != tc.wantResign {
				t.Fatalf("assertion re-signed = %v, want %v", resigned, tc.wantResign)
			}
			if tc.wantResign && got[0].URL.Query().Get("grant_type") != "client_credentials" {
				t.Fatal("other token request parameters were dropped")
			}
		})
	}
}

func TestGetV3ClientConfigAssertionRefreshWiring(t *testing.T) {
	_, keyPEM := testPrivateKeyPEM(t)
	keyFile := filepath.Join(t.TempDir(), "key.pem")
	if err := os.WriteFile(keyFile, []byte(keyPEM), 0o600); err != nil {
		t.Fatal(err)
	}

	for name, tc := range map[string]struct {
		mutate   func(*OktaAPIConfig)
		wantWrap bool
	}{
		"private key PEM":       {func(c *OktaAPIConfig) {}, true},
		"private key file path": {func(c *OktaAPIConfig) { c.PrivateKey = keyFile }, true},
		"api token":             {func(c *OktaAPIConfig) { c.PrivateKey = ""; c.ApiToken = "ssws" }, false},
		"access token wins":     {func(c *OktaAPIConfig) { c.AccessToken = "bearer" }, false},
		"backoff client":        {func(c *OktaAPIConfig) { c.Backoff = true }, true},
	} {
		t.Run(name, func(t *testing.T) {
			cfg := testAPIConfig(keyPEM, hclog.NewNullLogger())
			tc.mutate(cfg)
			config, _, err := GetV3ClientConfig(cfg)
			if err != nil {
				t.Fatal(err)
			}
			_, wrapped := config.HTTPClient.Transport.(*v3ClientAssertionTransport)
			if wrapped != tc.wantWrap {
				t.Fatalf("wrapped = %v, want %v", wrapped, tc.wantWrap)
			}
		})
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// dpopOrg is a fake Okta org that only issues DPoP-bound access tokens. It
// records every DPoP proof and client assertion jti it receives.
type dpopOrg struct {
	t *testing.T
	// tokenRateLimits is how many token requests that carry a DPoP proof and
	// nonce are rate limited before a token is issued.
	tokenRateLimits int
	// groupRateLimits is how many create group requests are rate limited
	// before the group is created.
	groupRateLimits int
	// rateLimitFirstAttempts rate limits the first attempt of every create
	// group request.
	rateLimitFirstAttempts bool

	mu            sync.Mutex
	tokenRequests int
	groupRequests int
	jtis          map[string]int
}

const testDpopNonce = "test-nonce"

func (o *dpopOrg) recordJTI(t *testing.T, token string) {
	parts := strings.Split(token, ".")
	require.Len(t, parts, 3)
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	require.NoError(t, err)
	var claims struct {
		ID string `json:"jti"`
	}
	require.NoError(t, json.Unmarshal(payload, &claims))
	require.NotEmpty(t, claims.ID)
	o.jtis[claims.ID]++
}

func rateLimited(w http.ResponseWriter) {
	now := time.Now()
	w.Header().Set("Date", now.UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT"))
	w.Header().Set("X-Rate-Limit-Reset", fmt.Sprint(now.Unix()))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusTooManyRequests)
	_, _ = w.Write([]byte(`{"errorCode":"E0000047","errorSummary":"API call exceeded rate limit due to too many requests."}`))
}

func (o *dpopOrg) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	o.mu.Lock()
	defer o.mu.Unlock()
	t := o.t
	dpop := r.Header.Get("Dpop")
	if dpop != "" {
		o.recordJTI(t, dpop)
	}

	switch r.URL.Path {
	case "/oauth2/v1/token":
		o.tokenRequests++
		require.NoError(t, r.ParseForm())
		o.recordJTI(t, r.PostForm.Get("client_assertion"))
		w.Header().Set("Content-Type", "application/json")
		switch {
		case dpop == "":
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":"invalid_dpop_proof","error_description":"The DPoP proof JWT header is missing."}`))
		case !strings.Contains(decodedClaims(t, dpop), testDpopNonce):
			w.Header().Set("Dpop-Nonce", testDpopNonce)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":"use_dpop_nonce","error_description":"Authorization server requires nonce in DPoP proof."}`))
		case o.tokenRateLimits > 0:
			o.tokenRateLimits--
			rateLimited(w)
		default:
			_, _ = w.Write([]byte(`{"token_type":"DPoP","access_token":"test-access-token","expires_in":3600,"scope":"okta.groups.manage"}`))
		}
	case "/api/v1/groups":
		o.groupRequests++
		require.Equal(t, "DPoP test-access-token", r.Header.Get("Authorization"))
		require.NotEmpty(t, dpop)
		if o.rateLimitFirstAttempts && r.Header.Get("X-Okta-Retry-Count") == "" {
			rateLimited(w)
			return
		}
		if o.groupRateLimits > 0 {
			o.groupRateLimits--
			rateLimited(w)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"00g1","profile":{"name":"test"}}`))
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func decodedClaims(t *testing.T, token string) string {
	parts := strings.Split(token, ".")
	require.Len(t, parts, 3)
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	require.NoError(t, err)
	return string(payload)
}

func (o *dpopOrg) requireUniqueJTIs(t *testing.T) {
	o.mu.Lock()
	defer o.mu.Unlock()
	for jti, count := range o.jtis {
		require.Equalf(t, 1, count, "jti %s was sent %d times", jti, count)
	}
}

func testPrivateKeyPEM(t *testing.T) string {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	return string(pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}))
}

func newDpopTestClient(t *testing.T, org *dpopOrg) *Client {
	org.t = t
	org.jtis = map[string]int{}
	server := httptest.NewServer(org)
	t.Cleanup(server.Close)

	_, client, err := NewClient(
		context.Background(),
		WithOrgUrl(server.URL),
		WithTestingDisableHttpsCheck(true),
		WithCache(false),
		WithAuthorizationMode("PrivateKey"),
		WithClientId("test-client-id"),
		WithScopes([]string{"okta.groups.manage"}),
		WithPrivateKey(testPrivateKeyPEM(t)),
		WithRateLimitMaxRetries(2),
		WithRateLimitMaxBackOff(0),
	)
	require.NoError(t, err)
	return client
}

func TestPrivateKeyAuthDpopTokenRateLimited(t *testing.T) {
	org := &dpopOrg{tokenRateLimits: 2}
	client := newDpopTestClient(t, org)

	group, _, err := client.Group.CreateGroup(context.Background(), Group{Profile: &GroupProfile{Name: "test"}})
	require.NoError(t, err)
	require.Equal(t, "00g1", group.Id)
	// invalid_dpop_proof, use_dpop_nonce, two rate limited requests, success
	require.Equal(t, 5, org.tokenRequests)
	org.requireUniqueJTIs(t)
}

func TestPrivateKeyAuthDpopTokenRateLimitExhausted(t *testing.T) {
	org := &dpopOrg{tokenRateLimits: 10}
	client := newDpopTestClient(t, org)

	_, _, err := client.Group.CreateGroup(context.Background(), Group{Profile: &GroupProfile{Name: "test"}})
	require.Error(t, err)
	require.NotContains(t, err.Error(), "empty access token")
	require.Contains(t, err.Error(), "status 429")
	// invalid_dpop_proof, use_dpop_nonce, then the first request plus two retries
	require.Equal(t, 5, org.tokenRequests)
	org.requireUniqueJTIs(t)
}

func TestPrivateKeyAuthDpopRequestRateLimited(t *testing.T) {
	org := &dpopOrg{groupRateLimits: 2}
	client := newDpopTestClient(t, org)

	group, _, err := client.Group.CreateGroup(context.Background(), Group{Profile: &GroupProfile{Name: "test"}})
	require.NoError(t, err)
	require.Equal(t, "00g1", group.Id)
	require.Equal(t, 3, org.groupRequests)
	// A rate limited request does not invalidate the access token.
	require.Equal(t, 3, org.tokenRequests)
	org.requireUniqueJTIs(t)
}

func TestPrivateKeyAuthDpopConcurrentRequestsShareToken(t *testing.T) {
	org := &dpopOrg{rateLimitFirstAttempts: true}
	client := newDpopTestClient(t, org)

	const requests = 20
	var wg sync.WaitGroup
	errs := make(chan error, requests)
	for range requests {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _, err := client.Group.CreateGroup(context.Background(), Group{Profile: &GroupProfile{Name: "test"}})
			errs <- err
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		require.NoError(t, err)
	}
	// One token for all requests: invalid_dpop_proof, use_dpop_nonce, success.
	require.Equal(t, 3, org.tokenRequests)
	require.Equal(t, 2*requests, org.groupRequests)
	org.requireUniqueJTIs(t)
}

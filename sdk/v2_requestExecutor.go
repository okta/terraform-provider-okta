// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK

package sdk

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	urlpkg "net/url"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/BurntSushi/toml"
	"github.com/cenkalti/backoff/v4"
	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/okta/terraform-provider-okta/sdk/cache"
	goCache "github.com/patrickmn/go-cache"
)

const (
	AccessTokenCacheKey = "OKTA_ACCESS_TOKEN"
	// maxDpopNonceAttempts bounds how many times a new DPoP nonce is accepted
	// from the token endpoint for a single token request.
	maxDpopNonceAttempts = 3
)

type RequestExecutor struct {
	httpClient *http.Client
	config     *config
	BaseUrl    *urlpkg.URL
	cache      cache.Cache
	tokenCache *goCache.Cache
	// tokenLock is shared by cloned executors, like tokenCache, so concurrent
	// requests that miss the cache make a single token request.
	tokenLock         *sync.Mutex
	binary            bool
	headerAccept      string
	headerContentType string
	freshCache        bool
}

type ClientAssertionClaims struct {
	Issuer   string           `json:"iss,omitempty"`
	Subject  string           `json:"sub,omitempty"`
	Audience string           `json:"aud,omitempty"`
	Expiry   *jwt.NumericDate `json:"exp,omitempty"`
	IssuedAt *jwt.NumericDate `json:"iat,omitempty"`
	ID       string           `json:"jti,omitempty"`
}

type RequestAccessToken struct {
	TokenType   string `json:"token_type,omitempty"`
	ExpiresIn   int    `json:"expires_in,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
	Scope       string `json:"scope,omitempty"`
}

type Authorization interface {
	Authorize(method, URL string) error
}

type SSWSAuth struct {
	token string
	req   *http.Request
}

func NewSSWSAuth(token string, req *http.Request) *SSWSAuth {
	return &SSWSAuth{token: token, req: req}
}

func (a *SSWSAuth) Authorize(method, URL string) error {
	a.req.Header.Add("Authorization", "SSWS "+a.token)
	return nil
}

type BearerAuth struct {
	token string
	req   *http.Request
}

func NewBearerAuth(token string, req *http.Request) *BearerAuth {
	return &BearerAuth{token: token, req: req}
}

func (a *BearerAuth) Authorize(method, URL string) error {
	a.req.Header.Add("Authorization", "Bearer "+a.token)
	return nil
}

type PrivateKeyAuth struct {
	tokenCache       *goCache.Cache
	tokenLock        *sync.Mutex
	httpClient       *http.Client
	privateKeySigner jose.Signer
	privateKey       string
	privateKeyId     string
	clientId         string
	orgURL           string
	scopes           []string
	maxRetries       int32
	maxBackoff       int64
	req              *http.Request
}

type PrivateKeyAuthConfig struct {
	TokenCache       *goCache.Cache
	TokenLock        *sync.Mutex
	HttpClient       *http.Client
	PrivateKeySigner jose.Signer
	PrivateKey       string
	PrivateKeyId     string
	ClientId         string
	OrgURL           string
	Scopes           []string
	MaxRetries       int32
	MaxBackoff       int64
	Req              *http.Request
}

func NewPrivateKeyAuth(config PrivateKeyAuthConfig) *PrivateKeyAuth {
	return &PrivateKeyAuth{
		tokenCache:       config.TokenCache,
		tokenLock:        config.TokenLock,
		httpClient:       config.HttpClient,
		privateKeySigner: config.PrivateKeySigner,
		privateKey:       config.PrivateKey,
		privateKeyId:     config.PrivateKeyId,
		clientId:         config.ClientId,
		orgURL:           config.OrgURL,
		scopes:           config.Scopes,
		maxRetries:       config.MaxRetries,
		maxBackoff:       config.MaxBackoff,
		req:              config.Req,
	}
}

func (a *PrivateKeyAuth) Authorize(method, URL string) error {
	return authorizeWithAccessToken(a.req, a.tokenCache, a.tokenLock, method, URL, func() (*RequestAccessToken, string, *rsa.PrivateKey, error) {
		if a.privateKeySigner == nil {
			var err error
			a.privateKeySigner, err = CreateKeySigner(a.privateKey, a.privateKeyId)
			if err != nil {
				return nil, "", nil, err
			}
		}

		clientAssertion, err := CreateClientAssertion(a.orgURL, a.clientId, a.privateKeySigner)
		if err != nil {
			return nil, "", nil, err
		}

		return getAccessTokenForPrivateKey(a.httpClient, a.orgURL, clientAssertion, a.scopes, a.maxRetries, a.maxBackoff, a.clientId, a.privateKeySigner)
	})
}

type JWTAuth struct {
	tokenCache      *goCache.Cache
	tokenLock       *sync.Mutex
	httpClient      *http.Client
	orgURL          string
	scopes          []string
	clientAssertion string
	maxRetries      int32
	maxBackoff      int64
	req             *http.Request
}

type JWTAuthConfig struct {
	TokenCache      *goCache.Cache
	TokenLock       *sync.Mutex
	HttpClient      *http.Client
	OrgURL          string
	Scopes          []string
	ClientAssertion string
	MaxRetries      int32
	MaxBackoff      int64
	Req             *http.Request
}

func NewJWTAuth(config JWTAuthConfig) *JWTAuth {
	return &JWTAuth{
		tokenCache:      config.TokenCache,
		tokenLock:       config.TokenLock,
		httpClient:      config.HttpClient,
		orgURL:          config.OrgURL,
		scopes:          config.Scopes,
		clientAssertion: config.ClientAssertion,
		maxRetries:      config.MaxRetries,
		maxBackoff:      config.MaxBackoff,
		req:             config.Req,
	}
}

func (a *JWTAuth) Authorize(method, URL string) error {
	return authorizeWithAccessToken(a.req, a.tokenCache, a.tokenLock, method, URL, func() (*RequestAccessToken, string, *rsa.PrivateKey, error) {
		return getAccessTokenForPrivateKey(a.httpClient, a.orgURL, a.clientAssertion, a.scopes, a.maxRetries, a.maxBackoff, "", nil)
	})
}

// cachedAccessToken is stored as a single tokenCache entry so an access token
// is never paired with the DPoP nonce or signing key of a different token.
type cachedAccessToken struct {
	tokenType   string
	accessToken string
	dpopNonce   string
	dpopKey     *rsa.PrivateKey
}

type accessTokenFetcher func() (*RequestAccessToken, string, *rsa.PrivateKey, error)

// authorizeWithAccessToken sets the Authorization header on req and, for
// DPoP-bound tokens, a DPoP proof. A new proof (with a unique jti) is generated
// on every call, so callers re-authorize a request before re-sending it.
func authorizeWithAccessToken(req *http.Request, tokenCache *goCache.Cache, tokenLock *sync.Mutex, method, URL string, fetch accessTokenFetcher) error {
	token, err := getOrFetchAccessToken(tokenCache, tokenLock, fetch)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("%v %v", token.tokenType, token.accessToken))
	if token.dpopKey != nil {
		dpopJWT, err := generateDpopJWT(token.dpopKey, method, URL, token.dpopNonce, token.accessToken)
		if err != nil {
			return err
		}
		req.Header.Set("Dpop", dpopJWT)
		req.Header.Set("x-okta-user-agent-extended", "isDPoP:true")
	}
	return nil
}

func getCachedAccessToken(tokenCache *goCache.Cache) (*cachedAccessToken, bool) {
	cached, ok := tokenCache.Get(AccessTokenCacheKey)
	if !ok {
		return nil, false
	}
	token, ok := cached.(*cachedAccessToken)
	return token, ok
}

// getOrFetchAccessToken returns the cached access token, fetching a new one on
// a cache miss. Requests that miss the cache at the same time wait for a single
// fetch instead of each calling the token endpoint.
func getOrFetchAccessToken(tokenCache *goCache.Cache, tokenLock *sync.Mutex, fetch accessTokenFetcher) (*cachedAccessToken, error) {
	if token, ok := getCachedAccessToken(tokenCache); ok {
		return token, nil
	}

	tokenLock.Lock()
	defer tokenLock.Unlock()
	if token, ok := getCachedAccessToken(tokenCache); ok {
		return token, nil
	}

	accessToken, nonce, privateKey, err := fetch()
	if err != nil {
		return nil, err
	}
	if accessToken == nil || accessToken.AccessToken == "" {
		return nil, errors.New("empty access token")
	}

	token := &cachedAccessToken{
		tokenType:   accessToken.TokenType,
		accessToken: accessToken.AccessToken,
	}
	if accessToken.TokenType == "DPoP" {
		if privateKey == nil {
			return nil, errors.New("using Dpop but signing key not found")
		}
		token.dpopNonce = nonce
		token.dpopKey = privateKey
	}

	// Trim a couple of seconds off calculated expiry so cache expiry
	// occures before Okta server side expiry.
	expiration := accessToken.ExpiresIn - 2
	if expiration > 0 {
		tokenCache.Set(AccessTokenCacheKey, token, time.Second*time.Duration(expiration))
	}
	return token, nil
}

func CreateKeySigner(privateKey, privateKeyID string) (jose.Signer, error) {
	var signerOptions *jose.SignerOptions
	if privateKeyID != "" {
		signerOptions = (&jose.SignerOptions{}).WithHeader("kid", privateKeyID)
	}

	priv := []byte(strings.ReplaceAll(privateKey, `\n`, "\n"))

	privPem, _ := pem.Decode(priv)
	if privPem == nil {
		return nil, errors.New("invalid private key")
	}

	if privPem.Type == "RSA PRIVATE KEY" {
		parsedKey, err := x509.ParsePKCS1PrivateKey(privPem.Bytes)
		if err != nil {
			return nil, err
		}
		return jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: parsedKey}, signerOptions)
	}
	if privPem.Type == "PRIVATE KEY" {
		parsedKey, err := x509.ParsePKCS8PrivateKey(privPem.Bytes)
		if err != nil {
			return nil, err
		}
		var alg jose.SignatureAlgorithm
		switch parsedKey.(type) {
		case *rsa.PrivateKey:
			alg = jose.RS256
		case *ecdsa.PrivateKey:
			alg = jose.ES256 // TODO handle ES384 or ES512 ?
		default:
			// TODO are either of these also valid?
			// ed25519.PrivateKey:
			// *ecdh.PrivateKey
			return nil, fmt.Errorf("private key %q is unknown pkcs#8 format type", privPem.Type)
		}
		return jose.NewSigner(jose.SigningKey{Algorithm: alg, Key: parsedKey}, signerOptions)
	}

	return nil, fmt.Errorf("private key %q is not pkcs#1 or pkcs#8 format", privPem.Type)
}

func CreateClientAssertion(orgURL, clientID string, privateKeySinger jose.Signer) (clientAssertion string, err error) {
	claims := ClientAssertionClaims{
		Subject:  clientID,
		IssuedAt: jwt.NewNumericDate(time.Now()),
		Expiry:   jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(1))),
		Issuer:   clientID,
		Audience: orgURL + "/oauth2/v1/token",
		ID:       uuid.New().String(),
	}
	jwtBuilder := jwt.Signed(privateKeySinger).Claims(claims)
	return jwtBuilder.Serialize()
}

// getAccessTokenForPrivateKey requests an access token with the client
// credentials grant. If the org requires DPoP-bound tokens the request is
// repeated with a DPoP proof, and again with the nonce the server asks for.
// Every attempt is a new request with a new DPoP proof and, when a signer is
// available, a new client assertion, because Okta rejects a jti it has already
// seen. Rate limited (429) attempts are retried up to maxRetries times, waiting
// for the rate limit reset but no longer than maxBackoff seconds.
func getAccessTokenForPrivateKey(httpClient *http.Client, orgURL, clientAssertion string, scopes []string, maxRetries int32, maxBackoff int64, clientID string, signer jose.Signer) (*RequestAccessToken, string, *rsa.PrivateKey, error) {
	tokenRequestURL := orgURL + "/oauth2/v1/token"
	var (
		dpopKey       *rsa.PrivateKey
		nonce         string
		nonceAttempts int
		retryCount    int32
		err           error
	)
	for attempt := 0; ; attempt++ {
		if attempt > 0 && signer != nil {
			clientAssertion, err = CreateClientAssertion(orgURL, clientID, signer)
			if err != nil {
				return nil, "", nil, err
			}
		}

		query := urlpkg.Values{}
		query.Add("grant_type", "client_credentials")
		query.Add("scope", strings.Join(scopes, " "))
		query.Add("client_assertion_type", "urn:ietf:params:oauth:client-assertion-type:jwt-bearer")
		query.Add("client_assertion", clientAssertion)
		tokenRequest, err := http.NewRequest("POST", tokenRequestURL, strings.NewReader(query.Encode()))
		if err != nil {
			return nil, "", nil, err
		}
		tokenRequest.Header.Add("Accept", "application/json")
		tokenRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		tokenRequest.Header.Add("User-Agent", NewUserAgent(&config{}).String())
		if dpopKey != nil {
			dpopJWT, err := generateDpopJWT(dpopKey, http.MethodPost, tokenRequestURL, nonce, "")
			if err != nil {
				return nil, "", nil, err
			}
			tokenRequest.Header.Set("DPoP", dpopJWT)
		}

		tokenResponse, err := httpClient.Do(tokenRequest)
		if err != nil {
			return nil, "", nil, err
		}
		respBody, err := io.ReadAll(tokenResponse.Body)
		_ = tokenResponse.Body.Close()
		if err != nil {
			return nil, "", nil, err
		}

		if tokenResponse.StatusCode < http.StatusMultipleChoices {
			var accessToken *RequestAccessToken
			if err := json.Unmarshal(respBody, &accessToken); err != nil {
				return nil, "", nil, fmt.Errorf("failed to parse access token response: %w", err)
			}
			return accessToken, nonce, dpopKey, nil
		}

		var tokenErr Error
		_ = json.Unmarshal(respBody, &tokenErr)
		switch {
		case tokenErr.ErrorMessage == "invalid_dpop_proof" && dpopKey == nil:
			// The org issues DPoP-bound tokens; generate the key the token
			// will be bound to and send a DPoP proof with the next attempt.
			dpopKey, err = generatePrivateKey(2048)
			if err != nil {
				return nil, "", nil, err
			}
			continue
		case tokenErr.ErrorMessage == "use_dpop_nonce" && dpopKey != nil && nonceAttempts < maxDpopNonceAttempts:
			nonce = tokenResponse.Header.Get("Dpop-Nonce")
			nonceAttempts++
			continue
		case tokenResponse.StatusCode == http.StatusTooManyRequests && retryCount < maxRetries:
			backoffDuration, err := Get429BackoffTime(tokenResponse)
			if err != nil || maxBackoff < backoffDuration {
				backoffDuration = maxBackoff
			}
			retryCount++
			// lintignore:R018
			time.Sleep(time.Second * time.Duration(backoffDuration))
			continue
		}
		return nil, "", nil, fmt.Errorf("failed to get access token, status %d: %w", tokenResponse.StatusCode, &tokenErr)
	}
}

func generatePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}
	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func generateDpopJWT(privateKey *rsa.PrivateKey, httpMethod, URL, nonce, accessToken string) (string, error) {
	set, err := jwk.New(privateKey.PublicKey)
	if err != nil {
		return "", err
	}
	err = jwk.AssignKeyID(set)
	if err != nil {
		return "", err
	}
	key := jose.SigningKey{Algorithm: jose.RS256, Key: privateKey}
	signerOpts := jose.SignerOptions{}
	signerOpts.WithType("dpop+jwt")
	signerOpts.WithHeader("jwk", set)
	rsaSigner, err := jose.NewSigner(key, &signerOpts)
	if err != nil {
		return "", err
	}
	dpopClaims := DpopClaims{
		ID:         uuid.New().String(),
		HTTPMethod: httpMethod,
		HTTPURI:    URL,
		IssuedAt:   jwt.NewNumericDate(time.Now()),
		Nonce:      nonce,
	}
	if accessToken != "" {
		h := sha256.New()
		h.Write(StringToAsciiBytes(accessToken))
		dpopClaims.AccessToken = base64.RawURLEncoding.EncodeToString(h.Sum(nil))
	}
	jwtBuilder := jwt.Signed(rsaSigner).Claims(dpopClaims)
	return jwtBuilder.Serialize()
}

type DpopClaims struct {
	HTTPMethod  string           `json:"htm,omitempty"`
	HTTPURI     string           `json:"htu,omitempty"`
	IssuedAt    *jwt.NumericDate `json:"iat,omitempty"`
	Nonce       string           `json:"nonce,omitempty"`
	ID          string           `json:"jti,omitempty"`
	AccessToken string           `json:"ath,omitempty"`
}

func StringToAsciiBytes(s string) []byte {
	t := make([]byte, utf8.RuneCountInString(s))
	i := 0
	for _, r := range s {
		t[i] = byte(r)
		i++
	}
	return t
}

func NewRequestExecutor(httpClient *http.Client, cache cache.Cache, config *config) *RequestExecutor {
	re := RequestExecutor{
		tokenCache: goCache.New(5*time.Minute, 10*time.Minute),
		tokenLock:  &sync.Mutex{},
	}

	re.httpClient = httpClient
	re.config = config
	re.cache = cache
	re.binary = false
	re.headerAccept = "application/json"
	re.headerContentType = "application/json"

	if httpClient == nil {
		tr := &http.Transport{
			IdleConnTimeout: 30 * time.Second,
		}
		re.httpClient = &http.Client{
			Transport: tr,
			Timeout:   time.Second * time.Duration(re.config.Okta.Client.ConnectionTimeout),
		}
	}

	return &re
}

func (re *RequestExecutor) NewRequest(method, url string, body interface{}) (*http.Request, error) {
	var buff io.ReadWriter
	if body != nil {
		switch v := body.(type) {
		case []byte:
			buff = bytes.NewBuffer(v)
		case *bytes.Buffer:
			buff = v
		default:
			buff = new(bytes.Buffer)
			encoder := json.NewEncoder(buff)
			encoder.SetEscapeHTML(false)
			err := encoder.Encode(body)
			if err != nil {
				return nil, err
			}
		}
	}
	url = re.config.Okta.Client.OrgUrl + url
	URLBuilder, err := urlpkg.Parse("/api/v1")
	if err != nil {
		return nil, err
	}
	URL, err := urlpkg.Parse(url)
	if err != nil {
		return nil, err
	}
	if URL.Scheme != "" {
		URLBuilder.Scheme = URL.Scheme
	}
	if URL.Host != "" {
		URLBuilder.Host = URL.Host
	}
	if URL.Path != "" {
		URLBuilder.Path = URL.Path
	}
	urlWithoutQuery := *URLBuilder

	if URL.RawQuery != "" {
		URLBuilder.RawQuery = URL.RawQuery
	}
	req, err := http.NewRequest(method, URLBuilder.String(), buff)
	if err != nil {
		return nil, err
	}

	var auth Authorization

	switch re.config.Okta.Client.AuthorizationMode {
	case "SSWS":
		auth = NewSSWSAuth(re.config.Okta.Client.Token, req)
	case "Bearer":
		auth = NewBearerAuth(re.config.Okta.Client.Token, req)
	case "PrivateKey":
		auth = NewPrivateKeyAuth(PrivateKeyAuthConfig{
			TokenCache:       re.tokenCache,
			TokenLock:        re.tokenLock,
			HttpClient:       re.httpClient,
			PrivateKeySigner: re.config.PrivateKeySigner,
			PrivateKey:       re.config.Okta.Client.PrivateKey,
			PrivateKeyId:     re.config.Okta.Client.PrivateKeyId,
			ClientId:         re.config.Okta.Client.ClientId,
			OrgURL:           re.config.Okta.Client.OrgUrl,
			Scopes:           re.config.Okta.Client.Scopes,
			MaxRetries:       re.config.Okta.Client.RateLimit.MaxRetries,
			MaxBackoff:       re.config.Okta.Client.RateLimit.MaxBackoff,
			Req:              req,
		})
	case "JWT":
		auth = NewJWTAuth(JWTAuthConfig{
			TokenCache:      re.tokenCache,
			TokenLock:       re.tokenLock,
			HttpClient:      re.httpClient,
			OrgURL:          re.config.Okta.Client.OrgUrl,
			Scopes:          re.config.Okta.Client.Scopes,
			ClientAssertion: re.config.Okta.Client.ClientAssertion,
			MaxRetries:      re.config.Okta.Client.RateLimit.MaxRetries,
			MaxBackoff:      re.config.Okta.Client.RateLimit.MaxBackoff,
			Req:             req,
		})
	default:
		return nil, fmt.Errorf("unknown authorization mode %v", re.config.Okta.Client.AuthorizationMode)
	}

	err = auth.Authorize(method, urlWithoutQuery.String())
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", NewUserAgent(re.config).String())
	req.Header.Add("Accept", re.headerAccept)

	if body != nil {
		req.Header.Set("Content-Type", re.headerContentType)
	}

	// Force reset defaults
	re.binary = false
	re.headerAccept = "application/json"
	re.headerContentType = "application/json"
	return req, nil
}

func (re *RequestExecutor) SetHTTPTransport(transport http.RoundTripper) {
	re.httpClient.Transport = transport
}

func (re *RequestExecutor) AsBinary() *RequestExecutor {
	re.binary = true
	return re
}

func (re *RequestExecutor) WithAccept(acceptHeader string) *RequestExecutor {
	re.headerAccept = acceptHeader
	return re
}

func (re *RequestExecutor) WithContentType(contentTypeHeader string) *RequestExecutor {
	re.headerContentType = contentTypeHeader
	return re
}

func (re *RequestExecutor) RefreshNext() *RequestExecutor {
	re.freshCache = true
	return re
}

func (re *RequestExecutor) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	cacheKey := cache.CreateCacheKey(req)
	if req.Method != http.MethodGet {
		re.cache.Delete(cacheKey)
	}
	inCache := re.cache.Has(cacheKey)
	if re.freshCache {
		re.cache.Delete(cacheKey)
		inCache = false
		re.freshCache = false
	}
	if !inCache {
		resp, err := re.doWithRetries(ctx, req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode >= 200 && resp.StatusCode <= 299 && req.Method == http.MethodGet && v != nil && reflect.TypeOf(v).Kind() != reflect.Slice {
			re.cache.Set(cacheKey, resp)
		}
		return buildResponse(resp, re, &v)
	}
	resp := re.cache.Get(cacheKey)
	return buildResponse(resp, re, &v)
}

type oktaBackoff struct {
	retryCount, maxRetries int32
	backoffDuration        time.Duration
	ctx                    context.Context
}

// NextBackOff returns the duration to wait before retrying the operation,
// or backoff. Stop to indicate that no more retries should be made.
func (o *oktaBackoff) NextBackOff() time.Duration {
	// stop retrying if operation reached retry limit
	if o.retryCount > o.maxRetries {
		return backoff.Stop
	}
	return o.backoffDuration
}

// Reset to initial state.
func (o *oktaBackoff) Reset() {}

func (o *oktaBackoff) Context() context.Context {
	return o.ctx
}

func (re *RequestExecutor) doWithRetries(ctx context.Context, req *http.Request) (*http.Response, error) {
	var bodyReader func() io.ReadCloser
	if req.Body != nil {
		buf, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		bodyReader = func() io.ReadCloser {
			return io.NopCloser(bytes.NewReader(buf))
		}
	}
	var (
		resp *http.Response
		err  error
	)
	if re.config.Okta.Client.RequestTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Second*time.Duration(re.config.Okta.Client.RequestTimeout))
		defer cancel()
	}
	bOff := &oktaBackoff{
		ctx:        ctx,
		maxRetries: re.config.Okta.Client.RateLimit.MaxRetries,
	}
	attempt := 0
	operation := func() error {
		attempt++
		// Always rewind the request body when non-nil.
		if bodyReader != nil {
			req.Body = bodyReader()
		}

		// Re-authorize every re-sent request (after a 429 or a network error)
		// so it never carries a DPoP proof JWT that Okta has already seen. The
		// cached access token is kept: it is shared by all requests and a
		// rate limited request does not invalidate it.
		if attempt > 1 && (re.config.Okta.Client.AuthorizationMode == "PrivateKey" || re.config.Okta.Client.AuthorizationMode == "JWT") {
			urlPath := req.URL.RequestURI()

			auth, err := re.NewRequest(req.Method, urlPath, nil)
			if err != nil {
				return err
			}

			req.Header = req.Header.Clone() // Start with original headers

			// Update only the authentication headers from the fresh auth request
			req.Header.Set("Authorization", auth.Header.Get("Authorization"))
			if dpopHeader := auth.Header.Get("Dpop"); dpopHeader != "" {
				req.Header.Set("Dpop", dpopHeader)
			}
			if userAgentExt := auth.Header.Get("x-okta-user-agent-extended"); userAgentExt != "" {
				req.Header.Set("x-okta-user-agent-extended", userAgentExt)
			}

			if bodyReader != nil {
				req.Body = bodyReader()
			}
		} else {
			// Reuse the existing request headers and body
			req.Header = req.Header.Clone()
			if bodyReader != nil {
				req.Body = bodyReader()
			}
		}
		resp, err = re.httpClient.Do(req.WithContext(ctx))
		if errors.Is(err, io.EOF) {
			// retry on EOF errors, which might be caused by network connectivity issues
			return fmt.Errorf("network error: %w", err)
		} else if err != nil {
			// this is error is considered to be permanent and should not be retried
			return backoff.Permanent(err)
		}
		if !tooManyRequests(resp) {
			return nil
		}
		if err = tryDrainBody(resp.Body); err != nil {
			return err
		}
		backoffDuration, err := Get429BackoffTime(resp)
		if err != nil {
			return err
		}
		if re.config.Okta.Client.RateLimit.MaxBackoff < backoffDuration {
			backoffDuration = re.config.Okta.Client.RateLimit.MaxBackoff
		}
		bOff.backoffDuration = time.Second * time.Duration(backoffDuration)
		bOff.retryCount++
		req.Header.Add("X-Okta-Retry-For", resp.Header.Get("X-Okta-Request-Id"))
		req.Header.Add("X-Okta-Retry-Count", fmt.Sprint(bOff.retryCount))
		return errors.New("too many requests")
	}
	err = backoff.Retry(operation, bOff)
	return resp, err
}

func tooManyRequests(resp *http.Response) bool {
	return resp != nil && resp.StatusCode == http.StatusTooManyRequests
}

func tryDrainBody(body io.ReadCloser) error {
	defer body.Close()
	_, err := io.Copy(io.Discard, io.LimitReader(body, 4096))
	return err
}

func Get429BackoffTime(resp *http.Response) (int64, error) {
	requestDate, err := time.Parse("Mon, 02 Jan 2006 15:04:05 GMT", resp.Header.Get("Date"))
	if err != nil {
		// this is error is considered to be permanent and should not be retried
		return 0, backoff.Permanent(fmt.Errorf("date header is missing or invalid: %w", err))
	}
	rateLimitReset, err := strconv.Atoi(resp.Header.Get("X-Rate-Limit-Reset"))
	if err != nil {
		// this is error is considered to be permanent and should not be retried
		return 0, backoff.Permanent(fmt.Errorf("X-Rate-Limit-Reset header is missing or invalid: %w", err))
	}
	return int64(rateLimitReset) - requestDate.Unix() + 1, nil
}

type Response struct {
	*http.Response
	re       *RequestExecutor
	Self     string
	NextPage string
}

func (r *Response) Next(ctx context.Context, v interface{}) (*Response, error) {
	if r.re == nil {
		return nil, errors.New("no initial response provided from previous request")
	}
	req, err := r.re.NewRequest("GET", r.NextPage, nil)
	if err != nil {
		return nil, err
	}
	return r.re.Do(ctx, req, v)
}

func (r *Response) HasNextPage() bool {
	return r.NextPage != ""
}

func newResponse(r *http.Response, re *RequestExecutor) *Response {
	response := &Response{Response: r, re: re}
	links := r.Header["Link"]

	if len(links) > 0 {
		for _, link := range links {
			splitLinkHeader := strings.Split(link, ";")
			if len(splitLinkHeader) < 2 {
				continue
			}
			rawLink := strings.TrimRight(strings.TrimLeft(splitLinkHeader[0], "<"), ">")
			rawURL, _ := urlpkg.Parse(rawLink)
			rawURL.Scheme = ""
			rawURL.Host = ""
			if r.Request != nil {
				q := r.Request.URL.Query()
				for k, v := range rawURL.Query() {
					q.Set(k, v[0])
				}
				rawURL.RawQuery = q.Encode()
			}
			if strings.Contains(link, `rel="self"`) {
				response.Self = rawURL.String()
			}
			if strings.Contains(link, `rel="next"`) {
				response.NextPage = rawURL.String()
			}
		}
	}

	return response
}

func CheckResponseForError(resp *http.Response) error {
	statusCode := resp.StatusCode
	if statusCode >= http.StatusOK && statusCode < http.StatusBadRequest {
		return nil
	}
	e := Error{}
	if (statusCode == http.StatusUnauthorized || statusCode == http.StatusForbidden) &&
		strings.Contains(resp.Header.Get("Www-Authenticate"), "Bearer") {
		for _, v := range strings.Split(resp.Header.Get("Www-Authenticate"), ", ") {
			if strings.Contains(v, "error_description") {
				_, err := toml.Decode(v, &e)
				if err != nil {
					e.ErrorSummary = "unauthorized"
				}
				return &e
			}
		}
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	copyBodyBytes := make([]byte, len(bodyBytes))
	copy(copyBodyBytes, bodyBytes)
	_ = resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	_ = json.NewDecoder(bytes.NewReader(copyBodyBytes)).Decode(&e)
	if statusCode == http.StatusInternalServerError {
		e.ErrorSummary += fmt.Sprintf(", x-okta-request-id=%s", resp.Header.Get("x-okta-request-id"))
	}
	return &e
}

func buildResponse(resp *http.Response, re *RequestExecutor, v interface{}) (*Response, error) {
	ct := resp.Header.Get("Content-Type")
	response := newResponse(resp, re)
	err := CheckResponseForError(resp)
	if err != nil {
		return response, err
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	copyBodyBytes := make([]byte, len(bodyBytes))
	copy(copyBodyBytes, bodyBytes)
	_ = resp.Body.Close()                                // close it to avoid memory leaks
	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // restore the original response body
	if len(copyBodyBytes) == 0 {
		return response, nil
	}
	switch {
	case strings.Contains(ct, "application/xml"):
		err = xml.NewDecoder(bytes.NewReader(copyBodyBytes)).Decode(v)
	case strings.Contains(ct, "application/json"):
		err = json.NewDecoder(bytes.NewReader(copyBodyBytes)).Decode(v)
	case strings.Contains(ct, "application/octet-stream"):
		// since the response is arbitrary binary data, we leave it to the user to decode it
		return response, nil
	default:
		return nil, errors.New("could not build a response for type: " + ct)
	}
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		return nil, err
	}
	return response, nil
}

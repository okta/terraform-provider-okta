package okta

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v4/okta"
	v5okta "github.com/okta/okta-sdk-golang/v5/okta"
	"github.com/okta/terraform-provider-okta/okta/internal/apimutex"
	"github.com/okta/terraform-provider-okta/okta/internal/transport"
	"github.com/okta/terraform-provider-okta/sdk"
)

const (
	OktaTerraformProviderVersion   = "4.18.0"
	OktaTerraformProviderUserAgent = "okta-terraform/" + OktaTerraformProviderVersion
)

var (
	// NOTE: Minor hack where runtime needs to know about testing environment.
	// Global clients are convenience for testing only
	sdkV5Client         *v5okta.APIClient
	sdkV3Client         *okta.APIClient
	sdkV2Client         *sdk.Client
	sdkSupplementClient *sdk.APISupplement
)

type (
	// Config contains our provider schema values and Okta clients
	Config struct {
		orgName                 string
		domain                  string
		httpProxy               string
		accessToken             string
		apiToken                string
		clientID                string
		privateKey              string
		privateKeyId            string
		scopes                  []string
		retryCount              int
		parallelism             int
		backoff                 bool
		minWait                 int
		maxWait                 int
		logLevel                int
		requestTimeout          int
		maxAPICapacity          int
		oktaSDKClientV2         *sdk.Client
		oktaSDKClientV3         *okta.APIClient
		oktaSDKClientV5         *v5okta.APIClient
		oktaSDKsupplementClient *sdk.APISupplement
		logger                  hclog.Logger
		queriedWellKnown        bool
		classicOrg              bool
		timeOperations          TimeOperations
	}
)

func NewConfig(d *schema.ResourceData) *Config {
	// defaults
	config := Config{
		backoff:        true,
		minWait:        30,
		maxWait:        300,
		retryCount:     5,
		parallelism:    1,
		logLevel:       int(hclog.Error),
		requestTimeout: 0,
		maxAPICapacity: 100,
	}
	logLevel := hclog.Level(config.logLevel)
	if os.Getenv("TF_LOG") != "" {
		logLevel = hclog.LevelFromString(os.Getenv("TF_LOG"))
	}
	config.logger = hclog.New(&hclog.LoggerOptions{
		Level:      logLevel,
		TimeFormat: "2006/01/02 03:04:05",
	})

	if val, ok := d.GetOk("org_name"); ok {
		config.orgName = val.(string)
	}
	if config.orgName == "" && os.Getenv("OKTA_ORG_NAME") != "" {
		config.orgName = os.Getenv("OKTA_ORG_NAME")
	}

	if val, ok := d.GetOk("base_url"); ok {
		config.domain = val.(string)
	}
	if config.domain == "" {
		if os.Getenv("OKTA_BASE_URL") != "" {
			config.domain = os.Getenv("OKTA_BASE_URL")
		}
	}

	if val, ok := d.GetOk("api_token"); ok {
		config.apiToken = val.(string)
	}
	if config.apiToken == "" && os.Getenv("OKTA_API_TOKEN") != "" {
		config.apiToken = os.Getenv("OKTA_API_TOKEN")
	}

	if val, ok := d.GetOk("access_token"); ok {
		config.accessToken = val.(string)
	}
	if config.accessToken == "" && os.Getenv("OKTA_ACCESS_TOKEN") != "" {
		config.accessToken = os.Getenv("OKTA_ACCESS_TOKEN")
	}

	if val, ok := d.GetOk("client_id"); ok {
		config.clientID = val.(string)
	}
	if config.clientID == "" && os.Getenv("OKTA_API_CLIENT_ID") != "" {
		config.clientID = os.Getenv("OKTA_API_CLIENT_ID")
	}

	if val, ok := d.GetOk("private_key"); ok {
		config.privateKey = val.(string)
	}
	if config.privateKey == "" && os.Getenv("OKTA_API_PRIVATE_KEY") != "" {
		config.privateKey = os.Getenv("OKTA_API_PRIVATE_KEY")
	}

	if val, ok := d.GetOk("private_key_id"); ok {
		config.privateKeyId = val.(string)
	}
	if config.privateKeyId == "" && os.Getenv("OKTA_API_PRIVATE_KEY_ID") != "" {
		config.privateKeyId = os.Getenv("OKTA_API_PRIVATE_KEY_ID")
	}

	if val, ok := d.GetOk("scopes"); ok {
		config.scopes = convertInterfaceToStringSet(val)
	}
	if v := os.Getenv("OKTA_API_SCOPES"); v != "" && len(config.scopes) == 0 {
		config.scopes = strings.Split(v, ",")
	}

	if val, ok := d.GetOk("max_retries"); ok {
		config.retryCount = val.(int)
	}

	if val, ok := d.GetOk("parallelism"); ok {
		config.parallelism = val.(int)
	}

	if val, ok := d.GetOk("backoff"); ok {
		config.backoff = val.(bool)
	}

	if val, ok := d.GetOk("min_wait_seconds"); ok {
		config.minWait = val.(int)
	}

	if val, ok := d.GetOk("max_wait_seconds"); ok {
		config.maxWait = val.(int)
	}

	if val, ok := d.GetOk("log_level"); ok {
		config.logLevel = val.(int)
	}

	if val, ok := d.GetOk("request_timeout"); ok {
		config.requestTimeout = val.(int)
	}

	if val, ok := d.GetOk("max_api_capacity"); ok {
		config.maxAPICapacity = val.(int)
	}
	if config.maxAPICapacity == 0 {
		if os.Getenv("MAX_API_CAPACITY") != "" {
			mac, err := strconv.ParseInt(os.Getenv("MAX_API_CAPACITY"), 10, 64)
			if err != nil {
				config.logger.Error("error with max_api_capacity value", err)
			} else {
				config.maxAPICapacity = int(mac)
			}
		}
	}

	if httpProxy, ok := d.Get("http_proxy").(string); ok {
		config.httpProxy = httpProxy
	}
	if config.httpProxy == "" && os.Getenv("OKTA_HTTP_PROXY") != "" {
		config.httpProxy = os.Getenv("OKTA_HTTP_PROXY")
	}

	if v := os.Getenv("OKTA_API_SCOPES"); v != "" && len(config.scopes) == 0 {
		config.scopes = strings.Split(v, ",")
	}

	return &config
}

// IsClassicOrg returns true if the org is a classic org. Does lazy evaluation
// of the well known endpoint.
func (c *Config) IsClassicOrg(ctx context.Context) bool {
	if !c.queriedWellKnown {
		// Discover if the Okta Org is Classic or OIE
		org, _, err := c.oktaSDKClientV3.OrgSettingAPI.GetWellknownOrgMetadata(ctx).Execute()
		if err != nil {
			c.logger.Error("error querying GET /.well-known/okta-organization", "error", err)
			return c.classicOrg
		}

		c.classicOrg = (org.GetPipeline() == "v1") // v1 == Classic, idx == OIE
		c.queriedWellKnown = true
	}

	return c.classicOrg
}

func (c *Config) IsOAuth20Auth() bool {
	return c.privateKey != "" || c.accessToken != ""
}

func (c *Config) SetTimeOperations(op TimeOperations) {
	c.timeOperations = op
}

func (c *Config) resetHttpTransport(transport *http.RoundTripper) {
	c.oktaSDKClientV5.GetConfig().HTTPClient.Transport = *transport
	c.oktaSDKClientV3.GetConfig().HTTPClient.Transport = *transport
	c.oktaSDKClientV2.GetConfig().HttpClient.Transport = *transport

	re := c.oktaSDKClientV2.CloneRequestExecutor()
	re.SetHTTPTransport(c.oktaSDKClientV3.GetConfig().HTTPClient.Transport)
	c.oktaSDKsupplementClient = &sdk.APISupplement{
		RequestExecutor: re,
	}
	// NOTE: global clients are convenience for testing only
	sdkSupplementClient = c.oktaSDKsupplementClient
}

// loadClients initializes the Okta SDK clients
func (c *Config) loadClients(ctx context.Context) error {
	v3Client, err := oktaV3SDKClient(c)
	if err != nil {
		return err
	}
	c.oktaSDKClientV3 = v3Client

	v5Client, err := oktaV5SDKClient(c)
	if err != nil {
		return err
	}
	c.oktaSDKClientV5 = v5Client

	// TODO: remove sdk client when v3 client is fully utilized within the provider
	client, err := oktaSDKClient(c)
	if err != nil {
		return err
	}
	c.oktaSDKClientV2 = client

	// TODO: remove supplement client when v3 client is fully utilized within the provider
	re := client.CloneRequestExecutor()
	re.SetHTTPTransport(c.oktaSDKClientV3.GetConfig().HTTPClient.Transport)
	c.oktaSDKsupplementClient = &sdk.APISupplement{
		RequestExecutor: re,
	}

	// NOTE: global clients are convenience for testing only; however do not
	// remove this code
	sdkV5Client = c.oktaSDKClientV5
	sdkV3Client = c.oktaSDKClientV3
	sdkV2Client = c.oktaSDKClientV2
	sdkSupplementClient = c.oktaSDKsupplementClient

	return nil
}

func (c *Config) verifyCredentials(ctx context.Context) error {
	// NOTE: validate credentials during initial config with a call to
	// GET /api/v1/users/me
	// only for SSWS API token. Should we keep doing this?
	if c.apiToken != "" {
		if _, _, err := c.oktaSDKClientV3.UserAPI.GetUser(ctx, "me").Execute(); err != nil {
			return fmt.Errorf("error with v3 SDK client: %v", err)
		}
		if _, _, err := c.oktaSDKClientV2.User.GetUser(ctx, "me"); err != nil {
			return fmt.Errorf("error with v2 SDK client: %v", err)
		}
	}

	return nil
}

func (c *Config) handleFrameworkDefaults(ctx context.Context, data *FrameworkProviderData) error {
	var err error
	if data.OrgName.IsNull() && os.Getenv("OKTA_ORG_NAME") != "" {
		data.OrgName = types.StringValue(os.Getenv("OKTA_ORG_NAME"))
	}
	if data.AccessToken.IsNull() && os.Getenv("OKTA_ACCESS_TOKEN") != "" {
		data.AccessToken = types.StringValue(os.Getenv("OKTA_ACCESS_TOKEN"))
	}
	if data.APIToken.IsNull() && os.Getenv("OKTA_API_TOKEN") != "" {
		data.APIToken = types.StringValue(os.Getenv("OKTA_API_TOKEN"))
	}
	if data.ClientID.IsNull() && os.Getenv("OKTA_API_CLIENT_ID") != "" {
		data.ClientID = types.StringValue(os.Getenv("OKTA_API_CLIENT_ID"))
	}
	if data.Scopes.IsNull() && os.Getenv("OKTA_API_SCOPES") != "" {
		v := os.Getenv("OKTA_API_SCOPES")
		scopes := strings.Split(v, ",")
		if len(scopes) > 0 {
			scopesTF := make([]attr.Value, 0)
			for _, scope := range scopes {
				scopesTF = append(scopesTF, types.StringValue(scope))
			}
			data.Scopes, _ = types.SetValue(types.StringType, scopesTF)
		}
	}
	if data.PrivateKey.IsNull() && os.Getenv("OKTA_API_PRIVATE_KEY") != "" {
		data.PrivateKey = types.StringValue(os.Getenv("OKTA_API_PRIVATE_KEY"))
	}
	if data.PrivateKeyID.IsNull() && os.Getenv("OKTA_API_PRIVATE_KEY_ID") != "" {
		data.PrivateKeyID = types.StringValue(os.Getenv("OKTA_API_PRIVATE_KEY_ID"))
	}
	if data.BaseURL.IsNull() {
		if os.Getenv("OKTA_BASE_URL") != "" {
			data.BaseURL = types.StringValue(os.Getenv("OKTA_BASE_URL"))
		}
	}
	if data.HTTPProxy.IsNull() && os.Getenv("OKTA_HTTP_PROXY") != "" {
		data.HTTPProxy = types.StringValue(os.Getenv("OKTA_HTTP_PROXY"))
	}
	if data.MaxAPICapacity.IsNull() {
		if os.Getenv("MAX_API_CAPACITY") != "" {
			mac, err := strconv.ParseInt(os.Getenv("MAX_API_CAPACITY"), 10, 64)
			if err != nil {
				return err
			}
			data.MaxAPICapacity = types.Int64Value(mac)
		} else {
			data.MaxAPICapacity = types.Int64Value(100)
		}
	}
	data.Backoff = types.BoolValue(true)
	data.MinWaitSeconds = types.Int64Value(30)
	data.MaxWaitSeconds = types.Int64Value(300)
	data.MaxRetries = types.Int64Value(5)
	data.Parallelism = types.Int64Value(1)
	data.LogLevel = types.Int64Value(int64(hclog.Error))
	data.RequestTimeout = types.Int64Value(0)

	if os.Getenv("TF_LOG") != "" {
		data.LogLevel = types.Int64Value(int64(hclog.LevelFromString(os.Getenv("TF_LOG"))))
	}
	c.logger = hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Level(data.LogLevel.ValueInt64()),
		TimeFormat: "2006/01/02 03:04:05",
	})

	return err
}

// oktaSDKClient should be called with a primary http client that is utilized
// throughout the provider
func oktaSDKClient(c *Config) (client *sdk.Client, err error) {
	httpClient := c.oktaSDKClientV3.GetConfig().HTTPClient
	var orgUrl string
	var disableHTTPS bool
	if c.httpProxy != "" {
		orgUrl = strings.TrimSuffix(c.httpProxy, "/")
		disableHTTPS = strings.HasPrefix(orgUrl, "http://")
	} else {
		orgUrl = fmt.Sprintf("https://%v.%v", c.orgName, c.domain)
	}
	_, err = url.Parse(orgUrl)
	if err != nil {
		return nil, fmt.Errorf("malformed Okta API URL (org_name+base_url value, or http_proxy value): %+v", err)
	}

	setters := []sdk.ConfigSetter{
		sdk.WithOrgUrl(orgUrl),
		sdk.WithCache(false),
		sdk.WithHttpClientPtr(httpClient),
		sdk.WithRateLimitMaxBackOff(int64(c.maxWait)),
		sdk.WithRequestTimeout(int64(c.requestTimeout)),
		sdk.WithRateLimitMaxRetries(int32(c.retryCount)),
		sdk.WithUserAgentExtra(OktaTerraformProviderUserAgent),
	}

	switch {
	case c.accessToken != "":
		setters = append(
			setters,
			sdk.WithToken(c.accessToken), sdk.WithAuthorizationMode("Bearer"),
		)

	case c.apiToken != "":
		setters = append(
			setters,
			sdk.WithToken(c.apiToken), sdk.WithAuthorizationMode("SSWS"),
		)

	case c.privateKey != "":
		setters = append(
			setters,
			sdk.WithPrivateKey(c.privateKey), sdk.WithPrivateKeyId(c.privateKeyId), sdk.WithScopes(c.scopes), sdk.WithClientId(c.clientID), sdk.WithAuthorizationMode("PrivateKey"),
		)
	}

	if disableHTTPS {
		setters = append(setters, sdk.WithTestingDisableHttpsCheck(true))
	}

	_, client, err = sdk.NewClient(
		context.Background(),
		setters...,
	)
	return
}

func oktaV3SDKClient(c *Config) (client *okta.APIClient, err error) {
	var httpClient *http.Client
	logLevel := strings.ToLower(os.Getenv("TF_LOG"))
	debugHttpRequests := (logLevel == "1" || logLevel == "debug" || logLevel == "trace")
	if c.backoff {
		retryableClient := retryablehttp.NewClient()
		retryableClient.RetryWaitMin = time.Second * time.Duration(c.minWait)
		retryableClient.RetryWaitMax = time.Second * time.Duration(c.maxWait)
		retryableClient.RetryMax = c.retryCount
		retryableClient.Logger = c.logger
		if debugHttpRequests {
			// Needed for pretty printing http protocol in a local developer environment, ignore deprecation warnings.
			//lint:ignore SA1019 used in developer mode only
			retryableClient.HTTPClient.Transport = logging.NewTransport("Okta", retryableClient.HTTPClient.Transport)
		} else {
			retryableClient.HTTPClient.Transport = logging.NewSubsystemLoggingHTTPTransport("Okta", retryableClient.HTTPClient.Transport)
		}
		retryableClient.ErrorHandler = errHandler
		retryableClient.CheckRetry = checkRetry
		httpClient = retryableClient.StandardClient()
		c.logger.Info(fmt.Sprintf("running with backoff http client, wait min %d, wait max %d, retry max %d", retryableClient.RetryWaitMin, retryableClient.RetryWaitMax, retryableClient.RetryMax))
	} else {
		httpClient = cleanhttp.DefaultClient()
		if debugHttpRequests {
			// Needed for pretty printing http protocol in a local developer environment, ignore deprecation warnings.
			//lint:ignore SA1019 used in developer mode only
			httpClient.Transport = logging.NewTransport("Okta", httpClient.Transport)
		} else {
			httpClient.Transport = logging.NewSubsystemLoggingHTTPTransport("Okta", httpClient.Transport)
		}
		c.logger.Info("running with default http client")
	}

	// adds transport governor to retryable or default client
	if c.maxAPICapacity > 0 && c.maxAPICapacity < 100 {
		c.logger.Info(fmt.Sprintf("running with experimental max_api_capacity configuration at %d%%", c.maxAPICapacity))
		apiMutex, err := apimutex.NewAPIMutex(c.maxAPICapacity)
		if err != nil {
			return nil, err
		}
		httpClient.Transport = transport.NewGovernedTransport(httpClient.Transport, apiMutex, c.logger)
	}
	var orgUrl string
	var disableHTTPS bool
	if c.httpProxy != "" {
		orgUrl = strings.TrimSuffix(c.httpProxy, "/")
		disableHTTPS = strings.HasPrefix(orgUrl, "http://")
	} else {
		orgUrl = fmt.Sprintf("https://%v.%v", c.orgName, c.domain)
	}
	_, err = url.Parse(orgUrl)
	if err != nil {
		return nil, fmt.Errorf("malformed Okta API URL (org_name+base_url value, or http_proxy value): %+v", err)
	}

	setters := []okta.ConfigSetter{
		okta.WithOrgUrl(orgUrl),
		okta.WithCache(false),
		okta.WithHttpClientPtr(httpClient),
		okta.WithRateLimitMaxBackOff(int64(c.maxWait)),
		okta.WithRequestTimeout(int64(c.requestTimeout)),
		okta.WithRateLimitMaxRetries(int32(c.retryCount)),
		okta.WithUserAgentExtra(OktaTerraformProviderUserAgent),
	}
	// v3 client also needs http proxy explicitly set
	if c.httpProxy != "" {
		_url, err := url.Parse(c.httpProxy)
		if err != nil {
			return nil, err
		}
		host := okta.WithProxyHost(_url.Hostname())
		setters = append(setters, host)

		sPort := _url.Port()
		if sPort == "" {
			sPort = "80"
		}
		iPort, err := strconv.Atoi(sPort)
		if err != nil {
			return nil, err
		}
		port := okta.WithProxyPort(int32(iPort))
		setters = append(setters, port)
	}

	switch {
	case c.accessToken != "":
		setters = append(
			setters,
			okta.WithToken(c.accessToken), okta.WithAuthorizationMode("Bearer"),
		)

	case c.apiToken != "":
		setters = append(
			setters,
			okta.WithToken(c.apiToken), okta.WithAuthorizationMode("SSWS"),
		)

	case c.privateKey != "":
		setters = append(
			setters,
			okta.WithPrivateKey(c.privateKey), okta.WithPrivateKeyId(c.privateKeyId), okta.WithScopes(c.scopes), okta.WithClientId(c.clientID), okta.WithAuthorizationMode("PrivateKey"),
		)
	}

	if disableHTTPS {
		setters = append(setters, okta.WithTestingDisableHttpsCheck(true))
	}

	config, err := okta.NewConfiguration(setters...)
	if err != nil {
		return nil, err
	}
	client = okta.NewAPIClient(config)
	return client, nil
}

func oktaV5SDKClient(c *Config) (client *v5okta.APIClient, err error) {
	var httpClient *http.Client
	logLevel := strings.ToLower(os.Getenv("TF_LOG"))
	debugHttpRequests := (logLevel == "1" || logLevel == "debug" || logLevel == "trace")
	if c.backoff {
		retryableClient := retryablehttp.NewClient()
		retryableClient.RetryWaitMin = time.Second * time.Duration(c.minWait)
		retryableClient.RetryWaitMax = time.Second * time.Duration(c.maxWait)
		retryableClient.RetryMax = c.retryCount
		retryableClient.Logger = c.logger
		if debugHttpRequests {
			// Needed for pretty printing http protocol in a local developer environment, ignore deprecation warnings.
			//lint:ignore SA1019 used in developer mode only
			retryableClient.HTTPClient.Transport = logging.NewTransport("Okta", retryableClient.HTTPClient.Transport)
		} else {
			retryableClient.HTTPClient.Transport = logging.NewSubsystemLoggingHTTPTransport("Okta", retryableClient.HTTPClient.Transport)
		}
		retryableClient.ErrorHandler = errHandler
		retryableClient.CheckRetry = checkRetry
		httpClient = retryableClient.StandardClient()
		c.logger.Info(fmt.Sprintf("running with backoff http client, wait min %d, wait max %d, retry max %d", retryableClient.RetryWaitMin, retryableClient.RetryWaitMax, retryableClient.RetryMax))
	} else {
		httpClient = cleanhttp.DefaultClient()
		if debugHttpRequests {
			// Needed for pretty printing http protocol in a local developer environment, ignore deprecation warnings.
			//lint:ignore SA1019 used in developer mode only
			httpClient.Transport = logging.NewTransport("Okta", httpClient.Transport)
		} else {
			httpClient.Transport = logging.NewSubsystemLoggingHTTPTransport("Okta", httpClient.Transport)
		}
		c.logger.Info("running with default http client")
	}

	// adds transport governor to retryable or default client
	if c.maxAPICapacity > 0 && c.maxAPICapacity < 100 {
		c.logger.Info(fmt.Sprintf("running with experimental max_api_capacity configuration at %d%%", c.maxAPICapacity))
		apiMutex, err := apimutex.NewAPIMutex(c.maxAPICapacity)
		if err != nil {
			return nil, err
		}
		httpClient.Transport = transport.NewGovernedTransport(httpClient.Transport, apiMutex, c.logger)
	}
	var orgUrl string
	var disableHTTPS bool
	if c.httpProxy != "" {
		orgUrl = strings.TrimSuffix(c.httpProxy, "/")
		disableHTTPS = strings.HasPrefix(orgUrl, "http://")
	} else {
		orgUrl = fmt.Sprintf("https://%v.%v", c.orgName, c.domain)
	}
	_, err = url.Parse(orgUrl)
	if err != nil {
		return nil, fmt.Errorf("malformed Okta API URL (org_name+base_url value, or http_proxy value): %+v", err)
	}

	setters := []v5okta.ConfigSetter{
		v5okta.WithOrgUrl(orgUrl),
		v5okta.WithCache(false),
		v5okta.WithHttpClientPtr(httpClient),
		v5okta.WithRateLimitMaxBackOff(int64(c.maxWait)),
		v5okta.WithRequestTimeout(int64(c.requestTimeout)),
		v5okta.WithRateLimitMaxRetries(int32(c.retryCount)),
		v5okta.WithUserAgentExtra(OktaTerraformProviderUserAgent),
	}
	// v3 client also needs http proxy explicitly set
	if c.httpProxy != "" {
		_url, err := url.Parse(c.httpProxy)
		if err != nil {
			return nil, err
		}
		host := v5okta.WithProxyHost(_url.Hostname())
		setters = append(setters, host)

		sPort := _url.Port()
		if sPort == "" {
			sPort = "80"
		}
		iPort, err := strconv.Atoi(sPort)
		if err != nil {
			return nil, err
		}
		port := v5okta.WithProxyPort(int32(iPort))
		setters = append(setters, port)
	}

	switch {
	case c.accessToken != "":
		setters = append(
			setters,
			v5okta.WithToken(c.accessToken), v5okta.WithAuthorizationMode("Bearer"),
		)

	case c.apiToken != "":
		setters = append(
			setters,
			v5okta.WithToken(c.apiToken), v5okta.WithAuthorizationMode("SSWS"),
		)

	case c.privateKey != "":
		setters = append(
			setters,
			v5okta.WithPrivateKey(c.privateKey), v5okta.WithPrivateKeyId(c.privateKeyId), v5okta.WithScopes(c.scopes), v5okta.WithClientId(c.clientID), v5okta.WithAuthorizationMode("PrivateKey"),
		)
	}

	if disableHTTPS {
		setters = append(setters, v5okta.WithTestingDisableHttpsCheck(true))
	}

	config, err := v5okta.NewConfiguration(setters...)
	if err != nil {
		return nil, err
	}
	client = v5okta.NewAPIClient(config)
	return client, nil
}

func errHandler(resp *http.Response, err error, numTries int) (*http.Response, error) {
	if err != nil {
		return resp, err
	}
	defer resp.Body.Close()
	err = sdk.CheckResponseForError(resp)
	if err != nil {
		var oErr *sdk.Error
		if errors.As(err, &oErr) {
			oErr.ErrorSummary = fmt.Sprintf("%s, giving up after %d attempt(s)", oErr.ErrorSummary, numTries)
			return resp, oErr
		}
		return resp, fmt.Errorf("%v: giving up after %d attempt(s)", err, numTries)
	}
	return resp, nil
}

type contextKey string

const retryOnStatusCodes contextKey = "retryOnStatusCodes"

// Used to make http client retry on provided list of response status codes
//
// To enable this check, inject `retryOnStatusCodes` key into the context with list of status codes you want to retry on
//
//	ctx = context.WithValue(ctx, retryOnStatusCodes, []int{404, 409})
func checkRetry(ctx context.Context, resp *http.Response, err error) (bool, error) {
	// do not retry on context.Canceled or context.DeadlineExceeded
	if ctx.Err() != nil {
		return false, ctx.Err()
	}
	retryCodes, ok := ctx.Value(retryOnStatusCodes).([]int)
	if ok && resp != nil && containsInt(retryCodes, resp.StatusCode) {
		return true, nil
	}
	// don't retry on internal server errors
	if resp != nil && resp.StatusCode == http.StatusInternalServerError {
		return false, nil
	}
	return retryablehttp.DefaultRetryPolicy(ctx, resp, err)
}

type TimeOperations interface {
	DoNotRetry(error) bool
	Sleep(time.Duration)
}

type ProductionTimeOperations struct{}

// DoNotRetry always retry in production
func (o *ProductionTimeOperations) DoNotRetry(err error) bool {
	return false
}

// Sleep facade to actual time.Sleep in production
func (o *ProductionTimeOperations) Sleep(d time.Duration) {
	//lintignore:R018
	time.Sleep(d)
}

// NewProductionTimeOperations new production time operations
func NewProductionTimeOperations() TimeOperations {
	return &ProductionTimeOperations{}
}

type TestTimeOperations struct{}

// DoNotRetry tests do not retry when there is an error and VCR is recording
func (o *TestTimeOperations) DoNotRetry(err error) bool {
	return err != nil && os.Getenv("OKTA_VCR_TF_ACC") == "record"
}

// Sleep no sleeping when test is in VCR play mode
func (o *TestTimeOperations) Sleep(d time.Duration) {
	if os.Getenv("OKTA_VCR_TF_ACC") != "play" {
		//lintignore:R018
		time.Sleep(d)
	}
}

// NewTestTimeOperations new test time operations
func NewTestTimeOperations() TimeOperations {
	return &TestTimeOperations{}
}

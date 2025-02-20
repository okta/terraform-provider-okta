package config

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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v4/okta"
	v5okta "github.com/okta/okta-sdk-golang/v5/okta"
	"github.com/okta/terraform-provider-okta/okta/internal/apimutex"
	"github.com/okta/terraform-provider-okta/okta/internal/transport"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
)

const (
var (
	// FIXME: move this to a package called `version`
	OktaTerraformProviderVersion   = "4.15.0"
	OktaTerraformProviderUserAgent = "okta-terraform/" + OktaTerraformProviderVersion
)

var (
	// NOTE: Minor hack where runtime needs to know about testing environment.
	// Global clients are convenience for testing only
	SdkV5ClientForTest         *v5okta.APIClient
	SdkV3ClientForTest         *okta.APIClient
	SdkV2ClientForTest         *sdk.Client
	SdkSupplementClientForTest *sdk.APISupplement
)

type (
	// Config contains our provider schema values and Okta clients
	Config struct {
		AccessToken             string
		ApiToken                string
		Backoff                 bool
		ClassicOrg              bool
		ClientID                string
		Domain                  string
		HttpProxy               string
		LogLevel                int
		Logger                  hclog.Logger
		MaxAPICapacity          int
		MaxWait                 int
		MinWait                 int
		OktaSDKClientV2         *sdk.Client
		OktaSDKClientV3         *okta.APIClient
		OktaSDKClientV5         *v5okta.APIClient
		OktaSDKsupplementClient *sdk.APISupplement
		OrgName                 string
		Parallelism             int
		PrivateKey              string
		PrivateKeyId            string
		QueriedWellKnown        bool
		RequestTimeout          int
		RetryCount              int
		Scopes                  []string
		TimeOperations          TimeOperations
	}
)

func NewConfig(d *schema.ResourceData) *Config {
	// defaults
	config := Config{
		Backoff:        true,
		LogLevel:       int(hclog.Error),
		MaxAPICapacity: 100,
		MaxWait:        300,
		MinWait:        30,
		Parallelism:    1,
		RequestTimeout: 0,
		RetryCount:     5,
	}
	logLevel := hclog.Level(config.LogLevel)
	if os.Getenv("TF_LOG") != "" {
		logLevel = hclog.LevelFromString(os.Getenv("TF_LOG"))
	}
	config.SetupLogger(logLevel)

	if val, ok := d.GetOk("org_name"); ok {
		config.OrgName = val.(string)
	}
	if config.OrgName == "" && os.Getenv("OKTA_ORG_NAME") != "" {
		config.OrgName = os.Getenv("OKTA_ORG_NAME")
	}

	if val, ok := d.GetOk("base_url"); ok {
		config.Domain = val.(string)
	}
	if config.Domain == "" {
		if os.Getenv("OKTA_BASE_URL") != "" {
			config.Domain = os.Getenv("OKTA_BASE_URL")
		}
	}

	if val, ok := d.GetOk("api_token"); ok {
		config.ApiToken = val.(string)
	}
	if config.ApiToken == "" && os.Getenv("OKTA_API_TOKEN") != "" {
		config.ApiToken = os.Getenv("OKTA_API_TOKEN")
	}

	if val, ok := d.GetOk("access_token"); ok {
		config.AccessToken = val.(string)
	}
	if config.AccessToken == "" && os.Getenv("OKTA_ACCESS_TOKEN") != "" {
		config.AccessToken = os.Getenv("OKTA_ACCESS_TOKEN")
	}

	if val, ok := d.GetOk("client_id"); ok {
		config.ClientID = val.(string)
	}
	if config.ClientID == "" && os.Getenv("OKTA_API_CLIENT_ID") != "" {
		config.ClientID = os.Getenv("OKTA_API_CLIENT_ID")
	}

	if val, ok := d.GetOk("private_key"); ok {
		config.PrivateKey = val.(string)
	}
	if config.PrivateKey == "" && os.Getenv("OKTA_API_PRIVATE_KEY") != "" {
		config.PrivateKey = os.Getenv("OKTA_API_PRIVATE_KEY")
	}

	if val, ok := d.GetOk("private_key_id"); ok {
		config.PrivateKeyId = val.(string)
	}
	if config.PrivateKeyId == "" && os.Getenv("OKTA_API_PRIVATE_KEY_ID") != "" {
		config.PrivateKeyId = os.Getenv("OKTA_API_PRIVATE_KEY_ID")
	}

	if val, ok := d.GetOk("scopes"); ok {
		config.Scopes = utils.ConvertInterfaceToStringSet(val)
	}
	if v := os.Getenv("OKTA_API_SCOPES"); v != "" && len(config.Scopes) == 0 {
		config.Scopes = strings.Split(v, ",")
	}

	if val, ok := d.GetOk("max_retries"); ok {
		config.RetryCount = val.(int)
	}

	if val, ok := d.GetOk("parallelism"); ok {
		config.Parallelism = val.(int)
	}

	if val, ok := d.GetOk("backoff"); ok {
		config.Backoff = val.(bool)
	}

	if val, ok := d.GetOk("min_wait_seconds"); ok {
		config.MinWait = val.(int)
	}

	if val, ok := d.GetOk("max_wait_seconds"); ok {
		config.MaxWait = val.(int)
	}

	if val, ok := d.GetOk("log_level"); ok {
		config.LogLevel = val.(int)
	}

	if val, ok := d.GetOk("request_timeout"); ok {
		config.RequestTimeout = val.(int)
	}

	if val, ok := d.GetOk("max_api_capacity"); ok {
		config.MaxAPICapacity = val.(int)
	}
	if config.MaxAPICapacity == 0 {
		if os.Getenv("MAX_API_CAPACITY") != "" {
			mac, err := strconv.ParseInt(os.Getenv("MAX_API_CAPACITY"), 10, 64)
			if err != nil {
				config.Logger.Error("error with max_api_capacity value", err)
			} else {
				config.MaxAPICapacity = int(mac)
			}
		}
	}

	if httpProxy, ok := d.Get("http_proxy").(string); ok {
		config.HttpProxy = httpProxy
	}
	if config.HttpProxy == "" && os.Getenv("OKTA_HTTP_PROXY") != "" {
		config.HttpProxy = os.Getenv("OKTA_HTTP_PROXY")
	}

	if v := os.Getenv("OKTA_API_SCOPES"); v != "" && len(config.Scopes) == 0 {
		config.Scopes = strings.Split(v, ",")
	}

	return &config
}

func (c *Config) SetupLogger(logLevel hclog.Level) {
	c.Logger = hclog.New(&hclog.LoggerOptions{
		Level:      logLevel,
		TimeFormat: "2006/01/02 03:04:05",
	})
}

// IsClassicOrg returns true if the org is a classic org. Does lazy evaluation
// of the well known endpoint.
func (c *Config) IsClassicOrg(ctx context.Context) bool {
	if !c.QueriedWellKnown {
		// Discover if the Okta Org is Classic or OIE
		org, _, err := c.OktaSDKClientV3.OrgSettingAPI.GetWellknownOrgMetadata(ctx).Execute()
		if err != nil {
			c.Logger.Error("error querying GET /.well-known/okta-organization", "error", err)
			return c.ClassicOrg
		}

		c.ClassicOrg = (org.GetPipeline() == "v1") // v1 == Classic, idx == OIE
		c.QueriedWellKnown = true
	}

	return c.ClassicOrg
}

func (c *Config) IsOAuth20Auth() bool {
	return c.PrivateKey != "" || c.AccessToken != ""
}

func (c *Config) SetTimeOperations(op TimeOperations) {
	c.TimeOperations = op
}

func (c *Config) ResetHttpTransport(transport *http.RoundTripper) {
	c.OktaSDKClientV5.GetConfig().HTTPClient.Transport = *transport
	c.OktaSDKClientV3.GetConfig().HTTPClient.Transport = *transport
	c.OktaSDKClientV2.GetConfig().HttpClient.Transport = *transport

	re := c.OktaSDKClientV2.CloneRequestExecutor()
	re.SetHTTPTransport(c.OktaSDKClientV3.GetConfig().HTTPClient.Transport)
	c.OktaSDKsupplementClient = &sdk.APISupplement{
		RequestExecutor: re,
	}
	// NOTE: global clients are a convenience for testing only
	SdkSupplementClientForTest = c.OktaSDKsupplementClient
}

// LoadClients initializes the Okta SDK clients
func (c *Config) LoadClients() error {
	v3Client, err := oktaV3SDKClient(c)
	if err != nil {
		return err
	}
	c.OktaSDKClientV3 = v3Client

	v5Client, err := oktaV5SDKClient(c)
	if err != nil {
		return err
	}
	c.OktaSDKClientV5 = v5Client

	// TODO: remove sdk client when v3 client is fully utilized within the provider
	client, err := oktaSDKClient(c)
	if err != nil {
		return err
	}
	c.OktaSDKClientV2 = client

	// TODO: remove supplement client when v3 client is fully utilized within the provider
	re := client.CloneRequestExecutor()
	re.SetHTTPTransport(c.OktaSDKClientV3.GetConfig().HTTPClient.Transport)
	c.OktaSDKsupplementClient = &sdk.APISupplement{
		RequestExecutor: re,
	}

	// NOTE: global clients are convenience for testing only; however do not
	// remove this code
	SdkV5ClientForTest = c.OktaSDKClientV5
	SdkV3ClientForTest = c.OktaSDKClientV3
	SdkV2ClientForTest = c.OktaSDKClientV2
	SdkSupplementClientForTest = c.OktaSDKsupplementClient

	return nil
}

func (c *Config) VerifyCredentials(ctx context.Context) error {
	// NOTE: validate credentials during initial config with a call to
	// GET /api/v1/users/me
	// only for SSWS API token. Should we keep doing this?
	if c.ApiToken != "" {
		if _, _, err := c.OktaSDKClientV3.UserAPI.GetUser(ctx, "me").Execute(); err != nil {
			return fmt.Errorf("error with v3 SDK client: %v", err)
		}
		if _, _, err := c.OktaSDKClientV2.User.GetUser(ctx, "me"); err != nil {
			return fmt.Errorf("error with v2 SDK client: %v", err)
		}
	}

	return nil
}

func oktaSDKClient(c *Config) (client *sdk.Client, err error) {
	httpClient := c.OktaSDKClientV3.GetConfig().HTTPClient
	var orgUrl string
	var disableHTTPS bool
	if c.HttpProxy != "" {
		orgUrl = strings.TrimSuffix(c.HttpProxy, "/")
		disableHTTPS = strings.HasPrefix(orgUrl, "http://")
	} else {
		orgUrl = fmt.Sprintf("https://%v.%v", c.OrgName, c.Domain)
	}
	_, err = url.Parse(orgUrl)
	if err != nil {
		return nil, fmt.Errorf("malformed Okta API URL (org_name+base_url value, or http_proxy value): %+v", err)
	}

	setters := []sdk.ConfigSetter{
		sdk.WithOrgUrl(orgUrl),
		sdk.WithCache(false),
		sdk.WithHttpClientPtr(httpClient),
		sdk.WithRateLimitMaxBackOff(int64(c.MaxWait)),
		sdk.WithRequestTimeout(int64(c.RequestTimeout)),
		sdk.WithRateLimitMaxRetries(int32(c.RetryCount)),
		sdk.WithUserAgentExtra(OktaTerraformProviderUserAgent),
	}

	switch {
	case c.AccessToken != "":
		setters = append(
			setters,
			sdk.WithToken(c.AccessToken), sdk.WithAuthorizationMode("Bearer"),
		)

	case c.ApiToken != "":
		setters = append(
			setters,
			sdk.WithToken(c.ApiToken), sdk.WithAuthorizationMode("SSWS"),
		)

	case c.PrivateKey != "":
		setters = append(
			setters,
			sdk.WithPrivateKey(c.PrivateKey), sdk.WithPrivateKeyId(c.PrivateKeyId), sdk.WithScopes(c.Scopes), sdk.WithClientId(c.ClientID), sdk.WithAuthorizationMode("PrivateKey"),
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
	if c.Backoff {
		retryableClient := retryablehttp.NewClient()
		retryableClient.RetryWaitMin = time.Second * time.Duration(c.MinWait)
		retryableClient.RetryWaitMax = time.Second * time.Duration(c.MaxWait)
		retryableClient.RetryMax = c.RetryCount
		retryableClient.Logger = c.Logger
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
		c.Logger.Info(fmt.Sprintf("running with backoff http client, wait min %d, wait max %d, retry max %d", retryableClient.RetryWaitMin, retryableClient.RetryWaitMax, retryableClient.RetryMax))
	} else {
		httpClient = cleanhttp.DefaultClient()
		if debugHttpRequests {
			// Needed for pretty printing http protocol in a local developer environment, ignore deprecation warnings.
			//lint:ignore SA1019 used in developer mode only
			httpClient.Transport = logging.NewTransport("Okta", httpClient.Transport)
		} else {
			httpClient.Transport = logging.NewSubsystemLoggingHTTPTransport("Okta", httpClient.Transport)
		}
		c.Logger.Info("running with default http client")
	}

	// adds transport governor to retryable or default client
	if c.MaxAPICapacity > 0 && c.MaxAPICapacity < 100 {
		c.Logger.Info(fmt.Sprintf("running with experimental max_api_capacity configuration at %d%%", c.MaxAPICapacity))
		apiMutex, err := apimutex.NewAPIMutex(c.MaxAPICapacity)
		if err != nil {
			return nil, err
		}
		httpClient.Transport = transport.NewGovernedTransport(httpClient.Transport, apiMutex, c.Logger)
	}
	var orgUrl string
	var disableHTTPS bool
	if c.HttpProxy != "" {
		orgUrl = strings.TrimSuffix(c.HttpProxy, "/")
		disableHTTPS = strings.HasPrefix(orgUrl, "http://")
	} else {
		orgUrl = fmt.Sprintf("https://%v.%v", c.OrgName, c.Domain)
	}
	_, err = url.Parse(orgUrl)
	if err != nil {
		return nil, fmt.Errorf("malformed Okta API URL (org_name+base_url value, or http_proxy value): %+v", err)
	}

	setters := []okta.ConfigSetter{
		okta.WithOrgUrl(orgUrl),
		okta.WithCache(false),
		okta.WithHttpClientPtr(httpClient),
		okta.WithRateLimitMaxBackOff(int64(c.MaxWait)),
		okta.WithRequestTimeout(int64(c.RequestTimeout)),
		okta.WithRateLimitMaxRetries(int32(c.RetryCount)),
		okta.WithUserAgentExtra(OktaTerraformProviderUserAgent),
	}
	// v3 client also needs http proxy explicitly set
	if c.HttpProxy != "" {
		_url, err := url.Parse(c.HttpProxy)
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
	case c.AccessToken != "":
		setters = append(
			setters,
			okta.WithToken(c.AccessToken), okta.WithAuthorizationMode("Bearer"),
		)

	case c.ApiToken != "":
		setters = append(
			setters,
			okta.WithToken(c.ApiToken), okta.WithAuthorizationMode("SSWS"),
		)

	case c.PrivateKey != "":
		setters = append(
			setters,
			okta.WithPrivateKey(c.PrivateKey), okta.WithPrivateKeyId(c.PrivateKeyId), okta.WithScopes(c.Scopes), okta.WithClientId(c.ClientID), okta.WithAuthorizationMode("PrivateKey"),
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
	if c.Backoff {
		retryableClient := retryablehttp.NewClient()
		retryableClient.RetryWaitMin = time.Second * time.Duration(c.MinWait)
		retryableClient.RetryWaitMax = time.Second * time.Duration(c.MaxWait)
		retryableClient.RetryMax = c.RetryCount
		retryableClient.Logger = c.Logger
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
		c.Logger.Info(fmt.Sprintf("running with backoff http client, wait min %d, wait max %d, retry max %d", retryableClient.RetryWaitMin, retryableClient.RetryWaitMax, retryableClient.RetryMax))
	} else {
		httpClient = cleanhttp.DefaultClient()
		if debugHttpRequests {
			// Needed for pretty printing http protocol in a local developer environment, ignore deprecation warnings.
			//lint:ignore SA1019 used in developer mode only
			httpClient.Transport = logging.NewTransport("Okta", httpClient.Transport)
		} else {
			httpClient.Transport = logging.NewSubsystemLoggingHTTPTransport("Okta", httpClient.Transport)
		}
		c.Logger.Info("running with default http client")
	}

	// adds transport governor to retryable or default client
	if c.MaxAPICapacity > 0 && c.MaxAPICapacity < 100 {
		c.Logger.Info(fmt.Sprintf("running with experimental max_api_capacity configuration at %d%%", c.MaxAPICapacity))
		apiMutex, err := apimutex.NewAPIMutex(c.MaxAPICapacity)
		if err != nil {
			return nil, err
		}
		httpClient.Transport = transport.NewGovernedTransport(httpClient.Transport, apiMutex, c.Logger)
	}
	var orgUrl string
	var disableHTTPS bool
	if c.HttpProxy != "" {
		orgUrl = strings.TrimSuffix(c.HttpProxy, "/")
		disableHTTPS = strings.HasPrefix(orgUrl, "http://")
	} else {
		orgUrl = fmt.Sprintf("https://%v.%v", c.OrgName, c.Domain)
	}
	_, err = url.Parse(orgUrl)
	if err != nil {
		return nil, fmt.Errorf("malformed Okta API URL (org_name+base_url value, or http_proxy value): %+v", err)
	}

	setters := []v5okta.ConfigSetter{
		v5okta.WithOrgUrl(orgUrl),
		v5okta.WithCache(false),
		v5okta.WithHttpClientPtr(httpClient),
		v5okta.WithRateLimitMaxBackOff(int64(c.MaxWait)),
		v5okta.WithRequestTimeout(int64(c.RequestTimeout)),
		v5okta.WithRateLimitMaxRetries(int32(c.RetryCount)),
		v5okta.WithUserAgentExtra(OktaTerraformProviderUserAgent),
	}
	// v3 client also needs http proxy explicitly set
	if c.HttpProxy != "" {
		_url, err := url.Parse(c.HttpProxy)
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
	case c.AccessToken != "":
		setters = append(
			setters,
			v5okta.WithToken(c.AccessToken), v5okta.WithAuthorizationMode("Bearer"),
		)

	case c.ApiToken != "":
		setters = append(
			setters,
			v5okta.WithToken(c.ApiToken), v5okta.WithAuthorizationMode("SSWS"),
		)

	case c.PrivateKey != "":
		setters = append(
			setters,
			v5okta.WithPrivateKey(c.PrivateKey), v5okta.WithPrivateKeyId(c.PrivateKeyId), v5okta.WithScopes(c.Scopes), v5okta.WithClientId(c.ClientID), v5okta.WithAuthorizationMode("PrivateKey"),
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

const RetryOnStatusCodes contextKey = "retryOnStatusCodes"

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
	retryCodes, ok := ctx.Value(RetryOnStatusCodes).([]int)
	if ok && resp != nil && utils.ContainsInt(retryCodes, resp.StatusCode) {
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

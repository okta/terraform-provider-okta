package okta

import (
	"context"
	"errors"
	"fmt"
	"net/http"
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
	"github.com/okta/okta-sdk-golang/v3/okta"
	"github.com/okta/terraform-provider-okta/okta/internal/apimutex"
	"github.com/okta/terraform-provider-okta/okta/internal/transport"
	"github.com/okta/terraform-provider-okta/sdk"
)

const OktaTerraformProviderVersion = "4.1.0"
const OktaTerraformProviderUserAgent = "okta-terraform/" + OktaTerraformProviderVersion

func (adt *AddHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("User-Agent", "Okta Terraform Provider")
	return adt.T.RoundTrip(req)
}

type (
	// AddHeaderTransport used to tack on default headers to outgoing requests
	AddHeaderTransport struct {
		T http.RoundTripper
	}

	// Config contains our provider schema values and Okta clients
	Config struct {
		orgName          string
		domain           string
		httpProxy        string
		accessToken      string
		apiToken         string
		clientID         string
		privateKey       string
		privateKeyId     string
		scopes           []string
		retryCount       int
		parallelism      int
		backoff          bool
		minWait          int
		maxWait          int
		logLevel         int
		requestTimeout   int
		maxAPICapacity   int // experimental
		oktaClient       *sdk.Client
		v3Client         *okta.APIClient
		supplementClient *sdk.APISupplement
		logger           hclog.Logger
		queriedWellKnown bool
		classicOrg       bool
		timeOperations   TimeOperations
	}
)

// IsClassicOrg returns true if the org is a classic org. Does lazy evaluation
// of the well known endpoint so that VCR can record the transaction.
func (c *Config) IsClassicOrg(ctx context.Context) bool {
	if !c.queriedWellKnown {
		// Discover if the Okta Org is Classic or OIE
		org, _, err := c.supplementClient.GetWellKnownOktaOrganization(ctx)
		if err != nil {
			c.logger.Error("error querying GET /.well-known/okta-organization", "error", err)
			return c.classicOrg
		}

		c.classicOrg = (org.Pipeline == "v1") // v1 == Classic, idx == OIE
		c.queriedWellKnown = true
	}

	return c.classicOrg
}

func (c *Config) SetTimeOperations(op TimeOperations) {
	c.timeOperations = op
}

func (c *Config) loadAndValidate(ctx context.Context) error {
	c.logger = providerLogger(c)

	v3Client, err := oktaV3SDKClient(c)
	if err != nil {
		return err
	}
	c.v3Client = v3Client
	// NOTE: we want to share one http client across all SDK clients

	// TODO: remove sdk client when v3 client is fully utilized within the provider
	client, err := oktaSDKClient(c)
	if err != nil {
		return err
	}
	c.oktaClient = client

	// TODO: remove supplement client when v3 client is fully utilized within the provider
	re := client.CloneRequestExecutor()
	re.SetHTTPTransport(c.v3Client.GetConfig().HTTPClient.Transport)
	c.supplementClient = &sdk.APISupplement{
		RequestExecutor: re,
	}

	// NOTE: Don't make this call when VCR is playing/recording as it will occur
	// outsite of the VCR transport
	if os.Getenv("OKTA_VCR_TF_ACC") == "" {
		// NOTE: validate credentials during initial config with a call to
		// /api/v1/users/me
		if c.apiToken != "" {
			if _, _, err := c.v3Client.UserApi.GetUser(ctx, "me").Execute(); err != nil {
				return fmt.Errorf("error with v3 SDK client: %v", err)
			}
			if _, _, err := c.oktaClient.User.GetUser(ctx, "me"); err != nil {
				return fmt.Errorf("error with v2 SDK client: %v", err)
			}
		}
	}

	return nil
}

func (c *Config) handlePluginDefaults(ctx context.Context) error {
	var err error
	if c.orgName == "" && os.Getenv("OKTA_ORG_NAME") != "" {
		c.orgName = os.Getenv("OKTA_ORG_NAME")
	}
	if c.accessToken == "" && os.Getenv("OKTA_ACCESS_TOKEN") != "" {
		c.accessToken = os.Getenv("OKTA_ACCESS_TOKEN")
	}
	if c.apiToken == "" && os.Getenv("OKTA_API_TOKEN") != "" {
		c.apiToken = os.Getenv("OKTA_API_TOKEN")
	}
	if c.clientID == "" && os.Getenv("OKTA_API_CLIENT_ID") != "" {
		c.clientID = os.Getenv("OKTA_API_CLIENT_ID")
	}
	if v := os.Getenv("OKTA_API_SCOPES"); v != "" && len(c.scopes) == 0 {
		c.scopes = strings.Split(v, ",")
	}
	if c.privateKey == "" && os.Getenv("OKTA_API_PRIVATE_KEY") != "" {
		c.privateKey = os.Getenv("OKTA_API_PRIVATE_KEY")
	}
	if c.privateKeyId == "" && os.Getenv("OKTA_API_PRIVATE_KEY_ID") != "" {
		c.privateKeyId = os.Getenv("OKTA_API_PRIVATE_KEY_ID")
	}
	if c.domain == "" {
		if os.Getenv("OKTA_BASE_URL") != "" {
			c.domain = os.Getenv("OKTA_BASE_URL")
		} else {
			c.domain = "okta.com"
		}
	}
	if c.httpProxy == "" && os.Getenv("OKTA_HTTP_PROXY") != "" {
		c.httpProxy = os.Getenv("OKTA_HTTP_PROXY")
	}
	if c.maxAPICapacity == 0 {
		if os.Getenv("MAX_API_CAPACITY") != "" {
			mac, err := strconv.ParseInt(os.Getenv("MAX_API_CAPACITY"), 10, 64)
			if err != nil {
				return err
			}
			c.maxAPICapacity = int(mac)
		} else {
			c.maxAPICapacity = 100
		}
	}
	c.backoff = true
	c.minWait = 30
	c.maxWait = 300
	c.retryCount = 5
	c.parallelism = 1
	c.logLevel = int(hclog.Error)
	c.requestTimeout = 0
	return err
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
		} else {
			data.BaseURL = types.StringValue("okta.com")
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

	return err
}

func providerLogger(c *Config) hclog.Logger {
	logLevel := hclog.Level(c.logLevel)
	if os.Getenv("TF_LOG") != "" {
		logLevel = hclog.LevelFromString(os.Getenv("TF_LOG"))
	}

	return hclog.New(&hclog.LoggerOptions{
		Level:      logLevel,
		TimeFormat: "2006/01/02 03:04:05",
	})
}

// oktaSDKClient should be called with a primary http client that is utilezed
// throughout the provider
func oktaSDKClient(c *Config) (client *sdk.Client, err error) {
	httpClient := c.v3Client.GetConfig().HTTPClient
	var orgUrl string
	var disableHTTPS bool
	if c.httpProxy != "" {
		orgUrl = strings.TrimSuffix(c.httpProxy, "/")
		disableHTTPS = strings.HasPrefix(orgUrl, "http://")
	} else {
		orgUrl = fmt.Sprintf("https://%v.%v", c.orgName, c.domain)
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

	setters := []okta.ConfigSetter{
		okta.WithOrgUrl(orgUrl),
		okta.WithCache(false),
		okta.WithHttpClientPtr(httpClient),
		okta.WithRateLimitMaxBackOff(int64(c.maxWait)),
		okta.WithRequestTimeout(int64(c.requestTimeout)),
		okta.WithRateLimitMaxRetries(int32(c.retryCount)),
		okta.WithUserAgentExtra(OktaTerraformProviderUserAgent),
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

	config := okta.NewConfiguration(setters...)
	client = okta.NewAPIClient(config)
	return
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

type ProductionTimeOperations struct {
}

// DoNotRetry always retry in production
func (o *ProductionTimeOperations) DoNotRetry(err error) bool {
	return false
}

// Sleep facade to actual time.Sleep in production
func (o *ProductionTimeOperations) Sleep(d time.Duration) {
	time.Sleep(d)
}

// NewProductionTimeOperations new production time operations
func NewProductionTimeOperations() TimeOperations {
	return &ProductionTimeOperations{}
}

type TestTimeOperations struct {
}

// DoNotRetry tests do not retry when there is an error and VCR is recording
func (o *TestTimeOperations) DoNotRetry(err error) bool {
	return err != nil && os.Getenv("OKTA_VCR_TF_ACC") == "record"
}

// Sleep no sleeping when test is in VCR play mode
func (o *TestTimeOperations) Sleep(d time.Duration) {
	if os.Getenv("OKTA_VCR_TF_ACC") != "play" {
		time.Sleep(d)
	}
}

// NewTestTimeOperations new test time operations
func NewTestTimeOperations() TimeOperations {
	return &TestTimeOperations{}
}

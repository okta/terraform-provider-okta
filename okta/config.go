package okta

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/terraform-provider-okta/okta/internal/apimutex"
	"github.com/okta/terraform-provider-okta/okta/internal/transport"
	"github.com/okta/terraform-provider-okta/sdk"
)

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
		orgName           string
		domain            string
		apiToken          string
		clientID          string
		privateKey        string
		scopes            []string
		retryCount        int
		parallelism       int
		backoff           bool
		minWait           int
		maxWait           int
		logLevel          int
		requestTimeout    int
		maxAPICapacity    int // experimental
		apiTokenClient    *okta.Client
		accessTokenClient *okta.Client
		supplementClient  *sdk.APISupplement
		logger            hclog.Logger
		primaryClient     string
	}
)

func (c *Config) loadAndValidate(ctx context.Context) error {
	logLevel := hclog.Level(c.logLevel)
	if os.Getenv("TF_LOG") != "" {
		logLevel = hclog.LevelFromString(os.Getenv("TF_LOG"))
	}

	c.logger = hclog.New(&hclog.LoggerOptions{
		Level:      logLevel,
		TimeFormat: "2006/01/02 03:04:05",
	})
	var httpClient *http.Client
	if c.backoff {
		retryableClient := retryablehttp.NewClient()
		retryableClient.RetryWaitMin = time.Second * time.Duration(c.minWait)
		retryableClient.RetryWaitMax = time.Second * time.Duration(c.maxWait)
		retryableClient.RetryMax = c.retryCount
		retryableClient.Logger = c.logger
		retryableClient.HTTPClient.Transport = logging.NewTransport("Okta", retryableClient.HTTPClient.Transport)
		retryableClient.ErrorHandler = errHandler
		retryableClient.CheckRetry = checkRetry
		httpClient = retryableClient.StandardClient()
		c.logger.Info(fmt.Sprintf("running with backoff http client, wait min %d, wait max %d, retry max %d", retryableClient.RetryWaitMin, retryableClient.RetryWaitMax, retryableClient.RetryMax))
	} else {
		httpClient = cleanhttp.DefaultClient()
		httpClient.Transport = logging.NewTransport("Okta", httpClient.Transport)
		c.logger.Info("running with default http client")
	}

	// adds transport governor to retryable or default client
	if c.maxAPICapacity > 0 && c.maxAPICapacity < 100 {
		c.logger.Info(fmt.Sprintf("running with experimental max_api_capacity configuration at %d%%", c.maxAPICapacity))
		apiMutex, err := apimutex.NewAPIMutex(c.maxAPICapacity)
		if err != nil {
			return err
		}
		httpClient.Transport = transport.NewGovernedTransport(httpClient.Transport, apiMutex, c.logger)
	}

	setters := []okta.ConfigSetter{
		okta.WithOrgUrl(fmt.Sprintf("https://%v.%v", c.orgName, c.domain)),
		okta.WithCache(false),
		okta.WithHttpClientPtr(httpClient),
		okta.WithRateLimitMaxBackOff(int64(c.maxWait)),
		okta.WithRequestTimeout(int64(c.requestTimeout)),
		okta.WithRateLimitMaxRetries(int32(c.retryCount)),
		okta.WithUserAgentExtra("okta-terraform/3.13.13"),
	}
	var err error
	if c.apiToken != "" {
		_, c.apiTokenClient, err = okta.NewClient(
			context.Background(),
			append(setters,
				okta.WithToken(c.apiToken))...,
		)
		if err != nil {
			return err
		}
	}
	if c.clientID != "" {
		_, c.accessTokenClient, err = okta.NewClient(
			context.Background(),
			append(setters,
				okta.WithAuthorizationMode("PrivateKey"),
				okta.WithClientId(c.clientID),
				okta.WithPrivateKey(c.privateKey),
				okta.WithScopes(c.scopes))...,
		)
		if err != nil {
			return err
		}
	}
	if c.apiTokenClient == nil && c.accessTokenClient == nil {
		return errors.New("please set either 'api_token' or 'client_id', 'privateKey' and 'scopes' (or both) in order to use the provider")
	}
	if c.primaryClient == "api_token_client" && c.apiTokenClient == nil {
		return errors.New("please set 'api_token' in order to use 'api_token_client' as a primary client")
	}
	if c.primaryClient == "access_token_client" && c.accessTokenClient == nil {
		return errors.New("please set 'client_id', 'privateKey' and 'scopes' in order to use 'access_token_client' as a primary client")
	}
	if c.primaryClient == "" {
		if c.apiTokenClient != nil {
			c.primaryClient = "api_token_client"
		} else {
			c.primaryClient = "access_token_client"
		}
	}
	if c.apiTokenClient != nil {
		if _, _, err := c.apiTokenClient.User.GetUser(ctx, "me"); err != nil {
			return err
		}
	}
	c.supplementClient = sdk.NewAPISupplement(c.apiTokenClient, c.accessTokenClient, c.primaryClient)
	return nil
}

func errHandler(resp *http.Response, err error, numTries int) (*http.Response, error) {
	if err != nil {
		return resp, err
	}
	defer resp.Body.Close()
	err = okta.CheckResponseForError(resp)
	if err != nil {
		var oErr *okta.Error
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
// 		ctx = context.WithValue(ctx, retryOnStatusCodes, []int{404, 409})
//
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

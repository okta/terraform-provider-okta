package api

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
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	v5okta "github.com/okta/okta-sdk-golang/v5/okta"
	v6okta "github.com/okta/okta-sdk-golang/v6/okta"
	"github.com/okta/terraform-provider-okta/okta/internal/apimutex"
	"github.com/okta/terraform-provider-okta/okta/internal/transport"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/okta/version"
	"github.com/okta/terraform-provider-okta/sdk"
)

func getV6ClientConfig(c *OktaAPIConfig) (*v6okta.Configuration, *v6okta.APIClient, error) {
	var httpClient *http.Client
	if c.Backoff {
		retryableClient := retryablehttp.NewClient()
		retryableClient.RetryWaitMin = time.Second * time.Duration(c.MinWait)
		retryableClient.RetryWaitMax = time.Second * time.Duration(c.MaxWait)
		retryableClient.RetryMax = c.RetryCount
		retryableClient.Logger = c.Logger
		retryableClient.HTTPClient.Transport = logging.NewSubsystemLoggingHTTPTransport("Okta", retryableClient.HTTPClient.Transport)
		retryableClient.ErrorHandler = errHandler
		retryableClient.CheckRetry = checkRetry
		httpClient = retryableClient.StandardClient()
		c.Logger.Info(fmt.Sprintf("v6 running with backoff http client, wait min %d, wait max %d, retry max %d", retryableClient.RetryWaitMin, retryableClient.RetryWaitMax, retryableClient.RetryMax))
	} else {
		httpClient = cleanhttp.DefaultClient()
		httpClient.Transport = logging.NewSubsystemLoggingHTTPTransport("Okta", httpClient.Transport)
		c.Logger.Info("v6 running with default http client")
	}

	// adds transport governor to retryable or default client
	if c.MaxAPICapacity > 0 && c.MaxAPICapacity < 100 {
		c.Logger.Info(fmt.Sprintf("v6 running with experimental max_api_capacity configuration at %d%%", c.MaxAPICapacity))
		apiMutex, err := apimutex.NewAPIMutex(c.MaxAPICapacity)
		if err != nil {
			return nil, nil, err
		}
		httpClient.Transport = transport.NewGovernedTransport(httpClient.Transport, apiMutex, c.Logger)
	}
	var orgURL string
	var disableHTTPS bool
	if c.HttpProxy != "" {
		orgURL = strings.TrimSuffix(c.HttpProxy, "/")
		disableHTTPS = strings.HasPrefix(orgURL, "http://")
	} else {
		orgURL = fmt.Sprintf("https://%v.%v", c.OrgName, c.Domain)
	}
	_, err := url.Parse(orgURL)
	if err != nil {
		return nil, nil, fmt.Errorf("malformed Okta API URL (org_name+base_url value, or http_proxy value): %+v", err)
	}

	setters := []v6okta.ConfigSetter{
		v6okta.WithOrgUrl(orgURL),
		v6okta.WithCache(false),
		v6okta.WithHttpClientPtr(httpClient),
		v6okta.WithRateLimitMaxBackOff(int64(c.MaxWait)),
		v6okta.WithRequestTimeout(int64(c.RequestTimeout)),
		v6okta.WithRateLimitMaxRetries(int32(c.RetryCount)),
		v6okta.WithUserAgentExtra(version.OktaTerraformProviderUserAgent),
	}
	// v6 client also needs http proxy explicitly set
	if c.HttpProxy != "" {
		_url, err := url.Parse(c.HttpProxy)
		if err != nil {
			return nil, nil, err
		}
		host := v6okta.WithProxyHost(_url.Hostname())
		setters = append(setters, host)

		sPort := _url.Port()
		if sPort == "" {
			sPort = "80"
		}
		iPort, err := strconv.Atoi(sPort)
		if err != nil {
			return nil, nil, err
		}
		port := v6okta.WithProxyPort(int32(iPort))
		setters = append(setters, port)
	}

	switch {
	case c.AccessToken != "":
		setters = append(
			setters,
			v6okta.WithToken(c.AccessToken), v6okta.WithAuthorizationMode("Bearer"),
		)

	case c.ApiToken != "":
		setters = append(
			setters,
			v6okta.WithToken(c.ApiToken), v6okta.WithAuthorizationMode("SSWS"),
		)

	case c.PrivateKey != "":
		setters = append(
			setters,
			v6okta.WithPrivateKey(c.PrivateKey), v6okta.WithPrivateKeyId(c.PrivateKeyId), v6okta.WithScopes(c.Scopes), v6okta.WithClientId(c.ClientID), v6okta.WithAuthorizationMode("PrivateKey"),
		)
	}

	if disableHTTPS {
		setters = append(setters, v6okta.WithTestingDisableHttpsCheck(true))
	}

	config, err := v6okta.NewConfiguration(setters...)
	if err != nil {
		return nil, nil, err
	}
	return config, nil, nil
}

func getV5ClientConfig(c *OktaAPIConfig) (*v5okta.Configuration, *v5okta.APIClient, error) {
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
		c.Logger.Info(fmt.Sprintf("v5 running with backoff http client, wait min %d, wait max %d, retry max %d", retryableClient.RetryWaitMin, retryableClient.RetryWaitMax, retryableClient.RetryMax))
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
			return nil, nil, err
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
	_, err := url.Parse(orgUrl)
	if err != nil {
		return nil, nil, fmt.Errorf("malformed Okta API URL (org_name+base_url value, or http_proxy value): %+v", err)
	}

	setters := []v5okta.ConfigSetter{
		v5okta.WithOrgUrl(orgUrl),
		v5okta.WithCache(false),
		v5okta.WithHttpClientPtr(httpClient),
		v5okta.WithRateLimitMaxBackOff(int64(c.MaxWait)),
		v5okta.WithRequestTimeout(int64(c.RequestTimeout)),
		v5okta.WithRateLimitMaxRetries(int32(c.RetryCount)),
		v5okta.WithUserAgentExtra(version.OktaTerraformProviderUserAgent),
	}
	// v3 client also needs http proxy explicitly set
	if c.HttpProxy != "" {
		_url, err := url.Parse(c.HttpProxy)
		if err != nil {
			return nil, nil, err
		}
		host := v5okta.WithProxyHost(_url.Hostname())
		setters = append(setters, host)

		sPort := _url.Port()
		if sPort == "" {
			sPort = "80"
		}
		iPort, err := strconv.Atoi(sPort)
		if err != nil {
			return nil, nil, err
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
		return nil, nil, err
	}
	return config, nil, nil
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

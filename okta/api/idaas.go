package api

import (
	"context"
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
	"github.com/okta/okta-sdk-golang/v4/okta"
	v5okta "github.com/okta/okta-sdk-golang/v5/okta"
	v6okta "github.com/okta/okta-sdk-golang/v6/okta"
	"github.com/okta/terraform-provider-okta/okta/internal/apimutex"
	"github.com/okta/terraform-provider-okta/okta/internal/transport"
	"github.com/okta/terraform-provider-okta/okta/version"
	"github.com/okta/terraform-provider-okta/sdk"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ OktaIDaaSClient = &iDaaSAPIClient{}
)

type contextKey string

const RetryOnStatusCodes contextKey = "retryOnStatusCodes"

type OktaIDaaSClient interface {
	OktaSDKClientV6() *v6okta.APIClient
	OktaSDKClientV5() *v5okta.APIClient
	OktaSDKClientV3() *okta.APIClient
	OktaSDKClientV2() *sdk.Client
	OktaSDKSupplementClient() *sdk.APISupplement
}

type OktaAPIConfig struct {
	AccessToken    string
	ApiToken       string
	Backoff        bool
	ClientID       string
	Domain         string
	HttpProxy      string
	Logger         hclog.Logger
	MaxAPICapacity int
	MaxWait        int
	MinWait        int
	OrgName        string
	PrivateKey     string
	PrivateKeyId   string
	RequestTimeout int
	RetryCount     int
	Scopes         []string
}

type iDaaSAPIClient struct {
	oktaSDKClientV6         *v6okta.APIClient
	oktaSDKClientV5         *v5okta.APIClient
	oktaSDKClientV3         *okta.APIClient
	oktaSDKClientV2         *sdk.Client
	oktaSDKSupplementClient *sdk.APISupplement
}

func (c *iDaaSAPIClient) OktaSDKClientV6() *v6okta.APIClient {
	return c.oktaSDKClientV6
}

func (c *iDaaSAPIClient) OktaSDKClientV5() *v5okta.APIClient {
	return c.oktaSDKClientV5
}

func (c *iDaaSAPIClient) OktaSDKClientV3() *okta.APIClient {
	return c.oktaSDKClientV3
}

func (c *iDaaSAPIClient) OktaSDKClientV2() *sdk.Client {
	return c.oktaSDKClientV2
}

func (c *iDaaSAPIClient) OktaSDKSupplementClient() *sdk.APISupplement {
	return c.oktaSDKSupplementClient
}

func NewOktaIDaaSAPIClient(c *OktaAPIConfig) (client OktaIDaaSClient, err error) {
	v6client, err := oktaV6SDKClient(c)
	if err != nil {
		return
	}

	v5Client, err := oktaV5SDKClient(c)
	if err != nil {
		return
	}

	v3Client, err := oktaV3SDKClient(c)
	if err != nil {
		return
	}

	httpClient := v3Client.GetConfig().HTTPClient
	v2Client, err := oktaV2SDKClient(httpClient, c)
	if err != nil {
		return
	}

	re := v2Client.CloneRequestExecutor()
	re.SetHTTPTransport(v3Client.GetConfig().HTTPClient.Transport)
	supClient := &sdk.APISupplement{
		RequestExecutor: re,
	}

	client = &iDaaSAPIClient{
		oktaSDKClientV6:         v6client,
		oktaSDKClientV5:         v5Client,
		oktaSDKClientV3:         v3Client,
		oktaSDKClientV2:         v2Client,
		oktaSDKSupplementClient: supClient,
	}

	return
}

func oktaV6SDKClient(c *OktaAPIConfig) (client *v6okta.APIClient, err error) {
	config, apiClient, err := getV6ClientConfig(c)
	if err != nil {
		return apiClient, err
	}
	client = v6okta.NewAPIClient(config)
	return client, nil
}

func oktaV5SDKClient(c *OktaAPIConfig) (client *v5okta.APIClient, err error) {
	config, apiClient, err := getV5ClientConfig(c)
	if err != nil {
		return apiClient, err
	}
	client = v5okta.NewAPIClient(config)
	return client, nil
}

func oktaV3SDKClient(c *OktaAPIConfig) (client *okta.APIClient, err error) {
	config, apiClient, configErr := GetV3ClientConfig(c)
	if configErr != nil {
		return apiClient, configErr
	}
	client = okta.NewAPIClient(config)
	return client, nil
}

func GetV3ClientConfig(c *OktaAPIConfig) (*okta.Configuration, *okta.APIClient, error) {
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
		c.Logger.Info(fmt.Sprintf("v3 running with backoff http client, wait min %d, wait max %d, retry max %d", retryableClient.RetryWaitMin, retryableClient.RetryWaitMax, retryableClient.RetryMax))
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

	setters := []okta.ConfigSetter{
		okta.WithOrgUrl(orgUrl),
		okta.WithCache(false),
		okta.WithHttpClientPtr(httpClient),
		okta.WithRateLimitMaxBackOff(int64(c.MaxWait)),
		okta.WithRequestTimeout(int64(c.RequestTimeout)),
		okta.WithRateLimitMaxRetries(int32(c.RetryCount)),
		okta.WithUserAgentExtra(version.OktaTerraformProviderUserAgent),
	}
	// v3 client also needs http proxy explicitly set
	if c.HttpProxy != "" {
		_url, err := url.Parse(c.HttpProxy)
		if err != nil {
			return nil, nil, err
		}
		host := okta.WithProxyHost(_url.Hostname())
		setters = append(setters, host)

		sPort := _url.Port()
		if sPort == "" {
			sPort = "80"
		}
		iPort, err := strconv.Atoi(sPort)
		if err != nil {
			return nil, nil, err
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
		return nil, nil, err
	}
	return config, nil, nil
}

func oktaV2SDKClient(httpClient *http.Client, c *OktaAPIConfig) (client *sdk.Client, err error) {
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
		sdk.WithUserAgentExtra(version.OktaTerraformProviderUserAgent),
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

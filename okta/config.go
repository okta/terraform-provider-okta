package okta

import (
	"context"
	"fmt"
	"net/http"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
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
		orgName          string
		domain           string
		apiToken         string
		retryCount       int
		parallelism      int
		backoff          bool
		maxWait          int
		logLevel         int
		oktaClient       *okta.Client
		supplementClient *sdk.ApiSupplement
		logger           hclog.Logger
	}
)

func (c *Config) loadAndValidate() error {
	c.logger = hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Level(c.logLevel),
		TimeFormat: "2006/01/02 03:04:05",
	})

	httpClient := retryablehttp.NewClient()
	httpClient.Logger = c.logger
	httpClient.HTTPClient.Transport = logging.NewTransport("Okta", httpClient.HTTPClient.Transport)

	_, client, err := okta.NewClient(
		context.Background(),
		okta.WithOrgUrl(fmt.Sprintf("https://%v.%v", c.orgName, c.domain)),
		okta.WithToken(c.apiToken),
		okta.WithCache(false),
		okta.WithHttpClient(*httpClient.StandardClient()),
		okta.WithRateLimitMaxBackOff(int64(c.maxWait)),
		okta.WithRateLimitMaxRetries(int32(c.retryCount)),
		okta.WithUserAgentExtra("okta-terraform/3.7.3"),
	)
	if err != nil {
		return err
	}
	c.oktaClient = client
	c.supplementClient = &sdk.ApiSupplement{
		BaseURL:         fmt.Sprintf("https://%s.%s", c.orgName, c.domain),
		Client:          httpClient.StandardClient(),
		Token:           c.apiToken,
		RequestExecutor: client.GetRequestExecutor(),
	}
	return nil
}

package okta

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-cleanhttp"
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
		oktaClient       *okta.Client
		supplementClient *sdk.ApiSupplement
		ctx              context.Context
	}
)

func (c *Config) loadAndValidate() error {
	httpClient := cleanhttp.DefaultClient()
	httpClient.Transport = logging.NewTransport("Okta", httpClient.Transport)

	orgURL := fmt.Sprintf("https://%v.%v", c.orgName, c.domain)
	ctx, client, err := okta.NewClient(
		context.Background(),
		okta.WithOrgUrl(orgURL),
		okta.WithToken(c.apiToken),
		okta.WithCache(false),
		okta.WithHttpClient(*httpClient),
		okta.WithRequestTimeout(int64(c.maxWait)),
		okta.WithRateLimitMaxRetries(int32(c.retryCount)),
		okta.WithUserAgentExtra("okta-terraform/3.6.0"),
	)
	if err != nil {
		return err
	}
	c.supplementClient = &sdk.ApiSupplement{
		BaseURL:         fmt.Sprintf("https://%s.%s", c.orgName, c.domain),
		Client:          httpClient,
		Token:           c.apiToken,
		RequestExecutor: client.GetRequestExecutor(),
	}

	// add the Okta SDK client object to Config
	c.oktaClient = client
	c.ctx = ctx
	return nil
}

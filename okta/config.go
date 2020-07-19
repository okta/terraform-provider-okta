package okta

import (
	"context"
	"fmt"
	"net/http"

	"github.com/terraform-providers/terraform-provider-okta/sdk"

	articulateOkta "github.com/articulate/oktasdk-go/okta"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/terraform-plugin-sdk/helper/logging"
	"github.com/okta/okta-sdk-golang/v2/okta"
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
		orgName      string
		domain       string
		apiToken     string
		retryCount   int
		parallelism  int
		waitForReset bool
		minWait      int
		maxWait      int

		articulateOktaClient *articulateOkta.Client
		oktaClient           *okta.Client
		ctx                  context.Context
		supplementClient     *sdk.ApiSupplement
	}
)

func (c *Config) loadAndValidate() error {
	httpClient := cleanhttp.DefaultClient()
	httpClient.Transport = logging.NewTransport("Okta", httpClient.Transport)

	articulateClient, err := articulateOkta.NewClientWithDomain(httpClient, c.orgName, c.domain, c.apiToken)

	// add the Articulate Okta client object to Config
	c.articulateOktaClient = articulateClient

	if err != nil {
		return fmt.Errorf("[ERROR] Error creating Articulate Okta client: %v", err)
	}

	orgUrl := fmt.Sprintf("https://%v.%v", c.orgName, c.domain)

	ctx, client, err := okta.NewClient(
		context.Background(),
		okta.WithOrgUrl(orgUrl),
		okta.WithToken(c.apiToken),
		okta.WithCache(false),
		okta.WithConnectionTimeout(int32(c.minWait)*60), // Using a i32 conversion due to times use of i64
		okta.WithRequestTimeout(int32(c.maxWait)*60),    // TODO change these types to i64 upstream
		okta.WithRateLimitMaxRetries(int32(c.retryCount)),
		okta.WithHttpClient(*httpClient),
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

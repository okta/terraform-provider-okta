// Package provider is configuration okta terraform provider
package provider

import (
	"context"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/okta/terraform-provider-okta/okta/config"
	"github.com/okta/terraform-provider-okta/okta/services/idaas"
)

// Provider establishes a client connection to an okta site
// determined by its schema string values
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"org_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The organization to manage in Okta.",
			},
			"access_token": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Bearer token granting privileges to Okta API.",
				ConflictsWith: []string{"api_token", "client_id", "scopes", "private_key"},
			},
			"api_token": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "API Token granting privileges to Okta API.",
				ConflictsWith: []string{"access_token", "client_id", "scopes", "private_key"},
			},
			"client_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "API Token granting privileges to Okta API.",
				ConflictsWith: []string{"access_token", "api_token"},
			},
			"scopes": {
				Type:          schema.TypeSet,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Description:   "API Token granting privileges to Okta API.",
				ConflictsWith: []string{"access_token", "api_token"},
			},
			"private_key": {
				Optional:      true,
				Type:          schema.TypeString,
				Description:   "API Token granting privileges to Okta API.",
				ConflictsWith: []string{"access_token", "api_token"},
			},
			"private_key_id": {
				Optional:      true,
				Type:          schema.TypeString,
				Description:   "API Token Id granting privileges to Okta API.",
				ConflictsWith: []string{"api_token"},
			},
			"base_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Okta url. (Use 'oktapreview.com' for Okta testing)",
			},
			"http_proxy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Alternate HTTP proxy of scheme://hostname or scheme://hostname:port format",
			},
			"backoff": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Use exponential back off strategy for rate limits.",
			},
			"min_wait_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "minimum seconds to wait when rate limit is hit. We use exponential backoffs when backoff is enabled.",
			},
			"max_wait_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "maximum seconds to wait when rate limit is hit. We use exponential backoffs when backoff is enabled.",
			},
			"max_retries": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateDiagFunc: intAtMost(100),
				Description:      "maximum number of retries to attempt before erroring out.",
			},
			"parallelism": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Number of concurrent requests to make within a resource where bulk operations are not possible. Take note of https://developer.okta.com/docs/api/getting_started/rate-limits.",
			},
			"log_level": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateDiagFunc: intBetween(1, 5),
				Description:      "providers log level. Minimum is 1 (TRACE), and maximum is 5 (ERROR)",
			},
			"max_api_capacity": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateDiagFunc: intBetween(1, 100),
				Description: "Sets what percentage of capacity the provider can use of the total rate limit " +
					"capacity while making calls to the Okta management API endpoints. Okta API operates in one minute buckets. " +
					"See Okta Management API Rate Limits: https://developer.okta.com/docs/reference/rl-global-mgmt/",
			},
			"request_timeout": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateDiagFunc: intBetween(0, 300),
				Description:      "Timeout for single request (in seconds) which is made to Okta, the default is `0` (means no limit is set). The maximum value can be `300`.",
			},
		},
		ResourcesMap:         idaas.ProviderResources(),
		DataSourcesMap:       idaas.ProviderDataSources(),
		ConfigureContextFunc: providerConfigure,
	}
}

// providerConfigure is only called once when a terraform command is run but it
// will be called many times while running different ACC tests
func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	log.Printf("[INFO] Initializing Okta client")
	cfg := config.NewConfig(d)
	if err := cfg.LoadAPIClient(); err != nil {
		return nil, diag.Errorf("[ERROR] failed to load sdk clients: %v", err)
	}
	cfg.SetTimeOperations(config.NewProductionTimeOperations())

	// NOTE: production runtime needs to know about VCR test environment for
	// this one case where the validate function calls GET /api/v1/users/me to
	// quickly verify if the operator's auth settings are correct.
	if os.Getenv("OKTA_VCR_TF_ACC") == "" {
		if err := cfg.VerifyCredentials(ctx); err != nil {
			return nil, diag.Errorf("[ERROR] failed validate configuration: %v", err)
		}
	}

	return cfg, nil
}

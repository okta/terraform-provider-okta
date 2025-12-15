package governance_test

import (
	"os"
	"testing"

	schema_sdk "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/api"
	"github.com/okta/terraform-provider-okta/okta/config"
	"github.com/okta/terraform-provider-okta/okta/provider"
)

var (
	// governanceAPIClientForTestUtil is a shared test client used by all governance tests
	// for CheckDestroy and other verification functions that need to make API calls
	governanceAPIClientForTestUtil api.OktaGovernanceClient
)

func init() {
	// Set up fake credentials for VCR playback mode
	if os.Getenv("OKTA_VCR_TF_ACC") == "play" {
		os.Setenv("OKTA_API_TOKEN", "token")
		os.Setenv("OKTA_BASE_URL", "dne-okta.com")
		if os.Getenv("OKTA_VCR_CASSETTE") != "" {
			os.Setenv("OKTA_ORG_NAME", os.Getenv("OKTA_VCR_CASSETTE"))
		}
	}

	// Initialize the shared governance client
	t := &testing.T{}
	governanceAPIClientForTestUtil = GovernanceClientForTest(t)
}

// GovernanceClientForTest creates a governance API client for testing
func GovernanceClientForTest(t *testing.T) api.OktaGovernanceClient {
	p := provider.Provider()
	d := resourceDataForTest(t, p.Schema)
	cfg := config.NewConfig(d)
	err := cfg.LoadAPIClient()
	if err != nil {
		t.Fatalf("Failed to load API client: %v", err)
	}
	return cfg.OktaGovernanceClient
}

// resourceDataForTest creates a ResourceData for testing with config values from environment
func resourceDataForTest(t *testing.T, s map[string]*schema_sdk.Schema) *schema_sdk.ResourceData {
	configValues := configValuesForTest()
	emptyConfigMap := map[string]interface{}{}
	d := schema_sdk.TestResourceDataRaw(t, s, emptyConfigMap)

	if len(configValues) > 0 {
		for k, v := range configValues {
			if err := d.Set(k, v); err != nil {
				t.Fatalf("Failed to set %s: %v", k, err)
			}
		}
	}
	return d
}

// configValuesForTest retrieves configuration values from environment variables
func configValuesForTest() map[string]interface{} {
	config := make(map[string]interface{})

	if v := os.Getenv("OKTA_ORG_NAME"); v != "" {
		config["org_name"] = v
	}
	if v := os.Getenv("OKTA_BASE_URL"); v != "" {
		config["base_url"] = v
	}
	if v := os.Getenv("OKTA_API_TOKEN"); v != "" {
		config["api_token"] = v
	}
	if v := os.Getenv("OKTA_API_CLIENT_ID"); v != "" {
		config["client_id"] = v
	}
	if v := os.Getenv("OKTA_API_SCOPES"); v != "" {
		config["scopes"] = v
	}
	if v := os.Getenv("OKTA_API_PRIVATE_KEY"); v != "" {
		config["private_key"] = v
	}
	if v := os.Getenv("OKTA_API_PRIVATE_KEY_ID"); v != "" {
		config["private_key_id"] = v
	}

	return config
}

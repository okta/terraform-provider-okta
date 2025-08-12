package governance_test

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	schema_sdk "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/api"
	"github.com/okta/terraform-provider-okta/okta/config"
	"github.com/okta/terraform-provider-okta/okta/provider"
	"net/http"
	"os"
	"strings"
	"testing"
)

var (
	sweeperLogger                  hclog.Logger
	sweeperLogLevel                hclog.Level
	governanceAPIClientForTestUtil api.OktaGovernanceClient
)

func init() {
	sweeperLogLevel = hclog.Warn
	if os.Getenv("TF_LOG") != "" {
		sweeperLogLevel = hclog.LevelFromString(os.Getenv("TF_LOG"))
	}
	sweeperLogger = hclog.New(&hclog.LoggerOptions{
		Level:      sweeperLogLevel,
		TimeFormat: "2006/01/02 03:04:05",
	})

	if os.Getenv("OKTA_VCR_TF_ACC") == "play" {
		os.Setenv("OKTA_API_TOKEN", "token")
		os.Setenv("OKTA_BASE_URL", "dne-okta.com")
		if os.Getenv("OKTA_VCR_CASSETTE") != "" {
			os.Setenv("OKTA_ORG_NAME", os.Getenv("OKTA_VCR_CASSETTE"))
		}
	}
	t := &testing.T{}

	fmt.Println("INSIDE GOVERNANCE_TEST")
	governanceAPIClientForTestUtil = GovernanceClientForTest(t)

	transport := governanceAPIClientForTestUtil.OktaIGSDKClientV5().GetConfig().HTTPClient.Transport
	governanceAPIClientForTestUtil.OktaIGSDKClientV5().GetConfig().HTTPClient.Transport = roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		fmt.Println("[VCR] HTTP Request:", req.Method, req.URL.String())
		return transport.RoundTrip(req)
	})

	fmt.Printf("Governance HTTPClient: %p\n", governanceAPIClientForTestUtil.OktaIGSDKClientV5().GetConfig().HTTPClient)
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func GovernanceClientForTest(t *testing.T) api.OktaGovernanceClient {
	p := provider.Provider()
	d := resourceDataForTest(t, p.Schema)
	cfg := config.NewConfig(d)
	_ = cfg.LoadAPIClient()
	return cfg.OktaGovernanceClient
}

func resourceDataForTest(t *testing.T, s map[string]*schema_sdk.Schema) *schema_sdk.ResourceData {
	configValues := configValuesForTest()
	emptyConfigMap := map[string]interface{}{}
	d := schema_sdk.TestResourceDataRaw(t, s, emptyConfigMap)

	if len(configValues) > 0 {
		for k, v := range configValues {
			// lintignore:R001
			_ = d.Set(k, v)
		}
	}

	return d
}

func configValuesForTest() map[string]interface{} {
	return map[string]interface{}{
		"access_token":   os.Getenv("OKTA_ACCESS_TOKEN"),
		"api_token":      os.Getenv("OKTA_API_TOKEN"),
		"org_name":       os.Getenv("OKTA_ORG_NAME"),
		"base_url":       os.Getenv("OKTA_BASE_URL"),
		"client_id":      os.Getenv("OKTA_API_CLIENT_ID"),
		"scopes":         strings.Split(os.Getenv("OKTA_API_SCOPES"), ","),
		"private_key":    os.Getenv("OKTA_API_PRIVATE_KEY"),
		"private_key_id": os.Getenv("OKTA_API_PRIVATE_KEY_ID"),
		"http_proxy":     os.Getenv("OKTA_HTTP_PROXY"),
		"log_level":      hclog.LevelFromString(os.Getenv("TF_LOG")),
	}
}

func TestMain(m *testing.M) {
	// TF_VAR_hostname allows the real hostname to be scripted into the config tests
	// see examples/resources/okta_resource_set/basic.tf
	if os.Getenv("TF_VAR_hostname") == "" {
		os.Setenv("TF_VAR_hostname", fmt.Sprintf("%s.%s", os.Getenv("OKTA_ORG_NAME"), os.Getenv("OKTA_BASE_URL")))
	}
	os.Setenv("TF_VAR_orgID", os.Getenv("OKTA_ORG_ID"))

	// NOTE: Acceptance test sweepers are necessary to prevent dangling
	// resources.
	// NOTE: Don't run sweepers if we are playing back VCR as nothing should be
	// going over the wire
	if os.Getenv("OKTA_VCR_TF_ACC") != "play" {
	}

	resource.TestMain(m)
}

package okta

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/terraform-provider-okta/sdk"
)

var (
	testAccProvidersFactories map[string]func() (*schema.Provider, error)
	testAccProvider           *schema.Provider
	testSDKClient             *okta.Client
	testSupplementClient      *sdk.APISupplement
)

func init() {
	testAccProvider = Provider()
	testAccProvidersFactories = map[string]func() (*schema.Provider, error){
		"okta": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}

	// We need to be able to query the SDK with an Okta SDK golang client that
	// is outside of the client that terraform provider creates. This is because
	// tests may need to query the okta API for status and the Terraform SDK
	// doesn't expose the provider's meta data where we store the provider's
	// config until after tests have completed.
	if os.Getenv("TF_ACC") != "" {
		// only set up for acceptance tests
		config := &Config{
			orgName: os.Getenv("OKTA_ORG_NAME"),
			domain:  os.Getenv("OKTA_BASE_URL"),
		}
		config.logger = providerLogger(config)
		testSDKClient, _ = oktaSDKClient(config)
		testSupplementClient = &sdk.APISupplement{
			RequestExecutor: testSDKClient.CloneRequestExecutor(),
		}
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	_ = Provider()
}

func oktaConfig() (*Config, error) {
	config := &Config{
		orgName:        os.Getenv("OKTA_ORG_NAME"),
		apiToken:       os.Getenv("OKTA_API_TOKEN"),
		httpProxy:      os.Getenv("OKTA_HTTP_PROXY"),
		clientID:       os.Getenv("OKTA_API_CLIENT_ID"),
		privateKey:     os.Getenv("OKTA_API_PRIVATE_KEY"),
		privateKeyId:   os.Getenv("OKTA_API_PRIVATE_KEY_ID"),
		scopes:         strings.Split(os.Getenv("OKTA_API_SCOPES"), ","),
		domain:         os.Getenv("OKTA_BASE_URL"),
		parallelism:    1,
		retryCount:     10,
		maxWait:        30,
		requestTimeout: 60,
		maxAPICapacity: 80,
	}
	if err := config.loadAndValidate(context.Background()); err != nil {
		return config, fmt.Errorf("error initializing Okta client: %v", err)
	}
	return config, nil
}

// testOIEOnlyAccPreCheck is a resource.test PreCheck function that will place a
// logical skip of OIE tests when tests are run against a classic org.
func testOIEOnlyAccPreCheck(t *testing.T) func() {
	return func() {
		err := accPreCheck()
		if err != nil {
			t.Fatalf("%v", err)
		}

		org, _, err := testSupplementClient.GetWellKnownOktaOrganization(context.TODO())
		if err != nil {
			t.Fatalf("%v", err)
		}

		// v1 == Classic, idx == OIE
		if org.Pipeline != "idx" {
			t.Skipf("%q test is for OIE orgs only", t.Name())
		}
	}
}

func testAccPreCheck(t *testing.T) func() {
	return func() {
		err := accPreCheck()
		if err != nil {
			t.Fatalf("%v", err)
		}
	}
}

func accPreCheck() error {
	if v := os.Getenv("OKTA_ORG_NAME"); v == "" {
		return errors.New("OKTA_ORG_NAME must be set for acceptance tests")
	}
	token := os.Getenv("OKTA_API_TOKEN")
	clientID := os.Getenv("OKTA_API_CLIENT_ID")
	privateKey := os.Getenv("OKTA_API_PRIVATE_KEY")
	privateKeyId := os.Getenv("OKTA_API_PRIVATE_KEY_ID")
	scopes := os.Getenv("OKTA_API_SCOPES")
	if token == "" && (clientID == "" || scopes == "" || privateKey == "" || privateKeyId == "") {
		return errors.New("either OKTA_API_TOKEN or OKTA_API_CLIENT_ID, OKTA_API_SCOPES and OKTA_API_PRIVATE_KEY must be set for acceptance tests")
	}
	return nil
}

func TestHTTPProxy(t *testing.T) {
	var handledUserRequest bool

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("x-rate-limit-reset", "0")
		w.Header().Set("x-rate-limit-limit", "0")
		w.Header().Set("x-rate-limit-limit", "0")
		w.Header().Set("x-rate-limit-remaining", "0")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&okta.User{
			Id:          "fake-user",
			LastLogin:   &now,
			LastUpdated: &now,
		})
		handledUserRequest = true
	}))

	defer ts.Close()

	oldHttpProxy := os.Getenv("OKTA_HTTP_PROXY")
	oldOrgName := os.Getenv("OKTA_ORG_NAME")
	oldApiToken := os.Getenv("OKTA_API_TOKEN")
	os.Setenv("OKTA_HTTP_PROXY", ts.URL)
	os.Setenv("OKTA_ORG_NAME", "unit-testing")
	os.Setenv("OKTA_API_TOKEN", "fake-token")
	t.Cleanup(func() {
		os.Setenv("OKTA_HTTP_PROXY", oldHttpProxy)
		os.Setenv("OKTA_ORG_NAME", oldOrgName)
		os.Setenv("OKTA_API_TOKEN", oldApiToken)
	})

	err := accPreCheck()
	if err != nil {
		t.Fatalf("Did not expect accPreCheck() to fail: %s", err)
	}

	c, err := oktaConfig()
	if err != nil {
		t.Fatalf("Did not expect oktaConfig() to fail: %s", err)
	}
	if c.httpProxy != ts.URL {
		t.Fatalf("Execpted httpProxy to be %q, got %q", ts.URL, c.httpProxy)
	}
	if !handledUserRequest {
		t.Fatal("Expected local server to handle user request, but it didn't")
	}
}

func TestProviderValidate(t *testing.T) {
	envKeys := []string{
		"OKTA_ACCESS_TOKEN",
		"OKTA_ALLOW_LONG_RUNNING_ACC_TEST",
		"OKTA_API_CLIENT_ID",
		"OKTA_API_PRIVATE_KEY",
		"OKTA_API_PRIVATE_KEY_ID",
		"OKTA_API_PRIVATE_KEY_IE",
		"OKTA_API_SCOPES",
		"OKTA_API_TOKEN",
		"OKTA_BASE_URL",
		"OKTA_DEFAULT",
		"OKTA_GROUP",
		"OKTA_HTTP_PROXY",
		"OKTA_ORG_NAME",
		"OKTA_UPDATE",
	}
	envVals := make(map[string]string)
	// save and clear OKTA env vars so config can be test cleanly
	for _, key := range envKeys {
		val := os.Getenv(key)
		if val == "" {
			continue
		}
		envVals[key] = val
		os.Unsetenv(key)
	}

	tests := []struct {
		name         string
		accessToken  string
		apiToken     string
		clientID     string
		privateKey   string
		privateKeyID string
		scopes       []interface{}
		expectError  bool
	}{
		{"simple pass", "", "", "", "", "", []interface{}{}, false},
		{"access_token pass", "accessToken", "", "", "", "", []interface{}{}, false},
		{"access_token fail 1", "accessToken", "apiToken", "", "", "", []interface{}{}, true},
		{"access_token fail 2", "accessToken", "", "cliendID", "", "", []interface{}{}, true},
		{"access_token fail 3", "accessToken", "", "", "privateKey", "", []interface{}{}, true},
		{"access_token fail 4", "accessToken", "", "", "", "", []interface{}{"scope1", "scope2"}, true},
		{"api_token pass", "", "apiToken", "", "", "", []interface{}{}, false},
		{"api_token fail 1", "accessToken", "apiToken", "", "", "", []interface{}{}, true},
		{"api_token fail 2", "", "apiToken", "clientID", "", "", []interface{}{}, true},
		{"api_token fail 3", "", "apiToken", "", "", "privateKey", []interface{}{}, true},
		{"api_token fail 4", "", "apiToken", "", "", "", []interface{}{"scope1", "scope2"}, true},
		{"client_id pass", "", "", "clientID", "", "", []interface{}{}, false},
		{"client_id fail 1", "accessToken", "", "clientID", "", "", []interface{}{}, true},
		{"client_id fail 2", "accessToken", "apiToken", "clientID", "", "", []interface{}{}, true},
		{"private_key pass", "", "", "", "privateKey", "", []interface{}{}, false},
		{"private_key fail 1", "accessToken", "", "", "privateKey", "", []interface{}{}, true},
		{"private_key fail 2", "", "apiToken", "", "privateKey", "", []interface{}{}, true},
		{"private_key_id pass", "", "", "", "", "privateKeyID", []interface{}{}, false},
		{"private_key_id fail 1", "", "apiToken", "", "", "privateKeyID", []interface{}{}, true},
		{"scopes pass", "", "", "", "", "", []interface{}{"scope1", "scope2"}, false},
		{"scopes fail 1", "accessToken", "", "", "", "", []interface{}{"scope1", "scope2"}, true},
		{"scopes fail 2", "", "apiToken", "", "", "", []interface{}{"scope1", "scope2"}, true},
	}

	for _, test := range tests {
		resourceConfig := map[string]interface{}{}
		if test.accessToken != "" {
			resourceConfig["access_token"] = test.accessToken
		}
		if test.apiToken != "" {
			resourceConfig["api_token"] = test.apiToken
		}
		if test.clientID != "" {
			resourceConfig["client_id"] = test.clientID
		}
		if test.privateKey != "" {
			resourceConfig["private_key"] = test.privateKey
		}
		if test.privateKeyID != "" {
			resourceConfig["private_key_id"] = test.privateKeyID
		}
		if len(test.scopes) > 0 {
			resourceConfig["scopes"] = test.scopes
		}

		config := terraform.NewResourceConfigRaw(resourceConfig)
		provider := Provider()
		err := provider.Validate(config)

		if test.expectError && err == nil {
			t.Errorf("test %q: expected error but received none", test.name)
		}
		if !test.expectError && err != nil {
			t.Errorf("test %q: did not expect error but received error: %+v", test.name, err)
			fmt.Println()
		}
	}

	for key, val := range envVals {
		os.Setenv(key, val)
	}
}

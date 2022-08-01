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
	"github.com/okta/okta-sdk-golang/v2/okta"
)

var (
	testAccProvidersFactories map[string]func() (*schema.Provider, error)
	testAccProvider           *schema.Provider
)

func init() {
	testAccProvider = Provider()
	testAccProvidersFactories = map[string]func() (*schema.Provider, error){
		"okta": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
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
		privateKeyId: 	os.Getenv("OKTA_API_PRIVATE_KEY_ID"),
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

func testAccPreCheck(t *testing.T) {
	err := accPreCheck()
	if err != nil {
		t.Fatalf("%v", err)
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

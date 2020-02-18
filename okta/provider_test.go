package okta

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"okta": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProviderImpl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func accPreCheck() error {
	if v := os.Getenv("OKTA_ORG_NAME"); v == "" {
		return errors.New("OKTA_ORG_NAME must be set for acceptance tests")
	}
	if v := os.Getenv("OKTA_API_TOKEN"); v == "" {
		return errors.New("OKTA_API_TOKEN must be set for acceptance tests")
	}

	return nil
}

func oktaConfig() (*Config, error) {
	config := &Config{
		orgName:     os.Getenv("OKTA_ORG_NAME"),
		apiToken:    os.Getenv("OKTA_API_TOKEN"),
		domain:      os.Getenv("OKTA_BASE_URL"),
		parallelism: 1,
		retryCount:  10,
		minWait:     30,
		maxWait:     600,
	}

	if err := config.loadAndValidate(); err != nil {
		return config, fmt.Errorf("Error initializing Okta client: %v", err)
	}

	return config, nil
}

func testAccPreCheck(t *testing.T) {
	err := accPreCheck()

	if err != nil {
		t.Fatalf("%v", err)
	}
}

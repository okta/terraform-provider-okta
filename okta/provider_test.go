package okta

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
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

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("OKTA_ORG_NAME"); v == "" {
		t.Fatal("OKTA_ORG_NAME must be set for acceptance tests")
	}
	if v := os.Getenv("OKTA_API_TOKEN"); v == "" {
		t.Fatal("OKTA_API_TOKEN must be set for acceptance tests")
	}
}

func testOktaConfig(t *testing.T) *Config {
	config := Config{
		orgName:  os.Getenv("OKTA_ORG_NAME"),
		apiToken: os.Getenv("OKTA_API_TOKEN"),
		domain:   os.Getenv("OKTA_BASE_URL"),
	}
	if err := config.loadAndValidate(); err != nil {
		t.Fatalf("Error initializing Okta client: %v", err)
	}
	return &config
}

func TestAccOktaProviderRegistration_articulateSDK(t *testing.T) {
	testAccPreCheck(t)
	c := testOktaConfig(t)

	// test credentials by listing our default user profile schema
	url := fmt.Sprintf("meta/schemas/user/default")

	req, err := c.articulateOktaClient.NewRequest("GET", url, nil)

	if err != nil {
		t.Fatalf("Error initializing test connection to Okta: %v", err)
	}
	_, err = c.articulateOktaClient.Do(req, nil)
	if err != nil {
		t.Fatalf("Error testing connection to Okta via the Articulate SDK. Please verify your credentials: %v", err)
	}
}

func TestAccOktaProviderRegistration_oktaSDK(t *testing.T) {
	testAccPreCheck(t)
	c := testOktaConfig(t)

  // test credentials by listing users in account
  // will limit to 1 user when this gets merged - https://github.com/okta/okta-sdk-golang/pull/28
	_, _, err := c.oktaClient.User.ListUsers(nil)

	if err != nil {
		t.Fatalf("Error testing connection to Okta via the Official Okta SDK. Please verify your credentials: %v", err)
	}
}

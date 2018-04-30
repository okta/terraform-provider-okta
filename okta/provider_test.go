package okta

import (
	"fmt"
	"os"
	"testing"

	"github.com/articulate/oktasdk-go/okta"
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
	testAccPreCheck(t)
	config := Config{
		orgName:  os.Getenv("OKTA_ORG_NAME"),
		apiToken: os.Getenv("OKTA_API_TOKEN"),
		domain:   os.Getenv("OKTA_BASE_URL"),
	}
	if err := config.loadAndValidate(); err != nil {
		t.Fatal("Error initializing Okta client: %v", err)
	}
	return &config
}

func TestAccOktaProviderRegistration(t *testing.T) {
	c := testOktaConfig(t)
	client, err := okta.NewClientWithDomain(nil, c.orgName, c.domain, c.apiToken)
	if err != nil {
		t.Fatalf("Error building Okta Client: %v", err)
	}
	// test credentials by listing our authorization servers
	url := fmt.Sprintf("meta/schemas/user/default")
	req, err := client.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("Error initializing test connection to Okta: %v", err)
	}
	_, err = client.Do(req, nil)
	if err != nil {
		t.Fatalf("Error testing connection to Okta. Please verify your credentials: %v", err)
	}
}

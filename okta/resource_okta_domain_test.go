package okta

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceOktaDomain(t *testing.T) {
	mgr := newFixtureManager("resources", domain, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", domain)
	domainName := "testacc.example.edu"
	updateDomainName := "testacctest.example.edu"

	oktaResourceTest(t, resource.TestCase{
		PreCheck:                 testAccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: testAccMergeProvidersFactories,
		CheckDestroy:             checkResourceDestroy(domain, domainExists),
		Steps: []resource.TestStep{
			{
				ExpectNonEmptyPlan: true,
				Config:             config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, domainExists),
					resource.TestCheckResourceAttr(resourceName, "name", domainName),
					resource.TestCheckResourceAttr(resourceName, "certificate_source_type", "MANUAL"),
					resource.TestCheckResourceAttr(resourceName, "validation_status", "FAILED_TO_VERIFY"),
					resource.TestCheckResourceAttr(resourceName, "dns_records.#", "2"),
				),
			},
			{
				ExpectNonEmptyPlan: true,
				Config:             updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, domainExists),
					resource.TestCheckResourceAttr(resourceName, "name", updateDomainName),
					resource.TestCheckResourceAttr(resourceName, "certificate_source_type", "OKTA_MANAGED"),
					resource.TestCheckResourceAttr(resourceName, "dns_records.#", "2"),
				),
			},
		},
	})
}

func domainExists(id string) (bool, error) {
	client := sdkV2ClientForTest()
	domain, resp, err := client.Domain.GetDomain(context.Background(), id)
	if err := suppressErrorOn404(resp, err); err != nil {
		return false, err
	}
	return domain != nil, nil
}

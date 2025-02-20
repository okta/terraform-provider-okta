package idaas_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/provider"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

func TestAccResourceOktaDomain_crud(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSDomain, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSDomain)
	domainName := "testacc.example.edu"
	updateDomainName := "testacctest.example.edu"

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.AccMergeProvidersFactoriesForTest(),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSDomain, domainExists),
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
	client := provider.SdkV2ClientForTest()
	domain, resp, err := client.Domain.GetDomain(context.Background(), id)
	if err := utils.SuppressErrorOn404(resp, err); err != nil {
		return false, err
	}
	return domain != nil, nil
}

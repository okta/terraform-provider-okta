package idaas_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

func TestAccResourceOktaDomain_crud(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSDomain, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSDomain)
	domainName := fmt.Sprintf("testacc-%d.example.com", mgr.Seed)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSDomain, domainExists),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, domainExists),
					resource.TestCheckResourceAttr(resourceName, "name", domainName),
					resource.TestCheckResourceAttr(resourceName, "certificate_source_type", "MANUAL"),
					resource.TestCheckResourceAttr(resourceName, "validation_status", "NOT_STARTED"),
					resource.TestCheckResourceAttr(resourceName, "dns_records.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "dns_records.1.record_type", "CNAME"),
					resource.TestCheckResourceAttr(resourceName, "dns_records.1.expiration", ""),
					resource.TestCheckResourceAttr(resourceName, "dns_records.1.fqdn", domainName),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, domainExists),
					resource.TestCheckResourceAttr(resourceName, "name", domainName),
					resource.TestCheckResourceAttr(resourceName, "certificate_source_type", "MANUAL"),
					resource.TestCheckResourceAttr(resourceName, "dns_records.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "dns_records.1.record_type", "CNAME"),
					resource.TestCheckResourceAttr(resourceName, "dns_records.1.expiration", ""),
					resource.TestCheckResourceAttr(resourceName, "dns_records.1.fqdn", domainName),
				),
			},
		},
	})
}

func domainExists(id string) (bool, error) {
	client := iDaaSAPIClientForTestUtil.OktaSDKClientV5()
	domain, resp, err := client.CustomDomainAPI.GetCustomDomain(context.Background(), id).Execute()
	if err := utils.SuppressErrorOn404_V5(resp, err); err != nil {
		return false, err
	}
	return domain != nil, nil
}

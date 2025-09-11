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

func TestAccResourceOktaEmailDomain_crud(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSEmailDomain, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSEmailDomain)
	domainName := fmt.Sprintf("testAcc-%d.example.com", mgr.Seed)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSEmailDomain, emailDomainExists),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, emailDomainExists),
					resource.TestCheckResourceAttrSet(resourceName, "brand_id"),
					resource.TestCheckResourceAttr(resourceName, "domain", domainName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "test"),
					resource.TestCheckResourceAttr(resourceName, "user_name", "fff"),
					resource.TestCheckResourceAttr(resourceName, "validation_subdomain", "mail"),
					resource.TestCheckResourceAttrSet(resourceName, "dns_validation_records.0.record_type"),
					resource.TestCheckResourceAttrSet(resourceName, "dns_validation_records.0.value"),
					resource.TestCheckResourceAttrSet(resourceName, "dns_validation_records.0.fqdn"),
				),
			},
		},
	})
}

func emailDomainExists(id string) (bool, error) {
	client := iDaaSAPIClientForTestUtil.OktaSDKClientV3()
	emailDomain, resp, err := client.EmailDomainAPI.GetEmailDomain(context.Background(), id).Execute()
	if err := utils.SuppressErrorOn404_V3(resp, err); err != nil {
		return false, err
	}
	return emailDomain != nil && emailDomain.GetValidationStatus() != "DELETED", nil
}

package okta

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceOktaEmailDomain(t *testing.T) {
	mgr := newFixtureManager(emailDomain, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	resourceName := fmt.Sprintf("%s.test", emailDomain)
	domainName := fmt.Sprintf("testAcc-%d.example.com", mgr.Seed)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkResourceDestroy(emailDomain, emailDomainExists),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, emailDomainExists),
					resource.TestCheckResourceAttrSet(resourceName, "brand_id"),
					resource.TestCheckResourceAttr(resourceName, "domain", domainName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "test"),
					resource.TestCheckResourceAttr(resourceName, "user_name", "fff"),
					resource.TestCheckResourceAttrSet(resourceName, "dns_validation_records.0.record_type"),
					resource.TestCheckResourceAttrSet(resourceName, "dns_validation_records.0.value"),
					resource.TestCheckResourceAttrSet(resourceName, "dns_validation_records.0.fqdn"),
				),
			},
		},
	})
}

func emailDomainExists(id string) (bool, error) {
	client := sdkV3ClientForTest()
	emailDomain, resp, err := client.EmailDomainApi.GetEmailDomain(context.Background(), id).Execute()
	if err := v3suppressErrorOn404(resp, err); err != nil {
		return false, err
	}
	return emailDomain != nil && emailDomain.GetValidationStatus() != "DELETED", nil
}

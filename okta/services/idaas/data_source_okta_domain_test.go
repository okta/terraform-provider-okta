package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaDomain_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSDomain, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	testAccDomain := fmt.Sprintf("testacc-%d.example.com", mgr.Seed)
	testAccDowncaseDomain := fmt.Sprintf("downcase-testacc-%d.example.com", mgr.Seed)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config:  config,
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					// Note Okta API down cases DNS names
					resource.TestCheckResourceAttr("data.okta_domain.by-id", "domain", testAccDomain),
					resource.TestCheckResourceAttr("data.okta_domain.by-name", "domain", testAccDomain),
					resource.TestCheckResourceAttr("data.okta_domain.by-id-downcase", "domain", testAccDowncaseDomain),
					resource.TestCheckResourceAttr("data.okta_domain.by-name-downcase", "domain", testAccDowncaseDomain),
				),
			},
		},
	})
}

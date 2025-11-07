package idaas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaDomain_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSDomain, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config:  config,
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.okta_domain.by-id-downcase", "domain", "testdowncase.example.com"),
					resource.TestCheckResourceAttr("data.okta_domain.by-id-downcase", "dns_records.1.record_type", "CNAME"),
					resource.TestCheckResourceAttr("data.okta_domain.by-id-downcase", "dns_records.1.expiration", ""),
					resource.TestCheckResourceAttr("data.okta_domain.by-id-downcase", "dns_records.1.fqdn", "testdowncase.example.com"),
				),
			},
		},
	})
}

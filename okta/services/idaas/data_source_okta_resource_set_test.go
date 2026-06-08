package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaResourceSet_read(t *testing.T) {
	datasourceName := fmt.Sprintf("data.%s.test", resources.OktaIDaaSResourceSet)

	mgr := newFixtureManager("data-sources", resources.OktaIDaaSResourceSet, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "id"),
					resource.TestCheckResourceAttr(datasourceName, "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(datasourceName, "description", "test resource set for data source"),
					resource.TestCheckResourceAttrSet(datasourceName, "resources.#"),
				),
				// VCR recording rewrites hostnames in API responses, causing
				// drift between the config hostname and the state hostname for
				// okta_resource_set resources. During VCR playback, hostnames
				// are consistent so the plan is empty.
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

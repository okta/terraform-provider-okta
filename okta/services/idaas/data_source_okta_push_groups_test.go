package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaPushGroups_read(t *testing.T) {
	datasourceNameSample := fmt.Sprintf("data.%s.sample", resources.OktaIDaaSPushGroups)

	mgr := newFixtureManager("data-sources", resources.OktaIDaaSPushGroups, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameSample, "app_id"),
					resource.TestCheckResourceAttrSet(datasourceNameSample, "id"),
				),
			},
		},
	})
}

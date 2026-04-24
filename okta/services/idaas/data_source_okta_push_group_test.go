package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaPushGroup_read(t *testing.T) {
	datasourceNameSample := fmt.Sprintf("data.%s.sample", resources.OktaIDaaSPushGroup)
	datasourceNameSampleTwo := fmt.Sprintf("data.%s.sample_two", resources.OktaIDaaSPushGroup)

	mgr := newFixtureManager("data-sources", resources.OktaIDaaSPushGroup, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameSample, "id"),
					resource.TestCheckResourceAttrSet(datasourceNameSample, "app_id"),
					resource.TestCheckResourceAttrSet(datasourceNameSample, "source_group_id"),
					resource.TestCheckResourceAttrSet(datasourceNameSample, "target_group_id"),
					resource.TestCheckResourceAttr(datasourceNameSample, "status", "ACTIVE"),
				),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameSampleTwo, "id"),
					resource.TestCheckResourceAttrSet(datasourceNameSampleTwo, "app_id"),
					resource.TestCheckResourceAttrSet(datasourceNameSampleTwo, "source_group_id"),
					resource.TestCheckResourceAttrSet(datasourceNameSampleTwo, "target_group_id"),
					resource.TestCheckResourceAttr(datasourceNameSampleTwo, "status", "ACTIVE"),
				),
			},
		},
	})
}

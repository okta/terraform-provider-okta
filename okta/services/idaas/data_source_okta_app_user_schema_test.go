package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaAppUserSchema_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSAppUserSchema, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)
	resourceName := fmt.Sprintf("data.%s.test", resources.OktaIDaaSAppUserSchema)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "app_id"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "custom_property.*", map[string]string{
						"index":       "testCustomProp1",
						"title":       "Test Custom Property 1",
						"type":        "string",
						"description": "Test description 1",
					}),
				),
			},
		},
	})
}

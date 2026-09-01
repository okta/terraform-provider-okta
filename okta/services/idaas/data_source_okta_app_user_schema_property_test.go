package idaas_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaAppUserSchemaProperty_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSAppUserSchemaProperty, t.Name())
	resourceName := fmt.Sprintf("data.%s.test", resources.OktaIDaaSAppUserSchemaProperty)
	config := mgr.GetFixtures("basic.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "index", "testCustomProperty"),
					resource.TestCheckResourceAttr(resourceName, "title", "Test Custom Property"),
					resource.TestCheckResourceAttr(resourceName, "type", "string"),
					resource.TestCheckResourceAttr(resourceName, "description", "Test description"),
					resource.TestCheckResourceAttr(resourceName, "permissions", "READ_ONLY"),
					resource.TestCheckResourceAttr(resourceName, "master", "PROFILE_MASTER"),
				),
			},
		},
	})
}

func TestAccDataSourceOktaAppUserSchemaProperty_notFound(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSAppUserSchemaProperty, t.Name())
	config := mgr.GetFixtures("not_found.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("application user schema property with index 'nonExistentProperty' not found"),
			},
		},
	})
}

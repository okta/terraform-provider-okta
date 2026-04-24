package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceHookKey_read(t *testing.T) {
	hookKeyName := "test_hook_key"
	dataSourceName := fmt.Sprintf("data.%s.test", resources.OktaIDaaSHookKey)
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSHookKey, t.Name())
	config := mgr.GetFixtures("data-source.tf", t)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", hookKeyName),
					resource.TestCheckResourceAttr(dataSourceName, "id", "abcdefghij0123456789"),
					resource.TestCheckResourceAttr(dataSourceName, "key_id", "074497ab-411a-44fd-b84d-676e1f6cb3c7"),
					resource.TestCheckResourceAttrSet(dataSourceName, "created"),
					resource.TestCheckResourceAttrSet(dataSourceName, "last_updated"),
					resource.TestCheckResourceAttr(dataSourceName, "is_used", "false"),
				),
			},
		},
	})
}

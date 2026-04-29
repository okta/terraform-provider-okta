package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccResourceOktaIdentitySourceBulkUpsert_basic(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSIdentitySourceBulkUpsert, t.Name())
	config := mgr.GetFixtures("resource.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSIdentitySourceBulkUpsert)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "identity_source_id", "0oaxc95befZNgrJl71d7"),
					resource.TestCheckResourceAttr(resourceName, "entity_type", "USERS"),
				),
			},
		},
	})
}

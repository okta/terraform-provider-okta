package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccResourceOktaAppTokenResource_basic(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppToken, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	resourceName := fmt.Sprintf("%s.example", resources.OktaIDaaSAppToken)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				ImportState:        true,
				ResourceName:       "okta_app_token.example",
				ImportStateId:      "0oardd5r32PWsF4421d7/oar1gyu6hmmw8bj6I1d7",
				ImportStatePersist: true,
				Config:             config,
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "client_id", "0oardd5r32PWsF4421d7"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "user_id", "00unkw1sfbTw08c0g1d7"),
				),
			},
		},
	})
}

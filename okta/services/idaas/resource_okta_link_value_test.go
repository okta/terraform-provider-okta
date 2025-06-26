package idaas_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccResourceOktaLinkValue_crud(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSLinkValue, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updated := mgr.GetFixtures("updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSLinkValue)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkLinkValueDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "associated_user_ids.#", "4"),
				),
			},
			{
				Config: updated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "associated_user_ids.#", "1"),
				),
			},
		},
	})
}

func checkLinkValueDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != resources.OktaIDaaSLinkValue {
			continue
		}
		client := iDaaSAPIClientForTestUtil.OktaSDKClientV2()
		lo, resp, err := client.LinkedObject.GetLinkedObjectDefinition(context.Background(), rs.Primary.Attributes["primary_name"])
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil
		} else if err != nil {
			return err
		}
		puID := rs.Primary.Attributes["primary_user_id"]
		los, resp, err := client.User.GetLinkedObjectsForUser(context.Background(), puID, lo.Associated.Name, nil)
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil
		} else if err != nil {
			return err
		}
		if len(los) == 0 {
			return nil
		}
		return fmt.Errorf("there are still %d relationships exists between primary and associated users", len(los))
	}
	return nil
}

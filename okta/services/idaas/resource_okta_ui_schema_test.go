package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccResourceOktaUISchema_crud(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSUISchema, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updated := mgr.GetFixtures("updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSUISchema)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ui_schema.type", "Group"),
					resource.TestCheckResourceAttr(resourceName, "ui_schema.button_label", "submit"),
					resource.TestCheckResourceAttr(resourceName, "ui_schema.label", "Sign in"),
					resource.TestCheckResourceAttr(resourceName, "ui_schema.elements.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ui_schema.elements.0.type", "Control"),
					resource.TestCheckResourceAttr(resourceName, "ui_schema.elements.0.label", "Last Name"),
				),
			},
			{
				Config: updated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ui_schema.elements.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ui_schema.elements.0.type", "Control"),
					resource.TestCheckResourceAttr(resourceName, "ui_schema.elements.0.label", "Last Name"),
					resource.TestCheckResourceAttr(resourceName, "ui_schema.elements.1.type", "Control"),
					resource.TestCheckResourceAttr(resourceName, "ui_schema.elements.1.label", "First Name"),
				),
			},
		},
	})
}

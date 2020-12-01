package okta

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaEmailTemplate_crud(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := fmt.Sprintf("%s.test", templateEmail)
	mgr := newFixtureManager(templateEmail)
	config := mgr.GetFixtures("basic.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(templateEmail, doesEmailTemplateExist),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "email.forgotPassword"),
				),
			},
		},
	})
}

func doesEmailTemplateExist(id string) (bool, error) {
	_, response, err := getSupplementFromMetadata(testAccProvider.Meta()).GetEmailTemplate(context.Background(), id)
	return doesResourceExist(response, err)
}

package okta

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaEmailTemplate_crud(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", templateEmail)
	mgr := newFixtureManager(templateEmail, t.Name())
	config := mgr.GetFixtures("basic.tf", t)

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
	templ, _, err := getSupplementFromMetadata(testAccProvider.Meta()).GetEmailTemplate(context.Background(), id)
	if err != nil {
		return false, err
	}
	if templ == nil || templ.Id == "" || templ.Id == "default" {
		return false, nil
	}
	return true, err
}

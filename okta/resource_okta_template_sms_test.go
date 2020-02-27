package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOktaSmsTemplate_crud(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := fmt.Sprintf("%s.test", templateSms)
	mgr := newFixtureManager(templateSms)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updated := mgr.GetFixtures("updated.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(templateSms, doesSmsTemplateExist),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "SMS_VERIFY_CODE"),
					resource.TestCheckResourceAttr(resourceName, "template", "Your ${org.name} code is: ${code}"),
				),
			},
			{
				Config: updated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "SMS_VERIFY_CODE"),
					resource.TestCheckResourceAttr(resourceName, "template", "Your ${org.name} updated code is: ${code}"),
				),
			},
		},
	})
}

func doesSmsTemplateExist(id string) (bool, error) {
	client := getSupplementFromMetadata(testAccProvider.Meta())
	_, response, err := client.GetSmsTemplate(id)

	return doesResourceExist(response, err)
}

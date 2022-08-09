package okta

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaSmsTemplate_crud(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", templateSms)
	mgr := newFixtureManager(templateSms, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updated := mgr.GetFixtures("updated.tf", t)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(templateSms, doesSmsTemplateExist),
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
	client := oktaClientForTest()
	_, response, err := client.SmsTemplate.GetSmsTemplate(context.Background(), id)
	return doesResourceExist(response, err)
}

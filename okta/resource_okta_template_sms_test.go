package okta

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceOktaSmsTemplate_crud(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", templateSms)
	mgr := newFixtureManager("resources", templateSms, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updated := mgr.GetFixtures("updated.tf", t)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkResourceDestroy(templateSms, doesSmsTemplateExist),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "SMS_VERIFY_CODE"),
					resource.TestCheckResourceAttr(resourceName, "template", "${org.name} code is: ${code}"),
					resource.TestCheckResourceAttr(resourceName, "template", "${org.name} code is: ${code}"),
					resource.TestCheckResourceAttr(resourceName, "translations.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "translations.0.language", "en"),
					resource.TestCheckResourceAttr(resourceName, "translations.1.language", "es"),
				),
			},
			{
				Config: updated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "SMS_VERIFY_CODE"),
					resource.TestCheckResourceAttr(resourceName, "template", "${org.name} updated code is: ${code}"),
					resource.TestCheckResourceAttr(resourceName, "translations.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "translations.0.language", "en"),
					resource.TestCheckResourceAttr(resourceName, "translations.1.language", "es"),
					resource.TestCheckResourceAttr(resourceName, "translations.2.language", "fr"),
				),
			},
		},
	})
}

func doesSmsTemplateExist(id string) (bool, error) {
	client := sdkV2ClientForTest()
	_, response, err := client.SmsTemplate.GetSmsTemplate(context.Background(), id)
	return doesResourceExist(response, err)
}

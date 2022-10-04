package okta

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaSmsTemplate_crud(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := fmt.Sprintf("%s.test", templateSms)
	mgr := newFixtureManager(templateSms)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updated := mgr.GetFixtures("updated.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(templateSms, doesSmsTemplateExist),
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
	_, response, err := getOktaClientFromMetadata(testAccProvider.Meta()).SmsTemplate.GetSmsTemplate(context.Background(), id)
	return doesResourceExist(response, err)
}

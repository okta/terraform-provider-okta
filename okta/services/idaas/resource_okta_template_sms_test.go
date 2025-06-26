package idaas_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

func TestAccResourceOktaSmsTemplate_crud(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSTemplateSms)
	mgr := newFixtureManager("resources", resources.OktaIDaaSTemplateSms, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updated := mgr.GetFixtures("updated.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSTemplateSms, doesSmsTemplateExist),
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
	client := iDaaSAPIClientForTestUtil.OktaSDKClientV2()
	_, response, err := client.SmsTemplate.GetSmsTemplate(context.Background(), id)
	return utils.DoesResourceExist(response, err)
}

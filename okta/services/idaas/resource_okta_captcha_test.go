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

func TestAccResourceOktaCaptcha_crud(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSCaptcha, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updated := mgr.GetFixtures("updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSCaptcha)
	acctest.OktaResourceTest(
		t, resource.TestCase{
			PreCheck:                 acctest.AccPreCheck(t),
			ErrorCheck:               testAccErrorChecks(t),
			ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
			CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSCaptcha, doesCaptchaExist),
			Steps: []resource.TestStep{
				{
					Config: config,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
						resource.TestCheckResourceAttr(resourceName, "type", "HCAPTCHA"),
						resource.TestCheckResourceAttr(resourceName, "site_key", "random_key"),
					),
				},
				{
					Config: updated,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)+"_updated"),
						resource.TestCheckResourceAttr(resourceName, "type", "HCAPTCHA"),
						resource.TestCheckResourceAttr(resourceName, "site_key", "random_key_updated")),
				},
			},
		})
}

func doesCaptchaExist(id string) (bool, error) {
	client := iDaaSAPIClientForTestUtil.OktaSDKSupplementClient()
	_, response, err := client.GetCaptcha(context.Background(), id)
	return utils.DoesResourceExist(response, err)
}

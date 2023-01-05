package okta

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func TestAccAppWSFedAppSettings_crud(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appWSFedAppSettings)
	preconfigured := mgr.GetFixtures("preconfigured.tf", ri, t)
	updated := mgr.GetFixtures("preconfigured_updated.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", appWSFedAppSettings)

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(appWSFed, createDoesAppExist(okta.NewWsFederationApplication())),
		Steps: []resource.TestStep{
			{
				Config: preconfigured,
				Check: resource.ComposeTestCheckFunc(
					checkAppWSFedAppSettingsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "settings", "{\"appFilter\":\"okta\",\"awsEnvironmentType\":\"aws.amazon\",\"groupFilter\":\"aws_(?{{accountid}}\\\\\\\\d+)_(?{{role}}[a-zA-Z0-9+=,.@\\\\\\\\-_]+)\",\"joinAllRoles\":false,\"loginURL\":\"https://console.aws.amazon.com/ec2/home\",\"roleValuePattern\":\"arn:aws:iam::${accountid}:saml-provider/OKTA,arn:aws:iam::${accountid}:role/${role}\",\"sessionDuration\":7600,\"useGroupMapping\":false}"),
				),
			},
			{
				Config: updated,
				Check: resource.ComposeTestCheckFunc(
					checkAppWSFedAppSettingsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "settings", "{\"appFilter\":\"okta\",\"awsEnvironmentType\":\"aws.amazon\",\"groupFilter\":\"aws_(?{{accountid}}\\\\\\\\d+)_(?{{role}}[a-zA-Z0-9+=,.@\\\\\\\\-_]+)\",\"joinAllRoles\":false,\"loginURL\":\"https://console.aws.amazon.com/ec2/home\",\"roleValuePattern\":\"arn:aws:iam::${accountid}:saml-provider/OKTA,arn:aws:iam::${accountid}:role/${role}\",\"sessionDuration\":3200,\"useGroupMapping\":false}"),
				),
			},
		},
	})
}

func checkAppWSFedAppSettingsExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		missingErr := fmt.Errorf("resource not found: %s", name)
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return missingErr
		}
		appID := rs.Primary.Attributes["app_id"]
		client := getOktaClientFromMetadata(testAccProvider.Meta())
		app := okta.NewWsFederationApplication()
		_, _, err := client.Application.GetApplication(context.Background(), appID, app, nil)
		if err != nil {
			return err
		}
		settings := make(okta.ApplicationSettingsApplication)
		_ = json.Unmarshal([]byte(rs.Primary.Attributes["settings"]), &settings)
		for k, v := range *app {
			if v == nil {
				delete(*app.Settings.App, k)
			}
		}
		e := reflect.DeepEqual(*app.Settings.App, settings)
		if !e {
			return fmt.Errorf("settings are not equal: actual: %+v , expected: %+v", *app.Settings.App, settings)
		}
		return nil
	}
}

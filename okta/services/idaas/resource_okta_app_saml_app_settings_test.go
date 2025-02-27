package idaas_test

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/sdk"
)

func TestAccResourceOktaAppSamlAppSettings_crud(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppSamlAppSettings, t.Name())
	preconfigured := mgr.GetFixtures("preconfigured.tf", t)
	updated := mgr.GetFixtures("preconfigured_updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppSamlAppSettings)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSAppSaml, createDoesAppExist(sdk.NewSamlApplication())),
		Steps: []resource.TestStep{
			{
				Config: preconfigured,
				Check: resource.ComposeTestCheckFunc(
					checkAppSamlAppSettingsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "settings", "{\"appFilter\":\"okta\",\"awsEnvironmentType\":\"aws.amazon\",\"groupFilter\":\"aws_(?{{accountid}}\\\\\\\\d+)_(?{{role}}[a-zA-Z0-9+=,.@\\\\\\\\-_]+)\",\"joinAllRoles\":false,\"loginURL\":\"https://console.aws.amazon.com/ec2/home\",\"roleValuePattern\":\"arn:aws:iam::${accountid}:saml-provider/OKTA,arn:aws:iam::${accountid}:role/${role}\",\"sessionDuration\":7600,\"useGroupMapping\":false}"),
				),
			},
			{
				Config: updated,
				Check: resource.ComposeTestCheckFunc(
					checkAppSamlAppSettingsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "settings", "{\"appFilter\":\"okta\",\"awsEnvironmentType\":\"aws.amazon\",\"groupFilter\":\"aws_(?{{accountid}}\\\\\\\\d+)_(?{{role}}[a-zA-Z0-9+=,.@\\\\\\\\-_]+)\",\"joinAllRoles\":false,\"loginURL\":\"https://console.aws.amazon.com/ec2/home\",\"roleValuePattern\":\"arn:aws:iam::${accountid}:saml-provider/OKTA,arn:aws:iam::${accountid}:role/${role}\",\"sessionDuration\":3200,\"useGroupMapping\":false}"),
				),
			},
		},
	})
}

func checkAppSamlAppSettingsExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		missingErr := fmt.Errorf("resource not found: %s", resourceName)
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return missingErr
		}
		appID := rs.Primary.Attributes["app_id"]
		client := iDaaSAPIClientForTestUtil.OktaSDKClientV2()
		app := sdk.NewSamlApplication()
		_, _, err := client.Application.GetApplication(context.Background(), appID, app, nil)
		if err != nil {
			return err
		}
		settings := make(sdk.ApplicationSettingsApplication)
		_ = json.Unmarshal([]byte(rs.Primary.Attributes["settings"]), &settings)
		for k, v := range *app.Settings.App {
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

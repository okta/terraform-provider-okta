package idaas_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/services/idaas"
	"github.com/okta/terraform-provider-okta/sdk"
)

// TestAccResourceOktaAppUserSchema_drift_detection
//
// Okta can add new app user schema properties outside of Terraform (e.g. when enabling provisioning).
// This test ensures we surface that as drift (non-empty plan), rather than silently reconciling it away.
func TestAccResourceOktaAppUserSchema_drift_detection(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppUserSchema, t.Name())
	baseConfig := mgr.GetFixtures("drift_base.tf", t)
	expectedConfig := mgr.GetFixtures("drift_expected.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppUserSchema)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSAppUserSchema, appUserSchemaExists),
		Steps: []resource.TestStep{
			{
				Config: baseConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, appUserSchemaExists),
					resource.TestCheckResourceAttr(resourceName, "custom_property.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "custom_property.*", map[string]string{
						"index":       "testCustomProp1",
						"title":       "Test Custom Property 1",
						"type":        "string",
						"description": "Test description 1",
						"permissions": "READ_ONLY",
						"scope":       "NONE",
					}),
				),
			},
			{
				// Re-apply the same config, but mutate Okta in-between.
				Config: baseConfig,
				Check: resource.ComposeTestCheckFunc(
					clickOpsAddAppUserSchemaProperty(resourceName, "autoAddedProp"),
				),
				// The test runner will refresh after Check funcs; the new schema property should be detected as drift.
				ExpectNonEmptyPlan: true,
			},
			{
				// Once we declare the auto-added property in Terraform config, drift should go away.
				Config: expectedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, appUserSchemaExists),
					resource.TestCheckResourceAttr(resourceName, "custom_property.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "custom_property.*", map[string]string{
						"index":       "autoAddedProp",
						"title":       "Auto Added Property",
						"type":        "string",
						"description": "Auto added by Okta",
						"permissions": "READ_ONLY",
						"scope":       "NONE",
					}),
				),
			},
		},
	})
}

func clickOpsAddAppUserSchemaProperty(appUserSchemaResourceName, index string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// During VCR playback, do not attempt real "click ops" mutations.
		// The cassette already represents the post-mutation Okta state.
		if os.Getenv("OKTA_VCR_TF_ACC") == "play" {
			return nil
		}

		rs, ok := s.RootModule().Resources[appUserSchemaResourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", appUserSchemaResourceName)
		}

		appID := rs.Primary.ID
		if appID == "" {
			return fmt.Errorf("resource %s has empty ID", appUserSchemaResourceName)
		}

		required := false
		attr := &sdk.UserSchemaAttribute{
			Title:       "Auto Added Property",
			Type:        "string",
			Description: "Auto added by Okta",
			Required:    &required,
			Scope:       "NONE",
			Master: &sdk.UserSchemaAttributeMaster{
				Type: "PROFILE_MASTER",
			},
			Permissions: []*sdk.UserSchemaAttributePermission{
				{
					Action:    "READ_ONLY",
					Principal: "SELF",
				},
			},
			Union: "DISABLE",
		}

		custom := idaas.BuildCustomUserSchema(index, attr)
		client := iDaaSAPIClientForTestUtil.OktaSDKClientV2()
		if _, _, err := client.UserSchema.UpdateApplicationUserProfile(context.Background(), appID, *custom); err != nil {
			return fmt.Errorf("API: unable to add app user schema property %q to app %q: %v", index, appID, err)
		}

		// Wait for eventual consistency: ensure the property shows up in reads before returning.
		for i := 0; i < 15; i++ {
			us, _, err := client.UserSchema.GetApplicationUserSchema(context.Background(), appID)
			if err == nil && idaas.UserSchemaCustomAttribute(us, index) != nil {
				return nil
			}
			// lintignore:R018
			time.Sleep(1 * time.Second)
		}

		return fmt.Errorf("API: app user schema property %q did not appear for app %q after waiting", index, appID)
	}
}

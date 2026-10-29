package idaas_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccResourceOktaAppUserSchema_crud(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppUserSchema, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppUserSchema)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSAppUserSchema, appUserSchemaExists),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, appUserSchemaExists),
					resource.TestCheckResourceAttrSet(resourceName, "app_id"),
					resource.TestCheckResourceAttr(resourceName, "custom_property.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "custom_property.*", map[string]string{
						"index":       "testCustomProp1",
						"title":       "Test Custom Property 1",
						"type":        "string",
						"description": "Test description 1",
						"required":    "false",
						"permissions": "READ_ONLY",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "custom_property.*", map[string]string{
						"index":       "testCustomProp2",
						"title":       "Test Custom Property 2",
						"type":        "string",
						"description": "Test description 2",
						"required":    "true",
						"scope":       "SELF",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "custom_property.*", map[string]string{
						"index":       "testArrayProp",
						"title":       "Test Array Property",
						"type":        "array",
						"array_type":  "string",
						"description": "Test array description",
					}),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, appUserSchemaExists),
					resource.TestCheckResourceAttrSet(resourceName, "app_id"),
					resource.TestCheckResourceAttr(resourceName, "custom_property.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "custom_property.*", map[string]string{
						"index":       "testCustomProp1",
						"title":       "Test Custom Property 1",
						"type":        "string",
						"description": "Test description 1 updated",
						"required":    "true",
						"permissions": "READ_WRITE",
					}),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					if len(s) != 1 {
						return errors.New("failed to import schema into state")
					}
					return nil
				},
			},
		},
	})
}

func appUserSchemaExists(id string) (bool, error) {
	client := iDaaSAPIClientForTestUtil.OktaSDKClientV2()
	_, resp, err := client.UserSchema.GetApplicationUserSchema(context.Background(), id)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("failed to get application user schema: %v", err)
	}
	return true, nil
}

package okta

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceOktaProfileMapping_crud(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", profileMapping)
	mgr := newFixtureManager("resources", profileMapping, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("updated.tf", t)
	preventDelete := mgr.GetFixtures("prevent_delete.tf", t)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkResourceDestroy(profileMapping, doesOktaProfileExist),
		Steps: []resource.TestStep{
			{
				Config: preventDelete,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "delete_when_absent", "false"),
				),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "delete_when_absent", "true"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "delete_when_absent", "true"),
				),
			},
		},
	})
}

func TestAccResourceOktaProfileMapping_existing(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", profileMapping)
	config := `
	resource "okta_profile_mapping" "test" {
		source_id          = okta_idp_social.google.id
		target_id          = data.okta_user_profile_mapping_source.user.id
		delete_when_absent = true
	  
		mappings {
		  id         = "firstName"
		  expression = "appuser.firstName"
		}
	  
		mappings {
		  id         = "lastName"
		  expression = "appuser.lastName"
		}
	  
		mappings {
		  id         = "email"
		  expression = "appuser.email"
		}
	  
		mappings {
		  id         = "login"
		  expression = "appuser.email"
		}
	  }
	  
	  resource "okta_idp_social" "google" {
		type          = "GOOGLE"
		protocol_type = "OIDC"
		name          = "testAcc_google_replace_with_uuid"
	  
		scopes = [
		  "profile",
		  "email",
		  "openid",
		]
	  
		client_id         = "abcd123"
		client_secret     = "abcd123"
		username_template = "idpuser.email"
	  }
	  
	  data "okta_user_profile_mapping_source" "user" {
		depends_on = [okta_idp_social.google]
	  }
	  
	`
	oktaResourceTest(t, resource.TestCase{
		ProtoV5ProviderFactories: testAccMergeProvidersFactories,
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "source_id"),
					resource.TestCheckResourceAttrSet(resourceName, "target_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
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

// TODO deprecated endpoint
func doesOktaProfileExist(profileID string) (bool, error) {
	client := sdkSupplementClientForTest()
	_, response, err := client.GetEmailTemplate(context.Background(), profileID)
	return doesResourceExist(response, err)
}

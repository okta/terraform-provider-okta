package idaas_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/provider"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

func TestAccResourceOktaProfileMapping_crud(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSProfileMapping)
	mgr := newFixtureManager("resources", resources.OktaIDaaSProfileMapping, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("updated.tf", t)
	preventDelete := mgr.GetFixtures("prevent_delete.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:          acctest.AccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: acctest.AccProvidersFactoriesForTest(),
		CheckDestroy:      checkResourceDestroy(resources.OktaIDaaSProfileMapping, doesOktaProfileExist),
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
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSProfileMapping)
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
	acctest.OktaResourceTest(t, resource.TestCase{
		ProtoV5ProviderFactories: acctest.AccMergeProvidersFactoriesForTest(),
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
	client := provider.SdkSupplementClientForTest()
	_, response, err := client.GetEmailTemplate(context.Background(), profileID)
	return utils.DoesResourceExist(response, err)
}

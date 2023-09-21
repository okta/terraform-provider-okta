package okta

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/sdk"
	"github.com/okta/terraform-provider-okta/sdk/query"
)

func TestAccResourceOktaAppOAuthRoleAssignment_basic(t *testing.T) {
	mgr := newFixtureManager("okta_app_oauth_role_assignment", t.Name())

	oktaResourceTest(t, resource.TestCase{
		PreCheck:                 testAccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: testAccMergeProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: mgr.GetFixtures("basic.tf"),
			},
			{
				Config: mgr.GetFixtures("basic_updated.tf"),
			}
		}
	})
}

func TestAccResourceOktaAppOAuthRoleAssignment_custom(t *testing.T) {
	mgr := newFixtureManager("okta_app_oauth_role_assignment", t.Name())

	oktaResourceTest(t, resource.TestCase{
		PreCheck:                 testAccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: testAccMergeProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: mgr.GetFixtures("custom.tf"),
			},
			{
				Config: mgr.GetFixtures("custom_updated.tf"),
			}
		}
	})
}

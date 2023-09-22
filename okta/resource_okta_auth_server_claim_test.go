package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceOktaAuthServerClaim_create(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", authServerClaim)
	mgr := newFixtureManager("", authServerClaim, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkResourceDestroy(authServer, authServerExists),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "name", "test"),
					resource.TestCheckResourceAttr(resourceName, "value_type", "EXPRESSION"),
					resource.TestCheckResourceAttr(resourceName, "value", "cool"),
					resource.TestCheckResourceAttr(resourceName, "claim_type", "RESOURCE"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", statusInactive),
					resource.TestCheckResourceAttr(resourceName, "name", "test_updated"),
					resource.TestCheckResourceAttr(resourceName, "value_type", "EXPRESSION"),
					resource.TestCheckResourceAttr(resourceName, "value", "cool_updated"),
					resource.TestCheckResourceAttr(resourceName, "claim_type", "RESOURCE"),
				),
			},
		},
	})
}

func TestAccResourceOktaAuthServerClaim_groupType(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", authServerClaim)
	swResourceName := fmt.Sprintf("%s.test_sw", authServerClaim)
	mgr := newFixtureManager("", authServerClaim, t.Name())
	config := mgr.GetFixtures("basic_group.tf", t)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkResourceDestroy(authServer, authServerExists),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "name", "test"),
					resource.TestCheckResourceAttr(resourceName, "group_filter_type", "EQUALS"),
					resource.TestCheckResourceAttr(resourceName, "value_type", "GROUPS"),
					resource.TestCheckResourceAttr(resourceName, "value", groupProfileEveryone),
					resource.TestCheckResourceAttr(resourceName, "claim_type", "RESOURCE"),

					resource.TestCheckResourceAttr(swResourceName, "status", statusActive),
					resource.TestCheckResourceAttr(swResourceName, "name", "test_sw"),
					resource.TestCheckResourceAttr(swResourceName, "group_filter_type", "STARTS_WITH"),
					resource.TestCheckResourceAttr(swResourceName, "value_type", "GROUPS"),
					resource.TestCheckResourceAttr(swResourceName, "value", "Every"),
					resource.TestCheckResourceAttr(swResourceName, "claim_type", "RESOURCE"),
				),
			},
		},
	})
}

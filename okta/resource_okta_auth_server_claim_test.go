package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOktaAuthServerClaim_create(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := fmt.Sprintf("%s.test", authServerClaim)
	mgr := newFixtureManager(authServerClaim)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(authServer, authServerExists),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "name", "test"),
					resource.TestCheckResourceAttr(resourceName, "value_type", "EXPRESSION"),
					resource.TestCheckResourceAttr(resourceName, "value", "cool"),
					resource.TestCheckResourceAttr(resourceName, "claim_type", "RESOURCE"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "name", "test_updated"),
					resource.TestCheckResourceAttr(resourceName, "value_type", "EXPRESSION"),
					resource.TestCheckResourceAttr(resourceName, "value", "cool_updated"),
					resource.TestCheckResourceAttr(resourceName, "claim_type", "RESOURCE"),
				),
			},
		},
	})
}

func TestAccOktaAuthServerClaim_groupType(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := fmt.Sprintf("%s.test", authServerClaim)
	swResourceName := fmt.Sprintf("%s.test_sw", authServerClaim)
	mgr := newFixtureManager(authServerClaim)
	config := mgr.GetFixtures("basic_group.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(authServer, authServerExists),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "name", "test"),
					resource.TestCheckResourceAttr(resourceName, "group_filter_type", "EQUALS"),
					resource.TestCheckResourceAttr(resourceName, "value_type", "GROUPS"),
					resource.TestCheckResourceAttr(resourceName, "value", "Everyone"),
					resource.TestCheckResourceAttr(resourceName, "claim_type", "RESOURCE"),

					resource.TestCheckResourceAttr(swResourceName, "status", "ACTIVE"),
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

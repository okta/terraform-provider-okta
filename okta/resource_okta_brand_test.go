package okta

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceOktaBrand_import_update(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(brand)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updatedConfig := mgr.GetFixtures("updated.tf", ri, t)
	importConfig := mgr.GetFixtures("import.tf", ri, t)

	// okta_brand is read and update only, so set up the test by importing the brand first
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy: func(s *terraform.State) error {
			// brand api doens't have real delete for a brand
			return nil
		},
		Steps: []resource.TestStep{
			{
				// this is set up only for import state test, ignore check as import.tf is for testing
				ExpectNonEmptyPlan: true,
				Config:             importConfig,
			},
			{
				ResourceName: "okta_brand.example",
				ImportState:  true,
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					// import should only net one brand
					if len(s) != 1 {
						return errors.New("failed to import into resource into state")
					}
					// simple check
					if len(s[0].Attributes["links"]) <= 2 {
						return fmt.Errorf("there should more than two attributes set on the instance %+v", s[0].Attributes)

					}
					return nil
				},
			},
			{
				Config:  config,
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("okta_brand.example", "id"),
					resource.TestCheckResourceAttr("okta_brand.example", "custom_privacy_policy_url", "https://example.com/privacy-policy"),
					resource.TestCheckResourceAttrSet("okta_brand.example", "links"),
					resource.TestCheckResourceAttr("okta_brand.example", "remove_powered_by_okta", "false"),
				),
			},
			{
				Config:  updatedConfig,
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("okta_brand.example", "id"),
					resource.TestCheckResourceAttr("okta_brand.example", "custom_privacy_policy_url", "https://example.com/privacy-policy-updated"),
					resource.TestCheckResourceAttrSet("okta_brand.example", "links"),
					resource.TestCheckResourceAttr("okta_brand.example", "remove_powered_by_okta", "true"),
				),
			},
		},
	})
}

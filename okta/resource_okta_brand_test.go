package okta

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceOktaBrand_import_update(t *testing.T) {
	mgr := newFixtureManager(brand, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("updated.tf", t)
	importConfig := mgr.GetFixtures("import.tf", t)

	// okta_brand is read and update only, so set up the test by importing the brand first

	// NOTE this test will only pass on an org that hasn't had any of its brand
	// values changed in the Admin UI. Need to look into making this more
	// robust.
	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: importConfig,
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
					resource.TestCheckResourceAttrSet("okta_brand.example", "name"),
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
					resource.TestCheckResourceAttrSet("okta_brand.example", "name"),
					resource.TestCheckResourceAttrSet("okta_brand.example", "links"),
					resource.TestCheckResourceAttr("okta_brand.example", "remove_powered_by_okta", "true"),
				),
			},
		},
	})
}

func TestAccResourceOktaBrand_default_brand(t *testing.T) {
	config := `
resource "okta_brand" "example" {
  brand_id = "default"
  lifecycle {
    ignore_changes = [
	  agree_to_custom_privacy_policy,
	  custom_privacy_policy_url,
	  remove_powered_by_okta
	]
  }
}
	`
	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config:             config,
				Destroy:            false,
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("okta_brand.example", "id"),
					resource.TestCheckResourceAttr("okta_brand.example", "brand_id", "default"),
				),
			},
		},
	})
}

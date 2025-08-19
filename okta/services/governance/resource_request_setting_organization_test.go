package governance_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"testing"
)

//func TestAccRequestSettingOrganizationResource_basic(t *testing.T) {
//	resourceName := "okta_request_setting_organization.rq_setting"
//
//	resource.Test(t, resource.TestCase{
//		PreCheck:                 acctest.AccPreCheck(t),
//		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
//		Steps: []resource.TestStep{
//			{
//				// Initial import step
//				Config: testAccRequestSettingOrganizationConfig(true),
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr(resourceName, "subprocessors_acknowledged", "true"),
//				),
//			},
//			{
//				// Update step: flip the attribute
//				Config: testAccRequestSettingOrganizationConfig(false),
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr(resourceName, "subprocessors_acknowledged", "false"),
//				),
//			},
//		},
//	})
//}
//
//func testAccRequestSettingOrganizationConfig(ack bool) string {
//	return fmt.Sprintf(`
//resource "okta_request_setting_organization" "rq_setting" {
//  id                          = "default"
//  subprocessors_acknowledged  = %t
//}
//`, ack)
//}

func TestAccRequestSettingOrganizationResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				// This step "reads" a pre-existing resource.
				// It does not create it.
				// The `Config` block would need to reference the existing resource ID.
				Config: testAccCheckMyResourceReadConfig(),
				Check: resource.ComposeTestCheckFunc(
					// Checks that the read operation was successful
					resource.TestCheckResourceAttr("my_resource.my_resource_id", "name", "original-name"),
				),
			},
			{
				// This step tests the update functionality.
				// It modifies the attribute of the pre-existing resource.
				Config: testAccCheckMyResourceUpdateConfig(),
				Check: resource.ComposeTestCheckFunc(
					// Verifies the update was successful
					resource.TestCheckResourceAttr("my_resource.my_resource_id", "name", "updated-name"),
				),
			},
		},
	})
}

// Assumes these functions return valid Terraform configurations that reference the pre-existing resource.
func testAccCheckMyResourceReadConfig() string {
	return `
resource "okta_request_setting_organization" "test" {
  // This would somehow point to the pre-existing resource, e.g., using an ID.
  id = "default" 
}
`
}

func testAccCheckMyResourceUpdateConfig() string {
	return `
resource "okta_request_setting_organization" "test" {
  id = "default" 
  subprocessors_acknowledged  = false
}
`
}

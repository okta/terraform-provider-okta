package okta

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	firstNameTestProp = "firstName"
	loginTestProp     = "login"
)

func TestAccResourceOktaUserBaseSchema_crud(t *testing.T) {
	mgr := newFixtureManager("resources", userBaseSchemaProperty, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updated := mgr.GetFixtures("updated.tf", t)
	nonDefault := mgr.GetFixtures("non_default_user_type.tf", t)
	resourceName := fmt.Sprintf("%s.%s", userBaseSchemaProperty, firstNameTestProp)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      nil, // can't delete base properties
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testOktaUserBaseSchemasExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "index", firstNameTestProp),
					resource.TestCheckResourceAttr(resourceName, "title", "First name"),
					resource.TestCheckResourceAttr(resourceName, "type", "string"),
					resource.TestCheckResourceAttr(resourceName, "permissions", "READ_ONLY"),
				),
			},
			{
				Config: updated,
				Check: resource.ComposeTestCheckFunc(
					testOktaUserBaseSchemasExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "index", firstNameTestProp),
					resource.TestCheckResourceAttr(resourceName, "title", "First name"),
					resource.TestCheckResourceAttr(resourceName, "type", "string"),
					resource.TestCheckResourceAttr(resourceName, "required", "true"),
					resource.TestCheckResourceAttr(resourceName, "permissions", "READ_WRITE"),
				),
			},
			{
				Config: nonDefault,
				Check: resource.ComposeTestCheckFunc(
					testOktaUserBaseSchemasExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "index", firstNameTestProp),
					resource.TestCheckResourceAttr(resourceName, "title", "First name"),
					resource.TestCheckResourceAttr(resourceName, "type", "string"),
					resource.TestCheckResourceAttr(resourceName, "permissions", "READ_ONLY"),
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

func TestAccResourceOktaUserBaseSchema_login_crud(t *testing.T) {
	mgr := newFixtureManager("resources", userBaseSchemaProperty, t.Name())
	config := mgr.GetFixtures("basic_login.tf", t)
	updated := mgr.GetFixtures("login_updated.tf", t)
	resourceName := fmt.Sprintf("%s.%s", userBaseSchemaProperty, loginTestProp)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      nil, // can't delete base properties
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testOktaUserBaseSchemasExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "index", loginTestProp),
					resource.TestCheckResourceAttr(resourceName, "title", "Username"),
					resource.TestCheckResourceAttr(resourceName, "type", "string"),
					resource.TestCheckResourceAttr(resourceName, "required", "true"),
					resource.TestCheckResourceAttr(resourceName, "permissions", "READ_ONLY"),
					resource.TestCheckResourceAttr(resourceName, "pattern", "[a-z]+"),
				),
			},
			{
				Config: updated,
				Check: resource.ComposeTestCheckFunc(
					testOktaUserBaseSchemasExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "index", loginTestProp),
					resource.TestCheckResourceAttr(resourceName, "title", "Username"),
					resource.TestCheckResourceAttr(resourceName, "type", "string"),
					resource.TestCheckResourceAttr(resourceName, "required", "true"),
					resource.TestCheckResourceAttr(resourceName, "pattern", ""),
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

func testOktaUserBaseSchemasExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		schemaUserType := "default"
		if rs.Primary.Attributes["user_type"] != "" {
			schemaUserType = rs.Primary.Attributes["user_type"]
		}
		exists, err := testUserSchemaPropertyExists(schemaUserType, rs.Primary.ID, baseSchema)
		if err != nil {
			return fmt.Errorf("failed to find: %v", err)
		}
		if !exists {
			return fmt.Errorf("base property %s does not exist in a profile subschema", rs.Primary.ID)
		}
		return nil
	}
}

// TestAccResourceOktaUserBaseSchema_login_multiple_properties_crud Test for issue 1217 fix.
// https://github.com/okta/terraform-provider-okta/issues/1217 This test would
// fail before the fix was implemented. The fix is to put a calling mutex on
// create and update for the `okta_user_base_schema_property` resource. The Okta
// management API ignores parallel calls to `POST
// /api/v1/meta/schemas/user/{userId}` and our fix is to use a calling mutex in
// the resource to impose the equivelent of `terraform -parallelism=1`
func TestAccResourceOktaUserBaseSchema_login_multiple_properties_crud(t *testing.T) {
	if skipVCRTest(t) {
		return
	}
	config := `
resource "okta_user_base_schema_property" "login" {
	index       = "login"
	title       = "Username"
	type        = "string"
	required    = true
	permissions = "%s"
}
resource "okta_user_base_schema_property" "firstname" {
	index       = "firstName"
	title       = "First name"
	type        = "string"
	permissions = "%s"
}
resource "okta_user_base_schema_property" "lastname" {
	index       = "lastName"
	title       = "Last name"
	type        = "string"
	permissions = "%s"
}
resource "okta_user_base_schema_property" "email" {
	index       = "email"
	title       = "Primary email"
	type        = "string"
	required    = true
	permissions = "%s"
}
resource "okta_user_base_schema_property" "mobilephone" {
	index       = "mobilePhone"
	title       = "Mobile phone"
	type        = "string"
	permissions = "%s"
}`
	ro := make([]interface{}, 5)
	for i := 0; i < 5; i++ {
		ro[i] = "READ_ONLY"
	}
	rw := make([]interface{}, 5)
	for i := 0; i < 5; i++ {
		rw[i] = "READ_WRITE"
	}
	roConfig := fmt.Sprintf(config, ro...)
	rwConfig := fmt.Sprintf(config, rw...)
	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		CheckDestroy:      nil,
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: roConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("okta_user_base_schema_property.login", "permissions", "READ_ONLY"),
					resource.TestCheckResourceAttr("okta_user_base_schema_property.firstname", "permissions", "READ_ONLY"),
					resource.TestCheckResourceAttr("okta_user_base_schema_property.lastname", "permissions", "READ_ONLY"),
					resource.TestCheckResourceAttr("okta_user_base_schema_property.email", "permissions", "READ_ONLY"),
					resource.TestCheckResourceAttr("okta_user_base_schema_property.mobilephone", "permissions", "READ_ONLY"),
				),
			},
			{
				Config: rwConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("okta_user_base_schema_property.login", "permissions", "READ_WRITE"),
					resource.TestCheckResourceAttr("okta_user_base_schema_property.firstname", "permissions", "READ_WRITE"),
					resource.TestCheckResourceAttr("okta_user_base_schema_property.lastname", "permissions", "READ_WRITE"),
					resource.TestCheckResourceAttr("okta_user_base_schema_property.email", "permissions", "READ_WRITE"),
					resource.TestCheckResourceAttr("okta_user_base_schema_property.mobilephone", "permissions", "READ_WRITE"),
				),
			},
		},
	})
}

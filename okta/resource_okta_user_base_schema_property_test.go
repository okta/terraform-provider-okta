package okta

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	firstNameTestProp = "firstName"
	loginTestProp     = "login"
)

func TestAccOktaUserBaseSchema_crud(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(userBaseSchemaProperty)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updated := mgr.GetFixtures("updated.tf", ri, t)
	nonDefault := mgr.GetFixtures("non_default_user_type.tf", ri, t)
	resourceName := fmt.Sprintf("%s.%s", userBaseSchemaProperty, firstNameTestProp)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
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

func TestAccOktaUserBaseSchemaLogin_crud(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(userBaseSchemaProperty)
	config := mgr.GetFixtures("basic_login.tf", ri, t)
	updated := mgr.GetFixtures("login_updated.tf", ri, t)
	resourceName := fmt.Sprintf("%s.%s", userBaseSchemaProperty, loginTestProp)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
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

func testOktaUserBaseSchemasExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
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

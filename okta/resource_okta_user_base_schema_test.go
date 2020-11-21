package okta

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const baseTestProp = "firstName"

func sweepUserBaseSchema(client *testClient) error {
	var errorList []error
	schemaUrl := "/api/v1/meta/schemas/user/default"

	schema, _, err := client.apiSupplement.GetUserSchema(schemaUrl)
	if err != nil {
		return err
	}

	for key := range schema.Definitions.Custom.Properties {
		if strings.HasPrefix(key, testResourcePrefix) {
			if _, err := client.apiSupplement.DeleteUserSchemaProperty(schemaUrl, key); err != nil {
				errorList = append(errorList, err)
			}
		}
	}

	return condenseError(errorList)
}

func TestAccOktaUserBaseSchema_crud(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(userBaseSchema)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updated := mgr.GetFixtures("updated.tf", ri, t)
	nonDefault := mgr.GetFixtures("non_default_user_type.tf", ri, t)
	resourceName := fmt.Sprintf("%s.%s", userBaseSchema, baseTestProp)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil, // can't delete base properties
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testOktaUserBaseSchemasExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "index", baseTestProp),
					resource.TestCheckResourceAttr(resourceName, "title", "First name"),
					resource.TestCheckResourceAttr(resourceName, "type", "string"),
					resource.TestCheckResourceAttr(resourceName, "permissions", "READ_ONLY"),
				),
			},
			{
				Config: updated,
				Check: resource.ComposeTestCheckFunc(
					testOktaUserBaseSchemasExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "index", baseTestProp),
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
					resource.TestCheckResourceAttr(resourceName, "index", baseTestProp),
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

func testOktaUserBaseSchemasExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}
		var schemaUserType = "default"
		if rs.Primary.Attributes["user_type"] != "" {
			schemaUserType = rs.Primary.Attributes["user_type"]
		}
		exists, err := testSchemaPropertyExists(schemaUserType, rs.Primary.ID, baseSchema)
		if err != nil {
			return fmt.Errorf("failed to find: %v", err)
		}
		if !exists {
			return fmt.Errorf("base property %s does not exist in a profile subschema", rs.Primary.ID)
		}
		return nil
	}
}

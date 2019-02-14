package okta

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func sweepUserSchema(client *testClient) error {
	schema, _, err := client.artClient.Schemas.GetUserSchema()
	if err != nil {
		return err
	}
	var errorList []error

	for _, schema := range schema.Definitions.Custom.Properties {
		if strings.HasPrefix(schema.Index, testResourcePrefix) {
			if _, _, err := client.artClient.Schemas.DeleteUserCustomSubSchema(schema.Index); err != nil {
				errorList = append(errorList, err)
			}
		}
	}
	return condenseError(errorList)
}

func TestAccOktaUserSchemas(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(userSchema)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updated := mgr.GetFixtures("updated.tf", ri, t)
	resourceName := buildResourceFQN(userSchema, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(userSchema, testUserSchemaExists),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testOktaUserSchemasExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "index", "testAcc_"+strconv.Itoa(ri)),
					resource.TestCheckResourceAttr(resourceName, "title", "terraform acceptance test"),
					resource.TestCheckResourceAttr(resourceName, "type", "string"),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform acceptance test"),
					resource.TestCheckResourceAttr(resourceName, "required", "false"),
					resource.TestCheckResourceAttr(resourceName, "min_length", "1"),
					resource.TestCheckResourceAttr(resourceName, "max_length", "50"),
					resource.TestCheckResourceAttr(resourceName, "permissions", "READ_ONLY"),
					resource.TestCheckResourceAttr(resourceName, "master", "PROFILE_MASTER"),
					resource.TestCheckResourceAttr(resourceName, "enum.0", "S"),
					resource.TestCheckResourceAttr(resourceName, "enum.1", "M"),
					resource.TestCheckResourceAttr(resourceName, "enum.2", "L"),
					resource.TestCheckResourceAttr(resourceName, "enum.3", "XL"),
					resource.TestCheckResourceAttr(resourceName, "one_of.#", "4"),
				),
			},
			{
				Config: updated,
				Check: resource.ComposeTestCheckFunc(
					testOktaUserSchemasExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "index", "testAcc_"+strconv.Itoa(ri)),
					resource.TestCheckResourceAttr(resourceName, "title", "terraform acceptance test updated"),
					resource.TestCheckResourceAttr(resourceName, "type", "string"),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform acceptance test updated"),
					resource.TestCheckResourceAttr(resourceName, "required", "true"),
					resource.TestCheckResourceAttr(resourceName, "min_length", "1"),
					resource.TestCheckResourceAttr(resourceName, "max_length", "70"),
					resource.TestCheckResourceAttr(resourceName, "permissions", "READ_WRITE"),
					resource.TestCheckResourceAttr(resourceName, "master", "OKTA"),
					resource.TestCheckResourceAttr(resourceName, "enum.0", "S"),
					resource.TestCheckResourceAttr(resourceName, "enum.1", "M"),
					resource.TestCheckResourceAttr(resourceName, "enum.2", "L"),
					resource.TestCheckResourceAttr(resourceName, "enum.3", "XXL"),
					resource.TestCheckResourceAttr(resourceName, "one_of.#", "4"),
				),
			},
		},
	})
}

func TestAccOktaUserSchemas_arrayString(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := buildResourceFQN(userSchema, ri)
	config := testOktaUserSchemas_arrayString(ri)
	updatedConfig := testOktaUserSchemas_arrayStringUpdated(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(userSchema, testUserSchemaExists),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testOktaUserSchemasExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "index", "testAcc_"+strconv.Itoa(ri)),
					resource.TestCheckResourceAttr(resourceName, "title", "terraform acceptance test"),
					resource.TestCheckResourceAttr(resourceName, "type", "array"),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform acceptance test"),
					resource.TestCheckResourceAttr(resourceName, "array_type", "string"),
					resource.TestCheckResourceAttr(resourceName, "required", "false"),
					resource.TestCheckResourceAttr(resourceName, "permissions", "READ_ONLY"),
					resource.TestCheckResourceAttr(resourceName, "master", "PROFILE_MASTER"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testOktaUserSchemasExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "index", "testAcc_"+strconv.Itoa(ri)),
					resource.TestCheckResourceAttr(resourceName, "title", "terraform acceptance test updated"),
					resource.TestCheckResourceAttr(resourceName, "type", "array"),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform acceptance test updated"),
					resource.TestCheckResourceAttr(resourceName, "array_type", "string"),
					resource.TestCheckResourceAttr(resourceName, "required", "true"),
					resource.TestCheckResourceAttr(resourceName, "permissions", "READ_WRITE"),
					resource.TestCheckResourceAttr(resourceName, "master", "OKTA"),
				),
			},
		},
	})
}

func TestAccUserSchemaImport(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := buildResourceFQN(userSchema, ri)
	mgr := newFixtureManager(userSchema)
	config := mgr.GetFixtures("basic.tf", ri, t)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					if len(s) != 1 {
						return errors.New("Failed to import schema into state")
					}

					return nil
				},
			},
		},
	})
}

func testOktaUserSchemasExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if exists, _ := testUserSchemaExists(rs.Primary.ID); !exists {
			return fmt.Errorf("Failed to find %s", rs.Primary.ID)
		}
		return nil
	}
}

func testUserSchemaExists(index string) (bool, error) {
	client := getClientFromMetadata(testAccProvider.Meta())
	subschema, _, err := client.Schemas.GetUserSubSchemaIndex(customSchema)
	if err != nil {
		return false, fmt.Errorf("Error Listing User Subschema in Okta: %v", err)
	}
	for _, key := range subschema {
		if key == index {
			return true, nil
		}
	}

	return false, nil
}

func testOktaUserSchemas_arrayString(rInt int) string {
	return fmt.Sprintf(`
resource "okta_user_schema" "testAcc_%[1]d" {
  index     = "testAcc_%[1]d"
  title     = "terraform acceptance test"
  type      = "array"
  description = "terraform acceptance test"
  array_type = "string"
  required    = false
  permissions = "READ_ONLY"
  master      = "PROFILE_MASTER"
}
`, rInt)
}

func testOktaUserSchemas_arrayStringUpdated(rInt int) string {
	return fmt.Sprintf(`
resource "okta_user_schema" "testAcc_%[1]d" {
  index     = "testAcc_%[1]d"
  title     = "terraform acceptance test updated"
  type      = "array"
  description = "terraform acceptance test updated"
  array_type = "string"
  required    = true
  permissions = "READ_WRITE"
  master      = "OKTA"
}
`, rInt)
}

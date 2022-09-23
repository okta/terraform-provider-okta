package okta

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAppUserSchemas_crud(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appUserSchemaProperty)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updated := mgr.GetFixtures("updated.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", appUserSchemaProperty)

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(appUserSchemaProperty, testAppUserSchemaExists),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAppUserSchemasExists(resourceName),
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
					testAppUserSchemasExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "index", "testAcc_"+strconv.Itoa(ri)),
					resource.TestCheckResourceAttr(resourceName, "title", "terraform acceptance test"),
					resource.TestCheckResourceAttr(resourceName, "type", "string"),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform acceptance test updated"),
					resource.TestCheckResourceAttr(resourceName, "required", "true"),
					resource.TestCheckResourceAttr(resourceName, "master", "OKTA"),
					resource.TestCheckResourceAttr(resourceName, "scope", "SELF"),
				),
			},
		},
	})
}

func TestAccAppUserSchemas_array_enum_number(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appUserSchemaProperty)
	resourceName := fmt.Sprintf("%s.test", appUserSchemaProperty)
	config := `
resource "okta_app_oauth" "test" {
	label          = "testAcc_replace_with_uuid"
	type           = "native"
	grant_types    = ["authorization_code"]
	redirect_uris  = ["http://d.com/"]
	response_types = ["code"]
	}
	
resource "okta_app_user_schema_property" "test" {
	app_id      = okta_app_oauth.test.id
	index       = "testAcc_replace_with_uuid"
	title       = "terraform acceptance test"
	type        = "array"
	description = "testing"
	required    = false
	permissions = "READ_ONLY"
	master      = "PROFILE_MASTER"
	array_type  = "number"
	array_enum  = ["0.011", "0.022", "0.033"]
	array_one_of {
	  title = "number point oh one one"
	  const = "0.011"
	}
	array_one_of {
	  title = "number point oh two two"
	  const = "0.022"
	}
	array_one_of {
	  title = "number point oh three three"
	  const = "0.033"
	}
}
`
	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(appUserSchemaProperty, testAppUserSchemaExists),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config, ri),
				Check: resource.ComposeTestCheckFunc(
					testAppUserSchemasExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", "array"),
					resource.TestCheckResourceAttr(resourceName, "array_type", "number"),
					resource.TestCheckResourceAttr(resourceName, "array_enum.0", "0.011"),
					resource.TestCheckResourceAttr(resourceName, "array_enum.1", "0.022"),
					resource.TestCheckResourceAttr(resourceName, "array_enum.2", "0.033"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.0.title", "number point oh one one"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.0.const", "0.011"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.1.title", "number point oh two two"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.1.const", "0.022"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.2.title", "number point oh three three"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.2.const", "0.033"),
				),
			},
		},
	})
}

func TestAccAppUserSchemas_enum_number(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appUserSchemaProperty)
	resourceName := fmt.Sprintf("%s.test", appUserSchemaProperty)
	config := `
resource "okta_app_oauth" "test" {
	label          = "testAcc_replace_with_uuid"
	type           = "native"
	grant_types    = ["authorization_code"]
	redirect_uris  = ["http://d.com/"]
	response_types = ["code"]
	}
	
resource "okta_app_user_schema_property" "test" {
	app_id      = okta_app_oauth.test.id
	index       = "testAcc_replace_with_uuid"
	title       = "terraform acceptance test"
	type        = "number"
	description = "testing"
	required    = false
	permissions = "READ_ONLY"
	master      = "PROFILE_MASTER"
	enum  		= ["0.011", "0.022", "0.033"]
	one_of {
	  title = "number point oh one one"
	  const = "0.011"
	}
	one_of {
	  title = "number point oh two two"
	  const = "0.022"
	}
	one_of {
	  title = "number point oh three three"
	  const = "0.033"
	}
}
`
	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(appUserSchemaProperty, testAppUserSchemaExists),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config, ri),
				Check: resource.ComposeTestCheckFunc(
					testAppUserSchemasExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", "number"),
					resource.TestCheckResourceAttr(resourceName, "enum.0", "0.011"),
					resource.TestCheckResourceAttr(resourceName, "enum.1", "0.022"),
					resource.TestCheckResourceAttr(resourceName, "enum.2", "0.033"),
					resource.TestCheckResourceAttr(resourceName, "one_of.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "one_of.0.title", "number point oh one one"),
					resource.TestCheckResourceAttr(resourceName, "one_of.0.const", "0.011"),
					resource.TestCheckResourceAttr(resourceName, "one_of.1.title", "number point oh two two"),
					resource.TestCheckResourceAttr(resourceName, "one_of.1.const", "0.022"),
					resource.TestCheckResourceAttr(resourceName, "one_of.2.title", "number point oh three three"),
					resource.TestCheckResourceAttr(resourceName, "one_of.2.const", "0.033"),
				),
			},
		},
	})
}

func TestAccAppUserSchemas_array_enum_integer(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appUserSchemaProperty)
	resourceName := fmt.Sprintf("%s.test", appUserSchemaProperty)
	config := `
resource "okta_app_oauth" "test" {
	label          = "testAcc_replace_with_uuid"
	type           = "native"
	grant_types    = ["authorization_code"]
	redirect_uris  = ["http://d.com/"]
	response_types = ["code"]
	}
	
resource "okta_app_user_schema_property" "test" {
	app_id      = okta_app_oauth.test.id
	index       = "testAcc_replace_with_uuid"
	title       = "terraform acceptance test"
	type        = "array"
	description = "testing"
	required    = false
	permissions = "READ_ONLY"
	master      = "PROFILE_MASTER"
	array_type  = "integer"
	array_enum  = [4, 5, 6]
	array_one_of {
		const = "4"
		title = "integer four"
	}
	array_one_of {
		const = "5"
		title = "integer five"
	}
	array_one_of {
		const = "6"
		title = "integer six"
	}
}
`
	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(appUserSchemaProperty, testAppUserSchemaExists),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config, ri),
				Check: resource.ComposeTestCheckFunc(
					testAppUserSchemasExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", "array"),
					resource.TestCheckResourceAttr(resourceName, "array_type", "integer"),
					resource.TestCheckResourceAttr(resourceName, "array_enum.0", "4"),
					resource.TestCheckResourceAttr(resourceName, "array_enum.1", "5"),
					resource.TestCheckResourceAttr(resourceName, "array_enum.2", "6"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.0.title", "integer four"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.0.const", "4"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.1.title", "integer five"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.1.const", "5"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.2.title", "integer six"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.2.const", "6"),
				),
			},
		},
	})
}

func TestAccAppUserSchemas_enum_integer(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appUserSchemaProperty)
	resourceName := fmt.Sprintf("%s.test", appUserSchemaProperty)
	config := `
resource "okta_app_oauth" "test" {
	label          = "testAcc_replace_with_uuid"
	type           = "native"
	grant_types    = ["authorization_code"]
	redirect_uris  = ["http://d.com/"]
	response_types = ["code"]
	}
	
resource "okta_app_user_schema_property" "test" {
	app_id      = okta_app_oauth.test.id
	index       = "testAcc_replace_with_uuid"
	title       = "terraform acceptance test"
	type        = "integer"
	description = "testing"
	required    = false
	permissions = "READ_ONLY"
	master      = "PROFILE_MASTER"
	enum        = [4, 5, 6]
	one_of {
		const = "4"
		title = "integer four"
	}
	one_of {
		const = "5"
		title = "integer five"
	}
	one_of {
		const = "6"
		title = "integer six"
	}
}
`
	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(appUserSchemaProperty, testAppUserSchemaExists),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config, ri),
				Check: resource.ComposeTestCheckFunc(
					testAppUserSchemasExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", "integer"),
					resource.TestCheckResourceAttr(resourceName, "enum.0", "4"),
					resource.TestCheckResourceAttr(resourceName, "enum.1", "5"),
					resource.TestCheckResourceAttr(resourceName, "enum.2", "6"),
					resource.TestCheckResourceAttr(resourceName, "one_of.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "one_of.0.title", "integer four"),
					resource.TestCheckResourceAttr(resourceName, "one_of.0.const", "4"),
					resource.TestCheckResourceAttr(resourceName, "one_of.1.title", "integer five"),
					resource.TestCheckResourceAttr(resourceName, "one_of.1.const", "5"),
					resource.TestCheckResourceAttr(resourceName, "one_of.2.title", "integer six"),
					resource.TestCheckResourceAttr(resourceName, "one_of.2.const", "6"),
				),
			},
		},
	})
}

func TestAccAppUserSchemas_array_enum_boolean(t *testing.T) {
	t.Skip("The test is failing due to core issue. Similar test TestAccResourceOktaGroupSchema_array_enum_boolean has passed in the past")
	ri := acctest.RandInt()
	mgr := newFixtureManager(appUserSchemaProperty)
	resourceName := fmt.Sprintf("%s.test", appUserSchemaProperty)
	config := `
resource "okta_app_oauth" "test" {
	label          = "testAcc_replace_with_uuid"
	type           = "native"
	grant_types    = ["authorization_code"]
	redirect_uris  = ["http://d.com/"]
	response_types = ["code"]
	}
	
resource "okta_app_user_schema_property" "test" {
	app_id      = okta_app_oauth.test.id
	index       = "testAcc_replace_with_uuid"
	title       = "terraform acceptance test"
	type        = "array"
	description = "testing"
	required    = false
	permissions = "READ_ONLY"
	master      = "PROFILE_MASTER"
	array_type  = "boolean"
	array_enum  = ["true", "false"]
	array_one_of {
	  title = "boolean True"
	  const = "true"
	}
	array_one_of {
	  title = "boolean False"
	  const = "false"
	}
}
`
	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(appUserSchemaProperty, testAppUserSchemaExists),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config, ri),
				Check: resource.ComposeTestCheckFunc(
					testAppUserSchemasExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", "array"),
					resource.TestCheckResourceAttr(resourceName, "array_type", "boolean"),
					resource.TestCheckResourceAttr(resourceName, "array_enum.0", "true"),
					resource.TestCheckResourceAttr(resourceName, "array_enum.1", "false"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.0.title", "boolean True"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.0.const", "true"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.1.title", "boolean False"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.1.const", "false"),
				),
			},
		},
	})
}

func TestAccAppUserSchemas_enum_boolean(t *testing.T) {
	t.Skip("The test is failing due to core issue. Similar test TestAccResourceOktaGroupSchema_enum_boolean has passed in the past")
	ri := acctest.RandInt()
	mgr := newFixtureManager(appUserSchemaProperty)
	resourceName := fmt.Sprintf("%s.test", appUserSchemaProperty)
	config := `
resource "okta_app_oauth" "test" {
	label          = "testAcc_replace_with_uuid"
	type           = "native"
	grant_types    = ["authorization_code"]
	redirect_uris  = ["http://d.com/"]
	response_types = ["code"]
	}
	
resource "okta_app_user_schema_property" "test" {
	app_id      = okta_app_oauth.test.id
	index       = "testAcc_replace_with_uuid"
	title       = "terraform acceptance test"
	type        = "boolean"
	description = "testing"
	required    = false
	permissions = "READ_ONLY"
	master      = "PROFILE_MASTER"
	enum  		= ["true", "false"]
	one_of {
	  title = "boolean True"
	  const = "true"
	}
	one_of {
	  title = "boolean False"
	  const = "false"
	}
}
`
	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(appUserSchemaProperty, testAppUserSchemaExists),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config, ri),
				Check: resource.ComposeTestCheckFunc(
					testAppUserSchemasExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", "boolean"),
					resource.TestCheckResourceAttr(resourceName, "enum.0", "true"),
					resource.TestCheckResourceAttr(resourceName, "enum.1", "false"),
					resource.TestCheckResourceAttr(resourceName, "one_of.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "one_of.0.title", "boolean True"),
					resource.TestCheckResourceAttr(resourceName, "one_of.0.const", "true"),
					resource.TestCheckResourceAttr(resourceName, "one_of.1.title", "boolean False"),
					resource.TestCheckResourceAttr(resourceName, "one_of.1.const", "false"),
				),
			},
		},
	})
}

func TestAccAppUserSchemas_array_enum_string(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appUserSchemaProperty)
	resourceName := fmt.Sprintf("%s.test", appUserSchemaProperty)
	config := `
resource "okta_app_oauth" "test" {
	label          = "testAcc_replace_with_uuid"
	type           = "native"
	grant_types    = ["authorization_code"]
	redirect_uris  = ["http://d.com/"]
	response_types = ["code"]
	}
	
resource "okta_app_user_schema_property" "test" {
	app_id      = okta_app_oauth.test.id
	index       = "testAcc_replace_with_uuid"
	title       = "terraform acceptance test"
	type        = "array"
	description = "testing"
	required    = false
	permissions = "READ_ONLY"
	master      = "PROFILE_MASTER"
	array_type  = "string"
	array_enum  = ["one", "two", "three"]
	array_one_of {
	  title = "string One"
	  const = "one"
	}
	array_one_of {
	  title = "string Two"
	  const = "two"
	}
	array_one_of {
	  title = "string Three"
	  const = "three"
	}
}
`
	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(appUserSchemaProperty, testAppUserSchemaExists),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config, ri),
				Check: resource.ComposeTestCheckFunc(
					testAppUserSchemasExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", "array"),
					resource.TestCheckResourceAttr(resourceName, "array_type", "string"),
					resource.TestCheckResourceAttr(resourceName, "array_enum.0", "one"),
					resource.TestCheckResourceAttr(resourceName, "array_enum.1", "two"),
					resource.TestCheckResourceAttr(resourceName, "array_enum.2", "three"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.0.title", "string One"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.0.const", "one"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.1.title", "string Two"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.1.const", "two"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.2.title", "string Three"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.2.const", "three"),
				),
			},
		},
	})
}

func TestAccAppUserSchemas_enum_string(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appUserSchemaProperty)
	resourceName := fmt.Sprintf("%s.test", appUserSchemaProperty)
	config := `
resource "okta_app_oauth" "test" {
	label          = "testAcc_replace_with_uuid"
	type           = "native"
	grant_types    = ["authorization_code"]
	redirect_uris  = ["http://d.com/"]
	response_types = ["code"]
	}
	
resource "okta_app_user_schema_property" "test" {
	app_id      = okta_app_oauth.test.id
	index       = "testAcc_replace_with_uuid"
	title       = "terraform acceptance test"
	description = "testing"
	required    = false
	permissions = "READ_ONLY"
	master      = "PROFILE_MASTER"
	type  = "string"
	enum  = ["one", "two", "three"]
	one_of {
	  title = "string One"
	  const = "one"
	}
	one_of {
	  title = "string Two"
	  const = "two"
	}
	one_of {
	  title = "string Three"
	  const = "three"
	}
}
`
	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(appUserSchemaProperty, testAppUserSchemaExists),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config, ri),
				Check: resource.ComposeTestCheckFunc(
					testAppUserSchemasExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", "string"),
					resource.TestCheckResourceAttr(resourceName, "enum.0", "one"),
					resource.TestCheckResourceAttr(resourceName, "enum.1", "two"),
					resource.TestCheckResourceAttr(resourceName, "enum.2", "three"),
					resource.TestCheckResourceAttr(resourceName, "one_of.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "one_of.0.title", "string One"),
					resource.TestCheckResourceAttr(resourceName, "one_of.0.const", "one"),
					resource.TestCheckResourceAttr(resourceName, "one_of.1.title", "string Two"),
					resource.TestCheckResourceAttr(resourceName, "one_of.1.const", "two"),
					resource.TestCheckResourceAttr(resourceName, "one_of.2.title", "string Three"),
					resource.TestCheckResourceAttr(resourceName, "one_of.2.const", "three"),
				),
			},
		},
	})
}

func testAppUserSchemasExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}

		if exists, _ := testAppUserSchemaExists(rs.Primary.ID); !exists {
			return fmt.Errorf("failed to find %s", rs.Primary.ID)
		}
		return nil
	}
}

func testAppUserSchemaExists(index string) (bool, error) {
	ids := strings.Split(index, "/")
	schema, resp, err := getOktaClientFromMetadata(testAccProvider.Meta()).UserSchema.
		GetApplicationUserSchema(context.Background(), ids[0])
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("failed to get application user schema: %v", err)
	}
	cu := userSchemaCustomAttribute(schema, ids[1])
	if cu != nil {
		return true, nil
	}
	return false, nil
}

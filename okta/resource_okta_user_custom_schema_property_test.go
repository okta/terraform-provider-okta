package okta

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceOktaUserSchema_crud(t *testing.T) {
	mgr := newFixtureManager(userSchemaProperty, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updated := mgr.GetFixtures("updated.tf", t)
	unique := mgr.GetFixtures("unique.tf", t)
	nonDefault := mgr.GetFixtures("non_default_user_type.tf", t)
	resourceName := buildResourceFQN(userSchemaProperty, mgr.Seed)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkOktaUserSchemasDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testOktaUserSchemasExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "index", "testAcc_"+strconv.Itoa(mgr.Seed)),
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
					resource.TestCheckResourceAttr(resourceName, "scope", "SELF"),
				),
			},
			{
				Config: updated,
				Check: resource.ComposeTestCheckFunc(
					testOktaUserSchemasExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "index", "testAcc_"+strconv.Itoa(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "title", "terraform acceptance test updated 004"),
					resource.TestCheckResourceAttr(resourceName, "type", "string"),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform acceptance test updated 004"),
					// FIXME when the schema is updated and set to true why does this cause the Admin login flow to be unusable?
					resource.TestCheckResourceAttr(resourceName, "required", "false"),
					resource.TestCheckResourceAttr(resourceName, "min_length", "1"),
					resource.TestCheckResourceAttr(resourceName, "max_length", "70"),
					resource.TestCheckResourceAttr(resourceName, "permissions", "READ_WRITE"),
					resource.TestCheckResourceAttr(resourceName, "master", "OKTA"),
					resource.TestCheckResourceAttr(resourceName, "enum.0", "S"),
					resource.TestCheckResourceAttr(resourceName, "enum.1", "M"),
					resource.TestCheckResourceAttr(resourceName, "enum.2", "L"),
					resource.TestCheckResourceAttr(resourceName, "enum.3", "XXL"),
					resource.TestCheckResourceAttr(resourceName, "one_of.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "pattern", ".+"),
					resource.TestCheckResourceAttr(resourceName, "scope", "NONE"),
				),
			},
			{
				Config: unique,
				Check: resource.ComposeTestCheckFunc(
					testOktaUserSchemasExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "index", "testAcc_"+strconv.Itoa(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "title", "terraform acceptance test setting unique attribute to UNIQUE_VALIDATED 007"),
					resource.TestCheckResourceAttr(resourceName, "type", "string"),
					// FIXME when the schema is updated and set to true why does this cause the a prompt:
					// "terraform acceptance test setting unique attribute to UNIQUE_VALIDATED 007" to be added to the flow. It doesn't block login
					// as there is a form and you can enter a number into it
					resource.TestCheckResourceAttr(resourceName, "description", "terraform acceptance test setting unique attribute to UNIQUE_VALIDATED 007"),
					resource.TestCheckResourceAttr(resourceName, "required", "false"),
					resource.TestCheckResourceAttr(resourceName, "min_length", "1"),
					resource.TestCheckResourceAttr(resourceName, "max_length", "70"),
					resource.TestCheckResourceAttr(resourceName, "permissions", "READ_WRITE"),
					resource.TestCheckResourceAttr(resourceName, "master", "OKTA"),
					resource.TestCheckResourceAttr(resourceName, "unique", "UNIQUE_VALIDATED"),
				),
			},
			{
				Config: nonDefault,
				Check: resource.ComposeTestCheckFunc(
					testOktaUserSchemasExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "index", "testAcc_"+strconv.Itoa(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "title", "terraform acceptance test"),
					resource.TestCheckResourceAttr(resourceName, "type", "string"),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform acceptance test"),
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
		},
	})
}

func TestAccResourceOktaUserSchema_array_enum(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", userSchemaProperty)
	mgr := newFixtureManager(userSchemaProperty, t.Name())
	config := mgr.GetFixtures("array_string.tf", t)
	updatedConfig := mgr.GetFixtures("array_string_updated.tf", t)
	arrayEnum := mgr.GetFixtures("array_enum.tf", t)
	arrayNumber := mgr.GetFixtures("array_number.tf", t)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkOktaUserSchemasDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testOktaUserSchemasExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "index", "testAcc_"+strconv.Itoa(mgr.Seed)),
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
					resource.TestCheckResourceAttr(resourceName, "index", "testAcc_"+strconv.Itoa(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "title", "terraform acceptance test updated 005"),
					resource.TestCheckResourceAttr(resourceName, "type", "array"),
					// FIXME when the schema is updated and set to true why does this cause the a prompt:
					// "terraform acceptance test updated 005" to be added to the flow. It doesn't block login
					// as there is a form and you can enter a value into it
					resource.TestCheckResourceAttr(resourceName, "description", "terraform acceptance test updated 005"),
					resource.TestCheckResourceAttr(resourceName, "array_type", "string"),
					resource.TestCheckResourceAttr(resourceName, "required", "false"),
					resource.TestCheckResourceAttr(resourceName, "permissions", "READ_WRITE"),
					resource.TestCheckResourceAttr(resourceName, "master", "OKTA"),
				),
			},
			{
				Config: arrayEnum,
				Check: resource.ComposeTestCheckFunc(
					testOktaUserSchemasExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "index", "testAcc_"+strconv.Itoa(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "title", "terraform acceptance test"),
					resource.TestCheckResourceAttr(resourceName, "type", "array"),
					resource.TestCheckResourceAttr(resourceName, "array_type", "string"),
					resource.TestCheckResourceAttr(resourceName, "description", "testing"),
					resource.TestCheckResourceAttr(resourceName, "required", "false"),
					resource.TestCheckResourceAttr(resourceName, "master", "OKTA"),
					resource.TestCheckResourceAttr(resourceName, "scope", "NONE"),
					resource.TestCheckResourceAttr(resourceName, "array_enum.0", "test"),
					resource.TestCheckResourceAttr(resourceName, "array_enum.1", "1"),
					resource.TestCheckResourceAttr(resourceName, "array_enum.2", "2"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.#", "3"),
				),
			},
			{
				Config: arrayNumber,
				Check: resource.ComposeTestCheckFunc(
					testOktaUserSchemasExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "index", "testAcc_"+strconv.Itoa(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "title", "terraform acceptance test"),
					resource.TestCheckResourceAttr(resourceName, "type", "array"),
					resource.TestCheckResourceAttr(resourceName, "array_type", "number"),
					resource.TestCheckResourceAttr(resourceName, "description", "testing"),
					resource.TestCheckResourceAttr(resourceName, "required", "false"),
					resource.TestCheckResourceAttr(resourceName, "master", "OKTA"),
					resource.TestCheckResourceAttr(resourceName, "scope", "SELF"),
					resource.TestCheckResourceAttr(resourceName, "array_enum.0", "0.01"),
					resource.TestCheckResourceAttr(resourceName, "array_enum.1", "0.02"),
					resource.TestCheckResourceAttr(resourceName, "array_enum.2", "0.03"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.#", "3"),
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

func TestAccResourceOktaUserSchema_array_enum_number(t *testing.T) {
	mgr := newFixtureManager(userSchemaProperty, t.Name())
	resourceName := fmt.Sprintf("%s.test", userSchemaProperty)
	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkOktaUserSchemasDestroy,
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(`
			resource "okta_user_schema_property" "test" {
			  index       = "testAcc_replace_with_uuid"
			  title       = "terraform acceptance test"
			  type        = "array"
			  description = "testing"
			  master      = "OKTA"
			  scope       = "SELF"
			  array_type  = "number"
			  array_enum  = ["0.01", "0.02", "0.03"]
			  array_one_of {
			    title = "number point oh one"
			    const = "0.01"
			  }
			  array_one_of {
			    title = "number point oh two"
			    const = "0.02"
			  }
			  array_one_of {
			    title = "number point oh three"
			    const = "0.03"
			  }
			}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "array"),
					resource.TestCheckResourceAttr(resourceName, "array_type", "number"),
					resource.TestCheckResourceAttr(resourceName, "array_enum.0", "0.01"),
					resource.TestCheckResourceAttr(resourceName, "array_enum.1", "0.02"),
					resource.TestCheckResourceAttr(resourceName, "array_enum.2", "0.03"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.0.title", "number point oh one"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.0.const", "0.01"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.1.title", "number point oh two"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.1.const", "0.02"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.2.title", "number point oh three"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.2.const", "0.03"),
				),
			},
			{
				Config: mgr.ConfigReplace(`
			resource "okta_user_schema_property" "test" {
			  index       = "testAcc_replace_with_uuid"
			  title       = "terraform acceptance test"
			  type        = "array"
			  description = "testing"
			  master      = "OKTA"
			  scope       = "SELF"
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
			}`),
				Check: resource.ComposeTestCheckFunc(
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

func TestAccResourceOktaUserSchema_enum_number(t *testing.T) {
	mgr := newFixtureManager(userSchemaProperty, t.Name())
	resourceName := fmt.Sprintf("%s.test", userSchemaProperty)
	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkOktaUserSchemasDestroy,
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(`
			resource "okta_user_schema_property" "test" {
			  index       = "testAcc_replace_with_uuid"
			  title       = "terraform acceptance test"
			  type        = "number"
			  description = "testing"
			  master      = "OKTA"
			  scope       = "SELF"
			  enum  = ["0.01", "0.02", "0.03"]
			  one_of {
			    title = "number point oh one"
			    const = "0.01"
			  }
			  one_of {
			    title = "number point oh two"
			    const = "0.02"
			  }
			  one_of {
			    title = "number point oh three"
			    const = "0.03"
			  }
			}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "number"),
					resource.TestCheckResourceAttr(resourceName, "enum.0", "0.01"),
					resource.TestCheckResourceAttr(resourceName, "enum.1", "0.02"),
					resource.TestCheckResourceAttr(resourceName, "enum.2", "0.03"),
					resource.TestCheckResourceAttr(resourceName, "one_of.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "one_of.0.title", "number point oh one"),
					resource.TestCheckResourceAttr(resourceName, "one_of.0.const", "0.01"),
					resource.TestCheckResourceAttr(resourceName, "one_of.1.title", "number point oh two"),
					resource.TestCheckResourceAttr(resourceName, "one_of.1.const", "0.02"),
					resource.TestCheckResourceAttr(resourceName, "one_of.2.title", "number point oh three"),
					resource.TestCheckResourceAttr(resourceName, "one_of.2.const", "0.03"),
				),
			},
			{
				Config: mgr.ConfigReplace(`
			resource "okta_user_schema_property" "test" {
			  index       = "testAcc_replace_with_uuid"
			  title       = "terraform acceptance test"
			  type        = "number"
			  description = "testing"
			  master      = "OKTA"
			  scope       = "SELF"
			  enum  = ["0.011", "0.022", "0.033"]
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
			}`),
				Check: resource.ComposeTestCheckFunc(
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

func TestAccResourceOktaUserSchema_array_enum_integer(t *testing.T) {
	mgr := newFixtureManager(userSchemaProperty, t.Name())
	resourceName := fmt.Sprintf("%s.test", userSchemaProperty)
	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkOktaUserSchemasDestroy,
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(`
			resource "okta_user_schema_property" "test" {
			  index       = "testAcc_replace_with_uuid"
			  title       = "terraform acceptance test"
			  type        = "array"
			  description = "testing"
			  master      = "OKTA"
			  scope       = "SELF"
			  array_type  = "integer"
			  array_enum  = ["1", "2", "3"]
			  array_one_of {
			    title = "integer one"
			    const = "1"
			  }
			  array_one_of {
			    title = "integer two"
			    const = "2"
			  }
			  array_one_of {
			    title = "integer three"
			    const = "3"
			  }
			}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "array"),
					resource.TestCheckResourceAttr(resourceName, "array_type", "integer"),
					resource.TestCheckResourceAttr(resourceName, "array_enum.0", "1"),
					resource.TestCheckResourceAttr(resourceName, "array_enum.1", "2"),
					resource.TestCheckResourceAttr(resourceName, "array_enum.2", "3"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.0.title", "integer one"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.0.const", "1"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.1.title", "integer two"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.1.const", "2"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.2.title", "integer three"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.2.const", "3"),
				),
			},
			{
				Config: mgr.ConfigReplace(`
			resource "okta_user_schema_property" "test" {
			  index       = "testAcc_replace_with_uuid"
			  title       = "terraform acceptance test"
			  type        = "array"
			  description = "testing"
			  master      = "OKTA"
			  scope       = "SELF"
			  array_type  = "integer"
			  array_enum  = ["4", "5", "6"]
			  array_one_of {
			    title = "integer four"
			    const = "4"
			  }
			  array_one_of {
			    title = "integer five"
			    const = "5"
			  }
			  array_one_of {
			    title = "integer six"
			    const = "6"
			  }
			}`),
				Check: resource.ComposeTestCheckFunc(
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

func TestAccResourceOktaUserSchema_enum_integer(t *testing.T) {
	mgr := newFixtureManager(userSchemaProperty, t.Name())
	resourceName := fmt.Sprintf("%s.test", userSchemaProperty)
	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkOktaUserSchemasDestroy,
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(`
			resource "okta_user_schema_property" "test" {
			  index       = "testAcc_replace_with_uuid"
			  title       = "terraform acceptance test"
			  type        = "integer"
			  description = "testing"
			  master      = "OKTA"
			  scope       = "SELF"
			  enum  = ["1", "2", "3"]
			  one_of {
			    title = "integer one"
			    const = "1"
			  }
			  one_of {
			    title = "integer two"
			    const = "2"
			  }
			  one_of {
			    title = "integer three"
			    const = "3"
			  }
			}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "integer"),
					resource.TestCheckResourceAttr(resourceName, "enum.0", "1"),
					resource.TestCheckResourceAttr(resourceName, "enum.1", "2"),
					resource.TestCheckResourceAttr(resourceName, "enum.2", "3"),
					resource.TestCheckResourceAttr(resourceName, "one_of.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "one_of.0.title", "integer one"),
					resource.TestCheckResourceAttr(resourceName, "one_of.0.const", "1"),
					resource.TestCheckResourceAttr(resourceName, "one_of.1.title", "integer two"),
					resource.TestCheckResourceAttr(resourceName, "one_of.1.const", "2"),
					resource.TestCheckResourceAttr(resourceName, "one_of.2.title", "integer three"),
					resource.TestCheckResourceAttr(resourceName, "one_of.2.const", "3"),
				),
			},
			{
				Config: mgr.ConfigReplace(`
			resource "okta_user_schema_property" "test" {
			  index       = "testAcc_replace_with_uuid"
			  title       = "terraform acceptance test"
			  type        = "integer"
			  description = "testing"
			  master      = "OKTA"
			  scope       = "SELF"
			  enum  = ["4", "5", "6"]
			  one_of {
			    title = "integer four"
			    const = "4"
			  }
			  one_of {
			    title = "integer five"
			    const = "5"
			  }
			  one_of {
			    title = "integer six"
			    const = "6"
			  }
			}`),
				Check: resource.ComposeTestCheckFunc(
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

func TestAccResourceOktaUserSchema_array_enum_boolean(t *testing.T) {
	t.Skip("TODO deal with apparent monolith bug")
	// TODO deal with apparent monolith bug:
	// "the API returned an error: Array specified in enum field must match const values specified in oneOf field."
	mgr := newFixtureManager(userSchemaProperty, t.Name())
	resourceName := fmt.Sprintf("%s.test", userSchemaProperty)
	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkOktaUserSchemasDestroy,
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(`
			resource "okta_user_schema_property" "test" {
			  index       = "testAcc_replace_with_uuid"
			  title       = "terraform acceptance test"
			  type        = "array"
			  description = "testing"
			  master      = "OKTA"
			  scope       = "SELF"
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
			}`),
				Check: resource.ComposeTestCheckFunc(
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
			{
				Config: mgr.ConfigReplace(`
			resource "okta_user_schema_property" "test" {
			  index       = "testAcc_replace_with_uuid"
			  title       = "terraform acceptance test"
			  type        = "array"
			  description = "testing"
			  master      = "OKTA"
			  scope       = "SELF"
			  array_type  = "boolean"
			  array_enum  = ["false", "true"]
			  array_one_of {
			    title = "boolean FALSE"
			    const = "false"
			  }
			  array_one_of {
			    title = "boolean TRUE"
			    const = "true"
			  }
			}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "array"),
					resource.TestCheckResourceAttr(resourceName, "array_type", "boolean"),
					resource.TestCheckResourceAttr(resourceName, "array_enum.0", "false"),
					resource.TestCheckResourceAttr(resourceName, "array_enum.1", "true"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.0.title", "boolean FALSE"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.0.const", "false"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.1.title", "boolean TRUE"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.1.const", "true"),
				),
			},
		},
	})
}

func TestAccResourceOktaUserSchema_enum_boolean(t *testing.T) {
	t.Skip("TODO deal with apparent monolith bug")
	// TODO deal with apparent monolith bug:
	// "the API returned an error: Array specified in enum field must match const values specified in oneOf field."
	mgr := newFixtureManager(userSchemaProperty, t.Name())
	resourceName := fmt.Sprintf("%s.test", userSchemaProperty)
	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkOktaUserSchemasDestroy,
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(`
			resource "okta_user_schema_property" "test" {
			  index       = "testAcc_replace_with_uuid"
			  title       = "terraform acceptance test"
			  type        = "boolean"
			  description = "testing"
			  master      = "OKTA"
			  scope       = "SELF"
			  enum  = ["true", "false"]
			  one_of {
			    title = "boolean True"
			    const = "true"
			  }
			  one_of {
			    title = "boolean False"
			    const = "false"
			  }
			}`),
				Check: resource.ComposeTestCheckFunc(
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
			{
				Config: mgr.ConfigReplace(`
			resource "okta_user_schema_property" "test" {
			  index       = "testAcc_replace_with_uuid"
			  title       = "terraform acceptance test"
			  type        = "boolean"
			  description = "testing"
			  master      = "OKTA"
			  scope       = "SELF"
			  enum  = ["false", "true"]
			  one_of {
			    title = "boolean FALSE"
			    const = "false"
			  }
			  one_of {
			    title = "boolean TRUE"
			    const = "true"
			  }
			}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "boolean"),
					resource.TestCheckResourceAttr(resourceName, "enum.0", "false"),
					resource.TestCheckResourceAttr(resourceName, "enum.1", "true"),
					resource.TestCheckResourceAttr(resourceName, "one_of.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "one_of.0.title", "boolean FALSE"),
					resource.TestCheckResourceAttr(resourceName, "one_of.0.const", "false"),
					resource.TestCheckResourceAttr(resourceName, "one_of.1.title", "boolean TRUE"),
					resource.TestCheckResourceAttr(resourceName, "one_of.1.const", "true"),
				),
			},
		},
	})
}

func TestAccResourceOktaUserSchema_array_enum_string(t *testing.T) {
	mgr := newFixtureManager(userSchemaProperty, t.Name())
	resourceName := fmt.Sprintf("%s.test", userSchemaProperty)
	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkOktaUserSchemasDestroy,
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(`
			resource "okta_user_schema_property" "test" {
			  index       = "testAcc_replace_with_uuid"
			  title       = "terraform acceptance test"
			  type        = "array"
			  description = "testing"
			  master      = "OKTA"
			  scope       = "SELF"
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
			}`),
				Check: resource.ComposeTestCheckFunc(
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
			{
				Config: mgr.ConfigReplace(`
			resource "okta_user_schema_property" "test" {
			  index       = "testAcc_replace_with_uuid"
			  title       = "terraform acceptance test"
			  type        = "array"
			  description = "testing"
			  master      = "OKTA"
			  scope       = "SELF"
			  array_type  = "string"
			  array_enum  = ["ONE", "TWO", "THREE"]
			  array_one_of {
			    title = "STRING ONE"
			    const = "ONE"
			  }
			  array_one_of {
			    title = "STRING TWO"
			    const = "TWO"
			  }
			  array_one_of {
			    title = "STRING THREE"
			    const = "THREE"
			  }
			}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "array"),
					resource.TestCheckResourceAttr(resourceName, "array_type", "string"),
					resource.TestCheckResourceAttr(resourceName, "array_enum.0", "ONE"),
					resource.TestCheckResourceAttr(resourceName, "array_enum.1", "TWO"),
					resource.TestCheckResourceAttr(resourceName, "array_enum.2", "THREE"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.0.title", "STRING ONE"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.0.const", "ONE"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.1.title", "STRING TWO"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.1.const", "TWO"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.2.title", "STRING THREE"),
					resource.TestCheckResourceAttr(resourceName, "array_one_of.2.const", "THREE"),
				),
			},
		},
	})
}

func TestAccResourceOktaUserSchema_enum_string(t *testing.T) {
	mgr := newFixtureManager(userSchemaProperty, t.Name())
	resourceName := fmt.Sprintf("%s.test", userSchemaProperty)
	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkOktaUserSchemasDestroy,
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(`
			resource "okta_user_schema_property" "test" {
			  index       = "testAcc_replace_with_uuid"
			  title       = "terraform acceptance test"
			  description = "testing"
			  master      = "OKTA"
			  scope       = "SELF"
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
			}`),
				Check: resource.ComposeTestCheckFunc(
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
			{
				Config: mgr.ConfigReplace(`
			resource "okta_user_schema_property" "test" {
			  index       = "testAcc_replace_with_uuid"
			  title       = "terraform acceptance test"
			  description = "testing"
			  master      = "OKTA"
			  scope       = "SELF"
			  type  = "string"
			  enum  = ["ONE", "TWO", "THREE"]
			  one_of {
			    title = "STRING ONE"
			    const = "ONE"
			  }
			  one_of {
			    title = "STRING TWO"
			    const = "TWO"
			  }
			  one_of {
			    title = "STRING THREE"
			    const = "THREE"
			  }
			}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "string"),
					resource.TestCheckResourceAttr(resourceName, "enum.0", "ONE"),
					resource.TestCheckResourceAttr(resourceName, "enum.1", "TWO"),
					resource.TestCheckResourceAttr(resourceName, "enum.2", "THREE"),
					resource.TestCheckResourceAttr(resourceName, "one_of.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "one_of.0.title", "STRING ONE"),
					resource.TestCheckResourceAttr(resourceName, "one_of.0.const", "ONE"),
					resource.TestCheckResourceAttr(resourceName, "one_of.1.title", "STRING TWO"),
					resource.TestCheckResourceAttr(resourceName, "one_of.1.const", "TWO"),
					resource.TestCheckResourceAttr(resourceName, "one_of.2.title", "STRING THREE"),
					resource.TestCheckResourceAttr(resourceName, "one_of.2.const", "THREE"),
				),
			},
		},
	})
}

func testOktaUserSchemasExists(resourceName string) resource.TestCheckFunc {
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
		exists, err := testUserSchemaPropertyExists(schemaUserType, rs.Primary.ID, customSchema)
		if err != nil {
			return fmt.Errorf("failed to find: %v", err)
		}
		if !exists {
			return fmt.Errorf("custom property %s does not exist in a user profile subschema", rs.Primary.ID)
		}
		return nil
	}
}

func testUserSchemaPropertyExists(schemaUserType, index, resolutionScope string) (bool, error) {
	client := sdkV2ClientForTest()
	typeSchemaID, err := getUserTypeSchemaID(context.Background(), client, schemaUserType)
	if err != nil {
		return false, err
	}
	us, _, err := client.UserSchema.GetUserSchema(context.Background(), typeSchemaID)
	if err != nil {
		return false, fmt.Errorf("failed to get user schema: %v", err)
	}
	switch resolutionScope {
	case baseSchema:
		bp := userSchemaBaseAttribute(us, index)
		return bp != nil, nil
	case customSchema:
		bp := userSchemaCustomAttribute(us, index)
		return bp != nil, nil
	default:
		return false, fmt.Errorf("resolution scope can be only 'base' or 'custom'")
	}
}

// TestAccResourceOktaUserSchema_parallel_api_calls test coverage to ensure
// backoff in create and update for okta_ser_schema_property resource is
// operating correctly.
func TestAccResourceOktaUserSchema_parallel_api_calls(t *testing.T) {
	mgr := newFixtureManager(userSchemaProperty, t.Name())
	config := `
resource "okta_user_schema_property" "one" {
	index       = "testAcc_replace_with_uuid_one"
	title       = "one"
	type        = "string"
	permissions = "%s"
}
resource "okta_user_schema_property" "two" {
	index       = "testAcc_replace_with_uuid_two"
	title       = "two"
	type        = "string"
	permissions = "%s"
}
resource "okta_user_schema_property" "three" {
	index       = "testAcc_replace_with_uuid_three"
	title       = "three"
	type        = "string"
	permissions = "%s"
}
resource "okta_user_schema_property" "four" {
	index       = "testAcc_replace_with_uuid_four"
	title       = "four"
	type        = "string"
	permissions = "%s"
}
resource "okta_user_schema_property" "five" {
	index       = "testAcc_replace_with_uuid_five"
	title       = "five"
	type        = "string"
	permissions = "%s"
}
`
	ro := make([]interface{}, 5)
	for i := 0; i < 5; i++ {
		ro[i] = "READ_ONLY"
	}
	rw := make([]interface{}, 5)
	for i := 0; i < 5; i++ {
		rw[i] = "READ_WRITE"
	}
	roConfig := mgr.ConfigReplace(fmt.Sprintf(config, ro...))
	rwConfig := mgr.ConfigReplace(fmt.Sprintf(config, rw...))
	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: roConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("okta_user_schema_property.one", "permissions", "READ_ONLY"),
					resource.TestCheckResourceAttr("okta_user_schema_property.two", "permissions", "READ_ONLY"),
					resource.TestCheckResourceAttr("okta_user_schema_property.three", "permissions", "READ_ONLY"),
					resource.TestCheckResourceAttr("okta_user_schema_property.four", "permissions", "READ_ONLY"),
					resource.TestCheckResourceAttr("okta_user_schema_property.five", "permissions", "READ_ONLY"),
				),
			},
			{
				Config: rwConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("okta_user_schema_property.one", "permissions", "READ_WRITE"),
					resource.TestCheckResourceAttr("okta_user_schema_property.two", "permissions", "READ_WRITE"),
					resource.TestCheckResourceAttr("okta_user_schema_property.three", "permissions", "READ_WRITE"),
					resource.TestCheckResourceAttr("okta_user_schema_property.four", "permissions", "READ_WRITE"),
					resource.TestCheckResourceAttr("okta_user_schema_property.five", "permissions", "READ_WRITE"),
				),
			},
		},
	})
}

func checkOktaUserSchemasDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		schemaUserType := "default"
		if rs.Primary.Attributes["user_type"] != "" {
			schemaUserType = rs.Primary.Attributes["user_type"]
		}
		exists, _ := testUserSchemaPropertyExists(schemaUserType, rs.Primary.ID, customSchema)
		if exists {
			return fmt.Errorf("resource still exists, ID: %s", rs.Primary.ID)
		}
	}
	return nil
}

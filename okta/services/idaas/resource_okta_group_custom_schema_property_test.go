package idaas_test

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/provider"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/services/idaas"
)

func TestAccResourceOktaGroupSchema_crud(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSGroupSchemaProperty, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updated := mgr.GetFixtures("basic_updated.tf", t)
	// unique := mgr.GetFixtures("unique.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSGroupSchemaProperty)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:          acctest.AccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: acctest.AccProvidersFactoriesForTest(),
		CheckDestroy:      checkOktaGroupSchemasDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testOktaGroupSchemasExists(t.Name(), resourceName),
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
					testOktaGroupSchemasExists(t.Name(), resourceName),
					resource.TestCheckResourceAttr(resourceName, "index", "testAcc_"+strconv.Itoa(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "title", "terraform acceptance test updated 002"),
					resource.TestCheckResourceAttr(resourceName, "type", "string"),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform acceptance test updated 002"),
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
					resource.TestCheckResourceAttr(resourceName, "scope", "NONE"),
				),
			},

			// NOTE this test will fail because of a bug in the monolith on step 3
			// "You cannot add the attribute with the variable name
			// 'testAcc_1749602242782417788' because the deletion process for an
			// attribute with the same variable name is incomplete. Wait until the data
			// clean up process finishes and then try again."
			/*
				{
					Config: unique,
					Check: resource.ComposeTestCheckFunc(
						testOktaGroupSchemasExists(resourceName),
						resource.TestCheckResourceAttr(resourceName, "index", "testAcc_"+strconv.Itoa(ri)),
						resource.TestCheckResourceAttr(resourceName, "title", "terraform acceptance test setting unique attribute to UNIQUE_VALIDATED 006"),
						resource.TestCheckResourceAttr(resourceName, "type", "string"),
						resource.TestCheckResourceAttr(resourceName, "description", "terraform acceptance test setting unique attribute to UNIQUE_VALIDATED 006"),
						resource.TestCheckResourceAttr(resourceName, "required", "true"),
						resource.TestCheckResourceAttr(resourceName, "min_length", "1"),
						resource.TestCheckResourceAttr(resourceName, "max_length", "70"),
						resource.TestCheckResourceAttr(resourceName, "permissions", "READ_WRITE"),
						resource.TestCheckResourceAttr(resourceName, "master", "OKTA"),
						resource.TestCheckResourceAttr(resourceName, "unique", "UNIQUE_VALIDATED"),
					),
				},
			*/
		},
	})
}

func TestAccResourceOktaGroupSchema_arrayString(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSGroupSchemaProperty)
	mgr := newFixtureManager("resources", resources.OktaIDaaSGroupSchemaProperty, t.Name())
	config := mgr.GetFixtures("array_string.tf", t)
	updatedConfig := mgr.GetFixtures("array_string_updated.tf", t)
	arrayEnum := mgr.GetFixtures("array_enum.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:          acctest.AccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: acctest.AccProvidersFactoriesForTest(),
		CheckDestroy:      checkOktaGroupSchemasDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testOktaGroupSchemasExists(t.Name(), resourceName),
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
					testOktaGroupSchemasExists(t.Name(), resourceName),
					resource.TestCheckResourceAttr(resourceName, "index", "testAcc_"+strconv.Itoa(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "title", "terraform acceptance test updated 003"),
					resource.TestCheckResourceAttr(resourceName, "type", "array"),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform acceptance test updated 003"),
					resource.TestCheckResourceAttr(resourceName, "array_type", "string"),
					resource.TestCheckResourceAttr(resourceName, "required", "true"),
					resource.TestCheckResourceAttr(resourceName, "permissions", "READ_WRITE"),
					resource.TestCheckResourceAttr(resourceName, "master", "OKTA"),
				),
			},
			{
				Config: arrayEnum,
				Check: resource.ComposeTestCheckFunc(
					testOktaGroupSchemasExists(t.Name(), resourceName),
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

func TestAccResourceOktaGroupSchema_array_enum_number(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSGroupSchemaProperty, t.Name())
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSGroupSchemaProperty)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:          acctest.AccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: acctest.AccProvidersFactoriesForTest(),
		CheckDestroy:      checkOktaGroupSchemasDestroy,
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(`
			resource "okta_group_schema_property" "test" {
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
			resource "okta_group_schema_property" "test" {
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

func TestAccResourceOktaGroupSchema_enum_number(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSGroupSchemaProperty, t.Name())
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSGroupSchemaProperty)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:          acctest.AccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: acctest.AccProvidersFactoriesForTest(),
		CheckDestroy:      checkOktaGroupSchemasDestroy,
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(`
			resource "okta_group_schema_property" "test" {
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
			resource "okta_group_schema_property" "test" {
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

func TestAccResourceOktaGroupSchema_array_enum_integer(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSGroupSchemaProperty, t.Name())
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSGroupSchemaProperty)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:          acctest.AccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: acctest.AccProvidersFactoriesForTest(),
		CheckDestroy:      checkOktaGroupSchemasDestroy,
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(`
			resource "okta_group_schema_property" "test" {
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
			resource "okta_group_schema_property" "test" {
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

func TestAccResourceOktaGroupSchema_enum_integer(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSGroupSchemaProperty, t.Name())
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSGroupSchemaProperty)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:          acctest.AccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: acctest.AccProvidersFactoriesForTest(),
		CheckDestroy:      checkOktaGroupSchemasDestroy,
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(`
			resource "okta_group_schema_property" "test" {
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
			resource "okta_group_schema_property" "test" {
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

func TestAccResourceOktaGroupSchema_array_enum_boolean(t *testing.T) {
	t.Skip("TODO deal with apparent monolith bug")
	// TODO deal with apparent monolith bug:
	// "the API returned an error: Array specified in enum field must match const values specified in oneOf field."
	mgr := newFixtureManager("resources", resources.OktaIDaaSGroupSchemaProperty, t.Name())
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSGroupSchemaProperty)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:          acctest.AccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: acctest.AccProvidersFactoriesForTest(),
		CheckDestroy:      checkOktaGroupSchemasDestroy,
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(`
			resource "okta_group_schema_property" "test" {
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
			resource "okta_group_schema_property" "test" {
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

func TestAccResourceOktaGroupSchema_enum_boolean(t *testing.T) {
	t.Skip("TODO deal with apparent monolith bug")
	// TODO deal with apparent monolith bug:
	// "the API returned an error: Array specified in enum field must match const values specified in oneOf field."
	mgr := newFixtureManager("resources", resources.OktaIDaaSGroupSchemaProperty, t.Name())
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSGroupSchemaProperty)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:          acctest.AccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: acctest.AccProvidersFactoriesForTest(),
		CheckDestroy:      checkOktaGroupSchemasDestroy,
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(`
			resource "okta_group_schema_property" "test" {
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
			resource "okta_group_schema_property" "test" {
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

func TestAccResourceOktaGroupSchema_array_enum_string(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSGroupSchemaProperty, t.Name())
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSGroupSchemaProperty)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:          acctest.AccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: acctest.AccProvidersFactoriesForTest(),
		CheckDestroy:      checkOktaGroupSchemasDestroy,
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(`
			resource "okta_group_schema_property" "test" {
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
			resource "okta_group_schema_property" "test" {
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

func TestAccResourceOktaGroupSchema_enum_string(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSGroupSchemaProperty, t.Name())
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSGroupSchemaProperty)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:          acctest.AccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: acctest.AccProvidersFactoriesForTest(),
		CheckDestroy:      checkOktaGroupSchemasDestroy,
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(`
			resource "okta_group_schema_property" "test" {
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
			resource "okta_group_schema_property" "test" {
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

// TestAccResourceOktaGroupSchema_parallel_api_calls test coverage to ensure
// backoff in create and update for okta_group_schema_property resource is
// operating correctly.
func TestAccResourceOktaGroupSchema_parallel_api_calls(t *testing.T) {
	if provider.SkipVCRTest(t) {
		return
	}
	mgr := newFixtureManager("resources", resources.OktaIDaaSGroupSchemaProperty, t.Name())
	config := `
resource "okta_group_schema_property" "one" {
	index       = "testAcc_replace_with_uuid_one"
	title       = "one"
	type        = "string"
	permissions = "%s"
}
resource "okta_group_schema_property" "two" {
	index       = "testAcc_replace_with_uuid_two"
	title       = "two"
	type        = "string"
	permissions = "%s"
}
resource "okta_group_schema_property" "three" {
	index       = "testAcc_replace_with_uuid_three"
	title       = "three"
	type        = "string"
	permissions = "%s"
}
resource "okta_group_schema_property" "four" {
	index       = "testAcc_replace_with_uuid_four"
	title       = "four"
	type        = "string"
	permissions = "%s"
}
resource "okta_group_schema_property" "five" {
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
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:          acctest.AccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		CheckDestroy:      nil,
		ProviderFactories: acctest.AccProvidersFactoriesForTest(),
		Steps: []resource.TestStep{
			{
				Config: roConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("okta_group_schema_property.one", "permissions", "READ_ONLY"),
					resource.TestCheckResourceAttr("okta_group_schema_property.two", "permissions", "READ_ONLY"),
					resource.TestCheckResourceAttr("okta_group_schema_property.three", "permissions", "READ_ONLY"),
					resource.TestCheckResourceAttr("okta_group_schema_property.four", "permissions", "READ_ONLY"),
					resource.TestCheckResourceAttr("okta_group_schema_property.five", "permissions", "READ_ONLY"),
				),
			},
			{
				Config: rwConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("okta_group_schema_property.one", "permissions", "READ_WRITE"),
					resource.TestCheckResourceAttr("okta_group_schema_property.two", "permissions", "READ_WRITE"),
					resource.TestCheckResourceAttr("okta_group_schema_property.three", "permissions", "READ_WRITE"),
					resource.TestCheckResourceAttr("okta_group_schema_property.four", "permissions", "READ_WRITE"),
					resource.TestCheckResourceAttr("okta_group_schema_property.five", "permissions", "READ_WRITE"),
				),
			},
		},
	})
}

func checkOktaGroupSchemasDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		exists, _ := testGroupSchemaPropertyExists(rs.Primary.ID)
		if exists {
			return fmt.Errorf("resource still exists, ID: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testGroupSchemaPropertyExists(index string) (bool, error) {
	client := provider.SdkV2ClientForTest()
	gs, _, err := client.GroupSchema.GetGroupSchema(context.Background())
	if err != nil {
		return false, fmt.Errorf("failed to get group schema: %v", err)
	}
	ca := idaas.GroupSchemaCustomAttribute(gs, index)
	return ca != nil, nil
}

func testOktaGroupSchemasExists(testName, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		exists, err := testGroupSchemaPropertyExists(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to find: %v", err)
		}
		if !exists {
			return fmt.Errorf("custom property %s does not exist in a group profile subschema", rs.Primary.ID)
		}
		return nil
	}
}

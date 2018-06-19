package okta

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOktaUserSchemas_subschemaCheck(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaUserSchemas(ri)
	updatedConfig := testOktaUserSchemas_subschemaCheck(ri)
	resourceName := "okta_user_schemas.test-" + strconv.Itoa(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOktaUserSchemasDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testOktaUserSchemasExists(resourceName),
				),
			},
			{
				Config:      updatedConfig,
				ExpectError: regexp.MustCompile("You cannot change the subschema field for an existing User SubSchema"),
				PlanOnly:    true,
				Check: resource.ComposeTestCheckFunc(
					testOktaUserSchemasExists(resourceName),
				),
			},
		},
	})
}

func TestAccOktaUserSchemas_indexCheck(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaUserSchemas(ri)
	updatedConfig := testOktaUserSchemas_indexCheck(ri)
	resourceName := "okta_user_schemas.test-" + strconv.Itoa(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOktaUserSchemasDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testOktaUserSchemasExists(resourceName),
				),
			},
			{
				Config:      updatedConfig,
				ExpectError: regexp.MustCompile("You cannot change the index field for an existing User SubSchema"),
				PlanOnly:    true,
				Check: resource.ComposeTestCheckFunc(
					testOktaUserSchemasExists(resourceName),
				),
			},
		},
	})
}

func TestAccOktaUserSchemas_typeCheck(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaUserSchemas(ri)
	updatedConfig := testOktaUserSchemas_typeCheck(ri)
	resourceName := "okta_user_schemas.test-" + strconv.Itoa(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOktaUserSchemasDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testOktaUserSchemasExists(resourceName),
				),
			},
			{
				Config:      updatedConfig,
				ExpectError: regexp.MustCompile("You cannot change the type field for an existing User SubSchema"),
				PlanOnly:    true,
				Check: resource.ComposeTestCheckFunc(
					testOktaUserSchemasExists(resourceName),
				),
			},
		},
	})
}

func TestAccOktaUserSchemas_arrayTypeDeny(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaUserSchemas_arrayTypeDeny(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("arraytype field not required if type field is not array"),
				PlanOnly:    true,
			},
		},
	})
}

func TestAccOktaUserSchemas_enumValid(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaUserSchemas_enumValid(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("enum field only valid if SubSchema type is string"),
				PlanOnly:    true,
			},
		},
	})
}

func TestAccOktaUserSchemas_oneofValid(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaUserSchemas_oneofValid(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("oneof field only valid if enum is defined"),
				PlanOnly:    true,
			},
		},
	})
}

func TestAccOktaUserSchemas_arrayTypeRequire(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaUserSchemas_arrayTypeRequire(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("arraytype field required if type field is array"),
				PlanOnly:    true,
			},
		},
	})
}

func TestAccOktaUserSchemas(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := "okta_user_schemas.test-" + strconv.Itoa(ri)
	config := testOktaUserSchemas(ri)
	updatedConfig := testOktaUserSchemas_updated(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOktaUserSchemasDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testOktaUserSchemasExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "subschema", "custom"),
					resource.TestCheckResourceAttr(resourceName, "index", "testAcc"+strconv.Itoa(ri)),
					resource.TestCheckResourceAttr(resourceName, "title", "terraform acceptance test"),
					resource.TestCheckResourceAttr(resourceName, "type", "string"),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform acceptance test"),
					resource.TestCheckResourceAttr(resourceName, "required", "false"),
					resource.TestCheckResourceAttr(resourceName, "minlength", "1"),
					resource.TestCheckResourceAttr(resourceName, "maxlength", "50"),
					resource.TestCheckResourceAttr(resourceName, "permissions", "READ_ONLY"),
					resource.TestCheckResourceAttr(resourceName, "master", "PROFILE_MASTER"),
					resource.TestCheckResourceAttr(resourceName, "enum.0", "S"),
					resource.TestCheckResourceAttr(resourceName, "enum.1", "M"),
					resource.TestCheckResourceAttr(resourceName, "enum.2", "L"),
					resource.TestCheckResourceAttr(resourceName, "enum.3", "XL"),
					resource.TestCheckResourceAttr(resourceName, "oneof", "[\n {\"const\": \"S\", \"title\": \"Small\"},\n {\"const\": \"M\", \"title\": \"Medium\"},\n {\"const\": \"L\", \"title\": \"Large\"},\n {\"const\": \"XL\", \"title\": \"Extra Large\"}\n]\n"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testOktaUserSchemasExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "subschema", "custom"),
					resource.TestCheckResourceAttr(resourceName, "index", "testAcc"+strconv.Itoa(ri)),
					resource.TestCheckResourceAttr(resourceName, "title", "terraform acceptance test updated"),
					resource.TestCheckResourceAttr(resourceName, "type", "string"),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform acceptance test updated"),
					resource.TestCheckResourceAttr(resourceName, "required", "true"),
					resource.TestCheckResourceAttr(resourceName, "minlength", "1"),
					resource.TestCheckResourceAttr(resourceName, "maxlength", "70"),
					resource.TestCheckResourceAttr(resourceName, "permissions", "READ_WRITE"),
					resource.TestCheckResourceAttr(resourceName, "master", "OKTA"),
					resource.TestCheckResourceAttr(resourceName, "enum.0", "S"),
					resource.TestCheckResourceAttr(resourceName, "enum.1", "M"),
					resource.TestCheckResourceAttr(resourceName, "enum.2", "L"),
					resource.TestCheckResourceAttr(resourceName, "enum.3", "XXL"),
					resource.TestCheckResourceAttr(resourceName, "oneof", "[\n {\"const\": \"S\", \"title\": \"Small\"},\n {\"const\": \"M\", \"title\": \"Medium\"},\n {\"const\": \"L\", \"title\": \"Large\"},\n {\"const\": \"XXL\", \"title\": \"Extra Extra Large\"}\n]\n"),
				),
			},
		},
	})
}

func TestAccOktaUserSchemas_arrayString(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := "okta_user_schemas.test-" + strconv.Itoa(ri)
	config := testOktaUserSchemas_arrayString(ri)
	updatedConfig := testOktaUserSchemas_arrayStringUpdated(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOktaUserSchemasDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testOktaUserSchemasExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "subschema", "custom"),
					resource.TestCheckResourceAttr(resourceName, "index", "testAcc"+strconv.Itoa(ri)),
					resource.TestCheckResourceAttr(resourceName, "title", "terraform acceptance test"),
					resource.TestCheckResourceAttr(resourceName, "type", "array"),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform acceptance test"),
					resource.TestCheckResourceAttr(resourceName, "arraytype", "string"),
					resource.TestCheckResourceAttr(resourceName, "required", "false"),
					resource.TestCheckResourceAttr(resourceName, "permissions", "READ_ONLY"),
					resource.TestCheckResourceAttr(resourceName, "master", "PROFILE_MASTER"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testOktaUserSchemasExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "subschema", "custom"),
					resource.TestCheckResourceAttr(resourceName, "index", "testAcc"+strconv.Itoa(ri)),
					resource.TestCheckResourceAttr(resourceName, "title", "terraform acceptance test updated"),
					resource.TestCheckResourceAttr(resourceName, "type", "array"),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform acceptance test updated"),
					resource.TestCheckResourceAttr(resourceName, "arraytype", "string"),
					resource.TestCheckResourceAttr(resourceName, "required", "true"),
					resource.TestCheckResourceAttr(resourceName, "permissions", "READ_WRITE"),
					resource.TestCheckResourceAttr(resourceName, "master", "OKTA"),
				),
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

		subschema, hasSchema := rs.Primary.Attributes["subschema"]
		if !hasSchema {
			return fmt.Errorf("Error: No subschema found in state")
		}
		index, hasIndex := rs.Primary.Attributes["index"]
		if !hasIndex {
			return fmt.Errorf("Error: no index found in state")
		}
		_, hasTitle := rs.Primary.Attributes["title"]
		if !hasTitle {
			return fmt.Errorf("Error: no title found in state")
		}
		_, hasType := rs.Primary.Attributes["type"]
		if !hasType {
			return fmt.Errorf("Error: no type found in state")
		}

		err := testUserSchemaExists(true, subschema, index)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func testOktaUserSchemasDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "okta_user_schemas" {
			continue
		}

		subschema, hasSchema := rs.Primary.Attributes["subschema"]
		if !hasSchema {
			return fmt.Errorf("Error: No subschema found in state")
		}
		index, hasIndex := rs.Primary.Attributes["index"]
		if !hasIndex {
			return fmt.Errorf("Error: no index found in state")
		}
		_, hasTitle := rs.Primary.Attributes["title"]
		if !hasTitle {
			return fmt.Errorf("Error: no title found in state")
		}
		_, hasType := rs.Primary.Attributes["type"]
		if !hasType {
			return fmt.Errorf("Error: no type found in state")
		}

		err := testUserSchemaExists(false, subschema, index)
		if err != nil {
			return err
		}
	}
	return nil
}

func testUserSchemaExists(expected bool, scope string, index string) error {
	client := testAccProvider.Meta().(*Config).oktaClient

	exists := false
	subschemas, _, err := client.Schemas.GetUserSubSchemaIndex(scope)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Listing User Subschemas in Okta: %v", err)
	}
	for _, key := range subschemas {
		if key == index {
			exists = true
			break
		}
	}

	if expected == true && exists == false {
		return fmt.Errorf("[ERROR] User Schema %v not found in Okta", index)
	} else if expected == false && exists == true {
		return fmt.Errorf("[ERROR] User Schema %v still exists in Okta", index)
	}
	return nil
}

func testOktaUserSchemas(rInt int) string {
	return fmt.Sprintf(`
resource "okta_user_schemas" "test-%d" {
  subschema   = "custom"
  index       = "testAcc%d"
  title       = "terraform acceptance test"
  type        = "string"
  description = "terraform acceptance test"
  required    = false
  minlength   = 1
  maxlength   = 50
  permissions = "READ_ONLY"
  master      = "PROFILE_MASTER"
  enum        = [ "S","M","L","XL" ]
  oneof = <<JSON
[
 {"const": "S", "title": "Small"},
 {"const": "M", "title": "Medium"},
 {"const": "L", "title": "Large"},
 {"const": "XL", "title": "Extra Large"}
]
JSON
}
`, rInt, rInt)
}

func testOktaUserSchemas_updated(rInt int) string {
	return fmt.Sprintf(`
resource "okta_user_schemas" "test-%d" {
  subschema   = "custom"
  index       = "testAcc%d"
  title       = "terraform acceptance test updated"
  type        = "string"
  description = "terraform acceptance test updated"
  required    = true
  minlength   = 1
  maxlength   = 70
  permissions = "READ_WRITE"
  master      = "OKTA"
  enum        = [ "S","M","L","XXL" ]
  oneof = <<JSON
[
 {"const": "S", "title": "Small"},
 {"const": "M", "title": "Medium"},
 {"const": "L", "title": "Large"},
 {"const": "XXL", "title": "Extra Extra Large"}
]
JSON
}
`, rInt, rInt)
}

func testOktaUserSchemas_arrayString(rInt int) string {
	return fmt.Sprintf(`
resource "okta_user_schemas" "test-%d" {
  subschema = "custom"
  index     = "testAcc%d"
  title     = "terraform acceptance test"
  type      = "array"
  description = "terraform acceptance test"
  arraytype = "string"
  required    = false
  permissions = "READ_ONLY"
  master      = "PROFILE_MASTER"
}
`, rInt, rInt)
}

func testOktaUserSchemas_arrayStringUpdated(rInt int) string {
	return fmt.Sprintf(`
resource "okta_user_schemas" "test-%d" {
  subschema = "custom"
  index     = "testAcc%d"
  title     = "terraform acceptance test updated"
  type      = "array"
  description = "terraform acceptance test updated"
  arraytype = "string"
  required    = true
  permissions = "READ_WRITE"
  master      = "OKTA"
}
`, rInt, rInt)
}

func testOktaUserSchemas_subschemaCheck(rInt int) string {
	return fmt.Sprintf(`
resource "okta_user_schemas" "test-%d" {
  subschema = "base"
  index     = "testAcc%d"
  title     = "terraform acceptance test"
  type      = "string"
}
`, rInt, rInt)
}

func testOktaUserSchemas_indexCheck(rInt int) string {
	return fmt.Sprintf(`
resource "okta_user_schemas" "test-%d" {
  subschema = "custom"
  index     = "testAccChanged%d"
  title     = "terraform acceptance test"
  type      = "string"
}
`, rInt, rInt)
}

func testOktaUserSchemas_typeCheck(rInt int) string {
	return fmt.Sprintf(`
resource "okta_user_schemas" "test-%d" {
  subschema = "custom"
  index     = "testAcc%d"
  title     = "terraform acceptance test"
  type      = "array"
}
`, rInt, rInt)
}

func testOktaUserSchemas_arrayTypeDeny(rInt int) string {
	return fmt.Sprintf(`
resource "okta_user_schemas" "test-%d" {
  subschema = "custom"
  index     = "testAcc%d"
  title     = "terraform acceptance test"
  type      = "string"
  arraytype = "string"
}
`, rInt, rInt)
}

func testOktaUserSchemas_arrayTypeRequire(rInt int) string {
	return fmt.Sprintf(`
resource "okta_user_schemas" "test-%d" {
  subschema = "custom"
  index     = "testAcc%d"
  title     = "terraform acceptance test"
  type      = "array"
}
`, rInt, rInt)
}

func testOktaUserSchemas_enumValid(rInt int) string {
	return fmt.Sprintf(`
resource "okta_user_schemas" "test-%d" {
  subschema = "custom"
  index     = "testAcc%d"
  title     = "terraform acceptance test"
  type      = "boolean"
  enum        = [ "S","M","L","XXL" ]
}
`, rInt, rInt)
}

func testOktaUserSchemas_oneofValid(rInt int) string {
	return fmt.Sprintf(`
resource "okta_user_schemas" "test-%d" {
  subschema = "custom"
  index     = "testAcc%d"
  title     = "terraform acceptance test"
  type      = "boolean"
  oneof = <<JSON
[
 {"const": "S", "title": "Small"},
 {"const": "M", "title": "Medium"},
 {"const": "L", "title": "Large"},
 {"const": "XL", "title": "Extra Large"}
]
JSON
}
`, rInt, rInt)
}

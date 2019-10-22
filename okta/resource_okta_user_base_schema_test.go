package okta

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

const baseTestProp = "firstName"

func sweepUserBaseSchema(client *testClient) error {
	_, _, err := client.artClient.Schemas.GetUserSchema()
	if err != nil {
		return err
	}
	var errorList []error

	return condenseError(errorList)
}

func TestAccOktaUserBaseSchema_crud(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(userBaseSchema)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updated := mgr.GetFixtures("updated.tf", ri, t)
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

func testOktaUserBaseSchemasExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if exists, _ := testUserBaseSchemaExists(rs.Primary.ID); !exists {
			return fmt.Errorf("Failed to find %s", rs.Primary.ID)
		}
		return nil
	}
}

func testUserBaseSchemaExists(index string) (bool, error) {
	client := getClientFromMetadata(testAccProvider.Meta())
	subschema, _, err := client.Schemas.GetUserSubSchemaIndex(baseSchema)
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

package okta

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

type (
	CheckUpstream func(string) (bool, error)
	mapIndexFunc  func(int, string) int
)

func ensureResourceExists(name string, checkUpstream CheckUpstream) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		missingErr := fmt.Errorf("resource not found: %s", name)
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return missingErr
		}

		ID := rs.Primary.ID
		exist, err := checkUpstream(ID)
		if err != nil {
			return err
		} else if !exist {
			return missingErr
		}

		return nil
	}
}

func createCheckResourceDestroy(typeName string, checkUpstream CheckUpstream) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != typeName {
				continue
			}

			ID := rs.Primary.ID
			exists, err := checkUpstream(ID)
			if err != nil {
				return err
			}

			if exists {
				return fmt.Errorf("resource still exists, ID: %s", ID)
			}
		}
		return nil
	}
}

// Composes a TestCheckFunc for a slice of strings.
func testCheckResourceSliceAttr(name string, field string, value []string) resource.TestCheckFunc {
	return composeSliceCheck(name, field, value, func(i int, field string) int {
		return i
	})
}

// Composes a TestCheckFunc for a slice of strings for TypeSet.
func testCheckResourceSliceAttrForSet(name string, field string, value []string) resource.TestCheckFunc {
	return composeSliceCheck(name, field, value, func(i int, val string) int {
		return schema.HashString(val)
	})
}

func composeSliceCheck(name string, field string, value []string, mapIndex mapIndexFunc) resource.TestCheckFunc {
	args := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(name, fmt.Sprintf("%s.#", field), strconv.Itoa(len(value))),
	}

	for i, val := range value {
		fieldName := fmt.Sprintf("%s.%d", field, mapIndex(i, val))
		args = append(args, resource.TestCheckResourceAttr(name, fieldName, val))
	}

	return resource.ComposeTestCheckFunc(args...)
}

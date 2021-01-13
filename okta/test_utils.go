package okta

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

type CheckUpstream func(string) (bool, error)

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

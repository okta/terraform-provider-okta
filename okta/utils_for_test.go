package okta

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

type checkUpstream func(string) (bool, error)

func ensureResourceExists(name string, checkUpstream checkUpstream) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if isVCRPlayMode() {
			return nil
		}
		missingErr := fmt.Errorf("resource not found: %s", name)
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return missingErr
		}
		exist, err := checkUpstream(rs.Primary.ID)
		if err != nil {
			return err
		} else if !exist {
			return missingErr
		}
		return nil
	}
}

func checkResourceDestroy(typeName string, checkUpstream checkUpstream) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if isVCRPlayMode() {
			return nil
		}
		for _, rs := range s.RootModule().Resources {
			if rs.Type != typeName {
				continue
			}
			exists, err := checkUpstream(rs.Primary.ID)
			if err != nil {
				return err
			}
			if exists {
				return fmt.Errorf("resource still exists, ID: %s", rs.Primary.ID)
			}
		}
		return nil
	}
}

func ensureResourceNotExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return nil
		}
		return fmt.Errorf("Resource found: %s", name)
	}
}

func condenseError(errorList []error) error {
	if len(errorList) < 1 {
		return nil
	}
	msgList := make([]string, len(errorList))
	for i, err := range errorList {
		if err != nil {
			msgList[i] = err.Error()
		}
	}
	return fmt.Errorf("series of errors occurred: %s", strings.Join(msgList, ", "))
}

const (
	ErrorCheckMissingPermission         = "You do not have permission to access the feature you are requesting"
	ErrorCheckCannotCreateSWA           = "Cannot create application instance template_swa"
	ErrorCheckCannotCreateBasicAuth     = "Cannot create application instance template_basic_auth"
	ErrorCheckCannotCreateSPS           = "Cannot create application instance template_sps"
	ErrorCheckCannotCreateAWSConole     = "Cannot create application instance aws_console"
	ErrorCheckCannotCreateSWAThreeField = "Cannot create application instance template_swa3field"
	ErrorCheckFFGroupMembershipRules    = "GROUP_MEMBERSHIP_RULES is not enabled"
	ErrorCheckFFMFAPolicy               = "Missing Required Feature Flag OKTA_MFA_POLICY"
	ErrorOnlyOIEOrgs                    = "for OIE Orgs only"
)

// testAccErrorChecks Intended for use with TF sdk TestCase ErrorCheck function.
// Ability to skip tests that have specific errors.
func testAccErrorChecks(t *testing.T) resource.ErrorCheckFunc {
	return func(err error) error {
		if err == nil {
			return nil
		}
		messages := []string{
			ErrorCheckMissingPermission,
			ErrorCheckCannotCreateSWA,
			ErrorCheckCannotCreateBasicAuth,
			ErrorCheckCannotCreateSPS,
			ErrorCheckCannotCreateAWSConole,
			ErrorCheckCannotCreateSWAThreeField,
			ErrorCheckFFGroupMembershipRules,
			ErrorCheckFFMFAPolicy,
		}
		for _, message := range messages {
			// if error check message containing matches the message it will
			// apply a skip to t
			if errorCheckMessageContaining(t, message, err) {
				return err
			}
		}

		// check for our error on resources that are OIE only but are running
		// against a classic test org
		if errorCheckOIEOnlyFeature(t, err) {
			return err
		}

		return err
	}
}

func errorCheckOIEOnlyFeature(t *testing.T, err error) bool {
	if strings.Contains(err.Error(), ErrorOnlyOIEOrgs) {
		t.Skipf("Attempt to run OIE feature test on a Classic org")
		return true
	}
	return false
}

func errorCheckMessageContaining(t *testing.T, message string, err error) bool {
	if err == nil {
		return false
	}

	errorMessage := err.Error()
	missingFlags := []string{}
	if message == ErrorCheckMissingPermission {
		missingFlags = append(missingFlags, "ADVANCED_SSO")
		missingFlags = append(missingFlags, "MAPPINGS_API")
	}
	if message == ErrorCheckCannotCreateSWA {
		missingFlags = append(missingFlags, "ALLOW_SWA")
	}
	if message == ErrorCheckCannotCreateBasicAuth {
		missingFlags = append(missingFlags, "ALLOW_SWA")
	}
	if message == ErrorCheckCannotCreateSPS {
		missingFlags = append(missingFlags, "ALLOW_SWA")
	}
	if message == ErrorCheckCannotCreateAWSConole {
		missingFlags = append(missingFlags, "ALLOW_SWA")
	}
	if message == ErrorCheckCannotCreateSWAThreeField {
		missingFlags = append(missingFlags, "ALLOW_SWA")
	}
	if message == ErrorCheckFFGroupMembershipRules {
		missingFlags = append(missingFlags, "GROUP_MEMBERSHIP_RULES")
	}
	if message == ErrorCheckFFMFAPolicy {
		missingFlags = append(missingFlags, "OKTA_MFA_POLICY")
	}
	if strings.Contains(errorMessage, message) {
		t.Skipf("Skipping test, org possibly missing flags:\n%s\nerror:\n%s", strings.Join(missingFlags, ", "), errorMessage)
		return true
	}

	return false
}

// allowLongRunningACCTest Test skip helper for allowing long running tests to
// be executed.
func allowLongRunningACCTest(t *testing.T) bool {
	envVar := "OKTA_ALLOW_LONG_RUNNING_ACC_TEST"
	allow := (os.Getenv(envVar) != "")
	if !allow {
		t.Skipf("%q not present, skipping test", envVar)
	}
	return allow
}

// orgAdminOnlyTest Test skip helper for tests that should only run with a token
// of org admin permissions, not super admin.
func orgAdminOnlyTest(t *testing.T) bool {
	envVar := "OKTA_API_TOKEN_ROLE"
	envVal := os.Getenv(envVar)
	allow := (envVal == "org-admin")
	if !allow {
		t.Skipf("%s=%s not, requires %q value, skipping test", envVar, envVal, "org-admin")
	}
	return allow
}

// testAttributeJSON Deep equal of the JSON at named resource attribute witht he
// expected JSON
func testAttributeJSON(name, attribute, expectedJSON string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}
		actualJSON := rs.Primary.Attributes[attribute]
		eq := areJSONStringsEqual(expectedJSON, actualJSON)
		if !eq {
			return fmt.Errorf("attribute '%s' in '%s' expected %q, got %q", attribute, name, expectedJSON, actualJSON)
		}
		return nil
	}
}

// sleepInSecondsForTest Add sleep in a test to allow for eventual consistency
// thanks github.com/hashicorp/terraform-provider-google/google/provider_test.go
func sleepInSecondsForTest(t int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if os.Getenv("OKTA_VCR_TF_ACC") != "play" {
			time.Sleep(time.Duration(t) * time.Second)
		}
		return nil
	}
}

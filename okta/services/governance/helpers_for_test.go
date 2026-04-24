package governance_test

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"io"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

const (
	ErrorCheckMissingPermission         = "You do not have permission to access the feature you are requesting"
	ErrorCheckCannotCreateSWA           = "Cannot create application instance template_swa"
	ErrorCheckCannotCreateBasicAuth     = "Cannot create application instance template_basic_auth"
	ErrorCheckCannotCreateSPS           = "Cannot create application instance template_sps"
	ErrorCheckCannotCreateAWSConole     = "Cannot create application instance aws_console"
	ErrorCheckCannotCreateSWAThreeField = "Cannot create application instance template_swa3field"
	ErrorCheckFFGroupMembershipRules    = "GROUP_MEMBERSHIP_RULES is not enabled"
	ErrorCheckFFMFAPolicy               = "Missing Required Feature Flag OKTA_MFA_POLICY"
	ErrorSelfServiceApplicationEnabled  = "Self service application assignment for organization managed apps must be enabled"
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
			ErrorSelfServiceApplicationEnabled,
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
	if message == ErrorSelfServiceApplicationEnabled {
		missingFlags = append(missingFlags, "Admin UI > Applications > Self Service > User App Requests > App Catalog Settings > Allow users to add org-managed apps (enabled)")
	}
	if strings.Contains(errorMessage, message) {
		t.Skipf("Skipping test, org possibly missing flags:\n%s\nerror:\n%s", strings.Join(missingFlags, ", "), errorMessage)
		return true
	}

	return false
}

type fixtureManager struct {
	Path     string
	Seed     int
	TestName string
}

const (
	uuidPattern = "replace_with_uuid"
)

// newFixtureManager Gets a new fixture manager for a particular resource.
func newFixtureManager(resourceType, resourceName, testName string) *fixtureManager {
	ri := acctest.RandInt()

	// If we are running in VCR mode make the random number be a hash of the
	// test name.
	if os.Getenv("OKTA_VCR_TF_ACC") != "" {
		h := fnv.New32a()
		h.Write([]byte(testName))
		ri = int(h.Sum32())
	}

	dir, _ := os.Getwd()
	return &fixtureManager{
		Path:     path.Join(dir, "../../../examples", resourceType, resourceName),
		TestName: testName,
		Seed:     ri,
	}
}

func (manager *fixtureManager) SeedStr() string {
	return fmt.Sprintf("%d", manager.Seed)
}

func (manager *fixtureManager) GetFixtures(fixtureName string, t *testing.T) string {
	file, err := os.Open(path.Join(manager.Path, fixtureName))
	if err != nil {
		t.Fatalf("failed to load terraform fixtures for ACC test, err: %v", err)
	}
	defer file.Close()
	var rawFile bytes.Buffer
	_, err = io.Copy(&rawFile, file)
	if err != nil {
		t.Fatalf("failed to load terraform fixtures for ACC test, err: %v", err)
	}
	tfConfig := rawFile.String()
	if strings.Count(tfConfig, uuidPattern) == 0 {
		return tfConfig
	}

	return manager.ConfigReplace(tfConfig)
}

func (manager *fixtureManager) ConfigReplace(tfConfig string) string {
	return strings.ReplaceAll(tfConfig, uuidPattern, fmt.Sprintf("%d", manager.Seed))
}

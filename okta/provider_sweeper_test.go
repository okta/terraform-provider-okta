package okta

import (
	"fmt"
	"strconv"
	"testing"

	articulateOkta "github.com/articulate/oktasdk-go/okta"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/okta/okta-sdk-golang/okta"
)

var testResourcePrefix = "testAcc"

// TestMain overridden main testing function. Package level BeforeAll and AfterAll.
// It also delineates between acceptance tests and unit tests
func TestMain(m *testing.M) {
	// Acceptance test sweepers necessary to prevent dangling resources
	setupSweeper(passwordPolicy, deletePasswordPolicies)
	setupSweeper(signOnPolicy, deleteSignOnPolicies)
	setupSweeper(signOnPolicyRule, deleteSignOnPolicyRules)
	setupSweeper(passwordPolicyRule, deletePasswordPolicyRules)
	// Cannot sweep application resources as there is a bug with listing applications.
	// setupSweeper(oAuthApp, deleteOAuthApps)
	resource.TestMain(m)
}

// Sets up sweeper to clean up dangling resources
func setupSweeper(resourceType string, del func(*articulateOkta.Client, *okta.Client) error) {
	resource.AddTestSweepers(resourceType, &resource.Sweeper{
		Name: resourceType,
		F: func(region string) error {
			articulateOktaClient, client, err := sharedClient(region)

			if err != nil {
				return err
			}

			return del(articulateOktaClient, client)
		},
	})
}

// Builds test specific resource name
func buildResourceFQN(resourceType string, testID int) string {
	return resourceType + "." + buildResourceName(testID)
}

func buildResourceName(testID int) string {
	return testResourcePrefix + "-" + strconv.Itoa(testID)
}

// sharedClient returns a common Okta Client for sweepers, which currently requires the original SDK and the official beta SDK
func sharedClient(region string) (*articulateOkta.Client, *okta.Client, error) {
	err := accPreCheck()
	if err != nil {
		return nil, nil, err
	}

	c, err := oktaConfig()
	if err != nil {
		return nil, nil, err
	}

	articulateClient, err := articulateOkta.NewClientWithDomain(nil, c.orgName, c.domain, c.apiToken)

	if err != nil {
		return nil, nil, fmt.Errorf("[ERROR] Error creating Articulate Okta client: %v", err)
	}

	orgURL := fmt.Sprintf("https://%v.%v", c.orgName, c.domain)

	config := okta.NewConfig().WithOrgUrl(orgURL).WithToken(c.apiToken)
	client := okta.NewClient(config, nil, nil)

	return articulateClient, client, nil
}

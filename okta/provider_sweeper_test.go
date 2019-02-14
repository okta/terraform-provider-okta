package okta

import (
	"fmt"
	"strconv"
	"testing"

	articulateOkta "github.com/articulate/oktasdk-go/okta"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/okta/okta-sdk-golang/okta"
)

type testClient struct {
	oktaClient    *okta.Client
	artClient     *articulateOkta.Client
	apiSupplement *ApiSupplement
}

var testResourcePrefix = "testAcc"

// TestMain overridden main testing function. Package level BeforeAll and AfterAll.
// It also delineates between acceptance tests and unit tests
func TestMain(m *testing.M) {
	// Acceptance test sweepers necessary to prevent dangling resources
	setupSweeper(passwordPolicy, deletePasswordPolicies)
	setupSweeper(signOnPolicy, deleteSignOnPolicies)
	setupSweeper(mfaPolicy, deleteMfaPolicies)
	setupSweeper(signOnPolicyRule, deleteSignOnPolicyRules)
	setupSweeper(passwordPolicyRule, deletePasswordPolicyRules)
	setupSweeper("okta_*_apps", deleteTestApps)
	setupSweeper(mfaPolicyRule, deleteMfaPolicyRules)
	setupSweeper(authServer, deleteAuthServers)
	resource.TestMain(m)
}

// Sets up sweeper to clean up dangling resources
func setupSweeper(resourceType string, del func(*testClient) error) {
	resource.AddTestSweepers(resourceType, &resource.Sweeper{
		Name: resourceType,
		F: func(region string) error {
			articulateOktaClient, client, apiSupplement, err := sharedClient(region)

			if err != nil {
				return err
			}

			return del(&testClient{client, articulateOktaClient, apiSupplement})
		},
	})
}

// Builds test specific resource name
func buildResourceFQN(resourceType string, testID int) string {
	return resourceType + "." + buildResourceName(testID)
}

func buildResourceName(testID int) string {
	return testResourcePrefix + "_" + strconv.Itoa(testID)
}

// sharedClient returns a common Okta Client for sweepers, which currently requires the original SDK and the official beta SDK
func sharedClient(region string) (*articulateOkta.Client, *okta.Client, *ApiSupplement, error) {
	err := accPreCheck()
	if err != nil {
		return nil, nil, nil, err
	}

	c, err := oktaConfig()
	if err != nil {
		return nil, nil, nil, err
	}

	articulateClient, err := articulateOkta.NewClientWithDomain(nil, c.orgName, c.domain, c.apiToken)

	if err != nil {
		return nil, nil, nil, fmt.Errorf("[ERROR] Error creating Articulate Okta client: %v", err)
	}

	orgURL := fmt.Sprintf("https://%v.%v", c.orgName, c.domain)

	config := okta.NewConfig().WithOrgUrl(orgURL).WithToken(c.apiToken).WithBackoff(true).WithRetries(5)
	client := okta.NewClient(config, nil, nil)
	api := &ApiSupplement{requestExecutor: client.GetRequestExecutor()}

	return articulateClient, client, api, nil
}

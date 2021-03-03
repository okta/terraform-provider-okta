package okta

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/terraform-provider-okta/sdk"
)

type testClient struct {
	oktaClient    *okta.Client
	apiSupplement *sdk.ApiSupplement
}

var testResourcePrefix = "testAcc"

// TestMain overridden main testing function. Package level BeforeAll and AfterAll.
// It also delineates between acceptance tests and unit tests
func TestMain(m *testing.M) {
	// Acceptance test sweepers necessary to prevent dangling resources
	setupSweeper(policyPassword, deletePasswordPolicies)
	setupSweeper(policySignOn, deleteSignOnPolicies)
	setupSweeper(policyRuleIdpDiscovery, deletePolicyRuleIdpDiscovery)
	setupSweeper(policyMfa, deleteMfaPolicies)
	setupSweeper(policyRuleSignOn, deleteSignOnPolicyRules)
	setupSweeper(policyRulePassword, deletePolicyRulePasswords)
	setupSweeper("okta_*_app", deleteTestApps)
	setupSweeper("okta_*_idp", deleteTestIdps)
	setupSweeper(policyRuleMfa, deleteMfaPolicyRules)
	setupSweeper(authServer, deleteAuthServers)
	setupSweeper(groupRule, sweepGroupRules)
	setupSweeper(oktaGroup, sweepGroups)
	setupSweeper(oktaUser, sweepUsers)
	setupSweeper(userSchema, sweepUserSchema)
	setupSweeper(userBaseSchema, sweepUserBaseSchema)
	setupSweeper(networkZone, sweepNetworkZones)
	setupSweeper(inlineHook, sweepInlineHooks)
	setupSweeper(userType, sweepUserTypes)

	// add zones sweeper
	resource.TestMain(m)
}

// Sets up sweeper to clean up dangling resources
func setupSweeper(resourceType string, del func(*testClient) error) {
	resource.AddTestSweepers(resourceType, &resource.Sweeper{
		Name: resourceType,
		F: func(_ string) error {
			client, apiSupplement, err := sharedClient()
			if err != nil {
				return err
			}

			return del(&testClient{oktaClient: client, apiSupplement: apiSupplement})
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
func sharedClient() (*okta.Client, *sdk.ApiSupplement, error) {
	err := accPreCheck()
	if err != nil {
		return nil, nil, err
	}
	c, err := oktaConfig()
	if err != nil {
		return nil, nil, err
	}
	orgURL := fmt.Sprintf("https://%v.%v", c.orgName, c.domain)
	_, client, err := okta.NewClient(
		context.Background(),
		okta.WithOrgUrl(orgURL),
		okta.WithToken(c.apiToken),
		okta.WithRateLimitMaxRetries(20),
	)
	if err != nil {
		return client, nil, err
	}
	api := &sdk.ApiSupplement{RequestExecutor: client.GetRequestExecutor()}

	return client, api, nil
}

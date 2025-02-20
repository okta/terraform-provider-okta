package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

// TestMain overridden main testing function. Package level BeforeAll and AfterAll.
// It also delineates between acceptance tests and unit tests
// FIXME TestMain in the correct location?
func TestMain(m *testing.M) {
	// TF_VAR_hostname allows the real hostname to be scripted into the config tests
	// see examples/okta_resource_set/basic.tf
	os.Setenv("TF_VAR_hostname", fmt.Sprintf("%s.%s", os.Getenv("OKTA_ORG_NAME"), os.Getenv("OKTA_BASE_URL")))

	// TODO move sweepers setup to a sweeprs package
	// NOTE: Acceptance test sweepers are necessary to prevent dangling
	// resources.
	// NOTE: Don't run sweepers if we are playing back VCR as nothing should be
	// going over the wire
	if os.Getenv("OKTA_VCR_TF_ACC") != "play" {
		setupSweeper("okta_*_app", sweepTestApps)
		setupSweeper("okta_*_idp", sweepTestIdps)
		// setupSweeper(policySignOn, sweepSignOnPolicies)
		setupSweeper(resources.OktaIDaaSAdminRoleCustom, sweepCustomRoles)
		setupSweeper(resources.OktaIDaaSAuthServer, sweepAuthServers)
		setupSweeper(resources.OktaIDaaSBehavior, sweepBehaviors)
		setupSweeper(resources.OktaIDaaSEmailCustomization, sweepEmailCustomization)
		setupSweeper(resources.OktaIDaaSGroupRule, sweepGroupRules)
		setupSweeper(resources.OktaIDaaSInlineHook, sweepInlineHooks)
		setupSweeper(resources.OktaIDaaSGroup, sweepGroups)
		setupSweeper(resources.OktaIDaaSGroupSchemaProperty, sweepGroupCustomSchema)
		setupSweeper(resources.OktaIDaaSLinkDefinition, sweepLinkDefinitions)
		setupSweeper(resources.OktaIDaaSLogStream, sweepLogStreams)
		setupSweeper(resources.OktaIDaaSNetworkZone, sweepNetworkZones)
		setupSweeper(resources.OktaIDaaSPolicyMfa, sweepPoliciesMFA)
		setupSweeper(resources.OktaIDaaSPolicyPassword, sweepPoliciesPassword)
		setupSweeper(resources.OktaIDaaSPolicyRuleIdpDiscovery, sweepPolicyRulesIdpDiscovery)
		setupSweeper(resources.OktaIDaaSPolicyRuleMfa, sweepPolicyRulesMFA)
		setupSweeper(resources.OktaIDaaSPolicyRulePassword, sweepPolicyRulesPassword)
		setupSweeper(resources.OktaIDaaSPolicyRuleSignOn, sweepPolicyRulesOktaSignOn)
		setupSweeper(resources.OktaIDaaSPolicySignOn, sweepPoliciesAccess)
		setupSweeper(resources.OktaIDaaSResourceSet, sweepResourceSets)
		setupSweeper(resources.OktaIDaaSUser, sweepUsers)
		setupSweeper(resources.OktaIDaaSUserSchemaProperty, sweepUserCustomSchema)
		setupSweeper(resources.OktaIDaaSUserType, sweepUserTypes)
	}

	resource.TestMain(m)
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	_ = Provider()
}

func TestProviderValidate(t *testing.T) {
	envKeys := []string{
		"OKTA_ACCESS_TOKEN",
		"OKTA_ALLOW_LONG_RUNNING_ACC_TEST",
		"OKTA_API_CLIENT_ID",
		"OKTA_API_PRIVATE_KEY",
		"OKTA_API_PRIVATE_KEY_ID",
		"OKTA_API_SCOPES",
		"OKTA_API_TOKEN",
		"OKTA_BASE_URL",
		"OKTA_DEFAULT",
		"OKTA_GROUP",
		"OKTA_HTTP_PROXY",
		"OKTA_ORG_NAME",
		"OKTA_UPDATE",
	}
	envVals := make(map[string]string)
	// save and clear OKTA env vars so config can be test cleanly
	for _, key := range envKeys {
		val := os.Getenv(key)
		if val == "" {
			continue
		}
		envVals[key] = val
		os.Unsetenv(key)
	}

	tests := []struct {
		name         string
		accessToken  string
		apiToken     string
		clientID     string
		privateKey   string
		privateKeyID string
		scopes       []interface{}
		expectError  bool
	}{
		{"simple pass", "", "", "", "", "", []interface{}{}, false},
		{"access_token pass", "accessToken", "", "", "", "", []interface{}{}, false},
		{"access_token fail 1", "accessToken", "apiToken", "", "", "", []interface{}{}, true},
		{"access_token fail 2", "accessToken", "", "cliendID", "", "", []interface{}{}, true},
		{"access_token fail 3", "accessToken", "", "", "privateKey", "", []interface{}{}, true},
		{"access_token fail 4", "accessToken", "", "", "", "", []interface{}{"scope1", "scope2"}, true},
		{"api_token pass", "", "apiToken", "", "", "", []interface{}{}, false},
		{"api_token fail 1", "accessToken", "apiToken", "", "", "", []interface{}{}, true},
		{"api_token fail 2", "", "apiToken", "clientID", "", "", []interface{}{}, true},
		{"api_token fail 3", "", "apiToken", "", "", "privateKey", []interface{}{}, true},
		{"api_token fail 4", "", "apiToken", "", "", "", []interface{}{"scope1", "scope2"}, true},
		{"client_id pass", "", "", "clientID", "", "", []interface{}{}, false},
		{"client_id fail 1", "accessToken", "", "clientID", "", "", []interface{}{}, true},
		{"client_id fail 2", "accessToken", "apiToken", "clientID", "", "", []interface{}{}, true},
		{"private_key pass", "", "", "", "privateKey", "", []interface{}{}, false},
		{"private_key fail 1", "accessToken", "", "", "privateKey", "", []interface{}{}, true},
		{"private_key fail 2", "", "apiToken", "", "privateKey", "", []interface{}{}, true},
		{"private_key_id pass", "", "", "", "", "privateKeyID", []interface{}{}, false},
		{"private_key_id fail 1", "", "apiToken", "", "", "privateKeyID", []interface{}{}, true},
		{"scopes pass", "", "", "", "", "", []interface{}{"scope1", "scope2"}, false},
		{"scopes fail 1", "accessToken", "", "", "", "", []interface{}{"scope1", "scope2"}, true},
		{"scopes fail 2", "", "apiToken", "", "", "", []interface{}{"scope1", "scope2"}, true},
	}

	for _, test := range tests {
		resourceConfig := map[string]interface{}{}
		if test.accessToken != "" {
			resourceConfig["access_token"] = test.accessToken
		}
		if test.apiToken != "" {
			resourceConfig["api_token"] = test.apiToken
		}
		if test.clientID != "" {
			resourceConfig["client_id"] = test.clientID
		}
		if test.privateKey != "" {
			resourceConfig["private_key"] = test.privateKey
		}
		if test.privateKeyID != "" {
			resourceConfig["private_key_id"] = test.privateKeyID
		}
		if len(test.scopes) > 0 {
			resourceConfig["scopes"] = test.scopes
		}

		config := terraform.NewResourceConfigRaw(resourceConfig)
		provider := Provider()
		err := provider.Validate(config)

		if test.expectError && err == nil {
			t.Errorf("test %q: expected error but received none", test.name)
		}
		if !test.expectError && err != nil {
			t.Errorf("test %q: did not expect error but received error: %+v", test.name, err)
			fmt.Println()
		}
	}

	for key, val := range envVals {
		os.Setenv(key, val)
	}
}

package idaas_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	schema_sdk "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/okta-sdk-golang/v4/okta"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/api"
	"github.com/okta/terraform-provider-okta/okta/config"
	"github.com/okta/terraform-provider-okta/okta/provider"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/services/idaas"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
	"github.com/okta/terraform-provider-okta/sdk/query"
)

var (
	sweeperLogger   hclog.Logger
	sweeperLogLevel hclog.Level

	iDaaSAPIClientForTestUtil api.OktaIDaaSClient
)

func init() {
	sweeperLogLevel = hclog.Warn
	if os.Getenv("TF_LOG") != "" {
		sweeperLogLevel = hclog.LevelFromString(os.Getenv("TF_LOG"))
	}
	sweeperLogger = hclog.New(&hclog.LoggerOptions{
		Level:      sweeperLogLevel,
		TimeFormat: "2006/01/02 03:04:05",
	})

	if os.Getenv("OKTA_VCR_TF_ACC") == "play" {
		os.Setenv("OKTA_API_TOKEN", "token")
		os.Setenv("OKTA_BASE_URL", "dne-okta.com")
		if os.Getenv("OKTA_VCR_CASSETTE") != "" {
			os.Setenv("OKTA_ORG_NAME", os.Getenv("OKTA_VCR_CASSETTE"))
		}
	}

	// Initialize the shared IDaaS client only when running acceptance
	// tests (TF_ACC=1). In unit-test or CI-lint runs where no Okta
	// credentials are available, skip initialization to avoid a panic from
	// calling Fatalf on a bare testing.T{}.
	if os.Getenv("TF_ACC") != "" {
		t := &testing.T{}
		iDaaSAPIClientForTestUtil = IDaaSClientForTest(t)
	}
}

// TestMain overridden main testing function. Package level BeforeAll and AfterAll.
// It also delineates between acceptance tests and unit tests
func TestMain(m *testing.M) {
	// TF_VAR_hostname allows the real hostname to be scripted into the config tests
	// see examples/resources/okta_resource_set/basic.tf
	if os.Getenv("TF_VAR_hostname") == "" {
		os.Setenv("TF_VAR_hostname", fmt.Sprintf("%s.%s", os.Getenv("OKTA_ORG_NAME"), os.Getenv("OKTA_BASE_URL")))
	}
	os.Setenv("TF_VAR_orgID", os.Getenv("OKTA_ORG_ID"))

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
		setupSweeper(resources.OktaIDaaSBrand, sweepBrands)
	}

	resource.TestMain(m)
}

func TestProvider(t *testing.T) {
	if err := provider.Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	_ = provider.Provider()
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
		provider := provider.Provider()
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

func logSweptResource(kind, id, nameOrLabel string) {
	sweeperLogger.Warn(fmt.Sprintf("sweeper found dangling %q %q %q", kind, id, nameOrLabel))
}

// TestRunForcedSweeper forces sweeping any tangling testAcc resources that it
// can find.
//
//	go clean -testcache && \
//	TF_LOG=warn OKTA_ACC_TEST_FORCE_SWEEPERS=1 TF_ACC=1 go test -tags unit -mod=readonly -test.v -run ^TestRunForcedSweeper$ ./okta
func TestRunForcedSweeper(t *testing.T) {
	if os.Getenv("OKTA_VCR_TF_ACC") != "" {
		t.Skip("forced sweeper is live and will never be run within VCR")
		return
	}
	if os.Getenv("OKTA_ACC_TEST_FORCE_SWEEPERS") == "" || os.Getenv("TF_ACC") == "" {
		t.Skipf("ENV vars %q and %q must not be blank to force running of the sweepers", "OKTA_ACC_TEST_FORCE_SWEEPERS", "TF_ACC")
		return
	}

	testClient := IDaaSClientForTest(t)

	sweepCustomRoles(testClient)
	sweepTestApps(testClient)
	sweepAuthServers(testClient)
	sweepBehaviors(testClient)
	sweepEmailCustomization(testClient)
	sweepGroupRules(testClient)
	sweepTestIdps(testClient)
	sweepInlineHooks(testClient)
	sweepGroups(testClient)
	sweepGroupCustomSchema(testClient)
	sweepLinkDefinitions(testClient)
	sweepLogStreams(testClient)
	sweepNetworkZones(testClient)
	sweepResourceSets(testClient)
	sweepUsers(testClient)
	sweepUserCustomSchema(testClient)
	sweepUserTypes(testClient)

	// policy rules clean up needs to occur before policies
	// policy rules
	sweepPolicyRulesAccess(testClient)
	sweepPolicyRulesIdpDiscovery(testClient)
	sweepPolicyRulesMFA(testClient)
	sweepPolicyRulesOauthAuthorization(testClient)
	sweepPolicyRulesOktaSignOn(testClient)
	sweepPolicyRulesPassword(testClient)
	sweepPolicyRulesProfileEnrollment(testClient)
	sweepPolicyRulesSignOn(testClient)

	// policies
	sweepPoliciesAccess(testClient)
	sweepPoliciesIDPDiscovery(testClient)
	sweepPoliciesMFA(testClient)
	sweepPoliciesOauthAuthorization(testClient)
	sweepPoliciesOktaSignOn(testClient)
	sweepPoliciesPassword(testClient)
	sweepPoliciesProfileEnrollment(testClient)
	sweepPoliciesSignOn(testClient)

	// brands
	sweepBrands(testClient)
}

// Sets up sweeper to clean up dangling resources
func setupSweeper(resourceType string, del func(api.OktaIDaaSClient) error) {
	resource.AddTestSweepers(resourceType, &resource.Sweeper{
		Name: resourceType,
		F: func(_ string) error {
			t := &testing.T{}
			testClient := IDaaSClientForTest(t)
			return del(testClient)
		},
	})
}

func sweepCustomRoles(client api.OktaIDaaSClient) error {
	var errorList []error
	customRoles, _, err := client.OktaSDKSupplementClient().ListCustomRoles(context.Background(), &query.Params{Limit: utils.DefaultPaginationLimit})
	if err != nil {
		return err
	}
	for _, role := range customRoles.Roles {
		if !strings.HasPrefix(role.Label, "testAcc_") {
			_, err := client.OktaSDKSupplementClient().DeleteCustomRole(context.Background(), role.Id)
			if err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource("custom role", role.Id, role.Label)
		}
	}
	return condenseError(errorList)
}

func sweepTestApps(client api.OktaIDaaSClient) error {
	appList, err := idaas.ListAppsV2(context.Background(), client.OktaSDKClientV2(), &idaas.AppFilters{LabelPrefix: acctest.ResourcePrefixForTest}, utils.DefaultPaginationLimit)
	if err != nil {
		return err
	}
	var warnings []string
	var errors []string
	for _, app := range appList {
		warn := fmt.Sprintf("failed to sweep an application, there may be dangling resources. ID %s, label %s", app.Id, app.Label)
		_, err := client.OktaSDKClientV2().Application.DeactivateApplication(context.Background(), app.Id)
		if err != nil {
			warnings = append(warnings, warn)
		}
		resp, err := client.OktaSDKClientV2().Application.DeleteApplication(context.Background(), app.Id)
		if utils.Is404(resp) {
			warnings = append(warnings, warn)
			continue
		} else if err != nil {
			errors = append(errors, fmt.Sprintf("app id: %q, error: %s", app.Id, err.Error()))
			continue
		}
		logSweptResource("app", app.Id, app.Name)
	}
	if len(warnings) > 0 {
		return fmt.Errorf("sweep failures: %s", strings.Join(warnings, ", "))
	}
	if len(errors) > 0 {
		return fmt.Errorf("sweep errors: %s", strings.Join(errors, ", "))
	}
	return nil
}

func sweepAuthServers(client api.OktaIDaaSClient) error {
	servers, _, err := client.OktaSDKClientV2().AuthorizationServer.ListAuthorizationServers(context.Background(), &query.Params{Q: acctest.ResourcePrefixForTest})
	if err != nil {
		return err
	}
	for _, s := range servers {
		if _, err := client.OktaSDKClientV2().AuthorizationServer.DeactivateAuthorizationServer(context.Background(), s.Id); err != nil {
			return err
		}
		if _, err := client.OktaSDKClientV2().AuthorizationServer.DeleteAuthorizationServer(context.Background(), s.Id); err != nil {
			return err
		}
		logSweptResource("authorization server", s.Id, s.Name)
	}
	return nil
}

func sweepBehaviors(client api.OktaIDaaSClient) error {
	var errorList []error
	listBehaviorDetectionRules := client.OktaSDKClientV5().BehaviorAPI.ListBehaviorDetectionRules(context.Background())
	behaviors, _, err := listBehaviorDetectionRules.Execute()
	if err != nil {
		return err
	}
	for _, behavior := range behaviors {
		type iBehavior interface {
			GetId() string
			GetName() string
		}
		ctx := context.Background()
		b := behavior.GetActualInstance().(iBehavior)
		deleteBehaviorDetectionRuleRequest := client.OktaSDKClientV5().BehaviorAPI.DeleteBehaviorDetectionRule(ctx, b.GetId())
		_, err := deleteBehaviorDetectionRuleRequest.Execute()
		if err != nil {
			errorList = append(errorList, err)
			continue
		}
		logSweptResource("behavior", b.GetId(), b.GetName())
	}
	return condenseError(errorList)
}

func sweepEmailCustomization(client api.OktaIDaaSClient) error {
	ctx := context.Background()
	brands, _, err := client.OktaSDKClientV3().CustomizationAPI.ListBrands(ctx).Execute()
	if err != nil {
		return err
	}
	for _, brand := range brands {
		templates, resp, err := client.OktaSDKClientV3().CustomizationAPI.ListEmailTemplates(ctx, brand.GetId()).Limit(int32(utils.DefaultPaginationLimit)).Execute()
		if err != nil {
			continue
		}
		for resp.HasNextPage() {
			var nextTemplates []okta.EmailTemplate
			resp, err = resp.Next(&nextTemplates)
			if err != nil {
				continue
			}
			templates = append(templates, nextTemplates...)
		}

		for _, template := range templates {
			_, _ = client.OktaSDKClientV3().CustomizationAPI.DeleteAllCustomizations(context.Background(), brand.GetId(), template.GetName()).Execute()
		}
	}

	return nil
}

func sweepGroupRules(client api.OktaIDaaSClient) error {
	var errorList []error
	// Should never need to deal with pagination
	rules, _, err := client.OktaSDKClientV2().Group.ListGroupRules(context.Background(), &query.Params{Limit: utils.DefaultPaginationLimit})
	if err != nil {
		return err
	}

	for _, s := range rules {
		if s.Status == idaas.StatusActive {
			if _, err := client.OktaSDKClientV2().Group.DeactivateGroupRule(context.Background(), s.Id); err != nil {
				errorList = append(errorList, err)
				continue
			}
		}
		if _, err := client.OktaSDKClientV2().Group.DeleteGroupRule(context.Background(), s.Id, nil); err != nil {
			errorList = append(errorList, err)
			continue
		}
		logSweptResource("group rule", s.Id, s.Name)
	}
	return condenseError(errorList)
}

func sweepTestIdps(client api.OktaIDaaSClient) error {
	providers, _, err := client.OktaSDKClientV2().IdentityProvider.ListIdentityProviders(context.Background(), &query.Params{Q: "testAcc_"})
	if err != nil {
		return err
	}
	for _, idp := range providers {
		_, err := client.OktaSDKClientV2().IdentityProvider.DeleteIdentityProvider(context.Background(), idp.Id)
		if err != nil {
			return err
		}
		logSweptResource("identity provider", idp.Id, idp.Name)

		if idp.Type == idaas.Saml2Idp {
			_, err := client.OktaSDKClientV2().IdentityProvider.DeleteIdentityProviderKey(context.Background(), idp.Protocol.Credentials.Trust.Kid)
			if err != nil {
				return err
			}
			logSweptResource("saml identity provider key", idp.Id, idp.Protocol.Credentials.Trust.Kid)
		}
	}
	return nil
}

func sweepInlineHooks(client api.OktaIDaaSClient) error {
	var errorList []error
	hooks, _, err := client.OktaSDKClientV2().InlineHook.ListInlineHooks(context.Background(), nil)
	if err != nil {
		return err
	}
	for _, hook := range hooks {
		if !strings.HasPrefix(hook.Name, acctest.ResourcePrefixForTest) {
			continue
		}
		if hook.Status == idaas.StatusActive {
			_, _, err = client.OktaSDKClientV2().InlineHook.DeactivateInlineHook(context.Background(), hook.Id)
			if err != nil {
				errorList = append(errorList, err)
			}
		}
		_, err = client.OktaSDKClientV2().InlineHook.DeleteInlineHook(context.Background(), hook.Id)
		if err != nil {
			errorList = append(errorList, err)
		}
		logSweptResource("inline hook", hook.Id, hook.Name)
	}
	return condenseError(errorList)
}

func sweepBrands(client api.OktaIDaaSClient) error {
	var errorList []error
	ctx := context.Background()
	brands, _, err := client.OktaSDKClientV3().CustomizationAPI.ListBrands(ctx).Execute()
	if err != nil {
		return err
	}

	for _, b := range brands {
		test := fmt.Sprintf("%s-", strings.ToLower(acctest.ResourceNamePrefixForTest))
		if !strings.Contains(*b.Name, test) {
			continue
		}
		if _, err := client.OktaSDKClientV3().CustomizationAPI.DeleteBrand(ctx, *b.Id).Execute(); err != nil {
			errorList = append(errorList, err)
			continue
		}
		logSweptResource("brand", *b.Id, *b.Name)
	}
	return condenseError(errorList)
}

func sweepGroups(client api.OktaIDaaSClient) error {
	var errorList []error
	// Should never need to deal with pagination, limit is 10,000 by default
	groups, _, err := client.OktaSDKClientV2().Group.ListGroups(context.Background(), &query.Params{Q: acctest.ResourcePrefixForTest})
	if err != nil {
		return err
	}

	for _, s := range groups {
		if _, err := client.OktaSDKClientV2().Group.DeleteGroup(context.Background(), s.Id); err != nil {
			errorList = append(errorList, err)
			continue
		}
		logSweptResource("group", s.Id, s.Profile.Name)
	}
	return condenseError(errorList)
}

func sweepGroupCustomSchema(client api.OktaIDaaSClient) error {
	schema, _, err := client.OktaSDKClientV2().GroupSchema.GetGroupSchema(context.Background())
	if err != nil {
		return err
	}
	for key := range schema.Definitions.Custom.Properties {
		if strings.HasPrefix(key, acctest.ResourcePrefixForTest) {
			custom := idaas.BuildCustomGroupSchema(key, nil)
			_, _, err = client.OktaSDKClientV2().GroupSchema.UpdateGroupSchema(context.Background(), *custom)
			if err != nil {
				return err
			}
			logSweptResource("update group schema", key, key)
		}
	}
	return nil
}

func sweepLinkDefinitions(client api.OktaIDaaSClient) error {
	var errorList []error
	linkedObjects, _, err := client.OktaSDKClientV2().LinkedObject.ListLinkedObjectDefinitions(context.Background())
	if err != nil {
		return err
	}
	for _, object := range linkedObjects {
		if strings.HasPrefix(object.Primary.Name, acctest.ResourcePrefixForTest) {
			if _, err := client.OktaSDKClientV2().LinkedObject.DeleteLinkedObjectDefinition(context.Background(), object.Primary.Name); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource("linked object definition", object.Primary.Name, object.Primary.Title)
		}
	}
	return condenseError(errorList)
}

func sweepLogStreams(client api.OktaIDaaSClient) error {
	var errorList []error
	streams, _, err := client.OktaSDKClientV3().LogStreamAPI.ListLogStreams(context.Background()).Execute()
	if err != nil {
		return err
	}
	for _, stream := range streams {
		var id, name string
		if stream.LogStreamAws != nil {
			name = stream.LogStreamAws.Name
			id = stream.LogStreamAws.Id
		}
		if stream.LogStreamSplunk != nil {
			name = stream.LogStreamSplunk.Name
			id = stream.LogStreamSplunk.Id
		}

		if strings.HasPrefix(name, acctest.ResourcePrefixForTest) {
			if _, err = client.OktaSDKClientV3().LogStreamAPI.DeleteLogStream(context.Background(), id).Execute(); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource("log stream", id, name)
		}
	}
	return condenseError(errorList)
}

func sweepNetworkZones(client api.OktaIDaaSClient) error {
	var errorList []error
	zones, _, err := client.OktaSDKClientV2().NetworkZone.ListNetworkZones(context.Background(), &query.Params{Limit: utils.DefaultPaginationLimit})
	if err != nil {
		return err
	}
	for _, zone := range zones {
		if strings.HasPrefix(zone.Name, acctest.ResourcePrefixForTest) {
			if _, err := client.OktaSDKClientV2().NetworkZone.DeleteNetworkZone(context.Background(), zone.Id); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource("network zone", zone.Id, zone.Name)
		}
	}
	return condenseError(errorList)
}

func sweepPoliciesAccess(client api.OktaIDaaSClient) error {
	return sweepPolicyByType(sdk.AccessPolicyType, client)
}

func sweepPolicyRulesAccess(client api.OktaIDaaSClient) error {
	return sweepPolicyRulesByType(sdk.AccessPolicyType, client)
}

func sweepPoliciesIDPDiscovery(client api.OktaIDaaSClient) error {
	return sweepPolicyByType(sdk.IdpDiscoveryType, client)
}

func sweepPolicyRulesIdpDiscovery(client api.OktaIDaaSClient) error {
	return sweepPolicyRulesByType(sdk.IdpDiscoveryType, client)
}

func sweepPoliciesMFA(client api.OktaIDaaSClient) error {
	return sweepPolicyByType(sdk.MfaPolicyType, client)
}

func sweepPolicyRulesMFA(client api.OktaIDaaSClient) error {
	return sweepPolicyRulesByType(sdk.MfaPolicyType, client)
}

func sweepPoliciesOauthAuthorization(client api.OktaIDaaSClient) error {
	return sweepPolicyByType(sdk.OauthAuthorizationPolicyType, client)
}

func sweepPolicyRulesOauthAuthorization(client api.OktaIDaaSClient) error {
	return sweepPolicyRulesByType(sdk.OauthAuthorizationPolicyType, client)
}

func sweepPoliciesOktaSignOn(client api.OktaIDaaSClient) error {
	return sweepPolicyByType(sdk.SignOnPolicyType, client)
}

func sweepPolicyRulesOktaSignOn(client api.OktaIDaaSClient) error {
	return sweepPolicyRulesByType(sdk.SignOnPolicyType, client)
}

func sweepPoliciesPassword(client api.OktaIDaaSClient) error {
	return sweepPolicyByType(sdk.PasswordPolicyType, client)
}

func sweepPolicyRulesPassword(client api.OktaIDaaSClient) error {
	return sweepPolicyRulesByType(sdk.PasswordPolicyType, client)
}

func sweepPoliciesProfileEnrollment(client api.OktaIDaaSClient) error {
	return sweepPolicyByType(sdk.ProfileEnrollmentPolicyType, client)
}

func sweepPolicyRulesProfileEnrollment(client api.OktaIDaaSClient) error {
	return sweepPolicyRulesByType(sdk.ProfileEnrollmentPolicyType, client)
}

func sweepPoliciesSignOn(client api.OktaIDaaSClient) error {
	return sweepPolicyByType(sdk.SignOnPolicyRuleType, client)
}

func sweepPolicyRulesSignOn(client api.OktaIDaaSClient) error {
	return sweepPolicyRulesByType(sdk.SignOnPolicyRuleType, client)
}

func sweepResourceSets(client api.OktaIDaaSClient) error {
	var errorList []error
	resourceSets, _, err := client.OktaSDKSupplementClient().ListResourceSets(context.Background())
	if err != nil {
		return err
	}
	for _, b := range resourceSets.ResourceSets {
		if !strings.HasPrefix(b.Label, "testAcc_") {
			if _, err := client.OktaSDKSupplementClient().DeleteResourceSet(context.Background(), b.Id); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource("resource set", b.Id, b.Label)
		}
	}
	return condenseError(errorList)
}

func sweepUsers(client api.OktaIDaaSClient) error {
	var errorList []error
	users, resp, err := client.OktaSDKClientV2().User.ListUsers(context.Background(), &query.Params{Limit: 200, Q: acctest.ResourcePrefixForTest})
	if err != nil {
		return err
	}
	for resp.HasNextPage() {
		var nextUsers []*sdk.User
		resp, err = resp.Next(context.Background(), &nextUsers)
		if err != nil {
			return err
		}
		users = append(users, nextUsers...)
	}

	for _, u := range users {
		if err := idaas.EnsureUserDelete(context.Background(), u.Id, u.Status, client.OktaSDKClientV2()); err != nil {
			errorList = append(errorList, err)
			continue
		}
		var label string
		for k, v := range *u.Profile {
			label += fmt.Sprintf("%s:%+v, ", k, v)
		}
		logSweptResource("user", u.Id, label)
	}
	return condenseError(errorList)
}

func sweepUserCustomSchema(client api.OktaIDaaSClient) error {
	userTypes, _, err := client.OktaSDKClientV2().UserType.ListUserTypes(context.Background())
	if err != nil {
		return err
	}
	for _, userType := range userTypes {
		typeSchemaID := idaas.UserTypeSchemaID(userType)
		schema, _, err := client.OktaSDKClientV2().UserSchema.GetUserSchema(context.Background(), typeSchemaID)
		if err != nil {
			return err
		}
		for key := range schema.Definitions.Custom.Properties {
			if strings.HasPrefix(key, acctest.ResourcePrefixForTest) {
				custom := idaas.BuildCustomUserSchema(key, nil)
				_, _, err = client.OktaSDKClientV2().UserSchema.UpdateUserProfile(context.Background(), typeSchemaID, *custom)
				if err != nil {
					return err
				}
				logSweptResource("custom schema", typeSchemaID, "-")
			}
		}
	}
	return nil
}

func sweepUserTypes(client api.OktaIDaaSClient) error {
	userTypeList, _, _ := client.OktaSDKClientV2().UserType.ListUserTypes(context.Background())
	var errorList []error
	for _, ut := range userTypeList {
		if strings.HasPrefix(ut.Name, acctest.ResourcePrefixForTest) {
			if _, err := client.OktaSDKClientV2().UserType.DeleteUserType(context.Background(), ut.Id); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource("user type", ut.Id, ut.Name)
		}
	}
	return condenseError(errorList)
}

func sweepPolicyByType(t string, client api.OktaIDaaSClient) error {
	ctx := context.Background()
	policies, _, err := client.OktaSDKClientV2().Policy.ListPolicies(ctx, &query.Params{Type: t})
	if err != nil {
		return fmt.Errorf("failed to list policies in order to properly destroy: %v", err)
	}
	for _, _policy := range policies {
		policy := _policy.(*sdk.Policy)
		if strings.HasPrefix(policy.Name, acctest.ResourcePrefixForTest) {
			_, err = client.OktaSDKClientV2().Policy.DeletePolicy(ctx, policy.Id)
			if err != nil {
				return err
			}
			logSweptResource("policy: "+t, policy.Id, policy.Name)
		}
	}
	return nil
}

func sweepPolicyRulesByType(ruleType string, client api.OktaIDaaSClient) error {
	ctx := context.Background()
	policies, _, err := client.OktaSDKClientV2().Policy.ListPolicies(ctx, &query.Params{Type: ruleType})
	if err != nil {
		return fmt.Errorf("failed to list policies in order to properly destroy rules: %v", err)
	}
	for _, _policy := range policies {
		policy := _policy.(*sdk.Policy)
		rules, _, err := client.OktaSDKSupplementClient().ListPolicyRules(ctx, policy.Id)
		if err != nil {
			return err
		}
		// Tests have always used default policy, I don't really think that is necessarily a good idea but
		// leaving for now, that means we only delete the rules and not the policy, we can keep it around.
		for i := range rules {
			if strings.HasPrefix(rules[i].Name, acctest.ResourcePrefixForTest) {
				_, err = client.OktaSDKClientV2().Policy.DeletePolicyRule(ctx, policy.Id, rules[i].Id)
				if err != nil {
					return err
				}
				logSweptResource("policy rule type: "+ruleType, policy.Id+"/"+rules[i].Id, rules[i].Name)
			}
		}
	}
	return nil
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

func IDaaSClientForTest(t *testing.T) api.OktaIDaaSClient {
	p := provider.Provider()
	d := resourceDataForTest(t, p.Schema)
	cfg := config.NewConfig(d)
	_ = cfg.LoadAPIClient()
	return cfg.OktaIDaaSClient
}

func resourceDataForTest(t *testing.T, s map[string]*schema_sdk.Schema) *schema_sdk.ResourceData {
	configValues := configValuesForTest()
	emptyConfigMap := map[string]interface{}{}
	d := schema_sdk.TestResourceDataRaw(t, s, emptyConfigMap)

	if len(configValues) > 0 {
		for k, v := range configValues {
			// lintignore:R001
			_ = d.Set(k, v)
		}
	}

	return d
}

func configValuesForTest() map[string]interface{} {
	return map[string]interface{}{
		"access_token":   os.Getenv("OKTA_ACCESS_TOKEN"),
		"api_token":      os.Getenv("OKTA_API_TOKEN"),
		"org_name":       os.Getenv("OKTA_ORG_NAME"),
		"base_url":       os.Getenv("OKTA_BASE_URL"),
		"client_id":      os.Getenv("OKTA_API_CLIENT_ID"),
		"scopes":         strings.Split(os.Getenv("OKTA_API_SCOPES"), ","),
		"private_key":    os.Getenv("OKTA_API_PRIVATE_KEY"),
		"private_key_id": os.Getenv("OKTA_API_PRIVATE_KEY_ID"),
		"http_proxy":     os.Getenv("OKTA_HTTP_PROXY"),
		"log_level":      hclog.LevelFromString(os.Getenv("TF_LOG")),
	}
}

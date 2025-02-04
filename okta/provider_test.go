package okta

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"github.com/hashicorp/terraform-plugin-mux/tf6to5server"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/okta-sdk-golang/v4/okta"
	"github.com/okta/terraform-provider-okta/sdk"
	"gopkg.in/dnaeon/go-vcr.v3/cassette"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
)

const TestDomainName = "dne-okta.com"

var (
	testAccProvidersFactories       map[string]func() (*schema.Provider, error)
	testAccProtoV5ProviderFactories map[string]func() (tfprotov5.ProviderServer, error)
	testAccMergeProvidersFactories  map[string]func() (tfprotov5.ProviderServer, error)
	testSdkV3Client                 *okta.APIClient
	testSdkV2Client                 *sdk.Client
	testSdkSupplementClient         *sdk.APISupplement

	// NOTE: Our VCR set up is inspired by terraform-provider-google
	providerConfigsLock = sync.RWMutex{}
	providerConfigs     map[string]*Config
)

func init() {
	pluginProvider := Provider()

	// v2 provider - terraform-plugin-sdk
	testAccProvidersFactories = map[string]func() (*schema.Provider, error){
		"okta": func() (*schema.Provider, error) {
			return pluginProvider, nil
		},
	}

	// v3 provider - terraform-plugin-framework
	// TODO: Uses v5 protocol for now, however lets swap to v6 when a drop of support for TF versions prior to 1.0 can be made
	frameworkProvider := NewFrameworkProvider("dev")
	framework, err := tf6to5server.DowngradeServer(context.Background(), providerserver.NewProtocol6(frameworkProvider))

	testAccProtoV5ProviderFactories = map[string]func() (tfprotov5.ProviderServer, error){
		"okta": func() (tfprotov5.ProviderServer, error) {
			return framework, err
		},
	}
	providers := []func() tfprotov5.ProviderServer{
		// v2 plugin
		pluginProvider.GRPCProvider,
		// v3 plugin
		providerserver.NewProtocol5(frameworkProvider),
	}

	// mux'd provider (v2 + v3) - terraform-plugin-mux
	muxServer, err := tf5muxserver.NewMuxServer(context.Background(), providers...)
	if err != nil {
		log.Fatal(err)
	}
	testAccMergeProvidersFactories = map[string]func() (tfprotov5.ProviderServer, error){
		"okta": func() (tfprotov5.ProviderServer, error) {
			return muxServer.ProviderServer(), nil
		},
	}

	providerConfigs = make(map[string]*Config)
}

// TestMain overridden main testing function. Package level BeforeAll and AfterAll.
// It also delineates between acceptance tests and unit tests
func TestMain(m *testing.M) {
	// TF_VAR_hostname allows the real hostname to be scripted into the config tests
	// see examples/okta_resource_set/basic.tf
	os.Setenv("TF_VAR_hostname", fmt.Sprintf("%s.%s", os.Getenv("OKTA_ORG_NAME"), os.Getenv("OKTA_BASE_URL")))

	// NOTE: Acceptance test sweepers are necessary to prevent dangling
	// resources.
	// NOTE: Don't run sweepers if we are playing back VCR as nothing should be
	// going over the wire
	if os.Getenv("OKTA_VCR_TF_ACC") != "play" {
		setupSweeper(adminRoleCustom, sweepCustomRoles)
		setupSweeper("okta_*_app", sweepTestApps)
		setupSweeper(authServer, sweepAuthServers)
		setupSweeper(behavior, sweepBehaviors)
		setupSweeper(emailCustomization, sweepEmailCustomization)
		setupSweeper(groupRule, sweepGroupRules)
		setupSweeper("okta_*_idp", sweepTestIdps)
		setupSweeper(inlineHook, sweepInlineHooks)
		setupSweeper(group, sweepGroups)
		setupSweeper(groupSchemaProperty, sweepGroupCustomSchema)
		setupSweeper(linkDefinition, sweepLinkDefinitions)
		setupSweeper(logStream, sweepLogStreams)
		setupSweeper(networkZone, sweepNetworkZones)
		setupSweeper(policyMfa, sweepMfaPolicies)
		setupSweeper(policyPassword, sweepPasswordPolicies)
		setupSweeper(policyRuleIdpDiscovery, sweepPolicyRuleIdpDiscovery)
		setupSweeper(policyRuleMfa, sweepMfaPolicyRules)
		setupSweeper(policyRulePassword, sweepPolicyRulePasswords)
		setupSweeper(policyRuleSignOn, sweepSignOnPolicyRules)
		setupSweeper(policySignOn, sweepAccessPolicies)
		// setupSweeper(policySignOn, sweepSignOnPolicies)
		setupSweeper(resourceSet, sweepResourceSets)
		setupSweeper(user, sweepUsers)
		setupSweeper(userSchemaProperty, sweepUserCustomSchema)
		setupSweeper(userType, sweepUserTypes)
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

func testAccPreCheck(t *testing.T) func() {
	return func() {
		err := accPreCheck()
		if err != nil {
			t.Fatalf("%v", err)
		}
	}
}

func accPreCheck() error {
	if v := os.Getenv("OKTA_ORG_NAME"); v == "" {
		return errors.New("OKTA_ORG_NAME must be set for acceptance tests")
	}
	token := os.Getenv("OKTA_API_TOKEN")
	clientID := os.Getenv("OKTA_API_CLIENT_ID")
	privateKey := os.Getenv("OKTA_API_PRIVATE_KEY")
	privateKeyID := os.Getenv("OKTA_API_PRIVATE_KEY_ID")
	scopes := os.Getenv("OKTA_API_SCOPES")
	if token == "" && (clientID == "" || scopes == "" || privateKey == "" || privateKeyID == "") {
		return errors.New("either OKTA_API_TOKEN or (OKTA_API_CLIENT_ID and OKTA_API_SCOPES and OKTA_API_PRIVATE_KEY and OKTA_API_PRIVATE_KEY_ID) must be set for acceptance tests")
	}
	return nil
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

func sdkV3ClientForTest() *okta.APIClient {
	if testSdkV3Client != nil {
		return testSdkV3Client
	}
	return sdkV3Client
}

func sdkV2ClientForTest() *sdk.Client {
	if testSdkV2Client != nil {
		return testSdkV2Client
	}
	return sdkV2Client
}

func sdkSupplementClientForTest() *sdk.APISupplement {
	if testSdkSupplementClient != nil {
		return testSdkSupplementClient
	}
	return sdkSupplementClient
}

// oktaResourceTest is the entry to overriding the Terraform SDKs Acceptance
// Test framework before the call to resource.Test
func oktaResourceTest(t *testing.T, c resource.TestCase) {
	// plug in the VCR
	mgr := newVCRManager(t.Name())

	if !mgr.VCREnabled() {
		// live ACC / non-VCR test
		resource.Test(t, c)
		return
	}

	if !mgr.ValidMode() {
		t.Fatalf("ENV variable OKTA_VCR_TF_ACC value should be %q or %q but was %q", "play", "record", mgr.VCRModeName)
		return
	}

	// sdk client clients are set up in provider factories for test
	if mgr.IsPlaying() {
		if !mgr.HasCassettesToPlay() {
			t.Skipf("%q test is missing VCR cassette(s) at %q, skipping test. See .github/CONTRIBUTING.md#acceptance-tests-with-vcr for more information about playing/recording cassettes.", t.Name(), mgr.CassettesPath)
			return
		}

		cassettes := mgr.Cassettes()
		for _, cassette := range cassettes {
			// need to artificially set expected OKTA env vars if VCR is playing
			// VCR re-writes the name [cassette].dne-okta.com
			os.Setenv("OKTA_ORG_NAME", cassette)
			os.Setenv("OKTA_BASE_URL", TestDomainName)
			os.Setenv("OKTA_API_TOKEN", "token")
			mgr.SetCurrentCassette(cassette)
			// FIXME: need to get our VCR lined up correctly tf sdk v2 and tf plugin framework
			c.ProviderFactories = nil
			c.ProtoV5ProviderFactories = vcrProviderFactoriesForTest(mgr)
			c.CheckDestroy = nil
			fmt.Printf("=== VCR PLAY CASSETTE %q for %s\n", cassette, t.Name())

			// FIXME: Once we get fully mux'd ACC tests recording with VCR
			// revisit if we can call ParallelTest when playing.
			resource.Test(t, c)
		}
		return
	}

	if mgr.IsRecording() {
		defer closeRecorder(t, mgr)
		if mgr.AttemptedWriteIsMissingCassetteName() {
			t.Fatalf("%q test is attempting to write cassette to %q, but OKTA_VCR_CASSETTE ENV var is missing. See .github/CONTRIBUTING.md#acceptance-tests-with-vcr for more information.", t.Name(), mgr.CassettePath())
			return
		}
		if mgr.AttemptedWriteOfExistingCassette() {
			t.Skipf("%q test is attempting to write %s.yaml cassette, delete it first before attempting new write, skipping test. See .github/CONTRIBUTING.md#acceptance-tests-with-vcr for more information.", t.Name(), mgr.CassettePath())
			return
		}
		c.ProviderFactories = nil
		c.ProtoV5ProviderFactories = vcrProviderFactoriesForTest(mgr)
		fmt.Printf("=== VCR RECORD CASSETTE %q for %s\n", mgr.CurrentCassette, t.Name())
		resource.Test(t, c)
		return
	}
}

// vcrProviderFactoriesForTest Returns the overridden provider factories used by
// the resource test case given the state of the VCR manager.  func
// vcrProviderFactoriesForTest(mgr *vcrManager) map[string]func()
// (*schema.Provider, error) {
func vcrProviderFactoriesForTest(mgr *vcrManager) map[string]func() (tfprotov5.ProviderServer, error) {
	return map[string]func() (tfprotov5.ProviderServer, error){
		"okta": func() (tfprotov5.ProviderServer, error) {
			provider := Provider()

			// v2
			oldConfigureContextFunc := provider.ConfigureContextFunc
			provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
				config, _diag := vcrCachedConfigV2(ctx, d, oldConfigureContextFunc, mgr)
				if _diag.HasError() {
					return nil, _diag
				}

				// this is needed so teardown api calls are recorded by VCR and
				// we don't run ACC tests in parallel
				testSdkV3Client = config.oktaSDKClientV3
				testSdkV2Client = config.oktaSDKClientV2
				testSdkSupplementClient = config.oktaSDKsupplementClient

				return config, _diag
			}

			// v3
			// FIXME: this needs to be cleaned up so vcr and the sdk clients are all set up in one place
			// FIXME: need to get our VCR lined up correctly tf sdk v2 and tf plugin framework
			// FIXME: as-is, any tests using a framework provider can't record for VCR
			frameworkProvider := NewFrameworkProvider("dev")
			testFrameworkProvider := frameworkProvider.(*FrameworkProvider)

			// mux - v2+v3
			providers := []func() tfprotov5.ProviderServer{
				// v2 plugin
				provider.GRPCProvider,
				// v3 plugin
				providerserver.NewProtocol5(testFrameworkProvider),
			}

			muxServer, err := tf5muxserver.NewMuxServer(context.Background(), providers...)
			if err != nil {
				log.Fatal(err)
			}

			return muxServer, nil
		},
	}
}

func newVCRRecorder(mgr *vcrManager, transport http.RoundTripper) (rec *recorder.Recorder, err error) {
	rec, err = recorder.NewWithOptions(&recorder.Options{
		CassetteName:       mgr.CassettePath(),
		Mode:               mgr.VCRMode(),
		SkipRequestLatency: true, // skip how vcr will mimic the real request latency that it can record allowing for fast playback
		RealTransport:      transport,
	})
	if err != nil {
		return
	}

	// Defines how VCR will match requests to responses.
	rec.SetMatcher(func(r *http.Request, i cassette.Request) bool {
		// Default matcher compares method and URL only
		if !cassette.DefaultMatcher(r, i) {
			return false
		}
		// TODO: there might be header information we could inspect to make this more precise
		if r.Body == nil {
			return true
		}

		var b bytes.Buffer
		if _, err := b.ReadFrom(r.Body); err != nil {
			log.Printf("[DEBUG] Failed to read request body from cassette: %v", err)
			return false
		}
		r.Body = io.NopCloser(&b)
		reqBody := b.String()
		// If body matches identically, we are done
		if reqBody == i.Body {
			return true
		}

		// JSON might be the same, but reordered. Try parsing json and comparing
		contentType := r.Header.Get("Content-Type")
		if strings.Contains(contentType, "application/json") {
			var reqJson, cassetteJson interface{}
			if err := json.Unmarshal([]byte(reqBody), &reqJson); err != nil {
				log.Printf("[DEBUG] Failed to unmarshall request json: %v", err)
				return false
			}
			if err := json.Unmarshal([]byte(i.Body), &cassetteJson); err != nil {
				log.Printf("[DEBUG] Failed to unmarshall cassette json: %v", err)
				return false
			}
			return reflect.DeepEqual(reqJson, cassetteJson)
		}

		return true
	})

	rec.AddHook(func(i *cassette.Interaction) error {
		// need to scrub OKTA_ORG_NAME+OKTA_BASE_URL strings and rewrite as
		// [cassette-name].dne-okta.com so that HTTP requests that escape VCR
		// are bad.

		// test-admin.dne-okta.com
		vcrAdminHostname := fmt.Sprintf("%s-admin.%s", mgr.CurrentCassette, TestDomainName)
		// test.dne-okta.com
		vcrHostname := fmt.Sprintf("%s.%s", mgr.CurrentCassette, TestDomainName)
		// example-admin.okta.com
		orgAdminHostname := fmt.Sprintf("%s-admin.%s", os.Getenv("OKTA_ORG_NAME"), os.Getenv("OKTA_BASE_URL"))
		// example.okta.com
		orgHostname := fmt.Sprintf("%s.%s", os.Getenv("OKTA_ORG_NAME"), os.Getenv("OKTA_BASE_URL"))

		// test-admin
		vcrAdminOrgName := fmt.Sprintf("%s-admin", mgr.CurrentCassette)
		// test
		vcrOrgName := mgr.CurrentCassette
		// example-admin
		adminOrgName := fmt.Sprintf("%s-admin", os.Getenv("OKTA_ORG_NAME"))
		// example
		orgName := os.Getenv("OKTA_ORG_NAME")

		// re-write the Authorization header
		authHeader := "Authorization"
		if auth, ok := firstHeaderValue(authHeader, i.Request.Headers); ok {
			i.Request.Headers.Del(authHeader)
			parts := strings.Split(auth, " ")
			i.Request.Headers.Set("Authorization", fmt.Sprintf("%s REDACTED", parts[0]))
		}

		// save disk space, clean up what gets written to disk
		i.Request.Headers.Del("User-Agent")
		deleteResponseHeaders := []string{
			"Cache-Control",
			"Content-Security-Policy",
			"Content-Security-Policy-Report-Only",
			"duration",
			"Expect-Ct",
			"Expires",
			"P3p",
			"Pragma",
			"Public-Key-Pins-Report-Only",
			"Server",
			"Set-Cookie",
			"Strict-Transport-Security",
			"Vary",
		}
		for _, header := range deleteResponseHeaders {
			i.Response.Headers.Del(header)
		}
		for name := range i.Response.Headers {
			// delete all X-headers
			if strings.HasPrefix(name, "X-") {
				i.Response.Headers.Del(name)
				continue
			}
		}

		// scrub client assertions out of token requests
		m := regexp.MustCompile("client_assertion=[^&]+")
		i.Request.URL = m.ReplaceAllString(i.Request.URL, "client_assertion=abc123")

		// replace admin based hostname before regular variations
		// %s/example-admin.okta.com/test-admin.dne-okta.com/
		i.Request.Host = strings.ReplaceAll(i.Request.Host, orgAdminHostname, vcrAdminHostname)
		// %s/example.okta.com/test.dne-okta.com/
		i.Request.Host = strings.ReplaceAll(i.Request.Host, orgHostname, vcrHostname)

		// %s/example-admin.okta.com/test-admin.dne-okta.com/
		i.Request.URL = strings.ReplaceAll(i.Request.URL, orgAdminHostname, vcrAdminHostname)
		// %s/example.okta.com/test.dne-okta.com/
		i.Request.URL = strings.ReplaceAll(i.Request.URL, orgHostname, vcrHostname)

		// %s/example-admin.okta.com/test-admin.dne-okta.com/
		i.Request.Body = strings.ReplaceAll(i.Request.Body, orgAdminHostname, vcrAdminHostname)
		// %s/example.okta.com/test.dne-okta.com/
		i.Request.Body = strings.ReplaceAll(i.Request.Body, orgHostname, vcrHostname)

		// %s/example-admin/test-admin/
		i.Request.Body = strings.ReplaceAll(i.Request.Body, adminOrgName, vcrAdminOrgName)
		// %s/example/test/
		i.Request.Body = strings.ReplaceAll(i.Request.Body, orgName, vcrOrgName)

		// %s/example-admin/test-admin/
		i.Response.Body = strings.ReplaceAll(i.Response.Body, adminOrgName, vcrAdminOrgName)
		// %s/example/test/
		i.Response.Body = strings.ReplaceAll(i.Response.Body, orgName, vcrOrgName)

		// %s/example-admin.okta.com/test-admin.dne-okta.com/
		headerLinks := replaceHeaderValues(i.Response.Headers["Link"], orgAdminHostname, vcrAdminHostname)
		// %s/example.okta.com/test.dne-okta.com/
		headerLinks = replaceHeaderValues(headerLinks, orgHostname, vcrHostname)
		i.Response.Headers.Del("Link")
		for _, val := range headerLinks {
			i.Response.Headers.Add("Link", val)
		}

		return nil
	}, recorder.AfterCaptureHook)

	return
}

func providerConfig(mgr *vcrManager) (*Config, bool) {
	c, ok := providerConfigs[mgr.TestAndCassetteNameKey()]
	return c, ok
}

// vcrCachedConfig V2 We need to hijack the provider's ConfigureContextFunc as
// it is called many times during an operation which has the side effect of
// resetting config values such as the http client that okta-sdk-golang
// utilizes. Instead, we want to create one dedicated VCR http transport for
// each test. The dedicated transport holds the API calls made specific to that
// test.
func vcrCachedConfigV2(ctx context.Context, d *schema.ResourceData, configureFunc schema.ConfigureContextFunc, mgr *vcrManager) (*Config, diag.Diagnostics) {
	providerConfigsLock.RLock()
	v, ok := providerConfig(mgr)
	providerConfigsLock.RUnlock()
	if ok {
		// config is cached, proceed
		return v, nil
	}

	c, diags := configureFunc(ctx, d)

	if diags.HasError() {
		return nil, diags
	}

	config := c.(*Config)
	config.SetTimeOperations(NewTestTimeOperations())

	transport := config.oktaSDKClientV3.GetConfig().HTTPClient.Transport
	rec, err := newVCRRecorder(mgr, transport)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	// VCR takes over http transport duties
	rt := http.RoundTripper(rec)
	config.resetHttpTransport(&rt)

	providerConfigsLock.Lock()
	providerConfigs[mgr.TestAndCassetteNameKey()] = config
	providerConfigsLock.Unlock()
	return config, nil
}

func replaceHeaderValues(currentValues []string, oldValue, newValue string) []string {
	result := []string{}
	for _, text := range currentValues {
		result = append(result, strings.ReplaceAll(text, oldValue, newValue))
	}
	return result
}

func firstHeaderValue(name string, headers http.Header) (string, bool) {
	if vals, ok := headers[name]; ok && len(vals) > 0 {
		return vals[0], true
	}
	return "", false
}

// closeRecorder closes the VCR recorder to save the cassette file
func closeRecorder(t *testing.T, vcr *vcrManager) {
	providerConfigsLock.RLock()
	config, ok := providerConfigs[vcr.TestAndCassetteNameKey()]
	providerConfigsLock.RUnlock()
	if ok {
		// don't record failing test runs
		if !t.Failed() {
			// If a test succeeds, write new seed/yaml to files
			err := config.oktaSDKClientV2.GetConfig().HttpClient.Transport.(*recorder.Recorder).Stop()
			if err != nil {
				t.Error(err)
			}
		}
		// Clean up test config
		providerConfigsLock.Lock()
		delete(providerConfigs, t.Name())
		providerConfigsLock.Unlock()
	}
}

// newVCRManager Returns a vcr manager
func newVCRManager(testName string) *vcrManager {
	dir, _ := os.Getwd()
	vcrFixturesHome := path.Join(dir, "../test/fixtures/vcr")
	cassettesPath := path.Join(vcrFixturesHome, testName)
	return &vcrManager{
		Name:            testName,
		FixturesHome:    vcrFixturesHome,
		CassettesPath:   cassettesPath,
		CurrentCassette: os.Getenv("OKTA_VCR_CASSETTE"),
		VCRModeName:     os.Getenv("OKTA_VCR_TF_ACC"),
	}
}

func (m *vcrManager) SetCurrentCassette(name string) {
	m.CurrentCassette = name
}

// TestAndCassetteNameKey is intended to be used as the key for configs that
// utilized for each test cassette.
func (m *vcrManager) TestAndCassetteNameKey() string {
	return fmt.Sprintf("%s-%s", m.Name, m.CurrentCassette)
}

// VcrEnabled test is considered to be in VCR mode if ENV var OKTA_VCR_TF_ACC
// is not empty. Valid values are "record" for recording VCR cassettes (test
// runs), and "play" for playing all cassettes of a test.
func (m *vcrManager) VCREnabled() bool {
	return m.VCRModeName != ""
}

// ValidMode Is this a valid vcr mode given vcr is enabled?
func (m *vcrManager) ValidMode() bool {
	return m.VCREnabled() && (m.VCRModeName == "play" || m.VCRModeName == "record")
}

// HasCassettesToPlay VCR is in play mode and there are cassette files to play.
func (m *vcrManager) HasCassettesToPlay() bool {
	cassetteDir := filepath.Dir(m.CassettesPath)
	files, err := os.ReadDir(cassetteDir)
	if err != nil {
		return false
	}
	return m.VCRModeName == "play" && len(files) > 0
}

// CassettePath the path to what would be the current cassette.
func (m *vcrManager) CassettePath() string {
	return path.Join(m.CassettesPath, m.CurrentCassette)
}

// AttemptedWriteOfExistingCassette given a vcr mode of record the current
// cassette file exists.
func (m *vcrManager) AttemptedWriteOfExistingCassette() bool {
	_, err := os.Stat(fmt.Sprintf("%s.yaml", m.CassettePath()))
	return m.VCRModeName == "record" && err == nil
}

// AttemptedWriteIsMissingCassetteName Tests if cassette name is present for recording
func (m *vcrManager) AttemptedWriteIsMissingCassetteName() bool {
	return m.VCRModeName == "record" && os.Getenv("OKTA_VCR_CASSETTE") == ""
}

// VCRMode the recorder.Mode value based on OKTA_VCR_TF_ACC
func (m *vcrManager) VCRMode() recorder.Mode {
	if m.VCRModeName == "record" {
		return recorder.ModeRecordOnly
	}

	return recorder.ModeReplayOnly
}

// IsPlaying VCR mode is play
func (m *vcrManager) IsPlaying() bool {
	return m.VCRModeName == "play"
}

// IsRecording VCR mode is record
func (m *vcrManager) IsRecording() bool {
	return m.VCRModeName == "record"
}

// Cassettes Slice of all cassettes for the given test.
func (m *vcrManager) Cassettes() []string {
	if os.Getenv("OKTA_VCR_CASSETTE") != "" {
		return []string{os.Getenv("OKTA_VCR_CASSETTE")}
	}

	cassettes := []string{}
	files, err := os.ReadDir(path.Join(m.CassettesPath))
	if err != nil {
		return cassettes
	}
	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}
		parts := strings.Split(file.Name(), ".")
		cassettes = append(cassettes, parts[0])
	}
	return cassettes
}

type vcrManager struct {
	Name            string
	FixturesHome    string
	CassettesPath   string
	CurrentCassette string
	VCRModeName     string
}

func skipVCRTest(t *testing.T) bool {
	skip := os.Getenv("OKTA_VCR_TF_ACC") != ""
	if skip {
		t.Skipf("test %q is not VCR compatible", t.Name())
	}
	return skip
}

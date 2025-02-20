package acctest

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
	"strconv"
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
	"github.com/okta/terraform-provider-okta/okta/config"
	"github.com/okta/terraform-provider-okta/okta/fwprovider"
	"github.com/okta/terraform-provider-okta/okta/provider"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"
)

const TestDomainName = "dne-okta.com"

var ResourceNamePrefixForTest = "testAcc"

func init() {
	providerConfigs = make(map[string]*config.Config)

	pluginProvider := provider.Provider()

	// v2 provider - terraform-plugin-sdk
	ProvidersFactoriesForTestAcc = map[string]func() (*schema.Provider, error){
		"okta": func() (*schema.Provider, error) {
			return pluginProvider, nil
		},
	}

	// v3 provider - terraform-plugin-framework
	// TODO: Uses v5 protocol for now, however lets swap to v6 when a drop of support for TF versions prior to 1.0 can be made
	frameworkProvider := fwprovider.NewFrameworkProvider("dev")
	framework, err := tf6to5server.DowngradeServer(context.Background(), providerserver.NewProtocol6(frameworkProvider))

	ProtoV5ProviderFactoriesForTestAcc = map[string]func() (tfprotov5.ProviderServer, error){
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
	MergeProvidersFactoriesForTestAcc = map[string]func() (tfprotov5.ProviderServer, error){
		"okta": func() (tfprotov5.ProviderServer, error) {
			return muxServer.ProviderServer(), nil
		},
	}
}

var (
	// NOTE: Our VCR set up is inspired by terraform-provider-google
	providerConfigsLock                = sync.RWMutex{}
	providerConfigs                    map[string]*config.Config
	ProvidersFactoriesForTestAcc       map[string]func() (*schema.Provider, error)
	ProtoV5ProviderFactoriesForTestAcc map[string]func() (tfprotov5.ProviderServer, error)
	MergeProvidersFactoriesForTestAcc  map[string]func() (tfprotov5.ProviderServer, error)
)

// OktaResourceTest is the entry to overriding the Terraform SDKs Acceptance
// Test framework before the call to resource.Test
func OktaResourceTest(t *testing.T, c resource.TestCase) {
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

type vcrManager struct {
	Name            string
	FixturesHome    string
	CassettesPath   string
	CurrentCassette string
	VCRModeName     string
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

func providerConfig(mgr *vcrManager) (*config.Config, bool) {
	c, ok := providerConfigs[mgr.TestAndCassetteNameKey()]
	return c, ok
}

// vcrCachedConfig V2 We need to hijack the provider's ConfigureContextFunc as
// it is called many times during an operation which has the side effect of
// resetting config values such as the http client that okta-sdk-golang
// utilizes. Instead, we want to create one dedicated VCR http transport for
// each test. The dedicated transport holds the API calls made specific to that
// test.
func vcrCachedConfigV2(ctx context.Context, d *schema.ResourceData, configureFunc schema.ConfigureContextFunc, mgr *vcrManager) (*config.Config, diag.Diagnostics) {
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

	cfg := c.(*config.Config)
	cfg.SetTimeOperations(config.NewTestTimeOperations())

	transport := cfg.OktaSDKClientV3.GetConfig().HTTPClient.Transport
	rec, err := newVCRRecorder(mgr, transport)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	// VCR takes over http transport duties
	rt := http.RoundTripper(rec)
	cfg.ResetHttpTransport(&rt)

	providerConfigsLock.Lock()
	providerConfigs[mgr.TestAndCassetteNameKey()] = cfg
	providerConfigsLock.Unlock()
	return cfg, nil
}

// vcrProviderFactoriesForTest Returns the overridden provider factories used by
// the resource test case given the state of the VCR manager.  func
// vcrProviderFactoriesForTest(mgr *vcrManager) map[string]func()
// (*schema.Provider, error) {
func vcrProviderFactoriesForTest(mgr *vcrManager) map[string]func() (tfprotov5.ProviderServer, error) {
	return map[string]func() (tfprotov5.ProviderServer, error){
		"okta": func() (tfprotov5.ProviderServer, error) {
			provider := provider.Provider()

			// v2
			oldConfigureContextFunc := provider.ConfigureContextFunc
			provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
				config, _diag := vcrCachedConfigV2(ctx, d, oldConfigureContextFunc, mgr)
				if _diag.HasError() {
					return nil, _diag
				}

				// this is needed so teardown api calls are recorded by VCR and
				// we don't run ACC tests in parallel
				// TODO FIXME
				/*
					testSdkV5Client = config.OktaSDKClientV5
					testSdkV3Client = config.OktaSDKClientV3
					testSdkV2Client = config.OktaSDKClientV2
					testSdkSupplementClient = config.OktaSDKsupplementClient
				*/

				return config, _diag
			}

			// v3
			// FIXME: this needs to be cleaned up so vcr and the sdk clients are all set up in one place
			// FIXME: need to get our VCR lined up correctly tf sdk v2 and tf plugin framework
			// FIXME: as-is, any tests using a framework provider can't record for VCR
			//frameworkProvider := fwprovider.NewFrameworkProvider("dev")
			//testFrameworkProvider := frameworkProvider.(*FrameworkProvider)
			testFrameworkProvider := fwprovider.NewFrameworkProvider("dev")

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
	// Defines how VCR will match requests to responses.

	vcrOpts := []recorder.Option{
		recorder.WithMatcher(vcrMatcher(mgr)),
		recorder.WithHook(vcrHook(mgr), recorder.AfterCaptureHook),
		recorder.WithMode(mgr.VCRMode()),
		recorder.WithSkipRequestLatency(true),
		recorder.WithRealTransport(transport),
	}
	rec, err = recorder.New(mgr.CassettePath(), vcrOpts...)

	return
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
			err := config.OktaSDKClientV2.GetConfig().HttpClient.Transport.(*recorder.Recorder).Stop()
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

func firstHeaderValue(name string, headers http.Header) (string, bool) {
	if vals, ok := headers[name]; ok && len(vals) > 0 {
		return vals[0], true
	}
	return "", false
}

func replaceHeaderValues(currentValues []string, oldValue, newValue string) []string {
	result := []string{}
	for _, text := range currentValues {
		result = append(result, strings.ReplaceAll(text, oldValue, newValue))
	}
	return result
}

func vcrMatcher(mgr *vcrManager) func(*http.Request, cassette.Request) bool {
	return func(r *http.Request, i cassette.Request) bool {
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
	}
}

func vcrHook(mgr *vcrManager) func(*cassette.Interaction) error {
	return func(i *cassette.Interaction) error {
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
	}
}

func AccPreCheck(t *testing.T) func() {
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

func AccProvidersFactoriesForTest() map[string]func() (*schema.Provider, error) {
	return ProvidersFactoriesForTestAcc
}

func AccMergeProvidersFactoriesForTest() map[string]func() (tfprotov5.ProviderServer, error) {
	return MergeProvidersFactoriesForTestAcc
}

func BuildResourceName(testID int) string {
	return ResourceNamePrefixForTest + "_" + strconv.Itoa(testID)
}

func BuildResourceNameWithPrefix(prefix string, testID int) string {
	return prefix + "_" + strconv.Itoa(testID)
}

func BuildResourceFQN(resourceType string, testID int) string {
	return resourceType + "." + BuildResourceName(testID)
}

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

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	schema_sdk "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-governance-sdk-golang/governance"
	v4SdkOkta "github.com/okta/okta-sdk-golang/v4/okta"
	v5SdkOkta "github.com/okta/okta-sdk-golang/v5/okta"
	v6SdkOkta "github.com/okta/okta-sdk-golang/v6/okta"
	"github.com/okta/terraform-provider-okta/okta/api"
	"github.com/okta/terraform-provider-okta/okta/config"
	"github.com/okta/terraform-provider-okta/okta/fwprovider"
	okta_provider "github.com/okta/terraform-provider-okta/okta/provider"
	oktaSdk "github.com/okta/terraform-provider-okta/sdk"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"
)

const (
	ResourcePrefixForTest = "testAcc"
	TestDomainName        = "dne-okta.com"
)

var (
	// ResourceNamePrefixForTest common resource name/lable used for sweeper
	// clean up of resources
	ResourceNamePrefixForTest = "testAcc"

	// NOTE: Our VCR set up is inspired by terraform-provider-google

	providerConfigsLock = sync.RWMutex{}
	providerConfigs     map[string]*config.Config
	vcrMgrsLock         = sync.RWMutex{}
	vcrMgrs             map[string]*vcrManager
)

func init() {
	providerConfigs = make(map[string]*config.Config)
	vcrMgrs = make(map[string]*vcrManager)
}

// OktaResourceTest is the entry to overriding the Terraform SDKs Acceptance
// Test framework before the call to resource.Test
func OktaResourceTest(t *testing.T, c resource.TestCase) {
	t.Helper()

	// plug in the VCR
	mgr := currentVCRManager(t.Name())

	if !mgr.IsVcrEnabled() {
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
		if !mgr.HasIdaasCassettesToPlay() {
			t.Skipf("%q test is missing VCR cassette(s) at %q, skipping test. See .github/CONTRIBUTING.md#acceptance-tests-with-vcr for more information about playing/recording cassettes.", t.Name(), mgr.IdaaSCassettesPath)
			return
		}

		if !mgr.HasGovernanceCassettesToPlay() {
			t.Skipf("%q test is missing VCR cassette(s) at %q, skipping test. See .github/CONTRIBUTING.md#acceptance-tests-with-vcr for more information about playing/recording cassettes.", t.Name(), mgr.GovernanceCassettesPath)
			return
		}

		// FIXME most of the skips we get are from the "classic-00" cassettes,
		// not the "oie-00" so reverse order the cassettes
		cassettes := mgr.Cassettes()
		// for _, cassette := range cassettes {
		for i := len(cassettes) - 1; i >= 0; i-- {
			cassette := cassettes[i]
			// need to artificially set expected OKTA env vars if VCR is playing
			// VCR re-writes the name [cassette].dne-okta.com
			os.Setenv("OKTA_ORG_NAME", cassette)
			os.Setenv("OKTA_BASE_URL", TestDomainName)
			os.Setenv("OKTA_API_TOKEN", "token")
			os.Setenv("TF_VAR_hostname", fmt.Sprintf("%s.%s", cassette, TestDomainName))
			mgr.SetCurrentCassette(cassette)

			// we disable check destroy when recording/playing vcr tests
			c.CheckDestroy = nil
			fmt.Printf("=== VCR PLAY CASSETTE %q for %s\n", cassette, t.Name())

			// FIXME: Once we get fully mux'd ACC tests recording with VCR
			// revisit if we can call ParallelTest when playing.

			// FIXME: if we get a skip from one cassette the next will not run.
			// Capturing the resource.Test result in a go routine does not give
			// a means of capturing and rerunning the test.
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
		// we disable check destroy when recording/playing vcr tests
		c.CheckDestroy = nil
		resource.Test(t, c)
		return
	}
}

type vcrManager struct {
	Name                    string
	IdaasFixturesHome       string
	IdaaSCassettesPath      string
	GovernanceFixturesHome  string
	GovernanceCassettesPath string
	CurrentCassette         string
	VCRModeName             string
}

// newVCRManager Returns a vcr manager
func newVCRManager(testName string) *vcrManager {
	dir, _ := os.Getwd()
	idaasVcrFixturesHome := path.Join(dir, "../../../test/fixtures/vcr/idaas")
	idaasCassettesPath := path.Join(idaasVcrFixturesHome, testName)
	governanceVcrFixturesHome := path.Join(dir, "../../../test/fixtures/vcr/governance")
	governanceCassettesPath := path.Join(governanceVcrFixturesHome, testName)

	return &vcrManager{
		Name:                    testName,
		IdaasFixturesHome:       idaasVcrFixturesHome,
		IdaaSCassettesPath:      idaasCassettesPath,
		GovernanceFixturesHome:  governanceVcrFixturesHome,
		GovernanceCassettesPath: governanceCassettesPath,
		CurrentCassette:         os.Getenv("OKTA_VCR_CASSETTE"),
		VCRModeName:             os.Getenv("OKTA_VCR_TF_ACC"),
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
func vcrCachedConfigV2(ctx context.Context, d *schema_sdk.ResourceData, configureFunc schema_sdk.ConfigureContextFunc, mgr *vcrManager) (*config.Config, diag.Diagnostics) {
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

	idaasTestClient := NewVcrIDaaSClient(d)
	cfg.SetIdaasAPIClient(idaasTestClient)

	idaaSRec, err := newVCRRecorder(mgr, idaasTestClient.Transport())
	if err != nil {
		return nil, diag.FromErr(err)
	}
	idaaSRt := http.RoundTripper(idaaSRec)
	idaasTestClient.SetTransport(idaaSRt)

	governanceTestClient := NewVcrGovernanceClient(d)
	cfg.SetGovernanceAPIClient(governanceTestClient)

	governanceRec, err := newVCRRecorder(mgr, governanceTestClient.Transport())
	if err != nil {
		return nil, diag.FromErr(err)
	}
	governanceRt := http.RoundTripper(governanceRec)
	governanceTestClient.SetTransport(governanceRt)

	providerConfigsLock.Lock()
	providerConfigs[mgr.TestAndCassetteNameKey()] = cfg
	providerConfigsLock.Unlock()
	return cfg, nil
}

func newVCRRecorder(mgr *vcrManager, transport http.RoundTripper) (rec *recorder.Recorder, err error) {
	// Defines how VCR will match requests to responses.

	vcrOpts := []recorder.Option{
		recorder.WithMatcher(vcrMatcher()),
		recorder.WithHook(vcrHook(mgr), recorder.AfterCaptureHook),
		recorder.WithHook(vcrSaveHook(mgr), recorder.BeforeSaveHook),
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
		// don't record failing test runs (unless explicitly allowed for error response testing)
		if !t.Failed() || os.Getenv("OKTA_VCR_RECORD_FAILURES") == "true" {
			// If a test succeeds, write new seed/yaml to files
			if strings.Contains(strings.ToLower(vcr.CassettePath()), "idaas") {
				rtIDaasHelper := config.OktaIDaaSClient.(HttpClientHelper)
				rt := rtIDaasHelper.Transport()
				err := rt.(*recorder.Recorder).Stop()
				if err != nil {
					t.Error(err)
				}
			} else {
				rtGovernanceHelper := config.OktaGovernanceClient.(HttpClientHelper)
				rtGovernance := rtGovernanceHelper.Transport()
				err := rtGovernance.(*recorder.Recorder).Stop()
				if err != nil {
					t.Error(err)
				}
			}
			fmt.Printf("=== VCR WROTE CASSETTE %s.yaml\n", vcr.CassettePath())
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

// IsVcrEnabled test is considered to be in VCR mode if ENV var OKTA_VCR_TF_ACC
// is not empty. Valid values are "record" for recording VCR cassettes (test
// runs), and "play" for playing all cassettes of a test.
func (m *vcrManager) IsVcrEnabled() bool {
	return m.IsPlaying() || m.IsRecording()
}

// ValidMode Is this a valid vcr mode given vcr is enabled?
func (m *vcrManager) ValidMode() bool {
	return m.VCRModeName == "play" || m.VCRModeName == "record"
}

// HasIdaasCassettesToPlay VCR is in play mode and there are cassette files to play.
func (m *vcrManager) HasIdaasCassettesToPlay() bool {
	cassetteDir := filepath.Dir(m.IdaaSCassettesPath)
	files, err := os.ReadDir(cassetteDir)
	if err != nil {
		return false
	}
	return m.VCRModeName == "play" && len(files) > 0
}

// HasGovernanceCassettesToPlay VCR is in play mode and there are cassette files to play.
func (m *vcrManager) HasGovernanceCassettesToPlay() bool {
	cassetteDir := filepath.Dir(m.GovernanceCassettesPath)
	files, err := os.ReadDir(cassetteDir)
	if err != nil {
		return false
	}
	return m.VCRModeName == "play" && len(files) > 0
}

// CassettePath the path to what would be the current cassette.
func (m *vcrManager) CassettePath() string {
	wd, _ := os.Getwd()
	if strings.Contains(strings.ToLower(wd), "governance") {
		return path.Join(m.GovernanceCassettesPath, m.CurrentCassette)
	}
	return path.Join(m.IdaaSCassettesPath, m.CurrentCassette)
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
	idaasFiles, err := os.ReadDir(path.Join(m.IdaaSCassettesPath))
	if err != nil {
		return cassettes
	}
	cassettes = append(cassettes, m.getCassettes(idaasFiles, cassettes)...)
	governanceFiles, err := os.ReadDir(path.Join(m.GovernanceCassettesPath))
	if err != nil {
		return cassettes
	}
	cassettes = append(cassettes, m.getCassettes(governanceFiles, cassettes)...)
	return cassettes
}

func (m *vcrManager) getCassettes(files []os.DirEntry, cassettes []string) []string {
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

// compareJSONWithUnorderedArrays compares two JSON objects, treating arrays as unordered
// when they contain primitive values (strings, numbers, booleans)
func compareJSONWithUnorderedArrays(reqJson, cassetteJson interface{}) bool {
	switch req := reqJson.(type) {
	case map[string]interface{}:
		cassette, ok := cassetteJson.(map[string]interface{})
		if !ok {
			return false
		}

		// Check if all keys in request exist in cassette
		for key, reqValue := range req {
			cassetteValue, exists := cassette[key]
			if !exists {
				return false
			}
			if !compareJSONWithUnorderedArrays(reqValue, cassetteValue) {
				return false
			}
		}

		// Check if all keys in cassette exist in request
		for key := range cassette {
			if _, exists := req[key]; !exists {
				return false
			}
		}
		return true

	case []interface{}:
		cassette, ok := cassetteJson.([]interface{})
		if !ok {
			return false
		}

		if len(req) != len(cassette) {
			return false
		}

		// If all elements are primitive types, treat as unordered
		if isPrimitiveArray(req) && isPrimitiveArray(cassette) {
			return compareUnorderedArrays(req, cassette)
		}

		// Otherwise, compare in order
		for i, reqValue := range req {
			if !compareJSONWithUnorderedArrays(reqValue, cassette[i]) {
				return false
			}
		}
		return true

	default:
		return reflect.DeepEqual(reqJson, cassetteJson)
	}
}

// isPrimitiveArray checks if all elements in the array are primitive types
func isPrimitiveArray(arr []interface{}) bool {
	for _, item := range arr {
		switch item.(type) {
		case string, float64, int, int64, bool, nil:
			continue
		default:
			return false
		}
	}
	return true
}

// compareUnorderedArrays compares two arrays as unordered sets
func compareUnorderedArrays(req, cassette []interface{}) bool {
	if len(req) != len(cassette) {
		return false
	}

	// Convert to sets for comparison
	reqSet := make(map[string]bool)
	cassetteSet := make(map[string]bool)

	for _, item := range req {
		reqSet[fmt.Sprintf("%v", item)] = true
	}

	for _, item := range cassette {
		cassetteSet[fmt.Sprintf("%v", item)] = true
	}

	return reflect.DeepEqual(reqSet, cassetteSet)
}

func vcrMatcher() func(*http.Request, cassette.Request) bool {
	return func(r *http.Request, i cassette.Request) bool {
		if r.Method != i.Method {
			log.Printf("[DEBUG] VCR matcher: method mismatch - request: %s, cassette: %s", r.Method, i.Method)
			return false
		}
		if r.URL.String() != i.URL {
			log.Printf("[DEBUG] VCR matcher: URL mismatch - request: %s, cassette: %s", r.URL.String(), i.URL)
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
			if reflect.DeepEqual(reqJson, cassetteJson) {
				return true
			}

			// Try comparing with unordered array handling
			if compareJSONWithUnorderedArrays(reqJson, cassetteJson) {
				return true
			}

			log.Printf("[DEBUG] VCR matcher: JSON body mismatch - request: %s, cassette: %s", reqBody, i.Body)
			return false
		}

		log.Printf("[DEBUG] VCR matcher: body mismatch - request: %s, cassette: %s", reqBody, i.Body)
		return false
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

		// Note: Response body and header sanitization is now handled in vcrSaveHook
		// to avoid causing Terraform drift detection during recording.

		return nil
	}
}

// vcrSaveHook sanitizes response bodies before saving the cassette to disk.
// This ensures that sensitive domain information is replaced with sanitized values
// while preserving the original response during the test execution.
func vcrSaveHook(mgr *vcrManager) func(*cassette.Interaction) error {
	return func(i *cassette.Interaction) error {
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

		// Sanitize response body - replace real domain names with sanitized ones
		// %s/example-admin.okta.com/test-admin.dne-okta.com/
		i.Response.Body = strings.ReplaceAll(i.Response.Body, orgAdminHostname, vcrAdminHostname)
		// %s/example.okta.com/test.dne-okta.com/
		i.Response.Body = strings.ReplaceAll(i.Response.Body, orgHostname, vcrHostname)

		// %s/example-admin/test-admin/
		i.Response.Body = strings.ReplaceAll(i.Response.Body, adminOrgName, vcrAdminOrgName)
		// %s/example/test/
		i.Response.Body = strings.ReplaceAll(i.Response.Body, orgName, vcrOrgName)

		// Sanitize response headers that might contain domain information
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

func BuildResourceName(testID int) string {
	return ResourceNamePrefixForTest + "_" + strconv.Itoa(testID)
}

func BuildResourceNameWithPrefix(prefix string, testID int) string {
	return prefix + "_" + strconv.Itoa(testID)
}

func BuildResourceFQN(resourceType string, testID int) string {
	return resourceType + "." + BuildResourceName(testID)
}

// ProtoV5ProviderFactoriesForTestAcc is for ProtoV5ProviderFactories argument in acc tests
func ProtoV5ProviderFactoriesForTestAcc(t *testing.T) map[string]func() (tfprotov5.ProviderServer, error) {
	return map[string]func() (tfprotov5.ProviderServer, error){
		"okta": func() (tfprotov5.ProviderServer, error) {
			provider, err := ProvidersForTest(t.Name())
			return provider(), err
		},
	}
}

func ProvidersForTest(testName string) (func() tfprotov5.ProviderServer, error) {
	ctx := context.Background()

	pluginSDKProvider := GetPluginSDKProvider(testName)
	providers := []func() tfprotov5.ProviderServer{
		// terraform-plugin-sdk/v2
		pluginSDKProvider.GRPCProvider,
		// terraform-plugin-framework/provider
		providerserver.NewProtocol5(NewFrameworkTestProvider(testName, pluginSDKProvider)),
	}
	muxServer, err := tf5muxserver.NewMuxServer(ctx, providers...)
	if err != nil {
		return nil, err
	}

	return muxServer.ProviderServer, nil
}

type frameworkTestProvider struct {
	fwprovider.FrameworkProvider
	TestName string
}

func NewFrameworkTestProvider(testName string, pluginSDKProvider *schema_sdk.Provider) *frameworkTestProvider {
	return &frameworkTestProvider{
		FrameworkProvider: fwprovider.FrameworkProvider{
			PluginSDKProvider: pluginSDKProvider,
			Version:           "test",
		},
		TestName: testName,
	}
}

func (p *frameworkTestProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	p.FrameworkProvider.Configure(ctx, req, resp)
}

func GetPluginSDKProvider(testName string) *schema_sdk.Provider {
	vcrMgr := currentVCRManager(testName)
	oktaProvider := okta_provider.Provider()
	if vcrMgr.IsVcrEnabled() {
		oldConfigureContextFunc := oktaProvider.ConfigureContextFunc
		oktaProvider.ConfigureContextFunc = func(ctx context.Context, d *schema_sdk.ResourceData) (interface{}, diag.Diagnostics) {
			return vcrCachedConfigV2(ctx, d, oldConfigureContextFunc, vcrMgr)
		}
	}
	return oktaProvider
}

func SkipVCRTest(t *testing.T) bool {
	skip := os.Getenv("OKTA_VCR_TF_ACC") != ""
	if skip {
		t.Skipf("test %q is not VCR compatible", t.Name())
	}
	return skip
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ api.OktaIDaaSClient      = &vcrIDaaSTestClient{}
	_ HttpClientHelper         = &vcrIDaaSTestClient{}
	_ api.OktaGovernanceClient = &vcrGovernanceTestClient{}
	_ HttpClientHelper         = &vcrGovernanceTestClient{}
)

type HttpClientHelper interface {
	Transport() http.RoundTripper
	SetTransport(http.RoundTripper)
}

type vcrIDaaSTestClient struct {
	sdkV6Client         *v6SdkOkta.APIClient
	sdkV5Client         *v5SdkOkta.APIClient
	sdkV3Client         *v4SdkOkta.APIClient
	sdkV2Client         *oktaSdk.Client
	sdkSupplementClient *oktaSdk.APISupplement
	transport           http.RoundTripper
}

type vcrGovernanceTestClient struct {
	oktaGovernanceSDKClient *governance.OktaGovernanceAPIClient
	transport               http.RoundTripper
}

func NewVcrIDaaSClient(d *schema_sdk.ResourceData) *vcrIDaaSTestClient {
	c := config.NewConfig(d)
	c.Backoff = false
	c.LoadAPIClient()

	// force all the API clients on a new config to use the same round tripper
	// for VCR recording/playback
	tripper := c.OktaIDaaSClient.OktaSDKClientV5().GetConfig().HTTPClient.Transport
	c.OktaIDaaSClient.OktaSDKClientV6().GetConfig().HTTPClient.Transport = tripper
	c.OktaIDaaSClient.OktaSDKClientV3().GetConfig().HTTPClient.Transport = tripper
	c.OktaIDaaSClient.OktaSDKClientV2().GetConfig().HttpClient.Transport = tripper
	re := c.OktaIDaaSClient.OktaSDKClientV2().CloneRequestExecutor()
	re.SetHTTPTransport(tripper)
	supClient := &oktaSdk.APISupplement{
		RequestExecutor: re,
	}

	client := &vcrIDaaSTestClient{
		sdkV6Client:         c.OktaIDaaSClient.OktaSDKClientV6(),
		sdkV5Client:         c.OktaIDaaSClient.OktaSDKClientV5(),
		sdkV3Client:         c.OktaIDaaSClient.OktaSDKClientV3(),
		sdkV2Client:         c.OktaIDaaSClient.OktaSDKClientV2(),
		sdkSupplementClient: supClient,
		transport:           tripper,
	}

	return client
}

func NewVcrGovernanceClient(d *schema_sdk.ResourceData) *vcrGovernanceTestClient {
	c := config.NewConfig(d)
	c.Backoff = false
	err := c.LoadAPIClient()
	if err != nil {
		fmt.Printf("Error loading API client: %v", err)
		return nil
	}

	// force all the API clients on a new config to use the same round tripper
	// for VCR recording/playback
	tripper := c.OktaGovernanceClient.OktaGovernanceSDKClient().GetConfig().HTTPClient.Transport

	client := &vcrGovernanceTestClient{
		oktaGovernanceSDKClient: c.OktaGovernanceClient.OktaGovernanceSDKClient(),
		transport:               tripper,
	}

	return client
}

func (c *vcrIDaaSTestClient) Transport() http.RoundTripper {
	return c.transport
}

func (c *vcrGovernanceTestClient) Transport() http.RoundTripper {
	return c.transport
}

func (c *vcrIDaaSTestClient) SetTransport(rt http.RoundTripper) {
	c.transport = rt

	c.sdkV6Client.GetConfig().HTTPClient.Transport = rt
	c.sdkV5Client.GetConfig().HTTPClient.Transport = rt
	c.sdkV3Client.GetConfig().HTTPClient.Transport = rt
	c.sdkV2Client.GetConfig().HttpClient.Transport = rt
	re := c.sdkV2Client.CloneRequestExecutor()
	re.SetHTTPTransport(rt)
	c.sdkSupplementClient = &oktaSdk.APISupplement{
		RequestExecutor: re,
	}
}

func (c *vcrGovernanceTestClient) SetTransport(rt http.RoundTripper) {
	c.transport = rt
	c.oktaGovernanceSDKClient.GetConfig().HTTPClient.Transport = rt
}

func (c *vcrIDaaSTestClient) OktaSDKClientV6() *v6SdkOkta.APIClient {
	return c.sdkV6Client
}

func (c *vcrIDaaSTestClient) OktaSDKClientV5() *v5SdkOkta.APIClient {
	return c.sdkV5Client
}

func (c *vcrIDaaSTestClient) OktaSDKClientV3() *v4SdkOkta.APIClient {
	return c.sdkV3Client
}

func (c *vcrIDaaSTestClient) OktaSDKClientV2() *oktaSdk.Client {
	return c.sdkV2Client
}

func (c *vcrIDaaSTestClient) OktaSDKSupplementClient() *oktaSdk.APISupplement {
	return c.sdkSupplementClient
}

func (c *vcrGovernanceTestClient) OktaGovernanceSDKClient() *governance.OktaGovernanceAPIClient {
	return c.oktaGovernanceSDKClient
}

func currentVCRManager(name string) *vcrManager {
	mgr, ok := vcrMgrs[name]
	if ok {
		return mgr
	}

	vcrMgrsLock.RLock()
	mgr = newVCRManager(name)
	vcrMgrs[name] = mgr
	vcrMgrsLock.RUnlock()

	return mgr
}

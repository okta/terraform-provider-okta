package config

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/api"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

type (
	// Config contains our provider schema values and Okta clients
	Config struct {
		AccessToken          string
		ApiToken             string
		Backoff              bool
		ClassicOrg           bool
		ClientID             string
		Domain               string
		HttpProxy            string
		HttpTransport        http.RoundTripper
		LogLevel             int
		Logger               hclog.Logger
		MaxAPICapacity       int
		MaxWait              int
		MinWait              int
		OktaIDaaSClient      api.OktaIDaaSClient
		OktaGovernanceClient api.OktaGovernanceClient
		OrgName              string
		Parallelism          int
		PrivateKey           string
		PrivateKeyId         string
		QueriedWellKnown     bool
		RequestTimeout       int
		RetryCount           int
		Scopes               []string
		TimeOperations       TimeOperations
	}
)

func NewConfig(d *schema.ResourceData) *Config {
	config := Config{
		Backoff:        true,
		LogLevel:       int(hclog.Error),
		MaxAPICapacity: 100,
		MaxWait:        300,
		MinWait:        30,
		Parallelism:    1,
		RequestTimeout: 0,
		RetryCount:     5,
	}

	if val, ok := d.GetOk("org_name"); ok {
		config.OrgName = val.(string)
	}
	if config.OrgName == "" && os.Getenv("OKTA_ORG_NAME") != "" {
		config.OrgName = os.Getenv("OKTA_ORG_NAME")
	}

	if val, ok := d.GetOk("base_url"); ok {
		config.Domain = val.(string)
	}
	if config.Domain == "" {
		if os.Getenv("OKTA_BASE_URL") != "" {
			config.Domain = os.Getenv("OKTA_BASE_URL")
		}
	}

	if val, ok := d.GetOk("api_token"); ok {
		config.ApiToken = val.(string)
	}
	if config.ApiToken == "" && os.Getenv("OKTA_API_TOKEN") != "" {
		config.ApiToken = os.Getenv("OKTA_API_TOKEN")
	}

	if val, ok := d.GetOk("access_token"); ok {
		config.AccessToken = val.(string)
	}
	if config.AccessToken == "" && os.Getenv("OKTA_ACCESS_TOKEN") != "" {
		config.AccessToken = os.Getenv("OKTA_ACCESS_TOKEN")
	}

	if val, ok := d.GetOk("client_id"); ok {
		config.ClientID = val.(string)
	}
	if config.ClientID == "" && os.Getenv("OKTA_API_CLIENT_ID") != "" {
		config.ClientID = os.Getenv("OKTA_API_CLIENT_ID")
	}

	if val, ok := d.GetOk("private_key"); ok {
		config.PrivateKey = val.(string)
	}
	if config.PrivateKey == "" && os.Getenv("OKTA_API_PRIVATE_KEY") != "" {
		config.PrivateKey = os.Getenv("OKTA_API_PRIVATE_KEY")
	}

	if val, ok := d.GetOk("private_key_id"); ok {
		config.PrivateKeyId = val.(string)
	}
	if config.PrivateKeyId == "" && os.Getenv("OKTA_API_PRIVATE_KEY_ID") != "" {
		config.PrivateKeyId = os.Getenv("OKTA_API_PRIVATE_KEY_ID")
	}

	if val, ok := d.GetOk("scopes"); ok {
		config.Scopes = utils.ConvertInterfaceToStringSet(val)
	}
	if v := os.Getenv("OKTA_API_SCOPES"); v != "" && len(config.Scopes) == 0 {
		config.Scopes = strings.Split(v, ",")
	}

	if val, ok := d.GetOk("max_retries"); ok {
		config.RetryCount = val.(int)
	}

	if val, ok := d.GetOk("parallelism"); ok {
		config.Parallelism = val.(int)
	}

	if val, ok := d.GetOk("backoff"); ok {
		config.Backoff = val.(bool)
	}

	if val, ok := d.GetOk("min_wait_seconds"); ok {
		config.MinWait = val.(int)
	}

	if val, ok := d.GetOk("max_wait_seconds"); ok {
		config.MaxWait = val.(int)
	}

	if val, ok := d.GetOk("log_level"); ok {
		config.LogLevel = val.(int)
	}

	if val, ok := d.GetOk("request_timeout"); ok {
		config.RequestTimeout = val.(int)
	}

	if val, ok := d.GetOk("max_api_capacity"); ok {
		config.MaxAPICapacity = val.(int)
	}
	if config.MaxAPICapacity == 0 {
		if os.Getenv("MAX_API_CAPACITY") != "" {
			mac, err := strconv.ParseInt(os.Getenv("MAX_API_CAPACITY"), 10, 64)
			if err != nil {
				config.Logger.Error("error with max_api_capacity value", err)
			} else {
				config.MaxAPICapacity = int(mac)
			}
		}
	}

	if httpProxy, ok := d.Get("http_proxy").(string); ok {
		config.HttpProxy = httpProxy
	}
	if config.HttpProxy == "" && os.Getenv("OKTA_HTTP_PROXY") != "" {
		config.HttpProxy = os.Getenv("OKTA_HTTP_PROXY")
	}

	if v := os.Getenv("OKTA_API_SCOPES"); v != "" && len(config.Scopes) == 0 {
		config.Scopes = strings.Split(v, ",")
	}

	config.SetupLogger()

	return &config
}

func (c *Config) SetupLogger() {
	logLevel := hclog.Level(c.LogLevel)
	if os.Getenv("TF_LOG") != "" {
		logLevel = hclog.LevelFromString(os.Getenv("TF_LOG"))
	}

	c.Logger = hclog.New(&hclog.LoggerOptions{
		Level:      logLevel,
		TimeFormat: "2006/01/02 03:04:05",
	})
}

// IsClassicOrg returns true if the org is a classic org. Does lazy evaluation
// of the well known endpoint.
func (c *Config) IsClassicOrg(ctx context.Context) bool {
	if !c.QueriedWellKnown {
		// Discover if the Okta Org is Classic or OIE
		org, _, err := c.OktaIDaaSClient.OktaSDKClientV3().OrgSettingAPI.GetWellknownOrgMetadata(ctx).Execute()
		if err != nil {
			c.Logger.Error("error querying GET /.well-known/okta-organization", "error", err)
			return c.ClassicOrg
		}

		c.ClassicOrg = (org.GetPipeline() == "v1") // v1 == Classic, idx == OIE
		c.QueriedWellKnown = true
	}

	return c.ClassicOrg
}

func (c *Config) IsOAuth20Auth() bool {
	return c.PrivateKey != "" || c.AccessToken != ""
}

func (c *Config) SetTimeOperations(op TimeOperations) {
	c.TimeOperations = op
}

// LoadAPIClient initializes the Okta SDK clients
func (c *Config) LoadAPIClient() (err error) {
	iDaaSConfig := &api.OktaAPIConfig{
		AccessToken:    c.AccessToken,
		ApiToken:       c.ApiToken,
		Backoff:        c.Backoff,
		ClientID:       c.ClientID,
		Domain:         c.Domain,
		HttpProxy:      c.HttpProxy,
		Logger:         c.Logger,
		MaxAPICapacity: c.MaxAPICapacity,
		MaxWait:        c.MaxWait,
		MinWait:        c.MinWait,
		OrgName:        c.OrgName,
		PrivateKey:     c.PrivateKey,
		PrivateKeyId:   c.PrivateKeyId,
		RequestTimeout: c.RequestTimeout,
		RetryCount:     c.RetryCount,
		Scopes:         c.Scopes,
	}

	idaasClient, _ := api.NewOktaIDaaSAPIClient(iDaaSConfig)
	governanceClient, err := api.NewOktaGovernanceAPIClient(iDaaSConfig)
	if err != nil {
		return err
	}
	c.SetIdaasAPIClient(idaasClient)
	c.SetGovernanceAPIClient(governanceClient)
	return
}

// SetIdaasAPIClient allow other environments to inject an alternative idaas client to the config
func (c *Config) SetIdaasAPIClient(client api.OktaIDaaSClient) {
	c.OktaIDaaSClient = client
}

func (c *Config) SetGovernanceAPIClient(client api.OktaGovernanceClient) {
	c.OktaGovernanceClient = client
}

func (c *Config) VerifyCredentials(ctx context.Context) error {
	// NOTE: validate credentials during initial config with a call to
	// GET /api/v1/users/me
	// only for SSWS API token. Should we keep doing this?
	if c.ApiToken != "" {
		if _, _, err := c.OktaIDaaSClient.OktaSDKClientV3().UserAPI.GetUser(ctx, "me").Execute(); err != nil {
			return fmt.Errorf("error with v3 SDK client: %v", err)
		}
		if _, _, err := c.OktaIDaaSClient.OktaSDKClientV2().User.GetUser(ctx, "me"); err != nil {
			return fmt.Errorf("error with v2 SDK client: %v", err)
		}
	}

	return nil
}

type TimeOperations interface {
	DoNotRetry(error) bool
	Sleep(time.Duration)
}

type ProductionTimeOperations struct{}

// DoNotRetry always retry in production
func (o *ProductionTimeOperations) DoNotRetry(err error) bool {
	return false
}

// Sleep facade to actual time.Sleep in production
func (o *ProductionTimeOperations) Sleep(d time.Duration) {
	// lintignore:R018
	time.Sleep(d)
}

// NewProductionTimeOperations new production time operations
func NewProductionTimeOperations() TimeOperations {
	return &ProductionTimeOperations{}
}

type TestTimeOperations struct{}

// DoNotRetry tests do not retry when there is an error and VCR is recording
func (o *TestTimeOperations) DoNotRetry(err error) bool {
	return err != nil && os.Getenv("OKTA_VCR_TF_ACC") == "record"
}

// Sleep no sleeping when test is in VCR play mode
func (o *TestTimeOperations) Sleep(d time.Duration) {
	if os.Getenv("OKTA_VCR_TF_ACC") != "play" {
		// lintignore:R018
		time.Sleep(d)
	}
}

// NewTestTimeOperations new test time operations
func NewTestTimeOperations() TimeOperations {
	return &TestTimeOperations{}
}

package fwprovider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	schema_sdk "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/config"
	"github.com/okta/terraform-provider-okta/okta/services/governance"
	"github.com/okta/terraform-provider-okta/okta/services/idaas"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &FrameworkProvider{}
)

// NewFrameworkProvider is a helper function to simplify provider server and
// testing implementation.
func NewFrameworkProvider(version string, primary *schema_sdk.Provider) provider.Provider {
	return &FrameworkProvider{
		Config:            config.Config{},
		PluginSDKProvider: primary,
		Version:           version,
	}
}

type FrameworkProvider struct {
	Config            config.Config
	PluginSDKProvider *schema_sdk.Provider
	Version           string
}

type FrameworkProviderData struct {
	OrgName        types.String `tfsdk:"org_name"`
	AccessToken    types.String `tfsdk:"access_token"`
	APIToken       types.String `tfsdk:"api_token"`
	ClientID       types.String `tfsdk:"client_id"`
	Scopes         types.Set    `tfsdk:"scopes"`
	PrivateKey     types.String `tfsdk:"private_key"`
	PrivateKeyID   types.String `tfsdk:"private_key_id"`
	BaseURL        types.String `tfsdk:"base_url"`
	HTTPProxy      types.String `tfsdk:"http_proxy"`
	Backoff        types.Bool   `tfsdk:"backoff"`
	MinWaitSeconds types.Int64  `tfsdk:"min_wait_seconds"`
	MaxWaitSeconds types.Int64  `tfsdk:"max_wait_seconds"`
	MaxRetries     types.Int64  `tfsdk:"max_retries"`
	Parallelism    types.Int64  `tfsdk:"parallelism"`
	LogLevel       types.Int64  `tfsdk:"log_level"`
	MaxAPICapacity types.Int64  `tfsdk:"max_api_capacity"`
	RequestTimeout types.Int64  `tfsdk:"request_timeout"`
}

// Metadata returns the provider type name.
func (p *FrameworkProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "okta"
	resp.Version = p.Version
}

// Schema defines the provider-level schema for configuration data.
func (p *FrameworkProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"org_name": schema.StringAttribute{
				Optional:    true,
				Description: "The organization to manage in Okta.",
			},
			"access_token": schema.StringAttribute{
				Optional:    true,
				Description: "Bearer token granting privileges to Okta API.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("api_token"),
						path.MatchRoot("client_id"),
						path.MatchRoot("scopes"),
						path.MatchRoot("private_key"),
					}...),
				},
			},
			"api_token": schema.StringAttribute{
				Optional:    true,
				Description: "API Token granting privileges to Okta API.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("access_token"),
						path.MatchRoot("client_id"),
						path.MatchRoot("scopes"),
						path.MatchRoot("private_key"),
					}...),
				},
			},
			"client_id": schema.StringAttribute{
				Optional:    true,
				Description: "API Token granting privileges to Okta API.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("access_token"),
						path.MatchRoot("api_token"),
					}...),
				},
			},
			"scopes": schema.SetAttribute{
				Optional:    true,
				Description: "API Token granting privileges to Okta API.",
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("access_token"),
						path.MatchRoot("api_token"),
					}...),
				},
			},
			"private_key": schema.StringAttribute{
				Optional:    true,
				Description: "API Token granting privileges to Okta API.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("access_token"),
						path.MatchRoot("api_token"),
					}...),
				},
			},
			"private_key_id": schema.StringAttribute{
				Optional:    true,
				Description: "API Token Id granting privileges to Okta API.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("api_token"),
					}...),
				},
			},
			"base_url": schema.StringAttribute{
				Optional:    true,
				Description: "The Okta url. (Use 'oktapreview.com' for Okta testing)",
			},
			"http_proxy": schema.StringAttribute{
				Optional:    true,
				Description: "Alternate HTTP proxy of scheme://hostname or scheme://hostname:port format",
			},
			"backoff": schema.BoolAttribute{
				Optional:    true,
				Description: "Use exponential back off strategy for rate limits.",
			},
			"min_wait_seconds": schema.Int64Attribute{
				Optional:    true,
				Description: "minimum seconds to wait when rate limit is hit. We use exponential backoffs when backoff is enabled.",
			},
			"max_wait_seconds": schema.Int64Attribute{
				Optional:    true,
				Description: "maximum seconds to wait when rate limit is hit. We use exponential backoffs when backoff is enabled.",
			},
			"max_retries": schema.Int64Attribute{
				Optional:    true,
				Description: "maximum number of retries to attempt before erroring out.",
				Validators: []validator.Int64{
					int64validator.AtMost(100),
				},
			},
			"parallelism": schema.Int64Attribute{
				Optional:    true,
				Description: "Number of concurrent requests to make within a resource where bulk operations are not possible. Take note of https://developer.okta.com/docs/api/getting_started/rate-limits.",
			},
			"log_level": schema.Int64Attribute{
				Optional:    true,
				Description: "providers log level. Minimum is 1 (TRACE), and maximum is 5 (ERROR)",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
					int64validator.AtMost(5),
				},
			},
			"max_api_capacity": schema.Int64Attribute{
				Optional: true,
				Description: "Sets what percentage of capacity the provider can use of the total rate limit " +
					"capacity while making calls to the Okta management API endpoints. Okta API operates in one minute buckets. " +
					"See Okta Management API Rate Limits: https://developer.okta.com/docs/reference/rl-global-mgmt/",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
					int64validator.AtMost(100),
				},
			},
			"request_timeout": schema.Int64Attribute{
				Optional:    true,
				Description: "Timeout for single request (in seconds) which is made to Okta, the default is `0` (means no limit is set). The maximum value can be `300`.",
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
					int64validator.AtMost(300),
				},
			},
		},
	}
}

func (p *FrameworkProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var data FrameworkProviderData
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// reuse config from sdk provider
	rawMeta := p.PluginSDKProvider.Meta()
	if rawMeta == nil {
		resp.Diagnostics.AddError(
			"Provider Configuration Error",
			"The Okta provider was not properly configured. Please ensure valid credentials are provided. "+
				"Set 'api_token' (or OKTA_API_TOKEN env var), 'access_token' (or OKTA_ACCESS_TOKEN env var), "+
				"or 'private_key' + 'client_id' (or OKTA_API_PRIVATE_KEY + OKTA_API_CLIENT_ID env vars). "+
				"See https://registry.terraform.io/providers/okta/okta/latest/docs for more information.",
		)
		return
	}
	meta, ok := rawMeta.(*config.Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Configuration Type",
			fmt.Sprintf("Expected *config.Config, got: %T. Please report this issue to the provider developers.", rawMeta),
		)
		return
	}
	resp.EphemeralResourceData = meta
	resp.DataSourceData = meta
	resp.ResourceData = meta
}

// DataSources defines the data sources implemented in the provider.
func (p *FrameworkProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	var sources []func() datasource.DataSource
	sources = append(sources, idaas.FWProviderDataSources()...)
	sources = append(sources, governance.FWProviderDataSources()...)
	return sources
}

// Resources defines the resources implemented in the provider.
func (p *FrameworkProvider) Resources(_ context.Context) []func() resource.Resource {
	var resources []func() resource.Resource

	// Append resources from various modules
	resources = append(resources, idaas.FWProviderResources()...)
	resources = append(resources, governance.FWProviderResources()...)

	return resources
}

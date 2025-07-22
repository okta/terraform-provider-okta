package idaas

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

// OAuth2ServerMetadata represents the OAuth 2.0 authorization server metadata
type OAuth2ServerMetadata struct {
	Issuer                                                    string   `json:"issuer"`
	AuthorizationEndpoint                                     string   `json:"authorization_endpoint"`
	TokenEndpoint                                             string   `json:"token_endpoint"`
	RegistrationEndpoint                                      string   `json:"registration_endpoint"`
	ResponseTypesSupported                                    []string `json:"response_types_supported"`
	ResponseModesSupported                                    []string `json:"response_modes_supported"`
	GrantTypesSupported                                       []string `json:"grant_types_supported"`
	SubjectTypesSupported                                     []string `json:"subject_types_supported"`
	ScopesSupported                                           []string `json:"scopes_supported"`
	TokenEndpointAuthMethodsSupported                         []string `json:"token_endpoint_auth_methods_supported"`
	ClaimsSupported                                           []string `json:"claims_supported"`
	CodeChallengeMethodsSupported                             []string `json:"code_challenge_methods_supported"`
	IntrospectionEndpoint                                     string   `json:"introspection_endpoint"`
	IntrospectionEndpointAuthMethodsSupported                 []string `json:"introspection_endpoint_auth_methods_supported"`
	RevocationEndpoint                                        string   `json:"revocation_endpoint"`
	RevocationEndpointAuthMethodsSupported                    []string `json:"revocation_endpoint_auth_methods_supported"`
	EndSessionEndpoint                                        string   `json:"end_session_endpoint"`
	RequestParameterSupported                                 bool     `json:"request_parameter_supported"`
	RequestObjectSigningAlgValuesSupported                    []string `json:"request_object_signing_alg_values_supported"`
	DeviceAuthorizationEndpoint                               string   `json:"device_authorization_endpoint"`
	PushedAuthorizationRequestEndpoint                        string   `json:"pushed_authorization_request_endpoint"`
	BackchannelTokenDeliveryModesSupported                    []string `json:"backchannel_token_delivery_modes_supported"`
	BackchannelAuthenticationRequestSigningAlgValuesSupported []string `json:"backchannel_authentication_request_signing_alg_values_supported"`
	DpopSigningAlgValuesSupported                             []string `json:"dpop_signing_alg_values_supported"`
}

// Ensure the implementation satisfies the expected interfaces
var _ datasource.DataSource = &OAuthAuthorizationServerDataSource{}

func NewOAuthAuthorizationServerDataSource() datasource.DataSource {
	return &OAuthAuthorizationServerDataSource{}
}

type OAuthAuthorizationServerDataSource struct {
	*config.Config
}

func (d *OAuthAuthorizationServerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_authorization_server"
}

// OAuthAuthorizationServerModel describes the data source data model.
type OAuthAuthorizationServerModel struct {
	ID                                                        types.String   `tfsdk:"id"`
	BaseURL                                                   types.String   `tfsdk:"base_url"`
	Issuer                                                    types.String   `tfsdk:"issuer"`
	AuthorizationEndpoint                                     types.String   `tfsdk:"authorization_endpoint"`
	TokenEndpoint                                             types.String   `tfsdk:"token_endpoint"`
	RegistrationEndpoint                                      types.String   `tfsdk:"registration_endpoint"`
	ResponseTypesSupported                                    []types.String `tfsdk:"response_types_supported"`
	ResponseModesSupported                                    []types.String `tfsdk:"response_modes_supported"`
	GrantTypesSupported                                       []types.String `tfsdk:"grant_types_supported"`
	SubjectTypesSupported                                     []types.String `tfsdk:"subject_types_supported"`
	ScopesSupported                                           []types.String `tfsdk:"scopes_supported"`
	TokenEndpointAuthMethodsSupported                         []types.String `tfsdk:"token_endpoint_auth_methods_supported"`
	ClaimsSupported                                           []types.String `tfsdk:"claims_supported"`
	CodeChallengeMethodsSupported                             []types.String `tfsdk:"code_challenge_methods_supported"`
	IntrospectionEndpoint                                     types.String   `tfsdk:"introspection_endpoint"`
	IntrospectionEndpointAuthMethodsSupported                 []types.String `tfsdk:"introspection_endpoint_auth_methods_supported"`
	RevocationEndpoint                                        types.String   `tfsdk:"revocation_endpoint"`
	RevocationEndpointAuthMethodsSupported                    []types.String `tfsdk:"revocation_endpoint_auth_methods_supported"`
	EndSessionEndpoint                                        types.String   `tfsdk:"end_session_endpoint"`
	RequestParameterSupported                                 types.Bool     `tfsdk:"request_parameter_supported"`
	RequestObjectSigningAlgValuesSupported                    []types.String `tfsdk:"request_object_signing_alg_values_supported"`
	DeviceAuthorizationEndpoint                               types.String   `tfsdk:"device_authorization_endpoint"`
	PushedAuthorizationRequestEndpoint                        types.String   `tfsdk:"pushed_authorization_request_endpoint"`
	BackchannelTokenDeliveryModesSupported                    []types.String `tfsdk:"backchannel_token_delivery_modes_supported"`
	BackchannelAuthenticationRequestSigningAlgValuesSupported []types.String `tfsdk:"backchannel_authentication_request_signing_alg_values_supported"`
	DpopSigningAlgValuesSupported                             []types.String `tfsdk:"dpop_signing_alg_values_supported"`
}

func (d *OAuthAuthorizationServerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Returns OpenID Connect metadata for the Okta org authorization server. Clients use this information to programmatically configure their interactions with Okta",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The authorization server's issuer identifier",
			},
			"base_url": schema.StringAttribute{
				Description: "The base URL of the Okta org",
				Optional:    true,
			},
			"issuer": schema.StringAttribute{
				Computed:    true,
				Description: "The authorization server's issuer identifier",
			},
			"authorization_endpoint": schema.StringAttribute{
				Computed:    true,
				Description: "URL of the authorization server's authorization endpoint",
			},
			"token_endpoint": schema.StringAttribute{
				Computed:    true,
				Description: "URL of the authorization server's token endpoint",
			},
			"registration_endpoint": schema.StringAttribute{
				Computed:    true,
				Description: "URL of the authorization server's dynamic client registration endpoint",
			},
			"response_types_supported": schema.ListAttribute{
				Computed:    true,
				Description: "JSON array containing a list of the OAuth 2.0 response_type values that this authorization server supports",
				ElementType: types.StringType,
			},
			"response_modes_supported": schema.ListAttribute{
				Computed:    true,
				Description: "JSON array containing a list of the OAuth 2.0 response_mode values that this authorization server supports",
				ElementType: types.StringType,
			},
			"grant_types_supported": schema.ListAttribute{
				Computed:    true,
				Description: "JSON array containing a list of the OAuth 2.0 grant type values that this authorization server supports",
				ElementType: types.StringType,
			},
			"subject_types_supported": schema.ListAttribute{
				Computed:    true,
				Description: "JSON array containing a list of the Subject Identifier types that this authorization server supports",
				ElementType: types.StringType,
			},
			"scopes_supported": schema.ListAttribute{
				Computed:    true,
				Description: "JSON array containing a list of the OAuth 2.0 scope values that this authorization server supports",
				ElementType: types.StringType,
			},
			"token_endpoint_auth_methods_supported": schema.ListAttribute{
				Computed:    true,
				Description: "JSON array containing a list of client authentication methods supported by this token endpoint",
				ElementType: types.StringType,
			},
			"claims_supported": schema.ListAttribute{
				Computed:    true,
				Description: "JSON array containing a list of the Claim Names of the Claims that the OpenID Provider MAY be able to supply values for",
				ElementType: types.StringType,
			},
			"code_challenge_methods_supported": schema.ListAttribute{
				Computed:    true,
				Description: "JSON array containing a list of Proof Key for Code Exchange (PKCE) code challenge methods supported by this authorization server",
				ElementType: types.StringType,
			},
			"introspection_endpoint": schema.StringAttribute{
				Computed:    true,
				Description: "URL of the authorization server's introspection endpoint",
			},
			"introspection_endpoint_auth_methods_supported": schema.ListAttribute{
				Computed:    true,
				Description: "JSON array containing a list of client authentication methods supported by this introspection endpoint",
				ElementType: types.StringType,
			},
			"revocation_endpoint": schema.StringAttribute{
				Computed:    true,
				Description: "URL of the authorization server's revocation endpoint",
			},
			"revocation_endpoint_auth_methods_supported": schema.ListAttribute{
				Computed:    true,
				Description: "JSON array containing a list of client authentication methods supported by this revocation endpoint",
				ElementType: types.StringType,
			},
			"end_session_endpoint": schema.StringAttribute{
				Computed:    true,
				Description: "URL at the authorization server to which an RP can perform a redirect to request that the End-User be logged out at the authorization server",
			},
			"request_parameter_supported": schema.BoolAttribute{
				Computed:    true,
				Description: "Boolean value specifying whether the authorization server supports use of the request parameter",
			},
			"request_object_signing_alg_values_supported": schema.ListAttribute{
				Computed:    true,
				Description: "JSON array containing a list of the JWS signing algorithms (alg values) supported by the authorization server for the request object",
				ElementType: types.StringType,
			},
			"device_authorization_endpoint": schema.StringAttribute{
				Computed:    true,
				Description: "URL of the authorization server's device authorization endpoint",
			},
			"pushed_authorization_request_endpoint": schema.StringAttribute{
				Computed:    true,
				Description: "URL of the authorization server's pushed authorization request endpoint",
			},
			"backchannel_token_delivery_modes_supported": schema.ListAttribute{
				Computed:    true,
				Description: "JSON array containing a list of the CIBA backchannel token delivery modes supported by this authorization server",
				ElementType: types.StringType,
			},
			"backchannel_authentication_request_signing_alg_values_supported": schema.ListAttribute{
				Computed:    true,
				Description: "JSON array containing a list of the JWS signing algorithms (alg values) supported by the authorization server for the backchannel authentication request",
				ElementType: types.StringType,
			},
			"dpop_signing_alg_values_supported": schema.ListAttribute{
				Computed:    true,
				Description: "JSON array containing a list of the JWS signing algorithms (alg values) supported by the authorization server for DPoP proof JWTs",
				ElementType: types.StringType,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *OAuthAuthorizationServerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

// Read refreshes the Terraform state with the latest data.
func (d *OAuthAuthorizationServerDataSource) Read(ctx context.Context, readReq datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state OAuthAuthorizationServerModel

	// Get the base URL from either the data source config or provider config
	var baseURL string
	var diags diag.Diagnostics

	// First, try to get base URL from data source config
	var dataSourceConfig OAuthAuthorizationServerModel
	diags = readReq.Config.Get(ctx, &dataSourceConfig)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !dataSourceConfig.BaseURL.IsNull() && !dataSourceConfig.BaseURL.IsUnknown() {
		// Use the provided base_url from the data source
		baseURL = dataSourceConfig.BaseURL.ValueString()
	} else {
		// Use the provider's configured base URL
		if d.OrgName == "" || d.Domain == "" {
			resp.Diagnostics.AddError("Missing base URL", "Either the provider's org_name and base_url must be configured, or the data source's base_url attribute must be set")
			return
		}
		baseURL = fmt.Sprintf("https://%s.%s", d.OrgName, d.Domain)
	}

	// Create a new HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/.well-known/oauth-authorization-server", nil)
	if err != nil {
		resp.Diagnostics.AddError("Error creating request", err.Error())
		return
	}

	// Make the request
	client := &http.Client{}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Error fetching OAuth authorization server metadata", err.Error())
		return
	}
	defer httpResp.Body.Close()

	// Check for non-200 status codes
	if httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError("Error response from server", fmt.Sprintf("Received status code %d", httpResp.StatusCode))
		return
	}

	var metadata OAuth2ServerMetadata
	if err := json.NewDecoder(httpResp.Body).Decode(&metadata); err != nil {
		resp.Diagnostics.AddError("Error decoding OAuth authorization server metadata", err.Error())
		return
	}

	state.ID = types.StringValue(metadata.Issuer)
	state.Issuer = types.StringValue(metadata.Issuer)
	state.AuthorizationEndpoint = types.StringValue(metadata.AuthorizationEndpoint)
	state.TokenEndpoint = types.StringValue(metadata.TokenEndpoint)
	state.RegistrationEndpoint = types.StringValue(metadata.RegistrationEndpoint)
	state.ResponseTypesSupported = utils.ConvertStringSlice(metadata.ResponseTypesSupported)
	state.ResponseModesSupported = utils.ConvertStringSlice(metadata.ResponseModesSupported)
	state.GrantTypesSupported = utils.ConvertStringSlice(metadata.GrantTypesSupported)
	state.SubjectTypesSupported = utils.ConvertStringSlice(metadata.SubjectTypesSupported)
	state.ScopesSupported = utils.ConvertStringSlice(metadata.ScopesSupported)
	state.TokenEndpointAuthMethodsSupported = utils.ConvertStringSlice(metadata.TokenEndpointAuthMethodsSupported)
	state.ClaimsSupported = utils.ConvertStringSlice(metadata.ClaimsSupported)
	state.CodeChallengeMethodsSupported = utils.ConvertStringSlice(metadata.CodeChallengeMethodsSupported)
	state.IntrospectionEndpoint = types.StringValue(metadata.IntrospectionEndpoint)
	state.IntrospectionEndpointAuthMethodsSupported = utils.ConvertStringSlice(metadata.IntrospectionEndpointAuthMethodsSupported)
	state.RevocationEndpoint = types.StringValue(metadata.RevocationEndpoint)
	state.RevocationEndpointAuthMethodsSupported = utils.ConvertStringSlice(metadata.RevocationEndpointAuthMethodsSupported)
	state.EndSessionEndpoint = types.StringValue(metadata.EndSessionEndpoint)
	state.RequestParameterSupported = types.BoolValue(metadata.RequestParameterSupported)
	state.RequestObjectSigningAlgValuesSupported = utils.ConvertStringSlice(metadata.RequestObjectSigningAlgValuesSupported)
	state.DeviceAuthorizationEndpoint = types.StringValue(metadata.DeviceAuthorizationEndpoint)
	state.PushedAuthorizationRequestEndpoint = types.StringValue(metadata.PushedAuthorizationRequestEndpoint)
	state.BackchannelTokenDeliveryModesSupported = utils.ConvertStringSlice(metadata.BackchannelTokenDeliveryModesSupported)
	state.BackchannelAuthenticationRequestSigningAlgValuesSupported = utils.ConvertStringSlice(metadata.BackchannelAuthenticationRequestSigningAlgValuesSupported)
	state.DpopSigningAlgValuesSupported = utils.ConvertStringSlice(metadata.DpopSigningAlgValuesSupported)

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

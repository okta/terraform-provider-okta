package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	tfpath "github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	v5okta "github.com/okta/okta-sdk-golang/v5/okta"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ resource.Resource                = &appConnections{}
	_ resource.ResourceWithConfigure   = &appConnections{}
	_ resource.ResourceWithImportState = &appConnections{}
)

type appConnections struct {
	*config.Config
}

func newAppConnectionsResource() resource.Resource {
	return &appConnections{}
}

type SigningModel struct {
	RotationMode types.String `tfsdk:"rotation_mode"`
}

type SettingsModel struct {
	AdminUsername types.String `tfsdk:"admin_username"`
	AdminPassword types.String `tfsdk:"admin_password"`
}

type ProfileModel struct {
	AuthScheme types.String   `tfsdk:"auth_scheme"`
	Token      types.String   `tfsdk:"token"`
	ClientID   types.String   `tfsdk:"client_id"`
	Signing    *SigningModel  `tfsdk:"signing"`
	Settings   *SettingsModel `tfsdk:"settings"`
}

type AppConnectionsModel struct {
	ID      types.String  `tfsdk:"id"`
	BaseURL types.String  `tfsdk:"base_url"`
	Profile *ProfileModel `tfsdk:"profile"`
	Status  types.String  `tfsdk:"status"`
	Action  types.String  `tfsdk:"action"`
}

func (r *appConnections) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_connection"
}

func (r *appConnections) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *appConnections) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, tfpath.Root("id"), req, resp)
}

func (r *appConnections) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages Okta App Connection configurations for provisioning.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The application ID.",
			},
			"base_url": schema.StringAttribute{
				Required:    true,
				Description: "The base URL for the provisioning connection.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Provisioning connection status.",
			},
			"action": schema.StringAttribute{
				Required:    true,
				Description: "The action to perform on the connection. Valid values are `activate` or `deactivate`.",
				Validators: []validator.String{
					stringvalidator.OneOf("activate", "deactivate"),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"profile": schema.SingleNestedBlock{
				Description: "Profile configuration for the app connection.",
				Attributes: map[string]schema.Attribute{
					"auth_scheme": schema.StringAttribute{
						Required:    true,
						Description: "Authentication scheme. Valid values are TOKEN or OAUTH2.",
						Validators: []validator.String{
							stringvalidator.OneOf("TOKEN", "OAUTH2"),
						},
					},
					"token": schema.StringAttribute{
						Optional:    true,
						Sensitive:   true,
						Description: "Authentication token (required for TOKEN auth scheme).",
					},
					"client_id": schema.StringAttribute{
						Optional:    true,
						Description: "OAuth2 client ID (required for OAUTH2 auth scheme).",
					},
				},
				Blocks: map[string]schema.Block{
					"signing": schema.SingleNestedBlock{
						Description: "Signing configuration.",
						Attributes: map[string]schema.Attribute{
							"rotation_mode": schema.StringAttribute{
								Optional:    true,
								Description: "Token rotation mode.",
								Validators: []validator.String{
									stringvalidator.OneOf("AUTO", "MANUAL"),
								},
							},
						},
					},
					"settings": schema.SingleNestedBlock{
						Description: "Additional settings for OAuth2 authentication.",
						Attributes: map[string]schema.Attribute{
							"admin_username": schema.StringAttribute{
								Optional:    true,
								Description: "Admin username for OAuth2.",
							},
							"admin_password": schema.StringAttribute{
								Optional:    true,
								Sensitive:   true,
								Description: "Admin password for OAuth2.",
							},
						},
					},
				},
			},
		},
	}
}

func (r *appConnections) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data AppConnectionsModel

	// Read Terraform plan data into the model - FIX: Use req.Plan.Get instead of req.Config.Get
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create/Update the app connection
	updateDefaultConnection, _, err := r.OktaIDaaSClient.OktaSDKClientV5().ApplicationConnectionsAPI.UpdateDefaultProvisioningConnectionForApplication(ctx, data.ID.ValueString()).UpdateDefaultProvisioningConnectionForApplicationRequest(buildDefaultProvisioningConnections(data)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating App Connection",
			"Could not create app connection for application ID "+data.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	if updateDefaultConnection == nil {
		resp.Diagnostics.AddError(
			"No App Connection Returned",
			"API call succeeded but no app connection was returned",
		)
		return
	}

	if data.Action.ValueString() == "activate" {
		_, err = r.OktaIDaaSClient.OktaSDKClientV5().ApplicationConnectionsAPI.ActivateDefaultProvisioningConnectionForApplication(ctx, data.ID.ValueString()).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Activating App Connection",
				"Could not create app connection for application ID "+data.ID.ValueString()+": "+err.Error(),
			)
			return
		}
		data.Status = types.StringValue("ENABLED")
		data.Action = types.StringValue("activate")
	}

	// Update state with response
	resp.Diagnostics.Append(updateAppConnectionState(&data, updateDefaultConnection)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func updateAppConnectionState(data *AppConnectionsModel, response *v5okta.ProvisioningConnectionResponse) diag.Diagnostics {
	var diags diag.Diagnostics

	// Update base URL
	data.BaseURL = types.StringValue(response.GetBaseUrl())

	// Initialize Profile if nil
	if data.Profile == nil {
		data.Profile = &ProfileModel{}
	}

	// Update auth scheme
	data.Profile.AuthScheme = types.StringValue(response.Profile.AuthScheme)

	// Note: Token and other sensitive fields are typically not returned by the API for security reasons
	// The state will maintain the values from the plan/config

	return diags
}

// FIX: Return 'req' instead of 'x'
func buildDefaultProvisioningConnections(data AppConnectionsModel) v5okta.UpdateDefaultProvisioningConnectionForApplicationRequest {
	var req v5okta.UpdateDefaultProvisioningConnectionForApplicationRequest
	authScheme := data.Profile.AuthScheme.ValueString()

	switch authScheme {
	case "TOKEN":
		req.ProvisioningConnectionTokenRequest = &v5okta.ProvisioningConnectionTokenRequest{
			BaseUrl: data.BaseURL.ValueStringPointer(),
			Profile: v5okta.ProvisioningConnectionTokenRequestProfile{
				AuthScheme: authScheme,
				Token:      data.Profile.Token.ValueStringPointer(),
			},
		}

	case "OAUTH2":
		req.ProvisioningConnectionOauthRequest = &v5okta.ProvisioningConnectionOauthRequest{
			Profile: v5okta.ProvisioningConnectionOauthRequestProfile{
				AuthScheme: authScheme,
				ClientId:   data.Profile.ClientID.ValueStringPointer(),
			},
		}

		// Add settings if provided
		if data.Profile.Settings != nil {
			req.ProvisioningConnectionOauthRequest.Profile.Settings.AdminUsername = data.Profile.Settings.AdminUsername.ValueString()
			req.ProvisioningConnectionOauthRequest.Profile.Settings.AdminPassword = data.Profile.Settings.AdminPassword.ValueString()
		}
	}

	return req // FIX: Return 'req' instead of 'x'
}

func (r *appConnections) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data AppConnectionsModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readDefaultProvisioningConnections, _, err := r.OktaIDaaSClient.OktaSDKClientV5().ApplicationConnectionsAPI.GetDefaultProvisioningConnectionForApplication(ctx, data.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading App Connection",
			"Could not read app connection for application ID "+data.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	if readDefaultProvisioningConnections == nil {
		resp.Diagnostics.AddError(
			"App Connection Not Found",
			"App connection not found for application ID "+data.ID.ValueString(),
		)
		return
	}

	resp.Diagnostics.Append(updateAppConnectionState(&data, readDefaultProvisioningConnections)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(setStatus(&data, readDefaultProvisioningConnections)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func setStatus(d *AppConnectionsModel, resp *v5okta.ProvisioningConnectionResponse) diag.Diagnostics {
	var diags diag.Diagnostics

	if resp.Status != "" {
		status := resp.GetStatus()
		d.Status = types.StringValue(status)

		// Set action based on current status
		switch status {
		case "ENABLED":
			d.Action = types.StringValue("activate")
		case "DISABLED":
			d.Action = types.StringValue("deactivate")
		default:
			// For unknown status, don't set action to avoid confusion
			d.Action = types.StringValue("")
		}
	} else {
		// If status is not available from API, try to determine from other fields
		// This is a fallback - you might need to adjust based on actual API response
		d.Status = types.StringValue("UNKNOWN")
		d.Action = types.StringValue("")
	}
	return diags
}

func (r *appConnections) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data AppConnectionsModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateDefaultConnection, _, err := r.OktaIDaaSClient.OktaSDKClientV5().ApplicationConnectionsAPI.UpdateDefaultProvisioningConnectionForApplication(ctx, data.ID.ValueString()).UpdateDefaultProvisioningConnectionForApplicationRequest(buildDefaultProvisioningConnections(data)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating App Connection",
			"Could not update app connection for application ID "+data.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	if updateDefaultConnection == nil {
		resp.Diagnostics.AddError(
			"No App Connection Returned",
			"API call succeeded but no app connection was returned",
		)
		return
	}

	if data.Action.ValueString() == "activate" {
		_, err = r.OktaIDaaSClient.OktaSDKClientV5().ApplicationConnectionsAPI.ActivateDefaultProvisioningConnectionForApplication(ctx, data.ID.ValueString()).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Activating App Connection",
				"Could not create app connection for application ID "+data.ID.ValueString()+": "+err.Error(),
			)
			return
		}
		data.Status = types.StringValue("ENABLED")
		data.Action = types.StringValue("activate")
	} else if data.Action.ValueString() == "deactivate" {
		_, err = r.OktaIDaaSClient.OktaSDKClientV5().ApplicationConnectionsAPI.DeactivateDefaultProvisioningConnectionForApplication(ctx, data.ID.ValueString()).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Deactivating App Connection",
				"Could not deactivate app connection for application ID "+data.ID.ValueString()+": "+err.Error(),
			)
			return
		}
		data.Status = types.StringValue("DISABLED")
		data.Action = types.StringValue("deactivate")
	}

	resp.Diagnostics.Append(updateAppConnectionState(&data, updateDefaultConnection)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *appConnections) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning(
		"Delete Not Supported",
		"This resource cannot be deleted via Terraform. App connections are managed through the application lifecycle.",
	)
}

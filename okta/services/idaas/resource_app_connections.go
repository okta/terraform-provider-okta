package idaas

import (
	"context"
	tfpath "github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	v5okta "github.com/okta/okta-sdk-golang/v5/okta"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ resource.Resource                = &appConnections{}
	_ resource.ResourceWithConfigure   = &appConnections{}
	_ resource.ResourceWithImportState = &appConnections{}
)

var _ resource.Resource = &appConnections{}

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
	ClientId   types.String   `tfsdk:"client_id"`
	Signing    *SigningModel  `tfsdk:"signing"`
	Settings   *SettingsModel `tfsdk:"settings"`
}

type AppConnectionsModel struct {
	Id      types.String  `tfsdk:"id"`
	BaseUrl types.String  `tfsdk:"base_url"`
	Profile *ProfileModel `tfsdk:"profile"`
}

func (r *appConnections) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_connections"
}

func (r *appConnections) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *appConnections) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, tfpath.Root("id"), req, resp)
}

func (r *appConnections) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required: true,
			},
			"base_url": schema.StringAttribute{
				Required: true,
			},
		},
		Blocks: map[string]schema.Block{
			"profile": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"auth_scheme": schema.StringAttribute{
						Required: true,
					},
					"token": schema.StringAttribute{
						Optional: true,
					},
					"client_id": schema.StringAttribute{
						Optional: true,
					},
				},
				Blocks: map[string]schema.Block{
					"signing": schema.SingleNestedBlock{
						Attributes: map[string]schema.Attribute{
							"rotation_mode": schema.StringAttribute{
								Optional: true,
							},
						},
					},
					"settings": schema.SingleNestedBlock{
						Attributes: map[string]schema.Attribute{
							"admin_username": schema.StringAttribute{
								Optional: true,
							},
							"admin_password": schema.StringAttribute{
								Optional: true,
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
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateDefaultConnection, _, err := r.OktaIDaaSClient.OktaSDKClientV5().ApplicationConnectionsAPI.UpdateDefaultProvisioningConnectionForApplication(ctx, data.Id.ValueString()).UpdateDefaultProvisioningConnectionForApplicationRequest(buildDefaultProvisioningConnections(data)).Execute()
	if err != nil {
		return
	}

	updateAppConnectionState(&data, updateDefaultConnection)

	// Save Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func updateAppConnectionState(data *AppConnectionsModel, updateDefaultConnection *v5okta.ProvisioningConnectionResponse) {
	data.BaseUrl = types.StringValue(updateDefaultConnection.GetBaseUrl())
	data.Profile.AuthScheme = types.StringValue(updateDefaultConnection.Profile.AuthScheme)
}

func buildDefaultProvisioningConnections(data AppConnectionsModel) v5okta.UpdateDefaultProvisioningConnectionForApplicationRequest {
	var x v5okta.UpdateDefaultProvisioningConnectionForApplicationRequest
	x.ProvisioningConnectionTokenRequest = &v5okta.ProvisioningConnectionTokenRequest{
		Profile: v5okta.ProvisioningConnectionTokenRequestProfile{
			AuthScheme: data.Profile.AuthScheme.ValueString(),
			Token:      data.Profile.Token.ValueStringPointer(),
		},
	}

	if data.BaseUrl.ValueStringPointer() != nil {
		x.ProvisioningConnectionTokenRequest.BaseUrl = data.BaseUrl.ValueStringPointer()
	}

	x.ProvisioningConnectionOauthRequest = &v5okta.ProvisioningConnectionOauthRequest{
		Profile: v5okta.ProvisioningConnectionOauthRequestProfile{
			AuthScheme: data.Profile.AuthScheme.ValueString(),
			ClientId:   data.Profile.ClientId.ValueStringPointer(),
		},
	}

	if data.Profile.Settings != nil {
		x.ProvisioningConnectionOauthRequest.Profile.Settings.AdminUsername = data.Profile.Settings.AdminUsername.ValueString()
		x.ProvisioningConnectionOauthRequest.Profile.Settings.AdminPassword = data.Profile.Settings.AdminPassword.ValueString()
	}
	return x
}

func (r *appConnections) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data AppConnectionsModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readDefaultProvisioningConnections, _, err := r.OktaIDaaSClient.OktaSDKClientV5().ApplicationConnectionsAPI.GetDefaultProvisioningConnectionForApplication(ctx, data.Id.ValueString()).Execute()
	if err != nil {
		return
	}

	updateAppConnectionState(&data, readDefaultProvisioningConnections)
	// Save Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *appConnections) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data AppConnectionsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateDefaultConnection, _, err := r.OktaIDaaSClient.OktaSDKClientV5().ApplicationConnectionsAPI.UpdateDefaultProvisioningConnectionForApplication(ctx, data.Id.ValueString()).UpdateDefaultProvisioningConnectionForApplicationRequest(buildDefaultProvisioningConnections(data)).Execute()
	if err != nil {
		return
	}

	updateAppConnectionState(&data, updateDefaultConnection)

	// Save Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *appConnections) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning(
		"Delete Not Supported",
		"This resource cannot be deleted via Terraform.",
	)
}

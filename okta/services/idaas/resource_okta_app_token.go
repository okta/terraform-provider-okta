package idaas

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ resource.Resource                = &appTokenResource{}
	_ resource.ResourceWithConfigure   = &appTokenResource{}
	_ resource.ResourceWithImportState = &appTokenResource{}
)

type appTokenResource struct {
	*config.Config
}

func newAppTokenResource() resource.Resource {
	return &appTokenResource{}
}

func (r *appTokenResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *appTokenResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Expected format: client_id/id",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("client_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}

type appTokenModel struct {
	ID       types.String `tfsdk:"id"`
	ClientID types.String `tfsdk:"client_id"`
	UserID   types.String `tfsdk:"user_id"`
	Status   types.String `tfsdk:"status"`
}

func (r *appTokenResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_token"
}

func (r *appTokenResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The unique Okta ID of this key record",
			},
			"client_id": schema.StringAttribute{
				Required:    true,
				Description: "The unique Okta ID of the application associated with this token.",
			},
			"user_id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique Okta ID of the user associated with this token.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the token.",
			},
		},
	}
}

func (r *appTokenResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	resp.Diagnostics.AddError(
		"Create Not Supported",
		"Application Token cannot be created via Terraform.",
	)
}

func (r *appTokenResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data appTokenModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...) // Read prior state data into the model
	if resp.Diagnostics.HasError() {
		return
	}

	getAppTokenResp, _, err := r.OktaIDaaSClient.OktaSDKClientV5().ApplicationTokensAPI.GetOAuth2TokenForApplication(ctx, data.ClientID.ValueString(), data.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error reading application token", "Could not read application token, unexpected error: "+err.Error())
		return
	}

	data.ID = types.StringValue(getAppTokenResp.GetId())
	data.ClientID = types.StringValue(getAppTokenResp.GetClientId())
	data.UserID = types.StringValue(getAppTokenResp.GetUserId())
	data.Status = types.StringValue(getAppTokenResp.GetStatus())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *appTokenResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddWarning(
		"Update Not Supported",
		"Application Token cannot be updated via Terraform. Terraform will retain the existing state.",
	)
}

func (r *appTokenResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data appTokenModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...) // Read prior state data into the model
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.OktaIDaaSClient.OktaSDKClientV5().ApplicationTokensAPI.RevokeOAuth2TokenForApplication(ctx, data.ClientID.ValueString(), data.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error revoking application token", "Could not revoke application token, unexpected error: "+err.Error())
		return
	}
}

package idaas

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	tfpath "github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	v5okta "github.com/okta/okta-sdk-golang/v5/okta"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ resource.Resource                = &identitySourceResource{}
	_ resource.ResourceWithConfigure   = &identitySourceResource{}
	_ resource.ResourceWithImportState = &identitySourceResource{}
)

func newIdentitySourceResource() resource.Resource {
	return &identitySourceResource{}
}

// pushProviderResource defines the resource implementation
type identitySourceResource struct {
	*config.Config
}

type identitySourceResourceModel struct {
	Id               types.String `tfsdk:"id"`
	IdentitySourceId types.String `tfsdk:"identity_source_id"` // Indicates the minimum required SKU to manage the campaign. Values can be `BASIC` and `PREMIUM`.
	ImportType       types.String `tfsdk:"import_type"`
	Status           types.String `tfsdk:"status"`
}

func (r *identitySourceResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(request, response)
}

func (r *identitySourceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_identity_source"
}

func (r *identitySourceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the identity source session.",
			},
			"identity_source_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the custom identity source for which the session is created.",
			},
			"import_type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of import. All imports are `INCREMENTAL` imports.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The current status of the identity source session.",
			},
		},
	}
}

func (r *identitySourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data identitySourceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	createIdentityResourceSession, _, err := r.OktaIDaaSClient.OktaSDKClientV5().IdentitySourceAPI.CreateIdentitySourceSession(ctx, data.IdentitySourceId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Identity Source Session",
			fmt.Sprintf("Could not create identity source session for identity source ID %s: %v", data.IdentitySourceId.ValueString(), err),
		)
		return
	}

	// Check if we got any sessions back
	if len(createIdentityResourceSession) == 0 {
		resp.Diagnostics.AddError(
			"No Identity Source Session Returned",
			"API call succeeded but no identity source session was returned",
		)
		return
	}

	fmt.Println("data len", len(createIdentityResourceSession))

	resp.Diagnostics.Append(applyIdentitySourceToState(ctx, createIdentityResourceSession[0], &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func applyIdentitySourceToState(ctx context.Context, session v5okta.IdentitySourceSession, data *identitySourceResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	data.Id = types.StringValue(session.GetId())
	data.ImportType = types.StringValue(session.GetImportType())
	data.Status = types.StringValue(session.GetStatus())
	data.IdentitySourceId = types.StringValue(session.GetIdentitySourceId())
	return diags
}

func (r *identitySourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data identitySourceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	getIdentityResourceSession, _, err := r.OktaIDaaSClient.OktaSDKClientV5().IdentitySourceAPI.GetIdentitySourceSession(ctx, data.IdentitySourceId.ValueString(), data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Identity Source Session",
			fmt.Sprintf("Could not read identity source session %s for identity source ID %s: %v", data.Id.ValueString(), data.IdentitySourceId.ValueString(), err),
		)
		return
	}

	if getIdentityResourceSession == nil {
		resp.Diagnostics.AddError(
			"Identity Source Session Not Found",
			fmt.Sprintf("Identity source session %s not found", data.Id.ValueString()),
		)
		return
	}

	resp.Diagnostics.Append(applyIdentitySourceToState(ctx, *getIdentityResourceSession, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *identitySourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddWarning(
		"Update Not Supported",
		"This resource cannot be updated via Terraform.",
	)
}

func (r *identitySourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data identitySourceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	_, err := r.OktaIDaaSClient.OktaSDKClientV5().IdentitySourceAPI.DeleteIdentitySourceSession(ctx, data.IdentitySourceId.ValueString(), data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Identity Source Session",
			fmt.Sprintf("Could not delete identity source session %s for identity source ID %s: %v", data.Id.ValueString(), data.IdentitySourceId.ValueString(), err),
		)
		return
	}
}

func (r *identitySourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, tfpath.Root("id"), req, resp)
}

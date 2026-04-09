package idaas

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	frameworkPath "github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	okta "github.com/okta/okta-sdk-golang/v6/okta"

	"github.com/okta/terraform-provider-okta/okta/config"
)

// Ensure interface compliance
var (
	_ resource.Resource                = &identitySourceGroupResource{}
	_ resource.ResourceWithConfigure   = &identitySourceGroupResource{}
	_ resource.ResourceWithImportState = &identitySourceGroupResource{}
)

// IdentitySourceGroupResource defines the resource implementation.
type identitySourceGroupResource struct {
	Config *config.Config
}

// IdentitySourceGroupModel describes the resource data model.
type identitySourceGroupModel struct {
	ID                 types.String `tfsdk:"id"`
	IdentitySourceId   types.String `tfsdk:"identity_source_id"`
	ExternalId         types.String `tfsdk:"external_id"`
	ProfileDescription types.String `tfsdk:"profile_description"`
	ProfileDisplayName types.String `tfsdk:"profile_display_name"`
}

func NewIdentitySourceGroupResource() resource.Resource {
	return &identitySourceGroupResource{}
}

func (r *identitySourceGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_identity_source_group"
}

func (r *identitySourceGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	cfg, ok := req.ProviderData.(*config.Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			"Expected *config.Config, got something else. Please report this issue to the provider developers.",
		)
		return
	}
	r.Config = cfg
}

func (r *identitySourceGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Okta IdentitySourceGroup resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier for the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"identity_source_id": schema.StringAttribute{
				Description: "ID of the identity source",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"external_id": schema.StringAttribute{
				Description: "The external ID of the identity source group",
				Optional:    true,
			},
			"profile_description": schema.StringAttribute{
				Description: "Description of the group profile",
				Optional:    true,
				Computed:    true,
			},
			"profile_display_name": schema.StringAttribute{
				Description: "Display name of the group profile",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *identitySourceGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, frameworkPath.Root("id"), req, resp)
}

func (r *identitySourceGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state identitySourceGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	identity_source_id := state.IdentitySourceId.ValueString()

	client := r.Config.OktaIDaaSClient.OktaSDKClientV6()
	result, httpResp, err := client.IdentitySourceAPI.GetIdentitySourceGroup(ctx, identity_source_id, id).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading identity_source_group", err.Error())
		return
	}
	// Map API response fields to state (scalar types only; WriteOnly fields are skipped — response type doesn't have them)
	state.ExternalId = types.StringValue(result.GetExternalId())
	state.ID = types.StringValue(result.GetId())

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *identitySourceGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan identitySourceGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	identitySourceId := plan.IdentitySourceId.ValueString()

	client := r.Config.OktaIDaaSClient.OktaSDKClientV6()

	// Build request body from plan
	createReq := client.IdentitySourceAPI.CreateIdentitySourceGroups(ctx, identitySourceId)
	body := okta.NewGroupsRequestSchemaWithDefaults()
	body.SetExternalId(plan.ExternalId.ValueString())
	createReq = createReq.GroupsRequestSchema(*body)

	result, _, err := createReq.Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error creating identity_source_group", err.Error())
		return
	}
	// Set ID from API response
	plan.ID = types.StringValue(result.GetId())

	// Map computed fields from response back to plan (scalar types only)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
func (r *identitySourceGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan identitySourceGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state identitySourceGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	identity_source_id := state.IdentitySourceId.ValueString()

	client := r.Config.OktaIDaaSClient.OktaSDKClientV6()

	// Build request body from plan — only send changed fields
	updateReq := client.IdentitySourceAPI.UpdateIdentitySourceGroups(ctx, identity_source_id, id)
	updateBody := okta.NewGroupsRequestSchemaWithDefaults()
	updateBody.SetExternalId(plan.ExternalId.ValueString())
	updateReq = updateReq.GroupsRequestSchema(*updateBody)

	res, _, err := updateReq.Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error updating identity_source_group", err.Error())
		return
	}
	// Map computed fields from response back to plan (scalar types only)

	plan.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &res)...)
}

func (r *identitySourceGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state identitySourceGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	identity_source_id := state.IdentitySourceId.ValueString()

	client := r.Config.OktaIDaaSClient.OktaSDKClientV6()
	httpResp, err := client.IdentitySourceAPI.DeleteIdentitySourceGroup(ctx, identity_source_id, id).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			return
		}
		resp.Diagnostics.AddError("Error deleting identity_source_group", err.Error())
		return
	}
}

// Ensure diag is used
var _ diag.Diagnostics

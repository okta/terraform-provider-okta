package idaas

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	frameworkPath "github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/okta/terraform-provider-okta/okta/config"
)

// Ensure interface compliance
var (
	_ resource.Resource                = &identitySourceSessionImportResource{}
	_ resource.ResourceWithConfigure   = &identitySourceSessionImportResource{}
	_ resource.ResourceWithImportState = &identitySourceSessionImportResource{}
)

// identitySourceSessionImportResource triggers the import of staged identity source data into Okta.
type identitySourceSessionImportResource struct {
	Config *config.Config
}

// identitySourceSessionImportModel describes the resource data model.
type identitySourceSessionImportModel struct {
	ID               types.String `tfsdk:"id"`
	IdentitySourceId types.String `tfsdk:"identity_source_id"`
	SessionId        types.String `tfsdk:"session_id"`
	Status           types.String `tfsdk:"status"`
	ImportType       types.String `tfsdk:"import_type"`
	Created          types.String `tfsdk:"created"`
	LastUpdated      types.String `tfsdk:"last_updated"`
}

func newIdentitySourceSessionImportResource() resource.Resource {
	return &identitySourceSessionImportResource{}
}

func (r *identitySourceSessionImportResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_identity_source_session_import"
}

func (r *identitySourceSessionImportResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *identitySourceSessionImportResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Triggers the import of staged identity source data into Okta. " +
			"This resource calls the startImportFromIdentitySource API, which applies all bulk-upsert and bulk-delete " +
			"operations previously uploaded to the given session. It must be applied after all bulk operation resources " +
			"for the same session have been created.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier for the resource (set to the session ID).",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"identity_source_id": schema.StringAttribute{
				Description: "ID of the identity source.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"session_id": schema.StringAttribute{
				Description: "ID of the identity source session to trigger import for.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				Description: "The current status of the identity source session after the import is triggered.",
				Computed:    true,
			},
			"import_type": schema.StringAttribute{
				Description: "The type of import. All imports are INCREMENTAL.",
				Computed:    true,
			},
			"created": schema.StringAttribute{
				Description: "Timestamp when the session was created (RFC3339).",
				Computed:    true,
			},
			"last_updated": schema.StringAttribute{
				Description: "Timestamp when the session was last updated (RFC3339).",
				Computed:    true,
			},
		},
	}
}

func (r *identitySourceSessionImportResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import ID format: {identity_source_id}/{session_id}
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Import ID must be in the format: {identity_source_id}/{session_id}",
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, frameworkPath.Root("identity_source_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, frameworkPath.Root("session_id"), parts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, frameworkPath.Root("id"), parts[1])...)
}

func (r *identitySourceSessionImportResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state identitySourceSessionImportModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Config.OktaIDaaSClient.OktaSDKClientV6()
	result, httpResp, err := client.IdentitySourceAPI.GetIdentitySourceSession(ctx, state.IdentitySourceId.ValueString(), state.SessionId.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading identity_source_session_import", err.Error())
		return
	}

	state.Status = types.StringValue(result.GetStatus())
	state.ImportType = types.StringValue(result.GetImportType())
	state.Created = types.StringValue(result.GetCreated().Format(time.RFC3339))
	state.LastUpdated = types.StringValue(result.GetLastUpdated().Format(time.RFC3339))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *identitySourceSessionImportResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan identitySourceSessionImportModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Config.OktaIDaaSClient.OktaSDKClientV6()
	result, _, err := client.IdentitySourceAPI.StartImportFromIdentitySource(ctx, plan.IdentitySourceId.ValueString(), plan.SessionId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error starting identity source import", err.Error())
		return
	}

	plan.ID = plan.SessionId
	plan.Status = types.StringValue(result.GetStatus())
	plan.ImportType = types.StringValue(result.GetImportType())
	plan.Created = types.StringValue(result.GetCreated().Format(time.RFC3339))
	plan.LastUpdated = types.StringValue(result.GetLastUpdated().Format(time.RFC3339))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *identitySourceSessionImportResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddWarning(
		"Update Not Supported",
		"This resource does not support in-place updates. Changes will require resource replacement.",
	)
}

func (r *identitySourceSessionImportResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning(
		"Delete Not Supported",
		"Removing this resource from configuration does not undo the import that was triggered in Okta.",
	)
}

// Ensure diag is used
var _ diag.Diagnostics

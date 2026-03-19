package idaas

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/okta-sdk-golang/v5/okta"
	"github.com/okta/terraform-provider-okta/okta/config"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &groupOwnerResource{}
	_ resource.ResourceWithConfigure   = &groupOwnerResource{}
	_ resource.ResourceWithImportState = &groupOwnerResource{}
)

func newGroupOwnerResource() resource.Resource {
	return &groupOwnerResource{}
}

type groupOwnerResource struct {
	*config.Config
}

type groupOwnerResourceModel struct {
	DisplayName    types.String `tfsdk:"display_name"`
	GroupID        types.String `tfsdk:"group_id"`
	IdOfGroupOwner types.String `tfsdk:"id_of_group_owner"`
	ID             types.String `tfsdk:"id"`
	OriginId       types.String `tfsdk:"origin_id"`
	OriginType     types.String `tfsdk:"origin_type"`
	Resolved       types.Bool   `tfsdk:"resolved"`
	Type           types.String `tfsdk:"type"`
}

func (r *groupOwnerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group_owner"
}

func (r *groupOwnerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `Manages group owner resource.`,
		Attributes: map[string]schema.Attribute{
			"display_name": schema.StringAttribute{
				Description: "The display name of the group owner",
				Computed:    true,
			},
			"group_id": schema.StringAttribute{
				Description: "The id of the group",
				Required:    true,
			},
			"id_of_group_owner": schema.StringAttribute{
				Description: "The user id of the group owner",
				Required:    true,
			},
			"id": schema.StringAttribute{
				Description: "The id of the group owner resource",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"origin_id": schema.StringAttribute{
				Description: "The ID of the app instance if the originType is APPLICATION. This value is NULL if originType is OKTA_DIRECTORY.",
				Computed:    true,
			},
			"origin_type": schema.StringAttribute{
				Description: "The source where group ownership is managed. Enum: \"APPLICATION\" \"OKTA_DIRECTORY\"",
				Computed:    true,
			},
			"resolved": schema.BoolAttribute{
				Description: "If originType is APPLICATION, this parameter is set to FALSE until the owner's originId is reconciled with an associated Okta ID.",
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "The entity type of the owner. Enum: \"GROUP\" \"USER\"",
				Required:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *groupOwnerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *groupOwnerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state groupOwnerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReqBody, err := buildCreateGroupOwnerRequest(state)
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to build group owner request",
			err.Error(),
		)
		return
	}

	createdGroupOwner, apiResp, err := r.OktaIDaaSClient.OktaSDKClientV5().GroupOwnerAPI.AssignGroupOwner(ctx, state.GroupID.ValueString()).AssignGroupOwnerRequestBody(createReqBody).Execute()
	if err != nil {
		// Handle the case where the owner is already assigned (HTTP 400 with specific error cause)
		if isAlreadyAssignedOwnerError(apiResp, err) {
			// Idempotent: owner already assigned. Set ID and let the next Read populate the rest.
			state.ID = types.StringValue(state.IdOfGroupOwner.ValueString())
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}

		resp.Diagnostics.AddError(
			fmt.Sprintf(
				"failed to add okta_group_owner with id of '%s' and type '%s' to okta_group '%s'",
				state.IdOfGroupOwner.ValueString(),
				state.Type.ValueString(),
				state.GroupID.ValueString(),
			),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(mapGroupOwnerToState(createdGroupOwner, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *groupOwnerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state groupOwnerResourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var grpOwner *okta.GroupOwner

	listGroupOwners, apiResp, err := r.OktaIDaaSClient.OktaSDKClientV5().GroupOwnerAPI.ListGroupOwners(ctx, state.GroupID.ValueString()).Execute()
	if err != nil {
		// If the group was deleted, Okta also deletes its owners automatically.
		// Treat 404 as resource gone and remove from state without error.
		if suppressErr := utils.SuppressErrorOn404_V5(apiResp, err); suppressErr == nil {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error retrieving the list of owners for okta_group '%s'", state.GroupID.ValueString()),
			fmt.Sprintf("Error returned: %s", err.Error()),
		)
		// Clear the state since we couldn't read the resource
		resp.State.RemoveResource(ctx)
		return
	}

	for _, groupOwner := range listGroupOwners {
		if groupOwner.GetId() == state.ID.ValueString() {
			grpOwner = &groupOwner
			break
		}
	}

	if grpOwner == nil {
		// The resource no longer exists; remove it from state
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(mapGroupOwnerToState(grpOwner, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *groupOwnerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state groupOwnerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := r.OktaIDaaSClient.OktaSDKClientV5().GroupOwnerAPI.DeleteGroupOwner(ctx, state.GroupID.ValueString(), state.ID.ValueString()).Execute()
	if err != nil {
		// remove 404 errors from err, if it's nil after, it means the err was a 404
		if err := utils.SuppressErrorOn404_V5(apiResp, err); err == nil {
			// If the group no longer exists, owners were already removed; treat as successful delete
			return
		}
		resp.Diagnostics.AddError(
			"failed to delete group owner "+state.ID.ValueString()+" from group",
			err.Error(),
		)
		return
	}
}

func (r *groupOwnerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state groupOwnerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReqBody, err := buildCreateGroupOwnerRequest(state)
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to build group owner request",
			err.Error(),
		)
		return
	}

	createdGroupOwner, apiResp, err := r.OktaIDaaSClient.OktaSDKClientV5().GroupOwnerAPI.AssignGroupOwner(ctx, state.GroupID.ValueString()).AssignGroupOwnerRequestBody(createReqBody).Execute()
	if err != nil {
		if isAlreadyAssignedOwnerError(apiResp, err) {
			// Delegate to Read to populate full state consistently
			r.Read(ctx, resource.ReadRequest{State: req.State}, &resource.ReadResponse{State: resp.State, Diagnostics: resp.Diagnostics})
			return
		}
		resp.Diagnostics.AddError(
			"failed to update/assign group owner",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(mapGroupOwnerToState(createdGroupOwner, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *groupOwnerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import ID format should be "group_id/group_owner_id"
	// Example: "group_123/group_owner_456"
	importID := req.ID
	if importID == "" {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Import ID cannot be empty. Expected format: group_id/group_owner_id",
		)
		return
	}

	// Split the import ID by forward slash
	parts := strings.Split(importID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import ID format",
			fmt.Sprintf("Expected import ID format 'group_id/group_owner_id', got '%s'", importID),
		)
		return
	}

	groupID := parts[0]
	groupOwnerID := parts[1]

	if groupID == "" || groupOwnerID == "" {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Both group_id and group_owner_id must be provided in import ID",
		)
		return
	}

	// Set both the group_id and id fields in the state
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("group_id"), groupID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), groupOwnerID)...)
}

func buildCreateGroupOwnerRequest(model groupOwnerResourceModel) (okta.AssignGroupOwnerRequestBody, error) {
	return okta.AssignGroupOwnerRequestBody{
		Id:   model.IdOfGroupOwner.ValueStringPointer(),
		Type: model.Type.ValueStringPointer(),
	}, nil
}

func mapGroupOwnerToState(data *okta.GroupOwner, state *groupOwnerResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	state.ID = types.StringPointerValue(data.Id)
	state.DisplayName = types.StringPointerValue(data.DisplayName)
	state.OriginId = types.StringPointerValue(data.OriginId)
	state.OriginType = types.StringPointerValue(data.OriginType)
	state.Resolved = types.BoolPointerValue(data.Resolved)
	state.Type = types.StringPointerValue(data.Type)

	return diags
}

// isAlreadyAssignedOwnerError returns true if the error response indicates that the provided owner is already assigned to the group.
// It checks the HTTP 400 status and looks for the full phrase: "Provided owner is already assigned to this group" in the error causes
// (case-insensitive). This avoids relying on generic error codes.
func isAlreadyAssignedOwnerError(apiResp *okta.APIResponse, err error) bool {
	if err == nil {
		return false
	}
	if apiResp == nil || apiResp.StatusCode != 400 {
		return false
	}
	needle := "provided owner is already assigned to this group"
	// The Okta V5 SDK wraps errors as GenericOpenAPIError where Model() may be an okta.Error
	var oae okta.GenericOpenAPIError
	if errors.As(err, &oae) {
		if m := oae.Model(); m != nil {
			if oe, ok := m.(okta.Error); ok {
				for _, cause := range oe.GetErrorCauses() {
					if strings.Contains(strings.ToLower(cause.GetErrorSummary()), needle) {
						return true
					}
				}
			}
		}
	}
	// Fallback to string search when we can't unpack a model
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, needle)
}

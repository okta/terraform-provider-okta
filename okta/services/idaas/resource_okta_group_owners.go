package idaas

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/okta-sdk-golang/v5/okta"
	"github.com/okta/terraform-provider-okta/okta/config"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &groupOwnersResource{}
	_ resource.ResourceWithConfigure   = &groupOwnersResource{}
	_ resource.ResourceWithImportState = &groupOwnersResource{}
)

func newGroupOwnersResource() resource.Resource {
	return &groupOwnersResource{}
}

type groupOwnersResource struct {
	*config.Config
}

type groupOwnersResourceModel struct {
	GroupID types.String `tfsdk:"group_id"`
	Owners  types.Set    `tfsdk:"owner"`
	ID      types.String `tfsdk:"id"`
}

type ownerEntryModel struct {
	ID   types.String `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
}

var ownerObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":   types.StringType,
		"type": types.StringType,
	},
}

var errGroupNotFound = errors.New("group not found")

func (r *groupOwnersResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group_owners"
}

func (r *groupOwnersResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage owners for a group in bulk. Uses the group_id as the resource ID. The resource is authoritative: any owners on the group not declared in configuration will be removed.",
		Attributes: map[string]schema.Attribute{
			"group_id": schema.StringAttribute{
				Description: "The ID of the Okta group.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Description: "Resource ID (same as group_id).",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"owner": schema.SetNestedBlock{
				Description: "Desired owners for the group.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The ID of the owner entity.",
							Required:    true,
						},
						"type": schema.StringAttribute{
							Description: "The entity type of the owner. Enum: \"GROUP\" \"USER\"",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.OneOf("GROUP", "USER"),
							},
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *groupOwnersResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *groupOwnersResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan groupOwnersResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupID := plan.GroupID.ValueString()

	// assign owners in the plan
	owners, diags := expandOwners(ctx, plan.Owners)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set ID early so partial failures still write state, allowing
	// the next plan/apply to reconcile via Read.
	plan.ID = types.StringValue(groupID)

	declaredMap := map[string]bool{}
	for _, o := range owners {
		declaredMap[o.ID] = true
		assignReq := okta.AssignGroupOwnerRequestBody{
			Id:   okta.PtrString(o.ID),
			Type: okta.PtrString(o.Type),
		}
		_, apiResp, err := r.OktaIDaaSClient.OktaSDKClientV5().GroupOwnerAPI.AssignGroupOwner(ctx, groupID).AssignGroupOwnerRequestBody(assignReq).Execute()
		if err != nil {
			if isAlreadyAssignedOwnerError(apiResp, err) {
				continue
			}
			resp.Diagnostics.AddError(
				fmt.Sprintf("failed to assign group owner (group_id=%s owner_id=%s type=%s)", groupID, o.ID, o.Type),
				assignOwnerErrorDetail(apiResp, err, o.ID, o.Type),
			)
			resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
			return
		}
	}

	// Authoritative: remove any pre-existing owners not in the declared set.
	// listAllGroupOwners confirms the group exists; a 404 on delete below
	// means the owner entity was removed out of band, not the group.
	apiOwners, err := r.listAllGroupOwners(ctx, groupID)
	if err != nil {
		resp.Diagnostics.AddError("failed to list owners after create", err.Error())
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		return
	}
	for _, existing := range apiOwners {
		if !declaredMap[existing.ID] {
			apiResp, err := r.OktaIDaaSClient.OktaSDKClientV5().GroupOwnerAPI.DeleteGroupOwner(ctx, groupID, existing.ID).Execute()
			if err != nil {
				if utils.SuppressErrorOn404_V5(apiResp, err) == nil {
					continue
				}
				resp.Diagnostics.AddError("failed to remove undeclared owner during create", fmt.Sprintf("group_id=%s owner_id=%s error=%s", groupID, existing.ID, err))
				resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
				return
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *groupOwnersResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state groupOwnersResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupID := state.GroupID.ValueString()

	apiOwners, err := r.listAllGroupOwners(ctx, groupID)
	if err != nil {
		if errors.Is(err, errGroupNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("error retrieving group owners", err.Error())
		return
	}

	setVal, diags := flattenOwners(apiOwners)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Owners = setVal

	state.ID = types.StringValue(groupID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *groupOwnersResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan groupOwnersResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupID := plan.GroupID.ValueString()

	newOwners, diags := expandOwners(ctx, plan.Owners)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Diff against actual API state for authoritative behavior
	apiOwners, err := r.listAllGroupOwners(ctx, groupID)
	if err != nil {
		if errors.Is(err, errGroupNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to list current owners", err.Error())
		return
	}

	oldMap := map[string]string{}
	for _, o := range apiOwners {
		oldMap[o.ID] = o.Type
	}
	newMap := map[string]string{}
	for _, o := range newOwners {
		newMap[o.ID] = o.Type
	}

	// additions (new owners not previously assigned)
	for id, typ := range newMap {
		if _, ok := oldMap[id]; !ok {
			assignReq := okta.AssignGroupOwnerRequestBody{Id: okta.PtrString(id), Type: okta.PtrString(typ)}
			_, apiResp, err := r.OktaIDaaSClient.OktaSDKClientV5().GroupOwnerAPI.AssignGroupOwner(ctx, groupID).AssignGroupOwnerRequestBody(assignReq).Execute()
			if err != nil {
				if isAlreadyAssignedOwnerError(apiResp, err) {
					continue
				}
				resp.Diagnostics.AddError(
					fmt.Sprintf("failed to assign group owner (group_id=%s owner_id=%s type=%s)", groupID, id, typ),
					assignOwnerErrorDetail(apiResp, err, id, typ),
				)
				plan.ID = types.StringValue(groupID)
				resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
				return
			}
		}
	}

	// removals — listAllGroupOwners above confirmed the group exists, so
	// a 404 here means the owner entity was deleted out of band, not the group.
	for id := range oldMap {
		if _, ok := newMap[id]; !ok {
			apiResp, err := r.OktaIDaaSClient.OktaSDKClientV5().GroupOwnerAPI.DeleteGroupOwner(ctx, groupID, id).Execute()
			if err != nil {
				if utils.SuppressErrorOn404_V5(apiResp, err) == nil {
					continue
				}
				resp.Diagnostics.AddError("failed to delete group owner", fmt.Sprintf("group_id=%s owner_id=%s error=%s", groupID, id, err))
				plan.ID = types.StringValue(groupID)
				resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
				return
			}
		}
	}

	plan.ID = types.StringValue(groupID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *groupOwnersResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state groupOwnersResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupID := state.GroupID.ValueString()

	// Query live API to remove all owners, not just state-tracked ones
	apiOwners, err := r.listAllGroupOwners(ctx, groupID)
	if err != nil {
		if errors.Is(err, errGroupNotFound) {
			return // group gone, owners gone
		}
		resp.Diagnostics.AddError("failed to list owners for deletion", err.Error())
		return
	}

	for _, o := range apiOwners {
		apiResp, err := r.OktaIDaaSClient.OktaSDKClientV5().GroupOwnerAPI.DeleteGroupOwner(ctx, groupID, o.ID).Execute()
		if err != nil {
			if utils.SuppressErrorOn404_V5(apiResp, err) == nil {
				continue
			}
			resp.Diagnostics.AddError(
				"failed to delete group owner",
				fmt.Sprintf("group_id=%s owner_id=%s error=%s", groupID, o.ID, err),
			)
		}
	}
}

func (r *groupOwnersResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	groupID := req.ID
	if groupID == "" {
		resp.Diagnostics.AddError("invalid import id", "group_id must not be empty")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), groupID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("group_id"), groupID)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Hydrate owner blocks from the API so the next plan does not error on
	// missing required attributes. Without this, the framework would leave
	// the "owner" set null and fail on the subsequent plan.
	apiOwners, err := r.listAllGroupOwners(ctx, groupID)
	if err != nil {
		if errors.Is(err, errGroupNotFound) {
			resp.Diagnostics.AddError("group not found", fmt.Sprintf("group %q does not exist", groupID))
			return
		}
		resp.Diagnostics.AddError("failed to list owners during import", err.Error())
		return
	}

	setVal, diags := flattenOwners(apiOwners)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("owner"), setVal)...)
}

// Internal helpers

type ownerFlat struct {
	ID   string
	Type string
}

func (r *groupOwnersResource) listAllGroupOwners(ctx context.Context, groupID string) ([]ownerFlat, error) {
	owners, apiResp, err := r.OktaIDaaSClient.OktaSDKClientV5().GroupOwnerAPI.ListGroupOwners(ctx, groupID).Execute()
	if err != nil {
		if utils.SuppressErrorOn404_V5(apiResp, err) == nil {
			return nil, errGroupNotFound
		}
		return nil, err
	}
	result := make([]ownerFlat, 0, len(owners))
	for _, o := range owners {
		id := safeString(o.Id)
		typ := safeString(o.Type)
		if id == "" || typ == "" {
			continue
		}
		result = append(result, ownerFlat{ID: id, Type: strings.ToUpper(typ)})
	}
	// attempt to page if available
	for apiResp != nil && apiResp.HasNextPage() {
		var next []okta.GroupOwner
		var err2 error
		apiResp, err2 = apiResp.Next(&next)
		if err2 != nil {
			return result, fmt.Errorf("error paginating group owners: %w", err2)
		}
		for _, o := range next {
			id := safeString(o.Id)
			typ := safeString(o.Type)
			if id == "" || typ == "" {
				continue
			}
			result = append(result, ownerFlat{ID: id, Type: strings.ToUpper(typ)})
		}
	}
	return result, nil
}

func expandOwners(ctx context.Context, set types.Set) ([]ownerFlat, diag.Diagnostics) {
	if set.IsNull() || set.IsUnknown() {
		return []ownerFlat{}, nil
	}
	var entries []ownerEntryModel
	diags := set.ElementsAs(ctx, &entries, false)
	if diags.HasError() {
		return nil, diags
	}
	result := make([]ownerFlat, 0, len(entries))
	seen := map[string]bool{}
	for _, e := range entries {
		if e.ID.IsNull() || e.ID.IsUnknown() || e.Type.IsNull() || e.Type.IsUnknown() {
			continue
		}
		id := e.ID.ValueString()
		if seen[id] {
			diags.AddError("duplicate owner id", fmt.Sprintf("owner id %q appears more than once", id))
			return nil, diags
		}
		seen[id] = true
		typ := strings.ToUpper(e.Type.ValueString())
		warnIfIDTypeMismatch(&diags, id, typ)
		result = append(result, ownerFlat{ID: id, Type: typ})
	}
	return result, nil
}

func flattenOwners(in []ownerFlat) (types.Set, diag.Diagnostics) {
	elements := make([]attr.Value, 0, len(in))
	for _, o := range in {
		obj := types.ObjectValueMust(ownerObjectType.AttrTypes, map[string]attr.Value{
			"id":   types.StringValue(o.ID),
			"type": types.StringValue(strings.ToUpper(o.Type)),
		})
		elements = append(elements, obj)
	}
	set, d := types.SetValue(ownerObjectType, elements)
	return set, d
}

package idaas

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
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
	GroupID        types.String `tfsdk:"group_id"`
	TrackAllOwners types.Bool   `tfsdk:"track_all_owners"`
	Owners         types.Set    `tfsdk:"owner"`
	ID             types.String `tfsdk:"id"`
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
		Description: "Manage owners for a group in bulk. Uses the group_id as the resource ID. By default, the resource manages all owners (removing any not specified). Set track_all_owners = false to opt-out of removing owners that aren't tracked in state.",
		Attributes: map[string]schema.Attribute{
			"group_id": schema.StringAttribute{
				Description: "The ID of the Okta group.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"track_all_owners": schema.BoolAttribute{
				Description: "If true (default), the resource tracks all owners on the group and will remove owners not declared in configuration. Set to false to opt-out of removing owners that aren't tracked in state.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
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
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
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

// isAlreadyAssignedOwnerError returns true if AssignGroupOwner returned the specific
// 400 error indicating the owner is already assigned to the group.
func isAlreadyAssignedOwnerError(apiResp *okta.APIResponse, err error) bool {
	if err == nil {
		return false
	}
	if apiResp == nil || apiResp.StatusCode != 400 {
		return false
	}
	needle := "provided owner is already assigned to this group"
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
	// Fallback: substring search in error
	return strings.Contains(strings.ToLower(err.Error()), needle)
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

	for _, o := range owners {
		assignReq := okta.AssignGroupOwnerRequestBody{
			Id:   okta.PtrString(o.ID),
			Type: okta.PtrString(o.Type),
		}
		_, apiResp, err := r.OktaIDaaSClient.OktaSDKClientV5().GroupOwnerAPI.AssignGroupOwner(ctx, groupID).AssignGroupOwnerRequestBody(assignReq).Execute()
		if err != nil {
			if isAlreadyAssignedOwnerError(apiResp, err) {
				continue
			}
			resp.Diagnostics.AddError("failed to assign group owner", fmt.Sprintf("group_id=%s owner_id=%s type=%s error=%s", groupID, o.ID, o.Type, err))
			return
		}
	}

	plan.ID = types.StringValue(groupID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *groupOwnersResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state groupOwnersResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupID := state.GroupID.ValueString()
	trackAll := true
	if !state.TrackAllOwners.IsNull() && !state.TrackAllOwners.IsUnknown() {
		trackAll = state.TrackAllOwners.ValueBool()
	}

	apiOwners, err := r.listAllGroupOwners(ctx, groupID)
	if err != nil {
		if errors.Is(err, errGroupNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("error retrieving group owners", err.Error())
		return
	}

	if trackAll {
		// set state owners to all from API
		setVal, diags := flattenOwners(ctx, apiOwners)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.Owners = setVal
	} else {
		// only remove entries that no longer exist in API, preserve others
		currentStateOwners, diags := expandOwners(ctx, state.Owners)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		apiMap := map[string]string{}
		for _, o := range apiOwners {
			apiMap[o.ID] = o.Type
		}
		filtered := []ownerFlat{}
		for _, o := range currentStateOwners {
			if _, ok := apiMap[o.ID]; ok {
				filtered = append(filtered, o)
			}
		}
		setVal, diags := flattenOwners(ctx, filtered)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.Owners = setVal
	}

	state.ID = types.StringValue(groupID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *groupOwnersResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state groupOwnersResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...) // plan contains desired owners and possibly track_all_owners change
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupID := plan.GroupID.ValueString()

	newOwners, diags := expandOwners(ctx, plan.Owners)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	oldOwners, diags := expandOwners(ctx, state.Owners)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	oldMap := map[string]string{}
	for _, o := range oldOwners {
		oldMap[o.ID] = o.Type
	}
	newMap := map[string]string{}
	for _, o := range newOwners {
		newMap[o.ID] = o.Type
	}

	// additions
	for id, typ := range newMap {
		if _, ok := oldMap[id]; !ok {
			assignReq := okta.AssignGroupOwnerRequestBody{Id: okta.PtrString(id), Type: okta.PtrString(typ)}
			_, apiResp, err := r.OktaIDaaSClient.OktaSDKClientV5().GroupOwnerAPI.AssignGroupOwner(ctx, groupID).AssignGroupOwnerRequestBody(assignReq).Execute()
			if err != nil {
				if isAlreadyAssignedOwnerError(apiResp, err) {
					continue
				}
				resp.Diagnostics.AddError("failed to assign group owner", fmt.Sprintf("group_id=%s owner_id=%s type=%s error=%s", groupID, id, typ, err))
				return
			}
		}
	}

	// removals
	for id := range oldMap {
		if _, ok := newMap[id]; !ok {
			apiResp, err := r.OktaIDaaSClient.OktaSDKClientV5().GroupOwnerAPI.DeleteGroupOwner(ctx, groupID, id).Execute()
			if err != nil {
				if utils.SuppressErrorOn404_V5(apiResp, err) == nil {
					continue
				}
				resp.Diagnostics.AddError("failed to delete group owner", fmt.Sprintf("group_id=%s owner_id=%s error=%s", groupID, id, err))
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
	owners, diags := expandOwners(ctx, state.Owners)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, o := range owners {
		apiResp, err := r.OktaIDaaSClient.OktaSDKClientV5().GroupOwnerAPI.DeleteGroupOwner(ctx, groupID, o.ID).Execute()
		if err != nil {
			if utils.SuppressErrorOn404_V5(apiResp, err) == nil {
				continue
			}
			resp.Diagnostics.AddError("failed to delete group owner", fmt.Sprintf("group_id=%s owner_id=%s error=%s", groupID, o.ID, err))
			return
		}
	}
}

func (r *groupOwnersResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// support formats: group_id or group_id/true|false
	parts := strings.Split(req.ID, "/")
	if len(parts) > 2 {
		resp.Diagnostics.AddError("invalid import id", "format must be 'group_id' or 'group_id/true|false'")
		return
	}
	groupID := parts[0]
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), groupID)...)       // set id
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("group_id"), groupID)...) // set group_id
	if len(parts) == 2 {
		track := strings.EqualFold(parts[1], "true")
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("track_all_owners"), track)...) // set flag
	}
}

// Internal helpers

type ownerFlat struct {
	ID   string
	Type string
}

func (r *groupOwnersResource) listAllGroupOwners(ctx context.Context, groupID string) ([]ownerFlat, error) {
	owners, apiResp, err := r.OktaIDaaSClient.OktaSDKClientV5().GroupOwnerAPI.ListGroupOwners(ctx, groupID).Execute()
	if err != nil {
		if apiResp != nil && apiResp.StatusCode == 404 {
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
		result = append(result, ownerFlat{ID: id, Type: typ})
	}
	// attempt to page if available
	for apiResp != nil && apiResp.HasNextPage() {
		var next []okta.GroupOwner
		var err2 error
		apiResp, err2 = apiResp.Next(&next)
		if err2 != nil {
			return result, nil // return what we have
		}
		for _, o := range next {
			id := safeString(o.Id)
			typ := safeString(o.Type)
			if id == "" || typ == "" {
				continue
			}
			result = append(result, ownerFlat{ID: id, Type: typ})
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
	for _, e := range entries {
		if e.ID.IsNull() || e.ID.IsUnknown() || e.Type.IsNull() || e.Type.IsUnknown() {
			continue
		}
		result = append(result, ownerFlat{ID: e.ID.ValueString(), Type: strings.ToUpper(e.Type.ValueString())})
	}
	return result, nil
}

func flattenOwners(ctx context.Context, in []ownerFlat) (types.Set, diag.Diagnostics) {
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

func safeString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

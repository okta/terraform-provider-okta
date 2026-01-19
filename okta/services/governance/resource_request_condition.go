package governance

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/okta-governance-sdk-golang/governance"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ resource.Resource                = &requestConditionResource{}
	_ resource.ResourceWithConfigure   = &requestConditionResource{}
	_ resource.ResourceWithImportState = &requestConditionResource{}
)

func newRequestConditionResource() resource.Resource {
	return &requestConditionResource{}
}

func (r *requestConditionResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	parts := strings.Split(request.ID, "/")
	if len(parts) != 2 {
		response.Diagnostics.AddError(
			"Invalid import ID",
			"Expected format: resource_id/sequence_id",
		)
		return
	}

	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("resource_id"), parts[0])...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}

func (r *requestConditionResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(request, response)
}

type requestConditionResource struct {
	*config.Config
}

type IdModel struct {
	Id types.String `tfsdk:"id"`
}
type Settings struct {
	Type types.String `tfsdk:"type"`
	Ids  []IdModel    `tfsdk:"ids"`
}

type AccessDurationSettings struct {
	Type     types.String `tfsdk:"type"`
	Duration types.String `tfsdk:"duration"`
}

type requestConditionResourceModel struct {
	Id                     types.String            `tfsdk:"id"`
	ResourceId             types.String            `tfsdk:"resource_id"`
	ApprovalSequenceId     types.String            `tfsdk:"approval_sequence_id"`
	Name                   types.String            `tfsdk:"name"`
	Description            types.String            `tfsdk:"description"`
	Priority               types.Int32             `tfsdk:"priority"`
	Created                types.String            `tfsdk:"created"`
	CreatedBy              types.String            `tfsdk:"created_by"`
	LastUpdated            types.String            `tfsdk:"last_updated"`
	LastUpdatedBy          types.String            `tfsdk:"last_updated_by"`
	Status                 types.String            `tfsdk:"status"`
	AccessScopeSettings    *Settings               `tfsdk:"access_scope_settings"`
	RequesterSettings      *Settings               `tfsdk:"requester_settings"`
	AccessDurationSettings *AccessDurationSettings `tfsdk:"access_duration_settings"`
}

func (r *requestConditionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_request_condition"
}

func (r *requestConditionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The id of the request condition.",
			},
			"resource_id": schema.StringAttribute{
				Required:    true,
				Description: "The id of the resource in Okta ID format.",
			},
			"approval_sequence_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the approval sequence.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the request condition.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the request condition.",
			},
			"priority": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The priority of the condition. The smaller the number, the higher the priority.",
			},
			"created":         schema.StringAttribute{Computed: true, Description: "The ISO 8601 formatted date and time when the resource was created."},
			"created_by":      schema.StringAttribute{Computed: true, Description: "The id of the Okta user who created the resource."},
			"last_updated":    schema.StringAttribute{Computed: true, Description: "The ISO 8601 formatted date and time when the object was last updated."},
			"last_updated_by": schema.StringAttribute{Computed: true, Description: "The id of the Okta user who last updated the object."},
			"status": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Status of the condition. Valid values: ACTIVE, INACTIVE. Default is INACTIVE.",
				Validators: []validator.String{
					stringvalidator.OneOf("ACTIVE", "INACTIVE"),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"access_scope_settings": schema.SingleNestedBlock{
				Blocks: map[string]schema.Block{
					"ids": schema.ListNestedBlock{
						Description: "Block list of groups/entitlement bundles ids.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Optional:    true,
									Description: "The group/entitlement bundle ID.",
								},
							},
						},
					},
				},
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Required: true,
					},
				},
			},
			"requester_settings": schema.SingleNestedBlock{
				Blocks: map[string]schema.Block{
					"ids": schema.ListNestedBlock{
						Description: "Block list of teams/groups ids.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Optional:    true,
									Description: "The group/team ID.",
								},
							},
						},
					},
				},
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Required: true,
					},
				},
			},
			"access_duration_settings": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("ADMIN_FIXED_DURATION",
								"REQUESTER_SPECIFIED_DURATION"),
						},
					},
					"duration": schema.StringAttribute{
						Optional:    true,
						Description: "The duration set by the admin for access durations. Use ISO8061 notation for duration values.",
					},
				},
				Description: "Settings that control who may specify the access duration allowed by this request condition or risk settings, as well as what duration may be requested.",
			},
		},
	}
}

func (r *requestConditionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data requestConditionResourceModel

	// Read Terraform plan Data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	requestConditionResp, _, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().RequestConditionsAPI.CreateResourceRequestConditionV2(ctx, data.ResourceId.ValueString()).RequestConditionCreatable(createRequestCondition(data)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Request conditions",
			"Could not create Request conditions, unexpected error: "+err.Error(),
		)
		return
	}

	// Activate the condition if status is set to ACTIVE
	if !data.Status.IsNull() && data.Status.ValueString() == "ACTIVE" {
		requestConditionResp, _, err = r.OktaGovernanceClient.OktaGovernanceSDKClient().
			RequestConditionsAPI.ActivateResourceRequestConditionV2(ctx,
			data.ResourceId.ValueString(),
			requestConditionResp.GetId()).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error activating Request condition",
				"Could not activate Request condition after creation: "+err.Error(),
			)
			return
		}
	}

	resp.Diagnostics.Append(applyRequestConditionToState(ctx, &data, requestConditionResp)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *requestConditionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data requestConditionResourceModel

	// Read Terraform prior state Data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	readRequestConditionResp, _, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().RequestConditionsAPI.GetResourceRequestConditionV2(ctx, data.ResourceId.ValueString(), data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Request conditions",
			"Could not read Request conditions, unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(applyRequestConditionToState(ctx, &data, readRequestConditionResp)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *requestConditionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state requestConditionResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform plan Data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic
	updatedRequestCondition, _, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().RequestConditionsAPI.UpdateResourceRequestConditionV2(ctx, data.ResourceId.ValueString(), state.Id.ValueString()).RequestConditionPatchable(createRequestConditionPatch(data)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Request conditions",
			"Could not update Request conditions, unexpected error: "+err.Error(),
		)
		return
	}

	// Handle status changes
	if !data.Status.IsNull() && !state.Status.IsNull() {
		oldStatus := state.Status.ValueString()
		newStatus := data.Status.ValueString()

		if oldStatus != newStatus {
			if newStatus == "ACTIVE" {
				updatedRequestCondition, _, err = r.OktaGovernanceClient.OktaGovernanceSDKClient().
					RequestConditionsAPI.ActivateResourceRequestConditionV2(ctx,
					data.ResourceId.ValueString(),
					state.Id.ValueString()).Execute()
				if err != nil {
					resp.Diagnostics.AddError(
						"Error activating Request condition",
						"Could not activate Request condition: "+err.Error(),
					)
					return
				}
			} else if newStatus == "INACTIVE" {
				updatedRequestCondition, _, err = r.OktaGovernanceClient.OktaGovernanceSDKClient().
					RequestConditionsAPI.DeactivateResourceRequestConditionV2(ctx,
					data.ResourceId.ValueString(),
					state.Id.ValueString()).Execute()
				if err != nil {
					resp.Diagnostics.AddError(
						"Error deactivating Request condition",
						"Could not deactivate Request condition: "+err.Error(),
					)
					return
				}
			}
		}
	}

	resp.Diagnostics.Append(applyRequestConditionToState(ctx, &data, updatedRequestCondition)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *requestConditionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data, state requestConditionResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform prior state Data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If the request condition is active, attempt to deactivate it first.
	// Request conditions must be INACTIVE before they can be deleted.
	if !state.Status.IsNull() && state.Status.ValueString() == "ACTIVE" {
		_, _, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().
			RequestConditionsAPI.DeactivateResourceRequestConditionV2(ctx,
			state.ResourceId.ValueString(),
			state.Id.ValueString()).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error deactivating Request condition before deletion",
				"Could not deactivate Request condition: "+err.Error(),
			)
			return
		}
	}

	// Delete API call logic
	_, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().RequestConditionsAPI.DeleteResourceRequestConditionV2(ctx, data.ResourceId.ValueString(), state.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Request conditions",
			"Could not delete Request conditions, unexpected error: "+err.Error(),
		)
		return
	}
}

func applyRequestConditionToState(ctx context.Context, data *requestConditionResourceModel, requestConditionResp *governance.RequestConditionFull) diag.Diagnostics {
	var diags diag.Diagnostics
	data.Id = types.StringValue(requestConditionResp.GetId())
	data.Name = types.StringValue(requestConditionResp.GetName())
	if requestConditionResp.Description != nil {
		data.Description = types.StringValue(requestConditionResp.GetDescription())
	}
	data.Priority = types.Int32Value(requestConditionResp.GetPriority())
	data.ApprovalSequenceId = types.StringValue(requestConditionResp.GetApprovalSequenceId())
	data.Created = types.StringValue(requestConditionResp.GetCreated().Format(time.RFC3339))
	data.CreatedBy = types.StringValue(requestConditionResp.GetCreatedBy())
	data.LastUpdated = types.StringValue(requestConditionResp.GetLastUpdated().Format(time.RFC3339))
	data.LastUpdatedBy = types.StringValue(requestConditionResp.GetLastUpdatedBy())
	data.Status = types.StringValue(string(requestConditionResp.GetStatus()))
	data.RequesterSettings, _ = setRequesterSettings(requestConditionResp.GetRequesterSettings())
	data.AccessScopeSettings, _ = setAccessScopeSettings(requestConditionResp.GetAccessScopeSettings())
	data.AccessDurationSettings = setAccessDurationSettings(requestConditionResp.GetAccessDurationSettings())
	return diags
}

func setAccessDurationSettings(settings governance.AccessDurationSettingsFull) *AccessDurationSettings {
	var accessDurationSettings AccessDurationSettings
	if settings.AccessDurationSettingsAdminFixedDuration != nil {
		accessDurationSettings.Type = types.StringValue(settings.AccessDurationSettingsAdminFixedDuration.Type)
		accessDurationSettings.Duration = types.StringValue(settings.AccessDurationSettingsAdminFixedDuration.Duration)
		return &accessDurationSettings
	} else if settings.AccessDurationSettingsRequesterSpecifiedDuration != nil {
		accessDurationSettings.Type = types.StringValue(settings.AccessDurationSettingsRequesterSpecifiedDuration.Type)
		accessDurationSettings.Duration = types.StringValue(settings.AccessDurationSettingsRequesterSpecifiedDuration.MaximumDuration)
		return &accessDurationSettings
	}
	return nil
}

func setAccessScopeSettings(settings governance.AccessScopeSettingsFullAccessScopeSettings) (*Settings, diag.Diagnostics) {
	var setting Settings
	if settings.GroupAccessScopeSettings != nil {
		setting.Type = types.StringValue(settings.GroupAccessScopeSettings.GetType())
		groupIds := settings.GroupAccessScopeSettings.GetGroups()
		var ids []IdModel
		for _, groupId := range groupIds {
			idModel := IdModel{
				Id: types.StringValue(groupId.GetId()),
			}
			ids = append(ids, idModel)
		}
		setting.Ids = ids
	} else if settings.EntitlementBundleAccessScopeSettings != nil {
		setting.Type = types.StringValue(settings.EntitlementBundleAccessScopeSettings.GetType())
		entitlementBundleIds := settings.EntitlementBundleAccessScopeSettings.GetEntitlementBundles()
		var ids []IdModel
		for _, entitlementBundleId := range entitlementBundleIds {
			idModel := IdModel{
				Id: types.StringValue(entitlementBundleId.GetId()),
			}
			ids = append(ids, idModel)
		}
		setting.Ids = ids
	} else if settings.ResourceDefaultAccessScopeSettings != nil {
		setting.Type = types.StringValue(settings.ResourceDefaultAccessScopeSettings.GetType())
	}
	return &setting, nil
}

func setRequesterSettings(settings governance.RequesterSettingsFullRequesterSettings) (*Settings, diag.Diagnostics) {
	var setting Settings
	if settings.GroupsRequesterSettings != nil {
		setting.Type = types.StringValue(settings.GroupsRequesterSettings.GetType())
		groupIds := settings.GroupsRequesterSettings.GetGroups()
		var ids []IdModel
		for _, groupId := range groupIds {
			idModel := IdModel{
				Id: types.StringValue(groupId.GetId()),
			}
			ids = append(ids, idModel)
		}
		setting.Ids = ids
	} else if settings.EveryoneRequesterSettings != nil {
		setting.Type = types.StringValue(settings.EveryoneRequesterSettings.GetType())
	}
	return &setting, nil
}

func createRequestCondition(data requestConditionResourceModel) governance.RequestConditionCreatable {
	req := governance.RequestConditionCreatable{}
	req.Name = data.Name.ValueString()
	req.ApprovalSequenceId = data.ApprovalSequenceId.ValueString()
	if !data.Description.IsNull() && data.Description.ValueString() != "" {
		req.Description = data.Description.ValueStringPointer()
	}
	if !data.Priority.IsNull() {
		req.Priority = data.Priority.ValueInt32Pointer()
	}

	accessScopeSettings := governance.AccessScopeSettingsCreatableAccessScopeSettings{}
	if data.AccessScopeSettings.Type.ValueString() == "GROUPS" {
		if accessScopeSettings.AccessScopeSettingsCreatableGroupAccessScopeSettings == nil {
			accessScopeSettings.AccessScopeSettingsCreatableGroupAccessScopeSettings =
				&governance.AccessScopeSettingsCreatableGroupAccessScopeSettings{}
		}
		accessScopeSettings.AccessScopeSettingsCreatableGroupAccessScopeSettings.Type = "GROUPS"
		var groupsIds []governance.GroupsArrayCreatableInner
		elems := data.AccessScopeSettings.Ids
		for _, elem := range elems {
			groupId := elem.Id
			groupsIds = append(groupsIds, governance.GroupsArrayCreatableInner{
				Id: groupId.ValueString(),
			})
		}
		accessScopeSettings.AccessScopeSettingsCreatableGroupAccessScopeSettings.Groups = groupsIds
		req.AccessScopeSettings = accessScopeSettings
	} else if data.AccessScopeSettings.Type.ValueString() == "ENTITLEMENT_BUNDLES" {
		if accessScopeSettings.AccessScopeSettingsCreatableEntitlementBundleAccessScopeSettings == nil {
			accessScopeSettings.AccessScopeSettingsCreatableEntitlementBundleAccessScopeSettings =
				&governance.AccessScopeSettingsCreatableEntitlementBundleAccessScopeSettings{}
		}
		accessScopeSettings.AccessScopeSettingsCreatableEntitlementBundleAccessScopeSettings.Type = "ENTITLEMENT_BUNDLES"
		var entitlementBundles []governance.EntitlementBundlesArrayCreatableInner
		elems := data.AccessScopeSettings.Ids
		for _, elem := range elems {
			groupId := elem.Id
			entitlementBundles = append(entitlementBundles, governance.EntitlementBundlesArrayCreatableInner{
				Id: groupId.ValueString(),
			})
		}
		accessScopeSettings.AccessScopeSettingsCreatableEntitlementBundleAccessScopeSettings.EntitlementBundles = entitlementBundles
		req.AccessScopeSettings = accessScopeSettings
	} else if data.AccessScopeSettings.Type.ValueString() == "RESOURCE_DEFAULT" {
		if accessScopeSettings.AccessScopeSettingsCreatableResourceDefaultAccessScopeSettings == nil {
			accessScopeSettings.AccessScopeSettingsCreatableResourceDefaultAccessScopeSettings = &governance.AccessScopeSettingsCreatableResourceDefaultAccessScopeSettings{}
		}
		accessScopeSettings.AccessScopeSettingsCreatableResourceDefaultAccessScopeSettings.Type = "RESOURCE_DEFAULT"
		req.AccessScopeSettings = accessScopeSettings
	}

	requestSettings := governance.RequesterSettingsCreatableRequesterSettings{}
	if data.RequesterSettings.Type.ValueString() == "GROUPS" {
		if requestSettings.RequesterSettingsCreatableGroupsRequesterSettings == nil {
			requestSettings.RequesterSettingsCreatableGroupsRequesterSettings = &governance.RequesterSettingsCreatableGroupsRequesterSettings{}
		}
		requestSettings.RequesterSettingsCreatableGroupsRequesterSettings.Type = "GROUPS"
		var groups []governance.GroupsArrayCreatableInner
		elems := data.RequesterSettings.Ids
		for _, elem := range elems {
			groupId := elem.Id
			groups = append(groups, governance.GroupsArrayCreatableInner{
				Id: groupId.ValueString(),
			})
		}
		requestSettings.RequesterSettingsCreatableGroupsRequesterSettings.Groups = groups
		req.RequesterSettings = requestSettings
	} else if data.RequesterSettings.Type.ValueString() == "EVERYONE" {
		if requestSettings.EveryoneRequesterSettings == nil {
			requestSettings.EveryoneRequesterSettings = &governance.EveryoneRequesterSettings{}
		}
		requestSettings.EveryoneRequesterSettings.Type = "EVERYONE"
		req.RequesterSettings = requestSettings
	}

	accessDurationSettings := &governance.AccessDurationSettingsCreatable{}
	if data.AccessDurationSettings != nil {
		if data.AccessDurationSettings.Type.ValueString() == "ADMIN_FIXED_DURATION" {
			durationSettings := &governance.AccessDurationSettingsAdminFixedDuration{}
			durationSettings.Duration = data.AccessDurationSettings.Duration.ValueString()
			durationSettings.Type = "ADMIN_FIXED_DURATION"
			accessDurationSettings.AccessDurationSettingsAdminFixedDuration = durationSettings
			req.AccessDurationSettings = accessDurationSettings
		} else if data.AccessDurationSettings.Type.ValueString() == "REQUESTER_SPECIFIED_DURATION" {
			durationSettings := &governance.AccessDurationSettingsRequesterSpecifiedDuration{}
			durationSettings.MaximumDuration = data.AccessDurationSettings.Duration.ValueString()
			durationSettings.Type = "REQUESTER_SPECIFIED_DURATION"
			accessDurationSettings.AccessDurationSettingsRequesterSpecifiedDuration = durationSettings
			req.AccessDurationSettings = accessDurationSettings
		}
	}

	return req
}

func createRequestConditionPatch(data requestConditionResourceModel) governance.RequestConditionPatchable {
	var patch governance.RequestConditionPatchable
	patch.Name = data.Name.ValueStringPointer()
	patch.ApprovalSequenceId = data.ApprovalSequenceId.ValueStringPointer()
	if !data.Description.IsNull() && data.Description.ValueString() != "" {
		patch.Description = data.Description.ValueStringPointer()
	}
	if !data.Priority.IsNull() {
		patch.Priority = data.Priority.ValueInt32Pointer()
	}
	patch.ApprovalSequenceId = data.ApprovalSequenceId.ValueStringPointer()
	var accessScopeSettings governance.AccessScopeSettingsCreatableAccessScopeSettings
	if data.AccessScopeSettings.Type.ValueString() == "GROUPS" {
		if accessScopeSettings.AccessScopeSettingsCreatableGroupAccessScopeSettings == nil {
			accessScopeSettings.AccessScopeSettingsCreatableGroupAccessScopeSettings = &governance.AccessScopeSettingsCreatableGroupAccessScopeSettings{}
		}
		accessScopeSettings.AccessScopeSettingsCreatableGroupAccessScopeSettings.Type = "GROUPS"
		var groupsIds []governance.GroupsArrayCreatableInner
		elems := data.AccessScopeSettings.Ids
		for _, elem := range elems {
			groupId := elem.Id.ValueString()
			groupsIds = append(groupsIds, governance.GroupsArrayCreatableInner{
				Id: groupId,
			})
		}
		accessScopeSettings.AccessScopeSettingsCreatableGroupAccessScopeSettings.Groups = groupsIds
		patch.AccessScopeSettings = &accessScopeSettings
	} else if data.AccessScopeSettings.Type.ValueString() == "ENTITLEMENT_BUNDLES" {
		if accessScopeSettings.AccessScopeSettingsCreatableEntitlementBundleAccessScopeSettings == nil {
			accessScopeSettings.AccessScopeSettingsCreatableEntitlementBundleAccessScopeSettings = &governance.AccessScopeSettingsCreatableEntitlementBundleAccessScopeSettings{}
		}
		accessScopeSettings.AccessScopeSettingsCreatableEntitlementBundleAccessScopeSettings.Type = "ENTITLEMENT_BUNDLES"
		var entitlementBundles []governance.EntitlementBundlesArrayCreatableInner
		elems := data.AccessScopeSettings.Ids
		for _, elem := range elems {
			groupId := elem.Id
			entitlementBundles = append(entitlementBundles, governance.EntitlementBundlesArrayCreatableInner{
				Id: groupId.ValueString(),
			})
		}
		accessScopeSettings.AccessScopeSettingsCreatableEntitlementBundleAccessScopeSettings.EntitlementBundles = entitlementBundles
		patch.AccessScopeSettings = &accessScopeSettings
	} else if data.AccessScopeSettings.Type.ValueString() == "RESOURCE_DEFAULT" {
		if accessScopeSettings.AccessScopeSettingsCreatableResourceDefaultAccessScopeSettings == nil {
			accessScopeSettings.AccessScopeSettingsCreatableResourceDefaultAccessScopeSettings = &governance.AccessScopeSettingsCreatableResourceDefaultAccessScopeSettings{}
		}
		accessScopeSettings.AccessScopeSettingsCreatableResourceDefaultAccessScopeSettings.Type = "RESOURCE_DEFAULT"
		patch.AccessScopeSettings = &accessScopeSettings
	}

	var requestSettings governance.RequesterSettingsCreatableRequesterSettings
	if data.RequesterSettings.Type.ValueString() == "GROUPS" {
		if requestSettings.RequesterSettingsCreatableGroupsRequesterSettings == nil {
			requestSettings.RequesterSettingsCreatableGroupsRequesterSettings = &governance.RequesterSettingsCreatableGroupsRequesterSettings{}
		}
		requestSettings.RequesterSettingsCreatableGroupsRequesterSettings.Type = "GROUPS"
		var groups []governance.GroupsArrayCreatableInner
		elems := data.RequesterSettings.Ids
		for _, elem := range elems {
			groupId := elem.Id
			groups = append(groups, governance.GroupsArrayCreatableInner{
				Id: groupId.ValueString(),
			})
		}
		requestSettings.RequesterSettingsCreatableGroupsRequesterSettings.Groups = groups
		patch.RequesterSettings = &requestSettings
	} else if data.RequesterSettings.Type.ValueString() == "EVERYONE" {
		if requestSettings.EveryoneRequesterSettings == nil {
			requestSettings.EveryoneRequesterSettings = &governance.EveryoneRequesterSettings{}
		}
		requestSettings.EveryoneRequesterSettings.Type = "EVERYONE"
		patch.RequesterSettings = &requestSettings
	}

	var accessDurationSettings governance.AccessDurationSettingsPatchable
	if data.AccessDurationSettings != nil {
		if data.AccessDurationSettings.Type.ValueString() == "ADMIN_FIXED_DURATION" {
			durationSettings := &governance.AccessDurationSettingsAdminFixedDuration{}
			durationSettings.Duration = data.AccessDurationSettings.Duration.ValueString()
			durationSettings.Type = "ADMIN_FIXED_DURATION"
			accessDurationSettings.AccessDurationSettingsAdminFixedDuration = durationSettings
			patch.AccessDurationSettings.Set(&accessDurationSettings)
		} else if data.AccessDurationSettings.Type.ValueString() == "REQUESTER_SPECIFIED_DURATION" {
			durationSettings := &governance.AccessDurationSettingsRequesterSpecifiedDuration{}
			durationSettings.MaximumDuration = data.AccessDurationSettings.Duration.ValueString()
			durationSettings.Type = "REQUESTER_SPECIFIED_DURATION"
			accessDurationSettings.AccessDurationSettingsRequesterSpecifiedDuration = durationSettings
			patch.AccessDurationSettings.Set(&accessDurationSettings)
		}
	}
	return patch
}

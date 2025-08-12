package governance

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"time"

	"example.com/aditya-okta/okta-ig-sdk-golang/governance"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
}

func (r *requestConditionResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(request, response)
}

type requestConditionResource struct {
	*config.Config
}

type Settings struct {
	Type types.String `tfsdk:"type"`
	Id   types.List   `tfsdk:"id"`
}

type AccessDurationSettings struct {
	Type     types.String `tfsdk:"type"`
	Duration types.String `tfsdk:"duration"`
}

type requestConditionResourceModel struct {
	Id                     types.String            `tfschema:"id"`
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
				Computed: true,
			},
			"resource_id": schema.StringAttribute{
				Required: true,
			},
			"approval_sequence_id": schema.StringAttribute{
				Required: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"priority": schema.Int32Attribute{
				Optional: true,
			},
			"created":         schema.StringAttribute{Computed: true},
			"created_by":      schema.StringAttribute{Computed: true},
			"last_updated":    schema.StringAttribute{Computed: true},
			"last_updated_by": schema.StringAttribute{Computed: true},
			"status":          schema.StringAttribute{Computed: true},
		},
		Blocks: map[string]schema.Block{
			"access_scope_settings": schema.SingleNestedBlock{
				Blocks: map[string]schema.Block{
					"ids": schema.ListNestedBlock{
						Description: "List of groups/entitlement bundles.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Required:    true,
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
						Description: "List of teams/groups ids.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Required:    true,
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
						Optional: true,
					},
				},
			},
		},
	}
}

func (r *requestConditionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data requestConditionResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic

	// Example data value setting
	requestConditionReq, diags := createRequestCondtion(data)
	if diags.HasError() {
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}
	requestConditionResp, _, err := r.OktaGovernanceClient.OktaIGSDKClient().RequestConditionsAPI.CreateResourceRequestConditionV2(ctx, data.ResourceId.ValueString()).RequestConditionCreatable(*requestConditionReq).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Request conditions",
			"Could not create Request conditions, unexpected error: "+err.Error(),
		)
		return
	}

	data.Id = types.StringValue(requestConditionResp.GetId())
	data.Name = types.StringValue(requestConditionResp.GetName())
	data.Description = types.StringValue(requestConditionResp.GetDescription())
	data.Priority = types.Int32Value(requestConditionResp.GetPriority())
	data.ApprovalSequenceId = types.StringValue(requestConditionResp.GetApprovalSequenceId())
	data.Created = types.StringValue(requestConditionResp.GetCreated().Format(time.RFC3339))
	data.CreatedBy = types.StringValue(requestConditionResp.GetCreatedBy())
	data.LastUpdated = types.StringValue(requestConditionResp.GetLastUpdated().Format(time.RFC3339))
	data.LastUpdatedBy = types.StringValue(requestConditionResp.GetLastUpdatedBy())
	data.Status = types.StringValue(string(requestConditionResp.GetStatus()))
	data.RequesterSettings, diags = setRequesterSettings(requestConditionResp.GetRequesterSettings())
	if diags.HasError() {
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}
	data.AccessScopeSettings, diags = setAccessScopeSettings(requestConditionResp.GetAccessScopeSettings())
	if diags.HasError() {
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	data.AccessDurationSettings = setAccessDurationSettings(requestConditionResp.GetAccessDurationSettings())

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func setAccessDurationSettings(settings governance.AccessDurationSettingsFull) *AccessDurationSettings {
	var accessDurationSettings AccessDurationSettings
	if settings.AccessDurationSettingsAdminFixedDuration != nil {
		accessDurationSettings.Type = types.StringValue(settings.AccessDurationSettingsAdminFixedDuration.Type)
		accessDurationSettings.Duration = types.StringValue(settings.AccessDurationSettingsAdminFixedDuration.Duration)
	} else if settings.AccessDurationSettingsRequesterSpecifiedDuration != nil {
		accessDurationSettings.Type = types.StringValue(settings.AccessDurationSettingsRequesterSpecifiedDuration.Type)
		accessDurationSettings.Duration = types.StringValue(settings.AccessDurationSettingsRequesterSpecifiedDuration.MaximumDuration)
	}
	return &accessDurationSettings
}

func setAccessScopeSettings(settings governance.AccessScopeSettingsFullAccessScopeSettings) (*Settings, diag.Diagnostics) {
	var setting Settings
	if settings.GroupAccessScopeSettings != nil {
		setting.Type = types.StringValue(settings.GroupAccessScopeSettings.GetType())
		groupIds := settings.GroupAccessScopeSettings.GetGroups()
		var ids []attr.Value
		for _, groupId := range groupIds {
			ids = append(ids, types.StringValue(groupId.GetId()))
		}
		listVal, diags := types.ListValue(types.StringType, ids)
		if diags.HasError() {
			return nil, diags
		}
		setting.Id = listVal
	} else if settings.EntitlementBundleAccessScopeSettings != nil {
		setting.Type = types.StringValue(settings.EntitlementBundleAccessScopeSettings.GetType())
		entitlementBundleIds := settings.EntitlementBundleAccessScopeSettings.GetEntitlementBundles()
		var ids []attr.Value
		for _, entitlementBundleId := range entitlementBundleIds {
			ids = append(ids, types.StringValue(entitlementBundleId.GetId()))
		}
		listVal, diags := types.ListValue(types.StringType, ids)
		if diags.HasError() {
			return nil, diags
		}
		setting.Id = listVal
	} else if settings.ResourceDefaultAccessScopeSettings != nil {
		setting.Type = types.StringValue(settings.EntitlementBundleAccessScopeSettings.GetType())
	}
	return &setting, nil
}

func setRequesterSettings(settings governance.RequesterSettingsFullRequesterSettings) (*Settings, diag.Diagnostics) {
	var setting Settings
	if settings.GroupsRequesterSettings != nil {
		setting.Type = types.StringValue(settings.GroupsRequesterSettings.GetType())
		groupIds := settings.GroupsRequesterSettings.GetGroups()
		var ids []attr.Value
		for _, groupId := range groupIds {
			ids = append(ids, types.StringValue(groupId.GetId()))
		}
		listVal, diags := types.ListValue(types.StringType, ids)
		if diags.HasError() {
			return nil, diags
		}
		setting.Id = listVal
	} else if settings.EveryoneRequesterSettings != nil {
		setting.Type = types.StringValue(settings.EveryoneRequesterSettings.GetType())
	}
	return &setting, nil
}

func createRequestCondtion(data requestConditionResourceModel) (*governance.RequestConditionCreatable, diag.Diagnostics) {
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
		accessScopeSettings.AccessScopeSettingsCreatableGroupAccessScopeSettings.Type = "GROUPS"
		var groupsIds []governance.GroupsArrayCreatableInner
		elems := data.AccessScopeSettings.Id.Elements()
		for _, elem := range elems {
			groupId := elem.(types.String)
			groupsIds = append(groupsIds, governance.GroupsArrayCreatableInner{
				Id: groupId.ValueString(),
			})
		}
		accessScopeSettings.AccessScopeSettingsCreatableGroupAccessScopeSettings.Groups = groupsIds
		req.AccessScopeSettings = accessScopeSettings
	} else if data.AccessScopeSettings.Type.ValueString() == "ENTITLEMENT_BUNDLES" {
		accessScopeSettings.AccessScopeSettingsCreatableEntitlementBundleAccessScopeSettings.Type = "ENTITLEMENT_BUNDLES"
		var entitlementBundles []governance.EntitlementBundlesArrayCreatableInner
		elems := data.AccessScopeSettings.Id.Elements()
		for _, elem := range elems {
			groupId := elem.(types.String)
			entitlementBundles = append(entitlementBundles, governance.EntitlementBundlesArrayCreatableInner{
				Id: groupId.ValueString(),
			})
		}
		accessScopeSettings.AccessScopeSettingsCreatableEntitlementBundleAccessScopeSettings.EntitlementBundles = entitlementBundles
		req.AccessScopeSettings = accessScopeSettings
	} else if data.AccessScopeSettings.Type.ValueString() == "RESOURCE_DEFAULT" {
		accessScopeSettings.AccessScopeSettingsCreatableResourceDefaultAccessScopeSettings.Type = "RESOURCE_DEFAULT"
		req.AccessScopeSettings = accessScopeSettings
	}

	requestSettings := governance.RequesterSettingsCreatableRequesterSettings{}
	if data.RequesterSettings.Type.ValueString() == "GROUPS" {
		requestSettings.RequesterSettingsCreatableGroupsRequesterSettings.Type = "GROUPS"
		var groups []governance.GroupsArrayCreatableInner
		elems := data.RequesterSettings.Id.Elements()
		for _, elem := range elems {
			groupId := elem.(types.String)
			groups = append(groups, governance.GroupsArrayCreatableInner{
				Id: groupId.ValueString(),
			})
		}
		requestSettings.RequesterSettingsCreatableGroupsRequesterSettings.Groups = groups
		req.RequesterSettings = requestSettings
	} else if data.RequesterSettings.Type.ValueString() == "EVERYONE" {
		requestSettings.EveryoneRequesterSettings.Type = "GROUPS"
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

	return &req, nil
}

func (r *requestConditionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data requestConditionResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	readRequestConditionResp, _, err := r.OktaGovernanceClient.OktaIGSDKClient().RequestConditionsAPI.GetResourceRequestConditionV2(ctx, data.ResourceId.ValueString(), data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Request conditions",
			"Could not read Request conditions, unexpected error: "+err.Error(),
		)
		return
	}
	data.Id = types.StringValue(readRequestConditionResp.GetId())
	data.Name = types.StringValue(readRequestConditionResp.GetName())
	data.Description = types.StringValue(readRequestConditionResp.GetDescription())
	data.Priority = types.Int32Value(readRequestConditionResp.GetPriority())
	data.ApprovalSequenceId = types.StringValue(readRequestConditionResp.GetApprovalSequenceId())
	data.Created = types.StringValue(readRequestConditionResp.GetCreated().Format(time.RFC3339))
	data.CreatedBy = types.StringValue(readRequestConditionResp.GetCreatedBy())
	data.LastUpdated = types.StringValue(readRequestConditionResp.GetLastUpdated().Format(time.RFC3339))
	data.LastUpdatedBy = types.StringValue(readRequestConditionResp.GetLastUpdatedBy())
	data.Status = types.StringValue(string(readRequestConditionResp.GetStatus()))
	data.RequesterSettings, _ = setRequesterSettings(readRequestConditionResp.GetRequesterSettings())
	data.AccessScopeSettings, _ = setAccessScopeSettings(readRequestConditionResp.GetAccessScopeSettings())
	data.AccessDurationSettings = setAccessDurationSettings(readRequestConditionResp.GetAccessDurationSettings())
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *requestConditionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data requestConditionResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	patch, diags := createRequestConditionPatch(data)
	if diags.HasError() {
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}
	// Update API call logic
	updatedRequestCondition, _, err := r.OktaGovernanceClient.OktaIGSDKClient().RequestConditionsAPI.UpdateResourceRequestConditionV2(ctx, data.ResourceId.ValueString(), data.Id.ValueString()).RequestConditionPatchable(patch).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Request conditions",
			"Could not update Request conditions, unexpected error: "+err.Error(),
		)
		return
	}

	data.Id = types.StringValue(updatedRequestCondition.GetId())
	data.Name = types.StringValue(updatedRequestCondition.GetName())
	data.Description = types.StringValue(updatedRequestCondition.GetDescription())
	data.Priority = types.Int32Value(updatedRequestCondition.GetPriority())
	data.ApprovalSequenceId = types.StringValue(updatedRequestCondition.GetApprovalSequenceId())
	data.Created = types.StringValue(updatedRequestCondition.GetCreated().Format(time.RFC3339))
	data.CreatedBy = types.StringValue(updatedRequestCondition.GetCreatedBy())
	data.LastUpdated = types.StringValue(updatedRequestCondition.GetLastUpdated().Format(time.RFC3339))
	data.LastUpdatedBy = types.StringValue(updatedRequestCondition.GetLastUpdatedBy())
	data.Status = types.StringValue(string(updatedRequestCondition.GetStatus()))
	data.RequesterSettings, diags = setRequesterSettings(updatedRequestCondition.GetRequesterSettings())
	if diags.HasError() {
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}
	data.AccessScopeSettings, diags = setAccessScopeSettings(updatedRequestCondition.GetAccessScopeSettings())
	if diags.HasError() {
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	data.AccessDurationSettings = setAccessDurationSettings(updatedRequestCondition.GetAccessDurationSettings())
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func createRequestConditionPatch(data requestConditionResourceModel) (governance.RequestConditionPatchable, diag.Diagnostics) {
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
		accessScopeSettings.AccessScopeSettingsCreatableGroupAccessScopeSettings.Type = "GROUPS"
		var groupsIds []governance.GroupsArrayCreatableInner
		elems := data.AccessScopeSettings.Id.Elements()
		for _, elem := range elems {
			groupId := elem.(types.String)
			groupsIds = append(groupsIds, governance.GroupsArrayCreatableInner{
				Id: groupId.ValueString(),
			})
		}
		accessScopeSettings.AccessScopeSettingsCreatableGroupAccessScopeSettings.Groups = groupsIds
		patch.AccessScopeSettings = &accessScopeSettings
	} else if data.AccessScopeSettings.Type.ValueString() == "ENTITLEMENT_BUNDLES" {
		accessScopeSettings.AccessScopeSettingsCreatableEntitlementBundleAccessScopeSettings.Type = "ENTITLEMENT_BUNDLES"
		var entitlementBundles []governance.EntitlementBundlesArrayCreatableInner
		elems := data.AccessScopeSettings.Id.Elements()
		for _, elem := range elems {
			groupId := elem.(types.String)
			entitlementBundles = append(entitlementBundles, governance.EntitlementBundlesArrayCreatableInner{
				Id: groupId.ValueString(),
			})
		}
		accessScopeSettings.AccessScopeSettingsCreatableEntitlementBundleAccessScopeSettings.EntitlementBundles = entitlementBundles
		patch.AccessScopeSettings = &accessScopeSettings
	} else if data.AccessScopeSettings.Type.ValueString() == "RESOURCE_DEFAULT" {
		accessScopeSettings.AccessScopeSettingsCreatableResourceDefaultAccessScopeSettings.Type = "RESOURCE_DEFAULT"
		patch.AccessScopeSettings = &accessScopeSettings
	}

	var requestSettings governance.RequesterSettingsCreatableRequesterSettings
	if data.RequesterSettings.Type.ValueString() == "GROUPS" {
		requestSettings.RequesterSettingsCreatableGroupsRequesterSettings.Type = "GROUPS"
		var groups []governance.GroupsArrayCreatableInner
		elems := data.RequesterSettings.Id.Elements()
		for _, elem := range elems {
			groupId := elem.(types.String)
			groups = append(groups, governance.GroupsArrayCreatableInner{
				Id: groupId.ValueString(),
			})
		}
		requestSettings.RequesterSettingsCreatableGroupsRequesterSettings.Groups = groups
		patch.RequesterSettings = &requestSettings
	} else if data.RequesterSettings.Type.ValueString() == "EVERYONE" {
		requestSettings.EveryoneRequesterSettings.Type = "GROUPS"
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
	return patch, nil
}

func (r *requestConditionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data requestConditionResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	_, err := r.OktaGovernanceClient.OktaIGSDKClient().RequestConditionsAPI.DeleteResourceRequestConditionV2(ctx, data.ResourceId.ValueString(), data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Request conditions",
			"Could not delete Request conditions, unexpected error: "+err.Error(),
		)
		return
	}
}

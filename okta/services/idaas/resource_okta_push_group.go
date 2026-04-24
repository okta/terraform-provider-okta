package idaas

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	v6okta "github.com/okta/okta-sdk-golang/v6/okta"
	"github.com/okta/terraform-provider-okta/okta/config"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

type pushGroupResourceModel struct {
	ID                         types.String `tfsdk:"id"`
	AppId                      types.String `tfsdk:"app_id"`
	SourceGroupId              types.String `tfsdk:"source_group_id"`
	TargetGroupId              types.String `tfsdk:"target_group_id"`
	TargetGroupName            types.String `tfsdk:"target_group_name"`
	Status                     types.String `tfsdk:"status"`
	DeleteTargetGroupOnDestroy types.Bool   `tfsdk:"delete_target_group_on_destroy"`
	AppConfig                  types.Object `tfsdk:"app_config"`
}

type pushGroupResource struct {
	config *config.Config
}

func newPushGroupResource() resource.Resource {
	return &pushGroupResource{}
}

func (r *pushGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_push_group"
}

func (r *pushGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Push Group ID",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"app_id": schema.StringAttribute{
				Required:    true,
				Computed:    false,
				Description: "The ID of the Okta Application.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"source_group_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the source group in Okta.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"target_group_name": schema.StringAttribute{
				Required:    false,
				Optional:    true,
				Computed:    false,
				Description: "The name of the target group for the group push mapping. This is used when creating a new downstream group. If the group already exists, it links to the existing group. If not specified, the name of the source group will be used as the name of the target group. Setting a target group name only works if you have unchecked 'Rename app groups to match group name in Okta' in the push groups settings UI.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"target_group_id": schema.StringAttribute{
				Required:    false,
				Optional:    true,
				Computed:    true,
				Description: "The ID of the target group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The status of the push group mapping. Valid values: `ACTIVE` and `INACTIVE`",
				Default:     stringdefault.StaticString("ACTIVE"),
				Validators: []validator.String{
					stringvalidator.OneOf("ACTIVE", "INACTIVE"),
				},
			},
			"delete_target_group_on_destroy": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to delete the target group when the push group mapping is destroyed. Default is true.",
				Default:     booldefault.StaticBool(true),
			},
			"app_config": schema.ObjectAttribute{
				Required:    false,
				Optional:    true,
				Computed:    false,
				Description: "Additional app configuration for group push mappings. Currently only required for Active Directory.",
				AttributeTypes: map[string]attr.Type{
					"distinguished_name": types.StringType,
					"group_scope":        types.StringType,
					"group_type":         types.StringType,
					"sam_account_name":   types.StringType,
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
					objectplanmodifier.RequiresReplace(),
				},
			},
		},
		Description: "Creates a Push Group assignment for an Application in Okta.",
	}
}

func (r *pushGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.config = resourceConfiguration(req, resp)
}

func (r *pushGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data pushGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var appConfig *v6okta.AppConfig
	appConfigAttrs := map[string]any{}
	for k, v := range data.AppConfig.Attributes() {
		value, err := v.ToTerraformValue(ctx)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Error converting app config attribute %s to terraform value", k), err.Error())
			return
		}
		if value.IsNull() {
			continue
		}

		var val string
		if err = value.As(&val); err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Error converting app config attribute %s to string", k), err.Error())
			return
		}

		appConfigAttrs[utils.UnderscoreToCamelCase(k)] = val
	}
	if len(appConfigAttrs) > 0 {
		appConfig = &v6okta.AppConfig{
			AdditionalProperties: appConfigAttrs,
		}
	}

	request := v6okta.CreateGroupPushMappingRequest{
		Status:          data.Status.ValueStringPointer(),
		SourceGroupId:   data.SourceGroupId.ValueString(),
		TargetGroupName: data.TargetGroupName.ValueStringPointer(),
		TargetGroupId:   data.TargetGroupId.ValueStringPointer(),
		AppConfig:       appConfig,
	}

	if data.TargetGroupName.ValueString() == "" && data.TargetGroupId.ValueString() == "" {
		sourceGroup, _, err := r.config.OktaIDaaSClient.OktaSDKClientV6().GroupAPI.GetGroup(ctx, data.SourceGroupId.ValueString()).Execute()
		if err != nil {
			resp.Diagnostics.AddError("Error creating Okta Push Group ", err.Error())
			return
		}
		if sourceGroup.Profile.OktaActiveDirectoryGroupProfile != nil {
			request.TargetGroupName = sourceGroup.Profile.OktaActiveDirectoryGroupProfile.Name
		} else if sourceGroup.Profile.OktaUserGroupProfile != nil {
			request.TargetGroupName = sourceGroup.Profile.OktaUserGroupProfile.Name
		} else {
			resp.Diagnostics.AddError("Error creating Okta Push Group ", "Target group name and target group id are both not set, and the provider failed to retrieve the source group name to use as the target group name. Please specify either the target group name or the target group id to fix this issue.")
			return
		}
	}
	if data.TargetGroupId.ValueString() == "" {
		request.TargetGroupId = nil
	}

	groupPushMapping, _, err := r.config.OktaIDaaSClient.OktaSDKClientV6().GroupPushMappingAPI.CreateGroupPushMapping(ctx, data.AppId.ValueString()).Body(request).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error creating Okta Push Group ", err.Error())
		return
	}

	resp.Diagnostics.Append(mapPushGroupResourceToState(groupPushMapping, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *pushGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state pushGroupResourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupPushMapping, _, err := r.config.OktaIDaaSClient.OktaSDKClientV6().GroupPushMappingAPI.GetGroupPushMapping(ctx, state.AppId.ValueString(), state.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error reading Okta push group mapping ", err.Error())
		return
	}

	resp.Diagnostics.Append(mapPushGroupResourceToState(groupPushMapping, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *pushGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state pushGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Status.ValueString() != "ACTIVE" && state.Status.ValueString() != "INACTIVE" {
		resp.Diagnostics.AddError("Group push mapping in disallowed status", "The Okta API doesn't allow updating the group push mapping when the status is not either ACTIVE or INACTIVE")
		return
	}

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupPushMapping, _, err := r.config.OktaIDaaSClient.OktaSDKClientV6().GroupPushMappingAPI.UpdateGroupPushMapping(ctx, state.AppId.ValueString(), state.ID.ValueString()).Body(v6okta.UpdateGroupPushMappingRequest{
		Status: state.Status.ValueString(),
	}).Execute()
	if err != nil {
		resp.Diagnostics.AddError("failed to update push group mapping: ", err.Error())
		return
	}

	resp.Diagnostics.Append(mapPushGroupResourceToState(groupPushMapping, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *pushGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state pushGroupResourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Status.ValueString() != "ACTIVE" && state.Status.ValueString() != "INACTIVE" {
		resp.Diagnostics.AddError("Group push mapping in disallowed status", "To delete a group push mapping, the status must be INACTIVE. Setting the status to INACTIVE however failed because the Okta API doesn't allow updating the group push mapping when the status is not either ACTIVE or INACTIVE")
		return
	}

	if state.Status.ValueString() != "INACTIVE" {
		_, _, err := r.config.OktaIDaaSClient.OktaSDKClientV6().GroupPushMappingAPI.UpdateGroupPushMapping(ctx, state.AppId.ValueString(), state.ID.ValueString()).Body(v6okta.UpdateGroupPushMappingRequest{
			Status: "INACTIVE",
		}).Execute()
		if err != nil {
			resp.Diagnostics.AddError("failed to delete push group mapping: ", err.Error())
			return
		}
	}

	_, err := r.config.OktaIDaaSClient.OktaSDKClientV6().GroupPushMappingAPI.DeleteGroupPushMapping(ctx, state.AppId.ValueString(), state.ID.ValueString()).DeleteTargetGroup(state.DeleteTargetGroupOnDestroy.ValueBool()).Execute()
	if err != nil {
		resp.Diagnostics.AddError("failed to delete push group mapping: ", err.Error())
		return
	}
}

func (r *pushGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importID := req.ID
	if importID == "" {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Import ID cannot be empty. Expected format: app_uid/mapping_id",
		)
		return
	}

	parts := strings.Split(importID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import ID format",
			fmt.Sprintf("Expected import ID format 'app_uid/mapping_id', got '%s'", importID),
		)
		return
	}

	appId := parts[0]
	mappingId := parts[1]

	if appId == "" || mappingId == "" {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Both app_uid and mapping_id must be provided in import ID",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("app_id"), appId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), mappingId)...)
}

func mapPushGroupResourceToState(groupPushMapping *v6okta.GroupPushMapping, state *pushGroupResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	state.ID = types.StringPointerValue(groupPushMapping.Id)
	state.SourceGroupId = types.StringPointerValue(groupPushMapping.SourceGroupId)
	state.TargetGroupId = types.StringPointerValue(groupPushMapping.TargetGroupId)
	state.Status = types.StringPointerValue(groupPushMapping.Status)
	return diags
}

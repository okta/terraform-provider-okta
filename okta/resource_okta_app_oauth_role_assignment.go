package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/sdk"
)

var _ resource.ResourceWithValidateConfig = &appOAuthRoleAssignmentResource{}

func NewAppOAuthRoleAssignmentResource() resource.Resource {
	return &appOAuthRoleAssignmentResource{}
}

type appOAuthRoleAssignmentResource struct {
	*Config
}

type OAuthRoleAssignmentResourceModel struct {
	ID          types.String `tfsdk:"id"`
	ClientID    types.String `tfsdk:"client_id"`
	Type        types.String `tfsdk:"type"`
	ResourceSet types.String `tfsdk:"resource_set"`
	Role        types.String `tfsdk:"role"`
	Status      types.String `tfsdk:"status"`
	Label       types.String `tfsdk:"label"`
}

func (r *appOAuthRoleAssignmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_oauth_role_assignment"
}

func (r *appOAuthRoleAssignmentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	config, ok := req.ProviderData.(*Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.Config = config
}

func (r *appOAuthRoleAssignmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages assignment of an admin role to an OAuth application",
		MarkdownDescription: "Manages assignment of an admin role to an OAuth application",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Role Assignment ID",
				MarkdownDescription: "Role Assignment ID",
				Computed:            true,
			},
			"client_id": schema.StringAttribute{
				Description:         "Client ID for the role to be assigned to",
				MarkdownDescription: "Client ID for the role to be assigned to",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Description:         "Role type to assign. This can be one of the standard Okta roles, such as `HELP_DESK_ADMIN` or `CUSTOM`. Using custom requires the `resource_set` and `role` attributes to be set.",
				MarkdownDescription: "Role type to assign. This can be one of the standard Okta roles, such as `HELP_DESK_ADMIN` or `CUSTOM`. Using custom requires the `resource_set` and `role` attributes to be set.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"resource_set": schema.StringAttribute{
				Description:         "Resource set for the custom role to assign, must be the ID of the created resource set.",
				MarkdownDescription: "Resource set for the custom role to assign, must be the ID of the created resource set.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"role": schema.StringAttribute{
				Description:         "Custom Role ID",
				MarkdownDescription: "Custom Role ID",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				Description:         "Status of the role assignment",
				MarkdownDescription: "Status of the role assignment",
				Computed:            true,
			},
			"label": schema.StringAttribute{
				Description:         "Label of the role assignment",
				MarkdownDescription: "Label of the role assignment",
				Computed:            true,
			},
		},
	}
}

func (r *appOAuthRoleAssignmentResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data *OAuthRoleAssignmentResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if data.Type.ValueString() == "CUSTOM" && (data.ResourceSet.IsNull() || data.Role.IsNull()) {
		resp.Diagnostics.AddAttributeError(
			path.Root("type"),
			"Missing attribute configuration",
			"When type is set to 'CUSTOM', the resource_set and role attributes must be set.",
		)
	}
}

func (r *appOAuthRoleAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *OAuthRoleAssignmentResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	roleAssignmentRequest := &sdk.ClientRoleAssignment{
		Type:        data.Type.ValueString(),
		ResourceSet: data.ResourceSet.ValueStringPointer(),
		Role:        data.Role.ValueStringPointer(),
	}

	role, _, err := r.Config.oktaSDKsupplementClient.AssignClientRole(ctx, data.ClientID.ValueString(), roleAssignmentRequest)
	if err != nil {
		resp.Diagnostics.AddError("Unable to assign role to client", err.Error())
		return
	}

	data.ID = types.StringPointerValue(role.Id)
	data.Status = types.StringPointerValue(role.Status)
	data.Label = types.StringPointerValue(role.Label)
	data.Type = types.StringPointerValue(role.Type)
	data.ResourceSet = types.StringPointerValue(role.ResourceSet)
	data.Role = types.StringPointerValue(role.Role)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *appOAuthRoleAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *OAuthRoleAssignmentResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	role, _, err := r.Config.oktaSDKsupplementClient.GetClientRole(ctx, data.ClientID.ValueString(), data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read role assignment", err.Error())
		return
	}

	data.ID = types.StringPointerValue(role.Id)
	data.Status = types.StringPointerValue(role.Status)
	data.Label = types.StringPointerValue(role.Label)
	data.Type = types.StringPointerValue(role.Type)
	data.ResourceSet = types.StringPointerValue(role.ResourceSet)
	data.Role = types.StringPointerValue(role.Role)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *appOAuthRoleAssignmentResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update not supported",
		"OAuth Role Assignments cannot be updated. If you get to this contact the provider maintainers as this should not be hit.",
	)
}

func (r *appOAuthRoleAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *OAuthRoleAssignmentResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.Config.oktaSDKsupplementClient.UnassignClientRole(ctx, data.ClientID.ValueString(), data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete role assignment", err.Error())
		return
	}
}

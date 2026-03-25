package idaas

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
	"github.com/okta/terraform-provider-okta/okta/utils"

	fwdiag "github.com/hashicorp/terraform-plugin-framework/diag"
	v6okta "github.com/okta/okta-sdk-golang/v6/okta"
)

var (
	_ resource.Resource                = &userRiskResource{}
	_ resource.ResourceWithConfigure   = &userRiskResource{}
	_ resource.ResourceWithImportState = &userRiskResource{}
)

func newUserRiskResource() resource.Resource {
	return &userRiskResource{}
}

type userRiskResource struct {
	*config.Config
}

type userRiskResourceModel struct {
	ID        types.String `tfsdk:"id"`
	UserID    types.String `tfsdk:"user_id"`
	RiskLevel types.String `tfsdk:"risk_level"`
}

func (r *userRiskResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_risk"
}

func (r *userRiskResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a user's risk level in Okta. This resource allows you to set and manage the risk level for a specific user.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of this resource (same as user_id).",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_id": schema.StringAttribute{
				Description: "ID of the user.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"risk_level": schema.StringAttribute{
				Description: "Risk level of the user. Valid values: `HIGH`, `LOW`.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("HIGH", "LOW"),
				},
			},
		},
	}
}

func (r *userRiskResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *userRiskResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state userRiskResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	userId := state.UserID.ValueString()
	riskLevel := state.RiskLevel.ValueString()

	r.Logger.Info("setting user risk", "user_id", userId, "risk_level", riskLevel)

	userRiskReq := v6okta.NewUserRiskRequest()
	userRiskReq.SetRiskLevel(riskLevel)

	_, _, err := r.OktaIDaaSClient.OktaSDKClientV6().UserRiskAPI.UpsertUserRisk(ctx, userId).UserRiskRequest(*userRiskReq).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to set user risk",
			utils.ErrorDetail_V6(err),
		)
		return
	}

	resp.Diagnostics.Append(r.readUserRisk(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *userRiskResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state userRiskResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.readUserRisk(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If the resource was removed (NONE risk level), signal removal from state
	if state.ID.IsNull() {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *userRiskResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state userRiskResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	userId := state.UserID.ValueString()
	riskLevel := state.RiskLevel.ValueString()

	r.Logger.Info("updating user risk", "user_id", userId, "risk_level", riskLevel)

	userRiskReq := v6okta.NewUserRiskRequest()
	userRiskReq.SetRiskLevel(riskLevel)

	_, _, err := r.OktaIDaaSClient.OktaSDKClientV6().UserRiskAPI.UpsertUserRisk(ctx, userId).UserRiskRequest(*userRiskReq).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to update user risk",
			utils.ErrorDetail_V6(err),
		)
		return
	}

	resp.Diagnostics.Append(r.readUserRisk(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *userRiskResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state userRiskResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// On destroy, we simply remove the resource from Terraform state.
	// The user's risk level will remain at whatever it was last set to.
	r.Logger.Info("removing user risk from state", "user_id", state.ID.ValueString())
}

func (r *userRiskResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	userId := req.ID

	if strings.TrimSpace(userId) == "" {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"user_id is required for import",
		)
		return
	}

	riskResp, _, err := r.OktaIDaaSClient.OktaSDKClientV6().UserRiskAPI.GetUserRisk(ctx, userId).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to get user risk",
			utils.ErrorDetail_V6(err),
		)
		return
	}

	// If the user has no risk level set (NONE), we cannot import
	if riskResp.UserRiskLevelNone != nil {
		resp.Diagnostics.AddError(
			"Cannot import user with no risk level",
			"User has no risk level set (NONE). Set a risk level (HIGH or LOW) before importing, or create a new okta_user_risk resource instead.",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), userId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("user_id"), userId)...)
}

func (r *userRiskResource) readUserRisk(ctx context.Context, state *userRiskResourceModel) (diags fwdiag.Diagnostics) {
	userId := state.UserID.ValueString()
	if userId == "" {
		userId = state.ID.ValueString()
	}

	r.Logger.Info("reading user risk", "user_id", userId)

	riskResp, _, err := r.OktaIDaaSClient.OktaSDKClientV6().UserRiskAPI.GetUserRisk(ctx, userId).Execute()
	if err != nil {
		diags.AddError(
			"failed to get user risk",
			utils.ErrorDetail_V6(err),
		)
		return
	}

	if riskResp.UserRiskLevelExists != nil {
		state.ID = types.StringValue(userId)
		state.UserID = types.StringValue(userId)
		state.RiskLevel = types.StringValue(riskResp.UserRiskLevelExists.GetRiskLevel())
	} else if riskResp.UserRiskLevelNone != nil {
		// User has no risk level set - signal removal from state
		state.ID = types.StringNull()
	}

	return
}

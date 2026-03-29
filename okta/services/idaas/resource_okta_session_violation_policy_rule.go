package idaas

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"

	fwdiag "github.com/hashicorp/terraform-plugin-framework/diag"
	v6okta "github.com/okta/okta-sdk-golang/v6/okta"
)

var (
	_ resource.Resource                = &sessionViolationPolicyRuleResource{}
	_ resource.ResourceWithConfigure   = &sessionViolationPolicyRuleResource{}
	_ resource.ResourceWithImportState = &sessionViolationPolicyRuleResource{}
)

func newSessionViolationPolicyRuleResource() resource.Resource {
	return &sessionViolationPolicyRuleResource{}
}

type sessionViolationPolicyRuleResource struct {
	*config.Config
}

type sessionViolationPolicyRuleResourceModel struct {
	ID                      types.String `tfsdk:"id"`
	PolicyID                types.String `tfsdk:"policy_id"`
	Name                    types.String `tfsdk:"name"`
	Status                  types.String `tfsdk:"status"`
	NetworkConnection       types.String `tfsdk:"network_connection"`
	NetworkIncludes         types.List   `tfsdk:"network_includes"`
	NetworkExcludes         types.List   `tfsdk:"network_excludes"`
	RiskScoreLevel          types.String `tfsdk:"risk_score_level"`
	MinRiskLevel            types.String `tfsdk:"min_risk_level"`
	PolicyEvaluationEnabled types.Bool   `tfsdk:"policy_evaluation_enabled"`
}

func (r *sessionViolationPolicyRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_session_violation_policy_rule"
}

func (r *sessionViolationPolicyRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages the Session Violation Detection Policy Rule. The Session Violation Detection Policy has exactly one modifiable rule (non-default). This resource allows you to configure that rule. Note: The rule cannot be created or deleted, only imported and modified.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ID of the policy rule.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"policy_id": schema.StringAttribute{
				Description: "ID of the Session Violation Detection Policy. Use the `okta_session_violation_policy` data source to get this ID.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the policy rule.",
				Optional:    true,
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Status of the rule: ACTIVE or INACTIVE.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("ACTIVE"),
				Validators: []validator.String{
					stringvalidator.OneOf("ACTIVE", "INACTIVE"),
				},
			},
			"network_connection": schema.StringAttribute{
				Description: "Network selection mode: ANYWHERE, ZONE, ON_NETWORK, or OFF_NETWORK.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("ANYWHERE", "ZONE", "ON_NETWORK", "OFF_NETWORK"),
				},
			},
			"network_includes": schema.ListAttribute{
				Description: "Required if `network_connection` is set to `ZONE`. List of network zone IDs to include.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"network_excludes": schema.ListAttribute{
				Description: "Required if `network_connection` is set to `ZONE`. List of network zone IDs to exclude.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"risk_score_level": schema.StringAttribute{
				Description: "The risk score level to match. Possible values: ANY, LOW, MEDIUM, HIGH.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("ANY", "LOW", "MEDIUM", "HIGH"),
				},
			},
			"min_risk_level": schema.StringAttribute{
				Description: "The minimum risk level to match. Only used in Session Violation Detection policy rules. Possible values: ANY, LOW, MEDIUM, HIGH.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("ANY", "LOW", "MEDIUM", "HIGH"),
				},
			},
			"policy_evaluation_enabled": schema.BoolAttribute{
				Description: "When true, the sign-on policies of the session are evaluated when a session violation is detected.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
		},
	}
}

func (r *sessionViolationPolicyRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *sessionViolationPolicyRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	resp.Diagnostics.AddError(
		"Session Violation Detection Policy Rule cannot be created",
		"The Session Violation Detection Policy Rule already exists and cannot be created, only imported and updated. "+
			"Use 'terraform import' with the policy ID and rule ID to import the existing rule. "+
			"You can find the rule ID via the `okta_session_violation_policy` data source or the Okta Admin Console or API.",
	)
}

func (r *sessionViolationPolicyRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state sessionViolationPolicyRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.readSessionViolationPolicyRule(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID.IsNull() {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *sessionViolationPolicyRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state sessionViolationPolicyRuleResourceModel
	var priorState sessionViolationPolicyRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &priorState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyId := state.PolicyID.ValueString()
	ruleId := priorState.ID.ValueString()
	state.ID = types.StringValue(ruleId)

	r.Logger.Info("updating session violation detection policy rule", "policy_id", policyId, "rule_id", ruleId)

	rule, diags := r.buildSessionViolationPolicyRule(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	rule.Id = &ruleId

	ruleRequest := v6okta.ListPolicyRules200ResponseInner{
		SessionViolationDetectionPolicyRule: rule,
	}

	_, _, err := r.OktaIDaaSClient.OktaSDKClientV6().PolicyAPI.ReplacePolicyRule(ctx, policyId, ruleId).PolicyRule(ruleRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to update session violation detection policy rule",
			err.Error(),
		)
		return
	}

	if !state.Status.Equal(priorState.Status) {
		status := state.Status.ValueString()
		if status == "ACTIVE" {
			_, err := r.OktaIDaaSClient.OktaSDKClientV6().PolicyAPI.ActivatePolicyRule(ctx, policyId, ruleId).Execute()
			if err != nil {
				resp.Diagnostics.AddError(
					"failed to activate session violation detection policy rule",
					err.Error(),
				)
				return
			}
		} else {
			_, err := r.OktaIDaaSClient.OktaSDKClientV6().PolicyAPI.DeactivatePolicyRule(ctx, policyId, ruleId).Execute()
			if err != nil {
				resp.Diagnostics.AddError(
					"failed to deactivate session violation detection policy rule",
					err.Error(),
				)
				return
			}
		}
	}

	resp.Diagnostics.Append(r.readSessionViolationPolicyRule(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *sessionViolationPolicyRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state sessionViolationPolicyRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyId := state.PolicyID.ValueString()
	ruleId := state.ID.ValueString()

	// The Session Violation Detection Policy rule cannot be deleted, only modified.
	// On destroy, we just remove it from Terraform state.
	// The rule will remain in Okta with its current configuration.
	r.Logger.Info("removing session violation detection policy rule from state (rule cannot be deleted)", "policy_id", policyId, "rule_id", ruleId)
}

func (r *sessionViolationPolicyRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import ID format",
			"Expected: <policy_id>/<rule_id>",
		)
		return
	}

	policyId := parts[0]
	ruleId := parts[1]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), ruleId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("policy_id"), policyId)...)
}

func (r *sessionViolationPolicyRuleResource) readSessionViolationPolicyRule(ctx context.Context, state *sessionViolationPolicyRuleResourceModel) (diags fwdiag.Diagnostics) {
	policyId := state.PolicyID.ValueString()
	ruleId := state.ID.ValueString()

	r.Logger.Info("reading session violation detection policy rule", "policy_id", policyId, "rule_id", ruleId)

	ruleResp, _, err := r.OktaIDaaSClient.OktaSDKClientV6().PolicyAPI.GetPolicyRule(ctx, policyId, ruleId).Execute()
	if err != nil {
		diags.AddError(
			"failed to get session violation detection policy rule",
			err.Error(),
		)
		return
	}

	if ruleResp.SessionViolationDetectionPolicyRule == nil {
		state.ID = types.StringNull()
		return
	}

	rule := ruleResp.SessionViolationDetectionPolicyRule

	// Check if this is the default rule (priority 99) - cannot be managed
	if rule.Priority.IsSet() {
		priority := rule.Priority.Get()
		if priority != nil && *priority == 99 {
			diags.AddError(
				"Cannot manage default policy rule",
				"The default Session Violation Detection Policy rule (priority 99) cannot be imported or modified. "+
					"Please import the non-default rule instead. Use the `okta_session_violation_policy` data source to get the correct rule ID.",
			)
			return
		}
	}

	state.Name = types.StringPointerValue(rule.Name)
	if rule.Status != nil {
		state.Status = types.StringValue(string(*rule.Status))
	}

	// Read network conditions
	if rule.Conditions != nil && rule.Conditions.Network != nil {
		network := rule.Conditions.Network
		if network.Connection != nil {
			state.NetworkConnection = types.StringValue(string(*network.Connection))
		}
		if len(network.Include) > 0 {
			includes, d := types.ListValueFrom(ctx, types.StringType, network.Include)
			diags.Append(d...)
			state.NetworkIncludes = includes
		}
		if len(network.Exclude) > 0 {
			excludes, d := types.ListValueFrom(ctx, types.StringType, network.Exclude)
			diags.Append(d...)
			state.NetworkExcludes = excludes
		}
	}

	// Read risk score conditions
	if rule.Conditions != nil && rule.Conditions.RiskScore != nil {
		riskScore := rule.Conditions.RiskScore
		state.RiskScoreLevel = types.StringValue(riskScore.Level)
		if riskScore.MinRiskLevel != nil {
			state.MinRiskLevel = types.StringPointerValue(riskScore.MinRiskLevel)
		}
	}

	// Read actions
	if rule.Actions != nil &&
		rule.Actions.SessionViolationDetection != nil &&
		rule.Actions.SessionViolationDetection.PolicyEvaluation != nil {
		policyEval := rule.Actions.SessionViolationDetection.PolicyEvaluation
		if policyEval.Enabled != nil {
			state.PolicyEvaluationEnabled = types.BoolPointerValue(policyEval.Enabled)
		}
	}

	return
}

func (r *sessionViolationPolicyRuleResource) buildSessionViolationPolicyRule(ctx context.Context, state *sessionViolationPolicyRuleResourceModel) (*v6okta.SessionViolationDetectionPolicyRule, fwdiag.Diagnostics) {
	var diags fwdiag.Diagnostics

	rule := v6okta.NewSessionViolationDetectionPolicyRule()

	if !state.Name.IsNull() && state.Name.ValueString() != "" {
		name := state.Name.ValueString()
		rule.Name = &name
	}

	ruleType := "SESSION_VIOLATION_DETECTION"
	rule.Type = &ruleType

	// Build conditions
	conditions := v6okta.NewSessionViolationDetectionPolicyRuleAllOfConditions()

	// Network condition
	if !state.NetworkConnection.IsNull() && state.NetworkConnection.ValueString() != "" {
		networkCondition := v6okta.NewPolicyNetworkCondition()
		conn := state.NetworkConnection.ValueString()
		networkCondition.Connection = &conn

		if !state.NetworkIncludes.IsNull() {
			var includes []string
			diags.Append(state.NetworkIncludes.ElementsAs(ctx, &includes, false)...)
			networkCondition.Include = includes
		}
		if !state.NetworkExcludes.IsNull() {
			var excludes []string
			diags.Append(state.NetworkExcludes.ElementsAs(ctx, &excludes, false)...)
			networkCondition.Exclude = excludes
		}

		conditions.Network = networkCondition
	}

	// Risk score condition
	if !state.RiskScoreLevel.IsNull() && state.RiskScoreLevel.ValueString() != "" {
		riskScore := v6okta.NewRiskScorePolicyRuleCondition(state.RiskScoreLevel.ValueString())
		if !state.MinRiskLevel.IsNull() && state.MinRiskLevel.ValueString() != "" {
			minLevel := state.MinRiskLevel.ValueString()
			riskScore.MinRiskLevel = &minLevel
		}
		conditions.RiskScore = riskScore
	}

	rule.Conditions = conditions

	// Build actions
	actions := v6okta.NewSessionViolationDetectionPolicyRuleAllOfActions()
	svdAction := v6okta.NewSessionViolationDetectionPolicyRuleAllOfActionsSessionViolationDetection()
	policyEval := v6okta.NewSessionViolationDetectionPolicyEvaluation()
	enabled := state.PolicyEvaluationEnabled.ValueBool()
	policyEval.SetEnabled(enabled)
	svdAction.SetPolicyEvaluation(*policyEval)
	actions.SetSessionViolationDetection(*svdAction)
	rule.Actions = actions

	return rule, diags
}

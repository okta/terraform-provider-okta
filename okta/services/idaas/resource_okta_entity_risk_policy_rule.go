package idaas

import (
	"context"
	"fmt"
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
	_ resource.Resource                = &entityRiskPolicyRuleResource{}
	_ resource.ResourceWithConfigure   = &entityRiskPolicyRuleResource{}
	_ resource.ResourceWithImportState = &entityRiskPolicyRuleResource{}
)

func newEntityRiskPolicyRuleResource() resource.Resource {
	return &entityRiskPolicyRuleResource{}
}

type entityRiskPolicyRuleResource struct {
	*config.Config
}

type entityRiskPolicyRuleResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	PolicyID             types.String `tfsdk:"policy_id"`
	Name                 types.String `tfsdk:"name"`
	Status               types.String `tfsdk:"status"`
	Priority             types.Int64  `tfsdk:"priority"`
	RiskLevel            types.String `tfsdk:"risk_level"`
	UsersIncluded        types.Set    `tfsdk:"users_included"`
	UsersExcluded        types.Set    `tfsdk:"users_excluded"`
	GroupsIncluded       types.Set    `tfsdk:"groups_included"`
	GroupsExcluded       types.Set    `tfsdk:"groups_excluded"`
	TerminateAllSessions types.Bool   `tfsdk:"terminate_all_sessions"`
	WorkflowID           types.String `tfsdk:"workflow_id"`
}

func (r *entityRiskPolicyRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entity_risk_policy_rule"
}

func (r *entityRiskPolicyRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Entity Risk Policy Rule. Entity Risk Policy rules define automated responses to identity threats detected by Okta's Identity Threat Protection (ITP).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ID of the policy rule.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"policy_id": schema.StringAttribute{
				Description: "ID of the Entity Risk Policy. Use the `okta_entity_risk_policy` data source to get this ID.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the policy rule.",
				Required:    true,
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
			"priority": schema.Int64Attribute{
				Description: "Priority of the rule. Rules are evaluated in priority order.",
				Optional:    true,
				Computed:    true,
			},
			"risk_level": schema.StringAttribute{
				Description: "Risk level to match. Valid values: HIGH, MEDIUM, LOW, ANY.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("HIGH", "MEDIUM", "LOW", "ANY"),
				},
			},
			"users_included": schema.SetAttribute{
				Description: "List of user IDs to include from this rule.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"users_excluded": schema.SetAttribute{
				Description: "List of user IDs to exclude from this rule.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"groups_included": schema.SetAttribute{
				Description: "List of group IDs to include in this rule.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"groups_excluded": schema.SetAttribute{
				Description: "List of group IDs to exclude from this rule.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"terminate_all_sessions": schema.BoolAttribute{
				Description: "When true, terminates all active sessions for the user when a risk event is detected.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"workflow_id": schema.StringAttribute{
				Description: "ID of the Okta Workflow to run when a risk event is detected.",
				Optional:    true,
			},
		},
	}
}

func (r *entityRiskPolicyRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *entityRiskPolicyRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state entityRiskPolicyRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyId := state.PolicyID.ValueString()
	name := state.Name.ValueString()

	r.Logger.Info("creating entity risk policy rule", "policy_id", policyId, "name", name)

	rule, diags := r.buildEntityRiskPolicyRule(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ruleRequest := v6okta.ListPolicyRules200ResponseInner{
		EntityRiskPolicyRule: rule,
	}

	ruleResp, _, err := r.OktaIDaaSClient.OktaSDKClientV6().PolicyAPI.CreatePolicyRule(ctx, policyId).PolicyRule(ruleRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to create entity risk policy rule",
			err.Error(),
		)
		return
	}

	if ruleResp.EntityRiskPolicyRule == nil || ruleResp.EntityRiskPolicyRule.Id == nil {
		resp.Diagnostics.AddError(
			"rule ID not found in response",
			"The API response did not contain a valid entity risk policy rule",
		)
		return
	}

	ruleId := *ruleResp.EntityRiskPolicyRule.Id
	state.ID = types.StringValue(ruleId)

	if state.Status.ValueString() == "INACTIVE" {
		_, err := r.OktaIDaaSClient.OktaSDKClientV6().PolicyAPI.DeactivatePolicyRule(ctx, policyId, ruleId).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"failed to deactivate entity risk policy rule",
				err.Error(),
			)
			return
		}
	}

	resp.Diagnostics.Append(r.readEntityRiskPolicyRule(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *entityRiskPolicyRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state entityRiskPolicyRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.readEntityRiskPolicyRule(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID.IsNull() {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *entityRiskPolicyRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state entityRiskPolicyRuleResourceModel
	var priorState entityRiskPolicyRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &priorState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyId := state.PolicyID.ValueString()
	ruleId := state.ID.ValueString()

	r.Logger.Info("updating entity risk policy rule", "policy_id", policyId, "rule_id", ruleId)

	rule, diags := r.buildEntityRiskPolicyRule(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	rule.Id = &ruleId

	ruleRequest := v6okta.ListPolicyRules200ResponseInner{
		EntityRiskPolicyRule: rule,
	}

	_, _, err := r.OktaIDaaSClient.OktaSDKClientV6().PolicyAPI.ReplacePolicyRule(ctx, policyId, ruleId).PolicyRule(ruleRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to update entity risk policy rule",
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
					"failed to activate entity risk policy rule",
					err.Error(),
				)
				return
			}
		} else {
			_, err := r.OktaIDaaSClient.OktaSDKClientV6().PolicyAPI.DeactivatePolicyRule(ctx, policyId, ruleId).Execute()
			if err != nil {
				resp.Diagnostics.AddError(
					"failed to deactivate entity risk policy rule",
					err.Error(),
				)
				return
			}
		}
	}

	resp.Diagnostics.Append(r.readEntityRiskPolicyRule(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *entityRiskPolicyRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state entityRiskPolicyRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyId := state.PolicyID.ValueString()
	ruleId := state.ID.ValueString()

	r.Logger.Info("deleting entity risk policy rule", "policy_id", policyId, "rule_id", ruleId)

	_, err := r.OktaIDaaSClient.OktaSDKClientV6().PolicyAPI.DeletePolicyRule(ctx, policyId, ruleId).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to delete entity risk policy rule",
			err.Error(),
		)
		return
	}
}

func (r *entityRiskPolicyRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *entityRiskPolicyRuleResource) readEntityRiskPolicyRule(ctx context.Context, state *entityRiskPolicyRuleResourceModel) (diags fwdiag.Diagnostics) {
	policyId := state.PolicyID.ValueString()
	ruleId := state.ID.ValueString()

	r.Logger.Info("reading entity risk policy rule", "policy_id", policyId, "rule_id", ruleId)

	ruleResp, _, err := r.OktaIDaaSClient.OktaSDKClientV6().PolicyAPI.GetPolicyRule(ctx, policyId, ruleId).Execute()
	if err != nil {
		diags.AddError(
			"failed to get entity risk policy rule",
			err.Error(),
		)
		return
	}

	if ruleResp.EntityRiskPolicyRule == nil {
		state.ID = types.StringNull()
		return
	}

	rule := ruleResp.EntityRiskPolicyRule

	// Check if this is the default rule (priority 99) - cannot be managed
	if rule.Priority.IsSet() {
		priority := rule.Priority.Get()
		if priority != nil && *priority == 99 {
			diags.AddError(
				"Cannot manage default policy rule",
				"The default Entity Risk Policy rule (priority 99) cannot be imported or modified. "+
					"Please create or import a non-default rule instead.",
			)
			return
		}
	}

	state.Name = types.StringPointerValue(rule.Name)
	if rule.Status != nil {
		state.Status = types.StringValue(string(*rule.Status))
	}
	if rule.Priority.IsSet() {
		state.Priority = types.Int64Value(int64(*rule.Priority.Get()))
	}

	if rule.Conditions != nil {
		if rule.Conditions.EntityRisk != nil {
			state.RiskLevel = types.StringValue(rule.Conditions.EntityRisk.Level)
		}

		if rule.Conditions.People != nil {
			if rule.Conditions.People.Users != nil && len(rule.Conditions.People.Users.Exclude) > 0 {
				usersExcluded, d := types.SetValueFrom(ctx, types.StringType, rule.Conditions.People.Users.Exclude)
				diags.Append(d...)
				state.UsersExcluded = usersExcluded
			}
			if rule.Conditions.People.Users != nil && len(rule.Conditions.People.Users.Include) > 0 {
				usersIncluded, d := types.SetValueFrom(ctx, types.StringType, rule.Conditions.People.Users.Include)
				diags.Append(d...)
				state.UsersIncluded = usersIncluded
			}
			if rule.Conditions.People.Groups != nil {
				if len(rule.Conditions.People.Groups.Include) > 0 {
					groupsIncluded, d := types.SetValueFrom(ctx, types.StringType, rule.Conditions.People.Groups.Include)
					diags.Append(d...)
					state.GroupsIncluded = groupsIncluded
				}
				if len(rule.Conditions.People.Groups.Exclude) > 0 {
					groupsExcluded, d := types.SetValueFrom(ctx, types.StringType, rule.Conditions.People.Groups.Exclude)
					diags.Append(d...)
					state.GroupsExcluded = groupsExcluded
				}
			}
		}
	}

	if rule.Actions != nil && rule.Actions.EntityRisk != nil {
		terminateAllSessions := false
		var workflowId string

		for _, action := range rule.Actions.EntityRisk.Actions {
			if action.Action != nil {
				switch *action.Action {
				case "TERMINATE_ALL_SESSIONS":
					terminateAllSessions = true
				case "RUN_WORKFLOW":
					if workflow, ok := action.AdditionalProperties["workflow"].(map[string]interface{}); ok {
						if id, ok := workflow["id"].(string); ok {
							workflowId = id
						} else if id, ok := workflow["id"].(float64); ok {
							workflowId = fmt.Sprintf("%.0f", id)
						}
					}
				}
			}
		}

		state.TerminateAllSessions = types.BoolValue(terminateAllSessions)
		if workflowId != "" {
			state.WorkflowID = types.StringValue(workflowId)
		}
	}

	return
}

func (r *entityRiskPolicyRuleResource) buildEntityRiskPolicyRule(ctx context.Context, state *entityRiskPolicyRuleResourceModel) (*v6okta.EntityRiskPolicyRule, fwdiag.Diagnostics) {
	var diags fwdiag.Diagnostics

	rule := v6okta.NewEntityRiskPolicyRule()

	name := state.Name.ValueString()
	rule.Name = &name

	ruleType := "ENTITY_RISK"
	rule.Type = &ruleType

	if !state.Priority.IsNull() {
		p := int32(state.Priority.ValueInt64())
		rule.Priority = *v6okta.NewNullableInt32(&p)
	}

	conditions := v6okta.NewEntityRiskPolicyRuleConditions()

	// Risk level condition (required)
	riskLevel := state.RiskLevel.ValueString()
	conditions.EntityRisk = v6okta.NewEntityRiskScorePolicyRuleCondition(riskLevel)

	// People conditions - always set with empty arrays if not specified
	// The SDK requires groups and users fields to be present
	var usersIncluded []string
	var usersExcluded []string
	var groupsIncluded []string
	var groupsExcluded []string

	if !state.UsersIncluded.IsNull() {
		diags.Append(state.UsersIncluded.ElementsAs(ctx, &usersIncluded, false)...)
	}
	if !state.UsersExcluded.IsNull() {
		diags.Append(state.UsersExcluded.ElementsAs(ctx, &usersExcluded, false)...)
	}
	if !state.GroupsIncluded.IsNull() {
		diags.Append(state.GroupsIncluded.ElementsAs(ctx, &groupsIncluded, false)...)
	}
	if !state.GroupsExcluded.IsNull() {
		diags.Append(state.GroupsExcluded.ElementsAs(ctx, &groupsExcluded, false)...)
	}

	if usersIncluded == nil {
		usersIncluded = []string{}
	}
	if usersExcluded == nil {
		usersExcluded = []string{}
	}
	if groupsIncluded == nil {
		groupsIncluded = []string{}
	}
	if groupsExcluded == nil {
		groupsExcluded = []string{}
	}

	groupCondition := v6okta.NewGroupCondition()
	groupCondition.SetExclude(groupsExcluded)
	groupCondition.SetInclude(groupsIncluded)

	userCondition := v6okta.NewUserCondition()
	userCondition.SetExclude(usersExcluded)
	userCondition.SetInclude(usersIncluded)

	peopleCondition := v6okta.NewPolicyPeopleCondition()
	peopleCondition.SetGroups(*groupCondition)
	peopleCondition.SetUsers(*userCondition)
	conditions.People = peopleCondition

	rule.Conditions = conditions

	actions := v6okta.NewEntityRiskPolicyRuleAllOfActions()
	entityRiskActions := v6okta.NewEntityRiskPolicyRuleAllOfActionsEntityRisk()
	var actionsList []v6okta.EntityRiskPolicyRuleActionsObject

	if state.TerminateAllSessions.ValueBool() {
		action := v6okta.NewEntityRiskPolicyRuleActionsObject()
		actionType := "TERMINATE_ALL_SESSIONS"
		action.Action = &actionType
		actionsList = append(actionsList, *action)
	}

	if !state.WorkflowID.IsNull() && state.WorkflowID.ValueString() != "" {
		action := v6okta.NewEntityRiskPolicyRuleActionsObject()
		actionType := "RUN_WORKFLOW"
		action.Action = &actionType
		action.AdditionalProperties = map[string]interface{}{
			"workflow": map[string]interface{}{
				"id": state.WorkflowID.ValueString(),
			},
		}
		actionsList = append(actionsList, *action)
	}

	entityRiskActions.Actions = actionsList
	actions.EntityRisk = entityRiskActions
	rule.Actions = actions

	return rule, diags
}

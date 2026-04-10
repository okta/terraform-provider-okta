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
	"github.com/okta/terraform-provider-okta/okta/utils"

	fwdiag "github.com/hashicorp/terraform-plugin-framework/diag"
	v6okta "github.com/okta/okta-sdk-golang/v6/okta"
)

var (
	_ resource.Resource                = &postAuthSessionPolicyRuleResource{}
	_ resource.ResourceWithConfigure   = &postAuthSessionPolicyRuleResource{}
	_ resource.ResourceWithImportState = &postAuthSessionPolicyRuleResource{}
)

func newPostAuthSessionPolicyRuleResource() resource.Resource {
	return &postAuthSessionPolicyRuleResource{}
}

type postAuthSessionPolicyRuleResource struct {
	*config.Config
}

type postAuthSessionPolicyRuleResourceModel struct {
	ID               types.String `tfsdk:"id"`
	PolicyID         types.String `tfsdk:"policy_id"`
	Name             types.String `tfsdk:"name"`
	Status           types.String `tfsdk:"status"`
	UsersExcluded    types.Set    `tfsdk:"users_excluded"`
	GroupsIncluded   types.Set    `tfsdk:"groups_included"`
	GroupsExcluded   types.Set    `tfsdk:"groups_excluded"`
	TerminateSession types.Bool   `tfsdk:"terminate_session"`
	WorkflowID       types.String `tfsdk:"workflow_id"`
}

func (r *postAuthSessionPolicyRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_post_auth_session_policy_rule"
}

func (r *postAuthSessionPolicyRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages the Post Auth Session Policy Rule. The Post Auth Session Policy has exactly one modifiable rule (non-default). This resource allows you to configure that rule. Note: The rule cannot be created or deleted, only modified.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ID of the policy rule.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"policy_id": schema.StringAttribute{
				Description: "ID of the Post Auth Session Policy. Use the `okta_post_auth_session_policy` data source to get this ID.",
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
			"terminate_session": schema.BoolAttribute{
				Description: "When true, terminates the user's session when a policy failure is detected.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"workflow_id": schema.StringAttribute{
				Description: "ID of the Okta Workflow to run when a policy failure is detected.",
				Optional:    true,
			},
		},
	}
}

func (r *postAuthSessionPolicyRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *postAuthSessionPolicyRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	resp.Diagnostics.AddError(
		"Post Auth Session Policy Rule cannot be created",
		"The Post Auth Session Policy Rule already exists and cannot be created, only imported and updated. "+
			"Use 'terraform import' with the policy ID and rule ID to import the existing rule. "+
			"You can find the rule ID via the Okta Admin Console or API.",
	)
}

func (r *postAuthSessionPolicyRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state postAuthSessionPolicyRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.readPostAuthSessionPolicyRule(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID.IsNull() {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *postAuthSessionPolicyRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state postAuthSessionPolicyRuleResourceModel
	var priorState postAuthSessionPolicyRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &priorState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyId := state.PolicyID.ValueString()
	ruleId := state.ID.ValueString()

	r.Logger.Info("updating post auth session policy rule", "policy_id", policyId, "rule_id", ruleId)

	rule, diags := r.buildPostAuthSessionPolicyRule(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	rule.Id = &ruleId

	ruleRequest := v6okta.ListPolicyRules200ResponseInner{
		PostAuthSessionPolicyRule: rule,
	}

	_, _, err := r.OktaIDaaSClient.OktaSDKClientV6().PolicyAPI.ReplacePolicyRule(ctx, policyId, ruleId).PolicyRule(ruleRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to update post auth session policy rule",
			utils.ErrorDetail_V6(err),
		)
		return
	}

	if !state.Status.Equal(priorState.Status) {
		status := state.Status.ValueString()
		if status == "ACTIVE" {
			_, err := r.OktaIDaaSClient.OktaSDKClientV6().PolicyAPI.ActivatePolicyRule(ctx, policyId, ruleId).Execute()
			if err != nil {
				resp.Diagnostics.AddError(
					"failed to activate post auth session policy rule",
					utils.ErrorDetail_V6(err),
				)
				return
			}
		} else {
			_, err := r.OktaIDaaSClient.OktaSDKClientV6().PolicyAPI.DeactivatePolicyRule(ctx, policyId, ruleId).Execute()
			if err != nil {
				resp.Diagnostics.AddError(
					"failed to deactivate post auth session policy rule",
					utils.ErrorDetail_V6(err),
				)
				return
			}
		}
	}

	resp.Diagnostics.Append(r.readPostAuthSessionPolicyRule(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *postAuthSessionPolicyRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state postAuthSessionPolicyRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyId := state.PolicyID.ValueString()
	ruleId := state.ID.ValueString()

	// The Post Auth Session Policy rule cannot be deleted, only modified.
	// On destroy, we just remove it from Terraform state.
	// The rule will remain in Okta with its current configuration.
	r.Logger.Info("removing post auth session policy rule from state (rule cannot be deleted)", "policy_id", policyId, "rule_id", ruleId)
}

func (r *postAuthSessionPolicyRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *postAuthSessionPolicyRuleResource) readPostAuthSessionPolicyRule(ctx context.Context, state *postAuthSessionPolicyRuleResourceModel) (diags fwdiag.Diagnostics) {
	policyId := state.PolicyID.ValueString()
	ruleId := state.ID.ValueString()

	r.Logger.Info("reading post auth session policy rule", "policy_id", policyId, "rule_id", ruleId)

	ruleResp, _, err := r.OktaIDaaSClient.OktaSDKClientV6().PolicyAPI.GetPolicyRule(ctx, policyId, ruleId).Execute()
	if err != nil {
		diags.AddError(
			"failed to get post auth session policy rule",
			utils.ErrorDetail_V6(err),
		)
		return
	}

	if ruleResp.PostAuthSessionPolicyRule == nil {
		state.ID = types.StringNull()
		return
	}

	rule := ruleResp.PostAuthSessionPolicyRule

	// Check if this is the default rule (priority 99) - cannot be managed
	if rule.Priority.IsSet() {
		priority := rule.Priority.Get()
		if priority != nil && *priority == 99 {
			diags.AddError(
				"Cannot manage default policy rule",
				"The default Post Auth Session Policy rule (priority 99) cannot be imported or modified. "+
					"Please import the non-default rule instead. Use the error from 'terraform apply' to get the correct rule ID.",
			)
			return
		}
	}

	state.Name = types.StringPointerValue(rule.Name)
	if rule.Status != nil {
		state.Status = types.StringValue(string(*rule.Status))
	}

	if rule.Conditions != nil {
		if rule.Conditions.People != nil {
			if rule.Conditions.People.Users != nil && len(rule.Conditions.People.Users.Exclude) > 0 {
				usersExcluded, d := types.SetValueFrom(ctx, types.StringType, rule.Conditions.People.Users.Exclude)
				diags.Append(d...)
				state.UsersExcluded = usersExcluded
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

	if rule.Actions != nil && rule.Actions.PostAuthSession != nil {
		terminateSession := false
		var workflowId string

		for _, action := range rule.Actions.PostAuthSession.FailureActions {
			if action.Action != nil {
				switch *action.Action {
				case "TERMINATE_SESSION":
					terminateSession = true
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

		state.TerminateSession = types.BoolValue(terminateSession)
		if workflowId != "" {
			state.WorkflowID = types.StringValue(workflowId)
		}
	}

	return
}

func (r *postAuthSessionPolicyRuleResource) buildPostAuthSessionPolicyRule(ctx context.Context, state *postAuthSessionPolicyRuleResourceModel) (*v6okta.PostAuthSessionPolicyRule, fwdiag.Diagnostics) {
	var diags fwdiag.Diagnostics

	rule := v6okta.NewPostAuthSessionPolicyRule()

	if !state.Name.IsNull() && state.Name.ValueString() != "" {
		name := state.Name.ValueString()
		rule.Name = &name
	}

	ruleType := "POST_AUTH_SESSION"
	rule.Type = &ruleType

	conditions := v6okta.NewPostAuthSessionPolicyRuleAllOfConditions()

	var usersExcluded []string
	var groupsIncluded []string
	var groupsExcluded []string

	if !state.UsersExcluded.IsNull() {
		diags.Append(state.UsersExcluded.ElementsAs(ctx, &usersExcluded, false)...)
	}
	if !state.GroupsIncluded.IsNull() {
		diags.Append(state.GroupsIncluded.ElementsAs(ctx, &groupsIncluded, false)...)
	}
	if !state.GroupsExcluded.IsNull() {
		diags.Append(state.GroupsExcluded.ElementsAs(ctx, &groupsExcluded, false)...)
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
	userCondition.SetInclude([]string{})

	peopleCondition := v6okta.NewPolicyPeopleCondition()
	peopleCondition.SetGroups(*groupCondition)
	peopleCondition.SetUsers(*userCondition)
	conditions.People = peopleCondition

	rule.Conditions = conditions

	actions := v6okta.NewPostAuthSessionPolicyRuleAllOfActions()
	postAuthSession := v6okta.NewPostAuthSessionPolicyRuleAllOfActionsPostAuthSession()
	var failureActions []v6okta.PostAuthSessionFailureActionsObject

	if state.TerminateSession.ValueBool() {
		action := v6okta.NewPostAuthSessionFailureActionsObject()
		actionType := "TERMINATE_SESSION"
		action.Action = &actionType
		failureActions = append(failureActions, *action)
	}

	if !state.WorkflowID.IsNull() && state.WorkflowID.ValueString() != "" {
		action := v6okta.NewPostAuthSessionFailureActionsObject()
		actionType := "RUN_WORKFLOW"
		action.Action = &actionType
		action.AdditionalProperties = map[string]interface{}{
			"workflow": map[string]interface{}{
				"id": state.WorkflowID.ValueString(),
			},
		}
		failureActions = append(failureActions, *action)
	}

	postAuthSession.FailureActions = failureActions
	actions.PostAuthSession = postAuthSession
	rule.Actions = actions

	return rule, diags
}

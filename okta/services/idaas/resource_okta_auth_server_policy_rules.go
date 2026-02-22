package idaas

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
)

// authServerPolicyRulesResource manages all rules of a single Authorization Server
// policy as one Terraform resource, coordinating priority ordering automatically.
type authServerPolicyRulesResource struct {
	*config.Config
}

func newAuthServerPolicyRulesResource() resource.Resource {
	return &authServerPolicyRulesResource{}
}

type authServerPolicyRulesModel struct {
	ID           types.String                `tfsdk:"id"`
	AuthServerID types.String                `tfsdk:"auth_server_id"`
	PolicyID     types.String                `tfsdk:"policy_id"`
	Rules        []authServerPolicyRuleModel `tfsdk:"rule"`
}

type authServerPolicyRuleModel struct {
	ID                          types.String `tfsdk:"id"`
	Name                        types.String `tfsdk:"name"`
	Priority                    types.Int64  `tfsdk:"priority"`
	Status                      types.String `tfsdk:"status"`
	System                      types.Bool   `tfsdk:"system"`
	GrantTypeWhitelist          types.Set    `tfsdk:"grant_type_whitelist"`
	ScopeWhitelist              types.Set    `tfsdk:"scope_whitelist"`
	GroupWhitelist              types.Set    `tfsdk:"group_whitelist"`
	GroupBlacklist              types.Set    `tfsdk:"group_blacklist"`
	UserWhitelist               types.Set    `tfsdk:"user_whitelist"`
	UserBlacklist               types.Set    `tfsdk:"user_blacklist"`
	AccessTokenLifetimeMinutes  types.Int64  `tfsdk:"access_token_lifetime_minutes"`
	RefreshTokenLifetimeMinutes types.Int64  `tfsdk:"refresh_token_lifetime_minutes"`
	RefreshTokenWindowMinutes   types.Int64  `tfsdk:"refresh_token_window_minutes"`
	InlineHookID                types.String `tfsdk:"inline_hook_id"`
}

func (r *authServerPolicyRulesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_auth_server_policy_rules"
}

func (r *authServerPolicyRulesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages all Authorization Server Policy Rules for a single policy as one resource, " +
			"automatically coordinating priority ordering and resolving conflicts.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of this resource (same as policy_id).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"auth_server_id": schema.StringAttribute{
				Required:    true,
				Description: "Authorization Server ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"policy_id": schema.StringAttribute{
				Required:    true,
				Description: "Authorization Server Policy ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"rule": schema.ListNestedBlock{
				Description: "Policy rules in priority order (lowest number = highest priority).",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Optional:    true,
							Description: "Rule ID.",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"name": schema.StringAttribute{
							Required:    true,
							Description: "Rule name.",
						},
						"priority": schema.Int64Attribute{
							Required:    true,
							Description: "Rule priority. Lower numbers are evaluated first.",
							PlanModifiers: []planmodifier.Int64{
								int64planmodifier.UseStateForUnknown(),
							},
						},
						"status": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Rule status: ACTIVE or INACTIVE. Default: ACTIVE.",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"system": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether this is a system (default) rule.",
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
						"grant_type_whitelist": schema.SetAttribute{
							Required:    true,
							ElementType: types.StringType,
							Description: "Accepted grant type values.",
						},
						"scope_whitelist": schema.SetAttribute{
							Optional:    true,
							ElementType: types.StringType,
							Description: "Scopes allowed for this policy rule.",
						},
						"group_whitelist": schema.SetAttribute{
							Optional:    true,
							ElementType: types.StringType,
							Description: "Groups whose users are included.",
						},
						"group_blacklist": schema.SetAttribute{
							Optional:    true,
							ElementType: types.StringType,
							Description: "Groups whose users are excluded.",
						},
						"user_whitelist": schema.SetAttribute{
							Optional:    true,
							ElementType: types.StringType,
							Description: "Users to include.",
						},
						"user_blacklist": schema.SetAttribute{
							Optional:    true,
							ElementType: types.StringType,
							Description: "Users to exclude.",
						},
						"access_token_lifetime_minutes": schema.Int64Attribute{
							Optional:    true,
							Computed:    true,
							Description: "Lifetime of access token in minutes (5–1440). Default: 60.",
						},
						"refresh_token_lifetime_minutes": schema.Int64Attribute{
							Optional:    true,
							Computed:    true,
							Description: "Lifetime of refresh token in minutes.",
						},
						"refresh_token_window_minutes": schema.Int64Attribute{
							Optional:    true,
							Computed:    true,
							Description: "Window for refresh token use in minutes (5–2628000). Default: 10080.",
						},
						"inline_hook_id": schema.StringAttribute{
							Optional:    true,
							Description: "ID of the inline hook to trigger.",
						},
					},
				},
			},
		},
	}
}

func (r *authServerPolicyRulesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *authServerPolicyRulesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan authServerPolicyRulesModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	authServerID := plan.AuthServerID.ValueString()
	policyID := plan.PolicyID.ValueString()

	// Pre-fetch existing rules to recover from a previous interrupted apply.
	existingByName, diags := r.authServerFetchExistingByName(ctx, authServerID, policyID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	sorted := authServerRulesSortByPriority(plan.Rules)
	created := make([]authServerPolicyRuleModel, 0, len(sorted))
	for _, rule := range sorted {
		if existingID, found := existingByName[rule.Name.ValueString()]; found && rule.ID.IsNull() {
			rule.ID = types.StringValue(existingID)
		}
		result, diags := r.authServerCreateOrAdoptRule(ctx, authServerID, policyID, rule)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		created = append(created, result)
	}

	plan.Rules = authServerReorderToMatchPlan(created, plan.Rules)
	plan.ID = plan.PolicyID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *authServerPolicyRulesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state authServerPolicyRulesModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	authServerID := state.AuthServerID.ValueString()
	policyID := state.PolicyID.ValueString()
	client := r.OktaIDaaSClient.OktaSDKClientV2()

	updated := make([]authServerPolicyRuleModel, 0, len(state.Rules))
	for _, rule := range state.Rules {
		if rule.ID.IsNull() || rule.ID.IsUnknown() {
			continue
		}
		apiRule, apiResp, err := client.AuthorizationServer.GetAuthorizationServerPolicyRule(
			ctx, authServerID, policyID, rule.ID.ValueString())
		if err != nil {
			if apiResp != nil && apiResp.StatusCode == http.StatusNotFound {
				continue
			}
			resp.Diagnostics.AddError("Error reading auth server policy rule",
				fmt.Sprintf("Could not read rule '%s': %s", rule.ID.ValueString(), err.Error()))
			return
		}
		updated = append(updated, authServerAPIToModel(ctx, apiRule))
	}

	state.Rules = updated
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *authServerPolicyRulesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan authServerPolicyRulesModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	authServerID := plan.AuthServerID.ValueString()
	policyID := plan.PolicyID.ValueString()
	client := r.OktaIDaaSClient.OktaSDKClientV2()

	stateByName := make(map[string]authServerPolicyRuleModel, len(state.Rules))
	for _, r := range state.Rules {
		stateByName[r.Name.ValueString()] = r
	}
	plannedNames := make(map[string]bool, len(plan.Rules))
	for _, r := range plan.Rules {
		plannedNames[r.Name.ValueString()] = true
	}

	// Delete removed rules.
	for _, stateRule := range state.Rules {
		if stateRule.System.ValueBool() || stateRule.ID.IsNull() {
			continue
		}
		if !plannedNames[stateRule.Name.ValueString()] {
			if err := r.authServerDeleteRule(ctx, authServerID, policyID, stateRule.ID.ValueString()); err != nil {
				resp.Diagnostics.AddError("Error deleting auth server policy rule",
					fmt.Sprintf("Could not delete rule '%s': %s", stateRule.Name.ValueString(), err.Error()))
				return
			}
		}
	}

	sorted := authServerRulesSortByPriority(plan.Rules)
	updated := make([]authServerPolicyRuleModel, 0, len(sorted))
	for _, planRule := range sorted {
		existingRule, exists := stateByName[planRule.Name.ValueString()]
		if exists && !existingRule.ID.IsNull() {
			apiRule := authServerModelToAPI(ctx, planRule)
			result, _, err := client.AuthorizationServer.UpdateAuthorizationServerPolicyRule(
				ctx, authServerID, policyID, existingRule.ID.ValueString(), apiRule)
			if err != nil {
				resp.Diagnostics.AddError("Error updating auth server policy rule",
					fmt.Sprintf("Could not update rule '%s': %s", planRule.Name.ValueString(), err.Error()))
				return
			}
			updated = append(updated, authServerAPIToModel(ctx, result))
		} else {
			result, diags := r.authServerCreateOrAdoptRule(ctx, authServerID, policyID, planRule)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			updated = append(updated, result)
		}
	}

	plan.Rules = authServerReorderToMatchPlan(updated, plan.Rules)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *authServerPolicyRulesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state authServerPolicyRulesModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	authServerID := state.AuthServerID.ValueString()
	policyID := state.PolicyID.ValueString()

	for _, rule := range state.Rules {
		if rule.System.ValueBool() || rule.ID.IsNull() {
			continue
		}
		if err := r.authServerDeleteRule(ctx, authServerID, policyID, rule.ID.ValueString()); err != nil {
			resp.Diagnostics.AddError("Error deleting auth server policy rule",
				fmt.Sprintf("Could not delete rule '%s': %s", rule.Name.ValueString(), err.Error()))
			return
		}
	}
}

func (r *authServerPolicyRulesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import format: {authServerID}/{policyID}
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid import ID",
			"Expected format: {auth_server_id}/{policy_id}")
		return
	}
	authServerID, policyID := parts[0], parts[1]
	client := r.OktaIDaaSClient.OktaSDKClientV2()

	apiRules, apiResp, err := client.AuthorizationServer.ListAuthorizationServerPolicyRules(ctx, authServerID, policyID)
	if err != nil {
		resp.Diagnostics.AddError("Error importing auth server policy rules",
			fmt.Sprintf("Could not list rules: %s", err.Error()))
		return
	}
	if apiResp != nil && apiResp.StatusCode == http.StatusNotFound {
		resp.Diagnostics.AddError("Policy not found",
			fmt.Sprintf("Policy '%s' was not found on auth server '%s'", policyID, authServerID))
		return
	}

	importedRules := make([]authServerPolicyRuleModel, 0, len(apiRules))
	for _, apiRule := range apiRules {
		if apiRule.System != nil && *apiRule.System {
			continue
		}
		importedRules = append(importedRules, authServerAPIToModel(ctx, apiRule))
	}

	state := authServerPolicyRulesModel{
		ID:           types.StringValue(policyID),
		AuthServerID: types.StringValue(authServerID),
		PolicyID:     types.StringValue(policyID),
		Rules:        importedRules,
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// --- helpers ---

// authServerDeleteRule deletes a rule, treating 404 as success (already deleted).
// This makes Delete and Update idempotent on re-run after an interrupted apply.
func (r *authServerPolicyRulesResource) authServerDeleteRule(ctx context.Context, authServerID, policyID, ruleID string) error {
	client := r.OktaIDaaSClient.OktaSDKClientV2()
	apiResp, err := client.AuthorizationServer.DeleteAuthorizationServerPolicyRule(ctx, authServerID, policyID, ruleID)
	if err != nil {
		if apiResp != nil && apiResp.StatusCode == http.StatusNotFound {
			return nil // Already deleted — treat as success.
		}
		return err
	}
	return nil
}

// authServerFetchExistingByName lists all non-system authorization server policy rules
// currently in Okta and returns a name→ID map. Used in Create to detect rules left
// behind by an interrupted previous apply so they can be adopted instead of re-created.
func (r *authServerPolicyRulesResource) authServerFetchExistingByName(
	ctx context.Context,
	authServerID, policyID string,
) (map[string]string, diag.Diagnostics) {
	var diags diag.Diagnostics
	client := r.OktaIDaaSClient.OktaSDKClientV2()
	apiRules, _, err := client.AuthorizationServer.ListAuthorizationServerPolicyRules(ctx, authServerID, policyID)
	if err != nil {
		diags.AddError("Error listing existing auth server policy rules",
			fmt.Sprintf("Could not list rules for policy '%s': %s", policyID, err.Error()))
		return nil, diags
	}
	byName := make(map[string]string, len(apiRules))
	for _, rule := range apiRules {
		if rule.System != nil && *rule.System {
			continue
		}
		byName[rule.Name] = rule.Id
	}
	return byName, diags
}

func (r *authServerPolicyRulesResource) authServerCreateOrAdoptRule(
	ctx context.Context,
	authServerID, policyID string,
	rule authServerPolicyRuleModel,
) (authServerPolicyRuleModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	client := r.OktaIDaaSClient.OktaSDKClientV2()
	apiRule := authServerModelToAPI(ctx, rule)

	if !rule.ID.IsNull() && !rule.ID.IsUnknown() && rule.ID.ValueString() != "" {
		result, _, err := client.AuthorizationServer.UpdateAuthorizationServerPolicyRule(
			ctx, authServerID, policyID, rule.ID.ValueString(), apiRule)
		if err != nil {
			diags.AddError("Error adopting auth server policy rule",
				fmt.Sprintf("Could not adopt rule '%s': %s", rule.Name.ValueString(), err.Error()))
			return authServerPolicyRuleModel{}, diags
		}
		return authServerAPIToModel(ctx, result), diags
	}

	result, err := r.authServerCreateWithRetry(ctx, authServerID, policyID, apiRule)
	if err != nil {
		diags.AddError("Error creating auth server policy rule",
			fmt.Sprintf("Could not create rule '%s': %s", rule.Name.ValueString(), err.Error()))
		return authServerPolicyRuleModel{}, diags
	}
	return authServerAPIToModel(ctx, result), diags
}

func (r *authServerPolicyRulesResource) authServerCreateWithRetry(
	ctx context.Context,
	authServerID, policyID string,
	rule sdk.AuthorizationServerPolicyRule,
) (*sdk.AuthorizationServerPolicyRule, error) {
	client := r.OktaIDaaSClient.OktaSDKClientV2()
	return backoff.Retry(ctx, func() (*sdk.AuthorizationServerPolicyRule, error) {
		created, apiResp, err := client.AuthorizationServer.CreateAuthorizationServerPolicyRule(
			ctx, authServerID, policyID, rule)
		if err != nil {
			if apiResp != nil && apiResp.StatusCode == http.StatusConflict {
				return nil, err
			}
			return nil, backoff.Permanent(err)
		}
		return created, nil
	}, backoff.WithMaxElapsedTime(apiRetryTimeout))
}

func authServerModelToAPI(ctx context.Context, rule authServerPolicyRuleModel) sdk.AuthorizationServerPolicyRule {
	var grantTypes, scopes, groupInclude, groupExclude, userInclude, userExclude []string
	_ = rule.GrantTypeWhitelist.ElementsAs(ctx, &grantTypes, false)
	_ = rule.ScopeWhitelist.ElementsAs(ctx, &scopes, false)
	_ = rule.GroupWhitelist.ElementsAs(ctx, &groupInclude, false)
	_ = rule.GroupBlacklist.ElementsAs(ctx, &groupExclude, false)
	_ = rule.UserWhitelist.ElementsAs(ctx, &userInclude, false)
	_ = rule.UserBlacklist.ElementsAs(ctx, &userExclude, false)

	atlm := int(rule.AccessTokenLifetimeMinutes.ValueInt64())
	if atlm == 0 {
		atlm = 60
	}
	rtwm := int(rule.RefreshTokenWindowMinutes.ValueInt64())
	if rtwm == 0 {
		rtwm = 10080
	}

	apiRule := sdk.AuthorizationServerPolicyRule{
		Name:        rule.Name.ValueString(),
		Status:      rule.Status.ValueString(),
		PriorityPtr: utils.Int64Ptr(int(rule.Priority.ValueInt64())),
		Type:        "RESOURCE_ACCESS",
		Actions: &sdk.AuthorizationServerPolicyRuleActions{
			Token: &sdk.TokenAuthorizationServerPolicyRuleAction{
				AccessTokenLifetimeMinutesPtr: utils.Int64Ptr(atlm),
				RefreshTokenWindowMinutesPtr:  utils.Int64Ptr(rtwm),
			},
		},
		Conditions: &sdk.AuthorizationServerPolicyRuleConditions{
			GrantTypes: &sdk.GrantTypePolicyRuleCondition{Include: grantTypes},
			Scopes:     &sdk.OAuth2ScopesMediationPolicyRuleCondition{Include: scopes},
			People: &sdk.PolicyPeopleCondition{
				Groups: &sdk.GroupCondition{Include: groupInclude, Exclude: groupExclude},
				Users:  &sdk.UserCondition{Include: userInclude, Exclude: userExclude},
			},
		},
	}

	if !rule.RefreshTokenLifetimeMinutes.IsNull() && rule.RefreshTokenLifetimeMinutes.ValueInt64() > 0 {
		apiRule.Actions.Token.RefreshTokenLifetimeMinutesPtr = utils.Int64Ptr(int(rule.RefreshTokenLifetimeMinutes.ValueInt64()))
	}

	if !rule.InlineHookID.IsNull() && rule.InlineHookID.ValueString() != "" {
		apiRule.Actions.Token.InlineHook = &sdk.TokenAuthorizationServerPolicyRuleActionInlineHook{
			Id: rule.InlineHookID.ValueString(),
		}
	}
	return apiRule
}

func authServerAPIToModel(ctx context.Context, apiRule *sdk.AuthorizationServerPolicyRule) authServerPolicyRuleModel {
	emptySet, _ := types.SetValueFrom(ctx, types.StringType, []string{})

	model := authServerPolicyRuleModel{
		ID:                         types.StringValue(apiRule.Id),
		Name:                       types.StringValue(apiRule.Name),
		Status:                     types.StringValue(apiRule.Status),
		System:                     types.BoolValue(utils.BoolFromBoolPtr(apiRule.System)),
		GrantTypeWhitelist:         emptySet,
		ScopeWhitelist:             emptySet,
		GroupWhitelist:             emptySet,
		GroupBlacklist:             emptySet,
		UserWhitelist:              emptySet,
		UserBlacklist:              emptySet,
		AccessTokenLifetimeMinutes: types.Int64Value(60),
		RefreshTokenWindowMinutes:  types.Int64Value(10080),
	}

	if apiRule.PriorityPtr != nil {
		model.Priority = types.Int64Value(*apiRule.PriorityPtr)
	}
	if apiRule.Actions != nil && apiRule.Actions.Token != nil {
		t := apiRule.Actions.Token
		if t.AccessTokenLifetimeMinutesPtr != nil {
			model.AccessTokenLifetimeMinutes = types.Int64Value(*t.AccessTokenLifetimeMinutesPtr)
		}
		if t.RefreshTokenLifetimeMinutesPtr != nil {
			model.RefreshTokenLifetimeMinutes = types.Int64Value(*t.RefreshTokenLifetimeMinutesPtr)
		}
		if t.RefreshTokenWindowMinutesPtr != nil {
			model.RefreshTokenWindowMinutes = types.Int64Value(*t.RefreshTokenWindowMinutesPtr)
		}
		if t.InlineHook != nil {
			model.InlineHookID = types.StringValue(t.InlineHook.Id)
		}
	}
	if apiRule.Conditions != nil {
		if apiRule.Conditions.GrantTypes != nil {
			v, _ := types.SetValueFrom(ctx, types.StringType, apiRule.Conditions.GrantTypes.Include)
			model.GrantTypeWhitelist = v
		}
		if apiRule.Conditions.Scopes != nil {
			v, _ := types.SetValueFrom(ctx, types.StringType, apiRule.Conditions.Scopes.Include)
			model.ScopeWhitelist = v
		}
		if apiRule.Conditions.People != nil {
			if apiRule.Conditions.People.Groups != nil {
				inc, _ := types.SetValueFrom(ctx, types.StringType, apiRule.Conditions.People.Groups.Include)
				exc, _ := types.SetValueFrom(ctx, types.StringType, apiRule.Conditions.People.Groups.Exclude)
				model.GroupWhitelist = inc
				model.GroupBlacklist = exc
			}
			if apiRule.Conditions.People.Users != nil {
				inc, _ := types.SetValueFrom(ctx, types.StringType, apiRule.Conditions.People.Users.Include)
				exc, _ := types.SetValueFrom(ctx, types.StringType, apiRule.Conditions.People.Users.Exclude)
				model.UserWhitelist = inc
				model.UserBlacklist = exc
			}
		}
	}
	return model
}

func authServerRulesSortByPriority(rules []authServerPolicyRuleModel) []authServerPolicyRuleModel {
	sorted := make([]authServerPolicyRuleModel, len(rules))
	copy(sorted, rules)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Priority.ValueInt64() < sorted[j].Priority.ValueInt64()
	})
	return sorted
}

func authServerReorderToMatchPlan(processed, plan []authServerPolicyRuleModel) []authServerPolicyRuleModel {
	byName := make(map[string]authServerPolicyRuleModel, len(processed))
	for _, r := range processed {
		byName[r.Name.ValueString()] = r
	}
	result := make([]authServerPolicyRuleModel, 0, len(plan))
	for _, planRule := range plan {
		if r, ok := byName[planRule.Name.ValueString()]; ok {
			result = append(result, r)
		}
	}
	return result
}

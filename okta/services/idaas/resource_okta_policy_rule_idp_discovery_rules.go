package idaas

import (
	"context"
	"fmt"
	"net/http"
	"sort"

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
	"github.com/okta/terraform-provider-okta/sdk"
)

// idpDiscoveryPolicyRulesResource manages all rules of a single IDP Discovery policy
// as one Terraform resource, coordinating priority ordering automatically.
type idpDiscoveryPolicyRulesResource struct {
	*config.Config
}

func newIdpDiscoveryPolicyRulesResource() resource.Resource {
	return &idpDiscoveryPolicyRulesResource{}
}

type idpDiscoveryPolicyRulesModel struct {
	ID       types.String                  `tfsdk:"id"`
	PolicyID types.String                  `tfsdk:"policy_id"`
	Rules    []idpDiscoveryPolicyRuleModel `tfsdk:"rule"`
}

type idpDiscoveryPolicyRuleModel struct {
	ID       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Priority types.Int64  `tfsdk:"priority"`
	Status   types.String `tfsdk:"status"`
	System   types.Bool   `tfsdk:"system"`
}

func (r *idpDiscoveryPolicyRulesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_rule_idp_discovery_rules"
}

func (r *idpDiscoveryPolicyRulesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages all IDP Discovery Policy Rules for a single policy as one resource, " +
			"automatically coordinating priority ordering and resolving conflicts.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of this resource (same as policy_id).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"policy_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the IDP Discovery policy.",
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
							Description: "Rule name. Must be unique within the policy.",
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
							Description: "Rule status: ACTIVE or INACTIVE.",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"system": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether this is a system (default catch-all) rule.",
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
					},
				},
			},
		},
	}
}

func (r *idpDiscoveryPolicyRulesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *idpDiscoveryPolicyRulesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan idpDiscoveryPolicyRulesModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyID := plan.PolicyID.ValueString()
	client := r.OktaIDaaSClient.OktaSDKSupplementClient()

	// Pre-fetch existing rules to recover from a previous interrupted apply.
	existingByName, diags := r.idpDiscoveryFetchExistingByName(ctx, client, policyID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	sorted := idpDiscoverySortByPriority(plan.Rules)
	created := make([]idpDiscoveryPolicyRuleModel, 0, len(sorted))
	for _, rule := range sorted {
		if existingID, found := existingByName[rule.Name.ValueString()]; found && rule.ID.IsNull() {
			rule.ID = types.StringValue(existingID)
		}
		result, diags := r.idpDiscoveryCreateOrAdoptRule(ctx, client, policyID, rule)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		created = append(created, result)
	}

	plan.Rules = idpDiscoveryReorderToMatchPlan(created, plan.Rules)
	plan.ID = plan.PolicyID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *idpDiscoveryPolicyRulesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state idpDiscoveryPolicyRulesModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyID := state.PolicyID.ValueString()
	client := r.OktaIDaaSClient.OktaSDKSupplementClient()

	updated := make([]idpDiscoveryPolicyRuleModel, 0, len(state.Rules))
	for _, rule := range state.Rules {
		if rule.ID.IsNull() || rule.ID.IsUnknown() {
			continue
		}
		apiRule, resp2, err := client.GetIdpDiscoveryRule(ctx, policyID, rule.ID.ValueString())
		if err != nil {
			if resp2 != nil && resp2.StatusCode == http.StatusNotFound {
				continue
			}
			resp.Diagnostics.AddError("Error reading IDP discovery policy rule",
				fmt.Sprintf("Could not read rule '%s': %s", rule.ID.ValueString(), err.Error()))
			return
		}
		updated = append(updated, idpDiscoveryAPIToModel(apiRule))
	}

	state.Rules = updated
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *idpDiscoveryPolicyRulesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan idpDiscoveryPolicyRulesModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyID := plan.PolicyID.ValueString()
	client := r.OktaIDaaSClient.OktaSDKSupplementClient()

	// Build name→ID map from state for matching.
	stateByName := make(map[string]idpDiscoveryPolicyRuleModel, len(state.Rules))
	for _, r := range state.Rules {
		stateByName[r.Name.ValueString()] = r
	}
	plannedNames := make(map[string]bool, len(plan.Rules))
	for _, r := range plan.Rules {
		plannedNames[r.Name.ValueString()] = true
	}

	// Delete rules removed from plan (skip system rules).
	for _, stateRule := range state.Rules {
		if stateRule.System.ValueBool() || stateRule.ID.IsNull() {
			continue
		}
		if !plannedNames[stateRule.Name.ValueString()] {
			if err := r.idpDiscoveryDeleteRule(ctx, client, policyID, stateRule.ID.ValueString()); err != nil {
				resp.Diagnostics.AddError("Error deleting IDP discovery policy rule",
					fmt.Sprintf("Could not delete rule '%s': %s", stateRule.Name.ValueString(), err.Error()))
				return
			}
		}
	}

	sorted := idpDiscoverySortByPriority(plan.Rules)
	updated := make([]idpDiscoveryPolicyRuleModel, 0, len(sorted))
	for _, planRule := range sorted {
		existingRule, exists := stateByName[planRule.Name.ValueString()]
		if exists && !existingRule.ID.IsNull() {
			apiRule := idpDiscoveryModelToAPI(planRule, existingRule.ID.ValueString())
			result, _, err := client.UpdateIdpDiscoveryRule(ctx, policyID, existingRule.ID.ValueString(), apiRule, nil)
			if err != nil {
				resp.Diagnostics.AddError("Error updating IDP discovery policy rule",
					fmt.Sprintf("Could not update rule '%s': %s", planRule.Name.ValueString(), err.Error()))
				return
			}
			updated = append(updated, idpDiscoveryAPIToModel(result))
		} else {
			result, diags := r.idpDiscoveryCreateOrAdoptRule(ctx, client, policyID, planRule)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			updated = append(updated, result)
		}
	}

	plan.Rules = idpDiscoveryReorderToMatchPlan(updated, plan.Rules)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *idpDiscoveryPolicyRulesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state idpDiscoveryPolicyRulesModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyID := state.PolicyID.ValueString()
	client := r.OktaIDaaSClient.OktaSDKSupplementClient()

	for _, rule := range state.Rules {
		if rule.System.ValueBool() || rule.ID.IsNull() {
			continue
		}
		if err := r.idpDiscoveryDeleteRule(ctx, client, policyID, rule.ID.ValueString()); err != nil {
			resp.Diagnostics.AddError("Error deleting IDP discovery policy rule",
				fmt.Sprintf("Could not delete rule '%s': %s", rule.Name.ValueString(), err.Error()))
			return
		}
	}
}

func (r *idpDiscoveryPolicyRulesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	policyID := req.ID
	client := r.OktaIDaaSClient.OktaSDKSupplementClient()

	apiRules, apiResp, err := client.ListIdpDiscoveryRules(ctx, policyID)
	if err != nil {
		resp.Diagnostics.AddError("Error importing IDP discovery policy rules",
			fmt.Sprintf("Could not list rules for policy '%s': %s", policyID, err.Error()))
		return
	}
	if apiResp != nil && apiResp.StatusCode == http.StatusNotFound {
		resp.Diagnostics.AddError("Policy not found",
			fmt.Sprintf("Policy '%s' was not found", policyID))
		return
	}

	importedRules := make([]idpDiscoveryPolicyRuleModel, 0, len(apiRules))
	for _, apiRule := range apiRules {
		if apiRule.System {
			continue
		}
		importedRules = append(importedRules, idpDiscoveryAPIToModel(apiRule))
	}

	state := idpDiscoveryPolicyRulesModel{
		ID:       types.StringValue(policyID),
		PolicyID: types.StringValue(policyID),
		Rules:    importedRules,
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// --- helpers ---

// idpDiscoveryFetchExistingByName lists all non-system IDP discovery rules currently
// in Okta and returns a name→ID map. Used in Create to detect rules left behind
// by an interrupted previous apply so they can be adopted instead of re-created.
func (r *idpDiscoveryPolicyRulesResource) idpDiscoveryFetchExistingByName(
	ctx context.Context,
	client *sdk.APISupplement,
	policyID string,
) (map[string]string, diag.Diagnostics) {
	var diags diag.Diagnostics
	apiRules, _, err := client.ListIdpDiscoveryRules(ctx, policyID)
	if err != nil {
		diags.AddError("Error listing existing IDP discovery rules",
			fmt.Sprintf("Could not list rules for policy '%s': %s", policyID, err.Error()))
		return nil, diags
	}
	byName := make(map[string]string, len(apiRules))
	for _, rule := range apiRules {
		if rule.System {
			continue
		}
		byName[rule.Name] = rule.ID
	}
	return byName, diags
}

func (r *idpDiscoveryPolicyRulesResource) idpDiscoveryCreateOrAdoptRule(
	ctx context.Context,
	client *sdk.APISupplement,
	policyID string,
	rule idpDiscoveryPolicyRuleModel,
) (idpDiscoveryPolicyRuleModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Adopt existing rule by ID if provided.
	if !rule.ID.IsNull() && !rule.ID.IsUnknown() && rule.ID.ValueString() != "" {
		apiRule := idpDiscoveryModelToAPI(rule, rule.ID.ValueString())
		result, _, err := client.UpdateIdpDiscoveryRule(ctx, policyID, rule.ID.ValueString(), apiRule, nil)
		if err != nil {
			diags.AddError("Error adopting IDP discovery policy rule",
				fmt.Sprintf("Could not adopt rule '%s': %s", rule.Name.ValueString(), err.Error()))
			return idpDiscoveryPolicyRuleModel{}, diags
		}
		return idpDiscoveryAPIToModel(result), diags
	}

	apiRule := idpDiscoveryModelToAPI(rule, "")
	result, err := r.idpDiscoveryCreateWithRetry(ctx, client, policyID, apiRule)
	if err != nil {
		diags.AddError("Error creating IDP discovery policy rule",
			fmt.Sprintf("Could not create rule '%s': %s", rule.Name.ValueString(), err.Error()))
		return idpDiscoveryPolicyRuleModel{}, diags
	}
	return idpDiscoveryAPIToModel(result), diags
}

func (r *idpDiscoveryPolicyRulesResource) idpDiscoveryCreateWithRetry(
	ctx context.Context,
	client *sdk.APISupplement,
	policyID string,
	rule sdk.IdpDiscoveryRule,
) (*sdk.IdpDiscoveryRule, error) {
	return backoff.Retry(ctx, func() (*sdk.IdpDiscoveryRule, error) {
		created, apiResp, err := client.CreateIdpDiscoveryRule(ctx, policyID, rule, nil)
		if err != nil {
			if apiResp != nil && apiResp.StatusCode == http.StatusConflict {
				return nil, err // retry on conflict
			}
			return nil, backoff.Permanent(err)
		}
		return created, nil
	}, backoff.WithMaxElapsedTime(apiRetryTimeout))
}

func (r *idpDiscoveryPolicyRulesResource) idpDiscoveryDeleteRule(
	ctx context.Context,
	client *sdk.APISupplement,
	policyID, ruleID string,
) error {
	_, err := backoff.Retry(ctx, func() (struct{}, error) {
		apiResp, err := client.DeleteIdpDiscoveryRule(ctx, policyID, ruleID)
		if err != nil {
			if apiResp != nil && apiResp.StatusCode == http.StatusConflict {
				return struct{}{}, err
			}
			if apiResp != nil && apiResp.StatusCode == http.StatusNotFound {
				return struct{}{}, nil
			}
			return struct{}{}, backoff.Permanent(err)
		}
		return struct{}{}, nil
	}, backoff.WithMaxElapsedTime(apiRetryTimeout))
	return err
}

func idpDiscoveryModelToAPI(rule idpDiscoveryPolicyRuleModel, id string) sdk.IdpDiscoveryRule {
	r := sdk.IdpDiscoveryRule{
		ID:       id,
		Name:     rule.Name.ValueString(),
		Priority: int(rule.Priority.ValueInt64()),
		Type:     sdk.IdpDiscoveryType,
		Actions: &sdk.IdpDiscoveryRuleActions{
			IDP: &sdk.IdpDiscoveryRuleIdp{
				Providers: []*sdk.IdpDiscoveryRuleProvider{{Type: "OKTA"}},
			},
		},
		Conditions: &sdk.IdpDiscoveryRuleConditions{
			App:     &sdk.IdpDiscoveryRuleApp{},
			Network: &sdk.IdpDiscoveryRuleNetwork{Connection: "ANYWHERE"},
		},
	}
	if !rule.Status.IsNull() && !rule.Status.IsUnknown() {
		r.Status = rule.Status.ValueString()
	}
	return r
}

func idpDiscoveryAPIToModel(apiRule *sdk.IdpDiscoveryRule) idpDiscoveryPolicyRuleModel {
	return idpDiscoveryPolicyRuleModel{
		ID:       types.StringValue(apiRule.ID),
		Name:     types.StringValue(apiRule.Name),
		Priority: types.Int64Value(int64(apiRule.Priority)),
		Status:   types.StringValue(apiRule.Status),
		System:   types.BoolValue(apiRule.System),
	}
}

func idpDiscoverySortByPriority(rules []idpDiscoveryPolicyRuleModel) []idpDiscoveryPolicyRuleModel {
	sorted := make([]idpDiscoveryPolicyRuleModel, len(rules))
	copy(sorted, rules)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Priority.ValueInt64() < sorted[j].Priority.ValueInt64()
	})
	return sorted
}

func idpDiscoveryReorderToMatchPlan(processed, plan []idpDiscoveryPolicyRuleModel) []idpDiscoveryPolicyRuleModel {
	byName := make(map[string]idpDiscoveryPolicyRuleModel, len(processed))
	for _, r := range processed {
		byName[r.Name.ValueString()] = r
	}
	result := make([]idpDiscoveryPolicyRuleModel, 0, len(plan))
	for _, planRule := range plan {
		if r, ok := byName[planRule.Name.ValueString()]; ok {
			result = append(result, r)
		}
	}
	return result
}

package idaas

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
)

const (
	// apiRetryTimeout is the maximum time to wait for API operations with retries.
	apiRetryTimeout = 30 * time.Second
	// ruleTypeAccessPolicy is the Okta API type for access policy rules.
	ruleTypeAccessPolicy = "ACCESS_POLICY"
)

// Validation values for rule attributes.
var (
	validStatuses           = []string{"ACTIVE", "INACTIVE"}
	validNetworkConnections = []string{"ANYWHERE", "ZONE", "ON_NETWORK", "OFF_NETWORK"}
	validAccessTypes        = []string{"ALLOW", "DENY"}
	validFactorModes        = []string{"1FA", "2FA"}
	validRiskScores         = []string{"ANY", "LOW", "MEDIUM", "HIGH"}
	validPlatformTypes      = []string{"ANY", "MOBILE", "DESKTOP"}
	// NOTE: ANY is intentionally excluded. The API silently maps os_type=ANY to
	// OTHER on read, which causes a post-apply inconsistency. Users should set
	// os_type = "OTHER" directly (with or without os_expression).
	validOSTypes = []string{"IOS", "ANDROID", "WINDOWS", "OSX", "MACOS", "CHROMEOS", "OTHER", "LINUX"}
)
var (
	_ resource.Resource                   = &appSignOnPolicyRulesResource{}
	_ resource.ResourceWithConfigure      = &appSignOnPolicyRulesResource{}
	_ resource.ResourceWithImportState    = &appSignOnPolicyRulesResource{}
	_ resource.ResourceWithValidateConfig = &appSignOnPolicyRulesResource{}
)

// NewAppSignOnPolicyRulesResource creates a new instance of the resource.
func NewAppSignOnPolicyRulesResource() resource.Resource {
	return &appSignOnPolicyRulesResource{}
}

// For backward compatibility with existing registration code.
func newAppSignOnPolicyRulesResource() resource.Resource {
	return NewAppSignOnPolicyRulesResource()
}

type appSignOnPolicyRulesResource struct {
	*config.Config
}

// appSignOnPolicyRulesModel represents the Terraform state for this resource.
type appSignOnPolicyRulesModel struct {
	ID       types.String `tfsdk:"id"`
	PolicyID types.String `tfsdk:"policy_id"`
	Rules    types.List   `tfsdk:"rule"`
}

// policyRuleModel represents a single policy rule in the Terraform state.
type policyRuleModel struct {
	ID                        types.String           `tfsdk:"id"`
	Name                      types.String           `tfsdk:"name"`
	System                    types.Bool             `tfsdk:"system"`
	Status                    types.String           `tfsdk:"status"`
	Priority                  types.Int64            `tfsdk:"priority"`
	GroupsIncluded            types.Set              `tfsdk:"groups_included"`
	GroupsExcluded            types.Set              `tfsdk:"groups_excluded"`
	UsersIncluded             types.Set              `tfsdk:"users_included"`
	UsersExcluded             types.Set              `tfsdk:"users_excluded"`
	NetworkConnection         types.String           `tfsdk:"network_connection"`
	NetworkIncludes           types.List             `tfsdk:"network_includes"`
	NetworkExcludes           types.List             `tfsdk:"network_excludes"`
	DeviceIsRegistered        types.Bool             `tfsdk:"device_is_registered"`
	DeviceIsManaged           types.Bool             `tfsdk:"device_is_managed"`
	DeviceAssurancesIncluded  types.Set              `tfsdk:"device_assurances_included"`
	UserTypesIncluded         types.Set              `tfsdk:"user_types_included"`
	UserTypesExcluded         types.Set              `tfsdk:"user_types_excluded"`
	CustomExpression          types.String           `tfsdk:"custom_expression"`
	Access                    types.String           `tfsdk:"access"`
	FactorMode                types.String           `tfsdk:"factor_mode"`
	Type                      types.String           `tfsdk:"type"`
	ReAuthenticationFrequency types.String           `tfsdk:"re_authentication_frequency"`
	InactivityPeriod          types.String           `tfsdk:"inactivity_period"`
	Constraints               types.List             `tfsdk:"constraints"`
	Chains                    types.List             `tfsdk:"chains"`
	RiskScore                 types.String           `tfsdk:"risk_score"`
	PlatformInclude           []platformIncludeModel `tfsdk:"platform_include"`
}

// platformIncludeModel represents platform conditions in the rule.
type platformIncludeModel struct {
	Type         types.String `tfsdk:"type"`
	OsType       types.String `tfsdk:"os_type"`
	OsExpression types.String `tfsdk:"os_expression"`
}

// reauthFrequencyModifier is a plan modifier that suppresses changes to
// re_authentication_frequency when chains contain reauthenticateIn.
type reauthFrequencyModifier struct{}

func (m reauthFrequencyModifier) Description(ctx context.Context) string {
	return "Suppresses re_authentication_frequency changes when chains contain reauthenticateIn"
}

func (m reauthFrequencyModifier) MarkdownDescription(ctx context.Context) string {
	return "Suppresses re_authentication_frequency changes when chains contain reauthenticateIn"
}

func (m reauthFrequencyModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// Get the parent rule object from the plan
	var planRule policyRuleModel
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, req.Path.ParentPath(), &planRule)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if chains contain reauthenticateIn
	if !planRule.Chains.IsNull() && !planRule.Chains.IsUnknown() {
		var chainStrings []string
		planRule.Chains.ElementsAs(ctx, &chainStrings, false)
		for _, chainStr := range chainStrings {
			if strings.Contains(chainStr, "reauthenticateIn") {
				// When chains have reauthenticateIn, the API computes re_authentication_frequency
				// If we have a state value, use it to suppress diff (like DiffSuppressFunc)
				// Otherwise mark as unknown so first apply accepts whatever API returns
				if !req.StateValue.IsNull() && !req.StateValue.IsUnknown() {
					resp.PlanValue = req.StateValue
				} else {
					resp.PlanValue = types.StringUnknown()
				}
				return
			}
		}
	}
}

// ruleIndex provides efficient lookups for rules by name and ID.
type ruleIndex struct {
	byName map[string]policyRuleModel
	byID   map[string]policyRuleModel
}

// newRuleIndex creates a new rule index from a slice of rules.
func newRuleIndex(rules []policyRuleModel) *ruleIndex {
	idx := &ruleIndex{
		byName: make(map[string]policyRuleModel, len(rules)),
		byID:   make(map[string]policyRuleModel, len(rules)),
	}
	for _, rule := range rules {
		idx.byName[rule.Name.ValueString()] = rule
		if !rule.ID.IsNull() && !rule.ID.IsUnknown() {
			idx.byID[rule.ID.ValueString()] = rule
		}
	}
	return idx
}

// findRule looks up a rule by name first, then by ID.
// Name is used as the primary key because it is a Required attribute set
// explicitly by the user in config. The ID is a Computed attribute that
// Terraform injects positionally from state, which can be wrong when the
// state order differs from the config order (e.g. after a rule is added or
// deleted). Falling back to ID handles the rename case where the name has
// changed but the same rule ID is referenced by the plan.
func (idx *ruleIndex) findRule(id, name string) (policyRuleModel, bool) {
	if name != "" {
		if rule, ok := idx.byName[name]; ok {
			return rule, true
		}
	}
	if id != "" {
		if rule, ok := idx.byID[id]; ok {
			return rule, true
		}
	}
	return policyRuleModel{}, false
}

// nameTracker tracks current name-to-ID mappings for conflict detection.
type nameTracker struct {
	nameToID map[string]string
	idToName map[string]string
}

// newNameTracker creates a tracker from existing rules.
func newNameTracker(rules []policyRuleModel) *nameTracker {
	t := &nameTracker{
		nameToID: make(map[string]string, len(rules)),
		idToName: make(map[string]string, len(rules)),
	}
	for _, rule := range rules {
		if rule.ID.IsNull() || rule.ID.IsUnknown() {
			continue
		}
		if rule.Name.IsNull() || rule.Name.IsUnknown() {
			continue
		}
		id := rule.ID.ValueString()
		name := rule.Name.ValueString()
		if name != "" {
			t.nameToID[name] = id
			t.idToName[id] = name
		}
	}
	return t
}

// hasConflict checks if the target name is held by a different rule.
func (t *nameTracker) hasConflict(targetName, ruleID string) (conflictingID string, hasConflict bool) {
	if existingID, ok := t.nameToID[targetName]; ok && existingID != ruleID {
		return existingID, true
	}
	return "", false
}

// updateMapping updates the tracker after a rule rename.
func (t *nameTracker) updateMapping(ruleID, newName string) {
	if oldName, ok := t.idToName[ruleID]; ok {
		delete(t.nameToID, oldName)
	}
	t.idToName[ruleID] = newName
	t.nameToID[newName] = ruleID
}
func (r *appSignOnPolicyRulesResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}
func (r *appSignOnPolicyRulesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_signon_policy_rules"
}
func (r *appSignOnPolicyRulesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = r.buildSchema()
}
func (r *appSignOnPolicyRulesResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data appSignOnPolicyRulesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var rules []policyRuleModel
	resp.Diagnostics.Append(data.Rules.ElementsAs(ctx, &rules, true)...)
	if resp.Diagnostics.HasError() {
		return
	}
	seen := make(map[string]int, len(rules))
	seenPriority := make(map[int64]int, len(rules))
	for i, rule := range rules {
		name := rule.Name.ValueString()
		if name == "" {
			continue
		}
		if prevIdx, exists := seen[name]; exists {
			resp.Diagnostics.AddError(
				"Duplicate rule name",
				fmt.Sprintf(
					"Rule name %q is used by both rule[%d] and rule[%d]. Each rule within a policy must have a unique name.",
					name, prevIdx, i,
				),
			)
		} else {
			seen[name] = i
		}
		if !rule.Priority.IsNull() && !rule.Priority.IsUnknown() {
			p := rule.Priority.ValueInt64()
			if prevIdx, exists := seenPriority[p]; exists {
				resp.Diagnostics.AddError(
					"Duplicate rule priority",
					fmt.Sprintf(
						"Priority %d is used by both rule[%d] and rule[%d]. Each rule must have a unique priority.",
						p, prevIdx, i,
					),
				)
			} else {
				seenPriority[p] = i
			}
		}
		// Warn if conditions are set on a system rule (Catch-all Rule).
		// The API rejects condition modifications on system rules.
		if name == "Catch-all Rule" && r.hasConditionsSet(rule) {
			resp.Diagnostics.AddWarning(
				"Conditions ignored for system rule",
				fmt.Sprintf(
					"rule[%d] %q is a system (Catch-all) rule. Conditions (network, platform, groups, users, "+
						"user_types, device, risk_score, custom_expression) cannot be modified on system rules "+
						"and will be ignored. Only actions (access, factor_mode, type, re_authentication_frequency, "+
						"inactivity_period, constraints) can be configured.",
					i, name,
				),
			)
		}
	}
}

// hasConditionsSet returns true if any condition attributes are set on the rule.
func (r *appSignOnPolicyRulesResource) hasConditionsSet(rule policyRuleModel) bool {
	return (!rule.NetworkConnection.IsNull() && !rule.NetworkConnection.IsUnknown()) ||
		(!rule.NetworkIncludes.IsNull() && !rule.NetworkIncludes.IsUnknown()) ||
		(!rule.NetworkExcludes.IsNull() && !rule.NetworkExcludes.IsUnknown()) ||
		(!rule.GroupsIncluded.IsNull() && !rule.GroupsIncluded.IsUnknown()) ||
		(!rule.GroupsExcluded.IsNull() && !rule.GroupsExcluded.IsUnknown()) ||
		(!rule.UsersIncluded.IsNull() && !rule.UsersIncluded.IsUnknown()) ||
		(!rule.UsersExcluded.IsNull() && !rule.UsersExcluded.IsUnknown()) ||
		(!rule.UserTypesIncluded.IsNull() && !rule.UserTypesIncluded.IsUnknown()) ||
		(!rule.UserTypesExcluded.IsNull() && !rule.UserTypesExcluded.IsUnknown()) ||
		(!rule.DeviceIsRegistered.IsNull() && !rule.DeviceIsRegistered.IsUnknown()) ||
		(!rule.DeviceIsManaged.IsNull() && !rule.DeviceIsManaged.IsUnknown()) ||
		(!rule.DeviceAssurancesIncluded.IsNull() && !rule.DeviceAssurancesIncluded.IsUnknown()) ||
		(!rule.RiskScore.IsNull() && !rule.RiskScore.IsUnknown()) ||
		(!rule.CustomExpression.IsNull() && !rule.CustomExpression.IsUnknown()) ||
		len(rule.PlatformInclude) > 0
}
func (r *appSignOnPolicyRulesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan appSignOnPolicyRulesModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var rules []policyRuleModel
	plan.Rules.ElementsAs(ctx, &rules, false)
	policyID := plan.PolicyID.ValueString()
	client := r.OktaIDaaSClient.OktaSDKSupplementClient()
	// Pre-fetch any rules that already exist in Okta for this policy.
	// This recovers from a previous interrupted apply where some rules were
	// created but state was never written. Without this, re-running apply
	// would hit 409 name conflicts on the already-created rules.
	existingByName, diags := r.fetchExistingRulesByName(ctx, client, policyID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Process rules in priority order to ensure correct ordering in Okta.
	sortedRules := r.sortRulesByPriority(rules)
	createdRules := make([]policyRuleModel, 0, len(sortedRules))
	for _, rule := range sortedRules {
		// If an existing rule was found by name, inject its ID so createOrAdoptRule
		// will update it rather than attempting to create a duplicate.
		// Check both IsNull() and IsUnknown() because during a fresh Create (no prior state),
		// Computed attributes come through as Unknown, not Null.
		if info, found := existingByName[rule.Name.ValueString()]; found {
			if rule.ID.IsNull() || rule.ID.IsUnknown() {
				rule.ID = types.StringValue(info.ID)
			}
		}
		resultRule, diags := r.createOrAdoptRule(ctx, client, policyID, rule)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		createdRules = append(createdRules, resultRule)
	}
	// Reorder to match config order (Terraform expects state order to match plan order).
	reorderedRules := r.reorderRulesToMatchPlan(createdRules, rules)
	plan.ID = plan.PolicyID
	// Marshal back to types.List
	plan.Rules, diags = types.ListValueFrom(ctx, r.policyRuleObjectType(), reorderedRules)
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
func (r *appSignOnPolicyRulesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state appSignOnPolicyRulesModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var rules []policyRuleModel
	resp.Diagnostics.Append(state.Rules.ElementsAs(ctx, &rules, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	policyID := state.PolicyID.ValueString()
	client := r.OktaIDaaSClient.OktaSDKSupplementClient()
	updatedRules := make([]policyRuleModel, 0, len(rules))
	for _, rule := range rules {
		if rule.ID.IsNull() || rule.ID.IsUnknown() {
			continue
		}
		apiRule, err := r.readRuleFromAPI(ctx, client, policyID, rule.ID.ValueString())
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				// Rule was deleted outside Terraform - skip it.
				continue
			}
			resp.Diagnostics.AddError("Error reading app sign-on policy rule",
				fmt.Sprintf("Could not read rule '%s': %s", rule.ID.ValueString(), err.Error()))
			return
		}
		updatedRules = append(updatedRules, r.updateRuleModelFromAPI(ctx, rule, apiRule))
	}
	// Marshal back to types.List
	state.Rules, resp.Diagnostics = types.ListValueFrom(ctx, r.policyRuleObjectType(), updatedRules)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
func (r *appSignOnPolicyRulesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan appSignOnPolicyRulesModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var stateRules, planRules []policyRuleModel
	resp.Diagnostics.Append(state.Rules.ElementsAs(ctx, &stateRules, false)...)
	resp.Diagnostics.Append(plan.Rules.ElementsAs(ctx, &planRules, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	policyID := plan.PolicyID.ValueString()
	client := r.OktaIDaaSClient.OktaSDKSupplementClient()
	// Build lookup structures.
	stateIndex := newRuleIndex(stateRules)
	nameTracker := newNameTracker(stateRules)
	plannedNames := r.buildPlannedNamesSet(planRules)
	plannedIDs := r.buildPlannedIDsSet(planRules)
	// Delete rules removed from plan.
	resp.Diagnostics.Append(r.deleteRemovedRules(ctx, client, policyID, stateRules, plannedNames, plannedIDs)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Process rules in priority order.
	sortedPlanRules := r.sortRulesByPriority(planRules)
	updatedRules := make([]policyRuleModel, 0, len(sortedPlanRules))
	for _, planRule := range sortedPlanRules {
		resultRule, diags := r.processRuleUpdate(ctx, client, policyID, planRule, stateIndex, nameTracker)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		updatedRules = append(updatedRules, resultRule)
	}
	// Reorder to match config order.
	reorderedRules := r.reorderRulesToMatchPlan(updatedRules, planRules)
	plan.ID = plan.PolicyID
	// Marshal back to types.List
	var diags diag.Diagnostics
	plan.Rules, diags = types.ListValueFrom(ctx, r.policyRuleObjectType(), reorderedRules)
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
func (r *appSignOnPolicyRulesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state appSignOnPolicyRulesModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var rules []policyRuleModel
	resp.Diagnostics.Append(state.Rules.ElementsAs(ctx, &rules, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	policyID := state.PolicyID.ValueString()
	client := r.OktaIDaaSClient.OktaSDKSupplementClient()
	for _, rule := range rules {
		// System rules cannot be deleted.
		if rule.System.ValueBool() || rule.ID.IsNull() || rule.ID.IsUnknown() {
			continue
		}
		if err := r.deleteRule(ctx, client, policyID, rule.ID.ValueString()); err != nil {
			resp.Diagnostics.AddError("Error deleting app sign-on policy rule",
				fmt.Sprintf("Could not delete rule '%s': %s", rule.Name.ValueString(), err.Error()))
			return
		}
	}
}
func (r *appSignOnPolicyRulesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	policyID := req.ID
	client := r.OktaIDaaSClient.OktaSDKSupplementClient()
	sdkRules, apiResp, err := client.ListPolicyRules(ctx, policyID)
	if err != nil {
		resp.Diagnostics.AddError("Error importing app sign-on policy rules",
			fmt.Sprintf("Could not list rules for policy '%s': %s", policyID, err.Error()))
		return
	}
	if apiResp != nil && apiResp.StatusCode == http.StatusNotFound {
		resp.Diagnostics.AddError("Policy not found",
			fmt.Sprintf("Policy '%s' was not found", policyID))
		return
	}
	// Fetch each rule individually to get full AccessPolicyRule details.
	importedRules := make([]policyRuleModel, 0, len(sdkRules))
	for _, sdkRule := range sdkRules {
		apiRule, err := r.readRuleFromAPI(ctx, client, policyID, sdkRule.Id)
		if err != nil {
			resp.Diagnostics.AddError("Error importing app sign-on policy rule",
				fmt.Sprintf("Could not read rule '%s': %s", sdkRule.Name, err.Error()))
			return
		}
		importedRules = append(importedRules, r.convertAPIRuleToModel(ctx, apiRule))
	}
	state := appSignOnPolicyRulesModel{
		ID:       types.StringValue(policyID),
		PolicyID: types.StringValue(policyID),
	}
	// Marshal to types.List
	state.Rules, resp.Diagnostics = types.ListValueFrom(ctx, r.policyRuleObjectType(), importedRules)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
func (r *appSignOnPolicyRulesResource) policyRuleObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":                          types.StringType,
			"name":                        types.StringType,
			"system":                      types.BoolType,
			"status":                      types.StringType,
			"priority":                    types.Int64Type,
			"groups_included":             types.SetType{ElemType: types.StringType},
			"groups_excluded":             types.SetType{ElemType: types.StringType},
			"users_included":              types.SetType{ElemType: types.StringType},
			"users_excluded":              types.SetType{ElemType: types.StringType},
			"network_connection":          types.StringType,
			"network_includes":            types.ListType{ElemType: types.StringType},
			"network_excludes":            types.ListType{ElemType: types.StringType},
			"device_is_registered":        types.BoolType,
			"device_is_managed":           types.BoolType,
			"device_assurances_included":  types.SetType{ElemType: types.StringType},
			"user_types_included":         types.SetType{ElemType: types.StringType},
			"user_types_excluded":         types.SetType{ElemType: types.StringType},
			"custom_expression":           types.StringType,
			"access":                      types.StringType,
			"factor_mode":                 types.StringType,
			"type":                        types.StringType,
			"re_authentication_frequency": types.StringType,
			"inactivity_period":           types.StringType,
			"constraints":                 types.ListType{ElemType: types.StringType},
			"chains":                      types.ListType{ElemType: types.StringType},
			"risk_score":                  types.StringType,
			"platform_include": types.ListType{
				ElemType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"type":          types.StringType,
						"os_type":       types.StringType,
						"os_expression": types.StringType,
					},
				},
			},
		},
	}
}
func (r *appSignOnPolicyRulesResource) buildSchema() schema.Schema {
	return schema.Schema{
		Description: "Manages multiple app sign-on policy rules for a single policy. " +
			"This resource allows you to define all rules for a policy in a single configuration block, " +
			"ensuring consistent priority ordering and avoiding drift issues.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of this resource (same as policy_id).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"policy_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the policy to manage rules for.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"rule": schema.ListNestedBlock{
				Description: "List of policy rules. Rules are processed in priority order (lowest number = highest priority).",
				NestedObject: schema.NestedBlockObject{
					Attributes: r.buildRuleAttributes(),
					Blocks: map[string]schema.Block{
						"platform_include": r.buildPlatformIncludeBlock(),
					},
				},
			},
		},
	}
}
func (r *appSignOnPolicyRulesResource) buildPlatformIncludeBlock() schema.ListNestedBlock {
	return schema.ListNestedBlock{
		Description: "Platform conditions to include.",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"type": schema.StringAttribute{
					Optional:    true,
					Description: "Platform type: ANY, MOBILE, or DESKTOP.",
					Validators:  []validator.String{stringvalidator.OneOf(validPlatformTypes...)},
				},
				"os_type": schema.StringAttribute{
					Optional:    true,
					Description: "OS type: ANY, IOS, ANDROID, WINDOWS, OSX, MACOS, CHROMEOS, or OTHER.",
					Validators:  []validator.String{stringvalidator.OneOf(validOSTypes...)},
				},
				"os_expression": schema.StringAttribute{
					Optional: true,
					Computed: true,
					// Default to "" so that omitting os_expression in config is equivalent
					// to os_expression = "". The Okta API requires the field to be present
					// (non-null) when os_type = "OTHER" but always returns null/empty on
					// read — mirroring the SDKv2 resource which stored "" implicitly.
					Default: stringdefault.StaticString(""),
					Description: "Custom OS expression for advanced matching. Required by the API when os_type is OTHER " +
						"(leave empty or omit to match any OTHER OS). " +
						"The API normalizes empty and wildcard values to null on read; the provider preserves \"\" in state.",
				},
			},
		},
	}
}
func (r *appSignOnPolicyRulesResource) buildRuleAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "ID of the rule. Can be specified to adopt an existing rule during migration.",
		},
		"name": schema.StringAttribute{
			Required:    true,
			Description: "Policy Rule Name. Must be unique within the policy.",
		},
		"system": schema.BoolAttribute{
			Computed:    true,
			Description: "Whether this is a system rule (e.g., Catch-all Rule). System rules cannot be modified.",
		},
		"status": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Default:     stringdefault.StaticString("ACTIVE"),
			Description: "Status of the rule: ACTIVE or INACTIVE.",
			Validators:  []validator.String{stringvalidator.OneOf(validStatuses...)},
		},
		"priority": schema.Int64Attribute{
			Optional:    true,
			Description: "Priority of the rule. Lower numbers are evaluated first.",
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"groups_included": schema.SetAttribute{
			Optional:    true,
			ElementType: types.StringType,
			Description: "Set of group IDs to include in this rule.",
		},
		"groups_excluded": schema.SetAttribute{
			Optional:    true,
			ElementType: types.StringType,
			Description: "Set of group IDs to exclude from this rule.",
		},
		"users_included": schema.SetAttribute{
			Optional:    true,
			ElementType: types.StringType,
			Description: "Set of user IDs to include in this rule.",
		},
		"users_excluded": schema.SetAttribute{
			Optional:    true,
			ElementType: types.StringType,
			Description: "Set of user IDs to exclude from this rule.",
		},
		"network_connection": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "Network selection mode: ANYWHERE, ZONE, ON_NETWORK, or OFF_NETWORK.",
			Default:     stringdefault.StaticString("ANYWHERE"),
			Validators:  []validator.String{stringvalidator.OneOf(validNetworkConnections...)},
		},
		"network_includes": schema.ListAttribute{
			Optional:    true,
			ElementType: types.StringType,
			Description: "List of network zone IDs to include.",
		},
		"network_excludes": schema.ListAttribute{
			Optional:    true,
			ElementType: types.StringType,
			Description: "List of network zone IDs to exclude.",
		},
		"device_is_registered": schema.BoolAttribute{
			Optional:    true,
			Description: "Require device to be registered with Okta Verify.",
		},
		"device_is_managed": schema.BoolAttribute{
			Optional:    true,
			Description: "Require device to be managed by a device management system.",
		},
		"device_assurances_included": schema.SetAttribute{
			Optional:    true,
			ElementType: types.StringType,
			Description: "Set of device assurance policy IDs to include.",
		},
		"user_types_included": schema.SetAttribute{
			Optional:    true,
			ElementType: types.StringType,
			Description: "Set of user type IDs to include.",
		},
		"user_types_excluded": schema.SetAttribute{
			Optional:    true,
			ElementType: types.StringType,
			Description: "Set of user type IDs to exclude.",
		},
		"custom_expression": schema.StringAttribute{
			Optional:    true,
			Description: "Custom Okta Expression Language condition for advanced matching.",
		},
		"access": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Default:     stringdefault.StaticString("ALLOW"),
			Description: "Access decision: ALLOW or DENY.",
			Validators:  []validator.String{stringvalidator.OneOf(validAccessTypes...)},
		},
		"factor_mode": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Default:     stringdefault.StaticString("2FA"),
			Description: "Number of factors required: 1FA or 2FA.",
			Validators:  []validator.String{stringvalidator.OneOf(validFactorModes...)},
		},
		"type": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Default:     stringdefault.StaticString("ASSURANCE"),
			Description: "Verification method type.",
		},
		"re_authentication_frequency": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "Re-authentication frequency in ISO 8601 duration format (e.g., PT2H for 2 hours). When using authentication chains with reauthenticateIn, this value is computed by the API based on the chain configuration.",
			PlanModifiers: []planmodifier.String{
				reauthFrequencyModifier{},
			},
		},
		"inactivity_period": schema.StringAttribute{
			Optional:    true,
			Description: "Inactivity period before re-authentication in ISO 8601 duration format.",
		},
		"constraints": schema.ListAttribute{
			Optional:    true,
			ElementType: types.StringType,
			Description: "List of authenticator constraints as JSON-encoded strings.",
		},
		"chains": schema.ListAttribute{
			Optional:    true,
			ElementType: types.StringType,
			Description: "List of authentication method chain objects as JSON-encoded strings. Use with `type = \"AUTH_METHOD_CHAIN\"` only.",
		},
		"risk_score": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Default:     stringdefault.StaticString("ANY"),
			Description: "Risk score level to match: ANY, LOW, MEDIUM, or HIGH.",
			Validators:  []validator.String{stringvalidator.OneOf(validRiskScores...)},
		},
	}
}

// sortRulesByPriority returns rules sorted by priority (ascending).
// Rules without priority are placed at the end.
func (r *appSignOnPolicyRulesResource) sortRulesByPriority(rules []policyRuleModel) []policyRuleModel {
	sorted := make([]policyRuleModel, len(rules))
	copy(sorted, rules)
	sort.Slice(sorted, func(i, j int) bool {
		iPriority := sorted[i].Priority
		jPriority := sorted[j].Priority
		if iPriority.IsNull() || iPriority.IsUnknown() {
			return false
		}
		if jPriority.IsNull() || jPriority.IsUnknown() {
			return true
		}
		return iPriority.ValueInt64() < jPriority.ValueInt64()
	})
	return sorted
}

// reorderRulesToMatchPlan reorders processed rules to match the plan's list order.
// This is critical because Terraform expects state order to match plan order.
func (r *appSignOnPolicyRulesResource) reorderRulesToMatchPlan(processedRules, planRules []policyRuleModel) []policyRuleModel {
	rulesByName := make(map[string]policyRuleModel, len(processedRules))
	for _, rule := range processedRules {
		rulesByName[rule.Name.ValueString()] = rule
	}
	result := make([]policyRuleModel, 0, len(planRules))
	for _, planRule := range planRules {
		if rule, ok := rulesByName[planRule.Name.ValueString()]; ok {
			result = append(result, rule)
		}
	}
	return result
}

// buildPlannedNamesSet creates a set of rule names from the plan.
func (r *appSignOnPolicyRulesResource) buildPlannedNamesSet(rules []policyRuleModel) map[string]bool {
	names := make(map[string]bool, len(rules))
	for _, rule := range rules {
		names[rule.Name.ValueString()] = true
	}
	return names
}

// existingRuleInfo holds identification data for an existing rule in Okta.
type existingRuleInfo struct {
	ID       string
	IsSystem bool
}

// fetchExistingRulesByName lists all rules currently in Okta for the policy
// and returns a map of rule name → existingRuleInfo. This is used in Create to
// detect rules that were partially created by an interrupted previous apply,
// and to discover system rules (e.g. Catch-all Rule) for adoption.
func (r *appSignOnPolicyRulesResource) fetchExistingRulesByName(ctx context.Context, client *sdk.APISupplement, policyID string) (map[string]existingRuleInfo, diag.Diagnostics) {
	var diags diag.Diagnostics
	sdkRules, _, err := client.ListPolicyRules(ctx, policyID)
	if err != nil {
		diags.AddError("Error listing existing policy rules",
			fmt.Sprintf("Could not list rules for policy '%s': %s", policyID, err.Error()))
		return nil, diags
	}
	byName := make(map[string]existingRuleInfo, len(sdkRules))
	for _, rule := range sdkRules {
		byName[rule.Name] = existingRuleInfo{
			ID:       rule.Id,
			IsSystem: rule.System != nil && *rule.System,
		}
	}
	return byName, diags
}

// createOrAdoptRule creates a new rule or adopts an existing one if ID is specified.
func (r *appSignOnPolicyRulesResource) createOrAdoptRule(ctx context.Context, client *sdk.APISupplement, policyID string, rule policyRuleModel) (policyRuleModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	// If ID is provided, adopt existing rule by updating it.
	if !rule.ID.IsNull() && !rule.ID.IsUnknown() && rule.ID.ValueString() != "" {
		ruleID := rule.ID.ValueString()
		// Read the existing rule from the API to check if it's a system rule.
		// This mirrors how okta_app_signon_policy_rule handles system rules:
		// build the full payload, then nil out Conditions for system rules.
		existingRule, err := r.readRuleFromAPI(ctx, client, policyID, ruleID)
		if err != nil {
			diags.AddError("Error adopting existing app sign-on policy rule",
				fmt.Sprintf("Could not read rule '%s' (ID: %s): %s", rule.Name.ValueString(), ruleID, err.Error()))
			return policyRuleModel{}, diags
		}
		isSystem := existingRule.System != nil && *existingRule.System
		apiRule := r.buildAPIRuleFromModel(ctx, rule)
		if isSystem {
			// System rules (e.g. Catch-all Rule) reject condition changes.
			// The API requires name, priority, and system=true to be present.
			apiRule.Conditions = nil
			apiRule.System = utils.BoolPtr(true)
			// Preserve the actual name and priority from the API (they can't be changed).
			apiRule.Name = existingRule.Name
			if existingRule.PriorityPtr != nil {
				apiRule.PriorityPtr = existingRule.PriorityPtr
			}
		}
		updatedRule, err := r.updateRuleInAPI(ctx, client, policyID, ruleID, apiRule)
		if err != nil {
			diags.AddError("Error adopting existing app sign-on policy rule",
				fmt.Sprintf("Could not adopt rule '%s' (ID: %s, isSystem: %t): %s", rule.Name.ValueString(), ruleID, isSystem, err.Error()))
			return policyRuleModel{}, diags
		}
		if err := r.syncRuleStatus(ctx, client, policyID, ruleID, updatedRule.Status, rule.Status.ValueString()); err != nil {
			diags.AddError("Error setting status on app sign-on policy rule",
				fmt.Sprintf("Could not set status for rule '%s': %s", rule.Name.ValueString(), err.Error()))
			return policyRuleModel{}, diags
		}
		result := r.updateRuleModelFromAPI(ctx, rule, updatedRule)
		result.Status = rule.Status
		return result, diags
	}
	// Create new rule.
	apiRule := r.buildAPIRuleFromModel(ctx, rule)
	createdRule, err := r.createRuleInAPI(ctx, client, policyID, apiRule)
	if err != nil {
		diags.AddError("Error creating app sign-on policy rule",
			fmt.Sprintf("Could not create rule '%s': %s", rule.Name.ValueString(), err.Error()))
		return policyRuleModel{}, diags
	}
	if err := r.syncRuleStatus(ctx, client, policyID, createdRule.Id, createdRule.Status, rule.Status.ValueString()); err != nil {
		diags.AddError("Error setting status on app sign-on policy rule",
			fmt.Sprintf("Could not set status for rule '%s': %s", rule.Name.ValueString(), err.Error()))
		return policyRuleModel{}, diags
	}
	result := r.updateRuleModelFromAPI(ctx, rule, createdRule)
	result.Status = rule.Status
	return result, diags
}

// normalizeAPIRuleForSystem strips conditions and preserves name/priority for system rules
func (r *appSignOnPolicyRulesResource) normalizeAPIRuleForSystem(
	apiRule *sdk.AccessPolicyRule,
	systemRule *sdk.AccessPolicyRule,
) {
	if systemRule.System == nil || !*systemRule.System {
		return
	}
	// System rules cannot have conditions modified
	apiRule.Conditions = nil
	apiRule.System = utils.BoolPtr(true)
	// Preserve immutable fields from API
	apiRule.Name = systemRule.Name
	apiRule.PriorityPtr = systemRule.PriorityPtr
}

// processRuleUpdate handles updating or creating a single rule during Update.
func (r *appSignOnPolicyRulesResource) processRuleUpdate(
	ctx context.Context,
	client *sdk.APISupplement,
	policyID string,
	planRule policyRuleModel,
	stateIndex *ruleIndex,
	nameTracker *nameTracker,
) (policyRuleModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	// Extract plan rule's ID if explicitly set
	planRuleID := ""
	if !planRule.ID.IsNull() && !planRule.ID.IsUnknown() {
		planRuleID = planRule.ID.ValueString()
	}
	// Never trust injected IDs from positional state mapping.
	// Only fall back to ID if name lookup fails AND user explicitly provided ID.
	planRuleName := planRule.Name.ValueString()
	existingRule, exists := stateIndex.findRule("", planRuleName)
	if !exists && planRuleID != "" {
		// Name lookup failed, but user explicitly provided ID in config.
		// Try ID fallback only in this case (rename scenario).
		existingRule, exists = stateIndex.findRule(planRuleID, "")
	}
	if !exists {
		// Rule not found in state - create new rule
		apiRule := r.buildAPIRuleFromModel(ctx, planRule)
		createdRule, err := r.createRuleInAPI(ctx, client, policyID, apiRule)
		if err != nil {
			diags.AddError("Error creating app sign-on policy rule",
				fmt.Sprintf("Could not create rule '%s': %s", planRuleName, err.Error()))
			return policyRuleModel{}, diags
		}
		if createdRule.Id != "" {
			nameTracker.updateMapping(createdRule.Id, planRuleName)
		}
		if err := r.syncRuleStatus(ctx, client, policyID, createdRule.Id, createdRule.Status, planRule.Status.ValueString()); err != nil {
			diags.AddError("Error setting status on app sign-on policy rule",
				fmt.Sprintf("Could not set status for rule '%s': %s", planRuleName, err.Error()))
			return policyRuleModel{}, diags
		}
		result := r.updateRuleModelFromAPI(ctx, planRule, createdRule)
		result.Status = planRule.Status
		return result, diags
	}
	// Rule exists in state - update it
	if existingRule.ID.IsNull() {
		diags.AddError("Error updating app sign-on policy rule",
			"Existing rule has no ID")
		return policyRuleModel{}, diags
	}
	existingRuleID := existingRule.ID.ValueString()
	targetName := planRuleName
	// discard the plan ID (it was from positional injection).
	// Use state's correct ID and priority.
	if planRuleID != "" && planRuleID != existingRuleID {
		planRule.ID = existingRule.ID
		planRule.Priority = existingRule.Priority
	}
	// Read current rule from API to determine if it's a system rule
	currentRule, err := r.readRuleFromAPI(ctx, client, policyID, existingRuleID)
	if err != nil {
		diags.AddError("Error updating app sign-on policy rule",
			fmt.Sprintf("Could not read rule '%s' (ID: %s): %s", targetName, existingRuleID, err.Error()))
		return policyRuleModel{}, diags
	}
	isSystem := currentRule.System != nil && *currentRule.System
	// Build the full payload, then normalize for system rules
	apiRule := r.buildAPIRuleFromModel(ctx, planRule)
	if isSystem {
		r.normalizeAPIRuleForSystem(&apiRule, currentRule)
	}
	updatedRule, err := r.updateRuleInAPI(ctx, client, policyID, existingRuleID, apiRule)
	if err != nil {
		diags.AddError("Error updating app sign-on policy rule",
			fmt.Sprintf("Could not update rule '%s': %s", targetName, err.Error()))
		return policyRuleModel{}, diags
	}
	if err := r.syncRuleStatus(ctx, client, policyID, existingRuleID, updatedRule.Status, planRule.Status.ValueString()); err != nil {
		diags.AddError("Error setting status on app sign-on policy rule",
			fmt.Sprintf("Could not set status for rule '%s': %s", targetName, err.Error()))
		return policyRuleModel{}, diags
	}
	nameTracker.updateMapping(existingRuleID, targetName)
	result := r.updateRuleModelFromAPI(ctx, planRule, updatedRule)
	result.Status = planRule.Status
	return result, diags
}

// buildPlannedIDsSet creates a set of rule IDs from the plan.
func (r *appSignOnPolicyRulesResource) buildPlannedIDsSet(rules []policyRuleModel) map[string]bool {
	ids := make(map[string]bool, len(rules))
	for _, rule := range rules {
		if !rule.ID.IsNull() && !rule.ID.IsUnknown() && rule.ID.ValueString() != "" {
			ids[rule.ID.ValueString()] = true
		}
	}
	return ids
}

// deleteRemovedRules deletes rules that are in state but not in plan.
// It checks both name and ID to avoid deleting rules that are being renamed.
func (r *appSignOnPolicyRulesResource) deleteRemovedRules(
	ctx context.Context,
	client *sdk.APISupplement,
	policyID string,
	stateRules []policyRuleModel,
	plannedNames map[string]bool,
	plannedIDs map[string]bool,
) diag.Diagnostics {
	var diags diag.Diagnostics
	for _, rule := range stateRules {
		if rule.System.ValueBool() || rule.ID.IsNull() {
			continue
		}
		// Keep the rule if its name is still in the plan OR if its ID is
		// referenced by a planned rule (i.e. it's being renamed, not removed).
		if plannedNames[rule.Name.ValueString()] || plannedIDs[rule.ID.ValueString()] {
			continue
		}
		if err := r.deleteRule(ctx, client, policyID, rule.ID.ValueString()); err != nil {
			diags.AddError("Error deleting app sign-on policy rule",
				fmt.Sprintf("Could not delete rule '%s': %s", rule.Name.ValueString(), err.Error()))
			return diags
		}
	}
	return diags
}
func (r *appSignOnPolicyRulesResource) createRuleInAPI(ctx context.Context, client *sdk.APISupplement, policyID string, rule sdk.AccessPolicyRule) (*sdk.AccessPolicyRule, error) {
	return backoff.Retry(ctx, func() (*sdk.AccessPolicyRule, error) {
		created, apiResp, err := client.CreateAppSignOnPolicyRule(ctx, policyID, rule)
		if err != nil {
			// Note: isRetryableStatusCode is only checked after ensuring apiResp != nil
			if apiResp != nil && isRetryableStatusCode(apiResp.StatusCode) {
				return nil, err
			}
			return nil, backoff.Permanent(err)
		}
		return created, nil
	}, backoff.WithMaxElapsedTime(apiRetryTimeout))
}
func (r *appSignOnPolicyRulesResource) readRuleFromAPI(ctx context.Context, client *sdk.APISupplement, policyID, ruleID string) (*sdk.AccessPolicyRule, error) {
	return backoff.Retry(ctx, func() (*sdk.AccessPolicyRule, error) {
		rule, apiResp, err := client.GetAppSignOnPolicyRule(ctx, policyID, ruleID)
		if err != nil {
			// Note: isRetryableStatusCode is only checked after ensuring apiResp != nil
			if apiResp != nil && isRetryableStatusCode(apiResp.StatusCode) {
				return nil, err
			}
			return nil, backoff.Permanent(err)
		}
		return rule, nil
	}, backoff.WithMaxElapsedTime(apiRetryTimeout))
}
func (r *appSignOnPolicyRulesResource) updateRuleInAPI(ctx context.Context, client *sdk.APISupplement, policyID, ruleID string, rule sdk.AccessPolicyRule) (*sdk.AccessPolicyRule, error) {
	return backoff.Retry(ctx, func() (*sdk.AccessPolicyRule, error) {
		updated, apiResp, err := client.UpdateAppSignOnPolicyRule(ctx, policyID, ruleID, rule)
		if err != nil {
			if apiResp != nil && isRetryableStatusCode(apiResp.StatusCode) {
				return nil, err
			}
			return nil, backoff.Permanent(err)
		}
		return updated, nil
	}, backoff.WithMaxElapsedTime(apiRetryTimeout))
}
func (r *appSignOnPolicyRulesResource) deleteRule(ctx context.Context, client *sdk.APISupplement, policyID, ruleID string) error {
	_, err := backoff.Retry(ctx, func() (struct{}, error) {
		apiResp, err := client.DeleteAppSignOnPolicyRule(ctx, policyID, ruleID)
		if err != nil {
			// Note: isRetryableStatusCode is only checked after ensuring apiResp != nil
			if apiResp != nil && isRetryableStatusCode(apiResp.StatusCode) {
				return struct{}{}, err
			}
			if apiResp != nil && apiResp.StatusCode == http.StatusNotFound {
				return struct{}{}, nil // Already deleted.
			}
			return struct{}{}, backoff.Permanent(err)
		}
		return struct{}{}, nil
	}, backoff.WithMaxElapsedTime(apiRetryTimeout))
	return err
}

// syncRuleStatus calls the lifecycle endpoint when the desired status differs
// from the current API status. The Okta API ignores the status field in the
// rule body; status changes must go through /lifecycle/activate or
// /lifecycle/deactivate.
func (r *appSignOnPolicyRulesResource) syncRuleStatus(ctx context.Context, client *sdk.APISupplement, policyID, ruleID, currentStatus, desiredStatus string) error {
	if desiredStatus == "" || desiredStatus == currentStatus {
		return nil
	}
	if desiredStatus == StatusActive {
		_, err := client.ActivateAppSignOnPolicyRule(ctx, policyID, ruleID)
		return err
	}
	_, err := client.DeactivateAppSignOnPolicyRule(ctx, policyID, ruleID)
	return err
}

// isRetryableStatusCode returns true for HTTP status codes that should trigger a retry.
func isRetryableStatusCode(statusCode int) bool {
	return statusCode == http.StatusConflict ||
		statusCode == http.StatusTooManyRequests ||
		statusCode == http.StatusServiceUnavailable
}

// buildAPIRuleFromModel converts a Terraform model to an Okta API rule.
func (r *appSignOnPolicyRulesResource) buildAPIRuleFromModel(ctx context.Context, rule policyRuleModel) sdk.AccessPolicyRule {
	apiRule := sdk.AccessPolicyRule{
		Name: rule.Name.ValueString(),
		Type: ruleTypeAccessPolicy,
		Actions: &sdk.AccessPolicyRuleActions{
			AppSignOn: &sdk.AccessPolicyRuleApplicationSignOn{
				Access: rule.Access.ValueString(),
				VerificationMethod: &sdk.VerificationMethod{
					FactorMode:       rule.FactorMode.ValueString(),
					ReauthenticateIn: rule.ReAuthenticationFrequency.ValueString(),
					Type:             rule.Type.ValueString(),
				},
			},
		},
		Conditions: &sdk.AccessPolicyRuleConditions{},
	}
	r.setAPIPriority(&apiRule, rule)
	r.setAPIInactivityPeriod(&apiRule, rule)
	r.setAPIConstraints(ctx, &apiRule, rule)
	r.setAPIChains(ctx, &apiRule, rule)
	r.setAPINetworkConditions(ctx, &apiRule, rule)
	r.setAPIPlatformConditions(&apiRule, rule)
	r.setAPICustomExpression(&apiRule, rule)
	r.setAPIRiskScore(&apiRule, rule)
	r.setAPIDeviceConditions(ctx, &apiRule, rule)
	r.setAPIPeopleConditions(ctx, &apiRule, rule)
	r.setAPIUserTypeConditions(ctx, &apiRule, rule)
	return apiRule
}
func (r *appSignOnPolicyRulesResource) setAPIPriority(apiRule *sdk.AccessPolicyRule, rule policyRuleModel) {
	if !rule.Priority.IsNull() && !rule.Priority.IsUnknown() {
		priority := rule.Priority.ValueInt64()
		apiRule.PriorityPtr = &priority
	}
}
func (r *appSignOnPolicyRulesResource) setAPIInactivityPeriod(apiRule *sdk.AccessPolicyRule, rule policyRuleModel) {
	if !rule.InactivityPeriod.IsNull() && !rule.InactivityPeriod.IsUnknown() {
		apiRule.Actions.AppSignOn.VerificationMethod.InactivityPeriod = rule.InactivityPeriod.ValueString()
	}
}
func (r *appSignOnPolicyRulesResource) setAPIConstraints(ctx context.Context, apiRule *sdk.AccessPolicyRule, rule policyRuleModel) {
	if rule.Constraints.IsNull() || rule.Constraints.IsUnknown() {
		return
	}
	var constraintStrings []string
	rule.Constraints.ElementsAs(ctx, &constraintStrings, false)
	var constraints []*sdk.AccessPolicyConstraints
	for _, c := range constraintStrings {
		var constraint sdk.AccessPolicyConstraints
		if err := json.Unmarshal([]byte(c), &constraint); err == nil {
			constraints = append(constraints, &constraint)
		}
	}
	apiRule.Actions.AppSignOn.VerificationMethod.Constraints = constraints
}

func (r *appSignOnPolicyRulesResource) setAPIChains(ctx context.Context, apiRule *sdk.AccessPolicyRule, rule policyRuleModel) {
	if rule.Chains.IsNull() || rule.Chains.IsUnknown() {
		return
	}
	var chainStrings []string
	rule.Chains.ElementsAs(ctx, &chainStrings, false)
	var chains []*sdk.AccessPolicyChains
	hasReauthenticateIn := false
	for _, c := range chainStrings {
		var chain sdk.AccessPolicyChains
		if err := json.Unmarshal([]byte(c), &chain); err == nil {
			chains = append(chains, &chain)
			if chain.ReauthenticateIn != "" {
				hasReauthenticateIn = true
			}
		}
	}
	apiRule.Actions.AppSignOn.VerificationMethod.Chains = chains
	// If any chain sets ReauthenticateIn, clear the top-level field to avoid
	// the API rejecting the combination (mirrors behaviour in singular resource).
	if hasReauthenticateIn {
		apiRule.Actions.AppSignOn.VerificationMethod.ReauthenticateIn = ""
	}
}

func (r *appSignOnPolicyRulesResource) setAPINetworkConditions(ctx context.Context, apiRule *sdk.AccessPolicyRule, rule policyRuleModel) {
	if rule.NetworkConnection.IsNull() || rule.NetworkConnection.IsUnknown() {
		return
	}
	apiRule.Conditions.Network = &sdk.PolicyNetworkCondition{
		Connection: rule.NetworkConnection.ValueString(),
	}
	if !rule.NetworkIncludes.IsNull() && !rule.NetworkIncludes.IsUnknown() {
		var includes []string
		rule.NetworkIncludes.ElementsAs(ctx, &includes, false)
		apiRule.Conditions.Network.Include = includes
	}
	if !rule.NetworkExcludes.IsNull() && !rule.NetworkExcludes.IsUnknown() {
		var excludes []string
		rule.NetworkExcludes.ElementsAs(ctx, &excludes, false)
		apiRule.Conditions.Network.Exclude = excludes
	}
}
func (r *appSignOnPolicyRulesResource) setAPIPlatformConditions(apiRule *sdk.AccessPolicyRule, rule policyRuleModel) {
	if len(rule.PlatformInclude) == 0 {
		return
	}
	var platforms []*sdk.PlatformConditionEvaluatorPlatform
	for _, p := range rule.PlatformInclude {
		platform := &sdk.PlatformConditionEvaluatorPlatform{}
		if !p.Type.IsNull() && !p.Type.IsUnknown() {
			platform.Type = p.Type.ValueString()
		}
		if !p.OsType.IsNull() && !p.OsType.IsUnknown() {
			osType := p.OsType.ValueString()
			platform.Os = &sdk.PlatformConditionEvaluatorPlatformOperatingSystem{
				Type: osType,
			}
			// The API requires os_expression when os_type is OTHER (returns 400 if absent).
			// It accepts (and ignores) an empty string, normalizing it to null on read.
			// Always send a non-nil pointer for OTHER — default to "" when the user
			// omitted it, passed null, or the value is still unknown at apply time
			// (e.g. try(..., null) in a dynamic block overrides the schema Default).
			if osType == "OTHER" {
				expr := ""
				if !p.OsExpression.IsNull() && !p.OsExpression.IsUnknown() && p.OsExpression.ValueString() != "" {
					expr = p.OsExpression.ValueString()
				}
				platform.Os.Expression = &expr
			}
		}
		platforms = append(platforms, platform)
	}
	apiRule.Conditions.Platform = &sdk.PlatformPolicyRuleCondition{Include: platforms}
}
func (r *appSignOnPolicyRulesResource) setAPICustomExpression(apiRule *sdk.AccessPolicyRule, rule policyRuleModel) {
	if !rule.CustomExpression.IsNull() && !rule.CustomExpression.IsUnknown() {
		apiRule.Conditions.ElCondition = &sdk.AccessPolicyRuleCustomCondition{
			Condition: rule.CustomExpression.ValueString(),
		}
	}
}
func (r *appSignOnPolicyRulesResource) setAPIRiskScore(apiRule *sdk.AccessPolicyRule, rule policyRuleModel) {
	if !rule.RiskScore.IsNull() && !rule.RiskScore.IsUnknown() {
		apiRule.Conditions.RiskScore = &sdk.RiskScorePolicyRuleCondition{
			Level: rule.RiskScore.ValueString(),
		}
	}
}
func (r *appSignOnPolicyRulesResource) setAPIDeviceConditions(ctx context.Context, apiRule *sdk.AccessPolicyRule, rule policyRuleModel) {
	hasRegistered := !rule.DeviceIsRegistered.IsNull() && !rule.DeviceIsRegistered.IsUnknown()
	hasManaged := !rule.DeviceIsManaged.IsNull() && !rule.DeviceIsManaged.IsUnknown()
	hasAssurances := !rule.DeviceAssurancesIncluded.IsNull() && !rule.DeviceAssurancesIncluded.IsUnknown()
	if !hasRegistered && !hasManaged && !hasAssurances {
		return
	}
	apiRule.Conditions.Device = &sdk.DeviceAccessPolicyRuleCondition{}
	if hasRegistered {
		apiRule.Conditions.Device.Registered = utils.BoolPtr(rule.DeviceIsRegistered.ValueBool())
	}
	if hasManaged {
		apiRule.Conditions.Device.Managed = utils.BoolPtr(rule.DeviceIsManaged.ValueBool())
	}
	if hasAssurances {
		var assurances []string
		rule.DeviceAssurancesIncluded.ElementsAs(ctx, &assurances, false)
		apiRule.Conditions.Device.Assurance = &sdk.DeviceAssurancePolicyRuleCondition{Include: assurances}
	}
}
func (r *appSignOnPolicyRulesResource) setAPIPeopleConditions(ctx context.Context, apiRule *sdk.AccessPolicyRule, rule policyRuleModel) {
	var usersIncluded, usersExcluded, groupsIncluded, groupsExcluded []string
	if !rule.UsersIncluded.IsNull() && !rule.UsersIncluded.IsUnknown() {
		rule.UsersIncluded.ElementsAs(ctx, &usersIncluded, false)
	}
	if !rule.UsersExcluded.IsNull() && !rule.UsersExcluded.IsUnknown() {
		rule.UsersExcluded.ElementsAs(ctx, &usersExcluded, false)
	}
	if !rule.GroupsIncluded.IsNull() && !rule.GroupsIncluded.IsUnknown() {
		rule.GroupsIncluded.ElementsAs(ctx, &groupsIncluded, false)
	}
	if !rule.GroupsExcluded.IsNull() && !rule.GroupsExcluded.IsUnknown() {
		rule.GroupsExcluded.ElementsAs(ctx, &groupsExcluded, false)
	}
	hasUserCondition := len(usersIncluded) > 0 || len(usersExcluded) > 0
	hasGroupCondition := len(groupsIncluded) > 0 || len(groupsExcluded) > 0
	if !hasUserCondition && !hasGroupCondition {
		return
	}
	apiRule.Conditions.People = &sdk.PolicyPeopleCondition{}
	if hasUserCondition {
		apiRule.Conditions.People.Users = &sdk.UserCondition{
			Include: usersIncluded,
			Exclude: usersExcluded,
		}
	}
	if hasGroupCondition {
		apiRule.Conditions.People.Groups = &sdk.GroupCondition{
			Include: groupsIncluded,
			Exclude: groupsExcluded,
		}
	}
}
func (r *appSignOnPolicyRulesResource) setAPIUserTypeConditions(ctx context.Context, apiRule *sdk.AccessPolicyRule, rule policyRuleModel) {
	var included, excluded []string
	if !rule.UserTypesIncluded.IsNull() && !rule.UserTypesIncluded.IsUnknown() {
		rule.UserTypesIncluded.ElementsAs(ctx, &included, false)
	}
	if !rule.UserTypesExcluded.IsNull() && !rule.UserTypesExcluded.IsUnknown() {
		rule.UserTypesExcluded.ElementsAs(ctx, &excluded, false)
	}
	if len(included) > 0 || len(excluded) > 0 {
		apiRule.Conditions.UserType = &sdk.UserTypeCondition{
			Include: included,
			Exclude: excluded,
		}
	}
}

// updateRuleModelFromAPI updates a Terraform model with data from an API response.
// It preserves the original plan/state value for optional collection fields that
// the user did not configure (null), so that Terraform's consistency check
// (plan vs post-apply state) is not violated for non-Computed attributes.
// Fields the user did configure (non-null) are updated from the API response.
func (r *appSignOnPolicyRulesResource) updateRuleModelFromAPI(ctx context.Context, rule policyRuleModel, apiRule *sdk.AccessPolicyRule) policyRuleModel {
	rule.ID = types.StringValue(apiRule.Id)
	rule.Name = types.StringValue(apiRule.Name)
	rule.System = types.BoolValue(apiRule.System != nil && *apiRule.System)
	rule.Status = types.StringValue(apiRule.Status)
	if apiRule.PriorityPtr != nil {
		rule.Priority = types.Int64Value(*apiRule.PriorityPtr)
	}
	r.updateRuleActionsFromAPI(ctx, &rule, apiRule)
	r.updateRuleConditionsFromAPI(ctx, &rule, apiRule)
	return rule
}
func (r *appSignOnPolicyRulesResource) updateRuleActionsFromAPI(ctx context.Context, rule *policyRuleModel, apiRule *sdk.AccessPolicyRule) {
	if apiRule.Actions == nil || apiRule.Actions.AppSignOn == nil {
		return
	}
	rule.Access = types.StringValue(apiRule.Actions.AppSignOn.Access)
	vm := apiRule.Actions.AppSignOn.VerificationMethod
	if vm == nil {
		return
	}
	// If API returns empty FactorMode (which can happen with ASSURANCE + chains), default to 2FA
	if vm.FactorMode != "" {
		rule.FactorMode = types.StringValue(vm.FactorMode)
	} else {
		rule.FactorMode = types.StringValue("2FA")
	}
	rule.Type = types.StringValue(vm.Type)
	rule.ReAuthenticationFrequency = types.StringValue(vm.ReauthenticateIn)
	if vm.InactivityPeriod != "" {
		rule.InactivityPeriod = types.StringValue(vm.InactivityPeriod)
	}

	// Convert chains to JSON strings.
	if len(vm.Chains) > 0 {
		var chainStrings []string
		for _, chain := range vm.Chains {
			if jsonBytes, err := json.Marshal(chain); err == nil {
				chainStrings = append(chainStrings, string(jsonBytes))
			}
		}
		rule.Chains, _ = types.ListValueFrom(ctx, types.StringType, chainStrings)
	}
}
func (r *appSignOnPolicyRulesResource) updateRuleConditionsFromAPI(ctx context.Context, rule *policyRuleModel, apiRule *sdk.AccessPolicyRule) {
	if apiRule.Conditions == nil {
		return
	}
	c := apiRule.Conditions
	// Network conditions
	if c.Network != nil {
		// Only set network_connection if the user configured it (non-null in state).
		// If null, leave it null to avoid "was null, but now ANYWHERE" inconsistency.
		if !rule.NetworkConnection.IsNull() {
			rule.NetworkConnection = types.StringValue(c.Network.Connection)
		}
		if len(c.Network.Include) > 0 {
			rule.NetworkIncludes, _ = types.ListValueFrom(ctx, types.StringType, c.Network.Include)
		} else if !rule.NetworkIncludes.IsNull() {
			rule.NetworkIncludes, _ = types.ListValueFrom(ctx, types.StringType, []string{})
		}
		if len(c.Network.Exclude) > 0 {
			rule.NetworkExcludes, _ = types.ListValueFrom(ctx, types.StringType, c.Network.Exclude)
		} else if !rule.NetworkExcludes.IsNull() {
			rule.NetworkExcludes, _ = types.ListValueFrom(ctx, types.StringType, []string{})
		}
	}
	// Risk score
	if c.RiskScore != nil {
		rule.RiskScore = types.StringValue(c.RiskScore.Level)
	}
	// Custom expression
	if c.ElCondition != nil && c.ElCondition.Condition != "" {
		rule.CustomExpression = types.StringValue(c.ElCondition.Condition)
	} else if !rule.CustomExpression.IsNull() {
		rule.CustomExpression = types.StringNull()
	}
	// Device conditions
	if c.Device != nil {
		if c.Device.Registered != nil {
			rule.DeviceIsRegistered = types.BoolValue(*c.Device.Registered)
		}
		if c.Device.Managed != nil {
			rule.DeviceIsManaged = types.BoolValue(*c.Device.Managed)
		}
		if c.Device.Assurance != nil && len(c.Device.Assurance.Include) > 0 {
			rule.DeviceAssurancesIncluded, _ = types.SetValueFrom(ctx, types.StringType, c.Device.Assurance.Include)
		} else if !rule.DeviceAssurancesIncluded.IsNull() {
			rule.DeviceAssurancesIncluded, _ = types.SetValueFrom(ctx, types.StringType, []string{})
		}
	}
	// People conditions
	if c.People != nil {
		if c.People.Users != nil {
			if len(c.People.Users.Include) > 0 {
				rule.UsersIncluded, _ = types.SetValueFrom(ctx, types.StringType, c.People.Users.Include)
			} else if !rule.UsersIncluded.IsNull() {
				rule.UsersIncluded, _ = types.SetValueFrom(ctx, types.StringType, []string{})
			}
			if len(c.People.Users.Exclude) > 0 {
				rule.UsersExcluded, _ = types.SetValueFrom(ctx, types.StringType, c.People.Users.Exclude)
			} else if !rule.UsersExcluded.IsNull() {
				rule.UsersExcluded, _ = types.SetValueFrom(ctx, types.StringType, []string{})
			}
		}
		if c.People.Groups != nil {
			if len(c.People.Groups.Include) > 0 {
				rule.GroupsIncluded, _ = types.SetValueFrom(ctx, types.StringType, c.People.Groups.Include)
			} else if !rule.GroupsIncluded.IsNull() {
				rule.GroupsIncluded, _ = types.SetValueFrom(ctx, types.StringType, []string{})
			}
			if len(c.People.Groups.Exclude) > 0 {
				rule.GroupsExcluded, _ = types.SetValueFrom(ctx, types.StringType, c.People.Groups.Exclude)
			} else if !rule.GroupsExcluded.IsNull() {
				rule.GroupsExcluded, _ = types.SetValueFrom(ctx, types.StringType, []string{})
			}
		}
	}
	// User type conditions
	if c.UserType != nil {
		if len(c.UserType.Include) > 0 {
			rule.UserTypesIncluded, _ = types.SetValueFrom(ctx, types.StringType, c.UserType.Include)
		} else if !rule.UserTypesIncluded.IsNull() {
			rule.UserTypesIncluded, _ = types.SetValueFrom(ctx, types.StringType, []string{})
		}
		if len(c.UserType.Exclude) > 0 {
			rule.UserTypesExcluded, _ = types.SetValueFrom(ctx, types.StringType, c.UserType.Exclude)
		} else if !rule.UserTypesExcluded.IsNull() {
			rule.UserTypesExcluded, _ = types.SetValueFrom(ctx, types.StringType, []string{})
		}
	}
	// Platform conditions
	if c.Platform != nil && len(c.Platform.Include) > 0 {
		var platforms []platformIncludeModel
		for _, p := range c.Platform.Include {
			platform := platformIncludeModel{Type: types.StringValue(p.Type)}
			if p.Os != nil {
				// The API silently maps os_type=ANY to OTHER on read.
				osType := p.Os.Type
				if osType == "ANY" {
					osType = "OTHER"
				}
				platform.OsType = types.StringValue(osType)
				if p.Os.Expression != nil && *p.Os.Expression != "" {
					// API returned a real expression — use it.
					platform.OsExpression = types.StringValue(*p.Os.Expression)
				} else {
					// API returns null/empty for os_expression on all os_types.
					// Always store "" to match the schema Default and avoid
					// "was cty.StringVal("\"\""), but now null" inconsistency.
					platform.OsExpression = types.StringValue("")
				}
			}
			platforms = append(platforms, platform)
		}
		rule.PlatformInclude = platforms
	} else if len(rule.PlatformInclude) > 0 {
		rule.PlatformInclude = nil
	}
}

// convertAPIRuleToModel creates a new Terraform model from an API response (used for import).
func (r *appSignOnPolicyRulesResource) convertAPIRuleToModel(ctx context.Context, apiRule *sdk.AccessPolicyRule) policyRuleModel {
	rule := policyRuleModel{
		ID:     types.StringValue(apiRule.Id),
		Name:   types.StringValue(apiRule.Name),
		System: types.BoolValue(apiRule.System != nil && *apiRule.System),
		Status: types.StringValue(apiRule.Status),
		// Initialize all List/Set fields with typed nulls to avoid "MISSING TYPE" errors.
		// Sub-functions will overwrite these with actual values when the API returns data.
		GroupsIncluded:           types.SetNull(types.StringType),
		GroupsExcluded:           types.SetNull(types.StringType),
		UsersIncluded:            types.SetNull(types.StringType),
		UsersExcluded:            types.SetNull(types.StringType),
		NetworkIncludes:          types.ListNull(types.StringType),
		NetworkExcludes:          types.ListNull(types.StringType),
		DeviceAssurancesIncluded: types.SetNull(types.StringType),
		UserTypesIncluded:        types.SetNull(types.StringType),
		UserTypesExcluded:        types.SetNull(types.StringType),
		Constraints:              types.ListNull(types.StringType),
		Chains:                   types.ListNull(types.StringType),
	}
	if apiRule.PriorityPtr != nil {
		rule.Priority = types.Int64Value(*apiRule.PriorityPtr)
	}
	r.convertAPIActionsToModel(ctx, &rule, apiRule)
	r.convertAPIConditionsToModel(ctx, &rule, apiRule)
	return rule
}
func (r *appSignOnPolicyRulesResource) convertAPIActionsToModel(ctx context.Context, rule *policyRuleModel, apiRule *sdk.AccessPolicyRule) {
	if apiRule.Actions == nil || apiRule.Actions.AppSignOn == nil {
		return
	}
	rule.Access = types.StringValue(apiRule.Actions.AppSignOn.Access)
	vm := apiRule.Actions.AppSignOn.VerificationMethod
	if vm == nil {
		return
	}
	// If API returns empty FactorMode (which can happen with ASSURANCE + chains), default to 2FA
	if vm.FactorMode != "" {
		rule.FactorMode = types.StringValue(vm.FactorMode)
	} else {
		rule.FactorMode = types.StringValue("2FA")
	}
	rule.Type = types.StringValue(vm.Type)
	rule.ReAuthenticationFrequency = types.StringValue(vm.ReauthenticateIn)
	if vm.InactivityPeriod != "" {
		rule.InactivityPeriod = types.StringValue(vm.InactivityPeriod)
	}
	// Convert constraints to JSON strings.
	if len(vm.Constraints) > 0 {
		var constraintStrings []string
		for _, constraint := range vm.Constraints {
			if jsonBytes, err := json.Marshal(constraint); err == nil {
				constraintStrings = append(constraintStrings, string(jsonBytes))
			}
		}
		rule.Constraints, _ = types.ListValueFrom(ctx, types.StringType, constraintStrings)
	}

	// Convert chains to JSON strings.
	if len(vm.Chains) > 0 {
		var chainStrings []string
		for _, chain := range vm.Chains {
			if jsonBytes, err := json.Marshal(chain); err == nil {
				chainStrings = append(chainStrings, string(jsonBytes))
			}
		}
		rule.Chains, _ = types.ListValueFrom(ctx, types.StringType, chainStrings)
	}
}
func (r *appSignOnPolicyRulesResource) convertAPIConditionsToModel(ctx context.Context, rule *policyRuleModel, apiRule *sdk.AccessPolicyRule) {
	if apiRule.Conditions == nil {
		return
	}
	c := apiRule.Conditions
	// Network
	if c.Network != nil {
		rule.NetworkConnection = types.StringValue(c.Network.Connection)
		if len(c.Network.Include) > 0 {
			rule.NetworkIncludes, _ = types.ListValueFrom(ctx, types.StringType, c.Network.Include)
		}
		if len(c.Network.Exclude) > 0 {
			rule.NetworkExcludes, _ = types.ListValueFrom(ctx, types.StringType, c.Network.Exclude)
		}
	}
	// Risk score
	if c.RiskScore != nil {
		rule.RiskScore = types.StringValue(c.RiskScore.Level)
	}
	// Custom expression
	if c.ElCondition != nil && c.ElCondition.Condition != "" {
		rule.CustomExpression = types.StringValue(c.ElCondition.Condition)
	}
	// Device
	if c.Device != nil {
		if c.Device.Registered != nil {
			rule.DeviceIsRegistered = types.BoolValue(*c.Device.Registered)
		}
		if c.Device.Managed != nil {
			rule.DeviceIsManaged = types.BoolValue(*c.Device.Managed)
		}
		if c.Device.Assurance != nil && len(c.Device.Assurance.Include) > 0 {
			rule.DeviceAssurancesIncluded, _ = types.SetValueFrom(ctx, types.StringType, c.Device.Assurance.Include)
		}
	}
	// People
	if c.People != nil {
		if c.People.Users != nil {
			if len(c.People.Users.Include) > 0 {
				rule.UsersIncluded, _ = types.SetValueFrom(ctx, types.StringType, c.People.Users.Include)
			}
			if len(c.People.Users.Exclude) > 0 {
				rule.UsersExcluded, _ = types.SetValueFrom(ctx, types.StringType, c.People.Users.Exclude)
			}
		}
		if c.People.Groups != nil {
			if len(c.People.Groups.Include) > 0 {
				rule.GroupsIncluded, _ = types.SetValueFrom(ctx, types.StringType, c.People.Groups.Include)
			}
			if len(c.People.Groups.Exclude) > 0 {
				rule.GroupsExcluded, _ = types.SetValueFrom(ctx, types.StringType, c.People.Groups.Exclude)
			}
		}
	}
	// User types
	if c.UserType != nil {
		if len(c.UserType.Include) > 0 {
			rule.UserTypesIncluded, _ = types.SetValueFrom(ctx, types.StringType, c.UserType.Include)
		}
		if len(c.UserType.Exclude) > 0 {
			rule.UserTypesExcluded, _ = types.SetValueFrom(ctx, types.StringType, c.UserType.Exclude)
		}
	}
	// Platforms
	if c.Platform != nil && len(c.Platform.Include) > 0 {
		var platforms []platformIncludeModel
		for _, p := range c.Platform.Include {
			platform := platformIncludeModel{Type: types.StringValue(p.Type)}
			if p.Os != nil {
				// The API silently maps os_type=ANY to OTHER on read.
				osType := p.Os.Type
				if osType == "ANY" {
					osType = "OTHER"
				}
				platform.OsType = types.StringValue(osType)
				if p.Os.Expression != nil && *p.Os.Expression != "" {
					platform.OsExpression = types.StringValue(*p.Os.Expression)
				} else {
					platform.OsExpression = types.StringValue("")
				}
			}
			platforms = append(platforms, platform)
		}
		rule.PlatformInclude = platforms
	}
}

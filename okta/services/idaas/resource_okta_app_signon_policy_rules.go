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
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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

	// tempRuleNamePrefix is used when temporarily renaming rules to avoid name conflicts.
	tempRuleNamePrefix = "__temp_"
)

// Validation values for rule attributes.
var (
	validStatuses           = []string{"ACTIVE", "INACTIVE"}
	validNetworkConnections = []string{"ANYWHERE", "ZONE", "ON_NETWORK", "OFF_NETWORK"}
	validAccessTypes        = []string{"ALLOW", "DENY"}
	validFactorModes        = []string{"1FA", "2FA"}
	validRiskScores         = []string{"ANY", "LOW", "MEDIUM", "HIGH"}
	validPlatformTypes      = []string{"ANY", "MOBILE", "DESKTOP"}
	validOSTypes            = []string{"ANY", "IOS", "ANDROID", "WINDOWS", "OSX", "MACOS", "CHROMEOS", "OTHER"}
)

var (
	_ resource.Resource                = &appSignOnPolicyRulesResource{}
	_ resource.ResourceWithConfigure   = &appSignOnPolicyRulesResource{}
	_ resource.ResourceWithImportState = &appSignOnPolicyRulesResource{}
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
	ID       types.String      `tfsdk:"id"`
	PolicyID types.String      `tfsdk:"policy_id"`
	Rules    []policyRuleModel `tfsdk:"rule"`
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
	RiskScore                 types.String           `tfsdk:"risk_score"`
	PlatformInclude           []platformIncludeModel `tfsdk:"platform_include"`
}

// platformIncludeModel represents platform conditions in the rule.
type platformIncludeModel struct {
	Type         types.String `tfsdk:"type"`
	OsType       types.String `tfsdk:"os_type"`
	OsExpression types.String `tfsdk:"os_expression"`
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

// findRule looks up a rule by ID first, then by name.
func (idx *ruleIndex) findRule(id, name string) (policyRuleModel, bool) {
	if id != "" {
		if rule, ok := idx.byID[id]; ok {
			return rule, true
		}
	}
	if rule, ok := idx.byName[name]; ok {
		return rule, true
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
		if !rule.ID.IsNull() && !rule.ID.IsUnknown() {
			id := rule.ID.ValueString()
			name := rule.Name.ValueString()
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

func (r *appSignOnPolicyRulesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan appSignOnPolicyRulesModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

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
	sortedRules := r.sortRulesByPriority(plan.Rules)
	createdRules := make([]policyRuleModel, 0, len(sortedRules))

	for _, rule := range sortedRules {
		// If an existing rule was found by name, inject its ID so createOrAdoptRule
		// will update it rather than attempting to create a duplicate.
		// Check both IsNull() and IsUnknown() because during a fresh Create (no prior state),
		// Computed attributes come through as Unknown, not Null.
		if existingID, found := existingByName[rule.Name.ValueString()]; found && (rule.ID.IsNull() || rule.ID.IsUnknown()) {
			rule.ID = types.StringValue(existingID)
		}
		resultRule, diags := r.createOrAdoptRule(ctx, client, policyID, rule)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		createdRules = append(createdRules, resultRule)
	}

	// Reorder to match config order (Terraform expects state order to match plan order).
	plan.Rules = r.reorderRulesToMatchPlan(createdRules, plan.Rules)
	plan.ID = plan.PolicyID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *appSignOnPolicyRulesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state appSignOnPolicyRulesModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyID := state.PolicyID.ValueString()
	client := r.OktaIDaaSClient.OktaSDKSupplementClient()

	updatedRules := make([]policyRuleModel, 0, len(state.Rules))
	for _, rule := range state.Rules {
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

	state.Rules = updatedRules
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *appSignOnPolicyRulesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan appSignOnPolicyRulesModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyID := plan.PolicyID.ValueString()
	client := r.OktaIDaaSClient.OktaSDKSupplementClient()

	// Build lookup structures.
	stateIndex := newRuleIndex(state.Rules)
	nameTracker := newNameTracker(state.Rules)
	plannedNames := r.buildPlannedNamesSet(plan.Rules)
	plannedIDs := r.buildPlannedIDsSet(plan.Rules)

	// Delete rules removed from plan.
	resp.Diagnostics.Append(r.deleteRemovedRules(ctx, client, policyID, state.Rules, plannedNames, plannedIDs)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build map of planned name changes for conflict detection.
	plannedNameByID := r.buildPlannedNameByID(plan.Rules)

	// Process rules in priority order.
	sortedPlanRules := r.sortRulesByPriority(plan.Rules)
	updatedRules := make([]policyRuleModel, 0, len(sortedPlanRules))

	for _, planRule := range sortedPlanRules {
		resultRule, diags := r.processRuleUpdate(ctx, client, policyID, planRule, stateIndex, nameTracker, plannedNameByID)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		updatedRules = append(updatedRules, resultRule)
	}

	// Reorder to match config order.
	plan.Rules = r.reorderRulesToMatchPlan(updatedRules, plan.Rules)
	plan.ID = plan.PolicyID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *appSignOnPolicyRulesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state appSignOnPolicyRulesModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyID := state.PolicyID.ValueString()
	client := r.OktaIDaaSClient.OktaSDKSupplementClient()

	for _, rule := range state.Rules {
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

	// Fetch each non-system rule individually to get full AccessPolicyRule details.
	importedRules := make([]policyRuleModel, 0, len(sdkRules))
	for _, sdkRule := range sdkRules {
		if sdkRule.System != nil && *sdkRule.System {
			continue // Skip system rules - they cannot be managed.
		}
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
		Rules:    importedRules,
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
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
					Blocks:     r.buildRuleBlocks(),
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
			// Note: We intentionally do NOT use UseStateForUnknown() here because
			// rules are matched by name, not by list position.
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
			Computed:    true,
			Description: "Priority of the rule. Lower numbers are evaluated first.",
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
			Default:     stringdefault.StaticString("ANYWHERE"),
			Description: "Network selection mode: ANYWHERE, ZONE, ON_NETWORK, or OFF_NETWORK.",
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
			Default:     stringdefault.StaticString("PT2H"),
			Description: "Re-authentication frequency in ISO 8601 duration format (e.g., PT2H for 2 hours).",
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
		"risk_score": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "Risk score level to match: ANY, LOW, MEDIUM, or HIGH.",
			Validators:  []validator.String{stringvalidator.OneOf(validRiskScores...)},
		},
	}
}

func (r *appSignOnPolicyRulesResource) buildRuleBlocks() map[string]schema.Block {
	return map[string]schema.Block{
		"platform_include": schema.ListNestedBlock{
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
						Optional:    true,
						Description: "Custom OS expression for advanced matching.",
					},
				},
			},
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

// buildPlannedNameByID creates a map of rule ID to planned name.
func (r *appSignOnPolicyRulesResource) buildPlannedNameByID(rules []policyRuleModel) map[string]string {
	m := make(map[string]string, len(rules))
	for _, rule := range rules {
		if !rule.ID.IsNull() && !rule.ID.IsUnknown() && rule.ID.ValueString() != "" {
			m[rule.ID.ValueString()] = rule.Name.ValueString()
		}
	}
	return m
}

// fetchExistingRulesByName lists all non-system rules currently in Okta for the policy
// and returns a map of rule name → rule ID. This is used in Create to detect rules
// that were partially created by an interrupted previous apply.
func (r *appSignOnPolicyRulesResource) fetchExistingRulesByName(ctx context.Context, client *sdk.APISupplement, policyID string) (map[string]string, diag.Diagnostics) {
	var diags diag.Diagnostics
	sdkRules, _, err := client.ListPolicyRules(ctx, policyID)
	if err != nil {
		diags.AddError("Error listing existing policy rules",
			fmt.Sprintf("Could not list rules for policy '%s': %s", policyID, err.Error()))
		return nil, diags
	}
	byName := make(map[string]string, len(sdkRules))
	for _, rule := range sdkRules {
		if rule.System != nil && *rule.System {
			continue // Never adopt system rules.
		}
		byName[rule.Name] = rule.Id
	}
	return byName, diags
}

// createOrAdoptRule creates a new rule or adopts an existing one if ID is specified.
func (r *appSignOnPolicyRulesResource) createOrAdoptRule(ctx context.Context, client *sdk.APISupplement, policyID string, rule policyRuleModel) (policyRuleModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	apiRule := r.buildAPIRuleFromModel(ctx, rule)

	// If ID is provided, adopt existing rule by updating it.
	if !rule.ID.IsNull() && !rule.ID.IsUnknown() && rule.ID.ValueString() != "" {
		updatedRule, err := r.updateRuleInAPI(ctx, client, policyID, rule.ID.ValueString(), apiRule)
		if err != nil {
			diags.AddError("Error adopting existing app sign-on policy rule",
				fmt.Sprintf("Could not adopt rule '%s' (ID: %s): %s", rule.Name.ValueString(), rule.ID.ValueString(), err.Error()))
			return policyRuleModel{}, diags
		}
		return r.updateRuleModelFromAPI(ctx, rule, updatedRule), diags
	}

	// Create new rule.
	createdRule, err := r.createRuleInAPI(ctx, client, policyID, apiRule)
	if err != nil {
		diags.AddError("Error creating app sign-on policy rule",
			fmt.Sprintf("Could not create rule '%s': %s", rule.Name.ValueString(), err.Error()))
		return policyRuleModel{}, diags
	}
	return r.updateRuleModelFromAPI(ctx, rule, createdRule), diags
}

// processRuleUpdate handles updating or creating a single rule during Update.
func (r *appSignOnPolicyRulesResource) processRuleUpdate(
	ctx context.Context,
	client *sdk.APISupplement,
	policyID string,
	planRule policyRuleModel,
	stateIndex *ruleIndex,
	nameTracker *nameTracker,
	plannedNameByID map[string]string,
) (policyRuleModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Find existing rule by ID or name.
	ruleID := ""
	if !planRule.ID.IsNull() && !planRule.ID.IsUnknown() {
		ruleID = planRule.ID.ValueString()
	}
	existingRule, exists := stateIndex.findRule(ruleID, planRule.Name.ValueString())

	if exists && !existingRule.ID.IsNull() {
		existingRuleID := existingRule.ID.ValueString()
		targetName := planRule.Name.ValueString()

		// Handle name conflicts (e.g., swapping names between rules).
		if err := r.resolveNameConflict(ctx, client, policyID, existingRuleID, targetName, nameTracker, plannedNameByID); err != nil {
			diags.AddError("Error updating app sign-on policy rule",
				fmt.Sprintf("Could not resolve name conflict: %s", err.Error()))
			return policyRuleModel{}, diags
		}

		// Update existing rule.
		apiRule := r.buildAPIRuleFromModel(ctx, planRule)
		updatedRule, err := r.updateRuleInAPI(ctx, client, policyID, existingRuleID, apiRule)
		if err != nil {
			diags.AddError("Error updating app sign-on policy rule",
				fmt.Sprintf("Could not update rule '%s': %s", planRule.Name.ValueString(), err.Error()))
			return policyRuleModel{}, diags
		}

		nameTracker.updateMapping(existingRuleID, targetName)
		return r.updateRuleModelFromAPI(ctx, planRule, updatedRule), diags
	}

	// Create new rule.
	apiRule := r.buildAPIRuleFromModel(ctx, planRule)
	createdRule, err := r.createRuleInAPI(ctx, client, policyID, apiRule)
	if err != nil {
		diags.AddError("Error creating app sign-on policy rule",
			fmt.Sprintf("Could not create rule '%s': %s", planRule.Name.ValueString(), err.Error()))
		return policyRuleModel{}, diags
	}

	if createdRule.Id != "" {
		nameTracker.updateMapping(createdRule.Id, planRule.Name.ValueString())
	}
	return r.updateRuleModelFromAPI(ctx, planRule, createdRule), diags
}

// resolveNameConflict handles cases where a rule's target name is held by another rule.
func (r *appSignOnPolicyRulesResource) resolveNameConflict(
	ctx context.Context,
	client *sdk.APISupplement,
	policyID, ruleID, targetName string,
	nameTracker *nameTracker,
	plannedNameByID map[string]string,
) error {
	conflictingID, hasConflict := nameTracker.hasConflict(targetName, ruleID)
	if !hasConflict {
		return nil
	}

	// Check if the conflicting rule will be renamed (part of a swap).
	plannedName, inPlan := plannedNameByID[conflictingID]
	if !inPlan || plannedName == targetName {
		return nil // No swap, let API handle the conflict.
	}

	// Rename conflicting rule to a temporary name.
	tempName := fmt.Sprintf("%s%s_%d", tempRuleNamePrefix, conflictingID, time.Now().UnixNano())
	if err := r.renameRuleToTemp(ctx, client, policyID, conflictingID, tempName); err != nil {
		return fmt.Errorf("failed to rename conflicting rule to temporary name: %w", err)
	}

	nameTracker.updateMapping(conflictingID, tempName)
	return nil
}

// renameRuleToTemp renames a rule to a temporary name while preserving all other fields.
func (r *appSignOnPolicyRulesResource) renameRuleToTemp(ctx context.Context, client *sdk.APISupplement, policyID, ruleID, tempName string) error {
	// Fetch current rule to preserve all fields.
	currentRule, err := r.readRuleFromAPI(ctx, client, policyID, ruleID)
	if err != nil {
		return fmt.Errorf("failed to read rule before renaming: %w", err)
	}

	currentRule.Name = tempName
	_, err = r.updateRuleInAPI(ctx, client, policyID, ruleID, *currentRule)
	return err
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
			platform.Os = &sdk.PlatformConditionEvaluatorPlatformOperatingSystem{
				Type: p.OsType.ValueString(),
			}
			if !p.OsExpression.IsNull() && !p.OsExpression.IsUnknown() {
				expr := p.OsExpression.ValueString()
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

	rule.FactorMode = types.StringValue(vm.FactorMode)
	rule.Type = types.StringValue(vm.Type)
	rule.ReAuthenticationFrequency = types.StringValue(vm.ReauthenticateIn)

	if vm.InactivityPeriod != "" {
		rule.InactivityPeriod = types.StringValue(vm.InactivityPeriod)
	}
}

func (r *appSignOnPolicyRulesResource) updateRuleConditionsFromAPI(ctx context.Context, rule *policyRuleModel, apiRule *sdk.AccessPolicyRule) {
	if apiRule.Conditions == nil {
		return
	}

	c := apiRule.Conditions

	// Network conditions
	if c.Network != nil {
		rule.NetworkConnection = types.StringValue(c.Network.Connection)
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
				platform.OsType = types.StringValue(p.Os.Type)
				if p.Os.Expression != nil && *p.Os.Expression != "" {
					platform.OsExpression = types.StringValue(*p.Os.Expression)
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
	// Initialize with empty collections to avoid nil pointer issues.
	emptySet, _ := types.SetValueFrom(ctx, types.StringType, []string{})
	emptyList, _ := types.ListValueFrom(ctx, types.StringType, []string{})

	rule := policyRuleModel{
		ID:                       types.StringValue(apiRule.Id),
		Name:                     types.StringValue(apiRule.Name),
		System:                   types.BoolValue(apiRule.System != nil && *apiRule.System),
		Status:                   types.StringValue(apiRule.Status),
		GroupsIncluded:           emptySet,
		GroupsExcluded:           emptySet,
		UsersIncluded:            emptySet,
		UsersExcluded:            emptySet,
		DeviceAssurancesIncluded: emptySet,
		UserTypesIncluded:        emptySet,
		UserTypesExcluded:        emptySet,
		NetworkIncludes:          emptyList,
		NetworkExcludes:          emptyList,
		Constraints:              emptyList,
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

	rule.FactorMode = types.StringValue(vm.FactorMode)
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
				platform.OsType = types.StringValue(p.Os.Type)
				if p.Os.Expression != nil && *p.Os.Expression != "" {
					platform.OsExpression = types.StringValue(*p.Os.Expression)
				}
			}
			platforms = append(platforms, platform)
		}
		rule.PlatformInclude = platforms
	}
}

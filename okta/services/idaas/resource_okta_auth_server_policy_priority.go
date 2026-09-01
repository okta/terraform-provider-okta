package idaas

import (
	"context"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
	"github.com/okta/terraform-provider-okta/sdk"
)

var (
	_ resource.Resource                = &authServerPolicyPriorityResource{}
	_ resource.ResourceWithConfigure   = &authServerPolicyPriorityResource{}
	_ resource.ResourceWithImportState = &authServerPolicyPriorityResource{}
)

func newAuthServerPolicyPriorityResource() resource.Resource {
	return &authServerPolicyPriorityResource{}
}

type authServerPolicyPriorityResource struct {
	*config.Config
}

type authServerPolicyPriorityModel struct {
	ID           types.String `tfsdk:"id"`
	AuthServerID types.String `tfsdk:"auth_server_id"`
	Priorities   types.List   `tfsdk:"priorities"`
}

func (r *authServerPolicyPriorityResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_auth_server_policy_priority"
}

func (r *authServerPolicyPriorityResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *authServerPolicyPriorityResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages the priority ordering of authorization server policies. Use this resource to guarantee consistent evaluation order and avoid race conditions when multiple okta_auth_server_policy resources are applied in parallel.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"auth_server_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the authorization server.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"priorities": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "Ordered list of authorization server policy IDs. The first entry receives priority 1, second receives priority 2, etc.",
			},
		},
	}
}

func (r *authServerPolicyPriorityResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan authServerPolicyPriorityModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var policyIDs []string
	resp.Diagnostics.Append(plan.Priorities.ElementsAs(ctx, &policyIDs, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if diags := r.validateNoDuplicates(policyIDs); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	authServerID := plan.AuthServerID.ValueString()
	resp.Diagnostics.Append(r.applyPriorities(ctx, authServerID, policyIDs)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = plan.AuthServerID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *authServerPolicyPriorityResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state authServerPolicyPriorityModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	authServerID := state.AuthServerID.ValueString()
	client := r.OktaIDaaSClient.OktaSDKClientV2()

	policies, _, err := client.AuthorizationServer.ListAuthorizationServerPolicies(ctx, authServerID)
	if err != nil {
		resp.Diagnostics.AddError("Error listing authorization server policies",
			fmt.Sprintf("Could not list policies for auth server '%s': %s", authServerID, err))
		return
	}

	// Build set of managed IDs from state.
	var managedIDs []string
	resp.Diagnostics.Append(state.Priorities.ElementsAs(ctx, &managedIDs, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	managedSet := make(map[string]bool, len(managedIDs))
	for _, id := range managedIDs {
		managedSet[id] = true
	}

	// Filter to managed policies only and sort by current priority.
	managed := make([]*sdk.AuthorizationServerPolicy, 0, len(managedIDs))
	for _, p := range policies {
		if managedSet[p.Id] {
			managed = append(managed, p)
		}
	}
	sortPoliciesByPriority(managed)

	// Silently drop any managed ID missing from API response (deleted outside Terraform).
	// This will surface as a plan diff on the next terraform plan.
	orderedIDs := make([]string, 0, len(managed))
	for _, p := range managed {
		orderedIDs = append(orderedIDs, p.Id)
	}

	var diags diag.Diagnostics
	state.Priorities, diags = types.ListValueFrom(ctx, types.StringType, orderedIDs)
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *authServerPolicyPriorityResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan authServerPolicyPriorityModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var policyIDs []string
	resp.Diagnostics.Append(plan.Priorities.ElementsAs(ctx, &policyIDs, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if diags := r.validateNoDuplicates(policyIDs); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	authServerID := plan.AuthServerID.ValueString()
	resp.Diagnostics.Append(r.applyPriorities(ctx, authServerID, policyIDs)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = plan.AuthServerID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *authServerPolicyPriorityResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// No-op: this resource manages ordering only; the policies themselves are not deleted.
}

func (r *authServerPolicyPriorityResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	authServerID := req.ID
	client := r.OktaIDaaSClient.OktaSDKClientV2()

	policies, _, err := client.AuthorizationServer.ListAuthorizationServerPolicies(ctx, authServerID)
	if err != nil {
		resp.Diagnostics.AddError("Error importing authorization server policy priority",
			fmt.Sprintf("Could not list policies for auth server '%s': %s", authServerID, err))
		return
	}

	sortPoliciesByPriority(policies)

	orderedIDs := make([]string, 0, len(policies))
	for _, p := range policies {
		orderedIDs = append(orderedIDs, p.Id)
	}

	priorities, diags := types.ListValueFrom(ctx, types.StringType, orderedIDs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := authServerPolicyPriorityModel{
		ID:           types.StringValue(authServerID),
		AuthServerID: types.StringValue(authServerID),
		Priorities:   priorities,
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// applyPriorities sets priorities bottom-up (highest index first) to avoid
// two policies sharing the same priority at any intermediate step.
func (r *authServerPolicyPriorityResource) applyPriorities(ctx context.Context, authServerID string, policyIDs []string) diag.Diagnostics {
	var diags diag.Diagnostics
	client := r.OktaIDaaSClient.OktaSDKClientV2()

	for i := len(policyIDs) - 1; i >= 0; i-- {
		policyID := policyIDs[i]
		priority := int64(i + 1)

		policy, _, err := client.AuthorizationServer.GetAuthorizationServerPolicy(ctx, authServerID, policyID)
		if err != nil {
			diags.AddError("Error reading authorization server policy",
				fmt.Sprintf("Could not read policy '%s': %s", policyID, err))
			return diags
		}

		// Skip update if priority is already correct
		if policy.PriorityPtr != nil && *policy.PriorityPtr == priority {
			continue
		}

		policy.PriorityPtr = &priority
		_, _, err = client.AuthorizationServer.UpdateAuthorizationServerPolicy(ctx, authServerID, policyID, *policy)
		if err != nil {
			diags.AddError("Error updating authorization server policy priority",
				fmt.Sprintf("Could not set priority %d on policy '%s': %s", priority, policyID, err))
			return diags
		}
	}
	return diags
}

// validateNoDuplicates returns an error diagnostic if any policy ID appears more than once.
func (r *authServerPolicyPriorityResource) validateNoDuplicates(policyIDs []string) diag.Diagnostics {
	var diags diag.Diagnostics
	seen := make(map[string]bool, len(policyIDs))
	for _, id := range policyIDs {
		if seen[id] {
			diags.AddError("Duplicate policy ID in priorities",
				fmt.Sprintf("Policy ID '%s' appears more than once in the priorities list.", id))
			return diags
		}
		seen[id] = true
	}
	return diags
}

// sortPoliciesByPriority sorts policies by PriorityPtr ascending.
// Policies with nil PriorityPtr are placed at the end.
func sortPoliciesByPriority(policies []*sdk.AuthorizationServerPolicy) {
	sort.Slice(policies, func(i, j int) bool {
		if policies[i].PriorityPtr == nil {
			return false
		}
		if policies[j].PriorityPtr == nil {
			return true
		}
		return *policies[i].PriorityPtr < *policies[j].PriorityPtr
	})
}

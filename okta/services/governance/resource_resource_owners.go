package governance

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/okta-governance-sdk-golang/governance"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ resource.Resource                = &resourceOwnersResource{}
	_ resource.ResourceWithConfigure   = &resourceOwnersResource{}
	_ resource.ResourceWithImportState = &resourceOwnersResource{}
)

var errResourceNotFound = errors.New("resource not found")

func newResourceOwnersResource() resource.Resource {
	return &resourceOwnersResource{}
}

type resourceOwnersResource struct {
	*config.Config
}

type resourceOwnersResourceModel struct {
	ID                types.String `tfsdk:"id"`
	ResourceOrn       types.String `tfsdk:"resource_orn"`
	ParentResourceOrn types.String `tfsdk:"parent_resource_orn"`
	PrincipalOrns     types.Set    `tfsdk:"principal_orns"`
}

func (r *resourceOwnersResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource_owners"
}

func (r *resourceOwnersResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages owners for a governance resource (entitlement bundle, entitlement value, collection, or app). " +
			"The resource is authoritative: any owners not declared in configuration will be removed.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Resource ID (same as resource_orn).",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"resource_orn": schema.StringAttribute{
				Description: "The ORN of the resource to manage owners for (e.g., an entitlement bundle, entitlement value, collection, or app).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"parent_resource_orn": schema.StringAttribute{
				Description: "The ORN of the parent resource (typically the app). " +
					"Automatically populated on create from the API response. " +
					"Required when importing. Can be found from the target_resource_orn attribute of okta_entitlement_bundle.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"principal_orns": schema.SetAttribute{
				Description: "The ORNs of the principals (users or groups) that own the resource. Maximum 5.",
				Required:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (r *resourceOwnersResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *resourceOwnersResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan resourceOwnersResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceOrn := plan.ResourceOrn.ValueString()
	principalOrns, diags := expandStringSet(ctx, plan.PrincipalOrns)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set ID early so partial failures still write state.
	plan.ID = types.StringValue(resourceOrn)

	configureReq := governance.ResourceOwnersUpdatable{
		ResourceOrns:  []string{resourceOrn},
		PrincipalOrns: principalOrns,
	}
	configureResp, _, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().ResourceOwnersAPI.
		ConfigureResourceOwners(ctx).
		ResourceOwnersUpdatable(configureReq).
		Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error configuring resource owners",
			fmt.Sprintf("Could not configure owners for resource %s: %s", resourceOrn, err.Error()),
		)
		// Don't write state on Create error — the resource may not exist,
		// and writing Unknown values (e.g., parent_resource_orn) would
		// cause a secondary framework error.
		return
	}

	// Capture parentResourceOrn from the response. We send exactly one
	// resourceOrn, so the first entry with a parentResourceOrn is ours.
	if configureResp != nil {
		for _, ro := range configureResp.GetData() {
			if pOrn := ro.GetParentResourceOrn(); pOrn != "" {
				plan.ParentResourceOrn = types.StringValue(pOrn)
				break
			}
		}
	}
	// Ensure parentResourceOrn is never Unknown in state (the framework
	// rejects Unknown values). This can happen if the user didn't set it
	// in config and the API response didn't include it (e.g., collections).
	if plan.ParentResourceOrn.IsUnknown() {
		plan.ParentResourceOrn = types.StringValue("")
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *resourceOwnersResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resourceOwnersResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceOrn := state.ResourceOrn.ValueString()
	parentOrn := state.ParentResourceOrn.ValueString()

	ro, err := r.findResourceOwner(ctx, resourceOrn, parentOrn)
	if err != nil {
		if errors.Is(err, errResourceNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading resource owners", err.Error())
		return
	}

	// Update state from API response
	state.ID = types.StringValue(resourceOrn)
	if ro.GetParentResourceOrn() != "" {
		state.ParentResourceOrn = types.StringValue(ro.GetParentResourceOrn())
	}

	principalOrns := make([]string, 0, len(ro.GetPrincipals()))
	for _, p := range ro.GetPrincipals() {
		principalOrns = append(principalOrns, p.GetOrn())
	}
	setVal, diags := flattenStringSet(principalOrns)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.PrincipalOrns = setVal

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *resourceOwnersResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state resourceOwnersResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceOrn := plan.ResourceOrn.ValueString()
	newPrincipalOrns, diags := expandStringSet(ctx, plan.PrincipalOrns)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use ConfigureResourceOwners for full replacement — simpler and more
	// reliable than PATCH since the API treats it as an authoritative set.
	configureReq := governance.ResourceOwnersUpdatable{
		ResourceOrns:  []string{resourceOrn},
		PrincipalOrns: newPrincipalOrns,
	}
	configureResp, _, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().ResourceOwnersAPI.
		ConfigureResourceOwners(ctx).
		ResourceOwnersUpdatable(configureReq).
		Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating resource owners",
			fmt.Sprintf("Could not update owners for resource %s: %s", resourceOrn, err.Error()),
		)
		plan.ID = state.ID
		plan.ParentResourceOrn = state.ParentResourceOrn
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		return
	}

	plan.ID = types.StringValue(resourceOrn)
	// Preserve or update parentResourceOrn from response
	plan.ParentResourceOrn = state.ParentResourceOrn
	if configureResp != nil {
		for _, ro := range configureResp.GetData() {
			if ro.GetParentResourceOrn() != "" {
				plan.ParentResourceOrn = types.StringValue(ro.GetParentResourceOrn())
				break
			}
		}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *resourceOwnersResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state resourceOwnersResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceOrn := state.ResourceOrn.ValueString()

	// Empty PrincipalOrns removes all owners.
	configureReq := governance.ResourceOwnersUpdatable{
		ResourceOrns:  []string{resourceOrn},
		PrincipalOrns: []string{},
	}
	_, apiResp, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().ResourceOwnersAPI.
		ConfigureResourceOwners(ctx).
		ResourceOwnersUpdatable(configureReq).
		Execute()
	if err != nil {
		if isHTTPNotFound(apiResp) {
			return
		}
		resp.Diagnostics.AddError(
			"Error removing resource owners",
			fmt.Sprintf("Could not remove owners for resource %s: %s", resourceOrn, err.Error()),
		)
	}
}

func (r *resourceOwnersResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import ID format: parent_resource_orn/resource_orn
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Import ID must be in the format: parent_resource_orn/resource_orn",
		)
		return
	}
	parentOrn := parts[0]
	resourceOrn := parts[1]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), resourceOrn)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("resource_orn"), resourceOrn)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("parent_resource_orn"), parentOrn)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Hydrate principal_orns from the API.
	ro, err := r.findResourceOwner(ctx, resourceOrn, parentOrn)
	if err != nil {
		resp.Diagnostics.AddError("Error reading resource owners during import", err.Error())
		return
	}

	principalOrns := make([]string, 0, len(ro.GetPrincipals()))
	for _, p := range ro.GetPrincipals() {
		principalOrns = append(principalOrns, p.GetOrn())
	}
	setVal, diags := flattenStringSet(principalOrns)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("principal_orns"), setVal)...)
}

// findResourceOwner looks up a specific resource owner entry by its ORN.
// It uses the parentResourceOrn filter combined with resource.orn when the parent is known,
// or falls back to resource.orn alone (works for collections and potentially other types).
func (r *resourceOwnersResource) findResourceOwner(ctx context.Context, resourceOrn, parentOrn string) (*governance.ResourceOwner, error) {
	var filter string
	if parentOrn != "" {
		filter = fmt.Sprintf("parentResourceOrn eq \"%s\" AND resource.orn eq \"%s\"", parentOrn, resourceOrn)
	} else {
		filter = fmt.Sprintf("resource.orn eq \"%s\"", resourceOrn)
	}

	listResp, apiResp, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().ResourceOwnersAPI.
		ListResourceOwners(ctx).
		Filter(filter).
		Execute()
	if err != nil {
		if isHTTPNotFound(apiResp) {
			return nil, errResourceNotFound
		}
		return nil, fmt.Errorf("error listing resource owners for %s: %w", resourceOrn, err)
	}

	// Find the matching resource in the response.
	data := listResp.GetData()
	for i := range data {
		res := data[i].GetResource()
		if res.GetOrn() == resourceOrn {
			return &data[i], nil
		}
	}

	return nil, errResourceNotFound
}

// isHTTPNotFound returns true if the API response indicates a 404.
func isHTTPNotFound(apiResp *governance.APIResponse) bool {
	return apiResp != nil && apiResp.Response != nil && apiResp.StatusCode == http.StatusNotFound
}

// expandStringSet converts a types.Set of strings to a []string.
func expandStringSet(ctx context.Context, set types.Set) ([]string, diag.Diagnostics) {
	if set.IsNull() || set.IsUnknown() {
		return []string{}, nil
	}
	var elements []types.String
	diags := set.ElementsAs(ctx, &elements, false)
	if diags.HasError() {
		return nil, diags
	}
	result := make([]string, 0, len(elements))
	for _, e := range elements {
		if !e.IsNull() && !e.IsUnknown() {
			result = append(result, e.ValueString())
		}
	}
	return result, nil
}

// flattenStringSet converts a []string to a types.Set of strings.
func flattenStringSet(values []string) (types.Set, diag.Diagnostics) {
	elements := make([]attr.Value, 0, len(values))
	for _, v := range values {
		elements = append(elements, types.StringValue(v))
	}
	return types.SetValue(types.StringType, elements)
}

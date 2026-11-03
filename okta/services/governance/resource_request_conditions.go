package governance

import (
	"context"
	"fmt"
	"net/http"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/okta-governance-sdk-golang/governance"
	"github.com/okta/terraform-provider-okta/okta/api"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ resource.Resource                   = &requestConditionsResource{}
	_ resource.ResourceWithConfigure      = &requestConditionsResource{}
	_ resource.ResourceWithImportState    = &requestConditionsResource{}
	_ resource.ResourceWithValidateConfig = &requestConditionsResource{}
)

func newRequestConditionsResource() resource.Resource {
	return &requestConditionsResource{}
}

// requestConditionsResource manages every request condition of a resource in
// a single Terraform resource. Individual okta_request_condition resources are
// created concurrently by Terraform, so the order in which they reach the API
// (and therefore their effective priority) is nondeterministic. This resource
// sorts the conditions by priority and sends the API calls one at a time.
//
// The API does not store the priority as given: it ranks the conditions of a
// resource 0, 1, 2, ... in order. The configured priority is therefore only
// used for ordering. Each condition is sent with its rank among the managed
// conditions, the configured value is kept in state, and drift is reported
// only when the order in Okta differs from the configured order.
type requestConditionsResource struct {
	*config.Config
}

type requestConditionsResourceModel struct {
	Id         types.String `tfsdk:"id"`
	ResourceId types.String `tfsdk:"resource_id"`
	Conditions types.List   `tfsdk:"condition"`
}

// requestConditionItemModel is a single condition block. It mirrors
// requestConditionResourceModel without resource_id, which is set once on
// the parent resource.
type requestConditionItemModel struct {
	Id                     types.String            `tfsdk:"id"`
	ApprovalSequenceId     types.String            `tfsdk:"approval_sequence_id"`
	Name                   types.String            `tfsdk:"name"`
	Description            types.String            `tfsdk:"description"`
	Priority               types.Int32             `tfsdk:"priority"`
	Created                types.String            `tfsdk:"created"`
	CreatedBy              types.String            `tfsdk:"created_by"`
	LastUpdated            types.String            `tfsdk:"last_updated"`
	LastUpdatedBy          types.String            `tfsdk:"last_updated_by"`
	Status                 types.String            `tfsdk:"status"`
	AccessScopeSettings    *Settings               `tfsdk:"access_scope_settings"`
	RequesterSettings      *Settings               `tfsdk:"requester_settings"`
	AccessDurationSettings *AccessDurationSettings `tfsdk:"access_duration_settings"`
}

// toResourceModel converts a condition block to the okta_request_condition
// model so the request builders of that resource can be reused. Unknown
// optional+computed values are turned into null so they are not sent to the
// API (an unknown priority would otherwise be sent as 0).
func (c requestConditionItemModel) toResourceModel(resourceId string) requestConditionResourceModel {
	if c.Priority.IsUnknown() {
		c.Priority = types.Int32Null()
	}
	if c.Status.IsUnknown() {
		c.Status = types.StringNull()
	}
	return requestConditionResourceModel{
		Id:                     c.Id,
		ResourceId:             types.StringValue(resourceId),
		ApprovalSequenceId:     c.ApprovalSequenceId,
		Name:                   c.Name,
		Description:            c.Description,
		Priority:               c.Priority,
		Created:                c.Created,
		CreatedBy:              c.CreatedBy,
		LastUpdated:            c.LastUpdated,
		LastUpdatedBy:          c.LastUpdatedBy,
		Status:                 c.Status,
		AccessScopeSettings:    c.AccessScopeSettings,
		RequesterSettings:      c.RequesterSettings,
		AccessDurationSettings: c.AccessDurationSettings,
	}
}

func requestConditionItemFromResourceModel(m requestConditionResourceModel) requestConditionItemModel {
	return requestConditionItemModel{
		Id:                     m.Id,
		ApprovalSequenceId:     m.ApprovalSequenceId,
		Name:                   m.Name,
		Description:            m.Description,
		Priority:               m.Priority,
		Created:                m.Created,
		CreatedBy:              m.CreatedBy,
		LastUpdated:            m.LastUpdated,
		LastUpdatedBy:          m.LastUpdatedBy,
		Status:                 m.Status,
		AccessScopeSettings:    m.AccessScopeSettings,
		RequesterSettings:      m.RequesterSettings,
		AccessDurationSettings: m.AccessDurationSettings,
	}
}

func (r *requestConditionsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_request_conditions"
}

func (r *requestConditionsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *requestConditionsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = requestConditionsSchema(ctx)
}

func requestConditionsSchema(ctx context.Context) schema.Schema {
	// Reuse the settings blocks of okta_request_condition so both resources
	// accept the same configuration for a condition.
	var single resource.SchemaResponse
	(&requestConditionResource{}).Schema(ctx, resource.SchemaRequest{}, &single)

	return schema.Schema{
		Description: "Manages all request conditions of a resource in a single configuration block. " +
			"Conditions are created and updated one at a time in priority order, so the resulting " +
			"priorities in Okta match the configuration.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of this resource (same as resource_id).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"resource_id": schema.StringAttribute{
				Required:    true,
				Description: "The id of the resource in Okta ID format.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"condition": schema.ListNestedBlock{
				Description: "List of request conditions. Conditions are processed in priority order (the smaller the number, the higher the priority).",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "The id of the request condition.",
						},
						"approval_sequence_id": schema.StringAttribute{
							Required:    true,
							Description: "The ID of the approval sequence.",
						},
						"name": schema.StringAttribute{
							Required:    true,
							Description: "The name of the request condition. Must be unique within the resource.",
						},
						"description": schema.StringAttribute{
							Optional:    true,
							Description: "The description of the request condition.",
						},
						"priority": schema.Int32Attribute{
							Optional: true,
							Computed: true,
							Description: "The priority of the condition. The smaller the number, the higher the priority. Must be unique within the resource. " +
								"Only the relative order matters: Okta ranks the conditions 0, 1, 2, ... in this order.",
						},
						"created":         schema.StringAttribute{Computed: true, Description: "The ISO 8601 formatted date and time when the resource was created."},
						"created_by":      schema.StringAttribute{Computed: true, Description: "The id of the Okta user who created the resource."},
						"last_updated":    schema.StringAttribute{Computed: true, Description: "The ISO 8601 formatted date and time when the object was last updated."},
						"last_updated_by": schema.StringAttribute{Computed: true, Description: "The id of the Okta user who last updated the object."},
						"status": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Status of the condition. Valid values: ACTIVE, INACTIVE. Default is INACTIVE.",
							Validators: []validator.String{
								stringvalidator.OneOf("ACTIVE", "INACTIVE"),
							},
						},
					},
					Blocks: single.Schema.Blocks,
				},
			},
		},
	}
}

// requestConditionItemType returns the object type of a single condition block.
func requestConditionItemType(ctx context.Context) attr.Type {
	return requestConditionsSchema(ctx).Blocks["condition"].Type().(types.ListType).ElemType
}

func (r *requestConditionsResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var list types.List
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("condition"), &list)...)
	if resp.Diagnostics.HasError() || list.IsNull() || list.IsUnknown() {
		return
	}
	var conditions []requestConditionItemModel
	resp.Diagnostics.Append(list.ElementsAs(ctx, &conditions, true)...)
	if resp.Diagnostics.HasError() {
		return
	}

	seenName := make(map[string]int, len(conditions))
	seenPriority := make(map[int32]int, len(conditions))
	for i, c := range conditions {
		conditionPath := path.Root("condition").AtListIndex(i)
		if !c.Name.IsNull() && !c.Name.IsUnknown() {
			name := c.Name.ValueString()
			if prev, exists := seenName[name]; exists {
				resp.Diagnostics.AddAttributeError(conditionPath.AtName("name"), "Duplicate request condition name",
					fmt.Sprintf("Name %q is used by both condition[%d] and condition[%d]. Each condition must have a unique name.", name, prev, i))
			} else {
				seenName[name] = i
			}
		}
		if !c.Priority.IsNull() && !c.Priority.IsUnknown() {
			p := c.Priority.ValueInt32()
			if prev, exists := seenPriority[p]; exists {
				resp.Diagnostics.AddAttributeError(conditionPath.AtName("priority"), "Duplicate request condition priority",
					fmt.Sprintf("Priority %d is used by both condition[%d] and condition[%d]. Each condition must have a unique priority.", p, prev, i))
			} else {
				seenPriority[p] = i
			}
		}
		if c.AccessScopeSettings == nil {
			resp.Diagnostics.AddAttributeError(conditionPath.AtName("access_scope_settings"), "Missing access_scope_settings",
				fmt.Sprintf("condition[%d] must define an access_scope_settings block.", i))
		}
		if c.RequesterSettings == nil {
			resp.Diagnostics.AddAttributeError(conditionPath.AtName("requester_settings"), "Missing requester_settings",
				fmt.Sprintf("condition[%d] must define a requester_settings block.", i))
		}
	}
}

func (r *requestConditionsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan requestConditionsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var conditions []requestConditionItemModel
	resp.Diagnostics.Append(plan.Conditions.ElementsAs(ctx, &conditions, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = context.WithValue(ctx, api.RetryOnStatusCodes, []int{http.StatusConflict})
	resourceId := plan.ResourceId.ValueString()
	plan.Id = plan.ResourceId

	// Conditions left behind by an interrupted apply have no state. Adopt
	// them by name instead of failing on the unique name constraint.
	existingByName, diags := r.listConditionIdsByName(ctx, resourceId)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	created := make([]requestConditionItemModel, 0, len(conditions))
	for rank, c := range sortRequestConditionsByPriority(conditions) {
		var result requestConditionItemModel
		if id, found := existingByName[c.Name.ValueString()]; found {
			result, diags = r.updateCondition(ctx, resourceId, id, c, int32(rank))
		} else {
			result, diags = r.createCondition(ctx, resourceId, c, int32(rank))
		}
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			// Save what was created so far. Terraform marks the resource as
			// tainted, so the next apply cleans up and recreates it.
			resp.Diagnostics.Append(setRequestConditions(ctx, &plan, reorderRequestConditionsToMatchPlan(created, conditions))...)
			resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
			return
		}
		created = append(created, result)
	}

	resp.Diagnostics.Append(setRequestConditions(ctx, &plan, reorderRequestConditionsToMatchPlan(created, conditions))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *requestConditionsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state requestConditionsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var conditions []requestConditionItemModel
	resp.Diagnostics.Append(state.Conditions.ElementsAs(ctx, &conditions, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceId := state.ResourceId.ValueString()
	refreshed := make([]requestConditionItemModel, 0, len(conditions))
	apiPriorities := make([]int32, 0, len(conditions))
	for _, c := range conditions {
		if c.Id.IsNull() || c.Id.IsUnknown() {
			continue
		}
		apiCondition, httpResp, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().RequestConditionsAPI.
			GetResourceRequestConditionV2(ctx, resourceId, c.Id.ValueString()).Execute()
		if err != nil {
			if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
				// Deleted outside Terraform: drop it so the next plan recreates it.
				continue
			}
			resp.Diagnostics.AddError("Error reading Request condition",
				fmt.Sprintf("Could not read request condition %q: %s", c.Name.ValueString(), err.Error()))
			return
		}
		result, diags := applyRequestConditionToItem(ctx, resourceId, c, apiCondition)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		refreshed = append(refreshed, result)
		apiPriorities = append(apiPriorities, apiCondition.GetPriority())
	}

	// State holds the configured priorities while Okta holds 0-based ranks.
	// Only when the order differs, show the ranks so the next plan reorders.
	if !requestConditionOrderMatches(refreshed, apiPriorities) {
		for i := range refreshed {
			refreshed[i].Priority = types.Int32Value(apiPriorities[i])
		}
	}

	resp.Diagnostics.Append(setRequestConditions(ctx, &state, refreshed)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *requestConditionsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state requestConditionsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var planConditions, stateConditions []requestConditionItemModel
	resp.Diagnostics.Append(plan.Conditions.ElementsAs(ctx, &planConditions, false)...)
	resp.Diagnostics.Append(state.Conditions.ElementsAs(ctx, &stateConditions, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = context.WithValue(ctx, api.RetryOnStatusCodes, []int{http.StatusConflict})
	resourceId := plan.ResourceId.ValueString()
	plan.Id = plan.ResourceId

	// Conditions are matched by name, which the API treats as a unique key.
	stateByName := make(map[string]requestConditionItemModel, len(stateConditions))
	for _, c := range stateConditions {
		stateByName[c.Name.ValueString()] = c
	}
	plannedNames := make(map[string]bool, len(planConditions))
	for _, c := range planConditions {
		plannedNames[c.Name.ValueString()] = true
	}

	// Delete removed conditions first so they don't take up priority slots.
	for _, c := range stateConditions {
		if plannedNames[c.Name.ValueString()] {
			continue
		}
		if err := r.deleteCondition(ctx, resourceId, c); err != nil {
			resp.Diagnostics.AddError("Error deleting Request condition",
				fmt.Sprintf("Could not delete request condition %q: %s", c.Name.ValueString(), err.Error()))
			r.saveUpdateProgress(ctx, resp, &plan, planConditions, stateByName, nil)
			return
		}
		delete(stateByName, c.Name.ValueString())
	}

	processed := make(map[string]requestConditionItemModel, len(planConditions))
	for rank, c := range sortRequestConditionsByPriority(planConditions) {
		var result requestConditionItemModel
		var diags diag.Diagnostics
		if existing, found := stateByName[c.Name.ValueString()]; found {
			result, diags = r.updateCondition(ctx, resourceId, existing.Id.ValueString(), c, int32(rank))
		} else {
			result, diags = r.createCondition(ctx, resourceId, c, int32(rank))
		}
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			r.saveUpdateProgress(ctx, resp, &plan, planConditions, stateByName, processed)
			return
		}
		processed[c.Name.ValueString()] = result
	}

	r.saveUpdateProgress(ctx, resp, &plan, planConditions, stateByName, processed)
}

// saveUpdateProgress writes the processed conditions in config order. A
// condition that was not processed yet keeps its prior state, and conditions
// that should have been deleted but still exist are kept as well, so a failed
// update does not lose track of conditions that still exist in Okta.
func (r *requestConditionsResource) saveUpdateProgress(ctx context.Context, resp *resource.UpdateResponse, plan *requestConditionsResourceModel,
	planConditions []requestConditionItemModel, stateByName, processed map[string]requestConditionItemModel,
) {
	result := make([]requestConditionItemModel, 0, len(planConditions))
	plannedNames := make(map[string]bool, len(planConditions))
	for _, c := range planConditions {
		name := c.Name.ValueString()
		plannedNames[name] = true
		if p, ok := processed[name]; ok {
			result = append(result, p)
		} else if s, ok := stateByName[name]; ok {
			result = append(result, s)
		}
	}
	for name, s := range stateByName {
		if !plannedNames[name] {
			result = append(result, s)
		}
	}
	resp.Diagnostics.Append(setRequestConditions(ctx, plan, result)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *requestConditionsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state requestConditionsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var conditions []requestConditionItemModel
	resp.Diagnostics.Append(state.Conditions.ElementsAs(ctx, &conditions, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceId := state.ResourceId.ValueString()
	for _, c := range conditions {
		if c.Id.IsNull() || c.Id.IsUnknown() {
			continue
		}
		if err := r.deleteCondition(ctx, resourceId, c); err != nil {
			resp.Diagnostics.AddError("Error deleting Request condition",
				fmt.Sprintf("Could not delete request condition %q: %s", c.Name.ValueString(), err.Error()))
			return
		}
	}
}

// ImportState imports every request condition of a resource. The import ID
// is the resource ID.
func (r *requestConditionsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceId := req.ID
	client := r.OktaGovernanceClient.OktaGovernanceSDKClient().RequestConditionsAPI
	list, _, err := client.ListResourceRequestConditionsV2(ctx, resourceId).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error importing Request conditions",
			fmt.Sprintf("Could not list request conditions for resource %q: %s", resourceId, err.Error()))
		return
	}

	imported := make([]requestConditionItemModel, 0, len(list.GetData()))
	for _, sparse := range list.GetData() {
		apiCondition, _, err := client.GetResourceRequestConditionV2(ctx, resourceId, sparse.GetId()).Execute()
		if err != nil {
			resp.Diagnostics.AddError("Error importing Request condition",
				fmt.Sprintf("Could not read request condition %q: %s", sparse.GetName(), err.Error()))
			return
		}
		result, diags := applyRequestConditionToItem(ctx, resourceId, requestConditionItemModel{}, apiCondition)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		imported = append(imported, result)
	}

	state := requestConditionsResourceModel{
		Id:         types.StringValue(resourceId),
		ResourceId: types.StringValue(resourceId),
	}
	resp.Diagnostics.Append(setRequestConditions(ctx, &state, sortRequestConditionsByPriority(imported))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// createCondition creates a condition at the given rank and activates it if
// requested.
func (r *requestConditionsResource) createCondition(ctx context.Context, resourceId string, c requestConditionItemModel, rank int32) (requestConditionItemModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	client := r.OktaGovernanceClient.OktaGovernanceSDKClient().RequestConditionsAPI
	data := c.toResourceModel(resourceId)
	data.Priority = types.Int32Value(rank)

	apiCondition, _, err := client.CreateResourceRequestConditionV2(ctx, resourceId).RequestConditionCreatable(createRequestCondition(data)).Execute()
	if err != nil {
		diags.AddError("Error creating Request condition",
			fmt.Sprintf("Could not create request condition %q: %s", c.Name.ValueString(), err.Error()))
		return c, diags
	}

	// The API may ignore the priority on create, so set it with a follow-up PATCH.
	if apiCondition.GetPriority() != rank {
		apiCondition, _, err = client.UpdateResourceRequestConditionV2(ctx, resourceId, apiCondition.GetId()).
			RequestConditionPatchable(createRequestConditionPatch(data)).Execute()
		if err != nil {
			diags.AddError("Error setting priority on Request condition",
				fmt.Sprintf("Could not update priority of request condition %q after creation: %s", c.Name.ValueString(), err.Error()))
			return c, diags
		}
	}

	if c.Status.ValueString() == string(governance.REQUESTCONDITIONSTATUS_ACTIVE) {
		apiCondition, _, err = client.ActivateResourceRequestConditionV2(ctx, resourceId, apiCondition.GetId()).Execute()
		if err != nil {
			diags.AddError("Error activating Request condition",
				fmt.Sprintf("Could not activate request condition %q after creation: %s", c.Name.ValueString(), err.Error()))
			return c, diags
		}
	}

	return applyRequestConditionToItem(ctx, resourceId, c, apiCondition)
}

// updateCondition patches an existing condition, moves it to the given rank
// and syncs its status.
func (r *requestConditionsResource) updateCondition(ctx context.Context, resourceId, id string, c requestConditionItemModel, rank int32) (requestConditionItemModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	client := r.OktaGovernanceClient.OktaGovernanceSDKClient().RequestConditionsAPI
	data := c.toResourceModel(resourceId)
	data.Priority = types.Int32Value(rank)

	apiCondition, _, err := client.UpdateResourceRequestConditionV2(ctx, resourceId, id).RequestConditionPatchable(createRequestConditionPatch(data)).Execute()
	if err != nil {
		diags.AddError("Error updating Request condition",
			fmt.Sprintf("Could not update request condition %q: %s", c.Name.ValueString(), err.Error()))
		return c, diags
	}

	if !c.Status.IsNull() && !c.Status.IsUnknown() && c.Status.ValueString() != string(apiCondition.GetStatus()) {
		switch c.Status.ValueString() {
		case string(governance.REQUESTCONDITIONSTATUS_ACTIVE):
			apiCondition, _, err = client.ActivateResourceRequestConditionV2(ctx, resourceId, id).Execute()
		case string(governance.REQUESTCONDITIONSTATUS_INACTIVE):
			apiCondition, _, err = client.DeactivateResourceRequestConditionV2(ctx, resourceId, id).Execute()
		}
		if err != nil {
			diags.AddError("Error setting status on Request condition",
				fmt.Sprintf("Could not set status of request condition %q: %s", c.Name.ValueString(), err.Error()))
			return c, diags
		}
	}

	return applyRequestConditionToItem(ctx, resourceId, c, apiCondition)
}

// deleteCondition deactivates (if needed) and deletes a condition. A condition
// that is already gone is treated as deleted.
func (r *requestConditionsResource) deleteCondition(ctx context.Context, resourceId string, c requestConditionItemModel) error {
	client := r.OktaGovernanceClient.OktaGovernanceSDKClient().RequestConditionsAPI
	id := c.Id.ValueString()

	// Request conditions must be INACTIVE before they can be deleted.
	if c.Status.ValueString() == string(governance.REQUESTCONDITIONSTATUS_ACTIVE) {
		_, httpResp, err := client.DeactivateResourceRequestConditionV2(ctx, resourceId, id).Execute()
		if err != nil {
			if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
				return nil
			}
			return fmt.Errorf("could not deactivate before deletion: %w", err)
		}
	}

	httpResp, err := client.DeleteResourceRequestConditionV2(ctx, resourceId, id).Execute()
	if err == nil || (httpResp != nil && httpResp.StatusCode == http.StatusNotFound) {
		return nil
	}
	if httpResp != nil && httpResp.StatusCode == http.StatusInternalServerError {
		// The governance API sometimes returns 500 even when the delete
		// succeeded. Treat the delete as successful if the condition is gone.
		_, verifyResp, _ := client.GetResourceRequestConditionV2(ctx, resourceId, id).Execute()
		if verifyResp != nil && verifyResp.StatusCode == http.StatusNotFound {
			return nil
		}
	}
	return err
}

// listConditionIdsByName returns the IDs of the conditions that already exist
// on the resource, keyed by name.
func (r *requestConditionsResource) listConditionIdsByName(ctx context.Context, resourceId string) (map[string]string, diag.Diagnostics) {
	var diags diag.Diagnostics
	list, _, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().RequestConditionsAPI.
		ListResourceRequestConditionsV2(ctx, resourceId).Execute()
	if err != nil {
		diags.AddError("Error listing Request conditions",
			fmt.Sprintf("Could not list request conditions for resource %q: %s", resourceId, err.Error()))
		return nil, diags
	}
	byName := make(map[string]string, len(list.GetData()))
	for _, c := range list.GetData() {
		byName[c.GetName()] = c.GetId()
	}
	return byName, diags
}

// applyRequestConditionToItem copies an API response onto a condition block.
// The API returns the rank of the condition rather than the configured
// priority, so a known priority is kept to avoid an inconsistent result after
// apply. Read compares the order separately.
func applyRequestConditionToItem(ctx context.Context, resourceId string, c requestConditionItemModel, apiCondition *governance.RequestConditionFull) (requestConditionItemModel, diag.Diagnostics) {
	knownPriority := c.Priority
	data := c.toResourceModel(resourceId)
	diags := applyRequestConditionToState(ctx, &data, apiCondition)
	if !knownPriority.IsNull() && !knownPriority.IsUnknown() {
		data.Priority = knownPriority
	}
	return requestConditionItemFromResourceModel(data), diags
}

// requestConditionOrderMatches reports whether sorting the conditions by
// their priority in state gives the same order as sorting them by their
// priority in Okta.
func requestConditionOrderMatches(conditions []requestConditionItemModel, apiPriorities []int32) bool {
	order := make([]int, len(conditions))
	for i, c := range conditions {
		if c.Priority.IsNull() || c.Priority.IsUnknown() {
			return false
		}
		order[i] = i
	}
	sort.SliceStable(order, func(i, j int) bool {
		return conditions[order[i]].Priority.ValueInt32() < conditions[order[j]].Priority.ValueInt32()
	})
	for k := 1; k < len(order); k++ {
		if apiPriorities[order[k-1]] >= apiPriorities[order[k]] {
			return false
		}
	}
	return true
}

func setRequestConditions(ctx context.Context, m *requestConditionsResourceModel, conditions []requestConditionItemModel) diag.Diagnostics {
	list, diags := types.ListValueFrom(ctx, requestConditionItemType(ctx), conditions)
	m.Conditions = list
	return diags
}

// sortRequestConditionsByPriority returns the conditions sorted by ascending
// priority. Conditions without a priority keep their relative order and are
// placed at the end.
func sortRequestConditionsByPriority(conditions []requestConditionItemModel) []requestConditionItemModel {
	sorted := make([]requestConditionItemModel, len(conditions))
	copy(sorted, conditions)
	sort.SliceStable(sorted, func(i, j int) bool {
		iPriority, jPriority := sorted[i].Priority, sorted[j].Priority
		if iPriority.IsNull() || iPriority.IsUnknown() {
			return false
		}
		if jPriority.IsNull() || jPriority.IsUnknown() {
			return true
		}
		return iPriority.ValueInt32() < jPriority.ValueInt32()
	})
	return sorted
}

// reorderRequestConditionsToMatchPlan returns the processed conditions in the
// order of the configuration, which Terraform expects state to follow.
func reorderRequestConditionsToMatchPlan(processed, planned []requestConditionItemModel) []requestConditionItemModel {
	byName := make(map[string]requestConditionItemModel, len(processed))
	for _, c := range processed {
		byName[c.Name.ValueString()] = c
	}
	result := make([]requestConditionItemModel, 0, len(planned))
	for _, c := range planned {
		if p, ok := byName[c.Name.ValueString()]; ok {
			result = append(result, p)
		}
	}
	return result
}

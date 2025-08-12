package governance

import (
	"context"
	"example.com/aditya-okta/okta-ig-sdk-golang/governance"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var _ resource.Resource = &riskRuleResource{}

func newRiskRuleResource() resource.Resource {
	return &riskRuleResource{}
}

type riskRuleResource struct {
	*config.Config
}

type riskRuleModel struct {
	Id               types.String            `tfsdk:"id"`
	Name             types.String            `tfsdk:"name"`
	Description      types.String            `tfsdk:"description"`
	Type             types.String            `tfsdk:"type"`
	Status           types.String            `tfsdk:"status"`
	Resources        []riskRuleResourceModel `tfsdk:"resources"`
	LastUpdated      types.String            `tfsdk:"last_updated"`
	Created          types.String            `tfsdk:"created"`
	CreatedBy        types.String            `tfsdk:"created_by"`
	LastUpdatedBy    types.String            `tfsdk:"last_updated_by"`
	ConflictCriteria *conflictCriteriaModel  `tfsdk:"conflict_criteria"`
	Notes            types.String            `tfsdk:"notes"`
}

type riskRuleResourceModel struct {
	ResourceOrn types.String `tfsdk:"resource_orn"`
}

type conflictCriteriaModel struct {
	And []conflictCriterionModel `tfsdk:"and"`
}

type conflictCriterionModel struct {
	Name      types.String        `tfsdk:"name"`
	Attribute types.String        `tfsdk:"attribute"`
	Operation types.String        `tfsdk:"operation"`
	Value     *conflictValueModel `tfsdk:"value"`
}

type conflictValueModel struct {
	Type  types.String           `tfsdk:"type"`
	Value []entitlementRuleModel `tfsdk:"value"`
}

type entitlementRuleModel struct {
	Id            types.String       `tfsdk:"id"`
	Name          types.String       `tfsdk:"name"`
	ExternalValue types.String       `tfsdk:"external_value"`
	Values        []entitlementValue `tfsdk:"values"`
}

type entitlementValue struct {
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	ExternalValue types.String `tfsdk:"external_value"`
}

func (r *riskRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_risk_rule"
}

func (r *riskRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Schema for a Separation of Duties Policy",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the policy.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description of the policy.",
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "Type of the policy, e.g., SEPARATION_OF_DUTIES.",
			},
			"created": schema.StringAttribute{
				Computed: true,
			},
			"created_by": schema.StringAttribute{
				Computed: true,
			},
			"last_updated_by": schema.StringAttribute{
				Computed: true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
			"notes": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"status": schema.StringAttribute{
				Computed: true,
			},
		},
		Blocks: map[string]schema.Block{
			"resources": schema.ListNestedBlock{
				Description: "List of resources this policy applies to.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"resource_orn": schema.StringAttribute{
							Required:    true,
							Description: "Object Resource Name of the resource.",
						},
					},
				},
			},
			"conflict_criteria": schema.SingleNestedBlock{
				Description: "Defines the logical conditions for policy conflict.",
				Blocks: map[string]schema.Block{
					"and": schema.ListNestedBlock{
						Description: "AND group of conditions.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required:    true,
									Description: "Name of the rule in the AND group.",
								},
								"attribute": schema.StringAttribute{
									Required:    true,
									Description: "Attribute path to check for conflict.",
								},
								"operation": schema.StringAttribute{
									Required:    true,
									Description: "Operation type, e.g., CONTAINS_ONE, CONTAINS_ALL.",
								},
							},
							Blocks: map[string]schema.Block{
								"value": schema.SingleNestedBlock{
									Description: "Values to match for the condition.",
									Attributes: map[string]schema.Attribute{
										"type": schema.StringAttribute{
											Required:    true,
											Description: "Type of the value, e.g., ENTITLEMENTS.",
										},
									},
									Blocks: map[string]schema.Block{
										"value": schema.ListNestedBlock{
											Description: "List of entitlements and their values.",
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"id": schema.StringAttribute{
														Required: true,
													},
													"name": schema.StringAttribute{
														Optional: true,
														Computed: true,
													},
													"external_value": schema.StringAttribute{
														Optional: true,
														Computed: true,
													},
												},
												Blocks: map[string]schema.Block{
													"values": schema.ListNestedBlock{
														NestedObject: schema.NestedBlockObject{
															Attributes: map[string]schema.Attribute{
																"id": schema.StringAttribute{
																	Required: true,
																},
																"name": schema.StringAttribute{
																	Optional: true,
																	Computed: true,
																},
																"external_value": schema.StringAttribute{
																	Optional: true,
																	Computed: true,
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *riskRuleResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
}

func (r *riskRuleResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(request, response)
}

func (r *riskRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data riskRuleModel
	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic
	createdRiskRule, _, err := r.OktaGovernanceClient.OktaIGSDKClient().RiskRulesAPI.CreateRiskRule(ctx).CreateRiskRuleRequest(buildRiskRule(data)).Execute()
	if err != nil {
		return
	}

	// Example data value setting
	applyToState(&data, createdRiskRule)
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *riskRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data riskRuleModel
	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	getRiskRuleResp, _, err := r.OktaGovernanceClient.OktaIGSDKClient().RiskRulesAPI.GetRiskRule(ctx, data.Id.ValueString()).Execute()
	if err != nil {
		return
	}
	applyToState(&data, getRiskRuleResp)
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *riskRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data riskRuleModel
	var state riskRuleModel
	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}
	data.Id = state.Id
	// Update API call logic
	fmt.Println("DAta ID", data.Id.ValueString())
	updateRiskRuleResp, _, err := r.OktaGovernanceClient.OktaIGSDKClient().RiskRulesAPI.ReplaceRiskRule(ctx, data.Id.ValueString()).UpdateRiskRuleRequest(buildUpdateRiskRule(data)).Execute()
	fmt.Println("Update Risk Rule Response:")
	if err != nil {
		return
	}

	applyToState(&data, updateRiskRuleResp)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func buildUpdateRiskRule(data riskRuleModel) governance.UpdateRiskRuleRequest {
	var r governance.UpdateRiskRuleRequest
	r.SetId(data.Id.ValueString())
	r.SetName(data.Name.ValueString())
	r.SetDescription(data.Description.ValueString())
	r.SetConflictCriteria(buildRiskRuleConflictCriteriaUpdatable(data.ConflictCriteria))
	r.SetNotes(data.Notes.ValueString())
	return r
}

func buildRiskRuleConflictCriteriaUpdatable(criteria *conflictCriteriaModel) governance.ConflictCriteriaUpdatable {
	var r governance.ConflictCriteriaUpdatable
	if criteria != nil {
		r.SetAnd(buildRiskRuleCriteria(criteria.And))
	}
	return r
}

func (r *riskRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data riskRuleModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	_, err := r.OktaGovernanceClient.OktaIGSDKClient().RiskRulesAPI.DeleteRiskRule(ctx, data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to delete risk rule",
			err.Error(),
		)
		return
	}
}

func applyToState(data *riskRuleModel, createdRiskRule *governance.RiskRuleResponse) {
	data.Id = types.StringValue(createdRiskRule.Id)
	data.Notes = types.StringPointerValue(createdRiskRule.Notes)
	data.Name = types.StringValue(createdRiskRule.Name)
	data.Description = types.StringPointerValue(createdRiskRule.Description)
	data.Type = types.StringValue(createdRiskRule.Type)
	data.Status = types.StringValue(string(createdRiskRule.Status))
	data.Resources = make([]riskRuleResourceModel, len(createdRiskRule.Resources))
	for i, resource := range createdRiskRule.Resources {
		data.Resources[i] = riskRuleResourceModel{
			ResourceOrn: types.StringPointerValue(resource.ResourceOrn),
		}
	}
	data.LastUpdated = types.StringValue(createdRiskRule.LastUpdated.Format(time.RFC3339))
	data.Created = types.StringValue(createdRiskRule.Created.Format(time.RFC3339))
	data.CreatedBy = types.StringValue(createdRiskRule.CreatedBy)
	data.LastUpdatedBy = types.StringValue(createdRiskRule.LastUpdatedBy)
	data.ConflictCriteria = &conflictCriteriaModel{
		And: make([]conflictCriterionModel, len(createdRiskRule.ConflictCriteria.And)),
	}

	for i, criterion := range createdRiskRule.ConflictCriteria.And {
		rule := conflictCriterionModel{
			Name:      types.StringPointerValue(criterion.Name),
			Attribute: types.StringPointerValue(criterion.Attribute),
			Operation: types.StringPointerValue(criterion.Operation),
		}

		if criterion.Value != nil {
			val := &conflictValueModel{
				Type:  types.StringPointerValue(criterion.Value.Type),
				Value: make([]entitlementRuleModel, len(criterion.Value.Value)),
			}

			for j, ent := range criterion.Value.Value {
				entitlement := entitlementRuleModel{
					Id:            types.StringValue(ent.Id),
					Name:          types.StringPointerValue(ent.Name),
					ExternalValue: types.StringPointerValue(ent.ExternalValue),
					Values:        make([]entitlementValue, len(ent.Values)),
				}

				for k, v := range ent.Values {
					entitlement.Values[k] = entitlementValue{
						Id:            types.StringPointerValue(v.Id),
						Name:          types.StringPointerValue(v.Name),
						ExternalValue: types.StringPointerValue(v.ExternalValue),
					}
				}

				val.Value[j] = entitlement
			}

			rule.Value = val
		}
		data.ConflictCriteria.And[i] = rule
	}
}

func buildRiskRule(data riskRuleModel) governance.CreateRiskRuleRequest {
	return governance.CreateRiskRuleRequest{
		Name:             data.Name.ValueString(),
		Description:      data.Description.ValueStringPointer(),
		Type:             data.Type.ValueString(),
		Resources:        buildRiskRuleResources(data.Resources),
		ConflictCriteria: buildRiskRuleConflictCriteria(data.ConflictCriteria),
		Notes:            data.Notes.ValueStringPointer(),
	}
}

func buildRiskRuleConflictCriteria(criteria *conflictCriteriaModel) governance.ConflictCriteriaCreatable {
	return governance.ConflictCriteriaCreatable{
		And: buildRiskRuleCriteria(criteria.And),
	}
}

func buildRiskRuleCriteria(and []conflictCriterionModel) []governance.CriteriaCreatable {
	result := make([]governance.CriteriaCreatable, len(and))
	for i, item := range and {
		result[i] = governance.CriteriaCreatable{
			Name:      item.Name.ValueStringPointer(),
			Attribute: item.Attribute.ValueStringPointer(),
			Operation: item.Operation.ValueStringPointer(),
			Value: &governance.CriteriaValueCreatable{
				Type:  item.Value.Type.ValueStringPointer(),
				Value: buildRiskRuleEntitlementMatch(item.Value.Value),
			},
		}
	}
	return result
}

func buildRiskRuleEntitlementMatch(value []entitlementRuleModel) []governance.EntitlementCreatable {
	var result []governance.EntitlementCreatable
	for _, item := range value {
		entitlement := governance.EntitlementCreatable{
			Id: item.Id.ValueStringPointer(),
		}
		for _, val := range item.Values {
			entitlement.Values = append(entitlement.Values, governance.EntitlementValueCreatable{
				Id: val.Id.ValueStringPointer(),
			})
		}
		result = append(result, entitlement)
	}
	return result
}

func buildRiskRuleResources(riskResources []riskRuleResourceModel) []governance.RuleConflictResource {
	var result []governance.RuleConflictResource
	for _, riskResource := range riskResources {
		result = append(result, governance.RuleConflictResource{
			ResourceOrn: riskResource.ResourceOrn.ValueStringPointer(),
		})
	}
	return result
}

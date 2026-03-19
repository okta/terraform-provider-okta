package governance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var _ datasource.DataSource = &EndUserMyRequestsDataSource{}

func newEndUserMyRequestsDataSource() datasource.DataSource {
	return &EndUserMyRequestsDataSource{}
}

type EndUserMyRequestsDataSource struct {
	*config.Config
}

type EndUserMyRequestsDataSourceModel struct {
	Id                   types.String `tfsdk:"id"`
	EntryId              types.String `tfsdk:"entry_id"`
	RequesterFieldValues types.List   `tfsdk:"requester_field_values"`
	Status               types.String `tfsdk:"status"`
	AccessDuration       types.String `tfsdk:"access_duration"`
	Created              types.String `tfsdk:"created"`
	CreatedBy            types.String `tfsdk:"created_by"`
	LastUpdated          types.String `tfsdk:"last_updated"`
	LastUpdatedBy        types.String `tfsdk:"last_updated_by"`
	Granted              types.String `tfsdk:"granted"`
	GrantStatus          types.String `tfsdk:"grant_status"`
	Requested            types.Object `tfsdk:"requested"`
	RequestedBy          types.Object `tfsdk:"requested_by"`
	RequestedFor         types.Object `tfsdk:"requested_for"`
	RiskAssessment       types.Object `tfsdk:"risk_assessment"`
	RevocationStatus     types.String `tfsdk:"revocation_status"`
	RevocationScheduled  types.String `tfsdk:"revocation_scheduled"`
	Revoked              types.String `tfsdk:"revoked"`
	Resolved             types.String `tfsdk:"resolved"`
}

func (r *EndUserMyRequestsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_end_user_my_requests"
}

func (d *EndUserMyRequestsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

func (r *EndUserMyRequestsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the request",
				Required:    true,
			},
			"entry_id": schema.StringAttribute{
				Description: "The ID of the catalog entry",
				Required:    true,
			},
			"status": schema.StringAttribute{
				Description: "The status of the request",
				Computed:    true,
			},
			"access_duration": schema.StringAttribute{
				Description: "How long the requester retains access after their request is approved and fulfilled.\nSpecified in ISO 8601 duration format.",
				Computed:    true,
			},
			"created": schema.StringAttribute{
				Description: "The ISO 8601 formatted date and time when the resource was created",
				Computed:    true,
			},
			"created_by": schema.StringAttribute{
				Description: "The id of the Okta user who created the resource",
				Computed:    true,
			},
			"last_updated": schema.StringAttribute{
				Description: "The ISO 8601 formatted date and time when the object was last updated",
				Computed:    true,
			},
			"last_updated_by": schema.StringAttribute{
				Description: "The id of the Okta user who last updated the object",
				Computed:    true,
			},
			"granted": schema.StringAttribute{
				Description: "The date the approved access was granted. Only set if request.status is APPROVED.",
				Computed:    true,
			},
			"grant_status": schema.StringAttribute{
				Description: "The grant status of the request",
				Computed:    true,
			},
			"revocation_scheduled": schema.StringAttribute{
				Description: "The date the granted access is scheduled for recovation. Only set if request.accessDuration exists, and request.grantStatus is GRANTED.",
				Computed:    true,
			},
			"revocation_status": schema.StringAttribute{
				Description: "The revocation status of the request, possible values are 'FAILED' 'PENDING' 'REVOKED'",
				Computed:    true,
			},
			"revoked": schema.StringAttribute{
				Description: "The date the granted access was revoked. Only set if request.grantStatus is GRANTED and request.revocationStatus is REVOKED.",
				Computed:    true,
			},
			"resolved": schema.StringAttribute{
				Description: "The date the request was resolved. The property may transition from having a value to null if the request is reopened.",
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"requested": schema.SingleNestedBlock{
				Description: "A representation of the resource in request",
				Attributes: map[string]schema.Attribute{
					"access_scope_id": schema.StringAttribute{
						Computed:    true,
						Description: "ID of the access scope",
					},
					"access_scope_type": schema.StringAttribute{
						Computed:    true,
						Description: "The access scope type",
					},
					"entry_id": schema.StringAttribute{
						Computed:    true,
						Description: "The ID of the resource catalog entry.",
					},
					"resource_id": schema.StringAttribute{
						Computed:    true,
						Description: "The requested resource ID",
					},
					"resource_type": schema.StringAttribute{
						Computed:    true,
						Description: "The requested resource type.",
					},
				},
			},
			"requested_by": schema.SingleNestedBlock{
				Description: "A representation of a principal",
				Attributes: map[string]schema.Attribute{
					"external_id": schema.StringAttribute{
						Computed:    true,
						Description: "The Okta user id",
					},
					"type": schema.StringAttribute{
						Computed:    true,
						Description: "The type of principal",
					},
				},
			},
			"requested_for": schema.SingleNestedBlock{
				Description: "A representation of a principal",
				Attributes: map[string]schema.Attribute{
					"external_id": schema.StringAttribute{
						Computed:    true,
						Description: "The Okta user id",
					},
					"type": schema.StringAttribute{
						Computed:    true,
						Description: "The type of principal",
					},
				},
			},
			"requester_field_values": schema.ListNestedBlock{
				Description: "The requester input fields required by the approval system.\nNote: The fields required are determined by the approval system.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Required:    true,
							Description: "The ID of a requesterField.",
						},
						"label": schema.StringAttribute{
							Optional:    true,
							Description: "A human-readable description of requesterField. It's used for display purposes and is optional",
						},
						"type": schema.StringAttribute{
							Description: "Type of value for the requester field.",
							Optional:    true,
							Validators: []validator.String{
								stringvalidator.OneOf("DURATION", "ISO_DATE", "MULTISELECT", "OKTA_USER_ID", "SELECT", "TEXT"),
							},
						},
						"value": schema.StringAttribute{
							Description: "The value of requesterField, which depends on the type of the field",
							Optional:    true,
						},
						"values": schema.ListAttribute{
							Description: "The values of requesterField with the type MULTISELECT.\nIf the field type is MULTISELECT, this property is required.",
							ElementType: types.StringType,
							Optional:    true,
						},
					},
				},
			},
			"risk_assessment": schema.SingleNestedBlock{
				Description: "A risk assessment indicates whether request submission is allowed or restricted and contains the risk rules that lead to possible conflicts for the requested resource.",
				Attributes: map[string]schema.Attribute{
					"request_submission_type": schema.StringAttribute{
						Optional:    true,
						Description: "Whether request submission is allowed or restricted in the risk settings.",
					},
				},
				Blocks: map[string]schema.Block{
					"risk_rules": schema.ListNestedBlock{
						Description: "An array of resources that are excluded from the review.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Optional:    true,
									Description: "The name of a resource rule causing a conflict.",
								},
								"description": schema.StringAttribute{
									Optional:    true,
									Description: "The human readable description.",
								},
								"resource_name": schema.StringAttribute{
									Optional:    true,
									Description: "Human readable name of the resource.",
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *EndUserMyRequestsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var stateData EndUserMyRequestsDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// step 1 : Retrieve End User's Request.
	getMyRequestV2Request := r.OktaGovernanceClient.OktaGovernanceSDKClient().MyRequestsAPI.GetMyRequestV2(ctx, stateData.EntryId.ValueString(), stateData.Id.ValueString())
	endUserMyRequest, _, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().MyRequestsAPI.GetMyRequestV2Execute(getMyRequestV2Request)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving End User My Request", "Could not retrieving End User My Request, unexpected error: "+err.Error())
		return
	}

	// step 2 : "Convert" API Compatible Type Back To Terraform Type.
	stateData.Id = types.StringValue(endUserMyRequest.GetId())
	stateData.EntryId = types.StringPointerValue(endUserMyRequest.GetRequested().EntryId)
	stateData.Status = types.StringPointerValue((*string)(endUserMyRequest.GetStatus().Ptr()))
	stateData.Created = types.StringValue(endUserMyRequest.GetCreated().String())
	stateData.CreatedBy = types.StringValue(endUserMyRequest.GetCreatedBy())
	stateData.LastUpdated = types.StringValue(endUserMyRequest.GetLastUpdated().String())
	stateData.LastUpdatedBy = types.StringValue(endUserMyRequest.GetLastUpdatedBy())
	stateData.GrantStatus = types.StringValue(string(*endUserMyRequest.GetGrantStatus().Ptr()))
	stateData.Granted = types.StringValue(endUserMyRequest.GetGranted().String())
	stateData.RevocationStatus = types.StringValue(string(endUserMyRequest.GetRevocationStatus()))
	stateData.RevocationScheduled = types.StringValue(string(endUserMyRequest.GetRevocationScheduled().String()))
	stateData.Revoked = types.StringValue(string(endUserMyRequest.GetRevoked().String()))
	stateData.Resolved = types.StringValue(string(endUserMyRequest.GetResolved().String()))
	requestedFieldsType := map[string]attr.Type{
		"access_scope_id":   types.StringType,
		"access_scope_type": types.StringType,
		"entry_id":          types.StringType,
		"resource_id":       types.StringType,
		"resource_type":     types.StringType,
	}
	requestedFieldsValue := map[string]attr.Value{
		"access_scope_id":   types.StringPointerValue(endUserMyRequest.GetRequested().AccessScopeId),
		"access_scope_type": types.StringPointerValue((*string)(endUserMyRequest.GetRequested().AccessScopeType.Ptr())),
		"entry_id":          types.StringPointerValue(endUserMyRequest.GetRequested().EntryId),
		"resource_id":       types.StringPointerValue(endUserMyRequest.GetRequested().ResourceId),
		"resource_type":     types.StringPointerValue((*string)(endUserMyRequest.GetRequested().ResourceType.Ptr())),
	}
	requested, diags := types.ObjectValue(requestedFieldsType, requestedFieldsValue)
	if diags != nil {
		resp.Diagnostics.AddError("Error retrieving End User My Request", "Could not retrieve Requested field: ")
		return
	}
	stateData.Requested = requested

	// fields for 'requested_by' and 'requested_for' objects are same
	requestedByForFieldsType := map[string]attr.Type{
		"external_id": types.StringType,
		"type":        types.StringType,
	}

	requestedByFieldsValue := map[string]attr.Value{
		"external_id": types.StringValue(endUserMyRequest.GetRequestedBy().ExternalId),
		"type":        types.StringValue(string(endUserMyRequest.GetRequestedBy().Type)),
	}
	requestedBy, diags := types.ObjectValue(requestedByForFieldsType, requestedByFieldsValue)
	if diags != nil {
		resp.Diagnostics.AddError("Error retrieving End User My Request", "Could not retrieve Requested By field: ")
		return
	}
	stateData.RequestedBy = requestedBy

	requestedForFieldsValue := map[string]attr.Value{
		"external_id": types.StringValue(endUserMyRequest.GetRequestedFor().ExternalId),
		"type":        types.StringValue(string(endUserMyRequest.GetRequestedFor().Type)),
	}
	requestedFor, diags := types.ObjectValue(requestedByForFieldsType, requestedForFieldsValue)
	if diags != nil {
		resp.Diagnostics.AddError("Error retrieving End User My Request", "Could not retrieve Requested For field: ")
		return
	}
	stateData.RequestedFor = requestedFor
	requesterFieldValuesType := map[string]attr.Type{
		"id":     types.StringType,
		"label":  types.StringType,
		"type":   types.StringType,
		"value":  types.StringType,
		"values": types.ListType{ElemType: types.StringType},
	}

	RequesterFieldValues := []attr.Value{}
	for _, requesterFieldValue := range endUserMyRequest.GetRequesterFieldValues() {
		fields := map[string]attr.Value{
			"id":     types.StringValue(requesterFieldValue.GetId()),
			"label":  types.StringNull(),
			"type":   types.StringNull(),
			"value":  types.StringNull(),
			"values": types.ListNull(types.StringType),
		}

		if labelField, ok := requesterFieldValue.GetLabelOk(); ok {
			fields["label"] = types.StringPointerValue(labelField)
		}

		if typeField, ok := requesterFieldValue.GetTypeOk(); ok {
			fields["type"] = types.StringPointerValue((*string)(typeField.Ptr()))
		}

		if valueField, ok := requesterFieldValue.GetValueOk(); ok {
			fields["value"] = types.StringPointerValue(valueField)
		}

		if valuesField, ok := requesterFieldValue.GetValuesOk(); ok {
			values := []attr.Value{}
			for _, value := range valuesField {
				values = append(values, types.StringValue(value))
			}
			fields["values"] = types.ListValueMust(types.StringType, values)
		}
		requesterFieldValue := types.ObjectValueMust(requesterFieldValuesType, fields)
		RequesterFieldValues = append(RequesterFieldValues, requesterFieldValue)
	}
	requesterFieldValuesListValue, diags := types.ListValue(types.ObjectType{AttrTypes: requesterFieldValuesType}, RequesterFieldValues)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	stateData.RequesterFieldValues = requesterFieldValuesListValue

	if riskAssessment, ok := endUserMyRequest.GetRiskAssessmentOk(); ok {
		// parse the risk rules
		riskRulesType := map[string]attr.Type{
			"name":          types.StringType,
			"description":   types.StringType,
			"resource_name": types.StringType,
		}

		riskRulesList := []attr.Value{}
		for _, riskRule := range riskAssessment.RiskRules {
			riskRuleEntry := map[string]attr.Value{
				"name":          types.StringValue(riskRule.GetName()),
				"description":   types.StringValue(riskRule.GetDescription()),
				"resource_name": types.StringValue(riskRule.GetResourceName()),
			}
			riskRuleObject, diags := types.ObjectValue(riskRulesType, riskRuleEntry)
			if diags.HasError() {
				resp.Diagnostics.Append(diags...)
				return
			}
			riskRulesList = append(riskRulesList, riskRuleObject)
		}

		riskRulesListValue, diags := types.ListValue(types.ObjectType{AttrTypes: riskRulesType}, riskRulesList)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}

		riskAssessmentFieldsType := map[string]attr.Type{
			"request_submission_type": types.StringType,
			"risk_rules":              types.ListType{ElemType: types.ObjectType{AttrTypes: riskRulesType}},
		}
		riskAssessmentFieldsValue := map[string]attr.Value{
			"request_submission_type": types.StringValue(string(riskAssessment.GetRequestSubmissionType())),
			"risk_rules":              riskRulesListValue,
		}
		riskAssessment, diags := types.ObjectValue(riskAssessmentFieldsType, riskAssessmentFieldsValue)
		if diags != nil {
			resp.Diagnostics.AddError("Error retrieving End User My Request", "Could not retrieve Requested field: ")
			return
		}
		stateData.RiskAssessment = riskAssessment
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

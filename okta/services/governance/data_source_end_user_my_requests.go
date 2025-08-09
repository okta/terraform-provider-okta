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
	RequestId            types.String `tfsdk:"request_id"`
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
			"request_id": schema.StringAttribute{
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
	getMyRequestV2Request := r.OktaGovernanceClient.OktaIGSDKClientV5().MyRequestsAPI.GetMyRequestV2(ctx, stateData.EntryId.ValueString(), stateData.RequestId.ValueString())
	endUserMyRequest, _, err := r.OktaGovernanceClient.OktaIGSDKClientV5().MyRequestsAPI.GetMyRequestV2Execute(getMyRequestV2Request)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving End User My Request", "Could not retrieving End User My Request, unexpected error: "+err.Error())
		return
	}

	// step 2 : "Convert" API Compatible Type Back To Terraform Type.
	stateData.RequestId = types.StringValue(endUserMyRequest.GetId())
	stateData.EntryId = types.StringPointerValue(endUserMyRequest.GetRequested().EntryId)
	stateData.Status = types.StringPointerValue((*string)(endUserMyRequest.GetStatus().Ptr()))
	stateData.Created = types.StringValue(endUserMyRequest.GetCreated().String())
	stateData.CreatedBy = types.StringValue(endUserMyRequest.GetCreatedBy())
	stateData.LastUpdated = types.StringValue(endUserMyRequest.GetLastUpdated().String())
	stateData.LastUpdatedBy = types.StringValue(endUserMyRequest.GetLastUpdatedBy())
	stateData.GrantStatus = types.StringValue(string(*endUserMyRequest.GetGrantStatus().Ptr()))
	stateData.Granted = types.StringValue(endUserMyRequest.GetGranted().String())
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
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

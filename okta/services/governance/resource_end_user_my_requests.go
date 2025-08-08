package governance

import (
	"context"

	"example.com/aditya-okta/okta-ig-sdk-golang/oktaInternalGovernance"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ resource.Resource                = &EndUserMyRequestsResource{}
	_ resource.ResourceWithConfigure   = &EndUserMyRequestsResource{}
	_ resource.ResourceWithImportState = &EndUserMyRequestsResource{}
)

func newMyRequestsResource() resource.Resource {
	return &EndUserMyRequestsResource{}
}

type EndUserMyRequestsResource struct {
	*config.Config
}

type EndUserMyRequestsResourceModel struct {
	Id                   types.String `tfsdk:"id"`
	EntryId              types.String `tfsdk:"entry_id"`
	RequesterFieldValues types.List   `tfsdk:"requester_field_values"`
	Status               types.String `tfsdk:"status"`
}

type RequesterFieldValue struct {
	Id     types.String `tfsdk:"id"`
	Label  types.String `tfsdk:"label"`
	Value  types.String `tfsdk:"value"`
	Values types.List   `tfsdk:"values"`
}

func (r *EndUserMyRequestsResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
}

func (r *EndUserMyRequestsResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(request, response)
}

func (r *EndUserMyRequestsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_end_user_my_requests"
}

func (r *EndUserMyRequestsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique identifier for the object",
				Computed:    true,
			},
			"entry_id": schema.StringAttribute{
				Description: "The ID of the catalog entry",
				Required:    true,
			},
			"status": schema.StringAttribute{
				Description: "The status of the request",
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{
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

func (r *EndUserMyRequestsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var stateData EndUserMyRequestsResourceModel // what our plan contains (assuming it read from state file)

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &stateData)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// step 1 : "Convert" Terraform Type to API Compatible Type.
	createMyRequestV2Request := r.OktaGovernanceClient.OktaIGSDKClientV5().MyRequestsAPI.CreateMyRequestV2(ctx, stateData.EntryId.ValueString())
	endUserMyRequest, _, err := r.OktaGovernanceClient.OktaIGSDKClientV5().MyRequestsAPI.CreateMyRequestV2Execute(createMyRequestV2Request.MyRequestCreatable(buildEndUserMyRequestCreateBody(ctx, stateData)))
	if err != nil {
		resp.Diagnostics.AddError("Error creating End User My Request", "Could not create End User My Request, unexpected error: "+err.Error())
		return
	}

	// step 2 : "Convert" API Compatible Type Back To Terraform Type.
	stateData.Id = types.StringValue(endUserMyRequest.GetId())
	stateData.EntryId = types.StringPointerValue(endUserMyRequest.GetRequested().EntryId)
	stateData.Status = types.StringPointerValue((*string)(endUserMyRequest.Status.Ptr()))

	if requesterFieldValues, ok := endUserMyRequest.GetRequesterFieldValuesOk(); ok && len(requesterFieldValues) > 0 {
		fieldTypes := map[string]attr.Type{
			"id":     types.StringType,
			"label":  types.StringType,
			"type":   types.StringType,
			"value":  types.StringType,
			"values": types.ListType{ElemType: types.StringType},
		}
		requesterFieldValuesList := []attr.Value{}
		for _, requesterFieldValue := range requesterFieldValues {
			fields := make(map[string]attr.Value)
			fields = map[string]attr.Value{
				"id": types.StringValue(requesterFieldValue.GetId()),
			}
			// Check before setting labe, type, value
			if label, ok := requesterFieldValue.GetLabelOk(); ok {
				fields["label"] = types.StringPointerValue(label)
			}
			if typ, ok := requesterFieldValue.GetTypeOk(); ok {
				fields["type"] = types.StringPointerValue((*string)(typ.Ptr()))
			}
			if value, ok := requesterFieldValue.GetValueOk(); ok {
				fields["value"] = types.StringPointerValue(value)
			}
			if valuesField, ok := requesterFieldValue.GetValuesOk(); ok {
				values := []attr.Value{}
				for _, value := range valuesField {
					values = append(values, types.StringValue(value))
				}
				valuesList, diags := types.ListValue(types.StringType, values)
				if diags.HasError() {
					resp.Diagnostics.Append(diags...)
					return
				}
				fields["values"] = valuesList
			} else {
				// // Initialize with empty list
				// emptyList, diags := types.ListValue(types.StringType, []attr.Value{})
				// if diags.HasError() {
				// 	resp.Diagnostics.Append(diags...)
				// 	return
				// }
				// fields["values"] = emptyList
			}
			fieldsObject := types.ObjectValueMust(fieldTypes, fields)
			requesterFieldValuesList = append(requesterFieldValuesList, fieldsObject)
		}

		// Set the requester field values in the state
		requesterFieldValuesListValue, diags := types.ListValue(types.ObjectType{
			AttrTypes: fieldTypes,
		}, requesterFieldValuesList)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		stateData.RequesterFieldValues = requesterFieldValuesListValue
	}
	// If API doesn't return requester field values, keep what was in the plan (stateData already has it from req.Plan.Get)
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *EndUserMyRequestsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var stateData EndUserMyRequestsResourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// step 1 : Retrieve End User's Request.
	getMyRequestV2Request := r.OktaGovernanceClient.OktaIGSDKClientV5().MyRequestsAPI.GetMyRequestV2(ctx, stateData.EntryId.ValueString(), stateData.Id.ValueString())
	endUserMyRequest, _, err := r.OktaGovernanceClient.OktaIGSDKClientV5().MyRequestsAPI.GetMyRequestV2Execute(getMyRequestV2Request)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving End User My Request", "Could not retrieving End User My Request, unexpected error: "+err.Error())
		return
	}

	// step 2 : "Convert" API Compatible Type Back To Terraform Type.
	stateData.Id = types.StringValue(endUserMyRequest.GetId())
	stateData.EntryId = types.StringPointerValue(endUserMyRequest.GetRequested().EntryId)
	stateData.Status = types.StringPointerValue((*string)(endUserMyRequest.GetStatus().Ptr()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *EndUserMyRequestsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddWarning(
		"Update Not Supported",
		"This resource does not support updates. Please cancel the request via the Okta Dashboard and recreate the resource instead.",
	)
}

func (r *EndUserMyRequestsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning(
		"Delete Not Supported",
		"This resource does not support deletions. Please cancel the request via the Okta Dashboard",
	)
}

func buildEndUserMyRequestCreateBody(ctx context.Context, data EndUserMyRequestsResourceModel) oktaInternalGovernance.MyRequestCreatable {
	if data.RequesterFieldValues.IsNull() || data.RequesterFieldValues.IsUnknown() {
		return oktaInternalGovernance.MyRequestCreatable{}
	}

	// Convert types.List to []RequesterFieldValue
	var requesterFieldValuesList []RequesterFieldValue
	diags := data.RequesterFieldValues.ElementsAs(ctx, &requesterFieldValuesList, false)
	if diags.HasError() {
		return oktaInternalGovernance.MyRequestCreatable{}
	}

	requesterFieldValues := []oktaInternalGovernance.RequestFieldValue{}
	for _, element := range requesterFieldValuesList {
		requestFieldValue := oktaInternalGovernance.RequestFieldValue{}
		requestFieldValue.Id = element.Id.ValueString()
		requestFieldValue.Label = element.Label.ValueStringPointer()
		requestFieldValue.Value = element.Value.ValueStringPointer()

		// Handle the values field for MULTISELECT type
		if !element.Values.IsNull() && !element.Values.IsUnknown() {
			var values []string
			diags := element.Values.ElementsAs(ctx, &values, false)
			if diags.HasError() {
				return oktaInternalGovernance.MyRequestCreatable{}
			}
			requestFieldValue.Values = values
		}
		requesterFieldValues = append(requesterFieldValues, requestFieldValue)
	}
	return oktaInternalGovernance.MyRequestCreatable{
		RequesterFieldValues: requesterFieldValues,
	}
}

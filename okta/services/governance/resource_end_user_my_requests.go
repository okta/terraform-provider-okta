package governance

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/okta-governance-sdk-golang/governance"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ resource.Resource                = &EndUserMyRequestsResource{}
	_ resource.ResourceWithConfigure   = &EndUserMyRequestsResource{}
	_ resource.ResourceWithImportState = &EndUserMyRequestsResource{}
)

func newEndUserMyRequestsResource() resource.Resource {
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

type RequesterFieldValueEntry struct {
	Id     types.String `tfsdk:"id"`
	Label  types.String `tfsdk:"label"`
	Type   types.String `tfsdk:"type"`
	Value  types.String `tfsdk:"value"`
	Values types.List   `tfsdk:"values"`
}

func (r *EndUserMyRequestsResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
	parts := strings.Split(request.ID, "/")
	if len(parts) != 2 {
		response.Diagnostics.AddError(
			"Invalid ID",
			"Expected format: request_id/entry_id",
		)
		return
	}
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("id"), parts[0])...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("entry_id"), parts[1])...)
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
				Description: "The id of the request",
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
	var stateData EndUserMyRequestsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &stateData)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// step 1 : "Convert" Terraform Type to API Compatible Type.
	createMyRequestV2Request := r.OktaGovernanceClient.OktaGovernanceSDKClient().MyRequestsAPI.CreateMyRequestV2(ctx, stateData.EntryId.ValueString())
	endRequestBody, diags := buildEndUserMyRequestCreateBody(ctx, stateData)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	endUserMyRequest, _, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().MyRequestsAPI.CreateMyRequestV2Execute(createMyRequestV2Request.MyRequestCreatable(endRequestBody))
	if err != nil {
		resp.Diagnostics.AddError("Error creating End User My Request", "Could not create End User My Request, unexpected error: "+err.Error())
		return
	}

	// step 2 : "Convert" API Compatible Type Back To Terraform Type.
	stateData.Id = types.StringValue(endUserMyRequest.GetId())
	stateData.EntryId = types.StringPointerValue(endUserMyRequest.GetRequested().EntryId)
	stateData.Status = types.StringPointerValue((*string)(endUserMyRequest.Status.Ptr()))

	if _, ok := endUserMyRequest.GetRequesterFieldValuesOk(); ok {
		// Define consistent field types that match the schema
		fieldTypes := map[string]attr.Type{
			"id":     types.StringType,
			"label":  types.StringType,
			"type":   types.StringType,
			"value":  types.StringType,
			"values": types.ListType{ElemType: types.StringType},
		}

		requesterFieldValues := []attr.Value{}
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
			fieldsObject := types.ObjectValueMust(fieldTypes, fields)
			requesterFieldValues = append(requesterFieldValues, fieldsObject)
		}
		RequesterFieldValues, diags := types.ListValue(types.ObjectType{AttrTypes: fieldTypes}, requesterFieldValues)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		stateData.RequesterFieldValues = RequesterFieldValues
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *EndUserMyRequestsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var stateData EndUserMyRequestsResourceModel
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
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *EndUserMyRequestsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddWarning(
		"Update Not Supported",
		"This resource does not support updates. Please cancel the request via the Okta Dashboard & recreate the resource instead.",
	)
}

func (r *EndUserMyRequestsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning(
		"Delete Not Supported",
		"This resource does not support deletions. Please cancel the request via the Okta Dashboard",
	)
}

func buildEndUserMyRequestCreateBody(ctx context.Context, data EndUserMyRequestsResourceModel) (governance.MyRequestCreatable, diag.Diagnostics) {
	if data.RequesterFieldValues.IsNull() || data.RequesterFieldValues.IsUnknown() {
		return governance.MyRequestCreatable{}, nil
	}

	var RequesterFieldValues []RequesterFieldValueEntry
	diags := data.RequesterFieldValues.ElementsAs(ctx, &RequesterFieldValues, true)
	if diags.HasError() {
		return governance.MyRequestCreatable{}, diags
	}

	requesterFieldValues := []governance.RequestFieldValue{}
	for _, element := range RequesterFieldValues {
		requestFieldValue := governance.RequestFieldValue{}
		requestFieldValue.Id = element.Id.ValueString()

		if !element.Label.IsNull() && !element.Label.IsUnknown() {
			requestFieldValue.Label = element.Label.ValueStringPointer()
		}
		if !element.Type.IsNull() && !element.Type.IsUnknown() {
			requestFieldValue.Type = (*governance.RequestFieldType)(element.Type.ValueStringPointer())
		}
		if !element.Value.IsNull() && !element.Value.IsUnknown() {
			requestFieldValue.Value = element.Value.ValueStringPointer()
		}
		if !element.Values.IsNull() && !element.Values.IsUnknown() {
			var values []string
			diags := element.Values.ElementsAs(ctx, &values, false)
			if diags.HasError() {
				return governance.MyRequestCreatable{}, diags
			}
			requestFieldValue.Values = values
		}
		requesterFieldValues = append(requesterFieldValues, requestFieldValue)
	}
	return governance.MyRequestCreatable{
		RequesterFieldValues: requesterFieldValues,
	}, nil
}

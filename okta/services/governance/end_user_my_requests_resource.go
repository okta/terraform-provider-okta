package governance

import (
	"context"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ resource.Resource              = &EndUserMyRequestsResource{}
	_ resource.ResourceWithConfigure = &EndUserMyRequestsResource{}
	// _ resource.ResourceWithImportState = &EndUserMyRequestsResource{}
)

func newMyRequestsResource() resource.Resource {
	return &EndUserMyRequestsResource{}
}

type EndUserMyRequestsResource struct {
	*config.Config
}

type EndUserMyRequestsResourceModel struct {
	RequestId            types.String `tfsdk:"request_id"`
	EntryId              types.String `tfsdk:"entry_id"`
	RequesterFieldValues types.List   `tfsdk:"requester_field_values"`
}

func (r *EndUserMyRequestsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_end_user_my_requests"
}

func (r *EndUserMyRequestsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"request_id": schema.StringAttribute{
				Description: "Unique identifier for the object",
				Computed:    true,
			},
			"entry_id": schema.StringAttribute{
				Description: "Creates a request for my catalog entry specified by entry_id",
				Required:    true,
			},
			"requester_field_values": schema.ListNestedAttribute{
				Description: "The requester input fields required by the approval system.",
				Required:    false,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Required:    true,
							Description: "The ID of a requesterField",
						},
						"label": schema.StringAttribute{
							Required:    false,
							Description: "A human-readable description of requester_field. It's used for display purposes and is optional",
						},
						"type": schema.StringAttribute{
							Required:    false,
							Description: "Type of value for the requester field.",
						},
						"value": schema.StringAttribute{
							Required:    true,
							Description: "A single string value for the field. Mutually exclusive with 'values'.",
						},
						"values": schema.ListAttribute{
							ElementType: types.StringType,
							Optional:    true,
							Description: "A list of string values for the field. Mutually exclusive with 'value'.",
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *EndUserMyRequestsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *EndUserMyRequestsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var stateData EndUserMyRequestsResourceModel // what our plan contains (assuming it read from state file)

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &stateData)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// step 1 : begin converting the terraform type to our api compatible type.
	fmt.Printf("stateData=%+v\n", spew.Sdump(stateData))
	createMyRequestV2Request := r.OktaGovernanceClient.OktaIGSDKClientV5().MyRequestsAPI.CreateMyRequestV2(ctx, stateData.EntryId.String())
	endUserMyRequest, _, err := r.OktaGovernanceClient.OktaIGSDKClientV5().MyRequestsAPI.CreateMyRequestV2Execute(createMyRequestV2Request)
	if err != nil {
		resp.Diagnostics.AddError("Error creating End User My Request", "Could not create End User My Request, unexpected error: "+err.Error())
		return
	}
	fmt.Printf("endUserMyRequest=%+v\n", spew.Sdump(endUserMyRequest))
	// step 1 : end

	// step 2 : converting the api compatible type back to terraform type
	stateData.RequestId = types.StringValue(endUserMyRequest.GetId())
	stateData.EntryId = types.StringPointerValue(endUserMyRequest.GetRequested().EntryId)

	if requesterFieldValues, ok := endUserMyRequest.GetRequesterFieldValuesOk(); ok {
		requesterFieldValuesList := []attr.Value{}
		// specify the type for each of the fields in requester_field_values
		fieldTypes := map[string]attr.Type{
			"id":     types.StringType,
			"label":  types.StringType,
			"type":   types.StringType,
			"value":  types.StringType,
			"values": types.ListType{ElemType: types.StringType},
		}

		// store the value for each of the elements in requester_field_values
		for _, requesterFieldValue := range requesterFieldValues {
			fields := make(map[string]attr.Value)
			fields = map[string]attr.Value{
				"id":    types.StringValue(requesterFieldValue.GetId()),
				"label": types.StringValue(requesterFieldValue.GetLabel()),
				"type":  types.StringValue(string(requesterFieldValue.GetType())),
			}
			// unsure if both value && values can be present together
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
					diags.Append(diags...)
					return
				}
				fields["value"] = valuesList
			}
			fieldsObject := types.ObjectValueMust(fieldTypes, fields)
			requesterFieldValuesList = append(requesterFieldValuesList, fieldsObject)
		}
		// list of objects, each object could very well be replaced by a map
		requesterFieldValues, diags := types.ListValue(types.ObjectType{AttrTypes: fieldTypes}, requesterFieldValuesList)
		if diags.HasError() {
			diags.Append(diags...)
			return
		}
		stateData.RequesterFieldValues = requesterFieldValues
	}
	// step 2 : end
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *EndUserMyRequestsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var stateData EndUserMyRequestsResourceModel // what our plan contains (assuming it read from state file)
	resp.Diagnostics.Append(resp.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// step 1 : retrieve the end user's my request
	createMyRequestV2Request := r.OktaGovernanceClient.OktaIGSDKClientV5().MyRequestsAPI.GetMyRequestV2(ctx, stateData.RequestId.String(), stateData.EntryId.String())
	endUserMyRequest, _, err := r.OktaGovernanceClient.OktaIGSDKClientV5().MyRequestsAPI.GetMyRequestV2Execute(createMyRequestV2Request)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving End User My Request", "Could not retrieving End User My Request, unexpected error: "+err.Error())
		return
	}
	// step 1 : end

	// step 2 : converting the api compatible type back to terraform type
	stateData.RequestId = types.StringValue(endUserMyRequest.GetId())
	stateData.EntryId = types.StringPointerValue(endUserMyRequest.GetRequested().EntryId)
	// step 2 : end

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *EndUserMyRequestsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// there is no update for this resource, we could just call create instead
}

func (r *EndUserMyRequestsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// there is no delete for this resource
}

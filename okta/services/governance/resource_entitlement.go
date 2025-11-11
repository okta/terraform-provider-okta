package governance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/okta-governance-sdk-golang/governance"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ resource.Resource                = &entitlementResource{}
	_ resource.ResourceWithConfigure   = &entitlementResource{}
	_ resource.ResourceWithImportState = &entitlementResource{}
)

func newEntitlementResource() resource.Resource {
	return &entitlementResource{}
}

type entitlementResource struct {
	*config.Config
}

type entitlementResourceModel struct {
	Id                types.String             `tfsdk:"id"`
	DataType          types.String             `tfsdk:"data_type"`
	ExternalValue     types.String             `tfsdk:"external_value"`
	MultiValue        types.Bool               `tfsdk:"multi_value"`
	Name              types.String             `tfsdk:"name"`
	Description       types.String             `tfsdk:"description"`
	Value             types.String             `tfsdk:"value"`
	Values            []entitlementValuesModel `tfsdk:"values"`
	Parent            *entitlementParentModel  `tfsdk:"parent"`
	ParentResourceOrn types.String             `tfsdk:"parent_resource_orn"`
}

type entitlementParentModel struct {
	ExternalID types.String `tfsdk:"external_id"`
	Type       types.String `tfsdk:"type"`
}

type entitlementValuesModel struct {
	Id            types.String `tfsdk:"id"`
	ExternalValue types.String `tfsdk:"external_value"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
}

func (r *entitlementResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entitlement"
}

func (r *entitlementResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The id property of an entitlement.",
			},
			"data_type": schema.StringAttribute{
				Required:    true,
				Description: "The data type of the entitlement property.",
			},
			"external_value": schema.StringAttribute{
				Required:    true,
				Description: "The value of an entitlement property.",
			},
			"multi_value": schema.BoolAttribute{
				Required:    true,
				Description: "The property that determines if the entitlement property can hold multiple values.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the entitlement property.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the entitlement property.",
			},
			"value": schema.StringAttribute{
				Optional:    true,
				Description: "The value of the entitlement property.",
			},
			"parent_resource_orn": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The Okta app instance, in ORN format.",
			},
		},
		Blocks: map[string]schema.Block{
			"parent": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"external_id": schema.StringAttribute{
						Required:    true,
						Description: "The Okta app.id of the resource.",
					},
					"type": schema.StringAttribute{
						Required:    true,
						Description: "The type of resource.",
					},
				},
				Description: "Representation of a resource.",
			},
			"values": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Collection of entitlement values.",
						},
						"external_value": schema.StringAttribute{
							Optional:    true,
							Description: "The value of an entitlement property value.",
						},
						"name": schema.StringAttribute{
							Optional:    true,
							Description: "The display name for an entitlement value.",
						},
						"description": schema.StringAttribute{
							Optional:    true,
							Description: "The description of the entitlement value.",
						},
					},
				},
			},
		},
	}
}

func (r *entitlementResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data entitlementResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.MultiValue.ValueBool() && data.DataType.ValueString() == "string" {
		resp.Diagnostics.AddAttributeError(
			path.Root("multi_value"),
			"Invalid Data Type Value",
			"Data type value should be array when multiValue is set to true.",
		)
		return
	}
	// Create API call logic
	createEntitlementResp, _, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().EntitlementsAPI.CreateEntitlement(ctx).EntitlementCreate(buildEntitlement(data)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Entitlement",
			"Could not create Entitlement, unexpected error: "+err.Error(),
		)
		return
	}

	applyEntitlementToState(&data, createEntitlementResp)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *entitlementResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data entitlementResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	readEntitlementResp, _, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().EntitlementsAPI.GetEntitlement(ctx, data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Entitlement",
			"Could not create Entitlement, unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(applyEntitlementToState(&data, readEntitlementResp)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *entitlementResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data entitlementResourceModel
	var state entitlementResourceModel
	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Id = state.Id
	// Update API call logic
	replaceEntitlementResp, _, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().EntitlementsAPI.ReplaceEntitlement(ctx, data.Id.ValueString()).EntitlementsFullWithParent(buildEntitlementReplace(data, state)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Entitlement",
			"An error occurred while updating the entitlement: "+err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(applyEntitlementToState(&data, replaceEntitlementResp)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *entitlementResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data entitlementResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	_, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().EntitlementsAPI.DeleteEntitlement(ctx, data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting entitlement",
			"Could not delete entitlement, unexpected error: "+err.Error(),
		)
		return
	}
}

func applyEntitlementToState(data *entitlementResourceModel, createEntitlementResp *governance.EntitlementsFullWithParent) diag.Diagnostics {
	data.Id = types.StringValue(createEntitlementResp.GetId())
	data.DataType = types.StringValue(string(createEntitlementResp.GetDataType()))
	data.ExternalValue = types.StringValue(createEntitlementResp.GetExternalValue())
	data.MultiValue = types.BoolValue(createEntitlementResp.GetMultiValue())
	data.Name = types.StringValue(createEntitlementResp.GetName())
	data.ParentResourceOrn = types.StringValue(createEntitlementResp.GetParentResourceOrn())
	if createEntitlementResp.Description != nil {
		data.Description = types.StringValue(createEntitlementResp.GetDescription())
	} else {
		data.Description = types.StringNull()
	}
	data.Parent = &entitlementParentModel{
		ExternalID: types.StringValue(createEntitlementResp.Parent.GetExternalId()),
		Type:       types.StringValue(string(createEntitlementResp.Parent.GetType())),
	}
	data.Values = make([]entitlementValuesModel, len(createEntitlementResp.GetValues()))
	for i, value := range createEntitlementResp.GetValues() {
		data.Values[i] = entitlementValuesModel{
			Id:            types.StringValue(value.GetId()),
			ExternalValue: types.StringValue(value.GetExternalValue()),
			Name:          types.StringValue(value.GetName()),
			Description:   types.StringValue(value.GetDescription()),
		}
	}
	return nil
}

func buildEntitlement(data entitlementResourceModel) governance.EntitlementCreate {
	return governance.EntitlementCreate{
		DataType:      governance.EntitlementPropertyDatatype(data.DataType.ValueString()),
		ExternalValue: data.ExternalValue.ValueString(),
		MultiValue:    data.MultiValue.ValueBool(),
		Name:          data.Name.ValueString(),
		Description:   data.Description.ValueStringPointer(),
		Parent: &governance.TargetResource{
			ExternalId: data.Parent.ExternalID.ValueString(),
			Type:       governance.ResourceType2(data.Parent.Type.ValueString()),
		},
		Values: buildEntitlementValues(data.Values),
	}
}

func buildEntitlementValues(data []entitlementValuesModel) []governance.EntitlementValueWritableProperties {
	x := make([]governance.EntitlementValueWritableProperties, len(data))
	for i, value := range data {
		x[i] = governance.EntitlementValueWritableProperties{
			ExternalValue: value.ExternalValue.ValueString(),
			Name:          value.Name.ValueString(),
			Description:   value.Description.ValueStringPointer(),
		}
	}
	return x
}

func buildEntitlementReplace(data, state entitlementResourceModel) governance.EntitlementsFullWithParent {
	return governance.EntitlementsFullWithParent{
		Id:            data.Id.ValueString(),
		DataType:      governance.EntitlementPropertyDatatype(data.DataType.ValueString()),
		ExternalValue: data.ExternalValue.ValueString(),
		MultiValue:    data.MultiValue.ValueBool(),
		Name:          data.Name.ValueString(),
		Description:   data.Description.ValueStringPointer(),
		Parent: governance.TargetResource{
			ExternalId: data.Parent.ExternalID.ValueString(),
			Type:       governance.ResourceType2(data.Parent.Type.ValueString()),
		},
		ParentResourceOrn: state.ParentResourceOrn.ValueString(),
		Values:            buildEntitlementValuesForReplace(data.Values, state.Values),
	}
}

func buildEntitlementValuesForReplace(values, state []entitlementValuesModel) []governance.EntitlementValueFull {
	x := make([]governance.EntitlementValueFull, len(values))
	for i, value := range values {
		var Id *string
		if i < len(state) {
			if !state[i].Id.IsNull() && state[i].Id.ValueString() != "" {
				Id = state[i].Id.ValueStringPointer()
			}
		}

		x[i] = governance.EntitlementValueFull{
			ExternalValue: value.ExternalValue.ValueStringPointer(),
			Name:          value.Name.ValueStringPointer(),
			Description:   value.Description.ValueStringPointer(),
			Id:            Id,
		}
	}
	return x
}

func (r *entitlementResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *entitlementResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(request, response)
}

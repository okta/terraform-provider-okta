package governance

import (
	"context"
	"example.com/aditya-okta/okta-ig-sdk-golang/oktaInternalGovernance"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

type entitlementUpdate struct {
	Op      types.String            `tfsdk:"op"`
	Path    types.String            `tfsdk:"path"`
	RefType types.String            `tfsdk:"ref_type"`
	Value   types.String            `tfsdk:"value"`
	Values  *entitlementValuesModel `tfsdk:"values"`
}

type entitlementResourceModel struct {
	Id            types.String             `tfsdk:"id"`
	DataType      types.String             `tfsdk:"data_type"`
	ExternalValue types.String             `tfsdk:"external_value"`
	MultiValue    types.Bool               `tfsdk:"multi_value"`
	Name          types.String             `tfsdk:"name"`
	Description   types.String             `tfsdk:"description"`
	Value         types.String             `tfsdk:"value"`
	Values        []entitlementValuesModel `tfsdk:"values"`
	Parent        *entitlementParentModel  `tfsdk:"parent"`
	//PatchOperations []entitlementUpdate      `tfsdk:"patch_operations"`
	ParentResourceOrn types.String `tfsdk:"parent_resource_orn"`
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
				Computed: true,
			},
			"data_type": schema.StringAttribute{
				Required: true,
			},
			"external_value": schema.StringAttribute{
				Required: true,
			},
			"multi_value": schema.BoolAttribute{
				Required: true,
			},
			"name": schema.StringAttribute{
				Optional: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"value": schema.StringAttribute{
				Optional: true,
			},
			"parent_resource_orn": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
		},
		Blocks: map[string]schema.Block{
			"parent": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"external_id": schema.StringAttribute{
						Required: true,
					},
					"type": schema.StringAttribute{
						Required: true,
					},
				},
			},
			"values": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Optional: true,
							Computed: true,
						},
						"external_value": schema.StringAttribute{
							Optional: true,
						},
						"name": schema.StringAttribute{
							Optional: true,
						},
						"description": schema.StringAttribute{
							Optional: true,
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
	createEntitlementResp, _, err := r.OktaGovernanceClient.OktaIGSDKClientV5().EntitlementsAPI.CreateEntitlement(ctx).EntitlementCreate(buildEntitlement(data)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Entitlement",
			"Could not create Entitlement, unexpected error: "+err.Error(),
		)
		return
	}

	r.applyEntitlementToState(&data, createEntitlementResp)
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
	readEntitlementResp, _, err := r.OktaGovernanceClient.OktaIGSDKClientV5().EntitlementsAPI.GetEntitlement(ctx, data.Id.ValueString()).Execute()
	if err != nil {
		return
	}

	resp.Diagnostics.Append(r.applyEntitlementToState(&data, readEntitlementResp)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError(
			"Unable to read entitlement",
			"Could not retrieve entitlement, unexpected error: "+err.Error(),
		)
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
	replaceEntitlementResp, _, err := r.OktaGovernanceClient.OktaIGSDKClientV5().EntitlementsAPI.ReplaceEntitlement(ctx, data.Id.ValueString()).EntitlementsFullWithParent(buildEntitlementReplace(data, state)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Entitlement",
			"An error occurred while updating the entitlement: "+err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(r.applyEntitlementToState(&data, replaceEntitlementResp)...)
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
	_, err := r.OktaGovernanceClient.OktaIGSDKClientV5().EntitlementsAPI.DeleteEntitlement(ctx, data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete entitlement",
			"Could not delete entitlement, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *entitlementResource) applyEntitlementToState(data *entitlementResourceModel, createEntitlementResp *oktaInternalGovernance.EntitlementsFullWithParent) diag.Diagnostics {
	data.Id = types.StringValue(createEntitlementResp.Id)
	data.DataType = types.StringValue(string(createEntitlementResp.DataType))
	data.ExternalValue = types.StringValue(createEntitlementResp.ExternalValue)
	data.MultiValue = types.BoolValue(createEntitlementResp.MultiValue)
	data.Name = types.StringValue(createEntitlementResp.Name)
	data.ParentResourceOrn = types.StringValue(createEntitlementResp.ParentResourceOrn)
	if createEntitlementResp.Description != nil {
		data.Description = types.StringValue(*createEntitlementResp.Description)
	} else {
		data.Description = types.StringNull()
	}
	data.Parent = &entitlementParentModel{
		ExternalID: types.StringValue(createEntitlementResp.Parent.ExternalId),
		Type:       types.StringValue(string(createEntitlementResp.Parent.Type)),
	}
	data.Values = make([]entitlementValuesModel, len(createEntitlementResp.Values))
	for i, value := range createEntitlementResp.Values {
		data.Values[i] = entitlementValuesModel{
			Id:            types.StringValue(*value.Id),
			ExternalValue: types.StringValue(*value.ExternalValue),
			Name:          types.StringValue(*value.Name),
			Description:   types.StringValue(*value.Description),
		}
	}
	return nil
}

func buildEntitlement(data entitlementResourceModel) oktaInternalGovernance.EntitlementCreate {
	return oktaInternalGovernance.EntitlementCreate{
		DataType:      oktaInternalGovernance.EntitlementPropertyDatatype(data.DataType.ValueString()),
		ExternalValue: data.ExternalValue.ValueString(),
		MultiValue:    data.MultiValue.ValueBool(),
		Name:          data.Name.ValueString(),
		Description:   data.Description.ValueStringPointer(),
		Parent: &oktaInternalGovernance.TargetResource{
			ExternalId: data.Parent.ExternalID.ValueString(),
			Type:       oktaInternalGovernance.ResourceType2(data.Parent.Type.ValueString()),
		},
		Values: buildEntitlementValues(data.Values),
	}
}

func buildEntitlementValues(data []entitlementValuesModel) []oktaInternalGovernance.EntitlementValueWritableProperties {
	x := make([]oktaInternalGovernance.EntitlementValueWritableProperties, len(data))
	for i, value := range data {
		x[i] = oktaInternalGovernance.EntitlementValueWritableProperties{
			ExternalValue: value.ExternalValue.ValueString(),
			Name:          value.Name.ValueString(),
			Description:   value.Description.ValueStringPointer(),
		}
	}
	return x
}

func buildEntitlementReplace(data, state entitlementResourceModel) oktaInternalGovernance.EntitlementsFullWithParent {
	return oktaInternalGovernance.EntitlementsFullWithParent{
		Id:            data.Id.ValueString(),
		DataType:      oktaInternalGovernance.EntitlementPropertyDatatype(data.DataType.ValueString()),
		ExternalValue: data.ExternalValue.ValueString(),
		MultiValue:    data.MultiValue.ValueBool(),
		Name:          data.Name.ValueString(),
		Description:   data.Description.ValueStringPointer(),
		Parent: oktaInternalGovernance.TargetResource{
			ExternalId: data.Parent.ExternalID.ValueString(),
			Type:       oktaInternalGovernance.ResourceType2(data.Parent.Type.ValueString()),
		},
		ParentResourceOrn: state.ParentResourceOrn.ValueString(),
		Values:            buildEntitlementValuesForReplace(data.Values, state.Values),
	}
}

func buildEntitlementValuesForReplace(values, state []entitlementValuesModel) []oktaInternalGovernance.EntitlementValueFull {
	x := make([]oktaInternalGovernance.EntitlementValueFull, len(values))
	for i, value := range values {
		var Id *string
		if i < len(state) {
			if !state[i].Id.IsNull() && state[i].Id.ValueString() != "" {
				Id = state[i].Id.ValueStringPointer()
			}
		}

		x[i] = oktaInternalGovernance.EntitlementValueFull{
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

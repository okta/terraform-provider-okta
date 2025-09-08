package governance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/okta-governance-sdk-golang/governance"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var _ datasource.DataSource = (*entitlementDataSource)(nil)

func newEntitlementDataSource() datasource.DataSource {
	return &entitlementDataSource{}
}

type entitlementDataSource struct {
	*config.Config
}

type parentBlockModel struct {
	ExternalId types.String `tfsdk:"external_id"`
	Type       types.String `tfsdk:"type"`
}

type entitlementValues struct {
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	ExternalValue types.String `tfsdk:"external_value"`
	ExternalId    types.String `tfsdk:"external_id"`
	Description   types.String `tfsdk:"description"`
}

type entitlementsDataSourceModel struct {
	Id                types.String        `tfsdk:"id"`
	DataType          types.String        `tfsdk:"data_type"`
	ExternalValue     types.String        `tfsdk:"external_value"`
	MultiValue        types.Bool          `tfsdk:"multi_value"`
	Name              types.String        `tfsdk:"name"`
	ParentResourceOrn types.String        `tfsdk:"parent_resource_orn"`
	Parent            *parentBlockModel   `tfsdk:"parent"`
	Values            []entitlementValues `tfsdk:"values"`
}

func (d *entitlementDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entitlement"
}

func (d *entitlementDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

func (d *entitlementDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required: true,
			},
			"data_type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of data for the entitlement, e.g., 'user', 'group', etc.",
			},
			"external_value": schema.StringAttribute{
				Computed:    true,
				Description: "The value of an entitlement property.",
			},
			"multi_value": schema.BoolAttribute{
				Computed:    true,
				Description: "The property that determines if the entitlement property can hold multiple values.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the entitlement property.",
			},
			"parent_resource_orn": schema.StringAttribute{
				Computed:    true,
				Description: "The Okta app instance, in ORN format.",
			},
		},
		Blocks: map[string]schema.Block{
			"parent": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"external_id": schema.StringAttribute{
						Computed:    true,
						Description: "The Okta app.id of the resource.",
					},
					"type": schema.StringAttribute{
						Computed:    true,
						Description: "The type of resource.",
					},
				},
				Description: "Representation of a resource.",
			},
			"values": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "The id of the entitlement value.",
						},
						"name": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "The name of the entitlement value.",
						},
						"external_value": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "The value of an entitlement property value.",
						},
						"external_id": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "The id of the entitlement value.",
						},
						"description": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "The description of an entitlement value.",
						},
					},
				},
			},
		},
	}
}

func (d *entitlementDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data entitlementsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	entitlementId := data.Id.ValueString()
	if entitlementId == "" {
		resp.Diagnostics.AddError("Missing entitlement Id", "The 'id' attribute must be set in the configuration.")
		return
	}
	readEntitlementResp, _, err := d.OktaGovernanceClient.OktaGovernanceSDKClient().EntitlementsAPI.GetEntitlement(ctx, entitlementId).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read entitlement",
			"Could not retrieve entitlement, unexpected error: "+err.Error(),
		)
		return
	}

	data = entitlementsDataSourceModel{
		Id:                types.StringValue(readEntitlementResp.GetId()),
		DataType:          types.StringValue(string(readEntitlementResp.GetDataType())),
		ExternalValue:     types.StringValue(readEntitlementResp.GetExternalValue()),
		MultiValue:        types.BoolValue(readEntitlementResp.GetMultiValue()),
		Name:              types.StringValue(readEntitlementResp.GetName()),
		ParentResourceOrn: types.StringValue(readEntitlementResp.GetParentResourceOrn()),
		Parent:            convertParent(&readEntitlementResp.Parent),
		Values:            convertValues(readEntitlementResp.GetValues()),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func convertValues(values []governance.EntitlementValueFull) []entitlementValues {
	var convertedValues []entitlementValues
	for _, value := range values {
		convertedValues = append(convertedValues, entitlementValues{
			Id:            types.StringValue(value.GetId()),
			Name:          types.StringValue(value.GetName()),
			ExternalValue: types.StringValue(value.GetExternalValue()),
			ExternalId:    types.StringValue(value.GetExternalId()),
			Description:   types.StringValue(value.GetDescription()),
		})
	}
	return convertedValues
}

func convertParent(parent *governance.TargetResource) *parentBlockModel {
	if parent == nil {
		return nil
	}
	return &parentBlockModel{
		ExternalId: types.StringValue(parent.GetExternalId()),
		Type:       types.StringValue(string(parent.GetType())),
	}
}

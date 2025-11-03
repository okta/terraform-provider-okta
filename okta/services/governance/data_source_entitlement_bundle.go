package governance

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var _ datasource.DataSource = (*entitlementBundleDataSource)(nil)

func newEntitlementBundleDataSource() datasource.DataSource {
	return &entitlementBundleDataSource{}
}

type entitlementBundleDataSource struct {
	*config.Config
}

type BundleEntitlementsDataSourceModel struct {
	Id            types.String        `tfsdk:"id"`
	DataType      types.String        `tfsdk:"data_type"`
	Description   types.String        `tfsdk:"description"`
	ExternalValue types.String        `tfsdk:"external_value"`
	MultiValue    types.Bool          `tfsdk:"multi_value"`
	Name          types.String        `tfsdk:"name"`
	Required      types.Bool          `tfsdk:"required"`
	Values        []entitlementValues `tfsdk:"values"`
}

type entitlementBundleDataSourceModel struct {
	Id                types.String                        `tfsdk:"id"`
	Description       types.String                        `tfsdk:"description"`
	Created           types.String                        `tfsdk:"created"`
	CreatedBy         types.String                        `tfsdk:"created_by"`
	LastUpdated       types.String                        `tfsdk:"last_updated"`
	LastUpdatedBy     types.String                        `tfsdk:"last_updated_by"`
	Name              types.String                        `tfsdk:"name"`
	Status            types.String                        `tfsdk:"status"`
	Target            *parentBlockModel                   `tfsdk:"target"`
	TargetResourceOrn types.String                        `tfsdk:"target_resource_orn"`
	Entitlements      []BundleEntitlementsDataSourceModel `tfsdk:"entitlements"`
}

func (d *entitlementBundleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entitlement_bundle"
}

func (d *entitlementBundleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

func (d *entitlementBundleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the Entitlement Bundle to retrieve.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "The human-readable description.",
			},
			"created": schema.StringAttribute{
				Computed:    true,
				Description: "The ISO 8601 formatted date and time when the resource was created.",
			},
			"created_by": schema.StringAttribute{
				Computed:    true,
				Description: "The id of the Okta user who created the resource.",
			},
			"last_updated": schema.StringAttribute{
				Computed:    true,
				Description: "The ISO 8601 formatted date and time when the object was last updated.",
			},
			"last_updated_by": schema.StringAttribute{
				Computed:    true,
				Description: "The id of the Okta user who last updated.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the entitlement bundle.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The lifecycle status of an entitlement bundle:.",
			},
			"target_resource_orn": schema.StringAttribute{
				Computed:    true,
				Description: "The Okta resource, in ORN format..",
			},
		},
		Blocks: map[string]schema.Block{
			"target": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"external_id": schema.StringAttribute{
						Computed:    true,
						Description: "The Okta app.id of the resource.",
					},
					"type": schema.StringAttribute{
						Computed:    true,
						Description: "The type of the resource.",
					},
				},
			},
			"entitlements": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "The id property of an entitlement.",
						},
						"data_type": schema.StringAttribute{
							Computed:    true,
							Description: "The data type of the entitlement property.",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "The description of the entitlement property.",
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
							Description: "The display name for an entitlement property.",
						},
						"required": schema.BoolAttribute{
							Computed:    true,
							Description: "The property that determines if the entitlement property is a required attribute.",
						},
					},
					Blocks: map[string]schema.Block{
						"values": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Computed:    true,
										Description: "The id of an entitlement value.",
									},
									"name": schema.StringAttribute{
										Computed:    true,
										Description: "The display name for an entitlement value.",
									},
									"external_value": schema.StringAttribute{
										Computed:    true,
										Description: "The value of an entitlement property value.",
									},
									"external_id": schema.StringAttribute{
										Computed:    true,
										Description: "The read-only id of an entitlement property value in the downstream application.",
									},
									"description": schema.StringAttribute{
										Computed:    true,
										Description: "The description of the entitlement value.",
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

func (d *entitlementBundleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data entitlementBundleDataSourceModel

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
	readEntitlementBundleResp, _, err := d.OktaGovernanceClient.OktaGovernanceSDKClient().EntitlementBundlesAPI.GetentitlementBundle(ctx, entitlementId).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read entitlement bundle",
			"Could not retrieve entitlement, unexpected error: "+err.Error(),
		)
		return
	}

	data = entitlementBundleDataSourceModel{
		Id:                types.StringValue(readEntitlementBundleResp.GetId()),
		Created:           types.StringValue(readEntitlementBundleResp.GetCreated().Format(time.RFC3339)),
		CreatedBy:         types.StringValue(readEntitlementBundleResp.GetCreatedBy()),
		LastUpdated:       types.StringValue(readEntitlementBundleResp.GetLastUpdated().Format(time.RFC3339)),
		LastUpdatedBy:     types.StringValue(readEntitlementBundleResp.GetLastUpdatedBy()),
		Status:            types.StringValue(string(readEntitlementBundleResp.GetStatus())),
		TargetResourceOrn: types.StringValue(readEntitlementBundleResp.GetTargetResourceOrn()),
		Name:              types.StringValue(readEntitlementBundleResp.GetName()),
		Description:       types.StringValue(readEntitlementBundleResp.GetDescription()),
	}

	if targetResource, ok := readEntitlementBundleResp.GetTargetOk(); ok {
		data.Target = &parentBlockModel{
			ExternalId: types.StringValue(targetResource.GetExternalId()),
			Type:       types.StringValue(string(targetResource.GetType())),
		}
	}

	for _, e := range readEntitlementBundleResp.GetEntitlements() {
		data.Entitlements = append(data.Entitlements, BundleEntitlementsDataSourceModel{
			Id:            types.StringValue(e.GetId()),
			Name:          types.StringValue(e.GetName()),
			ExternalValue: types.StringValue(e.GetExternalValue()),
			Description:   types.StringValue(e.GetDescription()),
			DataType:      types.StringValue(string(e.GetDataType())),
			MultiValue:    types.BoolValue(e.GetMultiValue()),
			Required:      types.BoolValue(e.GetRequired()),
			Values:        convertValues(e.GetValues()),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

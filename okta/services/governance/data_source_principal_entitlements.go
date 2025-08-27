package governance

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var _ datasource.DataSource = &principalEntitlementsDataSource{}

func newPrincipalEntitlementsDataSource() datasource.DataSource {
	return &principalEntitlementsDataSource{}
}

func (d *principalEntitlementsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

type principalEntitlementsDataSource struct {
	*config.Config
}

type principalEntitlementsDataSourceModel struct {
	Id              types.String                          `tfsdk:"id"`
	Parent          *parentModel                          `tfsdk:"parent"`           // Optional block input
	TargetPrincipal *principalModel                       `tfsdk:"target_principal"` // Optional block input
	Data            []principalEntitlementDataSourceModel `tfsdk:"data"`
}

type principalModel struct {
	ExternalId types.String `tfsdk:"external_id"`
	Type       types.String `tfsdk:"type"`
}

type parentModel struct {
	ExternalId types.String `tfsdk:"external_id"`
	Type       types.String `tfsdk:"type"`
}

type ValueModel struct {
	Id            types.String `tfsdk:"id"`
	ExternalValue types.String `tfsdk:"external_value"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
}

type principalEntitlementDataSourceModel struct {
	Id                 types.String    `tfsdk:"id"`
	Name               types.String    `tfsdk:"name"`
	ExternalValue      types.String    `tfsdk:"external_value"`
	Description        types.String    `tfsdk:"description"`
	MultiValue         types.Bool      `tfsdk:"multi_value"`
	Required           types.Bool      `tfsdk:"required"`
	DataType           types.String    `tfsdk:"data_type"`
	TargetPrincipalOrn types.String    `tfsdk:"target_principal_orn"`
	TargetPrincipal    *principalModel `tfsdk:"target_principal"`
	ParentResourceOrn  types.String    `tfsdk:"parent_resource_orn"`
	Parent             *parentModel    `tfsdk:"parent"`
	Values             []ValueModel    `tfsdk:"values"`
}

func (d *principalEntitlementsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_principal_entitlements"
}

func (d *principalEntitlementsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The internal identifier for this data source, required by Terraform to track state. This field does not exist in the Okta API response.",
			},
		},
		Blocks: map[string]schema.Block{
			"target_principal": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"external_id": schema.StringAttribute{
						Required: true,
					},
					"type": schema.StringAttribute{
						Required: true,
					},
				},
			},
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
			"data": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "The id property of an entitlement.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The display name for an entitlement property.",
						},
						"external_value": schema.StringAttribute{
							Computed:    true,
							Description: "The value of an entitlement property.",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "The description of an entitlement property.",
						},
						"multi_value": schema.BoolAttribute{
							Computed:    true,
							Description: "The property that determines if the entitlement property can hold multiple values.",
						},
						"required": schema.BoolAttribute{
							Computed:    true,
							Description: "The property that determines if the entitlement property is a required attribute",
						},
						"data_type": schema.StringAttribute{
							Computed:    true,
							Description: "The data type of the entitlement property.",
						},
						"target_principal_orn": schema.StringAttribute{
							Computed:    true,
							Description: "The Okta user id in ORN format.",
						},
						"parent_resource_orn": schema.StringAttribute{
							Computed:    true,
							Description: "The Okta app instance, in ORN format.",
						},
					},
					Blocks: map[string]schema.Block{
						"target_principal": schema.SingleNestedBlock{
							Attributes: map[string]schema.Attribute{
								"external_id": schema.StringAttribute{
									Computed:    true,
									Description: "The Okta user id.",
								},
								"type": schema.StringAttribute{
									Computed:    true,
									Description: "The type of principal.",
								},
							},
						},
						"parent": schema.SingleNestedBlock{
							Attributes: map[string]schema.Attribute{
								"external_id": schema.StringAttribute{
									Computed:    true,
									Description: "The Okta id of the resource.",
								},
								"type": schema.StringAttribute{
									Computed:    true,
									Description: "The type of the resource.",
								},
							},
						},
						"values": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Computed:    true,
										Description: "The id of an entitlement value.",
									},
									"external_value": schema.StringAttribute{
										Computed:    true,
										Description: "The value of an entitlement property value.",
									},
									"name": schema.StringAttribute{
										Computed:    true,
										Description: "The name of an entitlement value.",
									},
									"description": schema.StringAttribute{
										Computed:    true,
										Description: "The description of an entitlement property.",
									},
								},
							},
							Description: "Collection of entitlement values.",
						},
					},
				},
			},
		},
	}
}

func (d *principalEntitlementsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data principalEntitlementsDataSourceModel

	// Read Terraform configuration Data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	principalEntitlementsResp, _, err := d.OktaGovernanceClient.OktaGovernanceSDKClient().PrincipalEntitlementsAPI.GetPrincipalEntitlements(ctx).Filter(prepareFilter(data)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Principal Entitlements",
			"Could not read Principal Entitlements, unexpected error: "+err.Error(),
		)
		return
	}
	var entitlements []principalEntitlementDataSourceModel
	for _, item := range principalEntitlementsResp.Data {
		entitlement := principalEntitlementDataSourceModel{
			Id:                 types.StringPointerValue(item.Id),
			Name:               types.StringPointerValue(item.Name),
			ExternalValue:      types.StringPointerValue(item.ExternalValue),
			Description:        types.StringPointerValue(item.Description),
			MultiValue:         types.BoolPointerValue(item.MultiValue),
			Required:           types.BoolPointerValue(item.Required),
			DataType:           types.StringValue(string(*item.DataType)),
			TargetPrincipalOrn: types.StringPointerValue(item.TargetPrincipalOrn),
			ParentResourceOrn:  types.StringPointerValue(item.ParentResourceOrn),
		}

		if item.TargetPrincipal != nil {
			entitlement.TargetPrincipal = &principalModel{
				ExternalId: types.StringValue(item.TargetPrincipal.ExternalId),
				Type:       types.StringValue(string(item.TargetPrincipal.Type)),
			}
		}

		if item.Parent != nil {
			entitlement.Parent = &parentModel{
				ExternalId: types.StringValue(item.Parent.ExternalId),
				Type:       types.StringValue(string(item.Parent.Type)),
			}
		}

		for _, v := range item.Values {
			val := ValueModel{
				Id:            types.StringPointerValue(v.Id),
				ExternalValue: types.StringPointerValue(v.ExternalValue),
				Name:          types.StringPointerValue(v.Name),
				Description:   types.StringPointerValue(v.Description),
			}
			entitlement.Values = append(entitlement.Values, val)
		}

		entitlements = append(entitlements, entitlement)
	}

	// Set Data in model
	data.Data = entitlements
	data.Id = types.StringValue("principal-entitlements")
	// Save Data into state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Save Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func prepareFilter(d principalEntitlementsDataSourceModel) string {
	parentExternalId := d.Parent.ExternalId.ValueString()
	parentType := d.Parent.Type.ValueString()
	targetPrincipalExternalId := d.TargetPrincipal.ExternalId.ValueString()
	targetPrincipalType := d.TargetPrincipal.Type.ValueString()
	return fmt.Sprintf(
		`parent.externalId eq "%s" AND parent.type eq "%s" AND targetPrincipal.externalId eq "%s" AND targetPrincipal.type eq "%s"`,
		parentExternalId, // or extract from parent_resource_orn
		parentType,
		targetPrincipalExternalId, // or extract externalId from ORN
		targetPrincipalType,
	)
}

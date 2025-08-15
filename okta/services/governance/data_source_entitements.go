package governance

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
	"strings"
)

var _ datasource.DataSource = (*entitlementsDataSource)(nil)

func newEntitlementsDataSource() datasource.DataSource {
	return &entitlementsDataSource{}
}

type entitlementsDataSource struct {
	*config.Config
}

type EntitlementModel struct {
	Id                types.String `tfsdk:"id"`
	MultiValue        types.Bool   `tfsdk:"multi_value"`
	Name              types.String `tfsdk:"name"`
	ParentResourceOrn types.String `tfsdk:"parent_resource_orn"`
	Parent            parentModel  `tfsdk:"parent"`
}

type entitlementsDataSourceModel struct {
	ExternalId    types.String       `tfsdk:"external_id"`
	Type          types.String       `tfsdk:"type"`
	ResourceOrn   types.String       `tfsdk:"resource_orn"`
	Prefix        types.String       `tfsdk:"prefix"`
	Substring     types.String       `tfsdk:"substring"`
	EntitlementId types.List         `tfsdk:"entitlement_id"`
	ExternalValue types.String       `tfsdk:"external_value"`
	Data          []EntitlementModel `tfsdk:"data"`
}

func (d *entitlementsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entitlements"
}

func (d *entitlementsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

func (d *entitlementsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"external_id": schema.StringAttribute{
				Optional: true,
			},
			"type": schema.StringAttribute{
				Optional: true,
			},
			"resource_orn": schema.StringAttribute{
				Optional: true,
			},
			"prefix": schema.StringAttribute{
				Optional: true,
			},
			"substring": schema.StringAttribute{
				Optional: true,
			},
			"entitlement_id": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"external_value": schema.StringAttribute{
				Optional: true,
			},
		},
		Blocks: map[string]schema.Block{
			"data": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"multi_value": schema.BoolAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"parent_resource_orn": schema.StringAttribute{
							Computed: true,
						},
					},
					Blocks: map[string]schema.Block{
						"parent": schema.SingleNestedBlock{
							Attributes: map[string]schema.Attribute{
								"external_id": schema.StringAttribute{
									Computed: true,
								},
								"type": schema.StringAttribute{
									Computed: true,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *entitlementsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data entitlementsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	listEntitlementsResp, _, err := d.OktaGovernanceClient.OktaIGSDKClient().EntitlementsAPI.ListEntitlements(ctx).Filter(buildEntitlementFilter(data)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Entitlements",
			"Could not read Entitlements, unexpected error: "+err.Error(),
		)
		return
	}

	// Example data value setting
	var entitlementsList []EntitlementModel
	for _, entitlement := range listEntitlementsResp.Data {
		e := EntitlementModel{
			Id:                types.StringValue(entitlement.Id),
			MultiValue:        types.BoolValue(entitlement.GetMultiValue()),
			Name:              types.StringValue(entitlement.Name),
			ParentResourceOrn: types.StringValue(entitlement.ParentResourceOrn),
			Parent: parentModel{
				ExternalId: types.StringValue(entitlement.Parent.GetExternalId()),
				Type:       types.StringValue(string(entitlement.Parent.GetType())),
			},
		}
		entitlementsList = append(entitlementsList, e)
	}
	data.Data = entitlementsList

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func buildEntitlementFilter(data entitlementsDataSourceModel) string {
	if !data.ResourceOrn.IsNull() {
		return fmt.Sprintf("parentResourceOrn eq \"%s\"", data.ResourceOrn.ValueString())
	} else if !data.Prefix.IsNull() && !data.Type.IsNull() && !data.ExternalId.IsNull() {
		return fmt.Sprintf("parent.externalId eq \"%s\" AND parent.type eq \"%s\" AND name sw \"%s\"", data.ExternalId.ValueString(), data.Type.ValueString(), data.Prefix.ValueString())
	} else if !data.ExternalValue.IsNull() && !data.Type.IsNull() && !data.ExternalId.IsNull() {
		return fmt.Sprintf("parent.externalId eq \"%s\" AND parent.type eq \"%s\" AND externalValue eq \"%s\"", data.ExternalId.ValueString(), data.Type.ValueString(), data.ExternalValue.ValueString())
	} else if !data.EntitlementId.IsNull() {
		var ids []string
		if diags := data.EntitlementId.ElementsAs(context.Background(), &ids, false); diags.HasError() {
			// handle error
			return ""
		}

		var entitlementIds []string
		for _, id := range ids {
			entitlementIds = append(entitlementIds, fmt.Sprintf(`id eq "%s"`, id))
		}

		entitlementFilter := strings.Join(entitlementIds, " OR ")

		return fmt.Sprintf(
			`parent.externalId eq "%s" AND parent.type eq "APPLICATION" AND (%s)`,
			data.ExternalId.ValueString(),
			entitlementFilter,
		)
	} else if !data.Substring.IsNull() && !data.Type.IsNull() && !data.ExternalId.IsNull() {
		return fmt.Sprintf("parent.externalId eq \"%s\" AND parent.type eq \"%s\" AND name co \"%s\"", data.ExternalId.ValueString(), data.Type.ValueString(), data.Substring.ValueString())
	} else if !data.ExternalId.IsNull() && !data.Type.IsNull() {
		return fmt.Sprintf("parent.externalId eq \"%s\" AND parent.type eq \"%s\"", data.ExternalId.ValueString(), data.Type.ValueString())
	}

	return ""
}

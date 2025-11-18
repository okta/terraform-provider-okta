package governance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var _ datasource.DataSource = (*collectionResourceDataSource)(nil)

func newCollectionResourceDataSource() datasource.DataSource {
	return &collectionResourceDataSource{}
}

type collectionResourceDataSource struct {
	*config.Config
}

type collectionResourceDataSourceModel struct {
	CollectionId types.String                                   `tfsdk:"collection_id"`
	ResourceId   types.String                                   `tfsdk:"resource_id"`
	ResourceOrn  types.String                                   `tfsdk:"resource_orn"`
	Entitlements []collectionResourceDataSourceEntitlementModel `tfsdk:"entitlements"`
}

type collectionResourceDataSourceEntitlementModel struct {
	Id     types.String                                        `tfsdk:"id"`
	Values []collectionResourceDataSourceEntitlementValueModel `tfsdk:"values"`
}

type collectionResourceDataSourceEntitlementValueModel struct {
	Id types.String `tfsdk:"id"`
}

func (d *collectionResourceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_collection_resource"
}

func (d *collectionResourceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

func (d *collectionResourceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves information about a resource within a collection.",
		Attributes: map[string]schema.Attribute{
			"collection_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the collection.",
			},
			"resource_id": schema.StringAttribute{
				Required:    true,
				Description: "The resource ID within the collection.",
			},
			"resource_orn": schema.StringAttribute{
				Computed:    true,
				Description: "The ORN identifier for the app resource.",
			},
		},
		Blocks: map[string]schema.Block{
			"entitlements": schema.ListNestedBlock{
				Description: "List of entitlements with their values assigned to this resource.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "The entitlement ID.",
						},
					},
					Blocks: map[string]schema.Block{
						"values": schema.ListNestedBlock{
							Description: "List of entitlement value IDs.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Computed:    true,
										Description: "The entitlement value ID.",
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

func (d *collectionResourceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data collectionResourceDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	collectionId := data.CollectionId.ValueString()
	resourceId := data.ResourceId.ValueString()

	if collectionId == "" {
		resp.Diagnostics.AddError("Missing collection ID", "The 'collection_id' attribute must be set in the configuration.")
		return
	}
	if resourceId == "" {
		resp.Diagnostics.AddError("Missing resource ID", "The 'resource_id' attribute must be set in the configuration.")
		return
	}

	// Read API call
	resource, _, err := d.OktaGovernanceClient.OktaGovernanceSDKClient().CollectionsAPI.
		GetCollectionResource(ctx, collectionId, resourceId).
		Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read collection resource",
			"Could not retrieve collection resource, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response to state
	data.ResourceOrn = types.StringValue(resource.GetResourceOrn())
	if resource.HasResourceId() {
		data.ResourceId = types.StringValue(resource.GetResourceId())
	}

	// Map entitlements if present
	if resource.HasEntitlements() {
		entitlements := resource.GetEntitlements()
		data.Entitlements = make([]collectionResourceDataSourceEntitlementModel, len(entitlements))

		for i, ent := range entitlements {
			entModel := collectionResourceDataSourceEntitlementModel{
				Id: types.StringValue(ent.GetId()),
			}

			if ent.HasValues() {
				values := ent.GetValues()
				entModel.Values = make([]collectionResourceDataSourceEntitlementValueModel, len(values))
				for j, val := range values {
					entModel.Values[j] = collectionResourceDataSourceEntitlementValueModel{
						Id: types.StringValue(val.GetId()),
					}
				}
			}

			data.Entitlements[i] = entModel
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

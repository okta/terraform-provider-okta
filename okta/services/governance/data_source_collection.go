package governance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var _ datasource.DataSource = (*collectionDataSource)(nil)

func newCollectionDataSource() datasource.DataSource {
	return &collectionDataSource{}
}

type collectionDataSource struct {
	*config.Config
}

type collectionDataSourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Created     types.String `tfsdk:"created"`
	CreatedBy   types.String `tfsdk:"created_by"`
	LastUpdated types.String `tfsdk:"last_updated"`
	UpdatedBy   types.String `tfsdk:"updated_by"`
}

func (d *collectionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_collection"
}

func (d *collectionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

func (d *collectionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves information about a resource collection.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The unique identifier of the collection.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the resource collection.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "The human-readable description of the collection.",
			},
			"created": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the collection was created.",
			},
			"created_by": schema.StringAttribute{
				Computed:    true,
				Description: "User ID of who created the collection.",
			},
			"last_updated": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the collection was last updated.",
			},
			"updated_by": schema.StringAttribute{
				Computed:    true,
				Description: "User ID of who last updated the collection.",
			},
		},
	}
}

func (d *collectionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data collectionDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	collectionId := data.Id.ValueString()
	if collectionId == "" {
		resp.Diagnostics.AddError("Missing collection ID", "The 'id' attribute must be set in the configuration.")
		return
	}

	// Read API call
	collection, _, err := d.OktaGovernanceClient.OktaGovernanceSDKClient().CollectionsAPI.GetCollection(ctx, collectionId).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read collection",
			"Could not retrieve collection, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response to state
	data.Id = types.StringValue(collection.GetId())
	data.Name = types.StringValue(collection.GetName())

	if collection.HasDescription() {
		data.Description = types.StringValue(collection.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	data.Created = types.StringValue(collection.Created.String())
	data.CreatedBy = types.StringValue(collection.GetCreatedBy())
	data.LastUpdated = types.StringValue(collection.LastUpdated.String())
	data.UpdatedBy = types.StringValue(collection.GetLastUpdatedBy())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

package governance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var _ datasource.DataSource = &collectionDataSource{}

func newCollectionDataSource() datasource.DataSource {
	return &collectionDataSource{}
}

type collectionDataSource struct {
	*config.Config
}

type resourceCountsModel struct {
	Applications types.Int32 `tfsdk:"applications"`
}

type collectionCountsModel struct {
	PrincipalAssignmentCount types.Int32         `tfsdk:"principal_assignment_count"`
	ResourceCounts           resourceCountsModel `tfsdk:"resource_counts"`
}

type collectionDataSourceModel struct {
	Id            types.String          `tfsdk:"id"`
	Created       types.String          `tfsdk:"created"`
	CreatedBy     types.String          `tfsdk:"created_by"`
	LastUpdated   types.String          `tfsdk:"last_updated"`
	LastUpdatedBy types.String          `tfsdk:"last_updated_by"`
	Name          types.String          `tfsdk:"name"`
	Description   types.String          `tfsdk:"description"`
	Counts        collectionCountsModel `tfsdk:"counts"`
}

func (d *collectionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_collection"
}

func (d *collectionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"created": schema.StringAttribute{
				Computed: true,
			},
			"created_by": schema.StringAttribute{
				Computed: true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
			"last_updated_by": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Computed: true,
			},
			"description": schema.StringAttribute{
				Computed: true,
			},
		},
		Blocks: map[string]schema.Block{
			"counts": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"principal_assignment_count": schema.Int32Attribute{
						Computed: true,
					},
				},
				Blocks: map[string]schema.Block{
					"resource_counts": schema.SingleNestedBlock{
						Attributes: map[string]schema.Attribute{
							"applications": schema.Int32Attribute{
								Computed: true,
							},
						},
					},
				},
			},
		},
	}
}

func (d *collectionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data collectionDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	readCollectionResp, _, err := d.OktaGovernanceClient.OktaIGSDKClient().CollectionsAPI.GetCollection(ctx, data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Collections",
			"Could not read Collections, unexpected error: "+err.Error(),
		)
		return
	}

	// Example data value setting
	data.Id = types.StringValue(readCollectionResp.GetId())
	data.Created = types.StringValue(readCollectionResp.GetCreated().String())
	data.CreatedBy = types.StringValue(readCollectionResp.GetCreatedBy())
	data.LastUpdated = types.StringValue(readCollectionResp.GetLastUpdated().String())
	data.LastUpdatedBy = types.StringValue(readCollectionResp.GetLastUpdatedBy())
	data.Name = types.StringValue(readCollectionResp.GetName())
	data.Description = types.StringValue(readCollectionResp.GetDescription())
	if counts, ok := readCollectionResp.GetCountsOk(); ok {
		data.Counts.PrincipalAssignmentCount = types.Int32Value(counts.PrincipalAssignmentCount)
		data.Counts.ResourceCounts = resourceCountsModel{}
		if resourceCounts, ok := counts.GetResourceCountsOk(); ok {
			data.Counts.ResourceCounts.Applications = types.Int32Value(resourceCounts.Applications)
		}
	}
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *collectionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

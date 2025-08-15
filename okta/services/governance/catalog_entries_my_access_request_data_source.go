package governance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = (*catalogEntriesMyAccessRequestDataSource)(nil)

func NewCatalogEntriesMyAccessRequestDataSource() datasource.DataSource {
	return &catalogEntriesMyAccessRequestDataSource{}
}

type catalogEntriesMyAccessRequestDataSource struct{}

type catalogEntriesMyAccessRequestDataSourceModel struct {
	Id types.String `tfsdk:"id"`
}

func (d *catalogEntriesMyAccessRequestDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_catalog_entries_my_access_request"
}

func (d *catalogEntriesMyAccessRequestDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *catalogEntriesMyAccessRequestDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data catalogEntriesMyAccessRequestDataSourceModel

	// Read Terraform configuration Data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic

	// Example Data value setting
	data.Id = types.StringValue("example-id")

	// Save Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

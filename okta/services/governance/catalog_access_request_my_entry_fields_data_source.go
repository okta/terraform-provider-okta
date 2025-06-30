package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = (*catalogAccessRequestMyEntryFieldsDataSource)(nil)

func NewCatalogAccessRequestMyEntryFieldsDataSource() datasource.DataSource {
	return &catalogAccessRequestMyEntryFieldsDataSource{}
}

type catalogAccessRequestMyEntryFieldsDataSource struct{}

type catalogAccessRequestMyEntryFieldsDataSourceModel struct {
	Id types.String `tfsdk:"id"`
}

func (d *catalogAccessRequestMyEntryFieldsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_catalog_access_request_my_entry_fields"
}

func (d *catalogAccessRequestMyEntryFieldsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *catalogAccessRequestMyEntryFieldsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data catalogAccessRequestMyEntryFieldsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic

	// Example data value setting
	data.Id = types.StringValue("example-id")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

package governance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = (*catalogAccessRequestMyEntryUsersDataSource)(nil)

func NewCatalogAccessRequestMyEntryUsersDataSource() datasource.DataSource {
	return &catalogAccessRequestMyEntryUsersDataSource{}
}

type catalogAccessRequestMyEntryUsersDataSource struct{}

type catalogAccessRequestMyEntryUsersDataSourceModel struct {
	Id types.String `tfsdk:"id"`
}

func (d *catalogAccessRequestMyEntryUsersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_catalog_access_request_my_entry_users"
}

func (d *catalogAccessRequestMyEntryUsersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *catalogAccessRequestMyEntryUsersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data catalogAccessRequestMyEntryUsersDataSourceModel

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

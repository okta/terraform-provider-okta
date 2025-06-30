package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = (*catalogAccessRequestMyEntryUsersFieldsDataSource)(nil)

func NewCatalogAccessRequestMyEntryUsersFieldsDataSource() datasource.DataSource {
	return &catalogAccessRequestMyEntryUsersFieldsDataSource{}
}

type catalogAccessRequestMyEntryUsersFieldsDataSource struct{}

type catalogAccessRequestMyEntryUsersFieldsDataSourceModel struct {
	Id types.String `tfsdk:"id"`
}

func (d *catalogAccessRequestMyEntryUsersFieldsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_catalog_access_request_my_entry_users_fields"
}

func (d *catalogAccessRequestMyEntryUsersFieldsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *catalogAccessRequestMyEntryUsersFieldsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data catalogAccessRequestMyEntryUsersFieldsDataSourceModel

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

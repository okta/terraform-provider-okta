package governance

import (
	"context"
	"github.com/okta/terraform-provider-okta/okta/config"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &catalogEntryUserAccessRequestFieldsDataSource{}

func newCatalogEntryUserAccessRequestFieldsDataSource() datasource.DataSource {
	return &catalogEntryUserAccessRequestFieldsDataSource{}
}

func (d *catalogEntryUserAccessRequestFieldsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

type catalogEntryUserAccessRequestFieldsDataSource struct {
	*config.Config
}

type choices struct {
	Choice types.String `tfsdk:"choice"`
}

type requesterFields struct {
	Id types.String `tfsdk:"id"`
	//EntryId      types.String `tfsdk:"entry_Id"`
	Required     types.Bool   `tfsdk:"required"`
	Type         types.String `tfsdk:"type"`
	Choices      []choices    `tfsdk:"choices"`
	Label        types.String `tfsdk:"label"`
	MaximumValue types.String `tfsdk:"maximum_value"`
	ReadOnly     types.Bool   `tfsdk:"read_only"`
	Value        types.String `tfsdk:"value"`
}

type catalogEntryUserAccessRequestFieldsDataSourceModel struct {
	EntryId types.String      `tfsdk:"entry_id"`
	UserId  types.String      `tfsdk:"user_id"`
	Data    []requesterFields `tfsdk:"data"`
}

func (d *catalogEntryUserAccessRequestFieldsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_catalog_entry_user_access_request_fields"
}

func (d *catalogEntryUserAccessRequestFieldsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"entry_id": schema.StringAttribute{
				Required: true,
			},
			"user_id": schema.StringAttribute{
				Required: true,
			},
		},
		Blocks: map[string]schema.Block{
			"data": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"required": schema.BoolAttribute{
							Computed: true,
						},
						"type": schema.StringAttribute{
							Computed: true,
						},
						"label": schema.StringAttribute{
							Computed: true,
						},
						"maximum_value": schema.StringAttribute{
							Computed: true,
						},
						"read_only": schema.BoolAttribute{
							Computed: true,
						},
						"value": schema.StringAttribute{
							Computed: true,
						},
					},
					Blocks: map[string]schema.Block{
						"choices": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"choice": schema.StringAttribute{
										Computed: true,
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

func (d *catalogEntryUserAccessRequestFieldsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data catalogEntryUserAccessRequestFieldsDataSourceModel

	// Read Terraform configuration Data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	getUserRequestFieldResp, _, err := d.OktaGovernanceClient.OktaIGSDKClient().CatalogsAPI.GetCatalogEntryRequestFieldsV2(ctx, data.EntryId.ValueString(), data.UserId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Request Fields",
			"Could not read Request Fields, unexpected error: "+err.Error(),
		)
		return
	}

	for _, field := range getUserRequestFieldResp.GetData() {
		requesterField := requesterFields{
			Id: types.StringValue(field.GetId()),
			//EntryId:      data.EntryId,
			Required:     types.BoolValue(field.GetRequired()),
			Type:         types.StringValue(string(field.GetType())),
			Label:        types.StringValue(field.GetLabel()),
			MaximumValue: types.StringValue(field.GetMaximumValue()),
			ReadOnly:     types.BoolValue(field.GetReadOnly()),
			Value:        types.StringValue(field.GetValue()),
		}

		var c []choices
		for _, choice := range field.GetChoices() {
			c = append(c, choices{
				Choice: types.StringValue(choice),
			})
		}
		requesterField.Choices = c
		data.Data = append(data.Data, requesterField)
	}
	// Save Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

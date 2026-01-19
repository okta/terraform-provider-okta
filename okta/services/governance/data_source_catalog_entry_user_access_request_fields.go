package governance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
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
	Id           types.String `tfsdk:"id"`
	Required     types.Bool   `tfsdk:"required"`
	Type         types.String `tfsdk:"type"`
	Choices      []choices    `tfsdk:"choices"`
	Label        types.String `tfsdk:"label"`
	MaximumValue types.String `tfsdk:"maximum_value"`
	ReadOnly     types.Bool   `tfsdk:"read_only"`
	Value        types.String `tfsdk:"value"`
}

type catalogEntryUserAccessRequestFieldsDataSourceModel struct {
	Id      types.String      `tfsdk:"id"`
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
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The internal identifier for this data source, required by Terraform to track state. This field does not exist in the Okta API response.",
			},
			"entry_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the catalog entry.",
			},
			"user_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the user.",
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
							Computed:    true,
							Description: "Indicates whether a value to this field is required to advance the request.",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "Type of value for the requester field.",
						},
						"label": schema.StringAttribute{
							Computed:    true,
							Description: "Label of the requester field.",
						},
						"maximum_value": schema.StringAttribute{
							Computed:    true,
							Description: "The maximum value allowed for this field. Only applies to DURATION fields.",
						},
						"read_only": schema.BoolAttribute{
							Computed:    true,
							Description: "Indicates this field is immutable.",
						},
						"value": schema.StringAttribute{
							Computed:    true,
							Description: "An admin configured value for this field. Only applies to DURATION fields.",
						},
					},
					Blocks: map[string]schema.Block{
						"choices": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"choice": schema.StringAttribute{
										Computed:    true,
										Description: "Valid choice.",
									},
								},
							},
							Description: "Valid choices when type is SELECT or MULTISELECT.",
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
	getUserRequestFieldResp, _, err := d.OktaGovernanceClient.OktaGovernanceSDKClient().CatalogsAPI.GetCatalogEntryRequestFieldsV2(ctx, data.EntryId.ValueString(), data.UserId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Request Fields",
			"Could not read Request Fields, unexpected error: "+err.Error(),
		)
		return
	}

	for _, field := range getUserRequestFieldResp.GetData() {
		requesterField := requesterFields{
			Id:           types.StringValue(field.GetId()),
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
	data.Id = types.StringValue(data.EntryId.ValueString() + "-" + data.UserId.ValueString())
	// Save Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

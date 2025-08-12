package governance

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/okta/terraform-provider-okta/okta/config"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = (*requestSequenceDataSource)(nil)

func newRequestSequencesDataSource() datasource.DataSource {
	return &requestSequenceDataSource{}
}

type requestSequenceDataSource struct {
	*config.Config
}

type requestSequenceDataSourceModel struct {
	Id                      types.String `tfsdk:"id"`
	ResourceId              types.String `tfsdk:"resource_id"`
	Description             types.String `tfsdk:"description"`
	Link                    types.String `tfsdk:"link"`
	Name                    types.String `tfsdk:"name"`
	CompatibleResourceTypes types.List   `tfsdk:"compatible_resource_types"`
}

func (d *requestSequenceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_request_sequence"
}

func (d *requestSequenceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required: true,
			},
			"resource_id": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Computed: true,
			},
			"link": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Computed: true,
			},
			"compatible_resource_types": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(
						stringvalidator.OneOf("APP", "GROUP"),
					),
				},
			},
		},
	}
}

func (d *requestSequenceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data requestSequenceDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	readRequestSeqResp, _, err := d.OktaGovernanceClient.OktaIGSDKClient().RequestSequencesAPI.GetResourceRequestSequenceV2(ctx, data.ResourceId.ValueString(), data.Id.ValueString()).Execute()
	if err != nil {
		return
	}

	// Example data value setting
	data.Link = types.StringValue(readRequestSeqResp.Link)
	data.Description = types.StringValue(readRequestSeqResp.Description)
	data.Name = types.StringValue(readRequestSeqResp.Name)
	data.CompatibleResourceTypes = setCompatibleResourceType(readRequestSeqResp.CompatibleResourceTypes)
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

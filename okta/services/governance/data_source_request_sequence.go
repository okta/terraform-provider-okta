package governance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var _ datasource.DataSource = &requestSequenceDataSource{}

func newRequestSequencesDataSource() datasource.DataSource {
	return &requestSequenceDataSource{}
}

func (d *requestSequenceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_request_sequence"
}

func (d *requestSequenceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

func (d *requestSequenceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The unique identifier for the request sequence. This is typically the sequence ID in Okta.",
			},
			"resource_id": schema.StringAttribute{
				Required:    true,
				Description: "The id of the resource in Okta ID format.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "Description of the request sequence",
			},
			"link": schema.StringAttribute{
				Computed:    true,
				Description: "Link to edit the request sequence.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Name of the request sequence.",
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

func (d *requestSequenceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data requestSequenceDataSourceModel

	// Read Terraform configuration Data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	readRequestSeqResp, _, err := d.OktaGovernanceClient.OktaGovernanceSDKClient().RequestSequencesAPI.GetResourceRequestSequenceV2(ctx, data.ResourceId.ValueString(), data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Principal Entitlements",
			"Could not read Principal Entitlements, unexpected error: "+err.Error(),
		)
		return
	}

	// Example Data value setting
	data.Link = types.StringValue(readRequestSeqResp.Link)
	data.Description = types.StringValue(readRequestSeqResp.Description)
	data.Name = types.StringValue(readRequestSeqResp.Name)
	data.CompatibleResourceTypes = setCompatibleResourceType(readRequestSeqResp.CompatibleResourceTypes)
	// Save Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

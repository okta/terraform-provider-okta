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

var _ datasource.DataSource = &approvalSequenceDataSource{}

func newApprovalSequenceDataSource() datasource.DataSource {
	return &approvalSequenceDataSource{}
}

func (d *approvalSequenceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_approval_sequence"
}

func (d *approvalSequenceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

func (d *approvalSequenceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Gets an Okta Approval Sequence (also known as Request Sequence in the Okta API).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The unique identifier for the approval sequence.",
			},
			"resource_id": schema.StringAttribute{
				Required:    true,
				Description: "The id of the resource in Okta ID format.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "Description of the approval sequence.",
			},
			"link": schema.StringAttribute{
				Computed:    true,
				Description: "Link to edit the approval sequence.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Name of the approval sequence.",
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

type approvalSequenceDataSource struct {
	*config.Config
}

type approvalSequenceDataSourceModel struct {
	Id                      types.String `tfsdk:"id"`
	ResourceId              types.String `tfsdk:"resource_id"`
	Description             types.String `tfsdk:"description"`
	Link                    types.String `tfsdk:"link"`
	Name                    types.String `tfsdk:"name"`
	CompatibleResourceTypes types.List   `tfsdk:"compatible_resource_types"`
}

func (d *approvalSequenceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data approvalSequenceDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResp, _, err := d.OktaGovernanceClient.OktaGovernanceSDKClient().RequestSequencesAPI.GetResourceRequestSequenceV2(ctx, data.ResourceId.ValueString(), data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Approval Sequence",
			"Could not read Approval Sequence, unexpected error: "+err.Error(),
		)
		return
	}

	data.Link = types.StringValue(readResp.Link)
	data.Description = types.StringValue(readResp.Description)
	data.Name = types.StringValue(readResp.Name)
	data.CompatibleResourceTypes = setCompatibleResourceType(readResp.CompatibleResourceTypes)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

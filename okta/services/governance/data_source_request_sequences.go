package governance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var _ datasource.DataSource = &requestSequencesDataSource{}

func newRequestSequencesDataSource() datasource.DataSource {
	return &requestSequencesDataSource{}
}

type requestSequencesDataSource struct {
	*config.Config
}

var requestSequenceAttrTypes = map[string]attr.Type{
	"id":                        types.StringType,
	"name":                      types.StringType,
	"description":               types.StringType,
	"link":                      types.StringType,
	"compatible_resource_types": types.ListType{ElemType: types.StringType},
}

type requestSequencesDataSourceModel struct {
	Id               types.String `tfsdk:"id"`
	ResourceId       types.String `tfsdk:"resource_id"`
	RequestSequences types.List   `tfsdk:"request_sequences"`
}

func (d *requestSequencesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_request_sequences"
}

func (d *requestSequencesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

func (d *requestSequencesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all access request sequences for a resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The id of the resource whose request sequences are listed. Same as `resource_id`.",
			},
			"resource_id": schema.StringAttribute{
				Required:    true,
				Description: "The id of the resource in Okta ID format or ORN format.",
			},
			"request_sequences": schema.ListAttribute{
				Computed:    true,
				Description: "The list of request sequences for the resource. Each element has the attributes `id`, `name`, `description`, `link` and `compatible_resource_types`.",
				ElementType: types.ObjectType{AttrTypes: requestSequenceAttrTypes},
			},
		},
	}
}

func (d *requestSequencesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data requestSequencesDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	listRequestSeqResp, _, err := d.OktaGovernanceClient.OktaGovernanceSDKClient().RequestSequencesAPI.ListResourceRequestSequencesV2(ctx, data.ResourceId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Request Sequences",
			"Could not list request sequences, unexpected error: "+err.Error(),
		)
		return
	}

	sequences := make([]attr.Value, 0, len(listRequestSeqResp.GetData()))
	for _, seq := range listRequestSeqResp.GetData() {
		seqObj, diags := types.ObjectValue(requestSequenceAttrTypes, map[string]attr.Value{
			"id":                        types.StringValue(seq.GetId()),
			"name":                      types.StringValue(seq.GetName()),
			"description":               types.StringValue(seq.GetDescription()),
			"link":                      types.StringValue(seq.GetLink()),
			"compatible_resource_types": setCompatibleResourceType(seq.GetCompatibleResourceTypes()),
		})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		sequences = append(sequences, seqObj)
	}

	seqList, diags := types.ListValue(types.ObjectType{AttrTypes: requestSequenceAttrTypes}, sequences)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Id = data.ResourceId
	data.RequestSequences = seqList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

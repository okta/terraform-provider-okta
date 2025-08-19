package governance

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/okta/terraform-provider-okta/okta/config"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &requestV2DataSource{}

func newRequestV2DataSource() datasource.DataSource {
	return &requestV2DataSource{}
}

type requestV2DataSource struct {
	*config.Config
}

type requestV2DataSourceModel struct {
	Id            types.String            `tfsdk:"id"`
	Created       types.String            `tfsdk:"created"`
	CreatedBy     types.String            `tfsdk:"created_by"`
	LastUpdated   types.String            `tfsdk:"last_updated"`
	LastUpdatedBy types.String            `tfsdk:"last_updated_by"`
	Status        types.String            `tfsdk:"status"`
	Requested     *requested              `tfsdk:"requested"`
	RequestedFor  *entitlementParentModel `tfsdk:"requested_for"`
	RequestedBy   *entitlementParentModel `tfsdk:"requested_by"`
}

func (d *requestV2DataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_request_v2"
}

func (d *requestV2DataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

func (d *requestV2DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required: true,
			},
			"created": schema.StringAttribute{
				Computed: true,
			},
			"created_by": schema.StringAttribute{
				Computed: true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
			"last_updated_by": schema.StringAttribute{
				Computed: true,
			},
			"status": schema.StringAttribute{
				Computed: true,
			},
		},
		Blocks: map[string]schema.Block{
			"requested": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"entry_id": schema.StringAttribute{
						Computed: true,
					},
					"type": schema.StringAttribute{
						Computed: true,
					},
					"access_scope_id": schema.StringAttribute{
						Computed: true,
					},
					"access_scope_type": schema.StringAttribute{
						Computed: true,
					},
					"resource_id": schema.StringAttribute{
						Computed: true,
					},
					"resource_type": schema.StringAttribute{
						Computed: true,
					},
				},
			},
			"requested_for": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"external_id": schema.StringAttribute{
						Computed: true,
					},
					"type": schema.StringAttribute{
						Computed: true,
						Validators: []validator.String{
							stringvalidator.OneOf("OKTA_USER"),
						},
					},
				},
			},
			"requested_by": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"external_id": schema.StringAttribute{
						Computed: true,
					},
					"type": schema.StringAttribute{
						Computed: true,
						Validators: []validator.String{
							stringvalidator.OneOf("OKTA_USER"),
						},
					},
				},
			},
		},
	}
}

func (d *requestV2DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data requestV2DataSourceModel

	// Read Terraform configuration Data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	getRequestV2Resp, _, err := d.OktaGovernanceClient.OktaGovernanceSDKClient().RequestsAPI.GetRequestV2(ctx, data.Id.ValueString()).Execute()
	if err != nil {
		return
	}
	// Example Data value setting
	data.Id = types.StringValue(getRequestV2Resp.GetId())
	data.Created = types.StringValue(getRequestV2Resp.GetCreated().Format(time.RFC3339))
	data.CreatedBy = types.StringValue(getRequestV2Resp.GetCreatedBy())
	data.LastUpdated = types.StringValue(getRequestV2Resp.GetLastUpdated().Format(time.RFC3339))
	data.LastUpdatedBy = types.StringValue(getRequestV2Resp.GetLastUpdatedBy())
	data.Requested = setRequested(getRequestV2Resp.GetRequested())
	data.RequestedBy = setRequestedBy(getRequestV2Resp.GetRequestedBy())
	data.RequestedFor = setRequestedBy(getRequestV2Resp.GetRequestedFor())
	// Save Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

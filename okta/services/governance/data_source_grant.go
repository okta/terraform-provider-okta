package governance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var _ datasource.DataSource = &grantDataSource{}

func newGrantDataSource() datasource.DataSource {
	return &grantDataSource{}
}

type grantDataSource struct {
	*config.Config
}

type grantDataSourceModel struct {
	Id                 types.String `tfsdk:"id"`
	Action             types.String `tfsdk:"action"`
	Actor              types.String `tfsdk:"actor"`
	Created            types.String `tfsdk:"created"`
	CreatedBy          types.String `tfsdk:"created_by"`
	LastUpdated        types.String `tfsdk:"last_updated"`
	LastUpdatedBy      types.String `tfsdk:"last_updated_by"`
	GrantType          types.String `tfsdk:"grant_type"`
	Status             types.String `tfsdk:"status"`
	TargetPrincipalOrn types.String `tfsdk:"target_principal_orn"`
	TargetResourceOrn  types.String `tfsdk:"target_resource_orn"`

	TargetPrincipal *principalModel `tfsdk:"target_principal"`
	Target          *principalModel `tfsdk:"target"`
}

func (d *grantDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_grant"
}

func (d *grantDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

func (d *grantDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required: true,
			},
			"action": schema.StringAttribute{
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("ALLOW", "DENY"),
				},
			},
			"actor": schema.StringAttribute{
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("ACCESS_REQUEST", "ADMIN", "API", "NONE"),
				},
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
			"grant_type": schema.StringAttribute{
				Computed: true,
			},
			"status": schema.StringAttribute{
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("ACTIVE", "EXPIRED", "INACTIVE", "SCHEDULED"),
				},
			},
			"target_principal_orn": schema.StringAttribute{
				Computed: true,
			},
			"target_resource_orn": schema.StringAttribute{
				Computed: true,
			},
		},
		Blocks: map[string]schema.Block{
			"target_principal": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"external_id": schema.StringAttribute{
						Computed: true,
					},
					"type": schema.StringAttribute{
						Computed: true,
					},
				},
			},
			"target": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"external_id": schema.StringAttribute{
						Computed: true,
					},
					"type": schema.StringAttribute{
						Computed: true,
					},
				},
			},
		},
	}
}

func (d *grantDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data grantDataSourceModel

	// Read Terraform configuration Data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	getGrantResp, _, err := d.OktaGovernanceClient.OktaGovernanceSDKClient().GrantsAPI.GetGrant(ctx, data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Grant",
			"Could not retrieve Grant from Okta Governance API: "+err.Error(),
		)
		return
	}
	// Example Data value setting
	data.Id = types.StringValue(getGrantResp.GrantFull.Id)
	data.TargetPrincipalOrn = types.StringValue(getGrantResp.GrantFull.TargetPrincipalOrn)
	data.TargetResourceOrn = types.StringValue(getGrantResp.GrantFull.TargetResourceOrn)
	data.Action = types.StringValue(string(getGrantResp.GrantFull.GetAction()))
	data.Actor = types.StringValue(string(getGrantResp.GrantFull.Actor))
	data.Created = types.StringValue(getGrantResp.GrantFull.GetCreated().String())
	data.CreatedBy = types.StringValue(getGrantResp.GrantFull.GetCreatedBy())
	data.LastUpdated = types.StringValue(getGrantResp.GrantFull.GetLastUpdatedBy())
	data.LastUpdatedBy = types.StringValue(getGrantResp.GrantFull.GetLastUpdatedBy())
	data.GrantType = types.StringValue(string(getGrantResp.GrantFull.GetGrantType()))
	data.Status = types.StringValue(string(getGrantResp.GrantFull.GetStatus()))
	// Map the target principal and target resource to the model
	if targetPrincipal, ok := getGrantResp.GrantFull.GetTargetPrincipalOk(); ok {
		data.TargetPrincipal = &principalModel{
			ExternalId: types.StringValue(targetPrincipal.GetExternalId()),
			Type:       types.StringValue(string(targetPrincipal.GetType())),
		}
	} else {
		data.TargetPrincipal = nil
	}
	if target, ok := getGrantResp.GrantFull.GetTargetOk(); ok {
		data.Target = &principalModel{
			ExternalId: types.StringValue(target.GetExternalId()),
			Type:       types.StringValue(string(target.GetType())),
		}
	} else {
		data.Target = nil
	}

	// Save Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

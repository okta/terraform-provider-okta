package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var _ datasource.DataSource = &entityRiskPolicyDataSource{}

func newEntityRiskPolicyDataSource() datasource.DataSource {
	return &entityRiskPolicyDataSource{}
}

func (d *entityRiskPolicyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

type entityRiskPolicyDataSource struct {
	*config.Config
}

type entityRiskPolicyDataSourceModel struct {
	ID     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	Status types.String `tfsdk:"status"`
}

func (d *entityRiskPolicyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entity_risk_policy"
}

func (d *entityRiskPolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the Entity Risk Policy. This is a system policy that is automatically created when Identity Threat Protection (ITP) is enabled. There is exactly one Entity Risk Policy per organization.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the Entity Risk Policy.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Name of the Entity Risk Policy.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Status of the policy: ACTIVE or INACTIVE.",
			},
		},
	}
}

func (d *entityRiskPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data entityRiskPolicyDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	d.Logger.Info("reading entity risk policy")

	policies, _, err := d.OktaIDaaSClient.OktaSDKClientV6().PolicyAPI.ListPolicies(ctx).Type_("ENTITY_RISK").Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Entity Risk Policy",
			"Failed to list entity risk policies: "+err.Error(),
		)
		return
	}

	if len(policies) == 0 {
		resp.Diagnostics.AddError(
			"No Entity Risk Policy found",
			"Ensure Identity Threat Protection (ITP) is enabled in your organization.",
		)
		return
	}

	// There should be exactly one Entity Risk Policy
	policy := policies[0]

	if policy.EntityRiskPolicy == nil {
		resp.Diagnostics.AddError(
			"Unexpected policy type",
			"Expected EntityRiskPolicy but got a different policy type.",
		)
		return
	}

	entityRiskPolicy := policy.EntityRiskPolicy

	if entityRiskPolicy.Id == nil {
		resp.Diagnostics.AddError(
			"Policy ID is nil",
			"The Entity Risk Policy ID is unexpectedly nil.",
		)
		return
	}

	data.ID = types.StringPointerValue(entityRiskPolicy.Id)
	data.Name = types.StringValue(entityRiskPolicy.Name)
	if entityRiskPolicy.Status != nil {
		data.Status = types.StringValue(string(*entityRiskPolicy.Status))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

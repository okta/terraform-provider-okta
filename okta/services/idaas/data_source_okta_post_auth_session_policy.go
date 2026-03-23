package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var _ datasource.DataSource = &postAuthSessionPolicyDataSource{}

func newPostAuthSessionPolicyDataSource() datasource.DataSource {
	return &postAuthSessionPolicyDataSource{}
}

func (d *postAuthSessionPolicyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

type postAuthSessionPolicyDataSource struct {
	*config.Config
}

type postAuthSessionPolicyDataSourceModel struct {
	ID     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	Status types.String `tfsdk:"status"`
}

func (d *postAuthSessionPolicyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_post_auth_session_policy"
}

func (d *postAuthSessionPolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the Post Auth Session Policy. This is a system policy that is automatically created when Identity Threat Protection (ITP) with Okta AI is enabled. There is exactly one Post Auth Session Policy per organization.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the Post Auth Session Policy.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Name of the Post Auth Session Policy.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Status of the policy: ACTIVE or INACTIVE.",
			},
		},
	}
}

func (d *postAuthSessionPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data postAuthSessionPolicyDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	d.Logger.Info("reading post auth session policy")

	policies, _, err := d.OktaIDaaSClient.OktaSDKClientV6().PolicyAPI.ListPolicies(ctx).Type_("POST_AUTH_SESSION").Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Post Auth Session Policy",
			"Failed to list post auth session policies: "+err.Error(),
		)
		return
	}

	if len(policies) == 0 {
		resp.Diagnostics.AddError(
			"No Post Auth Session Policy found",
			"Ensure Identity Threat Protection (ITP) with Okta AI is enabled in your organization.",
		)
		return
	}

	// There should be exactly one Post Auth Session Policy
	policy := policies[0]

	if policy.PostAuthSessionPolicy == nil {
		resp.Diagnostics.AddError(
			"Unexpected policy type",
			"Expected PostAuthSessionPolicy but got a different policy type.",
		)
		return
	}

	postAuthSessionPolicy := policy.PostAuthSessionPolicy

	if postAuthSessionPolicy.Id == nil {
		resp.Diagnostics.AddError(
			"Policy ID is nil",
			"The Post Auth Session Policy ID is unexpectedly nil.",
		)
		return
	}

	data.ID = types.StringPointerValue(postAuthSessionPolicy.Id)
	data.Name = types.StringValue(postAuthSessionPolicy.Name)
	if postAuthSessionPolicy.Status != nil {
		data.Status = types.StringValue(string(*postAuthSessionPolicy.Status))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

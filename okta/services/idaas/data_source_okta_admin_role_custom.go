package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var _ datasource.DataSource = &adminRoleCustomDataSource{}

func newAdminRoleCustomDataSource() datasource.DataSource {
	return &adminRoleCustomDataSource{}
}

type adminRoleCustomDataSource struct {
	config *config.Config
}

type adminRoleCustomDataSourceModel struct {
	Id          types.String `tfsdk:"id"`
	Label       types.String `tfsdk:"label"`
	Description types.String `tfsdk:"description"`
	Permissions types.Set    `tfsdk:"permissions"`
}

func (d *adminRoleCustomDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_admin_role_custom"
}

func (d *adminRoleCustomDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.config = dataSourceConfiguration(req, resp)
}

func (d *adminRoleCustomDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a custom admin role from Okta.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The unique identifier for the custom role. Accepts the role ID or label.",
			},
			"label": schema.StringAttribute{
				Computed:    true,
				Description: "Unique label for the role.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "Description of the role.",
			},
			"permissions": schema.SetAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "Array of permissions that the role grants.",
			},
		},
	}
}

func (d *adminRoleCustomDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data adminRoleCustomDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := d.config.OktaIDaaSClient.OktaSDKSupplementClient()

	role, _, err := client.GetCustomRole(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading custom admin role",
			"Could not read custom admin role, unexpected error: "+err.Error(),
		)
		return
	}

	data.Id = types.StringValue(role.Id)
	data.Label = types.StringValue(role.Label)
	data.Description = types.StringValue(role.Description)

	perms, _, err := client.ListCustomRolePermissions(ctx, role.Id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading custom admin role permissions",
			"Could not read permissions for custom admin role, unexpected error: "+err.Error(),
		)
		return
	}

	permValues := make([]string, len(perms.Permissions))
	for i, p := range perms.Permissions {
		permValues[i] = p.Label
	}

	permSet, diags := types.SetValueFrom(ctx, types.StringType, permValues)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Permissions = permSet

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

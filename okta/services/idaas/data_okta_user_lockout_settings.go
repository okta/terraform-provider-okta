package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var _ datasource.DataSource = (*userLockoutSettingsDataSource)(nil)

func newUserLockoutSettingsDataSource() datasource.DataSource {
	return &userLockoutSettingsDataSource{}
}

type userLockoutSettingsDataSource struct {
	*config.Config
}

func (d *userLockoutSettingsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

type requestSettingOrganizationDataSourceModel struct {
	Id                                   types.String `tfsdk:"id"`
	VerifyKnowledgeSecondWhen2faRequired types.Bool   `tfsdk:"verify_knowledge_second_when_2fa_required"`
}

func (d *userLockoutSettingsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_request_setting_organization"
}

func (d *userLockoutSettingsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The id property of an entitlement.",
			},
			"verify_knowledge_second_when_2fa_required": schema.BoolAttribute{
				Required:    true,
				Description: "If true, requires users to verify a possession factor before verifying a knowledge factor when the assurance requires two-factor authentication (2FA).",
			},
		},
	}
}

func (d *userLockoutSettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data requestSettingOrganizationDataSourceModel

	// Read Terraform configuration Data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	getAuthSettingsResp, _, err := d.OktaIDaaSClient.OktaSDKClientV5().AttackProtectionAPI.GetAuthenticatorSettings(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Request Setting Organization",
			"Could not read Request Setting Organization, unexpected error: "+err.Error(),
		)
		return
	}

	data.Id = types.StringValue("authenticator_settings")
	data.VerifyKnowledgeSecondWhen2faRequired = types.BoolValue(getAuthSettingsResp[0].GetVerifyKnowledgeSecondWhen2faRequired())

	// Save Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

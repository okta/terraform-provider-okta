package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewDefaultSigninPageDataSource() datasource.DataSource {
	return &defaultSigninPageDataSource{}
}

type defaultSigninPageDataSource struct {
	*Config
}

type defaultSigninPageDatasourceModel struct {
	ID                           types.String `tfsdk:"id"`
	BrandID                      types.String `tfsdk:"brand_id"`
	PageContent                  types.String `tfsdk:"page_content"`
	WidgetVersion                types.String `tfsdk:"widget_version"`
	ContentSecurityPolicySetting types.Object `tfsdk:"content_security_policy_setting"`
	WidgetCustomizations         types.Object `tfsdk:"widget_customizations"`
}

// // TODU
// type contentSecurityPolicySettingModel struct {
// 	Mode      types.String `tfsdk:"mode"`
// 	ReportUri types.String `tfsdk:"report_uri"`
// 	SrcList   types.List   `tfsdk:"src_list"`
// }

type widgetCustomizationsModel struct {
	SignInLabel                             types.String `tfsdk:"sign_in_label"`
	UsernameLabel                           types.String `tfsdk:"username_label"`
	UsernameInfoTip                         types.String `tfsdk:"username_info_tip"`
	PasswordLabel                           types.String `tfsdk:"password_label"`
	PasswordInfoTip                         types.String `tfsdk:"password_info_tip"`
	ShowPasswordVisibilityToggle            types.Bool   `tfsdk:"show_password_visibility_toggle"`
	ShowUserIdentifier                      types.Bool   `tfsdk:"show_user_identifier"`
	ForgotPasswordLabel                     types.String `tfsdk:"forgot_password_label"`
	ForgotPasswordURL                       types.String `tfsdk:"forgot_password_url"`
	UnlockAccountLabel                      types.String `tfsdk:"unlock_account_label"`
	UnlockAccountURL                        types.String `tfsdk:"unlock_account_url"`
	HelpLabel                               types.String `tfsdk:"help_label"`
	HelpURL                                 types.String `tfsdk:"help_url"`
	CustomLink1Label                        types.String `tfsdk:"custom_link_1_label"`
	CustomLink1URL                          types.String `tfsdk:"custom_link_1_url"`
	CustomLink2Label                        types.String `tfsdk:"custom_link_2_label"`
	CustomLink2URL                          types.String `tfsdk:"custom_link_2_url"`
	AuthenticatorPageCustomLinkLabel        types.String `tfsdk:"authenticator_page_custom_link_label"`
	AuthenticatorPageCustomLinkURL          types.String `tfsdk:"authenticator_page_custom_link_url"`
	ClassicRecoveryFlowEmailOrUsernameLabel types.String `tfsdk:"classic_recovery_flow_email_or_username_label"`
	WidgetGeneration                        types.String `tfsdk:"widget_generation"`
}

func (d *defaultSigninPageDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_signin_page"
}

func (d *defaultSigninPageDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve the default signin page of a brand",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "placeholder id",
				Computed:    true,
			},
			"brand_id": schema.StringAttribute{
				Description: "brand id of the default sigin page",
				Required:    true,
			},
			"page_content": schema.StringAttribute{
				Description: "page content of the default sigin page",
				Computed:    true,
			},
			"widget_version": schema.StringAttribute{
				Description: "widget version specified as a Semver",
				Computed:    true,
			},
		},
		// NOTED: due to the provider using protocol v5, schema.SingleNestedAttribute is not support and we have to used schema.SingleNestedBlock instead
		Blocks: map[string]schema.Block{
			"content_security_policy_setting": schema.SingleNestedBlock{
				Description: "",
				Attributes: map[string]schema.Attribute{
					"mode": schema.StringAttribute{
						Description: "enforced or report_only",
						Computed:    true,
					},
					"report_uri": schema.StringAttribute{
						Description: "",
						Computed:    true,
					},
					// NOTED: ListAttribute when used inside Blocks will result in the following error
					// An unexpected error was encountered while verifying an attribute value
					// matched its expected type to prevent unexpected behavior or panics. This is
					// always an error in the provider. Please report the following to the provider
					// developer:
					// Expected framework type from provider logic:
					// types.ListType[basetypes.StringType] / underlying type:
					// tftypes.List[tftypes.String]
					// Received framework type from provider logic: types.ListType[!!! MISSING TYPE
					// !!!] / underlying type: tftypes.List[tftypes.DynamicPseudoType]
					// NOTED: This behavior does not happened when ListAttribute used inside Attributes
					"src_list": schema.ListAttribute{
						Description: "",
						Computed:    true,
						ElementType: types.StringType,
					},
				},
			},
			"widget_customizations": schema.SingleNestedBlock{
				Description: "",
				Attributes: map[string]schema.Attribute{
					"sign_in_label": schema.StringAttribute{
						Description: "",
						Computed:    true,
					},
					"username_label": schema.StringAttribute{
						Description: "",
						Computed:    true,
					},
					"username_info_tip": schema.StringAttribute{
						Description: "",
						Computed:    true,
					},
					"password_label": schema.StringAttribute{
						Description: "",
						Computed:    true,
					},
					"password_info_tip": schema.StringAttribute{
						Description: "",
						Computed:    true,
					},
					"show_password_visibility_toggle": schema.BoolAttribute{
						Description: "",
						Computed:    true,
					},
					"show_user_identifier": schema.BoolAttribute{
						Description: "",
						Computed:    true,
					},
					"forgot_password_label": schema.StringAttribute{
						Description: "",
						Computed:    true,
					},
					"forgot_password_url": schema.StringAttribute{
						Description: "",
						Computed:    true,
					},
					"unlock_account_label": schema.StringAttribute{
						Description: "",
						Computed:    true,
					},
					"unlock_account_url": schema.StringAttribute{
						Description: "",
						Computed:    true,
					},
					"help_label": schema.StringAttribute{
						Description: "",
						Computed:    true,
					},
					"help_url": schema.StringAttribute{
						Description: "",
						Computed:    true,
					},
					"custom_link_1_label": schema.StringAttribute{
						Description: "",
						Computed:    true,
					},
					"custom_link_1_url": schema.StringAttribute{
						Description: "",
						Computed:    true,
					},
					"custom_link_2_label": schema.StringAttribute{
						Description: "",
						Computed:    true,
					},
					"custom_link_2_url": schema.StringAttribute{
						Description: "",
						Computed:    true,
					},
					"authenticator_page_custom_link_label": schema.StringAttribute{
						Description: "",
						Computed:    true,
					},
					"authenticator_page_custom_link_url": schema.StringAttribute{
						Description: "",
						Computed:    true,
					},
					"classic_recovery_flow_email_or_username_label": schema.StringAttribute{
						Description: "",
						Computed:    true,
					},
					"widget_generation": schema.StringAttribute{
						Description: "",
						Computed:    true,
					},
				},
			},
		},
	}
}

func (d *defaultSigninPageDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	config, ok := req.ProviderData.(*Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Datasource Configure Type",
			fmt.Sprintf("Expected *Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.Config = config
}

// TODU
func (d *defaultSigninPageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data defaultSigninPageDatasourceModel
	var diags diag.Diagnostics

	resp.Diagnostics.Append(resp.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	defaultSigninPage, _, err := d.Config.oktaSDKClientV3.CustomizationAPI.GetDefaultSignInPage(ctx, data.BrandID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving default signin page",
			fmt.Sprintf("Error returned: %s", err.Error()),
		)
		return
	}

	data.ID = types.StringValue(data.BrandID.ValueString())
	data.PageContent = types.StringPointerValue(defaultSigninPage.PageContent)
	data.WidgetVersion = types.StringPointerValue(defaultSigninPage.WidgetVersion)

	// contentSecuritySetting := &contentSecurityPolicySettingModel{}
	// if setting, ok := defaultSigninPage.GetContentSecurityPolicySettingOk(); ok {
	// 	contentSecuritySetting.Mode = types.StringPointerValue(setting.Mode)
	// 	contentSecuritySetting.ReportUri = types.StringPointerValue(setting.ReportUri)
	// 	srcList, diags := types.ListValueFrom(ctx, types.StringType, []string{"abc", "xyz"})
	// 	if diags.HasError() {
	// 		resp.Diagnostics.Append(diags...)
	// 		return
	// 	}
	// 	contentSecuritySetting.SrcList = srcList
	// }
	// settingValue, diags := types.ObjectValueFrom(ctx, data.ContentSecurityPolicySetting.AttributeTypes(ctx), contentSecuritySetting)
	// if diags.HasError() {
	// 	resp.Diagnostics.Append(diags...)
	// 	return
	// }
	// data.ContentSecurityPolicySetting = settingValue

	widgetCustomizations := &widgetCustomizationsModel{}
	if customization, ok := defaultSigninPage.GetWidgetCustomizationsOk(); ok {
		widgetCustomizations.SignInLabel = types.StringPointerValue(customization.SignInLabel)
		widgetCustomizations.UsernameLabel = types.StringPointerValue(customization.UsernameLabel)
		widgetCustomizations.UsernameInfoTip = types.StringPointerValue(customization.UsernameInfoTip)
		widgetCustomizations.PasswordLabel = types.StringPointerValue(customization.PasswordLabel)
		widgetCustomizations.PasswordInfoTip = types.StringPointerValue(customization.PasswordInfoTip)
		widgetCustomizations.ShowPasswordVisibilityToggle = types.BoolPointerValue(customization.ShowPasswordVisibilityToggle)
		widgetCustomizations.ShowUserIdentifier = types.BoolPointerValue(customization.ShowUserIdentifier)
		widgetCustomizations.ForgotPasswordLabel = types.StringPointerValue(customization.ForgotPasswordLabel)
		widgetCustomizations.ForgotPasswordURL = types.StringPointerValue(customization.ForgotPasswordUrl)
		widgetCustomizations.UnlockAccountLabel = types.StringPointerValue(customization.UnlockAccountLabel)
		widgetCustomizations.UnlockAccountURL = types.StringPointerValue(customization.UnlockAccountUrl)
		widgetCustomizations.HelpLabel = types.StringPointerValue(customization.HelpLabel)
		widgetCustomizations.HelpURL = types.StringPointerValue(customization.HelpUrl)
		widgetCustomizations.CustomLink1Label = types.StringPointerValue(customization.CustomLink1Label)
		widgetCustomizations.CustomLink1URL = types.StringPointerValue(customization.CustomLink1Url)
		widgetCustomizations.CustomLink2Label = types.StringPointerValue(customization.CustomLink2Label)
		widgetCustomizations.CustomLink2URL = types.StringPointerValue(customization.CustomLink2Url)
		widgetCustomizations.AuthenticatorPageCustomLinkLabel = types.StringPointerValue(customization.AuthenticatorPageCustomLinkLabel)
		widgetCustomizations.AuthenticatorPageCustomLinkURL = types.StringPointerValue(customization.AuthenticatorPageCustomLinkUrl)
		widgetCustomizations.ClassicRecoveryFlowEmailOrUsernameLabel = types.StringPointerValue(customization.ClassicRecoveryFlowEmailOrUsernameLabel)
		widgetCustomizations.WidgetGeneration = types.StringPointerValue((*string)(customization.WidgetGeneration.Ptr()))
	}
	customizationValue, diags := types.ObjectValueFrom(ctx, data.WidgetCustomizations.AttributeTypes(ctx), widgetCustomizations)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	data.WidgetCustomizations = customizationValue
	// Save data into state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/okta/okta-sdk-golang/v4/okta"
)

type signinPageModel struct {
	ID                           types.String `tfsdk:"id"`
	BrandID                      types.String `tfsdk:"brand_id"`
	PageContent                  types.String `tfsdk:"page_content"`
	WidgetVersion                types.String `tfsdk:"widget_version"`
	ContentSecurityPolicySetting types.Object `tfsdk:"content_security_policy_setting"`
	WidgetCustomizations         types.Object `tfsdk:"widget_customizations"`
}

type contentSecurityPolicySettingModel struct {
	Mode      types.String `tfsdk:"mode"`
	ReportUri types.String `tfsdk:"report_uri"`
	SrcList   types.List   `tfsdk:"src_list"`
}

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

func mapSignInPageToState(ctx context.Context, data *okta.SignInPage, state *signinPageModel) diag.Diagnostics {
	var diags diag.Diagnostics
	state.ID = types.StringValue(state.BrandID.ValueString())
	state.PageContent = types.StringPointerValue(data.PageContent)
	state.WidgetVersion = types.StringPointerValue(data.WidgetVersion)
	if setting, ok := data.GetContentSecurityPolicySettingOk(); ok {
		srcList := make([]attr.Value, 0)
		for _, v := range setting.SrcList {
			srcList = append(srcList, types.StringValue(v))
		}
		listValues := types.ListValueMust(types.StringType, srcList)
		elements := map[string]attr.Value{
			"src_list":   listValues,
			"mode":       types.StringPointerValue(setting.Mode),
			"report_uri": types.StringPointerValue(setting.ReportUri),
		}
		elementTypes := map[string]attr.Type{
			"src_list":   types.ListType{ElemType: types.StringType},
			"mode":       types.StringType,
			"report_uri": types.StringType,
		}
		settingValue := types.ObjectValueMust(elementTypes, elements)
		state.ContentSecurityPolicySetting = settingValue
	}

	widgetCustomizations := &widgetCustomizationsModel{}
	if customization, ok := data.GetWidgetCustomizationsOk(); ok {
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
		widgetCustomizations.WidgetGeneration = types.StringPointerValue(customization.WidgetGeneration)
	}
	customizationValue, diags := types.ObjectValueFrom(ctx, state.WidgetCustomizations.AttributeTypes(ctx), widgetCustomizations)
	if diags.HasError() {
		diags.Append(diags...)
		return diags
	}
	state.WidgetCustomizations = customizationValue
	return diags
}

func buildSignInPageRequest(ctx context.Context, model signinPageModel) (okta.SignInPage, diag.Diagnostics) {
	sp := okta.SignInPage{}
	if !model.PageContent.IsNull() {
		sp.SetPageContent(model.PageContent.ValueString())
	}

	if !model.WidgetVersion.IsNull() {
		sp.SetWidgetVersion(model.WidgetVersion.ValueString())
	}

	wc := okta.SignInPageAllOfWidgetCustomizations{}
	wcm := &widgetCustomizationsModel{}
	diags := model.WidgetCustomizations.As(ctx, wcm, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return *okta.NewSignInPage(), diags
	}
	wc.SignInLabel = wcm.SignInLabel.ValueStringPointer()
	wc.UsernameLabel = wcm.UsernameLabel.ValueStringPointer()
	wc.UsernameInfoTip = wcm.UsernameInfoTip.ValueStringPointer()
	wc.PasswordLabel = wcm.PasswordLabel.ValueStringPointer()
	wc.PasswordInfoTip = wcm.PasswordInfoTip.ValueStringPointer()
	wc.ShowPasswordVisibilityToggle = wcm.ShowPasswordVisibilityToggle.ValueBoolPointer()
	wc.ShowUserIdentifier = wcm.ShowUserIdentifier.ValueBoolPointer()
	wc.ForgotPasswordLabel = wcm.ForgotPasswordLabel.ValueStringPointer()
	wc.ForgotPasswordUrl = wcm.ForgotPasswordURL.ValueStringPointer()
	wc.UnlockAccountLabel = wcm.UnlockAccountLabel.ValueStringPointer()
	wc.UnlockAccountUrl = wcm.UnlockAccountURL.ValueStringPointer()
	wc.HelpLabel = wcm.HelpLabel.ValueStringPointer()
	wc.HelpUrl = wcm.HelpURL.ValueStringPointer()
	wc.CustomLink1Label = wcm.CustomLink1Label.ValueStringPointer()
	wc.CustomLink1Url = wcm.CustomLink1URL.ValueStringPointer()
	wc.CustomLink2Label = wcm.CustomLink2Label.ValueStringPointer()
	wc.CustomLink2Url = wcm.CustomLink2URL.ValueStringPointer()
	wc.AuthenticatorPageCustomLinkLabel = wcm.AuthenticatorPageCustomLinkLabel.ValueStringPointer()
	wc.AuthenticatorPageCustomLinkUrl = wcm.AuthenticatorPageCustomLinkURL.ValueStringPointer()
	wc.ClassicRecoveryFlowEmailOrUsernameLabel = wcm.ClassicRecoveryFlowEmailOrUsernameLabel.ValueStringPointer()
	wc.SetWidgetGeneration(wcm.WidgetGeneration.ValueString())
	sp.SetWidgetCustomizations(wc)

	if !model.ContentSecurityPolicySetting.IsNull() {
		csp := okta.ContentSecurityPolicySetting{}
		cspm := &contentSecurityPolicySettingModel{}
		diags := model.ContentSecurityPolicySetting.As(ctx, cspm, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return *okta.NewSignInPage(), diags
		}
		csp.Mode = cspm.Mode.ValueStringPointer()
		csp.ReportUri = cspm.ReportUri.ValueStringPointer()
		elements := make([]types.String, 0, len(cspm.SrcList.Elements()))
		diags = cspm.SrcList.ElementsAs(ctx, &elements, false)
		if diags.HasError() {
			return *okta.NewSignInPage(), diags
		}
		convertElements := make([]string, 0)
		for _, v := range elements {
			convertElements = append(convertElements, v.ValueString())
		}
		csp.SrcList = convertElements

		sp.SetContentSecurityPolicySetting(csp)
	}

	return sp, nil
}

var dataSourceSignInSchema = datasourceSchema.Schema{
	Attributes: map[string]datasourceSchema.Attribute{
		"id": datasourceSchema.StringAttribute{
			Description: "placeholder id",
			Computed:    true,
		},
		"brand_id": datasourceSchema.StringAttribute{
			Description: "brand id of the preview signin page",
			Required:    true,
		},
		"page_content": datasourceSchema.StringAttribute{
			Description: "page content of the preview signin page",
			Computed:    true,
		},
		"widget_version": datasourceSchema.StringAttribute{
			Description: "widget version specified as a Semver",
			Computed:    true,
		},
	},
	// NOTED: due to the provider using protocol v5, schema.SingleNestedAttribute is not support and we have to used schema.SingleNestedBlock instead
	Blocks: map[string]datasourceSchema.Block{
		"content_security_policy_setting": datasourceSchema.SingleNestedBlock{
			Description: "",
			Attributes: map[string]datasourceSchema.Attribute{
				"mode": datasourceSchema.StringAttribute{
					Description: "enforced or report_only",
					Computed:    true,
				},
				"report_uri": datasourceSchema.StringAttribute{
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
				"src_list": datasourceSchema.ListAttribute{
					Description: "",
					Computed:    true,
					ElementType: types.StringType,
				},
			},
		},
		"widget_customizations": datasourceSchema.SingleNestedBlock{
			Description: "",
			Attributes: map[string]datasourceSchema.Attribute{
				"sign_in_label": datasourceSchema.StringAttribute{
					Description: "",
					Computed:    true,
				},
				"username_label": datasourceSchema.StringAttribute{
					Description: "",
					Computed:    true,
				},
				"username_info_tip": datasourceSchema.StringAttribute{
					Description: "",
					Computed:    true,
				},
				"password_label": datasourceSchema.StringAttribute{
					Description: "",
					Computed:    true,
				},
				"password_info_tip": datasourceSchema.StringAttribute{
					Description: "",
					Computed:    true,
				},
				"show_password_visibility_toggle": datasourceSchema.BoolAttribute{
					Description: "",
					Computed:    true,
				},
				"show_user_identifier": datasourceSchema.BoolAttribute{
					Description: "",
					Computed:    true,
				},
				"forgot_password_label": datasourceSchema.StringAttribute{
					Description: "",
					Computed:    true,
				},
				"forgot_password_url": datasourceSchema.StringAttribute{
					Description: "",
					Computed:    true,
				},
				"unlock_account_label": datasourceSchema.StringAttribute{
					Description: "",
					Computed:    true,
				},
				"unlock_account_url": datasourceSchema.StringAttribute{
					Description: "",
					Computed:    true,
				},
				"help_label": datasourceSchema.StringAttribute{
					Description: "",
					Computed:    true,
				},
				"help_url": datasourceSchema.StringAttribute{
					Description: "",
					Computed:    true,
				},
				"custom_link_1_label": datasourceSchema.StringAttribute{
					Description: "",
					Computed:    true,
				},
				"custom_link_1_url": datasourceSchema.StringAttribute{
					Description: "",
					Computed:    true,
				},
				"custom_link_2_label": datasourceSchema.StringAttribute{
					Description: "",
					Computed:    true,
				},
				"custom_link_2_url": datasourceSchema.StringAttribute{
					Description: "",
					Computed:    true,
				},
				"authenticator_page_custom_link_label": datasourceSchema.StringAttribute{
					Description: "",
					Computed:    true,
				},
				"authenticator_page_custom_link_url": datasourceSchema.StringAttribute{
					Description: "",
					Computed:    true,
				},
				"classic_recovery_flow_email_or_username_label": datasourceSchema.StringAttribute{
					Description: "",
					Computed:    true,
				},
				"widget_generation": datasourceSchema.StringAttribute{
					Description: "",
					Computed:    true,
				},
			},
		},
	},
}

var resourceSignInSchema = resourceSchema.Schema{
	Attributes: map[string]resourceSchema.Attribute{
		"id": resourceSchema.StringAttribute{
			Description: "placeholder id",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"brand_id": resourceSchema.StringAttribute{
			Description: "brand id of the preview signin page",
			Required:    true,
		},
		"page_content": resourceSchema.StringAttribute{
			Description: "page content of the preview signin page",
			Required:    true,
		},
		"widget_version": resourceSchema.StringAttribute{
			Description: `widget version specified as a Semver. The following are currently supported
			*, ^1, ^2, ^3, ^4, ^5, ^6, ^7, 1.6, 1.7, 1.8, 1.9, 1.10, 1.11, 1.12, 1.13, 2.1, 2.2, 2.3, 2.4,
			2.5, 2.6, 2.7, 2.8, 2.9, 2.10, 2.11, 2.12, 2.13, 2.14, 2.15, 2.16, 2.17, 2.18, 2.19, 2.20, 2.21,
			3.0, 3.1, 3.2, 3.3, 3.4, 3.5, 3.6, 3.7, 3.8, 3.9, 4.0, 4.1, 4.2, 4.3, 4.4, 4.5, 5.0, 5.1, 5.2, 5.3,
			5.4, 5.5, 5.6, 5.7, 5.8, 5.9, 5.10, 5.11, 5.12, 5.13, 5.14, 5.15, 5.16, 6.0, 6.1, 6.2, 6.3, 6.4, 6.5,
			6.6, 6.7, 6.8, 6.9, 7.0, 7.1, 7.2, 7.3, 7.4, 7.5, 7.6, 7.7, 7.8, 7.9, 7.10, 7.11, 7.12, 7.13.`,
			Required: true,
		},
	},

	// NOTED: due to the provider using protocol v5, schema.SingleNestedAttribute is not support and we have to used schema.SingleNestedBlock instead
	Blocks: map[string]resourceSchema.Block{
		"content_security_policy_setting": resourceSchema.SingleNestedBlock{
			Description: "",
			Attributes: map[string]resourceSchema.Attribute{
				"mode": resourceSchema.StringAttribute{
					Description: "enforced or report_only",
					Optional:    true,
				},
				"report_uri": resourceSchema.StringAttribute{
					Description: "",
					Optional:    true,
				},
				"src_list": resourceSchema.ListAttribute{
					Description: "",
					Optional:    true,
					ElementType: types.StringType,
				},
			},
		},
		"widget_customizations": resourceSchema.SingleNestedBlock{
			Description: "",
			Attributes: map[string]resourceSchema.Attribute{
				"sign_in_label": resourceSchema.StringAttribute{
					Description: "",
					Optional:    true,
				},
				"username_label": resourceSchema.StringAttribute{
					Description: "",
					Optional:    true,
				},
				"username_info_tip": resourceSchema.StringAttribute{
					Description: "",
					Optional:    true,
				},
				"password_label": resourceSchema.StringAttribute{
					Description: "",
					Optional:    true,
				},
				"password_info_tip": resourceSchema.StringAttribute{
					Description: "",
					Optional:    true,
				},
				// NOTED: Computed is behave funny here, if Computed is not set, this attr does not show up in the payload and return from API as true
				// if Computed is set, this attr show up in the payload as false
				"show_password_visibility_toggle": resourceSchema.BoolAttribute{
					Description: "",
					Optional:    true,
					Computed:    true,
				},
				"show_user_identifier": resourceSchema.BoolAttribute{
					Description: "",
					Optional:    true,
					Computed:    true,
				},
				"forgot_password_label": resourceSchema.StringAttribute{
					Description: "",
					Optional:    true,
				},
				"forgot_password_url": resourceSchema.StringAttribute{
					Description: "",
					Optional:    true,
				},
				"unlock_account_label": resourceSchema.StringAttribute{
					Description: "",
					Optional:    true,
				},
				"unlock_account_url": resourceSchema.StringAttribute{
					Description: "",
					Optional:    true,
				},
				"help_label": resourceSchema.StringAttribute{
					Description: "",
					Optional:    true,
				},
				"help_url": resourceSchema.StringAttribute{
					Description: "",
					Optional:    true,
				},
				"custom_link_1_label": resourceSchema.StringAttribute{
					Description: "",
					Optional:    true,
				},
				"custom_link_1_url": resourceSchema.StringAttribute{
					Description: "",
					Optional:    true,
				},
				"custom_link_2_label": resourceSchema.StringAttribute{
					Description: "",
					Optional:    true,
				},
				"custom_link_2_url": resourceSchema.StringAttribute{
					Description: "",
					Optional:    true,
				},
				"authenticator_page_custom_link_label": resourceSchema.StringAttribute{
					Description: "",
					Optional:    true,
				},
				"authenticator_page_custom_link_url": resourceSchema.StringAttribute{
					Description: "",
					Optional:    true,
				},
				"classic_recovery_flow_email_or_username_label": resourceSchema.StringAttribute{
					Description: "",
					Optional:    true,
				},
				"widget_generation": resourceSchema.StringAttribute{
					Description: "",
					Required:    true,
				},
			},
		},
	},
}

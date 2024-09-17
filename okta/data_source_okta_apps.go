package okta

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/okta-sdk-golang/v5/okta"
)

var _ datasource.DataSource = &AppsDataSource{}

func NewAppsDataSource() datasource.DataSource { return &AppsDataSource{} }

type AppsDataSource struct{ config *Config }

type AppsDataSourceModel struct {
	ActiveOnly        types.Bool     `tfsdk:"active_only"`
	Label             types.String   `tfsdk:"label"`
	LabelPrefix       types.String   `tfsdk:"label_prefix"`
	IncludeNonDeleted types.Bool     `tfsdk:"include_non_deleted"`
	UseOptimization   types.Bool     `tfsdk:"use_optimization"`
	Apps              []OktaAppModel `tfsdk:"apps"`
}

type OktaAppModel struct {
	ID          types.String `tfsdk:"id"`
	Created     types.String `tfsdk:"created"`
	LastUpdated types.String `tfsdk:"last_updated"`
	Name        types.String `tfsdk:"name"`
	Label       types.String `tfsdk:"label"`
	Status      types.String `tfsdk:"status"`
	SignOnMode  types.String `tfsdk:"sign_on_mode"`
	Features    types.List   `tfsdk:"features"`
	// LogoUrl     types.String `json:"logo_url" tfsdk:"logo_url"`
	// AdminNote    types.String `json:"admin_note" tfsdk:"admin_note"`
	// EndUserNote  types.String `json:"enduser_note" tfsdk:"enduser_note"`
	// AppLinksJson types.Map `json:"app_links_json" tfsdk:"app_links_json"`
	// Accessibility              types.Map    `json:"accessibility" tfsdk:"accessibility"`
	// AutoSubmitToolbar          types.Bool   `json:"auto_submit_toolbar" tfsdk:"auto_submit_toolbar"`
	// HideIos                    types.Bool   `json:"hide_ios" tfsdk:"hide_ios"`
	// HideWeb                    types.Bool   `json:"hide_web" tfsdk:"hide_web"`
	// UserNameTemplate           types.String `json:"user_name_template" tfsdk:"user_name_template"`
	// UserNameTemplateSuffix     types.String `json:"user_name_template_suffix" tfsdk:"user_name_template_suffix"`
	// UserNameTemplateType       types.String `json:"user_name_template_type" tfsdk:"user_name_template_type"`
	// UserNameTemplatePushStatus types.String `json:"user_name_template_push_status" tfsdk:"user_name_template_push_status"`
	// Profile                    types.String `json:"app_profile_json" tfsdk:"app_profile_json"`
	// Links                      types.Map    `json:"links" tfsdk:"links"`
}

type OktaApp interface {
	GetId() string
	GetCreated() time.Time
	GetLastUpdated() time.Time
	GetName() string
	GetLabel() string
	GetStatus() string
	GetSignOnMode() string
	GetFeatures() []string
	// GetEmbedded() map[string]map[string]interface{}
	// GetSettings() ApplicationSettings
	// GetNotes() string
	// GetProfile() string
	// GetAccessibility() ApplicationAccessibility
	// GetLicensing() ApplicationLicensing
	// GetVisibility() ApplicationVisibility
	// GetLinks() ApplicationLinks
	// GetCredentials() ApplicationCredentials
}

func (d *AppsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_apps"
}

func (d *AppsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"active_only": schema.BoolAttribute{
				Optional:    true,
				Description: "Search only active applications.",
			},
			"label": schema.StringAttribute{
				Optional:    true,
				Validators:  []validator.String{stringvalidator.ConflictsWith(path.MatchRoot("label_prefix"))},
				Description: "Searches for applications whose label or name property matches this value exactly.",
			},
			"label_prefix": schema.StringAttribute{
				Optional:    true,
				Validators:  []validator.String{stringvalidator.ConflictsWith(path.MatchRoot("label"))},
				Description: "Searches for applications whose label or name property begins with this value.",
			},
			"include_non_deleted": schema.BoolAttribute{
				Optional:    true,
				Description: "Specifies whether to include non-active, but not deleted apps in the results.",
			},
			"use_optimization": schema.BoolAttribute{
				Optional:    true,
				Description: "Specifies whether to use query optimization. If you specify `useOptimization=true` in the request query, the response contains a subset of app instance properties.",
			},
			"apps": schema.ListAttribute{
				Description: "The list of applications that match the search criteria.",
				Computed:    true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"id":           types.StringType,
						"created":      types.StringType,
						"last_updated": types.StringType,
						"name":         types.StringType,
						"label":        types.StringType,
						"status":       types.StringType,
						"sign_on_mode": types.StringType,
						"features":     types.ListType{ElemType: types.StringType},
						// "logo_url":       types.StringType,
						// "admin_note":     types.StringType,
						// "enduser_note":   types.StringType,
						// "app_links_json": types.MapType{ElemType: types.StringType},
						// "accessibility":                  types.MapType{ElemType: types.StringType},
						// "auto_submit_toolbar": types.BoolType,
						// "hide_ios":            types.BoolType,
						// "hide_web":            types.BoolType,
						// "user_name_template":             types.StringType,
						// "user_name_template_suffix":      types.StringType,
						// "user_name_template_type":        types.StringType,
						// "user_name_template_push_status": types.StringType,
						// "profile":  types.StringType,
						// "links":    types.MapType{ElemType: types.StringType},
					},
				},
			},
		},
	}
}

func (d *AppsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.config = dataSourceConfiguration(req, resp)
}

func (d *AppsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state AppsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	var filters []string
	if state.ActiveOnly == types.BoolValue(true) {
		filters = append(filters, `status eq "ACTIVE"`)
	}
	if label := state.Label.ValueString(); label != "" {
		filters = append(filters, fmt.Sprintf(`label eq "%s"`, label))
	} else if labelPrefix := state.LabelPrefix.ValueString(); labelPrefix != "" {
		filters = append(filters, fmt.Sprintf(`label sw "%s"`, labelPrefix))
	}
	filterValue := strings.Join(filters, " AND ")

	// Read the list of applications from Okta.
	apiRequest := d.config.oktaSDKClientV5.ApplicationAPI.ListApplications(ctx).Filter(filterValue)
	applicationList, apiResp, err := apiRequest.Execute()
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read Okta Apps", fmt.Sprintf("Error retrieving apps: %s", err.Error()))
		return
	}
	// Handle api pagination
	for apiResp.HasNextPage() {
		var nextApps []okta.ListApplications200ResponseInner
		apiResp, err = apiResp.Next(&nextApps)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Read Okta Apps",
				fmt.Sprintf("Error retrieving next page of apps: %s", err.Error()),
			)
			return
		}
		applicationList = append(applicationList, nextApps...)
	}

	// Convert the list of applications to the Terraform schema.
	for _, app := range applicationList {
		oktaApp := app.GetActualInstance().(OktaApp)

		FeaturesList, err := types.ListValueFrom(ctx, types.StringType, oktaApp.GetFeatures())
		resp.Diagnostics.Append(err...)

		state.Apps = append(state.Apps, OktaAppModel{
			ID:          types.StringValue(oktaApp.GetId()),
			Created:     types.StringValue(oktaApp.GetCreated().Format(time.RFC3339)),
			LastUpdated: types.StringValue(oktaApp.GetLastUpdated().Format(time.RFC3339)),
			Name:        types.StringValue(oktaApp.GetName()),
			Label:       types.StringValue(oktaApp.GetLabel()),
			Status:      types.StringValue(oktaApp.GetStatus()),
			SignOnMode:  types.StringValue(oktaApp.GetSignOnMode()),
			Features:    FeaturesList,
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// type ApplicationSettings interface {
// 	GetIdentityStoreId() string
// 	GetImplicitAssignment() bool
// 	GetInlineHookId() string
// 	GetNotes() ApplicationSettingsNotes
// 	GetNotifications() ApplicationSettingsNotifications
// 	GetApp() ApplicationSettingsApp
// }
// type ApplicationSettingsApp interface {
// 	GetAcsUrl() string
// 	GetAudRestriction() string
// 	GetBaseUrl()
// }
// type ApplicationSettingsNotes interface {
// 	GetAdmin() string
// 	GetEnduser() string
// }
// type ApplicationCredentials interface {
// 	GetSigning()
// 	GetUserNameTemplate() ApplicationCredentialsUsernameTemplate
// }
// type ApplicationCredentialsUsernameTemplate interface {
// 	GetPushStatus() string
// 	GetTemplate() string
// 	GetType() string
// 	GetUserSuffix() string
// }
// type ApplicationCredentialsSigning interface {
// 	GetKid() string
// 	GetLastRotated() time.Time
// 	GetNextRotation() time.Time
// 	GetRotationMode() string
// 	GetUse() string
// }
// type ApplicationSettingsNotifications interface {
// 	GetVpn()
// }
// type ApplicationSettingsNotificationsVpn interface {
// 	GetHelpUrl() string
// 	GetMessage() string
// 	GetNetwork() string
// }
// type ApplicationAccessibility interface {
// 	GetErrorRedirectUrl() string
// 	GetLoginRedirectUrl() string
// 	GetSelfService() bool
// }
// type ApplicationLicensing interface {
// 	GetSeatCount() int32
// }
// type ApplicationVisibility interface {
// 	GetAppLinks() map[string]bool
// 	GetAutoLaunch() bool
// 	GetAutoSubmitToolbar() bool
// 	GetHide() ApplicationVisibilityHide
// }
// type ApplicationVisibilityHide interface {
// 	GetIOS() bool
// 	GetWeb() bool
// }
// type ApplicationLinks interface {
// 	GetAccessPolicy() HrefObject
// 	GetActivate() HrefObject
// 	GetDeactivate() HrefObject
// 	GetGroups() HrefObject
// 	GetLogo() HrefObject
// 	GetMetadata() HrefObject
// 	GetSelf() HrefObject
// 	GetUsers() HrefObject
// }
// type HrefObject interface {
// 	GetHints() HrefHints
// 	GetHref() string
// 	GetName() string
// 	GetTemplated() bool
// 	GetType() string
// }
// type HrefHints interface {
// 	GetAllow() []string
// }
// type ApplicationSettingsNotificationsVpnNetwork interface {
// 	GetConnection() string
// 	GetExclude() []string
// 	GetInclude() []string
// }

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
	AdminNote   types.String `tfsdk:"admin_note"`
	EndUserNote types.String `tfsdk:"enduser_note"`
	Visibility  types.Object `tfsdk:"visibility"`
}

type OktaAppVisibility struct {
	AutoLaunch        types.Bool   `tfsdk:"auto_launch"`
	AutoSubmitToolbar types.Bool   `tfsdk:"auto_submit_toolbar"`
	Hide              types.Object `tfsdk:"hide"`
}

type OktaAppVisibilityHide struct {
	IOS types.Bool `tfsdk:"ios"`
	Web types.Bool `tfsdk:"web"`
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
	GetVisibility() okta.ApplicationVisibility
	GetLinks() okta.ApplicationLinks
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
						"admin_note":   types.StringType,
						"enduser_note": types.StringType,
						"visibility": types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"auto_launch":         types.BoolType,
								"auto_submit_toolbar": types.BoolType,
								"hide": types.ObjectType{
									AttrTypes: map[string]attr.Type{
										"ios": types.BoolType,
										"web": types.BoolType,
									},
								},
							},
						},
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
	apiRequest := d.config.oktaSDKClientV5.ApplicationAPI.ListApplications(ctx).Filter(filterValue).Limit(200)
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
		oktaApp, ok := app.GetActualInstance().(OktaApp)
		if !ok {
			continue
		}

		FeaturesList, err := types.ListValueFrom(ctx, types.StringType, oktaApp.GetFeatures())
		resp.Diagnostics.Append(err...)

		adminNotes, userNotes := getNotes(oktaApp)

		hideValue := map[string]attr.Value{
			"ios": types.BoolPointerValue(oktaApp.GetVisibility().Hide.IOS),
			"web": types.BoolPointerValue(oktaApp.GetVisibility().Hide.Web),
		}
		hideTypes := map[string]attr.Type{
			"ios": types.BoolType,
			"web": types.BoolType,
		}

		hide, _ := types.ObjectValue(hideTypes, hideValue)

		visibilityValue := map[string]attr.Value{
			"hide":                hide,
			"auto_launch":         types.BoolPointerValue(oktaApp.GetVisibility().AutoLaunch),
			"auto_submit_toolbar": types.BoolPointerValue(oktaApp.GetVisibility().AutoSubmitToolbar),
		}
		visibilityTypes := map[string]attr.Type{
			"hide": types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"ios": types.BoolType,
					"web": types.BoolType,
				},
			},
			"auto_launch":         types.BoolType,
			"auto_submit_toolbar": types.BoolType,
		}

		visibility, _ := types.ObjectValue(visibilityTypes, visibilityValue)

		state.Apps = append(state.Apps, OktaAppModel{
			ID:          types.StringValue(oktaApp.GetId()),
			Created:     types.StringValue(oktaApp.GetCreated().Format(time.RFC3339)),
			LastUpdated: types.StringValue(oktaApp.GetLastUpdated().Format(time.RFC3339)),
			Name:        types.StringValue(oktaApp.GetName()),
			Label:       types.StringValue(oktaApp.GetLabel()),
			Status:      types.StringValue(oktaApp.GetStatus()),
			SignOnMode:  types.StringValue(oktaApp.GetSignOnMode()),
			Features:    FeaturesList,
			AdminNote:   types.StringValue(adminNotes),
			EndUserNote: types.StringValue(userNotes),
			Visibility:  visibility,
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func getNotesFromSettings(notes *okta.ApplicationSettingsNotes) (string, string) {
	if notes == nil {
		return "", ""
	}

	admin := ""
	if notes.Admin != nil {
		admin = *notes.Admin
	}

	user := ""
	if notes.Enduser != nil {
		user = *notes.Enduser
	}

	return admin, user
}

func getNotes(app interface{}) (string, string) {
	if app == nil {
		return "", ""
	}

	switch v := app.(type) {
	case *okta.AutoLoginApplication:
		return getNotesFromSettings(v.Settings.Notes)
	case *okta.BasicAuthApplication:
		return getNotesFromSettings(v.Settings.Notes)
	case *okta.BookmarkApplication:
		return getNotesFromSettings(v.Settings.Notes)
	case *okta.BrowserPluginApplication:
		return getNotesFromSettings(v.Settings.Notes)
	case *okta.OpenIdConnectApplication:
		return getNotesFromSettings(v.Settings.Notes)
	case *okta.Saml11Application:
		return getNotesFromSettings(v.Settings.Notes)
	case *okta.SamlApplication:
		return getNotesFromSettings(v.Settings.Notes)
	case *okta.SecurePasswordStoreApplication:
		return getNotesFromSettings(v.Settings.Notes)
	case *okta.WsFederationApplication:
		return getNotesFromSettings(v.Settings.Notes)
	}

	return "", ""
}

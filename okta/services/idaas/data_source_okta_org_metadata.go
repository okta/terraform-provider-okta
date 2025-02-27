package idaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &logStreamDataSource{}
	_ datasource.DataSourceWithConfigure = &logStreamDataSource{}
)

func newOrgMetadataDataSource() datasource.DataSource {
	return &orgMetadataDataSource{}
}

type orgMetadataDataSource struct {
	*config.Config
}

type OrgMetadataDataSourceModel struct {
	ID       types.String `tfsdk:"id"`
	Pipeline types.String `tfsdk:"pipeline"`
	Settings types.Object `tfsdk:"settings"`
	Domains  types.Object `tfsdk:"domains"`
}

type OrgMetadataSettingsModel struct {
	AnalyticsCollectionEnabled types.Bool `tfsdk:"analytics_collection_enabled"`
	BugReportingEnabled        types.Bool `tfsdk:"bug_reporting_enabled"`
	OmEnabled                  types.Bool `tfsdk:"om_enabled"`
}

type OrgMetadataDomainsModel struct {
	Organization types.String `tfsdk:"organization"`
	Alternate    types.String `tfsdk:"alternate"`
}

func (d *orgMetadataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_org_metadata"
}

func (d *orgMetadataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves the well-known org metadata, which includes the id, configured custom domains, authentication pipeline, and various other org settings.",
		MarkdownDescription: "Retrieves the well-known org metadata, which includes the id, configured custom domains, authentication pipeline, and various other org settings.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "The unique identifier of the Org.",
				MarkdownDescription: "The unique identifier of the Org.",
				Computed:            true,
			},
			"pipeline": schema.StringAttribute{
				Description:         "The authentication pipeline of the org. idx means the org is using the Identity Engine, while v1 means the org is using the Classic authentication pipeline.",
				MarkdownDescription: "The authentication pipeline of the org. idx means the org is using the Identity Engine, while v1 means the org is using the Classic authentication pipeline.",
				Computed:            true,
			},
		},
		Blocks: map[string]schema.Block{
			"settings": schema.SingleNestedBlock{
				Description:         "The wellknown org settings (safe for public consumption).",
				MarkdownDescription: "The wellknown org settings (safe for public consumption).",
				Attributes: map[string]schema.Attribute{
					"analytics_collection_enabled": schema.BoolAttribute{
						Description:         "",
						MarkdownDescription: "",
						Computed:            true,
					},
					"bug_reporting_enabled": schema.BoolAttribute{
						Description:         "",
						MarkdownDescription: "",
						Computed:            true,
					},
					"om_enabled": schema.BoolAttribute{
						Description:         "Whether the legacy Okta Mobile application is enabled for the org",
						MarkdownDescription: "Whether the legacy Okta Mobile application is enabled for the org",
						Computed:            true,
					},
				},
			},
			"domains": schema.SingleNestedBlock{
				Description:         "The URIs for the org's configured domains.",
				MarkdownDescription: "The URIs for the org's configured domains.",
				Attributes: map[string]schema.Attribute{
					"organization": schema.StringAttribute{
						Description:         "Standard Org URI",
						MarkdownDescription: "Standard Org URI",
						Computed:            true,
					},
					"alternate": schema.StringAttribute{
						Description:         "Custom Domain Org URI",
						MarkdownDescription: "Custom Domain Org URI",
						Computed:            true,
					},
				},
			},
		},
	}
}

func (d *orgMetadataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

func (d *orgMetadataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OrgMetadataDataSourceModel

	resp.Diagnostics.Append(resp.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	org, _, err := d.OktaIDaaSClient.OktaSDKClientV3().OrgSettingAPI.GetWellknownOrgMetadata(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving org metadata",
			fmt.Sprintf("Error returned: %s", err.Error()),
		)
		return
	}

	data.ID = types.StringValue(org.GetId())
	data.Pipeline = types.StringValue(string(org.GetPipeline()))
	settings := &OrgMetadataSettingsModel{}
	orgSettings, ok := org.GetSettingsOk()
	if ok {
		ace, ok := orgSettings.GetAnalyticsCollectionEnabledOk()
		if ok {
			settings.AnalyticsCollectionEnabled = types.BoolValue(*ace)
		}
		bre, ok := orgSettings.GetBugReportingEnabledOk()
		if ok {
			settings.BugReportingEnabled = types.BoolValue(*bre)
		}
		ome, ok := orgSettings.GetOmEnabledOk()
		if ok {
			settings.OmEnabled = types.BoolValue(*ome)
		}
	}
	settingsValue, diags := types.ObjectValueFrom(ctx, data.Settings.AttributeTypes(ctx), settings)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	data.Settings = settingsValue

	domains := &OrgMetadataDomainsModel{}
	orgLinks, ok := org.GetLinksOk()
	if ok {
		al, ok := orgLinks.GetAlternateOk()
		if ok {
			href, ok := al.GetHrefOk()
			if ok {
				domains.Alternate = types.StringValue(*href)
			}
		}
		or, ok := orgLinks.GetOrganizationOk()
		if ok {
			href, ok := or.GetHrefOk()
			if ok {
				domains.Organization = types.StringValue(*href)
			}
		}
	}
	domainsValue, diags := types.ObjectValueFrom(ctx, data.Domains.AttributeTypes(ctx), domains)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	data.Domains = domainsValue

	// Save data into state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

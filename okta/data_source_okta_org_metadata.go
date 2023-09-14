package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewOrgMetadataDataSource() datasource.DataSource {
	return &OrgMetadataDataSource{}
}

type OrgMetadataDataSource struct {
	config *Config
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

func (d *OrgMetadataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_org_metadata"
}

func (d *OrgMetadataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves the well-known org metadata, which includes the id, configured custom domains, authentication pipeline, and various other org settings.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the Org.",
				Computed:            true,
			},
			"pipeline": schema.StringAttribute{
				MarkdownDescription: "The authentication pipeline of the org. idx means the org is using the Identity Engine, while v1 means the org is using the Classic authentication pipeline.",
				Computed:            true,
			},
			"settings": schema.ObjectAttribute{
				MarkdownDescription: "The wellknown org settings (safe for public consumption).",
				Computed:            true,
				AttributeTypes: map[string]attr.Type{
					"analytics_collection_enabled": types.BoolType,
					"bug_reporting_enabled":        types.BoolType,
					"om_enabled":                   types.BoolType,
				},
			},
			"domains": schema.ObjectAttribute{
				MarkdownDescription: "The URIs for the org's configured domains.",
				Computed:            true,
				AttributeTypes: map[string]attr.Type{
					"organization": types.StringType,
					"alternate":    types.StringType,
				},
			},
		},
	}
}

func (d *OrgMetadataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	config, ok := req.ProviderData.(*Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.config = config
}

func (d *OrgMetadataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OrgMetadataDataSourceModel

	resp.Diagnostics.Append(resp.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	org, _, err := d.config.oktaSDKsupplementClient.GetWellKnownOktaOrganization(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving org metadata",
			fmt.Sprintf("Error returned: %s", err.Error()),
		)
		return
	}

	data.ID = types.StringValue(org.Id)
	data.Pipeline = types.StringValue(org.Pipeline)
	data.Settings = types.ObjectValue()

	// Save data into state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

package idaas

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/okta/okta-sdk-golang/v6/okta"

	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ datasource.DataSource              = &identitySourceGroupsDataSource{}
	_ datasource.DataSourceWithConfigure = &identitySourceGroupsDataSource{}
)

// IdentitySourceGroupsDataSource defines the data source implementation.
type identitySourceGroupsDataSource struct {
	Config *config.Config
}

// identitySourceGroupProfileModel describes the profile block of the group.
type identitySourceGroupProfileModel struct {
	DisplayName types.String `tfsdk:"display_name"`
	Description types.String `tfsdk:"description"`
}

// IdentitySourceGroupsDataSourceModel describes the data source data model.
type identitySourceGroupsDataSourceModel struct {
	ID               types.String                     `tfsdk:"id"`
	IdentitySourceId types.String                     `tfsdk:"identity_source_id"`
	ExternalId       types.String                     `tfsdk:"external_id"`
	Profile          *identitySourceGroupProfileModel `tfsdk:"profile"`
}

func newIdentitySourceGroupsDataSource() datasource.DataSource {
	return &identitySourceGroupsDataSource{}
}

func (d *identitySourceGroupsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_identity_source_groups"
}

func (d *identitySourceGroupsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

func (d *identitySourceGroupsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for the Okta `identity_source_groups` resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Okta-assigned ID of the group.",
				Optional:            true,
				Computed:            true,
			},
			"identity_source_id": schema.StringAttribute{
				MarkdownDescription: "ID of the identity source",
				Required:            true,
			},
			"external_id": schema.StringAttribute{
				MarkdownDescription: "The external ID of the group in the identity source.",
				Required:            true,
			},
		},
		Blocks: map[string]schema.Block{
			"profile": schema.SingleNestedBlock{
				MarkdownDescription: "Profile of the group.",
				Attributes: map[string]schema.Attribute{
					"display_name": schema.StringAttribute{
						MarkdownDescription: "Display name of the group.",
						Computed:            true,
					},
					"description": schema.StringAttribute{
						MarkdownDescription: "Description of the group.",
						Computed:            true,
					},
				},
			},
		},
	}
}

func (d *identitySourceGroupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state identitySourceGroupsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The API looks up a group by its external ID (groupOrExternalId parameter).
	identitySourceId := state.IdentitySourceId.ValueString()
	groupExternalId := state.ExternalId.ValueString()

	client := d.Config.OktaIDaaSClient.OktaSDKClientV6()
	result, httpResp, err := client.IdentitySourceAPI.GetIdentitySourceGroup(ctx, identitySourceId, groupExternalId).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			resp.Diagnostics.AddError(
				"Identity source group not found",
				fmt.Sprintf("No group with external ID %q found in identity source %q.", groupExternalId, identitySourceId),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Failed to read identity_source_groups",
			fmt.Sprintf("Error reading identity source group: %s", err.Error()),
		)
		return
	}

	mapIdentitySourceGroupsResponseToState(ctx, result, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// mapIdentitySourceGroupsResponseToState maps the API response to the data source state model.
func mapIdentitySourceGroupsResponseToState(_ context.Context, response *okta.GroupsResponseSchema, state *identitySourceGroupsDataSourceModel, _ *diag.Diagnostics) {
	state.ID = types.StringValue(response.GetId())
	state.ExternalId = types.StringValue(response.GetExternalId())

	profileModel := &identitySourceGroupProfileModel{
		DisplayName: types.StringNull(),
		Description: types.StringNull(),
	}
	// The SDK's GroupsResponseSchemaProfile wraps an extra "profile" layer that the API
	// does not produce — the API returns displayName/description directly at the profile
	// level, which the SDK deserializer places into AdditionalProperties.
	if response.HasProfile() {
		ap := response.GetProfile().AdditionalProperties
		if dn, ok := ap["displayName"].(string); ok {
			profileModel.DisplayName = types.StringValue(dn)
		}
		if desc, ok := ap["description"].(string); ok {
			profileModel.Description = types.StringValue(desc)
		}
	}
	state.Profile = profileModel
}

package idaas

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/okta/terraform-provider-okta/okta/config"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

var (
	_ datasource.DataSource              = &identitySourceGroupMembershipsDataSource{}
	_ datasource.DataSourceWithConfigure = &identitySourceGroupMembershipsDataSource{}
)

// identitySourceGroupMembershipsDataSource defines the data source implementation.
type identitySourceGroupMembershipsDataSource struct {
	Config *config.Config
}

// identitySourceGroupMembershipsDataSourceModel describes the data source data model.
type identitySourceGroupMembershipsDataSourceModel struct {
	ID                types.String `tfsdk:"id"`
	IdentitySourceId  types.String `tfsdk:"identity_source_id"`
	GroupExternalId   types.String `tfsdk:"group_external_id"`
	MemberExternalIds types.List   `tfsdk:"member_external_ids"`
}

func newIdentitySourceGroupMembershipsDataSource() datasource.DataSource {
	return &identitySourceGroupMembershipsDataSource{}
}

func (d *identitySourceGroupMembershipsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_identity_source_group_memberships"
}

func (d *identitySourceGroupMembershipsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

func (d *identitySourceGroupMembershipsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves the group memberships for the given identity source group in the given identity source instance.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The composite ID of the group memberships resource.",
				Computed:            true,
			},
			"identity_source_id": schema.StringAttribute{
				MarkdownDescription: "ID of the identity source.",
				Required:            true,
			},
			"group_external_id": schema.StringAttribute{
				MarkdownDescription: "The external ID (or Okta-assigned ID) of the group whose memberships to retrieve.",
				Required:            true,
			},
			"member_external_ids": schema.ListAttribute{
				MarkdownDescription: "List of external IDs of members belonging to the group in the identity source.",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (d *identitySourceGroupMembershipsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state identitySourceGroupMembershipsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	identitySourceId := state.IdentitySourceId.ValueString()
	groupExternalId := state.GroupExternalId.ValueString()

	client := d.Config.OktaIDaaSClient.OktaSDKClientV6()

	// Paginate through all memberships using after/limit.
	var allMemberIds []string
	apiReq := client.IdentitySourceAPI.
		GetIdentitySourceGroupMemberships(ctx, identitySourceId, groupExternalId).
		Limit(200)

	for {
		result, httpResp, err := apiReq.Execute()
		if err != nil {
			if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
				resp.Diagnostics.AddError(
					"Identity source group not found",
					fmt.Sprintf("No group with external ID %q found in identity source %q.", groupExternalId, identitySourceId),
				)
				return
			}
			resp.Diagnostics.AddError(
				"Failed to read identity_source_group_memberships",
				fmt.Sprintf("Error reading identity source group memberships: %s", err.Error()),
			)
			return
		}

		allMemberIds = append(allMemberIds, result.GetMemberExternalIds()...)

		// Check for next page via Link header.
		after :=
			utils.ExtractAfterCursor(httpResp)
		if after == "" {
			break
		}
		apiReq = client.IdentitySourceAPI.
			GetIdentitySourceGroupMemberships(ctx, identitySourceId, groupExternalId).
			Limit(200).
			After(after)
	}

	mapIdentitySourceGroupMembershipsResponseToState(ctx, allMemberIds, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	state.ID = types.StringValue(fmt.Sprintf("%s/%s", identitySourceId, groupExternalId))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// mapIdentitySourceGroupMembershipsResponseToState maps the member external IDs to the data source state model.
func mapIdentitySourceGroupMembershipsResponseToState(ctx context.Context, memberExternalIds []string, state *identitySourceGroupMembershipsDataSourceModel, diags *diag.Diagnostics) {
	memberList, d := types.ListValueFrom(ctx, types.StringType, memberExternalIds)
	diags.Append(d...)
	state.MemberExternalIds = memberList
}

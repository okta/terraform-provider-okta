package idaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/okta-sdk-golang/v5/okta"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var _ datasource.DataSource = &appGroupAssignmentsDataSource{}

func newAppGroupAssignmentsDataSource() datasource.DataSource {
	return &appGroupAssignmentsDataSource{}
}

type appGroupAssignmentsDataSource struct {
	*config.Config
}

type appGroupAssignmentsDataSourceModel struct {
	ID     types.String   `tfsdk:"id"`
	Groups []types.String `tfsdk:"groups"`
}

func (d *appGroupAssignmentsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_group_assignments"
}

func (d *appGroupAssignmentsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get a set of groups assigned to an Okta application.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the Okta App being queried for groups",
			},
			"groups": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "List of groups IDs assigned to the app",
			},
		},
	}
}

func (d *appGroupAssignmentsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

func (d *appGroupAssignmentsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state appGroupAssignmentsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	client := d.OktaIDaaSClient.OktaSDKClientV5()
	apiRequest := client.ApplicationGroupsAPI.ListApplicationGroupAssignments(ctx, state.ID.ValueString()).Limit(200)
	groupAssignments, httpResp, err := apiRequest.Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Okta App Group Assignments",
			fmt.Sprintf("Error retrieving group assignments: %s", err.Error()),
		)
		return
	}

	// Handle pagination
	for {
		var moreAssignments []okta.ApplicationGroupAssignment
		if httpResp.HasNextPage() {
			httpResp, err = httpResp.Next(&moreAssignments)
			if err != nil {
				resp.Diagnostics.AddError(
					"Unable to Read next page of Okta App Group Assignments",
					fmt.Sprintf("Error retrieving group assignments: %s", err.Error()),
				)
				return
			}
			groupAssignments = append(groupAssignments, moreAssignments...)
		} else {
			break
		}
	}

	// Convert the API response into a list of group IDs
	groups := make([]types.String, 0, len(groupAssignments))
	for _, assignment := range groupAssignments {
		groups = append(groups, types.StringValue(*assignment.Id))
	}

	state.Groups = groups
	state.ID = types.StringValue(state.ID.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

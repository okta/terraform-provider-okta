package idaas

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	okta "github.com/okta/okta-sdk-golang/v6/okta"

	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ datasource.DataSource              = &identitySourceSessionsDataSource{}
	_ datasource.DataSourceWithConfigure = &identitySourceSessionsDataSource{}
)

// IdentitySourceSessionsDataSource defines the data source implementation.
type identitySourceSessionsDataSource struct {
	Config *config.Config
}

// IdentitySourceSessionsDataSourceModel describes the data source data model.
type identitySourceSessionsDataSourceModel struct {
	ID               types.String `tfsdk:"id"`
	IdentitySourceId types.String `tfsdk:"identity_source_id"`
	SessionId        types.String `tfsdk:"session_id"`
	Created          types.String `tfsdk:"created"`
	ImportType       types.String `tfsdk:"import_type"`
	LastUpdated      types.String `tfsdk:"last_updated"`
	Status           types.String `tfsdk:"status"`
}

func newIdentitySourceSessionsDataSource() datasource.DataSource {
	return &identitySourceSessionsDataSource{}
}

func (d *identitySourceSessionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_identity_source_sessions"
}

func (d *identitySourceSessionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

func (d *identitySourceSessionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for the Okta `identity_source_sessions` resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the identity source session.",
				Computed:            true,
			},
			"identity_source_id": schema.StringAttribute{
				MarkdownDescription: "ID of the identity source",
				Required:            true,
			},
			"session_id": schema.StringAttribute{
				MarkdownDescription: "ID of the identity source session.",
				Required:            true,
			},
			"created": schema.StringAttribute{
				MarkdownDescription: "The timestamp when the identity source session was created",
				Computed:            true,
			},
			"import_type": schema.StringAttribute{
				MarkdownDescription: "The type of import.  All imports are `INCREMENTAL` imports.",
				Computed:            true,
			},
			"last_updated": schema.StringAttribute{
				MarkdownDescription: "The timestamp when the identity source session was created",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The current status of the identity source session",
				Computed:            true,
			},
		},
	}
}

func (d *identitySourceSessionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state identitySourceSessionsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	identitySourceId := state.IdentitySourceId.ValueString()
	sessionId := state.SessionId.ValueString()
	client := d.Config.OktaIDaaSClient.OktaSDKClientV6()

	result, httpResp, err := client.IdentitySourceAPI.GetIdentitySourceSession(ctx, identitySourceId, sessionId).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			resp.Diagnostics.AddError(
				"Identity source session not found",
				fmt.Sprintf("No session with ID %q found in identity source %q.", sessionId, identitySourceId),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Failed to read identity_source_sessions",
			fmt.Sprintf("Error reading identity source session: %s", err.Error()),
		)
		return
	}
	mapIdentitySourceSessionsResponseToState(ctx, result, &state, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// mapIdentitySourceSessionsResponseToState maps the API response to the data source state model.
func mapIdentitySourceSessionsResponseToState(_ context.Context, response *okta.IdentitySourceSession, state *identitySourceSessionsDataSourceModel, _ *diag.Diagnostics) {
	state.ID = types.StringValue(response.GetId())
	state.Created = types.StringValue(response.GetCreated().Format(time.RFC3339))
	state.ImportType = types.StringValue(response.GetImportType())
	state.LastUpdated = types.StringValue(response.GetLastUpdated().Format(time.RFC3339))
	state.Status = types.StringValue(response.GetStatus())
}

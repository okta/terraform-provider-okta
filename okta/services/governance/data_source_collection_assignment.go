package governance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var _ datasource.DataSource = (*collectionAssignmentDataSource)(nil)

func newCollectionAssignmentDataSource() datasource.DataSource {
	return &collectionAssignmentDataSource{}
}

type collectionAssignmentDataSource struct {
	*config.Config
}

type collectionAssignmentDataSourceModel struct {
	Id             types.String `tfsdk:"id"`
	CollectionId   types.String `tfsdk:"collection_id"`
	PrincipalId    types.String `tfsdk:"principal_id"`
	PrincipalType  types.String `tfsdk:"principal_type"`
	Actor          types.String `tfsdk:"actor"`
	ExpirationTime types.String `tfsdk:"expiration_time"`
	TimeZone       types.String `tfsdk:"time_zone"`
	AssignmentType types.String `tfsdk:"assignment_type"`
}

func (d *collectionAssignmentDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_collection_assignment"
}

func (d *collectionAssignmentDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

func (d *collectionAssignmentDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves information about a collection assignment to a principal (user or group).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The assignment ID.",
			},
			"collection_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the collection.",
			},
			"principal_id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the principal (user ID or group ID).",
			},
			"principal_type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of principal (OKTA_USER or OKTA_GROUP).",
			},
			"actor": schema.StringAttribute{
				Computed:    true,
				Description: "The actor who made the assignment.",
			},
			"expiration_time": schema.StringAttribute{
				Computed:    true,
				Description: "The date/time when the assignment expires (ISO 8601 format).",
			},
			"time_zone": schema.StringAttribute{
				Computed:    true,
				Description: "The time zone in IANA format for the expiration date.",
			},
			"assignment_type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of assignment (e.g., INDIVIDUAL).",
			},
		},
	}
}

func (d *collectionAssignmentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data collectionAssignmentDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	collectionId := data.CollectionId.ValueString()
	assignmentId := data.Id.ValueString()

	if collectionId == "" {
		resp.Diagnostics.AddError("Missing collection ID", "The 'collection_id' attribute must be set in the configuration.")
		return
	}
	if assignmentId == "" {
		resp.Diagnostics.AddError("Missing assignment ID", "The 'id' attribute must be set in the configuration.")
		return
	}

	// Read API call - list all assignments and find the matching one
	assignments, _, err := d.OktaGovernanceClient.OktaGovernanceSDKClient().CollectionsAPI.
		ListCollectionAssignments(ctx, collectionId).
		Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read collection assignments",
			"Could not retrieve collection assignments, unexpected error: "+err.Error(),
		)
		return
	}

	// Find the specific assignment
	var found bool
	if assignments.HasData() {
		for _, assignment := range assignments.GetData() {
			if assignment.HasId() && assignment.GetId() == assignmentId {
				found = true

				// Map response to state
				data.Id = types.StringValue(assignment.GetId())

				if assignment.HasCollectionId() {
					data.CollectionId = types.StringValue(assignment.GetCollectionId())
				}

				if assignment.HasPrincipal() {
					principal := assignment.GetPrincipal()
					data.PrincipalId = types.StringValue(principal.GetExternalId())
					data.PrincipalType = types.StringValue(string(principal.GetType()))
				}

				if assignment.HasActor() {
					data.Actor = types.StringValue(string(assignment.GetActor()))
				}

				if assignment.HasExpirationTime() {
					data.ExpirationTime = types.StringValue(assignment.GetExpirationTime().Format("2006-01-02T15:04:05Z07:00"))
				} else {
					data.ExpirationTime = types.StringNull()
				}

				if assignment.HasTimeZone() {
					data.TimeZone = types.StringValue(assignment.GetTimeZone())
				} else {
					data.TimeZone = types.StringNull()
				}

				if assignment.HasAssignmentType() {
					data.AssignmentType = types.StringValue(string(assignment.GetAssignmentType()))
				}

				break
			}
		}
	}

	if !found {
		resp.Diagnostics.AddError(
			"Assignment not found",
			"Could not find assignment with ID "+assignmentId+" in collection "+collectionId,
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

package governance

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/okta-governance-sdk-golang/governance"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ resource.Resource                = &collectionAssignmentResource{}
	_ resource.ResourceWithConfigure   = &collectionAssignmentResource{}
	_ resource.ResourceWithImportState = &collectionAssignmentResource{}
)

func newCollectionAssignmentResource() resource.Resource {
	return &collectionAssignmentResource{}
}

type collectionAssignmentResource struct {
	*config.Config
}

type collectionAssignmentModel struct {
	Id             types.String `tfsdk:"id"`
	CollectionId   types.String `tfsdk:"collection_id"`
	PrincipalId    types.String `tfsdk:"principal_id"`
	PrincipalType  types.String `tfsdk:"principal_type"`
	Actor          types.String `tfsdk:"actor"`
	ExpirationTime types.String `tfsdk:"expiration_time"`
	TimeZone       types.String `tfsdk:"time_zone"`
	AssignmentType types.String `tfsdk:"assignment_type"`
}

func (r *collectionAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_collection_assignment"
}

func (r *collectionAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an assignment of a collection to a principal (user or group). This creates POLICY grants automatically.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The assignment ID.",
			},
			"collection_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the collection to assign.",
			},
			"principal_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the principal (user ID or group ID).",
			},
			"principal_type": schema.StringAttribute{
				Required:    true,
				Description: "The type of principal. Valid values: OKTA_USER, OKTA_GROUP.",
			},
			"actor": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The actor making the assignment. Valid values: ACCESS_REQUEST, ADMIN, API, NONE. Default: API.",
			},
			"expiration_time": schema.StringAttribute{
				Optional:    true,
				Description: "The date/time when the assignment expires (ISO 8601 format).",
			},
			"time_zone": schema.StringAttribute{
				Optional:    true,
				Description: "The time zone in IANA format for the expiration date.",
			},
			"assignment_type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of assignment (e.g., INDIVIDUAL).",
			},
		},
	}
}

func (r *collectionAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data collectionAssignmentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	assignment := governance.NewAssignedPrincipal()
	principal := governance.NewTargetPrincipalFull(
		data.PrincipalId.ValueString(),
		governance.PrincipalType(data.PrincipalType.ValueString()),
	)
	assignment.SetPrincipal(*principal)

	actor := "API"
	if !data.Actor.IsNull() {
		actor = data.Actor.ValueString()
	}
	assignment.SetActor(governance.GrantActor(actor))

	if !data.TimeZone.IsNull() && data.ExpirationTime.IsNull() {
		resp.Diagnostics.AddError("Invalid configuration", "time_zone requires expiration_time to be set")
		return
	}
	if !data.ExpirationTime.IsNull() {
		t, err := time.Parse(time.RFC3339, data.ExpirationTime.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Invalid expiration_time format", err.Error())
			return
		}
		if t.Before(time.Now().UTC()) {
			resp.Diagnostics.AddError("Invalid expiration_time", "expiration_time must be a future timestamp")
			return
		}
		assignment.SetExpirationTime(t)
	}
	if !data.TimeZone.IsNull() {
		assignment.SetTimeZone(data.TimeZone.ValueString())
	}

	assignments, _, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().CollectionsAPI.
		AssignCollection(ctx, data.CollectionId.ValueString()).
		AssignedPrincipal([]governance.AssignedPrincipal{*assignment}).
		Execute()

	if err != nil {
		resp.Diagnostics.AddError("Error assigning Collection", err.Error())
		return
	}

	if len(assignments) > 0 {
		applyCollectionAssignmentToState(&data, &assignments[0])
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *collectionAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data collectionAssignmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	escapedPrincipal := strings.ReplaceAll(data.PrincipalId.ValueString(), "\\", "\\\\")
	escapedPrincipal = strings.ReplaceAll(escapedPrincipal, "\"", "\\\"")
	filter := fmt.Sprintf("principal.externalId eq \"%s\" AND principal.type eq \"%s\"",
		escapedPrincipal, data.PrincipalType.ValueString())

	assignments, _, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().CollectionsAPI.
		ListCollectionAssignments(ctx, data.CollectionId.ValueString()).
		Filter(filter).
		Execute()

	if err != nil {
		resp.Diagnostics.AddError("Error reading Collection Assignment", err.Error())
		return
	}

	if assignments.HasData() && len(assignments.GetData()) > 0 {
		for _, assignment := range assignments.GetData() {
			if assignment.HasId() && assignment.GetId() == data.Id.ValueString() {
				applyCollectionAssignmentDetailsToState(&data, &assignment)
				break
			}
		}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *collectionAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data collectionAssignmentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var operations []governance.AssignmentPatchOperation
	if !data.TimeZone.IsNull() && data.ExpirationTime.IsNull() {
		resp.Diagnostics.AddError("Invalid configuration", "time_zone requires expiration_time to be set")
		return
	}
	if !data.ExpirationTime.IsNull() {
		t, err := time.Parse(time.RFC3339, data.ExpirationTime.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Invalid expiration_time format", err.Error())
			return
		}
		if t.Before(time.Now().UTC()) {
			resp.Diagnostics.AddError("Invalid expiration_time", "expiration_time must be a future timestamp")
			return
		}
		op := governance.NewAssignmentPatchOperation("REPLACE", "/expirationTime")
		op.SetValue(data.ExpirationTime.ValueString())
		operations = append(operations, *op)
	}
	if !data.TimeZone.IsNull() {
		op := governance.NewAssignmentPatchOperation("REPLACE", "/timeZone")
		op.SetValue(data.TimeZone.ValueString())
		operations = append(operations, *op)
	}

	_, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().CollectionsAPI.
		UpdatePrincipalAssignment(ctx, data.CollectionId.ValueString(), data.Id.ValueString()).
		AssignmentPatchOperation(operations).
		Execute()

	if err != nil {
		resp.Diagnostics.AddError("Error updating Collection Assignment", err.Error())
		return
	}

	// Re-read
	var readReq resource.ReadRequest
	readReq.State = req.State
	var readResp resource.ReadResponse
	readResp.State = resp.State
	r.Read(ctx, readReq, &readResp)
	resp.Diagnostics = readResp.Diagnostics
	resp.State = readResp.State
}

func (r *collectionAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data collectionAssignmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().CollectionsAPI.
		DeletePrincipalAssignment(ctx, data.CollectionId.ValueString(), data.Id.ValueString()).
		Execute()

	if err != nil {
		resp.Diagnostics.AddError("Error deleting Collection Assignment", err.Error())
	}
}

func (r *collectionAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid Import ID", "Format: collection_id/assignment_id")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("collection_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}

func (r *collectionAssignmentResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(request, response)
}

func applyCollectionAssignmentToState(data *collectionAssignmentModel, assignment *governance.AssignedPrincipalFull) {
	if assignment.HasId() {
		data.Id = types.StringValue(assignment.GetId())
	}
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
		data.ExpirationTime = types.StringValue(assignment.GetExpirationTime().Format(time.RFC3339))
	}
	if assignment.HasTimeZone() {
		data.TimeZone = types.StringValue(assignment.GetTimeZone())
	}
	if assignment.HasAssignmentType() {
		data.AssignmentType = types.StringValue(string(assignment.GetAssignmentType()))
	}
}

func applyCollectionAssignmentDetailsToState(data *collectionAssignmentModel, assignment *governance.AssignedPrincipalDetails) {
	if assignment.HasId() {
		data.Id = types.StringValue(assignment.GetId())
	}
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
		data.ExpirationTime = types.StringValue(assignment.GetExpirationTime().Format(time.RFC3339))
	}
	if assignment.HasTimeZone() {
		data.TimeZone = types.StringValue(assignment.GetTimeZone())
	}
	if assignment.HasAssignmentType() {
		data.AssignmentType = types.StringValue(string(assignment.GetAssignmentType()))
	}
}

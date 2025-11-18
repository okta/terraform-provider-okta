package governance

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/okta-governance-sdk-golang/governance"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ resource.Resource                = &collectionResource{}
	_ resource.ResourceWithConfigure   = &collectionResource{}
	_ resource.ResourceWithImportState = &collectionResource{}
)

func newCollectionResource() resource.Resource {
	return &collectionResource{}
}

type collectionResource struct {
	*config.Config
}

type collectionResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Created     types.String `tfsdk:"created"`
	CreatedBy   types.String `tfsdk:"created_by"`
	LastUpdated types.String `tfsdk:"last_updated"`
	UpdatedBy   types.String `tfsdk:"updated_by"`
}

func (r *collectionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_collection"
}

func (r *collectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a resource collection. Collections allow an admin to assign multiple resources at one time to a user or group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier of the collection.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the resource collection (1-255 characters).",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The human-readable description of the collection (1-1000 characters).",
			},
			"created": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the collection was created.",
			},
			"created_by": schema.StringAttribute{
				Computed:    true,
				Description: "User ID of who created the collection.",
			},
			"last_updated": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the collection was last updated.",
			},
			"updated_by": schema.StringAttribute{
				Computed:    true,
				Description: "User ID of who last updated the collection.",
			},
		},
	}
}

func (r *collectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data collectionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build the create request
	createReq := governance.NewCollectionCreatable(data.Name.ValueString())
	if !data.Description.IsNull() {
		createReq.SetDescription(data.Description.ValueString())
	}

	// Create API call
	collection, _, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().CollectionsAPI.CreateCollection(ctx).CollectionCreatable(*createReq).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Collection",
			"Could not create Collection, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response to state
	applyCollectionToState(&data, collection)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *collectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data collectionResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call
	collection, _, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().CollectionsAPI.GetCollection(ctx, data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Collection",
			"Could not read Collection, unexpected error: "+err.Error(),
		)
		return
	}

	applyCollectionToState(&data, collection)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *collectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data collectionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the update request
	updateReq := governance.NewCollectionUpdatable(data.Name.ValueString())
	if !data.Description.IsNull() {
		updateReq.SetDescription(data.Description.ValueString())
	}

	// Update API call
	collection, _, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().CollectionsAPI.ReplaceCollection(ctx, data.Id.ValueString()).CollectionUpdatable(*updateReq).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Collection",
			"Could not update Collection, unexpected error: "+err.Error(),
		)
		return
	}

	applyCollectionToState(&data, collection)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *collectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data collectionResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call
	_, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().CollectionsAPI.DeleteCollection(ctx, data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Collection",
			"Could not delete Collection, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *collectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *collectionResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(request, response)
}

func applyCollectionToState(data *collectionResourceModel, collection *governance.CollectionFull) {
	data.Id = types.StringValue(collection.GetId())
	data.Name = types.StringValue(collection.GetName())

	if collection.HasDescription() {
		data.Description = types.StringValue(collection.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	data.Created = types.StringValue(collection.Created.Format(time.RFC3339))
	data.CreatedBy = types.StringValue(collection.GetCreatedBy())
	data.LastUpdated = types.StringValue(collection.LastUpdated.Format(time.RFC3339))
	data.UpdatedBy = types.StringValue(collection.GetLastUpdatedBy())
}

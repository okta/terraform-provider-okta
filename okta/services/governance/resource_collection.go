package governance

import (
	"context"

	"example.com/aditya-okta/okta-ig-sdk-golang/governance"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

func (r *collectionResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(request, response)
}

type collectionResource struct {
	*config.Config
}

type collectionResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

func (r *collectionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_collection"
}

func (r *collectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (r *collectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data collectionResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic
	createCollectionResp, _, err := r.OktaGovernanceClient.OktaIGSDKClient().CollectionsAPI.CreateCollection(ctx).CollectionCreatable(createCollection(data)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Collections",
			"Could not create Collections, unexpected error: "+err.Error(),
		)
		return
	}
	r.applyCollectionsToState(ctx, &data, createCollectionResp)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *collectionResource) applyCollectionsToState(ctx context.Context, data *collectionResourceModel, createCollectinoResp *governance.CollectionFull) diag.Diagnostics {
	var diags diag.Diagnostics
	data.Id = types.StringValue(createCollectinoResp.Id)
	data.Name = types.StringValue(createCollectinoResp.Name)
	data.Description = types.StringPointerValue(createCollectinoResp.Description)
	return diags
}

func createCollection(data collectionResourceModel) governance.CollectionCreatable {
	return governance.CollectionCreatable{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueStringPointer(),
	}
}

func (r *collectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data collectionResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	getCollectionResp, _, err := r.OktaGovernanceClient.OktaIGSDKClient().CollectionsAPI.GetCollection(ctx, data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Collections",
			"Could not reading Collections, unexpected error: "+err.Error(),
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(r.applyCollectionsToState(ctx, &data, getCollectionResp)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *collectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data collectionResourceModel
	var state collectionResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic
	data.Id = state.Id
	updatedCollectionResp, _, err := r.OktaGovernanceClient.OktaIGSDKClient().CollectionsAPI.ReplaceCollection(ctx, state.Id.ValueString()).CollectionUpdatable(buildUpdateCollection(data)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Collections",
			"Could not update Collections, unexpected error: "+err.Error(),
		)
		return
	}
	// Save updated data into Terraform state
	r.applyCollectionsToState(ctx, &data, updatedCollectionResp)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func buildUpdateCollection(data collectionResourceModel) governance.CollectionUpdatable {
	return governance.CollectionUpdatable{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueStringPointer(),
	}
}

func (r *collectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data collectionResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	r.OktaGovernanceClient.OktaIGSDKClient().CollectionsAPI.DeleteCollection(ctx, data.Id.ValueString())
}

func (r *collectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

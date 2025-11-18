package governance

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/okta-governance-sdk-golang/governance"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ resource.Resource                = &collectionResourceResource{}
	_ resource.ResourceWithConfigure   = &collectionResourceResource{}
	_ resource.ResourceWithImportState = &collectionResourceResource{}
)

func newCollectionResourceResource() resource.Resource {
	return &collectionResourceResource{}
}

type collectionResourceResource struct {
	*config.Config
}

type collectionResourceResourceModel struct {
	CollectionId types.String                                 `tfsdk:"collection_id"`
	ResourceId   types.String                                 `tfsdk:"resource_id"`
	ResourceOrn  types.String                                 `tfsdk:"resource_orn"`
	Entitlements []collectionResourceResourceEntitlementModel `tfsdk:"entitlements"`
}

type collectionResourceResourceEntitlementModel struct {
	Id     types.String                                      `tfsdk:"id"`
	Values []collectionResourceResourceEntitlementValueModel `tfsdk:"values"`
}

type collectionResourceResourceEntitlementValueModel struct {
	Id types.String `tfsdk:"id"`
}

func (r *collectionResourceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_collection_resource"
}

func (r *collectionResourceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a resource within a collection. Adds applications with specific entitlements to a collection.",
		Attributes: map[string]schema.Attribute{
			"collection_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the collection.",
			},
			"resource_id": schema.StringAttribute{
				Computed:    true,
				Description: "The computed resource ID within the collection.",
			},
			"resource_orn": schema.StringAttribute{
				Required:    true,
				Description: "The ORN identifier for the app resource (e.g., 'orn:okta:idp:00o...:apps:salesforce:0oa...').",
			},
		},
		Blocks: map[string]schema.Block{
			"entitlements": schema.ListNestedBlock{
				Description: "List of entitlements with their values to assign to this resource.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Required:    true,
							Description: "The entitlement ID.",
						},
					},
					Blocks: map[string]schema.Block{
						"values": schema.ListNestedBlock{
							Description: "List of entitlement value IDs.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Required:    true,
										Description: "The entitlement value ID.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *collectionResourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data collectionResourceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Basic ORN validation
	orn := data.ResourceOrn.ValueString()
	if !strings.HasPrefix(orn, "orn:okta:idp:") || len(strings.Split(orn, ":")) < 6 {
		resp.Diagnostics.AddError("Invalid resource_orn", fmt.Sprintf("Unexpected format: %s", orn))
		return
	}
	// Build the create request
	createReq := buildCollectionResourceCreateRequest(data)

	// Add resources to collection
	resourcesList, _, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().CollectionsAPI.
		AddResourcesToCollection(ctx, data.CollectionId.ValueString()).
		CollectionResourceCreatable([]governance.CollectionResourceCreatable{createReq}).
		Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error adding resource to Collection",
			"Could not add resource to Collection, unexpected error: "+err.Error(),
		)
		return
	}

	// The API returns a list, we expect one item
	if resourcesList.HasData() && len(resourcesList.GetData()) > 0 {
		resourcesData := resourcesList.GetData()
		applyCollectionResourceToState(&data, &resourcesData[0])
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *collectionResourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data collectionResourceResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call
	resourceResp, _, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().CollectionsAPI.
		GetCollectionResource(ctx, data.CollectionId.ValueString(), data.ResourceId.ValueString()).
		Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Collection Resource",
			"Could not read Collection Resource, unexpected error: "+err.Error(),
		)
		return
	}

	applyCollectionResourceToState(&data, resourceResp)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *collectionResourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data collectionResourceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// ORN validation (immutable but ensure still valid)
	orn := data.ResourceOrn.ValueString()
	if !strings.HasPrefix(orn, "orn:okta:idp:") || len(strings.Split(orn, ":")) < 6 {
		resp.Diagnostics.AddError("Invalid resource_orn", fmt.Sprintf("Unexpected format: %s", orn))
		return
	}
	// Build the update request
	updateReq := buildCollectionResourceUpdateRequest(data)

	// Update API call
	resourceResp, _, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().CollectionsAPI.
		ReplaceCollectionResource(ctx, data.CollectionId.ValueString(), data.ResourceId.ValueString()).
		CollectionResourceUpdatable(updateReq).
		Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Collection Resource",
			"Could not update Collection Resource, unexpected error: "+err.Error(),
		)
		return
	}

	applyCollectionResourceToState(&data, resourceResp)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *collectionResourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data collectionResourceResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call
	_, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().CollectionsAPI.
		DeleteCollectionResource(ctx, data.CollectionId.ValueString(), data.ResourceId.ValueString()).
		Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Collection Resource",
			"Could not delete Collection Resource, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *collectionResourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import format: collection_id/resource_id
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be in the format 'collection_id/resource_id', got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("collection_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("resource_id"), parts[1])...)
}

func (r *collectionResourceResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(request, response)
}

func buildCollectionResourceCreateRequest(data collectionResourceResourceModel) governance.CollectionResourceCreatable {
	req := governance.NewCollectionResourceCreatable(data.ResourceOrn.ValueString())

	if len(data.Entitlements) > 0 {
		// deterministic ordering
		sort.SliceStable(data.Entitlements, func(i, j int) bool {
			return data.Entitlements[i].Id.ValueString() < data.Entitlements[j].Id.ValueString()
		})
		entitlements := make([]governance.EntitlementCreatable, len(data.Entitlements))
		for i, ent := range data.Entitlements {
			entitlement := governance.NewEntitlementCreatable()
			entitlement.SetId(ent.Id.ValueString())
			if len(ent.Values) > 0 {
				sort.SliceStable(ent.Values, func(a, b int) bool { return ent.Values[a].Id.ValueString() < ent.Values[b].Id.ValueString() })
				values := make([]governance.EntitlementValueCreatable, len(ent.Values))
				for j, val := range ent.Values {
					value := governance.NewEntitlementValueCreatable()
					value.SetId(val.Id.ValueString())
					values[j] = *value
				}
				entitlement.SetValues(values)
			}
			entitlements[i] = *entitlement
		}
		req.SetEntitlements(entitlements)
	}

	return *req
}

func buildCollectionResourceUpdateRequest(data collectionResourceResourceModel) governance.CollectionResourceUpdatable {
	req := governance.NewCollectionResourceUpdatable()

	if len(data.Entitlements) > 0 {
		// deterministic ordering
		sort.SliceStable(data.Entitlements, func(i, j int) bool {
			return data.Entitlements[i].Id.ValueString() < data.Entitlements[j].Id.ValueString()
		})
		entitlements := make([]governance.EntitlementCreatable, len(data.Entitlements))
		for i, ent := range data.Entitlements {
			entitlement := governance.NewEntitlementCreatable()
			entitlement.SetId(ent.Id.ValueString())
			if len(ent.Values) > 0 {
				sort.SliceStable(ent.Values, func(a, b int) bool { return ent.Values[a].Id.ValueString() < ent.Values[b].Id.ValueString() })
				values := make([]governance.EntitlementValueCreatable, len(ent.Values))
				for j, val := range ent.Values {
					value := governance.NewEntitlementValueCreatable()
					value.SetId(val.Id.ValueString())
					values[j] = *value
				}
				entitlement.SetValues(values)
			}
			entitlements[i] = *entitlement
		}
		req.SetEntitlements(entitlements)
	}

	return *req
}

func applyCollectionResourceToState(data *collectionResourceResourceModel, resource *governance.CollectionResourceFull) {
	data.ResourceOrn = types.StringValue(resource.GetResourceOrn())

	if resource.HasResourceId() {
		data.ResourceId = types.StringValue(resource.GetResourceId())
	}

	// Map entitlements if present
	if resource.HasEntitlements() {
		entitlements := resource.GetEntitlements()
		sort.SliceStable(entitlements, func(i, j int) bool { return entitlements[i].GetId() < entitlements[j].GetId() })
		data.Entitlements = make([]collectionResourceResourceEntitlementModel, len(entitlements))

		for i, ent := range entitlements {
			entModel := collectionResourceResourceEntitlementModel{
				Id: types.StringValue(ent.GetId()),
			}

			if ent.HasValues() {
				values := ent.GetValues()
				entModel.Values = make([]collectionResourceResourceEntitlementValueModel, len(values))
				for j, val := range values {
					entModel.Values[j] = collectionResourceResourceEntitlementValueModel{
						Id: types.StringValue(val.GetId()),
					}
				}
			}

			data.Entitlements[i] = entModel
		}
	}
}

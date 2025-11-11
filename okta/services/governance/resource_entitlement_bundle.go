package governance

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/okta-governance-sdk-golang/governance"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ resource.Resource                = &entitlementBundleResource{}
	_ resource.ResourceWithConfigure   = &entitlementBundleResource{}
	_ resource.ResourceWithImportState = &entitlementBundleResource{}
)

func newEntitlementBundleResource() resource.Resource {
	return &entitlementBundleResource{}
}

type entitlementBundleResource struct {
	*config.Config
}

func (r *entitlementBundleResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
}

func (r *entitlementBundleResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(request, response)
}

type TargetResourceModel struct {
	ExternalId types.String `tfsdk:"external_id"`
	Type       types.String `tfsdk:"type"`
}

type entitlementBundleResourceModel struct {
	Id                types.String         `tfsdk:"id"`
	Name              types.String         `tfsdk:"name"`
	Target            *TargetResourceModel `tfsdk:"target"`
	TargetResourceOrn types.String         `tfsdk:"target_resource_orn"`
	Description       types.String         `tfsdk:"description"`
	Entitlements      []entitlements       `tfsdk:"entitlements"`
	Status            types.String         `tfsdk:"status"`
}

type entitlements struct {
	Id     types.String `tfsdk:"id"`
	Values []valueBlock `tfsdk:"values"`
}

type valueBlock struct {
	Id types.String `tfsdk:"id"`
}

func (r *entitlementBundleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entitlement_bundle"
}

func (r *entitlementBundleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique identifier of the entitlement bundle",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the entitlement bundle",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the entitlement bundle",
				Optional:    true,
			},
			"target_resource_orn": schema.StringAttribute{
				Description: "The ORN of the target resource. Required when updating the entitlement bundle",
				Optional:    true,
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "status of the entitlement bundle",
				Computed:    true,
				Optional:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"target": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"external_id": schema.StringAttribute{
						Description: "External ID of the target resource",
						Required:    true,
					},
					"type": schema.StringAttribute{
						Description: "Type of the target resource",
						Required:    true,
					},
				},
			},
			"entitlements": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Entitlement ID",
							Required:    true,
						},
					},
					Blocks: map[string]schema.Block{
						"values": schema.SetNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Entitlement value ID",
										Required:    true,
									},
								},
							},
						},
					},
				},
				Description: "Collection of entitlements and their values",
			},
		},
	}
}

func (r *entitlementBundleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data entitlementBundleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic
	entitlementBundle, _, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().EntitlementBundlesAPI.CreateEntitlementBundle(ctx).EntitlementBundleCreatable(buildEntitlementBundleCreateBody(data)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Entitlement Bundles",
			"Could not create Entitlement Bundles, unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(applyEntitlementBundleToState2(ctx, entitlementBundle, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *entitlementBundleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data entitlementBundleResourceModel

	// Read Terraform prior state Data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	getEntitlementBundleResp, _, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().EntitlementBundlesAPI.GetentitlementBundle(ctx, data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading campaign",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(applyEntitlementBundleToState(ctx, getEntitlementBundleResp, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *entitlementBundleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state entitlementBundleResourceModel
	// Read Terraform prior state Data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Read Terraform prior state Data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Id = state.Id // Ensure the ID is set from the prior state
	data.Status = state.Status
	data.TargetResourceOrn = state.TargetResourceOrn

	// Update API call logic
	replaceEntitlementBundleResp, _, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().EntitlementBundlesAPI.ReplaceEntitlementBundle(ctx, state.Id.ValueString()).EntitlementBundleUpdatable(buildEntitlementBundleUpdateBody(data)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Entitlement Bundle",
			"Could not update Entitlement Bundle with ID "+data.Id.ValueString()+", unexpected error: "+err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(applyEntitlementBundleToState2(ctx, replaceEntitlementBundleResp, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Save updated Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *entitlementBundleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data entitlementBundleResourceModel

	// Read Terraform prior state Data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	_, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().EntitlementBundlesAPI.DeleteEntitlementBundle(ctx, data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Entitlement Bundle",
			"Could not delete Entitlement Bundle with ID "+data.Id.ValueString()+", unexpected error: "+err.Error(),
		)
		return
	}
}

func buildEntitlementBundleCreateBody(data entitlementBundleResourceModel) governance.EntitlementBundleCreatable {
	rt := governance.ResourceType2(data.Target.Type.ValueString())
	name := data.Name.ValueString()
	description := data.Description.ValueStringPointer()
	target := governance.TargetResource{
		ExternalId: data.Target.ExternalId.ValueString(),
		Type:       rt,
	}
	entitlements := make([]governance.EntitlementCreatable, 0, len(data.Entitlements))
	for _, ent := range data.Entitlements {
		values := make([]governance.EntitlementValueCreatable, 0, len(ent.Values))
		for _, val := range ent.Values {
			values = append(values, governance.EntitlementValueCreatable{
				Id: val.Id.ValueStringPointer(),
			})
		}
		entitlements = append(entitlements, governance.EntitlementCreatable{
			Id:     ent.Id.ValueStringPointer(),
			Values: values,
		})
	}

	return governance.EntitlementBundleCreatable{
		Name:         name,
		Description:  description,
		Target:       target,
		Entitlements: entitlements,
	}
}

func applyEntitlementBundleToState(ctx context.Context, data *governance.EntitlementBundleFullWithEntitlements, state *entitlementBundleResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	state.Id = types.StringValue(data.GetId())
	state.Name = types.StringValue(data.GetName())
	state.Description = types.StringValue(data.GetDescription())
	if target, ok := data.GetTargetOk(); ok {
		state.Target = &TargetResourceModel{
			ExternalId: types.StringValue(target.GetExternalId()),
			Type:       types.StringValue(string(target.GetType())),
		}
	}
	e := make([]entitlements, 0, len(data.GetEntitlements()))
	for _, ent := range data.GetEntitlements() {
		vals := make([]valueBlock, 0, len(ent.GetValues()))
		for _, v := range ent.GetValues() {
			vals = append(vals, valueBlock{
				Id: types.StringValue(v.GetId()),
			})
		}
		e = append(e, entitlements{
			Id:     types.StringValue(ent.GetId()),
			Values: vals,
		})
	}
	state.Entitlements = e
	state.TargetResourceOrn = types.StringValue(data.GetTargetResourceOrn())
	state.Status = types.StringValue(string(data.GetStatus()))
	return diags
}

func applyEntitlementBundleToState2(ctx context.Context, data *governance.EntitlementBundleFull, state *entitlementBundleResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	state.Id = types.StringValue(data.GetId())
	state.Name = types.StringValue(data.GetName())
	state.Description = types.StringValue(data.GetDescription())
	if target, ok := data.GetTargetOk(); ok {
		state.Target = &TargetResourceModel{
			ExternalId: types.StringValue(target.GetExternalId()),
			Type:       types.StringValue(string(target.GetType())),
		}
	}
	e := make([]entitlements, 0, len(data.GetEntitlements()))
	for _, ent := range data.GetEntitlements() {
		vals := make([]valueBlock, 0, len(ent.GetValues()))
		for _, v := range ent.GetValues() {
			vals = append(vals, valueBlock{
				Id: types.StringValue(v.GetId()),
			})
		}
		e = append(e, entitlements{
			Id:     types.StringValue(ent.GetId()),
			Values: vals,
		})
	}
	state.Entitlements = e
	state.TargetResourceOrn = types.StringValue(data.GetTargetResourceOrn())
	state.Status = types.StringValue(string(data.GetStatus()))
	return diags
}

func buildEntitlementBundleUpdateBody(data entitlementBundleResourceModel) governance.EntitlementBundleUpdatable {
	rt := governance.ResourceType2(data.Target.Type.ValueString())
	status := governance.EntitlementBundleStatus(data.Status.ValueString())
	name := data.Name.ValueString()
	description := data.Description.ValueStringPointer()
	target := governance.TargetResource{
		ExternalId: data.Target.ExternalId.ValueString(),
		Type:       rt,
	}
	targetResourceOrn := data.TargetResourceOrn.ValueString()
	entitlements := make([]governance.EntitlementCreatable, 0, len(data.Entitlements))
	if data.Entitlements != nil || len(data.Entitlements) > 0 {
		for _, ent := range data.Entitlements {
			values := make([]governance.EntitlementValueCreatable, 0, len(ent.Values))
			for _, val := range ent.Values {
				values = append(values, governance.EntitlementValueCreatable{
					Id: val.Id.ValueStringPointer(),
				})
			}
			entitlements = append(entitlements, governance.EntitlementCreatable{
				Id:     ent.Id.ValueStringPointer(),
				Values: values,
			})
		}
	}
	return governance.EntitlementBundleUpdatable{
		Id:                data.Id.ValueString(),
		Name:              name,
		Description:       description,
		Target:            target,
		TargetResourceOrn: targetResourceOrn,
		Entitlements:      entitlements,
		Status:            &status,
	}
}

// Copyright 2025 - Present Okta, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package idaas

import (
	"context"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	fwdiag "github.com/hashicorp/terraform-plugin-framework/diag"
	frameworkPath "github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"

	okta "github.com/okta/okta-sdk-golang/v6/okta"

	"github.com/okta/terraform-provider-okta/okta/config"
)

// Ensure interface compliance
var (
	_ resource.Resource                = &trustedOriginResource{}
	_ resource.ResourceWithConfigure   = &trustedOriginResource{}
	_ resource.ResourceWithImportState = &trustedOriginResource{}
)

// trustedOriginResource defines the resource implementation.
type trustedOriginResource struct {
	Config *config.Config
}

// trustedOriginModel describes the resource data model.
type trustedOriginModel struct {
	ID            types.String                    `tfsdk:"id"`
	Created       types.String                    `tfsdk:"created"`
	CreatedBy     types.String                    `tfsdk:"created_by"`
	LastUpdated   types.String                    `tfsdk:"last_updated"`
	LastUpdatedBy types.String                    `tfsdk:"last_updated_by"`
	Name          types.String                    `tfsdk:"name"`
	Origin        types.String                    `tfsdk:"origin"`
	Scopes        []TrustedOriginModelScopesModel `tfsdk:"scopes"`
	Status        types.String                    `tfsdk:"status"`
}

// TrustedOriginModelScopesModel is the nested model for scopes.
type TrustedOriginModelScopesModel struct {
	AllowedOktaApps types.List   `tfsdk:"allowed_okta_apps"`
	Type            types.String `tfsdk:"type"`
}

func newTrustedOriginResource() resource.Resource {
	return &trustedOriginResource{}
}

func (r *trustedOriginResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_trusted_origin"
}

func (r *trustedOriginResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *trustedOriginResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The Trusted Origins API provides operations to manage Trusted Origins and sources.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier for the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created": schema.StringAttribute{
				Description: "Timestamp when the trusted origin was created",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_by": schema.StringAttribute{
				Description: "The ID of the user who created the trusted origin",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Description: "Timestamp when the trusted origin was last updated",
				Computed:    true,
			},
			"last_updated_by": schema.StringAttribute{
				Description: "The ID of the user who last updated the trusted origin",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Unique name for the trusted origin",
				Required:    true,
			},
			"origin": schema.StringAttribute{
				Description: "Unique origin URL for the trusted origin.",
				Required:    true,
			},
			"status": schema.StringAttribute{
				Description: "Status of the trusted origin. Values: ACTIVE, INACTIVE",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"scopes": schema.ListNestedBlock{
				Description: "Array of scope types that this trusted origin is used for",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"allowed_okta_apps": schema.ListAttribute{
							Description: "The allowed Okta apps for the trusted origin scope",
							ElementType: types.StringType,
							Optional:    true,
						},
						"type": schema.StringAttribute{
							Description: "The scope type.",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func (r *trustedOriginResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, frameworkPath.Root("id"), req, resp)
}

// mapAPIScopes maps API TrustedOriginScope slice to TrustedOriginModelScopesModel slice.
func mapAPIScopes(apiScopes []okta.TrustedOriginScope, diagnostics *fwdiag.Diagnostics) []TrustedOriginModelScopesModel {
	var scopesState []TrustedOriginModelScopesModel
	for _, scope := range apiScopes {
		var appsListVal types.List
		if apiApps := scope.GetAllowedOktaApps(); len(apiApps) > 0 {
			var allowedApps []attr.Value
			for _, app := range apiApps {
				allowedApps = append(allowedApps, types.StringValue(app))
			}
			var diags fwdiag.Diagnostics
			appsListVal, diags = types.ListValue(types.StringType, allowedApps)
			diagnostics.Append(diags...)
		} else {
			appsListVal = types.ListNull(types.StringType)
		}
		scopesState = append(scopesState, TrustedOriginModelScopesModel{
			Type:            types.StringValue(scope.GetType()),
			AllowedOktaApps: appsListVal,
		})
	}
	return scopesState
}

// buildRequestScopes builds the SDK TrustedOriginScope slice from plan model.
func buildRequestScopes(planScopes []TrustedOriginModelScopesModel) []okta.TrustedOriginScope {
	var scopesSlice []okta.TrustedOriginScope
	for _, item := range planScopes {
		nestedItem := okta.NewTrustedOriginScopeWithDefaults()
		if !item.AllowedOktaApps.IsNull() && !item.AllowedOktaApps.IsUnknown() {
			var allowedOktaAppsItemSlice []string
			for _, elem := range item.AllowedOktaApps.Elements() {
				if sv, ok := elem.(types.String); ok {
					allowedOktaAppsItemSlice = append(allowedOktaAppsItemSlice, sv.ValueString())
				}
			}
			nestedItem.SetAllowedOktaApps(allowedOktaAppsItemSlice)
		}
		if !item.Type.IsNull() && !item.Type.IsUnknown() {
			nestedItem.SetType(item.Type.ValueString())
		}
		scopesSlice = append(scopesSlice, *nestedItem)
	}
	return scopesSlice
}

func (r *trustedOriginResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state trustedOriginModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id := state.ID.ValueString()

	client := r.Config.OktaIDaaSClient.OktaSDKClientV6()
	result, httpResp, err := client.TrustedOriginAPI.GetTrustedOrigin(ctx, id).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading trusted_origin", err.Error())
		return
	}
	state.Created = types.StringValue(result.GetCreated().Format(time.RFC3339))
	state.CreatedBy = types.StringValue(result.GetCreatedBy())
	state.LastUpdated = types.StringValue(result.GetLastUpdated().Format(time.RFC3339))
	state.LastUpdatedBy = types.StringValue(result.GetLastUpdatedBy())
	state.Name = types.StringValue(result.GetName())
	state.Origin = types.StringValue(result.GetOrigin())
	state.Status = types.StringValue(result.GetStatus())
	state.ID = types.StringValue(result.GetId())
	state.Scopes = mapAPIScopes(result.GetScopes(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *trustedOriginResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan trustedOriginModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Config.OktaIDaaSClient.OktaSDKClientV6()

	createReq := client.TrustedOriginAPI.CreateTrustedOrigin(ctx)
	body := okta.NewTrustedOriginWriteWithDefaults()
	body.SetName(plan.Name.ValueString())
	body.SetOrigin(plan.Origin.ValueString())
	if len(plan.Scopes) > 0 {
		body.SetScopes(buildRequestScopes(plan.Scopes))
	}
	createReq = createReq.TrustedOrigin(*body)
	result, _, err := createReq.Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error creating trusted_origin", err.Error())
		return
	}
	plan.ID = types.StringValue(result.GetId())
	plan.Created = types.StringValue(result.GetCreated().Format(time.RFC3339))
	plan.CreatedBy = types.StringValue(result.GetCreatedBy())
	plan.LastUpdated = types.StringValue(result.GetLastUpdated().Format(time.RFC3339))
	plan.LastUpdatedBy = types.StringValue(result.GetLastUpdatedBy())
	plan.Name = types.StringValue(result.GetName())
	plan.Origin = types.StringValue(result.GetOrigin())
	plan.Status = types.StringValue(result.GetStatus())
	plan.Scopes = mapAPIScopes(result.GetScopes(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *trustedOriginResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan trustedOriginModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state trustedOriginModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id := state.ID.ValueString()

	client := r.Config.OktaIDaaSClient.OktaSDKClientV6()

	updateReq := client.TrustedOriginAPI.ReplaceTrustedOrigin(ctx, id)
	updateBody := okta.NewTrustedOriginWithDefaults()
	updateBody.SetName(plan.Name.ValueString())
	updateBody.SetOrigin(plan.Origin.ValueString())
	if len(plan.Scopes) > 0 {
		updateBody.SetScopes(buildRequestScopes(plan.Scopes))
	}
	if !plan.Status.IsNull() && !plan.Status.IsUnknown() {
		updateBody.SetStatus(plan.Status.ValueString())
	}
	updateReq = updateReq.TrustedOrigin(*updateBody)
	result, _, err := updateReq.Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error updating trusted_origin", err.Error())
		return
	}
	state.ID = types.StringValue(result.GetId())
	state.Created = types.StringValue(result.GetCreated().Format(time.RFC3339))
	state.CreatedBy = types.StringValue(result.GetCreatedBy())
	state.LastUpdated = types.StringValue(result.GetLastUpdated().Format(time.RFC3339))
	state.LastUpdatedBy = types.StringValue(result.GetLastUpdatedBy())
	state.Name = types.StringValue(result.GetName())
	state.Origin = types.StringValue(result.GetOrigin())
	state.Status = types.StringValue(result.GetStatus())
	state.Scopes = mapAPIScopes(result.GetScopes(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *trustedOriginResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state trustedOriginModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id := state.ID.ValueString()

	client := r.Config.OktaIDaaSClient.OktaSDKClientV6()
	httpResp, err := client.TrustedOriginAPI.DeleteTrustedOrigin(ctx, id).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			return
		}
		resp.Diagnostics.AddError("Error deleting trusted_origin", err.Error())
		return
	}
}

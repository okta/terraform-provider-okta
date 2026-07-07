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

// Handwritten backwards-compatible resource for okta_trusted_origin.
// Maintains the original schema (active bool, flat scopes list) from the SDK-v2 era
// while using the v6 SDK and Terraform Plugin Framework internally.
// See resource_okta_trusted_origin_generated.go for the API-faithful generated version.
package idaas

import (
	"context"
	"net/http"

	frameworkPath "github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	okta "github.com/okta/okta-sdk-golang/v6/okta"

	"github.com/okta/terraform-provider-okta/okta/config"
)

// Ensure interface compliance
var (
	_ resource.Resource                = &trustedOriginCompatResource{}
	_ resource.ResourceWithConfigure   = &trustedOriginCompatResource{}
	_ resource.ResourceWithImportState = &trustedOriginCompatResource{}
	_ resource.ResourceWithUpgradeState = &trustedOriginCompatResource{}
)

// trustedOriginCompatResource is the backwards-compatible handwritten implementation.
type trustedOriginCompatResource struct {
	Config *config.Config
}

// trustedOriginCompatModel holds the backwards-compatible schema state.
type trustedOriginCompatModel struct {
	ID     types.String `tfsdk:"id"`
	Active types.Bool   `tfsdk:"active"`
	Name   types.String `tfsdk:"name"`
	Origin types.String `tfsdk:"origin"`
	Scopes types.List   `tfsdk:"scopes"` // list(string)
}

// newTrustedOriginCompatResource returns the backwards-compatible resource.
// Registered in FWProviderResources() instead of the generated NewTrustedOriginResource.
func newTrustedOriginCompatResource() resource.Resource {
	return &trustedOriginCompatResource{}
}

func (r *trustedOriginCompatResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_trusted_origin"
}

func (r *trustedOriginCompatResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *trustedOriginCompatResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Version:     1,
		Description: "Creates a Trusted Origin. This resource allows you to create and configure a Trusted Origin.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier for the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"active": schema.BoolAttribute{
				Description: "Whether the Trusted Origin is active or not - can only be issued post-creation. By default, it is `true`.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"name": schema.StringAttribute{
				Description: "Unique name for this trusted origin.",
				Required:    true,
			},
			"origin": schema.StringAttribute{
				Description: "Unique origin URL for this trusted origin.",
				Required:    true,
			},
			"scopes": schema.ListAttribute{
				Description: "Scopes of the Trusted Origin - can either be `CORS` and/or `REDIRECT`.",
				Required:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (r *trustedOriginCompatResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, frameworkPath.Root("id"), req, resp)
}

// UpgradeState migrates state from schema version 0 (old Plugin SDK v2) to version 1 (this resource).
// Since the attribute names and types are identical, this is a pass-through upgrade.
func (r *trustedOriginCompatResource) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
	v0Schema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":     schema.StringAttribute{Computed: true},
			"active": schema.BoolAttribute{Optional: true, Computed: true},
			"name":   schema.StringAttribute{Required: true},
			"origin": schema.StringAttribute{Required: true},
			"scopes": schema.ListAttribute{Required: true, ElementType: types.StringType},
		},
	}
	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema: &v0Schema,
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				var old trustedOriginCompatModel
				resp.Diagnostics.Append(req.State.Get(ctx, &old)...)
				if resp.Diagnostics.HasError() {
					return
				}
				resp.Diagnostics.Append(resp.State.Set(ctx, &old)...)
			},
		},
	}
}

func (r *trustedOriginCompatResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan trustedOriginCompatModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.Active.ValueBool() {
		resp.Diagnostics.AddError("Invalid configuration", "cannot create an inactive trusted origin; only existing trusted origins can be deactivated")
		return
	}

	body := okta.NewTrustedOriginWriteWithDefaults()
	body.SetName(plan.Name.ValueString())
	body.SetOrigin(plan.Origin.ValueString())
	body.SetScopes(buildTrustedOriginScopesV6(ctx, plan.Scopes))

	client := r.Config.OktaIDaaSClient.OktaSDKClientV6()
	result, _, err := client.TrustedOriginAPI.CreateTrustedOrigin(ctx).TrustedOrigin(*body).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error creating trusted origin", err.Error())
		return
	}

	plan.ID = types.StringValue(result.GetId())
	plan.Active = types.BoolValue(result.GetStatus() == StatusActive)
	plan.Name = types.StringValue(result.GetName())
	plan.Origin = types.StringValue(result.GetOrigin())
	plan.Scopes = flattenTrustedOriginScopesV6(result.GetScopes())

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *trustedOriginCompatResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state trustedOriginCompatModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Config.OktaIDaaSClient.OktaSDKClientV6()
	result, httpResp, err := client.TrustedOriginAPI.GetTrustedOrigin(ctx, state.ID.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading trusted origin", err.Error())
		return
	}

	state.Active = types.BoolValue(result.GetStatus() == StatusActive)
	state.Name = types.StringValue(result.GetName())
	state.Origin = types.StringValue(result.GetOrigin())
	state.Scopes = flattenTrustedOriginScopesV6(result.GetScopes())

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *trustedOriginCompatResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan trustedOriginCompatModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state trustedOriginCompatModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	client := r.Config.OktaIDaaSClient.OktaSDKClientV6()

	// Handle activate/deactivate if active changed
	if !plan.Active.Equal(state.Active) {
		var err error
		if plan.Active.ValueBool() {
			_, _, err = client.TrustedOriginAPI.ActivateTrustedOrigin(ctx, id).Execute()
		} else {
			_, _, err = client.TrustedOriginAPI.DeactivateTrustedOrigin(ctx, id).Execute()
		}
		if err != nil {
			resp.Diagnostics.AddError("Error changing trusted origin status", err.Error())
			return
		}
	}

	updateBody := okta.NewTrustedOriginWithDefaults()
	updateBody.SetName(plan.Name.ValueString())
	updateBody.SetOrigin(plan.Origin.ValueString())
	updateBody.SetScopes(buildTrustedOriginScopesV6(ctx, plan.Scopes))

	result, _, err := client.TrustedOriginAPI.ReplaceTrustedOrigin(ctx, id).TrustedOrigin(*updateBody).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error updating trusted origin", err.Error())
		return
	}

	state.Active = types.BoolValue(result.GetStatus() == StatusActive)
	state.Name = types.StringValue(result.GetName())
	state.Origin = types.StringValue(result.GetOrigin())
	state.Scopes = flattenTrustedOriginScopesV6(result.GetScopes())

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *trustedOriginCompatResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state trustedOriginCompatModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Config.OktaIDaaSClient.OktaSDKClientV6()
	httpResp, err := client.TrustedOriginAPI.DeleteTrustedOrigin(ctx, state.ID.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			return
		}
		resp.Diagnostics.AddError("Error deleting trusted origin", err.Error())
	}
}

// buildTrustedOriginScopesV6 converts a flat list(string) of scope types to v6 SDK scope objects.
func buildTrustedOriginScopesV6(ctx context.Context, scopesList types.List) []okta.TrustedOriginScope {
	var scopeStrings []string
	_ = scopesList.ElementsAs(ctx, &scopeStrings, false)
	scopes := make([]okta.TrustedOriginScope, len(scopeStrings))
	for i, s := range scopeStrings {
		scope := okta.NewTrustedOriginScopeWithDefaults()
		scope.SetType(s)
		scopes[i] = *scope
	}
	return scopes
}

// flattenTrustedOriginScopesV6 converts v6 SDK scope objects back to a flat list(string) of scope types.
func flattenTrustedOriginScopesV6(scopes []okta.TrustedOriginScope) types.List {
	scopeStrings := make([]string, len(scopes))
	for i, s := range scopes {
		scopeStrings[i] = s.GetType()
	}
	result, _ := types.ListValueFrom(context.Background(), types.StringType, scopeStrings)
	return result
}


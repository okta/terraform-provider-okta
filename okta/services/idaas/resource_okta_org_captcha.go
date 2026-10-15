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

	frameworkPath "github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	okta "github.com/okta/okta-sdk-golang/v6/okta"

	"github.com/okta/terraform-provider-okta/okta/config"
)

// Ensure interface compliance
var (
	_ resource.Resource                = &orgCaptchaResource{}
	_ resource.ResourceWithConfigure   = &orgCaptchaResource{}
	_ resource.ResourceWithImportState = &orgCaptchaResource{}
)

// OrgCaptchaResource defines the resource implementation.
type orgCaptchaResource struct {
	Config *config.Config
}

// OrgCaptchaModel describes the resource data model.
type orgCaptchaModel struct {
	ID           types.String `tfsdk:"id"`
	CaptchaId    types.String `tfsdk:"captcha_id"`
	EnabledPages types.List   `tfsdk:"enabled_pages"`
}

func newOrgCaptchaResource() resource.Resource {
	return &orgCaptchaResource{}
}

func (r *orgCaptchaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_org_captcha"
}

func (r *orgCaptchaResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *orgCaptchaResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "As an option to increase org security, Okta supports CAPTCHA services to prevent automated sign-in attempts.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier for the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"captcha_id": schema.StringAttribute{
				Description: "The unique key of the associated CAPTCHA instance",
				Optional:    true,
			},
			"enabled_pages": schema.ListAttribute{
				Description: "An array of pages that have CAPTCHA enabled",
				ElementType: types.StringType,
				Optional:    true,
			},
		},
	}
}

func (r *orgCaptchaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, frameworkPath.Root("id"), req, resp)
}

func (r *orgCaptchaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state orgCaptchaModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Config.OktaIDaaSClient.OktaSDKClientV6()
	result, httpResp, err := client.CAPTCHAAPI.GetOrgCaptchaSettings(ctx).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading org_captcha", err.Error())
		return
	}
	// Map API response fields to state (scalar types only; WriteOnly fields are skipped — response type doesn't have them)
	state.CaptchaId = types.StringValue(string(result.GetCaptchaId()))

	state.ID = types.StringValue("org_captcha")

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *orgCaptchaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan orgCaptchaModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Config.OktaIDaaSClient.OktaSDKClientV6()

	// Build request body from plan
	createReq := client.CAPTCHAAPI.ReplacesOrgCaptchaSettings(ctx)
	body := okta.NewOrgCAPTCHASettingsWithDefaults()
	body.SetCaptchaId(plan.CaptchaId.ValueString())
	if !plan.EnabledPages.IsNull() && !plan.EnabledPages.IsUnknown() {
		var enabledPagesSlice []string
		for _, item := range plan.EnabledPages.Elements() {
			if sv, ok := item.(types.String); ok {
				enabledPagesSlice = append(enabledPagesSlice, sv.ValueString())
			}
		}
		body.SetEnabledPages(enabledPagesSlice)
	}
	createReq = createReq.OrgCAPTCHASettings(*body)
	result, _, err := createReq.Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error creating org_captcha", err.Error())
		return
	}
	plan.ID = types.StringValue("org_captcha")
	// Map response fields back to plan (scalar types only; WriteOnly and SkipRead fields skipped)
	plan.CaptchaId = types.StringValue(string(result.GetCaptchaId()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
func (r *orgCaptchaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan orgCaptchaModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state orgCaptchaModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Config.OktaIDaaSClient.OktaSDKClientV6()

	// Build request body from plan — only send changed fields
	updateReq := client.CAPTCHAAPI.ReplacesOrgCaptchaSettings(ctx)
	updateBody := okta.NewOrgCAPTCHASettingsWithDefaults()
	updateBody.SetCaptchaId(plan.CaptchaId.ValueString())
	if !plan.EnabledPages.IsNull() && !plan.EnabledPages.IsUnknown() {
		var enabledPagesSlice []string
		for _, item := range plan.EnabledPages.Elements() {
			if sv, ok := item.(types.String); ok {
				enabledPagesSlice = append(enabledPagesSlice, sv.ValueString())
			}
		}
		updateBody.SetEnabledPages(enabledPagesSlice)
	}
	updateReq = updateReq.OrgCAPTCHASettings(*updateBody)
	result, _, err := updateReq.Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error updating org_captcha", err.Error())
		return
	}
	// Map API response fields to state (scalar types only; WriteOnly and SkipRead fields skipped)
	state.ID = types.StringValue("org_captcha")
	state.CaptchaId = types.StringValue(string(result.GetCaptchaId()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *orgCaptchaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state orgCaptchaModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Config.OktaIDaaSClient.OktaSDKClientV6()
	httpResp, err := client.CAPTCHAAPI.DeleteOrgCaptchaSettings(ctx).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			return
		}
		resp.Diagnostics.AddError("Error deleting org_captcha", err.Error())
		return
	}
}

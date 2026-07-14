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
	_ resource.Resource                = &inlineHookResource{}
	_ resource.ResourceWithConfigure   = &inlineHookResource{}
	_ resource.ResourceWithImportState = &inlineHookResource{}
)

// InlineHookResource defines the resource implementation.
type inlineHookResource struct {
	Config *config.Config
}

// InlineHookModel describes the resource data model.
type inlineHookModel struct {
	ID          types.String                 `tfsdk:"id"`
	Channel     *InlineHookModelChannelModel `tfsdk:"channel"`
	Created     types.String                 `tfsdk:"created"`
	LastUpdated types.String                 `tfsdk:"last_updated"`
	Name        types.String                 `tfsdk:"name"`
	Status      types.String                 `tfsdk:"status"`
	Type        types.String                 `tfsdk:"type"`
	Version     types.String                 `tfsdk:"version"`
}

// InlineHookModelChannelModel is the nested model for channel.
type InlineHookModelChannelModel struct {
	Type    types.String `tfsdk:"type"`
	Version types.String `tfsdk:"version"`
}

func newInlineHookResource() resource.Resource {
	return &inlineHookResource{}
}

func (r *inlineHookResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_inline_hook"
}

func (r *inlineHookResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *inlineHookResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The Inline Hooks API provides operations to manage inline hooks for your organization.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier for the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created": schema.StringAttribute{
				Description: "Date of the inline hook creation",
				Computed:    true,
			},
			"last_updated": schema.StringAttribute{
				Description: "Date of the last inline hook update",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The display name of the inline hook",
				Optional:    true,
			},
			"status": schema.StringAttribute{
				Description: "Status",
				Optional:    true,
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "One of the inline hook types",
				Optional:    true,
			},
			"version": schema.StringAttribute{
				Description: "Version of the inline hook type.",
				Optional:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"channel": schema.SingleNestedBlock{
				Description: "Channel",
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description: "Type",
						Optional:    true,
					},
					"version": schema.StringAttribute{
						Description: "Version of the inline hook type.",
						Optional:    true,
					},
				},
			},
		},
	}
}
func (r *inlineHookResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, frameworkPath.Root("id"), req, resp)
}

func (r *inlineHookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state inlineHookModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id := state.ID.ValueString()

	client := r.Config.OktaIDaaSClient.OktaSDKClientV6()
	result, httpResp, err := client.InlineHookAPI.GetInlineHook(ctx, id).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading inline_hook", err.Error())
		return
	}
	// Map API response fields to state (scalar types only; WriteOnly fields are skipped — response type doesn't have them)
	state.Created = types.StringValue(result.GetCreated().Format(time.RFC3339))
	state.LastUpdated = types.StringValue(result.GetLastUpdated().Format(time.RFC3339))
	state.Name = types.StringValue(string(result.GetName()))
	state.Status = types.StringValue(string(result.GetStatus()))
	state.Type = types.StringValue(string(result.GetType()))
	state.Version = types.StringValue(string(result.GetVersion()))
	if channelRaw0, ok := result.GetChannelOk(); ok {
		channelModel0 := &InlineHookModelChannelModel{}
		channelModel0.Type = types.StringValue(string(channelRaw0.GetType()))
		channelModel0.Version = types.StringValue(string(channelRaw0.GetVersion()))
		state.Channel = channelModel0
	}

	state.ID = types.StringValue(string(result.GetId()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *inlineHookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan inlineHookModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Config.OktaIDaaSClient.OktaSDKClientV6()

	// Build request body from plan
	createReq := client.InlineHookAPI.CreateInlineHook(ctx)
	body := okta.NewInlineHookCreateWithDefaults()
	if plan.Channel != nil {
		nestedChannel := okta.NewInlineHookChannelCreateWithDefaults()
		if !plan.Channel.Type.IsNull() && !plan.Channel.Type.IsUnknown() {
			nestedChannel.SetType(plan.Channel.Type.ValueString())
		}
		if !plan.Channel.Version.IsNull() && !plan.Channel.Version.IsUnknown() {
			nestedChannel.SetVersion(plan.Channel.Version.ValueString())
		}
		body.SetChannel(*nestedChannel)
	}
	body.SetName(plan.Name.ValueString())
	body.SetType(plan.Type.ValueString())
	body.SetVersion(plan.Version.ValueString())
	createReq = createReq.InlineHookCreate(*body)
	result, _, err := createReq.Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error creating inline_hook", err.Error())
		return
	}
	// Set ID from API response
	plan.ID = types.StringValue(string(result.GetId()))
	// Map response fields back to plan (scalar types only; WriteOnly and SkipRead fields skipped)
	plan.Created = types.StringValue(result.GetCreated().Format(time.RFC3339))
	plan.LastUpdated = types.StringValue(result.GetLastUpdated().Format(time.RFC3339))
	plan.Name = types.StringValue(string(result.GetName()))
	plan.Status = types.StringValue(string(result.GetStatus()))
	plan.Type = types.StringValue(string(result.GetType()))
	plan.Version = types.StringValue(string(result.GetVersion()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
func (r *inlineHookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan inlineHookModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state inlineHookModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id := state.ID.ValueString()

	client := r.Config.OktaIDaaSClient.OktaSDKClientV6()

	// Build request body from plan — only send changed fields
	updateReq := client.InlineHookAPI.ReplaceInlineHook(ctx, id)
	updateBody := okta.NewInlineHookReplaceWithDefaults()
	if plan.Channel != nil {
		nestedChannel := okta.NewInlineHookChannelCreateWithDefaults()
		if !plan.Channel.Type.IsNull() && !plan.Channel.Type.IsUnknown() {
			nestedChannel.SetType(plan.Channel.Type.ValueString())
		}
		if !plan.Channel.Version.IsNull() && !plan.Channel.Version.IsUnknown() {
			nestedChannel.SetVersion(plan.Channel.Version.ValueString())
		}
		updateBody.SetChannel(*nestedChannel)
	}
	updateBody.SetName(plan.Name.ValueString())
	updateBody.SetVersion(plan.Version.ValueString())
	updateReq = updateReq.InlineHook(*updateBody)
	result, _, err := updateReq.Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error updating inline_hook", err.Error())
		return
	}
	// Map API response fields to state (scalar types only; WriteOnly and SkipRead fields skipped)
	state.ID = types.StringValue(string(result.GetId()))
	state.Created = types.StringValue(result.GetCreated().Format(time.RFC3339))
	state.LastUpdated = types.StringValue(result.GetLastUpdated().Format(time.RFC3339))
	state.Name = types.StringValue(string(result.GetName()))
	state.Status = types.StringValue(string(result.GetStatus()))
	state.Type = types.StringValue(string(result.GetType()))
	state.Version = types.StringValue(string(result.GetVersion()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *inlineHookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state inlineHookModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id := state.ID.ValueString()

	client := r.Config.OktaIDaaSClient.OktaSDKClientV6()
	httpResp, err := client.InlineHookAPI.DeleteInlineHook(ctx, id).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			return
		}
		resp.Diagnostics.AddError("Error deleting inline_hook", err.Error())
		return
	}
}

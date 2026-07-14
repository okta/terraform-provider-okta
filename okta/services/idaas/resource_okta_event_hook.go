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
	_ resource.Resource                = &eventHookResource{}
	_ resource.ResourceWithConfigure   = &eventHookResource{}
	_ resource.ResourceWithImportState = &eventHookResource{}
)

// EventHookResource defines the resource implementation.
type eventHookResource struct {
	Config *config.Config
}

// EventHookModel describes the resource data model.
type eventHookModel struct {
	ID                 types.String                `tfsdk:"id"`
	Channel            *EventHookModelChannelModel `tfsdk:"channel"`
	Created            types.String                `tfsdk:"created"`
	CreatedBy          types.String                `tfsdk:"created_by"`
	Description        types.String                `tfsdk:"description"`
	Events             *EventHookModelEventsModel  `tfsdk:"events"`
	LastUpdated        types.String                `tfsdk:"last_updated"`
	Name               types.String                `tfsdk:"name"`
	Status             types.String                `tfsdk:"status"`
	VerificationStatus types.String                `tfsdk:"verification_status"`
}

// EventHookModelChannelModel is the nested model for channel.
type EventHookModelChannelModel struct {
	Config  *EventHookModelChannelModelConfigModel `tfsdk:"config"`
	Type    types.String                           `tfsdk:"type"`
	Version types.String                           `tfsdk:"version"`
}

// EventHookModelChannelModelConfigModel is the nested model for config.
type EventHookModelChannelModelConfigModel struct {
	AuthScheme *EventHookModelChannelModelConfigModelAuthSchemeModel `tfsdk:"auth_scheme"`
	Headers    []EventHookModelChannelModelConfigModelHeadersModel   `tfsdk:"headers"`
	Method     types.String                                          `tfsdk:"method"`
	Uri        types.String                                          `tfsdk:"uri"`
}

// EventHookModelChannelModelConfigModelAuthSchemeModel is the nested model for auth_scheme.
type EventHookModelChannelModelConfigModelAuthSchemeModel struct {
	Key   types.String `tfsdk:"key"`
	Type  types.String `tfsdk:"type"`
	Value types.String `tfsdk:"value"`
}

// EventHookModelChannelModelConfigModelHeadersModel is the nested model for headers.
type EventHookModelChannelModelConfigModelHeadersModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

// EventHookModelEventsModel is the nested model for events.
type EventHookModelEventsModel struct {
	Filter *EventHookModelEventsModelFilterModel `tfsdk:"filter"`
	Items  types.List                            `tfsdk:"items"`
	Type   types.String                          `tfsdk:"type"`
}

// EventHookModelEventsModelFilterModel is the nested model for filter.
type EventHookModelEventsModelFilterModel struct {
	EventFilterMap []EventHookModelEventsModelFilterModelEventFilterMapModel `tfsdk:"event_filter_map"`
	Type           types.String                                              `tfsdk:"type"`
}

// EventHookModelEventsModelFilterModelEventFilterMapModel is the nested model for event_filter_map.
type EventHookModelEventsModelFilterModelEventFilterMapModel struct {
	Condition *EventHookModelEventsModelFilterModelEventFilterMapModelConditionModel `tfsdk:"condition"`
	Event     types.String                                                           `tfsdk:"event"`
}

// EventHookModelEventsModelFilterModelEventFilterMapModelConditionModel is the nested model for condition.
type EventHookModelEventsModelFilterModelEventFilterMapModelConditionModel struct {
	Expression types.String `tfsdk:"expression"`
	Version    types.String `tfsdk:"version"`
}

func newEventHookResource() resource.Resource {
	return &eventHookResource{}
}

func (r *eventHookResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_event_hook"
}

func (r *eventHookResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *eventHookResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The Event Hooks API provides operations to manage event hooks for your organization.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier for the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created": schema.StringAttribute{
				Description: "Timestamp of the event hook creation",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_by": schema.StringAttribute{
				Description: "The ID of the user who created the event hook",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description of the event hook",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Description: "Date of the last event hook update",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Display name for the event hook",
				Required:    true,
			},
			"status": schema.StringAttribute{
				Description: "Status of the event hook",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"verification_status": schema.StringAttribute{
				Description: "Verification status of the event hook.",
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"channel": schema.SingleNestedBlock{
				Description: "Channel",
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description: "The channel type.",
						Required:    true,
					},
					"version": schema.StringAttribute{
						Description: "Version of the channel.",
						Required:    true,
					},
				},
				Blocks: map[string]schema.Block{
					"config": schema.SingleNestedBlock{
						Description: "Config",
						Attributes: map[string]schema.Attribute{
							"method": schema.StringAttribute{
								Description: "The method of the Okta event hook request",
								Computed:    true,
							},
							"uri": schema.StringAttribute{
								Description: "The external service endpoint called to execute the event hook handler",
								Required:    true,
							},
						},
						Blocks: map[string]schema.Block{
							"auth_scheme": schema.SingleNestedBlock{
								Description: "The authentication scheme used for this request.",
								Attributes: map[string]schema.Attribute{
									"key": schema.StringAttribute{
										Description: "The name for the authorization header",
										Optional:    true,
									},
									"type": schema.StringAttribute{
										Description: "The authentication scheme type.",
										Optional:    true,
									},
									"value": schema.StringAttribute{
										Description: "The header value.",
										Optional:    true,
									},
								},
							},
							"headers": schema.ListNestedBlock{
								Description: "Optional list of key/value pairs for headers that can be sent with the request to the external service.",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"key": schema.StringAttribute{
											Description: "The optional field or header name",
											Optional:    true,
										},
										"value": schema.StringAttribute{
											Description: "The value for the key",
											Optional:    true,
										},
									},
								},
							},
						},
					},
				},
			},
			"events": schema.SingleNestedBlock{
				Description: "Events",
				Attributes: map[string]schema.Attribute{
					"items": schema.ListAttribute{
						Description: "The subscribed event types that trigger the event hook.",
						ElementType: types.StringType,
						Required:    true,
					},
					"type": schema.StringAttribute{
						Description: "The events object type.",
						Required:    true,
					},
				},
				Blocks: map[string]schema.Block{
					"filter": schema.SingleNestedBlock{
						Description: "The optional filter defined on a specific event type  > **Note:** Event hook filters is a [self-service Early Access (EA)](/openapi/okta-management/guides/release-lifecycle/#early-access-ea) to enable.",
						Attributes: map[string]schema.Attribute{
							"type": schema.StringAttribute{
								Description: "The type of filter.",
								Computed:    true,
							},
						},
						Blocks: map[string]schema.Block{
							"event_filter_map": schema.ListNestedBlock{
								Description: "The object that maps the filter to the event type",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"event": schema.StringAttribute{
											Description: "The filtered event type",
											Optional:    true,
										},
									},
									Blocks: map[string]schema.Block{
										"condition": schema.SingleNestedBlock{
											Description: "Condition",
											Attributes: map[string]schema.Attribute{
												"expression": schema.StringAttribute{
													Description: "The Okta Expression language statement that filters the event type",
													Optional:    true,
												},
												"version": schema.StringAttribute{
													Description: "Internal field",
													Computed:    true,
												},
											},
										},
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
func (r *eventHookResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, frameworkPath.Root("id"), req, resp)
}

func (r *eventHookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state eventHookModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id := state.ID.ValueString()

	client := r.Config.OktaIDaaSClient.OktaSDKClientV6()
	result, httpResp, err := client.EventHookAPI.GetEventHook(ctx, id).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading event_hook", err.Error())
		return
	}
	// Map API response fields to state (scalar types only; WriteOnly fields are skipped — response type doesn't have them)
	state.Created = types.StringValue(result.GetCreated().Format(time.RFC3339))
	state.CreatedBy = types.StringValue(string(result.GetCreatedBy()))
	state.Description = types.StringValue(string(result.GetDescription()))
	state.LastUpdated = types.StringValue(result.GetLastUpdated().Format(time.RFC3339))
	state.Name = types.StringValue(string(result.GetName()))
	state.Status = types.StringValue(string(result.GetStatus()))
	state.VerificationStatus = types.StringValue(string(result.GetVerificationStatus()))
	if channelRaw0, ok := result.GetChannelOk(); ok {
		channelModel0 := &EventHookModelChannelModel{}
		if configRaw1, ok := channelRaw0.GetConfigOk(); ok {
			configModel1 := &EventHookModelChannelModelConfigModel{}
			if authSchemeRaw2, ok := configRaw1.GetAuthSchemeOk(); ok {
				authSchemeModel2 := &EventHookModelChannelModelConfigModelAuthSchemeModel{}
				authSchemeModel2.Key = types.StringValue(string(authSchemeRaw2.GetKey()))
				authSchemeModel2.Type = types.StringValue(string(authSchemeRaw2.GetType()))
				authSchemeModel2.Value = types.StringValue(string(authSchemeRaw2.GetValue()))
				configModel1.AuthScheme = authSchemeModel2
			}
			configModel1.Method = types.StringValue(string(configRaw1.GetMethod()))
			configModel1.Uri = types.StringValue(string(configRaw1.GetUri()))
			channelModel0.Config = configModel1
		}
		channelModel0.Type = types.StringValue(string(channelRaw0.GetType()))
		channelModel0.Version = types.StringValue(string(channelRaw0.GetVersion()))
		state.Channel = channelModel0
	}
	if eventsRaw0, ok := result.GetEventsOk(); ok {
		eventsModel0 := &EventHookModelEventsModel{}
		if filterRaw1, ok := eventsRaw0.GetFilterOk(); ok && filterRaw1 != nil {
			filterModel1 := &EventHookModelEventsModelFilterModel{}
			filterModel1.Type = types.StringValue(string(filterRaw1.GetType()))
			eventsModel0.Filter = filterModel1
		}
		{
			listVal, listDiags := types.ListValueFrom(ctx, types.StringType, eventsRaw0.GetItems())
			resp.Diagnostics.Append(listDiags...)
			eventsModel0.Items = listVal
		}
		eventsModel0.Type = types.StringValue(string(eventsRaw0.GetType()))
		state.Events = eventsModel0
	}

	state.ID = types.StringValue(string(result.GetId()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *eventHookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan eventHookModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Config.OktaIDaaSClient.OktaSDKClientV6()

	// Build request body from plan
	createReq := client.EventHookAPI.CreateEventHook(ctx)
	body := okta.NewEventHookWithDefaults()
	if plan.Channel != nil {
		nestedChannel := okta.NewEventHookChannelWithDefaults()
		if !plan.Channel.Type.IsNull() && !plan.Channel.Type.IsUnknown() {
			nestedChannel.SetType(plan.Channel.Type.ValueString())
		}
		if !plan.Channel.Version.IsNull() && !plan.Channel.Version.IsUnknown() {
			nestedChannel.SetVersion(plan.Channel.Version.ValueString())
		}
		body.SetChannel(*nestedChannel)
	}
	body.SetDescription(plan.Description.ValueString())
	if plan.Events != nil {
		nestedEvents := okta.NewEventSubscriptionsWithDefaults()
		if !plan.Events.Items.IsNull() && !plan.Events.Items.IsUnknown() {
			var itemsSlice []string
			for _, elem := range plan.Events.Items.Elements() {
				if sv, ok := elem.(types.String); ok {
					itemsSlice = append(itemsSlice, sv.ValueString())
				}
			}
			nestedEvents.SetItems(itemsSlice)
		}
		if !plan.Events.Type.IsNull() && !plan.Events.Type.IsUnknown() {
			nestedEvents.SetType(plan.Events.Type.ValueString())
		}
		body.SetEvents(*nestedEvents)
	}
	body.SetName(plan.Name.ValueString())
	createReq = createReq.EventHook(*body)
	result, _, err := createReq.Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error creating event_hook", err.Error())
		return
	}
	// Set ID from API response
	plan.ID = types.StringValue(string(result.GetId()))
	// Map response fields back to plan (scalar types only; WriteOnly and SkipRead fields skipped)
	plan.Created = types.StringValue(result.GetCreated().Format(time.RFC3339))
	plan.CreatedBy = types.StringValue(string(result.GetCreatedBy()))
	plan.Description = types.StringValue(string(result.GetDescription()))
	plan.LastUpdated = types.StringValue(result.GetLastUpdated().Format(time.RFC3339))
	plan.Name = types.StringValue(string(result.GetName()))
	plan.Status = types.StringValue(string(result.GetStatus()))
	plan.VerificationStatus = types.StringValue(string(result.GetVerificationStatus()))
	if channelRaw0, ok := result.GetChannelOk(); ok {
		channelModel0 := &EventHookModelChannelModel{}
		if configRaw1, ok := channelRaw0.GetConfigOk(); ok {
			configModel1 := &EventHookModelChannelModelConfigModel{}
			if authSchemeRaw2, ok := configRaw1.GetAuthSchemeOk(); ok {
				authSchemeModel2 := &EventHookModelChannelModelConfigModelAuthSchemeModel{}
				authSchemeModel2.Key = types.StringValue(string(authSchemeRaw2.GetKey()))
				authSchemeModel2.Type = types.StringValue(string(authSchemeRaw2.GetType()))
				authSchemeModel2.Value = types.StringValue(string(authSchemeRaw2.GetValue()))
				configModel1.AuthScheme = authSchemeModel2
			}
			configModel1.Method = types.StringValue(string(configRaw1.GetMethod()))
			configModel1.Uri = types.StringValue(string(configRaw1.GetUri()))
			var headersSlice []EventHookModelChannelModelConfigModelHeadersModel
			for _, h := range configRaw1.GetHeaders() {
				headersSlice = append(headersSlice, EventHookModelChannelModelConfigModelHeadersModel{
					Key:   types.StringValue(string(h.GetKey())),
					Value: types.StringValue(string(h.GetValue())),
				})
			}
			configModel1.Headers = headersSlice
			channelModel0.Config = configModel1
		}
		channelModel0.Type = types.StringValue(string(channelRaw0.GetType()))
		channelModel0.Version = types.StringValue(string(channelRaw0.GetVersion()))
		plan.Channel = channelModel0
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
func (r *eventHookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan eventHookModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state eventHookModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id := state.ID.ValueString()

	client := r.Config.OktaIDaaSClient.OktaSDKClientV6()

	// Build request body from plan — only send changed fields
	updateReq := client.EventHookAPI.ReplaceEventHook(ctx, id)
	updateBody := okta.NewEventHookWithDefaults()
	if plan.Channel != nil {
		nestedChannel := okta.NewEventHookChannelWithDefaults()
		if !plan.Channel.Type.IsNull() && !plan.Channel.Type.IsUnknown() {
			nestedChannel.SetType(plan.Channel.Type.ValueString())
		}
		if !plan.Channel.Version.IsNull() && !plan.Channel.Version.IsUnknown() {
			nestedChannel.SetVersion(plan.Channel.Version.ValueString())
		}
		updateBody.SetChannel(*nestedChannel)
	}
	updateBody.SetDescription(plan.Description.ValueString())
	if plan.Events != nil {
		nestedEvents := okta.NewEventSubscriptionsWithDefaults()
		if !plan.Events.Items.IsNull() && !plan.Events.Items.IsUnknown() {
			var itemsSlice []string
			for _, elem := range plan.Events.Items.Elements() {
				if sv, ok := elem.(types.String); ok {
					itemsSlice = append(itemsSlice, sv.ValueString())
				}
			}
			nestedEvents.SetItems(itemsSlice)
		}
		if !plan.Events.Type.IsNull() && !plan.Events.Type.IsUnknown() {
			nestedEvents.SetType(plan.Events.Type.ValueString())
		}
		updateBody.SetEvents(*nestedEvents)
	}
	updateBody.SetName(plan.Name.ValueString())
	updateReq = updateReq.EventHook(*updateBody)
	result, _, err := updateReq.Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error updating event_hook", err.Error())
		return
	}
	// Map API response fields to state (scalar types only; WriteOnly and SkipRead fields skipped)
	state.ID = types.StringValue(string(result.GetId()))
	state.Created = types.StringValue(result.GetCreated().Format(time.RFC3339))
	state.CreatedBy = types.StringValue(string(result.GetCreatedBy()))
	state.Description = types.StringValue(string(result.GetDescription()))
	state.LastUpdated = types.StringValue(result.GetLastUpdated().Format(time.RFC3339))
	state.Name = types.StringValue(string(result.GetName()))
	state.Status = types.StringValue(string(result.GetStatus()))
	state.VerificationStatus = types.StringValue(string(result.GetVerificationStatus()))
	if channelRaw0, ok := result.GetChannelOk(); ok {
		channelModel0 := &EventHookModelChannelModel{}
		if configRaw1, ok := channelRaw0.GetConfigOk(); ok {
			configModel1 := &EventHookModelChannelModelConfigModel{}
			if authSchemeRaw2, ok := configRaw1.GetAuthSchemeOk(); ok {
				authSchemeModel2 := &EventHookModelChannelModelConfigModelAuthSchemeModel{}
				authSchemeModel2.Key = types.StringValue(string(authSchemeRaw2.GetKey()))
				authSchemeModel2.Type = types.StringValue(string(authSchemeRaw2.GetType()))
				authSchemeModel2.Value = types.StringValue(string(authSchemeRaw2.GetValue()))
				configModel1.AuthScheme = authSchemeModel2
			}
			configModel1.Method = types.StringValue(string(configRaw1.GetMethod()))
			configModel1.Uri = types.StringValue(string(configRaw1.GetUri()))
			var headersSlice []EventHookModelChannelModelConfigModelHeadersModel
			for _, h := range configRaw1.GetHeaders() {
				headersSlice = append(headersSlice, EventHookModelChannelModelConfigModelHeadersModel{
					Key:   types.StringValue(string(h.GetKey())),
					Value: types.StringValue(string(h.GetValue())),
				})
			}
			configModel1.Headers = headersSlice
			channelModel0.Config = configModel1
		}
		channelModel0.Type = types.StringValue(string(channelRaw0.GetType()))
		channelModel0.Version = types.StringValue(string(channelRaw0.GetVersion()))
		state.Channel = channelModel0
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *eventHookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state eventHookModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id := state.ID.ValueString()

	client := r.Config.OktaIDaaSClient.OktaSDKClientV6()
	httpResp, err := client.EventHookAPI.DeleteEventHook(ctx, id).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			return
		}
		resp.Diagnostics.AddError("Error deleting event_hook", err.Error())
		return
	}
}

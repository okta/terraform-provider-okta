package idaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/okta/okta-sdk-golang/v4/okta"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ resource.Resource                = &eventHookResource{}
	_ resource.ResourceWithConfigure   = &eventHookResource{}
	_ resource.ResourceWithImportState = &eventHookResource{}
)

func newEventHookResource() resource.Resource {
	return &eventHookResource{}
}

type eventHookResource struct {
	*config.Config
}

type eventHookResourceModel struct {
	ID      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Status  types.String `tfsdk:"status"`
	Events  types.Set    `tfsdk:"events"`
	Headers types.List   `tfsdk:"headers"`
	Auth    types.Object `tfsdk:"auth"`
	Channel types.Object `tfsdk:"channel"`
}

type eventHookAuthModel struct {
	Type  types.String `tfsdk:"type"`
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

type eventHookChannelModel struct {
	Type    types.String `tfsdk:"type"`
	Version types.String `tfsdk:"version"`
	URI     types.String `tfsdk:"uri"`
}

type eventHookHeaderModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

func (r *eventHookResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_event_hook"
}

func (r *eventHookResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates an event hook. This resource allows you to create and configure an event hook.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the event hook.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The event hook display name.",
				Required:    true,
			},
			"status": schema.StringAttribute{
				Description: "The status of the event hook. Valid values: ACTIVE, INACTIVE.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("ACTIVE"),
				Validators: []validator.String{
					stringvalidator.OneOf("ACTIVE", "INACTIVE"),
				},
			},
			"events": schema.SetAttribute{
				Description: "The events that will be delivered to this hook. See https://developer.okta.com/docs/reference/api/event-types/?q=event-hook-eligible for a list of supported events.",
				Required:    true,
				ElementType: types.StringType,
			},
		},
		Blocks: map[string]schema.Block{
			"headers": schema.ListNestedBlock{
				Description: "Map of headers to send along in event hook request.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Description: "The header key.",
							Required:    true,
						},
						"value": schema.StringAttribute{
							Description: "The header value.",
							Required:    true,
						},
					},
				},
			},
			"auth": schema.SingleNestedBlock{
				Description: "Authentication configuration for the event hook.",
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description: "The type of authentication. Currently only 'HEADER' is supported.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("HEADER"),
						Validators: []validator.String{
							stringvalidator.OneOf("HEADER"),
						},
					},
					"key": schema.StringAttribute{
						Description: "The authentication key (e.g., header name).",
						Optional:    true,
					},
					"value": schema.StringAttribute{
						Description: "The authentication value (e.g., API token). This field supports ephemeral values and will not be stored in state.",
						Optional:    true,
						WriteOnly:   true, // This enables ephemeral support
					},
				},
			},
			"channel": schema.SingleNestedBlock{
				Description: "Details of the endpoint the event hook will hit.",
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description: "The type of channel. Currently only 'HTTP' is supported.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("HTTP"),
						Validators: []validator.String{
							stringvalidator.OneOf("HTTP"),
						},
					},
					"version": schema.StringAttribute{
						Description: "The version of the channel. Currently only '1.0.0' is supported.",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.OneOf("1.0.0"),
						},
					},
					"uri": schema.StringAttribute{
						Description: "The URI the hook will hit.",
						Required:    true,
					},
				},
			},
		},
	}
}

func (r *eventHookResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *eventHookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan eventHookResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hook, err := r.buildEventHookFromModel(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to build event hook", err.Error())
		return
	}

	client := r.OktaIDaaSClient.OktaSDKClientV3()
	newHook, _, err := client.EventHookAPI.CreateEventHook(ctx).EventHook(*hook).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Failed to create event hook", err.Error())
		return
	}

	plan.ID = types.StringValue(*newHook.Id)

	// Set status if needed
	if !plan.Status.IsNull() && !plan.Status.IsUnknown() && *newHook.Status != plan.Status.ValueString() {
		err = r.setEventHookStatus(ctx, client, *newHook.Id, plan.Status.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Failed to set event hook status", err.Error())
			return
		}
	}

	// Read the created resource to get the final state
	r.readEventHook(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *eventHookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state eventHookResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.readEventHook(ctx, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *eventHookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan eventHookResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hook, err := r.buildEventHookFromModel(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to build event hook", err.Error())
		return
	}

	client := r.OktaIDaaSClient.OktaSDKClientV3()
	newHook, _, err := client.EventHookAPI.ReplaceEventHook(ctx, plan.ID.ValueString()).EventHook(*hook).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Failed to update event hook", err.Error())
		return
	}

	// Set status if needed
	if !plan.Status.IsNull() && !plan.Status.IsUnknown() && *newHook.Status != plan.Status.ValueString() {
		err = r.setEventHookStatus(ctx, client, *newHook.Id, plan.Status.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Failed to set event hook status", err.Error())
			return
		}
	}

	// Read the updated resource to get the final state
	r.readEventHook(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *eventHookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state eventHookResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.OktaIDaaSClient.OktaSDKClientV3()

	// Deactivate the event hook first
	_, _, err := client.EventHookAPI.DeactivateEventHook(ctx, state.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Failed to deactivate event hook", err.Error())
		return
	}

	// Delete the event hook
	_, err = client.EventHookAPI.DeleteEventHook(ctx, state.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete event hook", err.Error())
		return
	}
}

func (r *eventHookResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *eventHookResource) buildEventHookFromModel(ctx context.Context, model *eventHookResourceModel) (*okta.EventHook, error) {
	// Convert events set to slice
	var events []string
	if !model.Events.IsNull() && !model.Events.IsUnknown() {
		eventElements := make([]types.String, 0, len(model.Events.Elements()))
		diags := model.Events.ElementsAs(ctx, &eventElements, false)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to convert events: %v", diags)
		}
		for _, event := range eventElements {
			events = append(events, event.ValueString())
		}
	}

	// Build channel
	var channelModel eventHookChannelModel
	if !model.Channel.IsNull() && !model.Channel.IsUnknown() {
		diags := model.Channel.As(ctx, &channelModel, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return nil, fmt.Errorf("failed to convert channel: %v", diags)
		}
	}

	channel := okta.EventHookChannel{
		Type:    channelModel.Type.ValueString(),
		Version: channelModel.Version.ValueString(),
		Config: okta.EventHookChannelConfig{
			Uri: channelModel.URI.ValueString(),
		},
	}

	// Build auth if provided
	if !model.Auth.IsNull() && !model.Auth.IsUnknown() {
		var authModel eventHookAuthModel
		diags := model.Auth.As(ctx, &authModel, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return nil, fmt.Errorf("failed to convert auth: %v", diags)
		}

		// Only create auth scheme if both key and value are provided
		if !authModel.Key.IsNull() && !authModel.Key.IsUnknown() &&
			!authModel.Value.IsNull() && !authModel.Value.IsUnknown() {
			authType := authModel.Type.ValueString()
			authKey := authModel.Key.ValueString()
			authValue := authModel.Value.ValueString()
			channel.Config.AuthScheme = &okta.EventHookChannelConfigAuthScheme{
				Type:  &authType,
				Key:   &authKey,
				Value: &authValue,
			}
		}
	}

	// Build headers if provided
	if !model.Headers.IsNull() && !model.Headers.IsUnknown() {
		var headerModels []eventHookHeaderModel
		diags := model.Headers.ElementsAs(ctx, &headerModels, false)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to convert headers: %v", diags)
		}

		var headers []okta.EventHookChannelConfigHeader
		for _, headerModel := range headerModels {
			key := headerModel.Key.ValueString()
			value := headerModel.Value.ValueString()
			headers = append(headers, okta.EventHookChannelConfigHeader{
				Key:   &key,
				Value: &value,
			})
		}
		channel.Config.Headers = headers
	}

	status := model.Status.ValueString()
	return &okta.EventHook{
		Name:    model.Name.ValueString(),
		Status:  &status,
		Events:  okta.EventSubscriptions{Type: "EVENT_TYPE", Items: events},
		Channel: channel,
	}, nil
}

func (r *eventHookResource) readEventHook(ctx context.Context, model *eventHookResourceModel, diags *diag.Diagnostics) {
	client := r.OktaIDaaSClient.OktaSDKClientV3()
	hook, resp, err := client.EventHookAPI.GetEventHook(ctx, model.ID.ValueString()).Execute()
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			// Resource was deleted outside of Terraform
			model.ID = types.StringNull()
			return
		}
		diags.AddError("Failed to read event hook", err.Error())
		return
	}

	if hook == nil {
		model.ID = types.StringNull()
		return
	}

	model.Name = types.StringValue(hook.Name)
	if hook.Status != nil {
		model.Status = types.StringValue(*hook.Status)
	}

	// Convert events
	if hook.Events.Items != nil {
		eventValues := make([]attr.Value, len(hook.Events.Items))
		for i, event := range hook.Events.Items {
			eventValues[i] = types.StringValue(event)
		}
		model.Events = types.SetValueMust(types.StringType, eventValues)
	}

	// Convert channel
	channelAttrs := map[string]attr.Value{
		"type":    types.StringValue(hook.Channel.Type),
		"version": types.StringValue(hook.Channel.Version),
		"uri":     types.StringValue(hook.Channel.Config.Uri),
	}
	model.Channel = types.ObjectValueMust(
		map[string]attr.Type{
			"type":    types.StringType,
			"version": types.StringType,
			"uri":     types.StringType,
		},
		channelAttrs,
	)

	// Convert auth (but don't include the value since it's write-only)
	if hook.Channel.Config.AuthScheme != nil {
		var authType, authKey string
		if hook.Channel.Config.AuthScheme.Type != nil {
			authType = *hook.Channel.Config.AuthScheme.Type
		}
		if hook.Channel.Config.AuthScheme.Key != nil {
			authKey = *hook.Channel.Config.AuthScheme.Key
		}
		authAttrs := map[string]attr.Value{
			"type":  types.StringValue(authType),
			"key":   types.StringValue(authKey),
			"value": types.StringNull(), // Write-only field, don't read back
		}
		model.Auth = types.ObjectValueMust(
			map[string]attr.Type{
				"type":  types.StringType,
				"key":   types.StringType,
				"value": types.StringType,
			},
			authAttrs,
		)
	}

	// Convert headers
	if hook.Channel.Config.Headers != nil {
		headerValues := make([]attr.Value, len(hook.Channel.Config.Headers))
		for i, header := range hook.Channel.Config.Headers {
			var key, value string
			if header.Key != nil {
				key = *header.Key
			}
			if header.Value != nil {
				value = *header.Value
			}
			headerAttrs := map[string]attr.Value{
				"key":   types.StringValue(key),
				"value": types.StringValue(value),
			}
			headerValues[i] = types.ObjectValueMust(
				map[string]attr.Type{
					"key":   types.StringType,
					"value": types.StringType,
				},
				headerAttrs,
			)
		}
		model.Headers = types.ListValueMust(
			types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"key":   types.StringType,
					"value": types.StringType,
				},
			},
			headerValues,
		)
	}
}

func (r *eventHookResource) setEventHookStatus(ctx context.Context, client *okta.APIClient, id, desiredStatus string) error {
	if desiredStatus == "INACTIVE" {
		_, _, err := client.EventHookAPI.DeactivateEventHook(ctx, id).Execute()
		return err
	} else {
		_, _, err := client.EventHookAPI.ActivateEventHook(ctx, id).Execute()
		return err
	}
}

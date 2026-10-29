package idaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	v5okta "github.com/okta/okta-sdk-golang/v5/okta"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

const (
	filterTypeExpressionLanguage = "EXPRESSION_LANGUAGE"
	eventTypeSubscription        = "EVENT_TYPE"
	defaultChannelType           = "HTTP"
	defaultChannelVersion        = "1.0.0"
	defaultAuthType              = "HEADER"
)

var eventHookHeaderSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"key": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The header key.",
		},
		"value": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The header value.",
		},
	},
}

func resourceEventHook() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEventHookCreate,
		ReadContext:   resourceEventHookRead,
		UpdateContext: resourceEventHookUpdate,
		DeleteContext: resourceEventHookDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "Creates an event hook. This resource allows you to create and configure an event hook.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The event hook display name.",
			},
			"status": statusSchema,
			"events": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The events that will be delivered to this hook. [See here for a list of supported events](https://developer.okta.com/docs/reference/api/event-types/?q=event-hook-eligible).",
			},
			"filter": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Optional filter to reduce the number of event hook calls using Okta Expression Language. This is a self-service Early Access (EA) feature. See [Event hook filters](https://developer.okta.com/docs/concepts/event-hooks/#create-an-event-hook-filter) for more information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The type of filter to use. Currently only 'EXPRESSION_LANGUAGE' is supported by the API.",
							Default:     filterTypeExpressionLanguage,
						},
						"event": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The event type to filter.",
						},
						"condition": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The Okta Expression Language statement that filters the event type.",
						},
					},
				},
			},
			"headers": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        HeaderSchema,
				Description: "Map of headers to send along in event hook request.",
			},
			"auth": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// Suppress diffs for default values when not explicitly set
					switch k {
					case "auth.type":
						return (old == defaultAuthType && new == "") || (old == "" && new == defaultAuthType)
					}
					return false
				},
				Description: `Authentication scheme for the event hook endpoint.   
	- 'key' - (Required) The key to use for authentication.
	- 'value' - (Required) The value or secret for authentication.
	- 'type' - (Optional) The type of authentication. Currently, the only supported type is 'HEADER'. Default: 'HEADER'.`,
			},
			"channel": {
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// Suppress diffs for default values when not explicitly set
					switch k {
					case "channel.type":
						return (old == defaultChannelType && new == "") || (old == "" && new == defaultChannelType)
					case "channel.version":
						return (old == defaultChannelVersion && new == "") || (old == "" && new == defaultChannelVersion)
					}
					return false
				},
				Description: `Details of the endpoint the event hook will hit.   
	- 'uri' - (Required) The URI the hook will hit.
	- 'type' - (Optional) The type of hook to trigger. Currently, the only supported type is 'HTTP'. Default: 'HTTP'.
	- 'version' - (Optional) The version of the channel. The currently-supported version is '1.0.0'. Default: '1.0.0'.`,
			},
		},
	}
}

func resourceEventHookCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := getOktaV5ClientFromMetadata(meta)
	hook := buildEventHook(d)
	newHook, _, err := client.EventHookAPI.CreateEventHook(ctx).EventHook(*hook).Execute()
	if err != nil {
		return diag.Errorf("failed to create event hook: %v", err)
	}
	d.SetId(*newHook.Id)
	resp, err := setEventHookStatus(ctx, d, client, newHook.Status)
	if err != nil && utils.SuppressErrorOn404_V5(resp, err) == nil {
		// if we get a 404 when creating the hook, we need to taint the resource
		d.SetId("")
		return diag.Errorf("failed to set event hook status: %v", err)
	} else if err != nil {
		return diag.Errorf("failed to set event hook status: %v", err)
	}
	return resourceEventHookRead(ctx, d, meta)
}

func resourceEventHookRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := getOktaV5ClientFromMetadata(meta)
	hook, resp, err := client.EventHookAPI.GetEventHook(ctx, d.Id()).Execute()
	if err := utils.SuppressErrorOn404_V5(resp, err); err != nil {
		return diag.Errorf("failed to get event hook: %v", err)
	}
	if hook == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("name", hook.GetName())
	_ = d.Set("status", hook.GetStatus())
	events := hook.GetEvents()
	_ = d.Set("events", schema.NewSet(schema.HashString, utils.ConvertStringSliceToInterfaceSlice(events.Items)))

	err = utils.SetNonPrimitives(d, map[string]interface{}{
		"channel": flattenEventHookChannel(&hook.Channel),
		"headers": flattenEventHookHeaders(&hook.Channel),
		"auth":    flattenEventHookAuth(d, &hook.Channel),
		"filter":  flattenEventHookFilter(hook.Events.Filter),
	})
	if err != nil {
		return diag.Errorf("failed to set event hook properties: %v", err)
	}
	return nil
}

func resourceEventHookUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := getOktaV5ClientFromMetadata(meta)
	hook := buildEventHook(d)
	newHook, resp, err := client.EventHookAPI.ReplaceEventHook(ctx, d.Id()).EventHook(*hook).Execute()
	if err != nil && utils.SuppressErrorOn404_V5(resp, err) == nil {
		// if we get a 404 when updating the hook, we need to taint the resource
		d.SetId("")
		return diag.Errorf("failed to update event hook: %v", err)
	} else if err != nil {
		return diag.Errorf("failed to update event hook: %v", err)
	}
	resp, err = setEventHookStatus(ctx, d, client, newHook.Status)
	if err != nil && utils.SuppressErrorOn404_V5(resp, err) == nil {
		// if we get a 404 when updating the hook's status, we need to taint the resource
		d.SetId("")
		return diag.Errorf("failed to set event hook status: %v", err)
	} else if err != nil {
		return diag.Errorf("failed to set event hook status: %v", err)
	}
	return resourceEventHookRead(ctx, d, meta)
}

func resourceEventHookDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := getOktaV5ClientFromMetadata(meta)

	_, resp, err := client.EventHookAPI.DeactivateEventHook(ctx, d.Id()).Execute()
	if err != nil && utils.SuppressErrorOn404_V5(resp, err) == nil {
		// If we get a 404, we can assume the hook is already deleted
		return nil
	} else if err != nil {
		return diag.Errorf("failed to deactivate event hook: %v", err)
	}
	resp, err = client.EventHookAPI.DeleteEventHook(ctx, d.Id()).Execute()
	if err != nil && utils.SuppressErrorOn404_V5(resp, err) == nil {
		// If we get a 404, we can assume the hook is already deleted
		return nil
	} else if err != nil {
		return diag.Errorf("failed to delete event hook: %v", err)
	}
	return nil
}

func buildEventHook(d *schema.ResourceData) *v5okta.EventHook {
	// Build events
	eventSet := d.Get("events").(*schema.Set).List()
	events := make([]string, len(eventSet))
	for i, v := range eventSet {
		events[i] = v.(string)
	}

	eventSubscriptions := v5okta.NewEventSubscriptions(events, eventTypeSubscription)

	// Build filters if present
	var filterMapObjects []v5okta.EventHookFilterMapObject

	// Handle filter syntax
	if filterSet, ok := d.GetOk("filter"); ok {
		filters := filterSet.(*schema.Set).List()
		for _, filterItem := range filters {
			filterData := filterItem.(map[string]interface{})

			filterMapObject := v5okta.NewEventHookFilterMapObject()
			filterMapObject.SetEvent(filterData["event"].(string))

			// Note: filter type is readonly in the API and will be set automatically to "EXPRESSION_LANGUAGE"

			// Handle condition expression
			if conditionExpr, ok := filterData["condition"].(string); ok && conditionExpr != "" {
				condition := v5okta.NewEventHookFilterMapObjectCondition()
				condition.SetExpression(conditionExpr)
				// Note: Version field should be null as per API documentation
				// but we don't set it explicitly as the SDK handles null values
				filterMapObject.SetCondition(*condition)
			}

			filterMapObjects = append(filterMapObjects, *filterMapObject)
		}
	}

	// Set filters if any were found
	if len(filterMapObjects) > 0 {
		filter := v5okta.NewEventHookFilters()
		// Note: Type field is readonly after creation, but required during creation
		filter.SetType(filterTypeExpressionLanguage)
		filter.SetEventFilterMap(filterMapObjects)
		eventSubscriptions.SetFilter(*filter)
	}

	// Build channel
	channel := buildEventChannel(d)

	return v5okta.NewEventHook(*channel, *eventSubscriptions, d.Get("name").(string))
}

func buildEventChannel(d *schema.ResourceData) *v5okta.EventHookChannel {
	// Get channel config
	rawChannel := d.Get("channel").(map[string]interface{})
	uri, ok := rawChannel["uri"].(string)
	if !ok || uri == "" {
		// This should not happen with proper user input, but be defensive
		// Note: The API will validate URI requirements
		uri = ""
	}
	config := v5okta.NewEventHookChannelConfig(uri)

	// Add headers if present
	if raw, ok := d.GetOk("headers"); ok {
		var headerList []v5okta.EventHookChannelConfigHeader
		for _, header := range raw.(*schema.Set).List() {
			h := header.(map[string]interface{})
			key, keyOk := h["key"].(string)
			value, valueOk := h["value"].(string)
			if keyOk && valueOk && key != "" && value != "" {
				headerObj := v5okta.NewEventHookChannelConfigHeader()
				headerObj.SetKey(key)
				headerObj.SetValue(value)
				headerList = append(headerList, *headerObj)
			}
		}
		config.SetHeaders(headerList)
	}

	// Add auth if present
	if rawAuth, ok := d.GetOk("auth"); ok {
		a := rawAuth.(map[string]interface{})

		// Validate required auth fields
		key, keyOk := a["key"].(string)
		value, valueOk := a["value"].(string)
		if !keyOk || !valueOk || key == "" || value == "" {
			// Skip auth if key or value is missing/empty
			// This should be caught by schema validation, but be defensive
		} else {
			authScheme := v5okta.NewEventHookChannelConfigAuthScheme()

			// Apply default auth type (works with DiffSuppressFunc in schema)
			authType := defaultAuthType
			if t, ok := a["type"].(string); ok && t != "" {
				authType = t
			}
			authScheme.SetType(authType)
			authScheme.SetKey(key)
			authScheme.SetValue(value)
			config.SetAuthScheme(*authScheme)
		}
	}

	// Get channel type and version (apply defaults for optional fields)
	// These defaults work in conjunction with DiffSuppressFunc in the schema
	channelType := defaultChannelType
	if t, ok := rawChannel["type"].(string); ok && t != "" {
		channelType = t
	}

	version := defaultChannelVersion
	if v, ok := rawChannel["version"].(string); ok && v != "" {
		version = v
	}

	return v5okta.NewEventHookChannel(*config, channelType, version)
}

// EventSet converts an EventSubscriptions object to a Terraform Set for the events field
func EventSet(e *v5okta.EventSubscriptions) *schema.Set {
	events := make([]interface{}, 0, len(e.Items))
	for _, event := range e.Items {
		events = append(events, event)
	}
	return schema.NewSet(schema.HashString, events)
}

func flattenEventHookChannel(c *v5okta.EventHookChannel) map[string]interface{} {
	channelType := c.GetType()
	if channelType == "" {
		channelType = defaultChannelType
	}

	version := c.GetVersion()
	if version == "" {
		version = defaultChannelVersion
	}

	return map[string]interface{}{
		"type":    channelType,
		"version": version,
		"uri":     c.Config.GetUri(),
	}
}

func flattenEventHookHeaders(c *v5okta.EventHookChannel) *schema.Set {
	var headers []interface{}
	if c.Config.HasHeaders() {
		headers = make([]interface{}, len(c.Config.GetHeaders()))
		for i, header := range c.Config.GetHeaders() {
			headers[i] = map[string]interface{}{
				"key":   header.GetKey(),
				"value": header.GetValue(),
			}
		}
	}
	return schema.NewSet(schema.HashResource(eventHookHeaderSchema), headers)
}

func flattenEventHookAuth(d *schema.ResourceData, c *v5okta.EventHookChannel) map[string]interface{} {
	auth := map[string]interface{}{}
	if c.Config.HasAuthScheme() {
		authScheme := c.Config.GetAuthScheme()

		// Get the stored value from terraform state (sensitive value not returned by API)
		var storedValue string
		if authMap, ok := d.Get("auth").(map[string]interface{}); ok {
			if val, exists := authMap["value"]; exists {
				storedValue = val.(string)
			}
		}

		authType := authScheme.GetType()
		if authType == "" {
			authType = defaultAuthType
		}

		auth = map[string]interface{}{
			"key":   authScheme.GetKey(),
			"type":  authType,
			"value": storedValue,
		}
	}
	return auth
}

func flattenEventHookFilter(filter v5okta.NullableEventHookFilters) *schema.Set {
	// Define schema inline since it's only used here - match the resource schema
	filterSchema := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type":      {Type: schema.TypeString},
			"event":     {Type: schema.TypeString},
			"condition": {Type: schema.TypeString},
		},
	}

	if !filter.IsSet() || filter.Get() == nil {
		return schema.NewSet(schema.HashResource(filterSchema), []interface{}{})
	}

	filterObj := filter.Get()
	var filters []interface{}

	if filterObj.HasEventFilterMap() {
		// Get the filter type to include in each filter item
		filterType := filterTypeExpressionLanguage // default
		if filterObj.HasType() {
			filterType = filterObj.GetType()
		}

		filters = make([]interface{}, len(filterObj.GetEventFilterMap()))
		for i, filterMapObject := range filterObj.GetEventFilterMap() {
			filterData := map[string]interface{}{
				"type": filterType, // Include type in each filter item
			}

			if filterMapObject.HasEvent() {
				filterData["event"] = filterMapObject.GetEvent()
			}

			if filterMapObject.HasCondition() {
				condition := filterMapObject.GetCondition()
				filterData["condition"] = condition.GetExpression()
			}

			filters[i] = filterData
		}
	}

	return schema.NewSet(schema.HashResource(filterSchema), filters)
}

func setEventHookStatus(ctx context.Context, d *schema.ResourceData, client *v5okta.APIClient, status *string) (*v5okta.APIResponse, error) {
	desiredStatus := d.Get("status").(string)
	currentStatus := ""
	if status != nil {
		currentStatus = *status
	}

	if currentStatus == desiredStatus {
		return nil, nil
	}

	var resp *v5okta.APIResponse
	var err error
	if desiredStatus == StatusInactive {
		_, resp, err = client.EventHookAPI.DeactivateEventHook(ctx, d.Id()).Execute()
		if err != nil {
			return resp, fmt.Errorf("failed to deactivate event hook (current: %s, desired: %s): %w", currentStatus, desiredStatus, err)
		}
	} else {
		_, resp, err = client.EventHookAPI.ActivateEventHook(ctx, d.Id()).Execute()
		if err != nil {
			return resp, fmt.Errorf("failed to activate event hook (current: %s, desired: %s): %w", currentStatus, desiredStatus, err)
		}
	}
	return resp, err
}

package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
)

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
			"headers": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        headerSchema,
				Description: "Map of headers to send along in event hook request.",
			},
			"auth": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if k == "auth.type" && new == "" {
						return true
					}
					return false
				},
				Description: "Authentication required for event hook request.",
			},
			"channel": {
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if k == "channel.type" && new == "" {
						return true
					}
					return false
				},
				Description: `Details of the endpoint the event hook will hit.   
	- 'version' - (Required) The version of the channel. The currently-supported version is '1.0.0'.
	- 'uri' - (Required) The URI the hook will hit.
	- 'type' - (Optional) The type of hook to trigger. Currently, the only supported type is 'HTTP'.`,
			},
		},
	}
}

func resourceEventHookCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(m)
	hook := buildEventHook(d)
	newHook, _, err := client.EventHook.CreateEventHook(ctx, *hook)
	if err != nil {
		return diag.Errorf("failed to create event hook: %v", err)
	}
	d.SetId(newHook.Id)
	err = setEventHookStatus(ctx, d, client, newHook.Status)
	if err != nil {
		return diag.Errorf("failed to set event hook status: %v", err)
	}
	return resourceEventHookRead(ctx, d, m)
}

func resourceEventHookRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hook, resp, err := getOktaClientFromMetadata(m).EventHook.GetEventHook(ctx, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get event hook: %v", err)
	}
	if hook == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("name", hook.Name)
	_ = d.Set("status", hook.Status)
	_ = d.Set("events", eventSet(hook.Events))
	err = setNonPrimitives(d, map[string]interface{}{
		"channel": flattenEventHookChannel(hook.Channel),
		"headers": flattenEventHookHeaders(hook.Channel),
		"auth":    flattenEventHookAuth(d, hook.Channel),
	})
	if err != nil {
		return diag.Errorf("failed to set event hook properties: %v", err)
	}
	return nil
}

func resourceEventHookUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(m)
	hook := buildEventHook(d)
	newHook, _, err := client.EventHook.UpdateEventHook(ctx, d.Id(), *hook)
	if err != nil {
		return diag.Errorf("failed to update auth event hook: %v", err)
	}
	err = setEventHookStatus(ctx, d, client, newHook.Status)
	if err != nil {
		return diag.Errorf("failed to set event hook status: %v", err)
	}
	return resourceEventHookRead(ctx, d, m)
}

func resourceEventHookDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(m)

	_, _, err := client.EventHook.DeactivateEventHook(ctx, d.Id())
	if err != nil {
		return diag.Errorf("failed to deactivate event hook: %v", err)
	}
	_, err = client.EventHook.DeleteEventHook(ctx, d.Id())
	if err != nil {
		return diag.Errorf("failed to delete event hook: %v", err)
	}
	return nil
}

func buildEventHook(d *schema.ResourceData) *sdk.EventHook {
	eventSet := d.Get("events").(*schema.Set).List()
	events := make([]string, len(eventSet))
	for i, v := range eventSet {
		events[i] = v.(string)
	}
	return &sdk.EventHook{
		Name:    d.Get("name").(string),
		Status:  d.Get("status").(string),
		Events:  &sdk.EventSubscriptions{Type: "EVENT_TYPE", Items: events},
		Channel: buildEventChannel(d),
	}
}

func buildEventChannel(d *schema.ResourceData) *sdk.EventHookChannel {
	var headerList []*sdk.EventHookChannelConfigHeader
	if raw, ok := d.GetOk("headers"); ok {
		for _, header := range raw.(*schema.Set).List() {
			h, ok := header.(map[string]interface{})
			if ok {
				headerList = append(headerList, &sdk.EventHookChannelConfigHeader{Key: h["key"].(string), Value: h["value"].(string)})
			}
		}
	}
	var auth *sdk.EventHookChannelConfigAuthScheme
	if rawAuth, ok := d.GetOk("auth"); ok {
		a := rawAuth.(map[string]interface{})
		_, ok := a["type"]
		if !ok {
			a["type"] = "HEADER"
		}
		auth = &sdk.EventHookChannelConfigAuthScheme{
			Key:   a["key"].(string),
			Type:  a["type"].(string),
			Value: a["value"].(string),
		}
	}
	rawChannel := d.Get("channel").(map[string]interface{})
	_, ok := rawChannel["type"]
	if !ok {
		rawChannel["type"] = "HTTP"
	}
	return &sdk.EventHookChannel{
		Config: &sdk.EventHookChannelConfig{
			Uri:        rawChannel["uri"].(string),
			AuthScheme: auth,
			Headers:    headerList,
		},
		Type:    rawChannel["type"].(string),
		Version: rawChannel["version"].(string),
	}
}

func flattenEventHookAuth(d *schema.ResourceData, c *sdk.EventHookChannel) map[string]interface{} {
	auth := map[string]interface{}{}
	if c.Config.AuthScheme != nil {
		auth = map[string]interface{}{
			"key":   c.Config.AuthScheme.Key,
			"type":  c.Config.AuthScheme.Type,
			"value": d.Get("auth").(map[string]interface{})["value"],
		}
	}
	return auth
}

func flattenEventHookChannel(c *sdk.EventHookChannel) map[string]interface{} {
	return map[string]interface{}{
		"type":    c.Type,
		"version": c.Version,
		"uri":     c.Config.Uri,
	}
}

func flattenEventHookHeaders(c *sdk.EventHookChannel) *schema.Set {
	headers := make([]interface{}, len(c.Config.Headers))
	for i, header := range c.Config.Headers {
		headers[i] = map[string]interface{}{
			"key":   header.Key,
			"value": header.Value,
		}
	}
	return schema.NewSet(schema.HashResource(headerSchema), headers)
}

func eventSet(e *sdk.EventSubscriptions) *schema.Set {
	events := make([]interface{}, len(e.Items))
	for i, event := range e.Items {
		events[i] = event
	}
	return schema.NewSet(schema.HashString, events)
}

func setEventHookStatus(ctx context.Context, d *schema.ResourceData, client *sdk.Client, status string) error {
	desiredStatus := d.Get("status").(string)
	if status == desiredStatus {
		return nil
	}
	var err error
	if desiredStatus == statusInactive {
		_, _, err = client.EventHook.DeactivateEventHook(ctx, d.Id())
	} else {
		_, _, err = client.EventHook.ActivateEventHook(ctx, d.Id())
	}
	return err
}

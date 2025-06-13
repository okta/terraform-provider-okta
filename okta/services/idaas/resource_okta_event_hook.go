package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v5/okta"
	"github.com/okta/terraform-provider-okta/okta/utils"
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
					if k == "auth.type" && new == "" {
						return true
					}
					return false
				},
				Description: `Details of the endpoint the event hook will hit.   
	- 'version' - (Required) The version of the channel. The currently-supported version is '1.0.0'.
	- 'uri' - (Required) The URI the hook will hit.
	- 'type' - (Optional) The type of hook to trigger. Currently, the only supported type is 'HTTP'.`,
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
				Description: "Details of the endpoint the event hook will hit.",
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
	d.SetId(newHook.GetId())
	err = setEventHookStatus(ctx, d, client, newHook.GetStatus())
	if err != nil {
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

	channel := hook.GetChannel()
	config := channel.GetConfig()
	channelMap := map[string]interface{}{
		"type":    channel.GetType(),
		"version": channel.GetVersion(),
		"uri":     config.GetUri(),
	}
	_ = d.Set("channel", channelMap)

	if channel.GetConfig().AuthScheme != nil {
		auth := map[string]interface{}{
			"key":   channel.GetConfig().AuthScheme.GetKey(),
			"type":  channel.GetConfig().AuthScheme.GetType(),
			"value": d.Get("auth").(map[string]interface{})["value"],
		}
		_ = d.Set("auth", auth)
	}

	headers := make([]interface{}, len(channel.GetConfig().Headers))
	for i, header := range channel.GetConfig().Headers {
		headers[i] = map[string]interface{}{
			"key":   header.GetKey(),
			"value": header.GetValue(),
		}
	}
	_ = d.Set("headers", schema.NewSet(schema.HashResource(HeaderSchema), headers))

	return nil
}

func resourceEventHookUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := getOktaV5ClientFromMetadata(meta)
	hook := buildEventHook(d)
	newHook, _, err := client.EventHookAPI.ReplaceEventHook(ctx, d.Id()).EventHook(*hook).Execute()
	if err != nil {
		return diag.Errorf("failed to update auth event hook: %v", err)
	}
	err = setEventHookStatus(ctx, d, client, newHook.GetStatus())
	if err != nil {
		return diag.Errorf("failed to set event hook status: %v", err)
	}
	return resourceEventHookRead(ctx, d, meta)
}

func resourceEventHookDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := getOktaV5ClientFromMetadata(meta)

	_, _, err := client.EventHookAPI.DeactivateEventHook(ctx, d.Id()).Execute()
	if err != nil {
		return diag.Errorf("failed to deactivate event hook: %v", err)
	}
	_, err = client.EventHookAPI.DeleteEventHook(ctx, d.Id()).Execute()
	if err != nil {
		return diag.Errorf("failed to delete event hook: %v", err)
	}
	return nil
}

func buildEventHook(d *schema.ResourceData) *okta.EventHook {
	eventSet := d.Get("events").(*schema.Set).List()
	events := make([]string, len(eventSet))
	for i, v := range eventSet {
		events[i] = v.(string)
	}

	rawChannel := d.Get("channel").(map[string]interface{})
	if _, ok := rawChannel["type"]; !ok {
		rawChannel["type"] = "HTTP"
	}

	var headerList []okta.EventHookChannelConfigHeader
	if raw, ok := d.GetOk("headers"); ok {
		for _, header := range raw.(*schema.Set).List() {
			h, ok := header.(map[string]interface{})
			if ok {
				headerObj := okta.EventHookChannelConfigHeader{}
				headerObj.SetKey(h["key"].(string))
				headerObj.SetValue(h["value"].(string))
				headerList = append(headerList, headerObj)
			}
		}
	}

	var auth *okta.EventHookChannelConfigAuthScheme
	if rawAuth, ok := d.GetOk("auth"); ok {
		a := rawAuth.(map[string]interface{})
		authType := "HEADER"
		if t, ok := a["type"]; ok {
			authType = t.(string)
		}
		auth = &okta.EventHookChannelConfigAuthScheme{}
		auth.SetKey(a["key"].(string))
		auth.SetType(authType)
		auth.SetValue(a["value"].(string))
	}

	config := &okta.EventHookChannelConfig{}
	config.SetUri(rawChannel["uri"].(string))
	if auth != nil {
		config.SetAuthScheme(*auth)
	}
	if len(headerList) > 0 {
		config.SetHeaders(headerList)
	}

	channel := &okta.EventHookChannel{}
	channel.SetConfig(*config)
	channel.SetType(rawChannel["type"].(string))
	channel.SetVersion(rawChannel["version"].(string))

	eventSubs := &okta.EventSubscriptions{}
	eventSubs.SetItems(events)
	eventSubs.SetType("EVENT_TYPE")

	hook := &okta.EventHook{}
	hook.SetChannel(*channel)
	hook.SetEvents(*eventSubs)
	hook.SetName(d.Get("name").(string))
	if status, ok := d.GetOk("status"); ok {
		hook.SetStatus(status.(string))
	}

	return hook
}

func setEventHookStatus(ctx context.Context, d *schema.ResourceData, client *okta.APIClient, status string) error {
	desiredStatus := d.Get("status").(string)
	if status == desiredStatus {
		return nil
	}
	var err error
	if desiredStatus == StatusInactive {
		_, _, err = client.EventHookAPI.DeactivateEventHook(ctx, d.Id()).Execute()
	} else {
		_, _, err = client.EventHookAPI.ActivateEventHook(ctx, d.Id()).Execute()
	}
	return err
}

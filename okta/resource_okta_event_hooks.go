package okta

import (
	"context"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
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
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": statusSchema,
			"events": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"headers": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     headerSchema,
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
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					var errs diag.Diagnostics
					m := i.(map[string]interface{})
					if _, ok := m["key"]; !ok {
						errs = append(errs, diag.Errorf("auth 'key' should not be empty")...)
					}
					if t, ok := m["type"]; ok {
						dErr := stringInSlice([]string{"HEADER"})(t, cty.GetAttrPath("type"))
						if dErr != nil {
							errs = append(errs, dErr...)
						}
					} else {
						m["type"] = "HEADER"
					}
					if _, ok := m["value"]; !ok {
						errs = append(errs, diag.Errorf("auth 'value' should not be empty")...)
					}
					return errs
				},
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
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					var errs diag.Diagnostics
					m := i.(map[string]interface{})
					if t, ok := m["type"]; ok {
						dErr := stringInSlice([]string{"HTTP"})(t, cty.GetAttrPath("type"))
						if dErr != nil {
							errs = append(errs, dErr...)
						}
					}
					dErr := stringIsVersion(m["version"], cty.GetAttrPath("version"))
					if dErr != nil {
						errs = append(errs, dErr...)
					}
					dErr = stringIsURL("https")(m["uri"], cty.GetAttrPath("uri"))
					if dErr != nil {
						errs = append(errs, dErr...)
					}
					return errs
				},
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

func buildEventHook(d *schema.ResourceData) *okta.EventHook {
	eventSet := d.Get("events").(*schema.Set).List()
	events := make([]string, len(eventSet))
	for i, v := range eventSet {
		events[i] = v.(string)
	}
	return &okta.EventHook{
		Name:    d.Get("name").(string),
		Status:  d.Get("status").(string),
		Events:  &okta.EventSubscriptions{Type: "EVENT_TYPE", Items: events},
		Channel: buildEventChannel(d),
	}
}

func buildEventChannel(d *schema.ResourceData) *okta.EventHookChannel {
	var headerList []*okta.EventHookChannelConfigHeader
	if raw, ok := d.GetOk("headers"); ok {
		for _, header := range raw.(*schema.Set).List() {
			h, ok := header.(map[string]interface{})
			if ok {
				headerList = append(headerList, &okta.EventHookChannelConfigHeader{Key: h["key"].(string), Value: h["value"].(string)})
			}
		}
	}
	var auth *okta.EventHookChannelConfigAuthScheme
	if rawAuth, ok := d.GetOk("auth"); ok {
		a := rawAuth.(map[string]interface{})
		_, ok := a["type"]
		if !ok {
			a["type"] = "HEADER"
		}
		auth = &okta.EventHookChannelConfigAuthScheme{
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
	return &okta.EventHookChannel{
		Config: &okta.EventHookChannelConfig{
			Uri:        rawChannel["uri"].(string),
			AuthScheme: auth,
			Headers:    headerList,
		},
		Type:    rawChannel["type"].(string),
		Version: rawChannel["version"].(string),
	}
}

func flattenEventHookAuth(d *schema.ResourceData, c *okta.EventHookChannel) map[string]interface{} {
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

func flattenEventHookChannel(c *okta.EventHookChannel) map[string]interface{} {
	return map[string]interface{}{
		"type":    c.Type,
		"version": c.Version,
		"uri":     c.Config.Uri,
	}
}

func flattenEventHookHeaders(c *okta.EventHookChannel) *schema.Set {
	headers := make([]interface{}, len(c.Config.Headers))
	for i, header := range c.Config.Headers {
		headers[i] = map[string]interface{}{
			"key":   header.Key,
			"value": header.Value,
		}
	}
	return schema.NewSet(schema.HashResource(headerSchema), headers)
}

func eventSet(e *okta.EventSubscriptions) *schema.Set {
	events := make([]interface{}, len(e.Items))
	for i, event := range e.Items {
		events[i] = event
	}
	return schema.NewSet(schema.HashString, events)
}

func setEventHookStatus(ctx context.Context, d *schema.ResourceData, client *okta.Client, status string) error {
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

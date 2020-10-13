package okta

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"

	"net/http"
)

func resourceEventHook() *schema.Resource {
	return &schema.Resource{
		Create: resourceEventHookCreate,
		Read:   resourceEventHookRead,
		Update: resourceEventHookUpdate,
		Delete: resourceEventHookDelete,
		Exists: resourceEventHookExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							Default:      "HEADER",
							ValidateFunc: validation.StringInSlice([]string{"HEADER"}, false),
						},
						"value": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
					},
				},
			},
			"channel": {
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							Default:  "HTTP",
						},
						"version": {
							Type:     schema.TypeString,
							Required: true,
							Default:  "1.0.0",
						},
						"uri": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceEventHookCreate(d *schema.ResourceData, m interface{}) error {
	client := getSupplementFromMetadata(m)
	hook := buildEventHook(d, m)
	newHook, _, err := client.CreateEventHook(*hook, nil)
	if err != nil {
		return err
	}

	d.SetId(newHook.ID)
	desiredStatus := d.Get("status").(string)
	err = setEventHookStatus(d, client, newHook.Status, desiredStatus)
	if err != nil {
		return err
	}

	return resourceEventHookRead(d, m)
}

func resourceEventHookExists(d *schema.ResourceData, m interface{}) (bool, error) {
	_, res, err := getSupplementFromMetadata(m).GetEventHook(d.Id())
	return err == nil && res.StatusCode != http.StatusNotFound, err
}

func resourceEventHookRead(d *schema.ResourceData, m interface{}) error {
	hook, resp, err := getSupplementFromMetadata(m).GetEventHook(d.Id())

	if is404(resp.StatusCode) {
		d.SetId("")
		return nil
	}

	if err != nil {
		return err
	}

	_ = d.Set("name", hook.Name)
	_ = d.Set("status", hook.Status)
	_ = d.Set("events", eventSet(hook.Events))

	return setNonPrimitives(d, map[string]interface{}{
		"channel": flattenEventHookChannel(hook.Channel),
		"headers": flattenEventHookHeaders(hook.Channel),
		"auth":    flattenEventHookAuth(d, hook.Channel),
	})
}

func resourceEventHookUpdate(d *schema.ResourceData, m interface{}) error {
	client := getSupplementFromMetadata(m)
	hook := buildEventHook(d, m)
	newHook, _, err := client.UpdateEventHook(d.Id(), *hook, nil)

	if err != nil {
		return err
	}

	desiredStatus := d.Get("status").(string)
	err = setEventHookStatus(d, client, newHook.Status, desiredStatus)
	if err != nil {
		return err
	}

	return resourceEventHookRead(d, m)
}

func resourceEventHookDelete(d *schema.ResourceData, m interface{}) error {
	client := getSupplementFromMetadata(m)
	res, err := client.DeactivateEventHook(d.Id())
	if err != nil {
		return responseErr(res, err)
	}

	_, err = client.DeleteEventHook(d.Id())

	return err
}

func buildEventHook(d *schema.ResourceData, m interface{}) *sdk.EventHook {
	eventSet := d.Get("events").(*schema.Set).List()
	events := make([]string, len(eventSet))
	for i, v := range eventSet {
		events[i] = v.(string)
	}
	return &sdk.EventHook{
		Name:    d.Get("name").(string),
		Status:  d.Get("status").(string),
		Events:  &sdk.EventHookEvents{Type: "EVENT_TYPE", Items: events},
		Channel: buildEventChannel(d, m),
	}
}

func buildEventChannel(d *schema.ResourceData, m interface{}) *sdk.EventHookChannel {
	if _, ok := d.GetOk("channel"); !ok {
		return nil
	}

	var headerList []*sdk.EventHookHeader
	if raw, ok := d.GetOk("headers"); ok {
		for _, header := range raw.(*schema.Set).List() {
			h, ok := header.(map[string]interface{})
			if ok {
				headerList = append(headerList, &sdk.EventHookHeader{Key: h["key"].(string), Value: h["value"].(string)})
			}
		}
	}

	var auth *sdk.EventHookAuthScheme
	if _, ok := d.GetOk("auth.key"); ok {
		auth = &sdk.EventHookAuthScheme{
			Key:   getStringValue(d, "auth.key"),
			Type:  getStringValue(d, "auth.type"),
			Value: getStringValue(d, "auth.value"),
		}
	}

	return &sdk.EventHookChannel{
		Config: &sdk.EventHookChannelConfig{
			URI:        getStringValue(d, "channel.uri"),
			AuthScheme: auth,
			Headers:    headerList,
		},
		Type:    getStringValue(d, "channel.type"),
		Version: getStringValue(d, "channel.version"),
	}
}

func flattenEventHookAuth(d *schema.ResourceData, c *sdk.EventHookChannel) map[string]interface{} {
	auth := map[string]interface{}{}

	if c.Config.AuthScheme != nil {
		auth = map[string]interface{}{
			"key":  c.Config.AuthScheme.Key,
			"type": c.Config.AuthScheme.Type,
			// Read only
			"value": getStringValue(d, "auth.value"),
		}
	}
	return auth
}

func flattenEventHookChannel(c *sdk.EventHookChannel) map[string]interface{} {
	return map[string]interface{}{
		"type":    c.Type,
		"version": c.Version,
		"uri":     c.Config.URI,
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

func eventSet(e *sdk.EventHookEvents) *schema.Set {
	events := make([]interface{}, len(e.Items))
	for i, event := range e.Items {
		events[i] = event
	}

	return schema.NewSet(schema.HashString, events)
}

func setEventHookStatus(d *schema.ResourceData, client *sdk.ApiSupplement, status string, desiredStatus string) error {
	if status != desiredStatus {
		if desiredStatus == "INACTIVE" {
			return responseErr(client.DeactivateEventHook(d.Id()))
		} else if desiredStatus == "ACTIVE" {
			return responseErr(client.ActivateEventHook(d.Id()))
		}
	}

	return nil
}

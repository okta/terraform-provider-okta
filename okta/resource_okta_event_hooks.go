package okta

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-okta/sdk"

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
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"status": statusSchema,
			"events": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"headers": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     headerSchema,
			},
			"auth": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"type": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							Default:      "HEADER",
							ValidateFunc: validation.StringInSlice([]string{"HEADER"}, false),
						},
						"value": &schema.Schema{
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
					},
				},
			},
			"channel": &schema.Schema{
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							Default:  "HTTP",
						},
						"version": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							Default:  "1.0.0",
						},
						"uri": &schema.Schema{
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

	d.Set("name", hook.Name)
	d.Set("status", hook.Status)
	d.Set("events", eventSet(hook.Events))

	return setNonPrimitives(d, map[string]interface{}{
		"channel": flattenEventHookChannel(hook.Channel),
		"headers": flattenHeaders(hook.Channel),
		"auth":    flattenAuth(d, hook.Channel),
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
		Events:  &sdk.Events{Items: events},
		Channel: buildEventChannel(d, m),
	}
}

func buildEventChannel(d *schema.ResourceData, m interface{}) *sdk.Channel {
	if _, ok := d.GetOk("channel"); !ok {
		return nil
	}

	headerList := []*sdk.Header{}
	if raw, ok := d.GetOk("headers"); ok {
		for _, header := range raw.(*schema.Set).List() {
			h, ok := header.(map[string]interface{})
			if ok {
				headerList = append(headerList, &sdk.Header{Key: h["key"].(string), Value: h["value"].(string)})
			}
		}
	}

	var auth *sdk.AuthScheme
	if _, ok := d.GetOk("auth.key"); ok {
		auth = &sdk.AuthScheme{
			Key:   getStringValue(d, "auth.key"),
			Type:  getStringValue(d, "auth.type"),
			Value: getStringValue(d, "auth.value"),
		}
	}

	return &sdk.Channel{
		Config: &sdk.HookConfig{
			URI:        getStringValue(d, "channel.uri"),
			AuthScheme: auth,
			Headers:    headerList,
		},
		Type:    getStringValue(d, "channel.type"),
		Version: getStringValue(d, "channel.version"),
	}
}

func flattenEventHookChannel(c *sdk.Channel) map[string]interface{} {
	return map[string]interface{}{
		"type":    c.Type,
		"version": c.Version,
		"uri":     c.Config.URI,
	}
}

func eventSet(e *sdk.Events) *schema.Set {
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

package okta

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-okta/sdk"

	"net/http"
)

var headerSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"key": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"value": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
	},
}

func resourceInlineHook() *schema.Resource {
	return &schema.Resource{
		Create: resourceInlineHookCreate,
		Read:   resourceInlineHookRead,
		Update: resourceInlineHookUpdate,
		Delete: resourceInlineHookDelete,
		Exists: resourceInlineHookExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		// For those familiar with Terraform schemas be sure to check the base hook schema and/or
		// the examples in the documentation
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"status": statusSchema,
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice(
					[]string{
						"com.okta.oauth2.tokens.transform",
						"com.okta.import.transform",
						"com.okta.saml.tokens.transform",
						"com.okta.user.pre-registration",
					},
					false,
				),
			},
			"version": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
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
							Optional: true,
						},
						"type": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "HEADER",
							ValidateFunc: validation.StringInSlice([]string{"HEADER"}, false),
						},
						"value": &schema.Schema{
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
					},
				},
			},
			"channel": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:     schema.TypeBool,
							Default:  "HTTP",
							Optional: true,
						},
						"version": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"uri": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"method": &schema.Schema{
							Type:     schema.TypeString,
							Default:  "POST",
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceInlineHookCreate(d *schema.ResourceData, m interface{}) error {
	client := getSupplementFromMetadata(m)
	hook := buildInlineHook(d, m)
	newHook, _, err := client.CreateInlineHook(*hook, nil)
	if err != nil {
		return err
	}

	d.SetId(newHook.ID)
	desiredStatus := d.Get("status").(string)
	err = setHookStatus(d, client, newHook.Status, desiredStatus)
	if err != nil {
		return err
	}

	return resourceInlineHookRead(d, m)
}

func resourceInlineHookExists(d *schema.ResourceData, m interface{}) (bool, error) {
	_, res, err := getSupplementFromMetadata(m).GetInlineHook(d.Id())
	return err == nil && res.StatusCode != http.StatusNotFound, err
}

func resourceInlineHookRead(d *schema.ResourceData, m interface{}) error {
	hook, resp, err := getSupplementFromMetadata(m).GetInlineHook(d.Id())

	if is404(resp.StatusCode) {
		d.SetId("")
		return nil
	}

	if err != nil {
		return err
	}

	d.Set("name", hook.Name)
	d.Set("status", hook.Status)
	d.Set("type", hook.Type)
	d.Set("version", hook.Version)

	return setNonPrimitives(d, map[string]interface{}{
		"channel": flattenHookChannel(hook.Channel),
		"headers": flattenHeaders(hook.Channel),
		"auth":    flattenAuth(d, hook.Channel),
	})
}

func resourceInlineHookUpdate(d *schema.ResourceData, m interface{}) error {
	client := getSupplementFromMetadata(m)
	hook := buildInlineHook(d, m)
	newHook, _, err := client.UpdateInlineHook(d.Id(), *hook, nil)

	if err != nil {
		return err
	}

	desiredStatus := d.Get("status").(string)
	err = setHookStatus(d, client, newHook.Status, desiredStatus)
	if err != nil {
		return err
	}

	return resourceInlineHookRead(d, m)
}

func resourceInlineHookDelete(d *schema.ResourceData, m interface{}) error {
	client := getSupplementFromMetadata(m)
	res, err := client.DeactivateInlineHook(d.Id())
	if err != nil {
		return responseErr(res, err)
	}

	_, err = client.DeleteInlineHook(d.Id())

	return err
}

func buildInlineHook(d *schema.ResourceData, m interface{}) *sdk.InlineHook {
	return &sdk.InlineHook{
		Name:    d.Get("name").(string),
		Status:  d.Get("status").(string),
		Type:    d.Get("type").(string),
		Version: d.Get("version").(string),
		Channel: buildInlineChannel(d, m),
	}
}

func buildInlineChannel(d *schema.ResourceData, m interface{}) *sdk.Channel {
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
			Method:     getStringValue(d, "channel.method"),
		},
		Type:    getStringValue(d, "channel.type"),
		Version: getStringValue(d, "channel.version"),
	}
}

func flattenAuth(d *schema.ResourceData, c *sdk.Channel) map[string]interface{} {
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

func flattenHookChannel(c *sdk.Channel) map[string]interface{} {
	return map[string]interface{}{
		"type":    c.Type,
		"version": c.Version,
		"uri":     c.Config.URI,
		"method":  c.Config.Method,
	}
}

func flattenHeaders(c *sdk.Channel) *schema.Set {
	headers := make([]interface{}, len(c.Config.Headers))
	for i, header := range c.Config.Headers {
		headers[i] = map[string]interface{}{
			"key":   header.Key,
			"value": header.Value,
		}
	}

	return schema.NewSet(schema.HashResource(headerSchema), headers)
}

func setHookStatus(d *schema.ResourceData, client *sdk.ApiSupplement, status string, desiredStatus string) error {
	if status != desiredStatus {
		if desiredStatus == "INACTIVE" {
			return responseErr(client.DeactivateInlineHook(d.Id()))
		} else if desiredStatus == "ACTIVE" {
			return responseErr(client.ActivateInlineHook(d.Id()))
		}
	}

	return nil
}

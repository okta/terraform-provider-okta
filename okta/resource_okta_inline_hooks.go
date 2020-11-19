package okta

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"

	"net/http"
)

var headerSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"key": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"value": {
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
			StateContext: schema.ImportStatePassthroughContext,
		},

		// For those familiar with Terraform schemas be sure to check the base hook schema and/or
		// the examples in the documentation
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": statusSchema,
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice(
					[]string{
						"com.okta.oauth2.tokens.transform",
						"com.okta.import.transform",
						"com.okta.saml.tokens.transform",
						"com.okta.user.pre-registration",
						"com.okta.user.credential.password.import",
					},
					false,
				),
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
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
							Optional: true,
						},
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "HEADER",
							ValidateFunc: validation.StringInSlice([]string{"HEADER"}, false),
						},
						"value": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
					},
				},
			},
			"channel": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeBool,
							Default:  "HTTP",
							Optional: true,
						},
						"version": {
							Type:     schema.TypeString,
							Required: true,
						},
						"uri": {
							Type:     schema.TypeString,
							Required: true,
						},
						"method": {
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
	hook := buildInlineHook(d)
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

	if resp != nil && is404(resp.StatusCode) {
		d.SetId("")
		return nil
	}

	if err != nil {
		return err
	}

	_ = d.Set("name", hook.Name)
	_ = d.Set("status", hook.Status)
	_ = d.Set("type", hook.Type)
	_ = d.Set("version", hook.Version)

	return setNonPrimitives(d, map[string]interface{}{
		"channel": flattenInlineHookChannel(hook.Channel),
		"headers": flattenInlineHookHeaders(hook.Channel),
		"auth":    flattenInlineHookAuth(d, hook.Channel),
	})
}

func resourceInlineHookUpdate(d *schema.ResourceData, m interface{}) error {
	client := getSupplementFromMetadata(m)
	hook := buildInlineHook(d)
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

func buildInlineHook(d *schema.ResourceData) *sdk.InlineHook {
	return &sdk.InlineHook{
		Name:    d.Get("name").(string),
		Status:  d.Get("status").(string),
		Type:    d.Get("type").(string),
		Version: d.Get("version").(string),
		Channel: buildInlineChannel(d),
	}
}

func buildInlineChannel(d *schema.ResourceData) *sdk.InlineHookChannel {
	if _, ok := d.GetOk("channel"); !ok {
		return nil
	}

	var headerList []*sdk.InlineHookHeader
	if raw, ok := d.GetOk("headers"); ok {
		for _, header := range raw.(*schema.Set).List() {
			h, ok := header.(map[string]interface{})
			if ok {
				headerList = append(headerList, &sdk.InlineHookHeader{Key: h["key"].(string), Value: h["value"].(string)})
			}
		}
	}

	var auth *sdk.InlineHookAuthScheme
	if _, ok := d.GetOk("auth.key"); ok {
		auth = &sdk.InlineHookAuthScheme{
			Key:   getStringValue(d, "auth.key"),
			Type:  getStringValue(d, "auth.type"),
			Value: getStringValue(d, "auth.value"),
		}
	}

	return &sdk.InlineHookChannel{
		Config: &sdk.InlineHookChannelConfig{
			URI:        getStringValue(d, "channel.uri"),
			AuthScheme: auth,
			Headers:    headerList,
			Method:     getStringValue(d, "channel.method"),
		},
		Type:    getStringValue(d, "channel.type"),
		Version: getStringValue(d, "channel.version"),
	}
}

func flattenInlineHookAuth(d *schema.ResourceData, c *sdk.InlineHookChannel) map[string]interface{} {
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

func flattenInlineHookChannel(c *sdk.InlineHookChannel) map[string]interface{} {
	return map[string]interface{}{
		"type":    c.Type,
		"version": c.Version,
		"uri":     c.Config.URI,
		"method":  c.Config.Method,
	}
}

func flattenInlineHookHeaders(c *sdk.InlineHookChannel) *schema.Set {
	headers := make([]interface{}, len(c.Config.Headers))
	for i, header := range c.Config.Headers {
		headers[i] = map[string]interface{}{
			"key":   header.Key,
			"value": header.Value,
		}
	}

	return schema.NewSet(schema.HashResource(headerSchema), headers)
}

func setHookStatus(d *schema.ResourceData, client *sdk.ApiSupplement, status, desiredStatus string) error {
	if status != desiredStatus {
		if desiredStatus == statusInactive {
			return responseErr(client.DeactivateInlineHook(d.Id()))
		} else if desiredStatus == statusActive {
			return responseErr(client.ActivateInlineHook(d.Id()))
		}
	}

	return nil
}

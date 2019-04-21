package okta

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/okta/okta-sdk-golang/okta/query"

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

		// For those familiar with Terraform schemas be sure to check the base hooklication schema and/or
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
				ValidateFunc: validation.StringInSlice(
					[]string{
						"com.okta.oauth2.tokens.transform",
						"com.okta.import.transform",
						"com.okta.saml.tokens.transform",
					},
					false,
				),
			},
			"version": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
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
						"auth": &schema.Schema{
							Type:     schema.TypeMap,
							Optional: true,
							Elem: schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Header name for authentication",
									},
									"type": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice(
											[]string{"HEADER"},
											false,
										),
									},
									"value": &schema.Schema{
										Type:      schema.TypeString,
										Optional:  true,
										Sensitive: true,
									},
								},
							},
						},
						"headers": &schema.Schema{
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &headerSchema,
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
	activate := d.Get("status").(string) == "ACTIVE"
	params := &query.Params{Activate: &activate}
	newHook, _, err := client.CreateInlineHook(*hook, params)
	if err != nil {
		return err
	}

	d.SetId(newHook.ID)

	return resourceInlineHookRead(d, m)
}

func resourceInlineHookExists(d *schema.ResourceData, m interface{}) (bool, error) {
	_, res, err := getSupplementFromMetadata(m).GetInlineHook(d.Id())
	return err == nil && res.StatusCode != http.StatusNotFound, err
}

func resourceInlineHookRead(d *schema.ResourceData, m interface{}) error {
	hook, _, err := getSupplementFromMetadata(m).GetInlineHook(d.Id())
	if err != nil {
		return err
	}
	d.Set("name", hook.Name)
	d.Set("status", hook.Status)
	d.Set("type", hook.Type)
	d.Set("version", hook.Version)

	return setNonPrimitives(d, map[string]interface{}{
		"channel": flattenHookChannel(hook.Channel),
	})
}

func resourceInlineHookUpdate(d *schema.ResourceData, m interface{}) error {
	client := getSupplementFromMetadata(m)
	hook := buildInlineHook(d, m)
	_, _, err := client.UpdateInlineHook(d.Id(), *hook, nil)

	if err != nil {
		return err
	}

	desiredStatus := d.Get("status").(string)
	err = setHookStatus(d, client, hook.Status, desiredStatus)
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

func buildInlineHook(d *schema.ResourceData, m interface{}) *InlineHook {
	return &InlineHook{
		Name:    d.Get("name").(string),
		Status:  d.Get("status").(string),
		Type:    d.Get("type").(string),
		Version: d.Get("version").(string),
		Channel: buildInlineChannel(d, m),
	}
}

func buildInlineChannel(d *schema.ResourceData, m interface{}) *Channel {
	if _, ok := d.GetOk("channel"); !ok {
		return nil
	}

	headerList := []*Header{}
	if raw, ok := d.GetOk("channel.headers"); ok {
		for _, header := range raw.(*schema.Set).List() {
			h, ok := header.(map[string]string)
			if ok {
				headerList = append(headerList, &Header{Key: h["key"], Value: h["type"]})
			}
		}
	}
	var auth *AuthScheme
	if _, ok := d.GetOk("channel.auth.type"); ok {
		auth = &AuthScheme{
			Key:   getStringValue(d, "channel.auth.key"),
			Type:  getStringValue(d, "channel.auth.type"),
			Value: getStringValue(d, "channel.auth.value"),
		}
	}

	return &Channel{
		Config: &HookConfig{
			URI:        getStringValue(d, "channel.uri"),
			AuthScheme: auth,
			Headers:    headerList,
		},
		Type:    getStringValue(d, "channel.type"),
		Version: getStringValue(d, "channel.version"),
	}
}

func flattenHookChannel(c *Channel) map[string]interface{} {
	headers := make([]interface{}, len(c.Config.Headers))
	for i, header := range c.Config.Headers {
		headers[i] = map[string]interface{}{
			"key":   header.Key,
			"value": header.Value,
		}
	}

	return map[string]interface{}{
		"type":    c.Type,
		"version": c.Version,
		"uri":     c.Config.URI,
		"auth": map[string]interface{}{
			"key":   c.Config.AuthScheme.Key,
			"type":  c.Config.AuthScheme.Type,
			"value": c.Config.AuthScheme.Value,
		},
		"headers": schema.NewSet(schema.HashResource(headerSchema), headers),
	}
}

func setHookStatus(d *schema.ResourceData, client *ApiSupplement, status string, desiredStatus string) error {
	if status != desiredStatus {
		if desiredStatus == "INACTIVE" {
			return responseErr(client.DeactivateInlineHook(d.Id()))
		} else if desiredStatus == "ACTIVE" {
			return responseErr(client.ActivateInlineHook(d.Id()))
		}
	}

	return nil
}

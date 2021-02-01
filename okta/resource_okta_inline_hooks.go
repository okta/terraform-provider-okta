package okta

import (
	"context"
	"net/http"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
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
		CreateContext: resourceInlineHookCreate,
		ReadContext:   resourceInlineHookRead,
		UpdateContext: resourceInlineHookUpdate,
		DeleteContext: resourceInlineHookDelete,
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
				ValidateDiagFunc: stringInSlice([]string{
					"com.okta.oauth2.tokens.transform",
					"com.okta.import.transform",
					"com.okta.saml.tokens.transform",
					"com.okta.user.pre-registration",
					"com.okta.user.credential.password.import",
				}),
			},
			"version": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: stringIsVersion,
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
					if k == "channel.method" && new == "" {
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
					dErr := stringIsURL("https")(m["uri"], cty.GetAttrPath("uri"))
					if dErr != nil {
						errs = append(errs, dErr...)
					}
					dErr = stringIsVersion(m["version"], cty.GetAttrPath("version"))
					if dErr != nil {
						errs = append(errs, dErr...)
					}
					if method, ok := m["method"]; ok {
						dErr = stringInSlice([]string{
							http.MethodPost,
							http.MethodGet,
							http.MethodPut,
							http.MethodDelete,
							http.MethodPatch,
						})(method, cty.GetAttrPath("method"))
						if dErr != nil {
							errs = append(errs, dErr...)
						}
					}
					return errs
				},
			},
		},
	}
}

func resourceInlineHookCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hook := buildInlineHook(d)
	newHook, _, err := getOktaClientFromMetadata(m).InlineHook.CreateInlineHook(ctx, hook)
	if err != nil {
		return diag.Errorf("failed to create inline hook: %v", err)
	}
	d.SetId(newHook.Id)
	err = setInlineHookStatus(ctx, d, getOktaClientFromMetadata(m), newHook.Status)
	if err != nil {
		return diag.Errorf("failed to change inline hook's status: %v", err)
	}
	return resourceInlineHookRead(ctx, d, m)
}

func resourceInlineHookRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hook, resp, err := getOktaClientFromMetadata(m).InlineHook.GetInlineHook(ctx, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get inline hook: %v", err)
	}
	if hook == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("name", hook.Name)
	_ = d.Set("status", hook.Status)
	_ = d.Set("type", hook.Type)
	_ = d.Set("version", hook.Version)
	err = setNonPrimitives(d, map[string]interface{}{
		"channel": flattenInlineHookChannel(hook.Channel),
		"headers": flattenInlineHookHeaders(hook.Channel),
		"auth":    flattenInlineHookAuth(d, hook.Channel),
	})
	if err != nil {
		return diag.Errorf("failed to set inline hook properties: %v", err)
	}
	return nil
}

func resourceInlineHookUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(m)
	hook := buildInlineHook(d)
	newHook, _, err := client.InlineHook.UpdateInlineHook(ctx, d.Id(), hook)
	if err != nil {
		return diag.Errorf("failed to update inline hook: %v", err)
	}
	err = setInlineHookStatus(ctx, d, client, newHook.Status)
	if err != nil {
		return diag.Errorf("failed to change inline hook's status: %v", err)
	}
	return resourceInlineHookRead(ctx, d, m)
}

func resourceInlineHookDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(m)
	_, resp, err := client.InlineHook.DeactivateInlineHook(ctx, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to deactivate inline hook: %v", err)
	}
	resp, err = client.InlineHook.DeleteInlineHook(ctx, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to delete inline hook: %v", err)
	}
	return nil
}

func buildInlineHook(d *schema.ResourceData) okta.InlineHook {
	return okta.InlineHook{
		Name:    d.Get("name").(string),
		Status:  d.Get("status").(string),
		Type:    d.Get("type").(string),
		Version: d.Get("version").(string),
		Channel: buildInlineChannel(d),
	}
}

func buildInlineChannel(d *schema.ResourceData) *okta.InlineHookChannel {
	var headerList []*okta.InlineHookChannelConfigHeaders
	if raw, ok := d.GetOk("headers"); ok {
		for _, header := range raw.(*schema.Set).List() {
			h, ok := header.(map[string]interface{})
			if ok {
				headerList = append(headerList, &okta.InlineHookChannelConfigHeaders{Key: h["key"].(string), Value: h["value"].(string)})
			}
		}
	}
	var auth *okta.InlineHookChannelConfigAuthScheme
	if rawAuth, ok := d.GetOk("auth"); ok {
		a := rawAuth.(map[string]interface{})
		_, ok := a["type"]
		if !ok {
			a["type"] = "HEADER"
		}
		auth = &okta.InlineHookChannelConfigAuthScheme{
			Key:   a["key"].(string),
			Type:  a["type"].(string),
			Value: a["value"].(string),
		}
	}
	rawChannel := d.Get("channel").(map[string]interface{})
	_, ok := rawChannel["method"]
	if !ok {
		rawChannel["method"] = "POST"
	}
	_, ok = rawChannel["type"]
	if !ok {
		rawChannel["type"] = "HTTP"
	}
	return &okta.InlineHookChannel{
		Config: &okta.InlineHookChannelConfig{
			Uri:        rawChannel["uri"].(string),
			AuthScheme: auth,
			Headers:    headerList,
			Method:     rawChannel["method"].(string),
		},
		Type:    rawChannel["type"].(string),
		Version: rawChannel["version"].(string),
	}
}

func flattenInlineHookAuth(d *schema.ResourceData, c *okta.InlineHookChannel) map[string]interface{} {
	auth := map[string]interface{}{}
	if c.Config.AuthScheme != nil {
		auth = map[string]interface{}{
			"key":  c.Config.AuthScheme.Key,
			"type": c.Config.AuthScheme.Type,
			// Read only
			"value": d.Get("auth").(map[string]interface{})["value"],
		}
	}
	return auth
}

func flattenInlineHookChannel(c *okta.InlineHookChannel) map[string]interface{} {
	return map[string]interface{}{
		"type":    c.Type,
		"version": c.Version,
		"uri":     c.Config.Uri,
		"method":  c.Config.Method,
	}
}

func flattenInlineHookHeaders(c *okta.InlineHookChannel) *schema.Set {
	headers := make([]interface{}, len(c.Config.Headers))
	for i, header := range c.Config.Headers {
		headers[i] = map[string]interface{}{
			"key":   header.Key,
			"value": header.Value,
		}
	}
	return schema.NewSet(schema.HashResource(headerSchema), headers)
}

func setInlineHookStatus(ctx context.Context, d *schema.ResourceData, client *okta.Client, status string) error {
	desiredStatus := d.Get("status").(string)
	if status == desiredStatus {
		return nil
	}
	var err error
	if desiredStatus == statusInactive {
		_, _, err = client.InlineHook.DeactivateInlineHook(ctx, d.Id())
	} else {
		_, _, err = client.InlineHook.ActivateInlineHook(ctx, d.Id())
	}
	return err
}

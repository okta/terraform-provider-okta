package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
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
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if k == "auth.type" && new == "" {
						return true
					}
					return false
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

func buildInlineHook(d *schema.ResourceData) sdk.InlineHook {
	return sdk.InlineHook{
		Name:    d.Get("name").(string),
		Status:  d.Get("status").(string),
		Type:    d.Get("type").(string),
		Version: d.Get("version").(string),
		Channel: buildInlineChannel(d),
	}
}

func buildInlineChannel(d *schema.ResourceData) *sdk.InlineHookChannel {
	var headerList []*sdk.InlineHookChannelConfigHeaders
	if raw, ok := d.GetOk("headers"); ok {
		for _, header := range raw.(*schema.Set).List() {
			h, ok := header.(map[string]interface{})
			if ok {
				headerList = append(headerList, &sdk.InlineHookChannelConfigHeaders{Key: h["key"].(string), Value: h["value"].(string)})
			}
		}
	}
	var auth *sdk.InlineHookChannelConfigAuthScheme
	if rawAuth, ok := d.GetOk("auth"); ok {
		a := rawAuth.(map[string]interface{})
		_, ok := a["type"]
		if !ok {
			a["type"] = "HEADER"
		}
		auth = &sdk.InlineHookChannelConfigAuthScheme{}
		if key, ok := a["key"]; ok && key != nil {
			auth.Key = key.(string)
		}
		if _type, ok := a["type"]; ok && _type != nil {
			auth.Type = _type.(string)
		}
		if value, ok := a["value"]; ok && value != nil {
			auth.Value = value.(string)
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
	return &sdk.InlineHookChannel{
		Config: &sdk.InlineHookChannelConfig{
			Uri:        rawChannel["uri"].(string),
			AuthScheme: auth,
			Headers:    headerList,
			Method:     rawChannel["method"].(string),
		},
		Type:    rawChannel["type"].(string),
		Version: rawChannel["version"].(string),
	}
}

func flattenInlineHookAuth(d *schema.ResourceData, c *sdk.InlineHookChannel) map[string]interface{} {
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

func flattenInlineHookChannel(c *sdk.InlineHookChannel) map[string]interface{} {
	return map[string]interface{}{
		"type":    c.Type,
		"version": c.Version,
		"uri":     c.Config.Uri,
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

func setInlineHookStatus(ctx context.Context, d *schema.ResourceData, client *sdk.Client, status string) error {
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

package idaas

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
)

var HeaderSchema = &schema.Resource{
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
		Description: "Creates an inline hook. This resource allows you to create and configure an inline hook.",
		// For those familiar with Terraform schemas be sure to check the base hook schema and/or
		// the examples in the documentation
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The inline hook display name.",
			},
			"status": statusSchema,
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The type of hook to create. [See here for supported types](https://developer.okta.com/docs/reference/api/inline-hooks/#supported-inline-hook-types).",
			},
			"version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The version of the hook. The currently-supported version is `1.0.0`.",
			},
			"headers": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        HeaderSchema,
				Description: "Map of headers to send along in inline hook request.",
			},
			// channel and auth presumed to work together
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
				ConflictsWith: []string{"channel_json"},
			},
			"channel": {
				Type:     schema.TypeMap,
				Optional: true,
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
				ConflictsWith: []string{"channel_json"},
			},
			"channel_json": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "true channel object for the inline hook API contract",
				ValidateDiagFunc: stringIsJSON,
				StateFunc:        utils.NormalizeDataJSON,
				DiffSuppressFunc: noChangeInObjectFromUnmarshaledChannelJSON,
				ConflictsWith:    []string{"channel", "auth"},
			},
		},
	}
}

func resourceInlineHookCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	hook := buildInlineHook(d)
	newHook, _, err := getOktaClientFromMetadata(meta).InlineHook.CreateInlineHook(ctx, hook)
	if err != nil {
		return diag.Errorf("failed to create inline hook: %v", err)
	}
	d.SetId(newHook.Id)
	err = setInlineHookStatus(ctx, d, getOktaClientFromMetadata(meta), newHook.Status)
	if err != nil {
		return diag.Errorf("failed to change inline hook's status: %v", err)
	}
	return resourceInlineHookRead(ctx, d, meta)
}

func resourceInlineHookRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	hook, resp, err := getOktaClientFromMetadata(meta).InlineHook.GetInlineHook(ctx, d.Id())
	if err := utils.SuppressErrorOn404(resp, err); err != nil {
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

	if oldChannelJson, ok := d.GetOk("channel_json"); ok {
		// NOTE: Okta responses don't include config.clientSecret so copy the
		// secret over if it exists the existing channel json
		var oldChannel sdk.InlineHookChannel
		if err = json.Unmarshal([]byte(oldChannelJson.(string)), &oldChannel); err == nil {
			if oldChannel.Config != nil && oldChannel.Config.ClientSecret != "" {
				if hook.Channel != nil && hook.Channel.Config != nil {
					hook.Channel.Config.ClientSecret = oldChannel.Config.ClientSecret
				}
			}
			if oldChannel.Config != nil && oldChannel.Config.AuthScheme != nil {
				if hook.Channel != nil && hook.Channel.Config != nil {
					hook.Channel.Config.AuthScheme.Value = oldChannel.Config.AuthScheme.Value
				}
			}
		}

		channelJson, err := json.Marshal(hook.Channel)
		if err != nil {
			return diag.Errorf("error marshaling channel json: %v", err)
		}
		_ = d.Set("channel_json", string(channelJson))
	} else {
		err = utils.SetNonPrimitives(d, map[string]interface{}{
			"channel": flattenInlineHookChannel(hook.Channel),
			"headers": flattenInlineHookHeaders(hook.Channel),
			"auth":    flattenInlineHookAuth(d, hook.Channel),
		})
	}
	if err != nil {
		return diag.Errorf("failed to set inline hook properties: %v", err)
	}
	return nil
}

func resourceInlineHookUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(meta)
	hook := buildInlineHook(d)
	newHook, _, err := client.InlineHook.UpdateInlineHook(ctx, d.Id(), hook)
	if err != nil {
		return diag.Errorf("failed to update inline hook: %v", err)
	}
	err = setInlineHookStatus(ctx, d, client, newHook.Status)
	if err != nil {
		return diag.Errorf("failed to change inline hook's status: %v", err)
	}
	return resourceInlineHookRead(ctx, d, meta)
}

func resourceInlineHookDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(meta)
	_, resp, err := client.InlineHook.DeactivateInlineHook(ctx, d.Id())
	if err := utils.SuppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to deactivate inline hook: %v", err)
	}
	resp, err = client.InlineHook.DeleteInlineHook(ctx, d.Id())
	if err := utils.SuppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to delete inline hook: %v", err)
	}
	return nil
}

func buildInlineHook(d *schema.ResourceData) sdk.InlineHook {
	inlineHook := sdk.InlineHook{
		Name:    d.Get("name").(string),
		Status:  d.Get("status").(string),
		Type:    d.Get("type").(string),
		Version: d.Get("version").(string),
	}
	if channelJson, ok := d.GetOk("channel_json"); ok {
		var channel sdk.InlineHookChannel
		_ = json.Unmarshal([]byte(channelJson.(string)), &channel)
		inlineHook.Channel = &channel
	} else {
		inlineHook.Channel = buildInlineChannel(d)
	}
	return inlineHook
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
	return schema.NewSet(schema.HashResource(HeaderSchema), headers)
}

func setInlineHookStatus(ctx context.Context, d *schema.ResourceData, client *sdk.Client, status string) error {
	desiredStatus := d.Get("status").(string)
	if status == desiredStatus {
		return nil
	}
	var err error
	if desiredStatus == StatusInactive {
		_, _, err = client.InlineHook.DeactivateInlineHook(ctx, d.Id())
	} else {
		_, _, err = client.InlineHook.ActivateInlineHook(ctx, d.Id())
	}
	return err
}

// noChangeInObjectFromUnmarshaledChannelJSON is a DiffSuppressFunc returns and
// true if old and new JSONs are equivalent object representations ...  It is
// true, there is no change!  Edge chase if newJSON is blank, will also return
// true which cover the new resource case.  Okta does not return
// config.clientSecret, config.authScheme.value, auth.value in the response so ignore these values.
// The above mentioned fields are mutually exclusive.
// TODO Explore usage of Terraform's write only attributes for sensitive values.
func noChangeInObjectFromUnmarshaledChannelJSON(k, oldJSON, newJSON string, d *schema.ResourceData) bool {
	if newJSON == "" {
		return true
	}
	var oldObj map[string]any
	var newObj map[string]any
	if err := json.Unmarshal([]byte(oldJSON), &oldObj); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(newJSON), &newObj); err != nil {
		return false
	}
	return reflect.DeepEqual(oldObj, newObj)
}

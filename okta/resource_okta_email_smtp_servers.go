package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v4/okta"
)

func resourceEmailSmtp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEmailSmtpCreate,
		ReadContext:   resourceEmailSmtpRead,
		UpdateContext: resourceEmailSmtpUpdate,
		DeleteContext: resourceEmailSmtpDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: `configure a custom external email provider to send email notifications. 
		By default, notifications such as the welcome email or an account recovery email are sent through an Okta-managed SMTP server.`,
		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Hostname or IP address of your SMTP server.",
			},
			"port": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Port number of your SMTP server.",
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name of the email domain.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "User name of the email domain.",
			},
			"alias": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Human-readable name for your SMTP server.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If true, routes all email traffic through your SMTP server.",
			},
		},
	}
}

func resourceEmailSmtpCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	emailSmtpResp, _, err := getOktaV3ClientFromMetadata(meta).EmailServerAPI.CreateEmailServer(ctx).EmailServerPost(buildEmailSmtp(d)).Execute()
	if err != nil {
		return nil
	}
	d.SetId(*emailSmtpResp.Id)
	return resourceEmailSmtpRead(ctx, d, meta)
}

func resourceEmailSmtpRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	emailSmtp, resp, err := getOktaV3ClientFromMetadata(meta).EmailServerAPI.GetEmailServer(ctx, d.Id()).Execute()
	if err := v3suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get email domain: %v", err)
	}
	if resp == nil || resp.StatusCode == 404 {
		d.SetId("")
		return nil
	}

	if emailSmtp == nil {
		return diag.Errorf("emailSmtp is nil but no 404 returned")
	}

	properties := emailSmtp.AdditionalProperties

	logger(meta).Info("emailSmtp found", "configs", *emailSmtp)

	_ = d.Set("host", properties["host"].(string))
	_ = d.Set("alias", properties["alias"].(string))
	_ = d.Set("enabled", properties["enabled"].(bool))
	_ = d.Set("username", properties["username"].(string))
	_ = d.Set("port", properties["port"].(float64))
	return nil
}

func resourceEmailSmtpUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logger(meta).Info("resourceEmailSmtpUpdate", "server id", d.Id(), "SERVER ALIAS", d.Get("alias"), "SERVER HOST", d.Get("host"), "SERVER PORT", d.Get("port"))
	req := buildEmailServerRequest(d)
	logger(meta).Info("resourceEmailSmtpUpdate", "update Request Body", "alias", *req.Alias, "port", req.Port)
	_, _, err := getOktaV3ClientFromMetadata(meta).EmailServerAPI.UpdateEmailServer(ctx, d.Id()).EmailServerRequest(req).Execute()
	if err != nil {
		return nil
	}
	return resourceEmailSmtpRead(ctx, d, meta)
}

func resourceEmailSmtpDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logger(meta).Info("resourceEmailSmtpDelete", "server configs", d.Id())
	_, err := getOktaV3ClientFromMetadata(meta).EmailServerAPI.DeleteEmailServer(ctx, d.Id()).Execute()
	if err != nil {
		return diag.Errorf("failed to delete email domain: %v", err)
	}
	return nil
}

func buildEmailSmtp(d *schema.ResourceData) okta.EmailServerPost {
	return okta.EmailServerPost{
		Alias:    d.Get("alias").(string),
		Host:     d.Get("host").(string),
		Port:     int32(d.Get("port").(int)),
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
		Enabled:  boolPtr(d.Get("enabled").(bool)),
	}
}

func interfaceToStringPointer(value interface{}) *string {
	if str, ok := value.(string); ok {
		return &str
	}
	return nil
}

func interfaceToInt32Pointer(value interface{}) *int32 {
	fmt.Println("port number", value)
	if i, ok := value.(int); ok {
		x := int32(i)
		return &x
	}
	return nil
}

func buildEmailServerRequest(d *schema.ResourceData) okta.EmailServerRequest {
	return okta.EmailServerRequest{
		Alias:    interfaceToStringPointer(d.Get("alias")),
		Host:     interfaceToStringPointer(d.Get("host")),
		Port:     interfaceToInt32Pointer(d.Get("port")),
		Username: interfaceToStringPointer(d.Get("username")),
		Password: interfaceToStringPointer(d.Get("password")),
		Enabled:  boolPtr(d.Get("enabled").(bool)),
	}
}

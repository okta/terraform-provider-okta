package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	v5okta "github.com/okta/okta-sdk-golang/v5/okta"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

func resourceEmailSMTP() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEmailSMTPCreate,
		ReadContext:   resourceEmailSMTPRead,
		UpdateContext: resourceEmailSMTPUpdate,
		DeleteContext: resourceEmailSMTPDelete,
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

func resourceEmailSMTPCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	emailSMTPResp, _, err := getOktaV5ClientFromMetadata(meta).EmailServerAPI.CreateEmailServer(ctx).EmailServerPost(buildEmailSMTP(d)).Execute()
	if err != nil {
		return nil
	}
	d.SetId(*emailSMTPResp.Id)
	return resourceEmailSMTPRead(ctx, d, meta)
}

func resourceEmailSMTPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	emailSMTP, resp, err := getOktaV5ClientFromMetadata(meta).EmailServerAPI.GetEmailServer(ctx, d.Id()).Execute()
	if err := utils.SuppressErrorOn404_V5(resp, err); err != nil {
		return diag.Errorf("failed to get email domain: %v", err)
	}
	if resp == nil || resp.StatusCode == 404 {
		d.SetId("")
		return nil
	}

	if emailSMTP == nil {
		return diag.Errorf("emailSMTPServer is nil but no 404 returned")
	}

	properties := emailSMTP.AdditionalProperties

	_ = d.Set("host", properties["host"].(string))
	_ = d.Set("alias", properties["alias"].(string))
	_ = d.Set("enabled", properties["enabled"].(bool))
	_ = d.Set("username", properties["username"].(string))
	_ = d.Set("port", properties["port"].(float64))
	return nil
}

func resourceEmailSMTPUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	req := buildEmailServerRequest(d)
	_, _, err := getOktaV5ClientFromMetadata(meta).EmailServerAPI.UpdateEmailServer(ctx, d.Id()).EmailServerRequest(req).Execute()
	if err != nil {
		return nil
	}
	return resourceEmailSMTPRead(ctx, d, meta)
}

func resourceEmailSMTPDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_, err := getOktaV5ClientFromMetadata(meta).EmailServerAPI.DeleteEmailServer(ctx, d.Id()).Execute()
	if err != nil {
		return diag.Errorf("failed to delete email domain: %v", err)
	}
	return nil
}

func buildEmailSMTP(d *schema.ResourceData) v5okta.EmailServerPost {
	return v5okta.EmailServerPost{
		Alias:    d.Get("alias").(string),
		Host:     d.Get("host").(string),
		Port:     int32(d.Get("port").(int)),
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
		Enabled:  utils.BoolPtr(d.Get("enabled").(bool)),
	}
}

func buildEmailServerRequest(d *schema.ResourceData) v5okta.EmailServerRequest {
	return v5okta.EmailServerRequest{
		Alias:    interfaceToStringPointer(d.Get("alias")),
		Host:     interfaceToStringPointer(d.Get("host")),
		Port:     interfaceToInt32Pointer(d.Get("port")),
		Username: interfaceToStringPointer(d.Get("username")),
		Password: interfaceToStringPointer(d.Get("password")),
		Enabled:  utils.BoolPtr(d.Get("enabled").(bool)),
	}
}

func interfaceToStringPointer(value interface{}) *string {
	if str, ok := value.(string); ok {
		return &str
	}
	return nil
}

func interfaceToInt32Pointer(value interface{}) *int32 {
	if i, ok := value.(int); ok {
		x := int32(i)
		return &x
	}
	return nil
}

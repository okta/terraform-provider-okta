package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceEmailSMTPServers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEmailSMTPServersRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the SMTP server.",
			},
			"username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the SMTP server.",
			},
			"host": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SMTP server host name.",
			},
			"port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "SMTP server port number.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the SMTP server requires a secure connection.",
			},
			"alias": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Human-readable name for your SMTP server.",
			},
		},
		Description: "Get the enrolled email SMTP server.",
	}
}

func dataSourceEmailSMTPServersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	emailSMTPServerId, ok := d.GetOk("id")
	if !ok {
		return diag.Errorf("id required for email SMTP servers")
	}
	emailSMTPServers, _, _ := getOktaV5ClientFromMetadata(meta).EmailServerAPI.GetEmailServer(ctx, emailSMTPServerId.(string)).Execute()
	properties := emailSMTPServers.AdditionalProperties
	d.SetId(emailSMTPServerId.(string))
	_ = d.Set("host", properties["host"].(string))
	_ = d.Set("alias", properties["alias"].(string))
	_ = d.Set("enabled", properties["enabled"].(bool))
	_ = d.Set("username", properties["username"].(string))
	_ = d.Set("port", properties["port"].(float64))
	return nil
}

package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

func dataSourceAuthServer() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAuthServerRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"audiences": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"kid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"credentials_last_rotated": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"credentials_next_rotation": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"credentials_rotation_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAuthServerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	servers, _, err := getOktaClientFromMetadata(m).AuthorizationServer.ListAuthorizationServers(context.Background(), &query.Params{Q: name})
	if err != nil {
		return diag.Errorf("failed to find auth server '%s': %v", name, err)
	}
	if len(servers) < 1 {
		diag.Errorf("authorization server with name '%s' does not exist", name)
	}
	var authServer *okta.AuthorizationServer
	if len(servers) > 1 {
		for _, server := range servers {
			if server.Name == name {
				authServer = server
				break
			}
		}
	}
	if authServer == nil {
		return diag.Errorf("authorization server with name '%s' does not exist", name)
	}
	d.SetId(authServer.Id)
	_ = d.Set("name", authServer.Name)
	_ = d.Set("description", authServer.Description)
	_ = d.Set("audiences", convertStringSetToInterface(authServer.Audiences))
	_ = d.Set("credentials_rotation_mode", authServer.Credentials.Signing.RotationMode)
	_ = d.Set("kid", authServer.Credentials.Signing.Kid)
	_ = d.Set("credentials_next_rotation", authServer.Credentials.Signing.NextRotation.String())
	_ = d.Set("credentials_last_rotated", authServer.Credentials.Signing.LastRotated.String())
	_ = d.Set("status", authServer.Status)
	return nil
}

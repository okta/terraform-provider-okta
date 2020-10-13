package okta

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

func dataSourceAuthServer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAuthServerRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"audiences": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"kid": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"credentials_last_rotated": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"credentials_next_rotation": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"credentials_rotation_mode": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAuthServerRead(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	authServer, err := getSupplementFromMetadata(m).FindAuthServer(name, &query.Params{})
	if err != nil {
		return err
	}
	if authServer == nil {
		return fmt.Errorf("No authorization server found with provided name %s", name)
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

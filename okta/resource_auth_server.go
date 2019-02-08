package okta

import (
	"net/http"

	"github.com/hashicorp/terraform/helper/validation"

	"github.com/okta/okta-sdk-golang/okta"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAuthServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceAuthServerCreate,
		Exists: resourceAuthServerExists,
		Read:   resourceAuthServerRead,
		Update: resourceAuthServerUpdate,
		Delete: resourceAuthServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "AuthServer name",
			},
			"credentials": &schema.Schema{
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "AuthServer credentials",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kid": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_rotated": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"next_rotation": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"rotation_mode": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"AUTO", "MANUAL"}, false),
							Default:      "AUTO",
						},
					},
				},
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AuthServer description",
			},
		},
	}
}

func buildAuthServer(d *schema.ResourceData) *okta.AuthorizationServer {
	return &okta.AuthorizationServer{
		Audiences: convertInterfaceToStringArrNullable(d.Get("audiences")),
		Credentials: &okta.AuthServerCredentials{
			Signing: &okta.ApplicationCredentialsSigning{
				RotationMode: d.Get("credentials.rotation_mode").(string),
			},
		},
		Description: d.Get("description").(string),
		Name:        d.Get("name").(string),
	}
}

func resourceAuthServerCreate(d *schema.ResourceData, m interface{}) error {
	authServer := buildAuthServer(d)
	responseAuthServer, _, err := getOktaClientFromMetadata(m).AuthorizationServer.CreateAuthorizationServer(*authServer, nil)
	if err != nil {
		return err
	}

	d.SetId(responseAuthServer.Id)

	return resourceAuthServerRead(d, m)
}

func resourceAuthServerExists(d *schema.ResourceData, m interface{}) (bool, error) {
	g, err := fetchAuthServer(d, m)

	return err == nil && g != nil, err
}

func resourceAuthServerRead(d *schema.ResourceData, m interface{}) error {
	authServer, err := fetchAuthServer(d, m)
	if err != nil {
		return err
	}

	d.Set("name", authServer.Name)
	d.Set("description", authServer.Description)
	d.Set("credentials", flattenCredentials(authServer.Credentials))
	d.Set("audiences", authServer.Audiences)

	return nil
}

func resourceAuthServerUpdate(d *schema.ResourceData, m interface{}) error {
	authServer := buildAuthServer(d)
	_, _, err := getOktaClientFromMetadata(m).AuthorizationServer.UpdateAuthorizationServer(d.Id(), *authServer, nil)
	if err != nil {
		return err
	}

	return resourceAuthServerRead(d, m)
}

func resourceAuthServerDelete(d *schema.ResourceData, m interface{}) error {
	_, err := getOktaClientFromMetadata(m).AuthorizationServer.DeleteAuthorizationServer(d.Id())

	return err
}

func fetchAuthServer(d *schema.ResourceData, m interface{}) (*okta.AuthorizationServer, error) {
	auth, resp, err := getOktaClientFromMetadata(m).AuthorizationServer.GetAuthorizationServer(d.Id(), okta.AuthorizationServer{})

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	return auth, err
}

func flattenCredentials(creds *okta.AuthServerCredentials) interface{} {
	return map[string]interface{}{
		"kid":           creds.Signing.Kid,
		"last_rotated":  creds.Signing.LastRotated,
		"next_rotation": creds.Signing.NextRotation,
		"rotation_mode": creds.Signing.RotationMode,
	}
}

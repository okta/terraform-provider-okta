package okta

import (
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	"github.com/okta/okta-sdk-golang/okta"
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
			"audiences": &schema.Schema{
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Currently Okta only supports a single value here",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"credentials": &schema.Schema{
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Auth Server credentials",
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
				Description: "Auth Server description",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Auth Server name",
			},
		},
	}
}

func handleAuthServerLifecycle(d *schema.ResourceData, m interface{}) error {
	client := getSupplementFromMetadata(m)

	if d.Get("status").(string) == "ACTIVE" {
		_, err := client.ActivateAuthorizationServer(d.Id())
		return err
	}

	_, err := client.DeactivateAuthorizationServer(d.Id())
	return err
}

func buildAuthServer(d *schema.ResourceData) *AuthorizationServer {
	return &AuthorizationServer{
		Audiences: convertInterfaceToStringSet(d.Get("audiences")),
		Credentials: &AuthServerCredentials{
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
	responseAuthServer, _, err := getSupplementFromMetadata(m).CreateAuthorizationServer(*authServer, nil)
	if err != nil {
		return err
	}

	if err := handleAuthServerLifecycle(d, m); err != nil {
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

	d.Set("audiences", convertStringSetToInterface(authServer.Audiences))
	d.Set("credentials", flattenCredentials(authServer.Credentials))
	d.Set("description", authServer.Description)
	d.Set("name", authServer.Name)

	return nil
}

func resourceAuthServerUpdate(d *schema.ResourceData, m interface{}) error {
	client := getSupplementFromMetadata(m)
	if d.HasChange("status") {
		handleAuthServerLifecycle(d, m)
	}

	authServer := buildAuthServer(d)
	_, _, err := client.UpdateAuthorizationServer(d.Id(), *authServer, nil)
	if err != nil {
		return err
	}

	return resourceAuthServerRead(d, m)
}

func resourceAuthServerDelete(d *schema.ResourceData, m interface{}) error {
	client := getSupplementFromMetadata(m)
	if _, err := client.DeactivateAuthorizationServer(d.Id()); err != nil {
		return err
	}
	_, err := client.DeleteAuthorizationServer(d.Id())

	return err
}

func fetchAuthServer(d *schema.ResourceData, m interface{}) (*AuthorizationServer, error) {
	auth, resp, err := getSupplementFromMetadata(m).GetAuthorizationServer(d.Id(), AuthorizationServer{})

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	return auth, err
}

func flattenCredentials(creds *AuthServerCredentials) interface{} {
	return map[string]interface{}{
		"kid":           creds.Signing.Kid,
		"last_rotated":  creds.Signing.LastRotated,
		"next_rotation": creds.Signing.NextRotation,
		"rotation_mode": creds.Signing.RotationMode,
	}
}

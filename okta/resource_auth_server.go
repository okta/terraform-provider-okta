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
			"status": statusSchema,
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
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"AUTO", "MANUAL"}, false),
				Default:      "AUTO",
				Description:  "Credential rotation mode, in many cases you cannot set this to MANUAL, the API will ignore the value and you will get a perpetual diff. This should rarely be used.",
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
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
				RotationMode: d.Get("credentials_rotation_mode").(string),
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
	d.Set("credentials_rotation_mode", authServer.Credentials.Signing.RotationMode)
	d.Set("kid", authServer.Credentials.Signing.Kid)
	d.Set("credentials_next_rotation", authServer.Credentials.Signing.NextRotation)
	d.Set("credentials_last_rotated", authServer.Credentials.Signing.LastRotated)
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

package okta

import (
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	"github.com/okta/okta-sdk-golang/okta"
	"github.com/terraform-providers/terraform-provider-okta/sdk"
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
			"issuer": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "EA Feature: allows you to use a custom issuer URL",
			},
			"issuer_mode": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "EA Feature: allows you to use a custom issuer URL",
				Default:      "ORG_URL",
				ValidateFunc: validation.StringInSlice([]string{"CUSTOM_URL", "ORG_URL"}, false),
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

func buildAuthServer(d *schema.ResourceData) *sdk.AuthorizationServer {
	return &sdk.AuthorizationServer{
		Audiences: convertInterfaceToStringSet(d.Get("audiences")),
		Credentials: &sdk.AuthServerCredentials{
			Signing: &okta.ApplicationCredentialsSigning{
				RotationMode: d.Get("credentials_rotation_mode").(string),
			},
		},
		Description: d.Get("description").(string),
		Name:        d.Get("name").(string),
		IssuerMode:  d.Get("issuer_mode").(string),
	}
}

func resourceAuthServerCreate(d *schema.ResourceData, m interface{}) error {
	authServer := buildAuthServer(d)

	responseAuthServer, _, err := getSupplementFromMetadata(m).CreateAuthorizationServer(*authServer, nil)
	if err != nil {
		return err
	}

	d.SetId(responseAuthServer.Id)

	if d.Get("credentials_rotation_mode").(string) == "MANUAL" {
		// Auth servers can only be set to manual on update. No clue why.
		err = resourceAuthServerUpdate(d, m)

		if err != nil {
			return err
		}
	}

	return resourceAuthServerRead(d, m)
}

func resourceAuthServerExists(d *schema.ResourceData, m interface{}) (bool, error) {
	g, err := fetchAuthServer(d, m)

	return err == nil && g != nil, err
}

func resourceAuthServerRead(d *schema.ResourceData, m interface{}) error {
	authServer, err := fetchAuthServer(d, m)

	if authServer == nil {
		d.SetId("")
		return nil
	}

	if err != nil {
		return err
	}

	d.Set("audiences", convertStringSetToInterface(authServer.Audiences))
	d.Set("kid", authServer.Credentials.Signing.Kid)

	if authServer.Credentials != nil && authServer.Credentials.Signing != nil {
		d.Set("credentials_rotation_mode", authServer.Credentials.Signing.RotationMode)

		if authServer.Credentials.Signing.NextRotation != nil {
			d.Set("credentials_next_rotation", authServer.Credentials.Signing.NextRotation.String())
		}

		if authServer.Credentials.Signing.LastRotated != nil {
			d.Set("credentials_last_rotated", authServer.Credentials.Signing.LastRotated.String())
		}
	}

	d.Set("description", authServer.Description)
	d.Set("name", authServer.Name)
	d.Set("status", authServer.Status)
	d.Set("issuer", authServer.Issuer)

	// Do not sync these unless the issuer mode is specified since it is an EA feature and is computed in some cases
	if authServer.IssuerMode != "" {
		d.Set("issuer_mode", authServer.IssuerMode)
	}

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

func fetchAuthServer(d *schema.ResourceData, m interface{}) (*sdk.AuthorizationServer, error) {
	auth, resp, err := getSupplementFromMetadata(m).GetAuthorizationServer(d.Id())

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	return auth, err
}

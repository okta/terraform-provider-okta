package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func resourceAuthServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAuthServerCreate,
		ReadContext:   resourceAuthServerRead,
		UpdateContext: resourceAuthServerUpdate,
		DeleteContext: resourceAuthServerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"audiences": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Currently Okta only supports a single value here",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"status": statusSchema,
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
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: stringInSlice([]string{"AUTO", "MANUAL"}),
				Default:          "AUTO",
				Description:      "Credential rotation mode, in many cases you cannot set this to MANUAL, the API will ignore the value and you will get a perpetual diff. This should rarely be used.",
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"issuer": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "allows you to use a custom issuer URL",
			},
			"issuer_mode": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "*Early Access Property*. Indicates which value is specified in the issuer of the tokens that a Custom Authorization Server returns: the original Okta org domain URL or a custom domain URL",
				Default:          "ORG_URL",
				ValidateDiagFunc: stringInSlice([]string{"CUSTOM_URL", "ORG_URL"}),
			},
		},
	}
}

func resourceAuthServerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	authServer := buildAuthServer(d)
	responseAuthServer, _, err := getOktaClientFromMetadata(m).AuthorizationServer.CreateAuthorizationServer(ctx, *authServer)
	if err != nil {
		return diag.Errorf("failed to create authorization server: %v", err)
	}
	d.SetId(responseAuthServer.Id)
	if d.Get("credentials_rotation_mode").(string) == "MANUAL" {
		// Auth servers can only be set to manual on update. No clue why.
		dErr := resourceAuthServerUpdate(ctx, d, m)
		if dErr != nil {
			return dErr
		}
	}
	return resourceAuthServerRead(ctx, d, m)
}

func resourceAuthServerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	authServer, resp, err := getOktaClientFromMetadata(m).AuthorizationServer.GetAuthorizationServer(ctx, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get authorization server: %v", err)
	}
	if authServer == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("audiences", convertStringSetToInterface(authServer.Audiences))

	if authServer.Credentials != nil && authServer.Credentials.Signing != nil {
		_ = d.Set("kid", authServer.Credentials.Signing.Kid)
		_ = d.Set("credentials_rotation_mode", authServer.Credentials.Signing.RotationMode)
		if authServer.Credentials.Signing.NextRotation != nil {
			_ = d.Set("credentials_next_rotation", authServer.Credentials.Signing.NextRotation.String())
		}
		if authServer.Credentials.Signing.LastRotated != nil {
			_ = d.Set("credentials_last_rotated", authServer.Credentials.Signing.LastRotated.String())
		}
	}

	_ = d.Set("description", authServer.Description)
	_ = d.Set("name", authServer.Name)
	_ = d.Set("status", authServer.Status)
	_ = d.Set("issuer", authServer.Issuer)

	// Do not sync these unless the issuer mode is specified since it is an EA feature and is computed in some cases
	if authServer.IssuerMode != "" {
		_ = d.Set("issuer_mode", authServer.IssuerMode)
	}
	return nil
}

func resourceAuthServerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if d.HasChange("status") {
		err := handleAuthServerLifecycle(ctx, d, m)
		if err != nil {
			return err
		}
	}
	authServer := buildAuthServer(d)
	_, _, err := getOktaClientFromMetadata(m).AuthorizationServer.UpdateAuthorizationServer(ctx, d.Id(), *authServer)
	if err != nil {
		return diag.Errorf("failed to update authorization server: %v", err)
	}
	return resourceAuthServerRead(ctx, d, m)
}

func handleAuthServerLifecycle(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(m)
	if d.Get("status").(string) == statusActive {
		_, err := client.AuthorizationServer.ActivateAuthorizationServer(ctx, d.Id())
		if err != nil {
			return diag.Errorf("failed to activate authorization server: %v", err)
		}
		return nil
	}
	_, err := client.AuthorizationServer.DeactivateAuthorizationServer(ctx, d.Id())
	if err != nil {
		return diag.Errorf("failed to deactivate authorization server: %v", err)
	}
	return nil
}

func resourceAuthServerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(m)
	resp, err := client.AuthorizationServer.DeactivateAuthorizationServer(ctx, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to deactivate authorization server: %v", err)
	}
	resp, err = client.AuthorizationServer.DeleteAuthorizationServer(ctx, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to delete authorization server: %v", err)
	}
	return nil
}

func buildAuthServer(d *schema.ResourceData) *okta.AuthorizationServer {
	return &okta.AuthorizationServer{
		Audiences: convertInterfaceToStringSet(d.Get("audiences")),
		Credentials: &okta.AuthorizationServerCredentials{
			Signing: &okta.AuthorizationServerCredentialsSigningConfig{
				RotationMode: d.Get("credentials_rotation_mode").(string),
			},
		},
		Description: d.Get("description").(string),
		Name:        d.Get("name").(string),
		IssuerMode:  d.Get("issuer_mode").(string),
	}
}

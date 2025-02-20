package idaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
)

func ResourceAuthServerDefault() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAuthServerDefaultUpdate,
		ReadContext:   resourceAuthServerDefaultRead,
		UpdateContext: resourceAuthServerDefaultUpdate,
		DeleteContext: utils.ResourceFuncNoOp,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "Configures Default Authorization Server. This resource allows you to configure Default Authorization Server.",
		Schema: map[string]*schema.Schema{
			"audiences": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "The recipients that the tokens are intended for. This becomes the `aud` claim in an access token. Currently Okta only supports a single value here.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{"api://default"}, nil
				},
			},
			"status": statusSchema,
			"kid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the JSON Web Key used for signing tokens issued by the authorization server.",
			},
			"credentials_last_rotated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp when the authorization server started to use the `kid` for signing tokens.",
			},
			"credentials_next_rotation": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp when the authorization server changes the key for signing tokens. Only returned when `credentials_rotation_mode` is `AUTO`.",
			},
			"credentials_rotation_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The key rotation mode for the authorization server. Can be `AUTO` or `MANUAL`. Default: `MANUAL`.Credential rotation mode, in many cases you cannot set this to MANUAL, the API will ignore the value and you will get a perpetual diff. This should rarely be used.",
				Default:     "MANUAL",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Default Authorization Server for your Applications",
				Description: "The description of the authorization server.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "default",
				Description: "The name of the authorization server. Not necessary but left for backwards capacity with legacy implementation.",
			},
			"issuer": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The complete URL for a Custom Authorization Server. This becomes the `iss` claim in an access token.",
			},
			"issuer_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "*Early Access Property*. Allows you to use a custom issuer URL. It can be set to `CUSTOM_URL`, `ORG_URL`, or `DYNAMIC`. Default: `ORG_URL`",
			},
		},
	}
}

func resourceAuthServerDefaultRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	authServer, resp, err := GetOktaClientFromMetadata(meta).AuthorizationServer.GetAuthorizationServer(ctx, d.Id())
	if err := utils.SuppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get authorization server: %v", err)
	}
	if authServer == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("audiences", utils.ConvertStringSliceToSet(authServer.Audiences))
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

func getDefaultAuthServer(ctx context.Context, meta interface{}, serverID string) (authServer *sdk.AuthorizationServer, err error) {
	if serverID != "" {
		authServer, _, err = GetOktaClientFromMetadata(meta).AuthorizationServer.GetAuthorizationServer(ctx, serverID)
		return
	}

	authServers, _, err := GetOktaClientFromMetadata(meta).AuthorizationServer.ListAuthorizationServers(ctx, nil)
	if err != nil {
		return nil, err
	}
	for _, server := range authServers {
		if *server.Default {
			authServer = server
			return
		}
	}

	return nil, fmt.Errorf("failed to find default authorization server")
}

func resourceAuthServerDefaultUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id := d.Id()

	// when id is blank this is the create case as this function is used for
	// both create and update
	if id == "" {
		name := d.Get("name").(string)
		if name != "" {
			id = name
			// Legacy implementation allowed name to be set to something other
			// than "default" which was to loose of a constraint. In order to
			// not have breaking changes to existing configurations we will warn
			// when we find this situation.
			if name != "default" {
				Logger(meta).Warn("\"name\" argument is not \"default\". Allowing this legacy behavior.")
			}
		}
	}

	authServer, err := getDefaultAuthServer(ctx, meta, id)
	if err != nil {
		return diag.Errorf("failed to get default authorization server: %v", err)
	}
	id = authServer.Id

	if status, ok := d.GetOk("status"); ok {
		client := GetOktaClientFromMetadata(meta)
		if status.(string) == StatusActive && authServer.Status != StatusActive {
			_, err := client.AuthorizationServer.ActivateAuthorizationServer(ctx, id)
			if err != nil {
				return diag.Errorf("failed to activate default authorization server: %v", err)
			}
		}
		if status.(string) == StatusInactive && authServer.Status != StatusInactive {
			_, err := client.AuthorizationServer.DeactivateAuthorizationServer(ctx, id)
			if err != nil {
				return diag.Errorf("failed to deactivate default authorization server: %v", err)
			}
		}
	}

	_name, ok := d.GetOk("name")
	if ok {
		authServer.Name = _name.(string)
	}
	if *authServer.Default && authServer.Name == "" {
		authServer.Name = "default"
	}
	authServer.Audiences = utils.ConvertInterfaceToStringSet(d.Get("audiences"))
	authServer.Credentials.Signing.RotationMode = d.Get("credentials_rotation_mode").(string)
	authServer.Description = d.Get("description").(string)
	authServer.IssuerMode = d.Get("issuer_mode").(string)
	_, _, err = GetOktaClientFromMetadata(meta).AuthorizationServer.UpdateAuthorizationServer(ctx, id, *authServer)
	if err != nil {
		return diag.Errorf("failed to update default authorization server: %v", err)
	}
	d.SetId(authServer.Id)
	return resourceAuthServerDefaultRead(ctx, d, meta)
}

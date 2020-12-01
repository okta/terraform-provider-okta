package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
)

func resourceAuthServerScope() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAuthServerScopeCreate,
		ReadContext:   resourceAuthServerScopeRead,
		UpdateContext: resourceAuthServerScopeUpdate,
		DeleteContext: resourceAuthServerScopeDelete,
		Importer:      createNestedResourceImporter([]string{"auth_server_id", "id"}),
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Auth server scope name",
			},
			"auth_server_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Auth server ID",
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"consent": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "IMPLICIT",
				Description: "EA Feature and thus it is simply ignored if the feature is off",
			},
			"metadata_publish": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "ALL_CLIENTS",
				Description:      "Whether to publish metadata or not, matching API type despite the fact it could just be a boolean",
				ValidateDiagFunc: stringInSlice([]string{"ALL_CLIENTS", "NO_CLIENTS"}),
			},
			"default": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "A default scope will be returned in an access token when the client omits the scope parameter in a token request, provided this scope is allowed as part of the access policy rule.",
			},
		},
	}
}

func resourceAuthServerScopeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	authServerScope := buildAuthServerScope(d)
	responseAuthServerScope, _, err := getSupplementFromMetadata(m).CreateAuthorizationServerScope(
		ctx,
		d.Get("auth_server_id").(string),
		*authServerScope,
		nil,
	)
	if err != nil {
		return diag.Errorf("failed to create auth server scope: %v", err)
	}
	d.SetId(responseAuthServerScope.Id)
	return resourceAuthServerScopeRead(ctx, d, m)
}

func resourceAuthServerScopeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	authServerScope, resp, err := getSupplementFromMetadata(m).GetAuthorizationServerScope(
		ctx,
		d.Get("auth_server_id").(string),
		d.Id(),
		sdk.AuthorizationServerScope{},
	)
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get auth server scope: %v", err)
	}
	if authServerScope == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("name", authServerScope.Name)
	_ = d.Set("description", authServerScope.Description)
	_ = d.Set("metadata_publish", authServerScope.MetadataPublish)
	_ = d.Set("default", authServerScope.Default)
	if authServerScope.Consent != "" {
		_ = d.Set("consent", authServerScope.Consent)
	}
	return nil
}

func resourceAuthServerScopeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	authServerScope := buildAuthServerScope(d)
	_, _, err := getSupplementFromMetadata(m).UpdateAuthorizationServerScope(
		ctx, d.Get("auth_server_id").(string),
		d.Id(),
		*authServerScope,
		nil,
	)
	if err != nil {
		return diag.Errorf("failed to update auth server scope: %v", err)
	}
	return resourceAuthServerScopeRead(ctx, d, m)
}

func resourceAuthServerScopeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	_, err := getSupplementFromMetadata(m).DeleteAuthorizationServerScope(
		ctx,
		d.Get("auth_server_id").(string),
		d.Id(),
	)
	if err != nil {
		return diag.Errorf("failed to delete auth server scope: %v", err)
	}
	return nil
}

func buildAuthServerScope(d *schema.ResourceData) *sdk.AuthorizationServerScope {
	return &sdk.AuthorizationServerScope{
		Consent:         d.Get("consent").(string),
		Description:     d.Get("description").(string),
		MetadataPublish: d.Get("metadata_publish").(string),
		Name:            d.Get("name").(string),
		Default:         d.Get("default").(bool),
	}
}

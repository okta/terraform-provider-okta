package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
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
	scope := buildAuthServerScope(d)
	respScope, _, err := getOktaClientFromMetadata(m).AuthorizationServer.CreateOAuth2Scope(ctx, d.Get("auth_server_id").(string), scope)
	if err != nil {
		return diag.Errorf("failed to create auth server scope: %v", err)
	}
	d.SetId(respScope.Id)
	return resourceAuthServerScopeRead(ctx, d, m)
}

func resourceAuthServerScopeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	scope, resp, err := getOktaClientFromMetadata(m).AuthorizationServer.GetOAuth2Scope(ctx, d.Get("auth_server_id").(string), d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get auth server scope: %v", err)
	}
	if scope == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("name", scope.Name)
	_ = d.Set("description", scope.Description)
	_ = d.Set("metadata_publish", scope.MetadataPublish)
	_ = d.Set("default", scope.Default)
	if scope.Consent != "" {
		_ = d.Set("consent", scope.Consent)
	}
	return nil
}

func resourceAuthServerScopeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	scope := buildAuthServerScope(d)
	_, _, err := getOktaClientFromMetadata(m).AuthorizationServer.UpdateOAuth2Scope(ctx, d.Get("auth_server_id").(string), d.Id(), scope)
	if err != nil {
		return diag.Errorf("failed to update auth server scope: %v", err)
	}
	return resourceAuthServerScopeRead(ctx, d, m)
}

func resourceAuthServerScopeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resp, err := getOktaClientFromMetadata(m).AuthorizationServer.DeleteOAuth2Scope(ctx, d.Get("auth_server_id").(string), d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to delete auth server scope: %v", err)
	}
	return nil
}

func buildAuthServerScope(d *schema.ResourceData) okta.OAuth2Scope {
	return okta.OAuth2Scope{
		Consent:         d.Get("consent").(string),
		Description:     d.Get("description").(string),
		MetadataPublish: d.Get("metadata_publish").(string),
		Name:            d.Get("name").(string),
		Default:         boolPtr(d.Get("default").(bool)),
	}
}

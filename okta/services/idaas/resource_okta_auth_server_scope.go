package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v4/okta"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

func resourceAuthServerScope() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAuthServerScopeCreate,
		ReadContext:   resourceAuthServerScopeRead,
		UpdateContext: resourceAuthServerScopeUpdate,
		DeleteContext: resourceAuthServerScopeDelete,
		Importer:      utils.CreateNestedResourceImporter([]string{"auth_server_id", "id"}),
		Description:   "Creates an Authorization Server Scope. This resource allows you to create and configure an Authorization Server Scope.",
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
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the Auth Server Scope.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the end user displayed in a consent dialog box",
			},
			"consent": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "IMPLICIT",
				Description: "Indicates whether a consent dialog is needed for the scope. It can be set to `REQUIRED` or `IMPLICIT`. Default: `IMPLICIT`",
			},
			"metadata_publish": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "ALL_CLIENTS",
				Description: "Whether to publish metadata or not. It can be set to `ALL_CLIENTS` or `NO_CLIENTS`. Default: `ALL_CLIENTS`",
			},
			"default": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "A default scope will be returned in an access token when the client omits the scope parameter in a token request, provided this scope is allowed as part of the access policy rule.",
			},
			"system": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether Okta created the Scope",
			},
			"optional": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether the scope optional",
			},
		},
	}
}

func resourceAuthServerScopeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	scope, err := buildAuthServerScope(d)
	if err != nil {
		return diag.Errorf("failed to build auth server scope: %v", err)
	}
	respScope, _, err := getOktaV3ClientFromMetadata(meta).AuthorizationServerAPI.CreateOAuth2Scope(ctx, d.Get("auth_server_id").(string)).OAuth2Scope(scope).Execute()
	if err != nil {
		return diag.Errorf("failed to create auth server scope: %v", err)
	}
	d.SetId(respScope.GetId())
	return resourceAuthServerScopeRead(ctx, d, meta)
}

func resourceAuthServerScopeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	scope, resp, err := getOktaV3ClientFromMetadata(meta).AuthorizationServerAPI.GetOAuth2Scope(ctx, d.Get("auth_server_id").(string), d.Id()).Execute()
	if err := utils.SuppressErrorOn404_V3(resp, err); err != nil {
		return diag.Errorf("failed to get auth server scope: %v", err)
	}
	if scope == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("name", scope.GetName())
	_ = d.Set("description", scope.GetDescription())
	_ = d.Set("display_name", scope.GetDisplayName())
	_ = d.Set("metadata_publish", scope.GetMetadataPublish())
	_ = d.Set("default", scope.GetDefault())
	_ = d.Set("system", scope.GetSystem())
	_ = d.Set("optional", scope.GetOptional())
	if scope.GetConsent() != "" {
		_ = d.Set("consent", scope.GetConsent())
	}
	return nil
}

func resourceAuthServerScopeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	scope, err := buildAuthServerScope(d)
	if err != nil {
		return diag.Errorf("failed to build auth server scope: %v", err)
	}
	_, _, err = getOktaV3ClientFromMetadata(meta).AuthorizationServerAPI.ReplaceOAuth2Scope(ctx, d.Get("auth_server_id").(string), d.Id()).OAuth2Scope(scope).Execute()
	if err != nil {
		return diag.Errorf("failed to update auth server scope: %v", err)
	}
	return resourceAuthServerScopeRead(ctx, d, meta)
}

func resourceAuthServerScopeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	resp, err := getOktaV3ClientFromMetadata(meta).AuthorizationServerAPI.DeleteOAuth2Scope(ctx, d.Get("auth_server_id").(string), d.Id()).Execute()
	if err := utils.SuppressErrorOn404_V3(resp, err); err != nil {
		return diag.Errorf("failed to delete auth server scope: %v", err)
	}
	return nil
}

func buildAuthServerScope(d *schema.ResourceData) (okta.OAuth2Scope, error) {
	scope := okta.OAuth2Scope{}
	if consent, ok := d.GetOk("consent"); ok {
		scope.SetConsent(consent.(string))
	}
	scope.SetDescription(d.Get("description").(string))
	if metadataPublish, ok := d.GetOk("metadata_publish"); ok {
		scope.SetMetadataPublish(metadataPublish.(string))
	}
	scope.SetName(d.Get("name").(string))
	scope.SetDisplayName(d.Get("display_name").(string))
	scope.SetDefault(d.Get("default").(bool))
	scope.SetOptional(d.Get("optional").(bool))
	return scope, nil
}

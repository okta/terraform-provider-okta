package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v3/okta"
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
			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the end user displayed in a consent dialog box",
			},
			"consent": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "IMPLICIT",
				Description: "EA Feature and thus it is simply ignored if the feature is off",
			},
			"metadata_publish": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "ALL_CLIENTS",
				Description: "Whether to publish metadata or not, matching API type despite the fact it could just be a boolean",
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

func resourceAuthServerScopeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	scope, err := buildAuthServerScope(d)
	if err != nil {
		return diag.Errorf("failed to build auth server scope: %v", err)
	}
	respScope, _, err := getOktaV3ClientFromMetadata(m).AuthorizationServerApi.CreateOAuth2Scope(ctx, d.Get("auth_server_id").(string)).OAuth2Scope(scope).Execute()
	if err != nil {
		return diag.Errorf("failed to create auth server scope: %v", err)
	}
	d.SetId(respScope.GetId())
	return resourceAuthServerScopeRead(ctx, d, m)
}

func resourceAuthServerScopeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	scope, resp, err := getOktaV3ClientFromMetadata(m).AuthorizationServerApi.GetOAuth2Scope(ctx, d.Get("auth_server_id").(string), d.Id()).Execute()
	if err := v3suppressErrorOn404(resp, err); err != nil {
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

func resourceAuthServerScopeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	scope, err := buildAuthServerScope(d)
	if err != nil {
		return diag.Errorf("failed to build auth server scope: %v", err)
	}
	_, _, err = getOktaV3ClientFromMetadata(m).AuthorizationServerApi.ReplaceOAuth2Scope(ctx, d.Get("auth_server_id").(string), d.Id()).OAuth2Scope(scope).Execute()
	if err != nil {
		return diag.Errorf("failed to update auth server scope: %v", err)
	}
	return resourceAuthServerScopeRead(ctx, d, m)
}

func resourceAuthServerScopeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resp, err := getOktaV3ClientFromMetadata(m).AuthorizationServerApi.DeleteOAuth2Scope(ctx, d.Get("auth_server_id").(string), d.Id()).Execute()
	if err := v3suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to delete auth server scope: %v", err)
	}
	return nil
}

func buildAuthServerScope(d *schema.ResourceData) (okta.OAuth2Scope, error) {
	scope := okta.OAuth2Scope{}
	if consent, ok := d.GetOk("consent"); ok {
		consentStr, err := okta.NewOAuth2ScopeConsentTypeFromValue(consent.(string))
		if err != nil {
			return okta.OAuth2Scope{}, err
		}
		scope.SetConsent(*consentStr)
	}
	scope.SetDescription(d.Get("description").(string))
	if metadataPublish, ok := d.GetOk("metadata_publish"); ok {
		metadataPublishStr, err := okta.NewOAuth2ScopeMetadataPublishFromValue(metadataPublish.(string))
		if err != nil {
			return okta.OAuth2Scope{}, err
		}
		scope.SetMetadataPublish(*metadataPublishStr)
	}
	scope.SetName(d.Get("name").(string))
	scope.SetDisplayName(d.Get("display_name").(string))
	scope.SetDefault(d.Get("default").(bool))
	scope.SetOptional(d.Get("optional").(bool))
	return scope, nil
}

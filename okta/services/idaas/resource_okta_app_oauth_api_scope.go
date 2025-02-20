package idaas

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
)

func ResourceAppOAuthAPIScope() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppOAuthAPIScopeCreate,
		ReadContext:   resourceAppOAuthAPIScopeRead,
		UpdateContext: resourceAppOAuthAPIScopeUpdate,
		DeleteContext: resourceAppOAuthAPIScopeDelete,
		Description: `Manages API scopes for OAuth applications. 
This resource allows you to grant or revoke API scopes for OAuth2 applications within your organization.
Note: you have to create an application before using this resource.`,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				scopes, _, err := GetOktaClientFromMetadata(meta).Application.ListScopeConsentGrants(ctx, d.Id(), nil)
				if err != nil {
					return nil, err
				}
				_ = d.Set("app_id", d.Id())
				if len(scopes) > 0 {
					// Assume issuer is the same for all granted scopes, taking the first
					_ = d.Set("issuer", scopes[0].Issuer)
				} else {
					return nil, errors.New("no application scope found")
				}
				err = setOAuthApiScopes(d, scopes)
				if err != nil {
					return nil, err
				}
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"app_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ID of the application.",
				ForceNew:    true,
			},
			"issuer": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The issuer of your Org Authorization Server, your Org URL.",
			},
			"scopes": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Scopes of the application for which consent is granted.",
			},
		},
	}
}

func resourceAppOAuthAPIScopeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	scopes := utils.ConvertInterfaceToStringSetNullable(d.Get("scopes"))
	grantScopeList := getOAuthApiScopeList(scopes, d.Get("issuer").(string))
	err := grantOAuthApiScopes(ctx, d, meta, grantScopeList)
	if err != nil {
		return diag.Errorf("failed to create application scope consent grant: %v", err)
	}
	return resourceAppOAuthAPIScopeRead(ctx, d, meta)
}

func resourceAppOAuthAPIScopeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	scopes, _, err := GetOktaClientFromMetadata(meta).Application.ListScopeConsentGrants(ctx, d.Get("app_id").(string), nil)
	if err != nil {
		return diag.Errorf("failed to get application scope consent grants: %v", err)
	}

	if scopes == nil {
		d.SetId("")
		return nil
	}

	err = setOAuthApiScopes(d, scopes)
	if err != nil {
		return diag.Errorf("failed to set application scope consent grant: %v", err)
	}

	return nil
}

func resourceAppOAuthAPIScopeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	scopes, _, err := GetOktaClientFromMetadata(meta).Application.ListScopeConsentGrants(ctx, d.Get("app_id").(string), nil)
	if err != nil {
		return diag.Errorf("failed to get application scope consent grants: %v", err)
	}

	grantList, revokeList := getOAuthApiScopeUpdateLists(d, scopes)
	grantScopeList := getOAuthApiScopeList(grantList, d.Get("issuer").(string))
	err = grantOAuthApiScopes(ctx, d, meta, grantScopeList)
	if err != nil {
		return diag.Errorf("failed to create application scope consent grant: %v", err)
	}

	scopeMap, err := getOAuthApiScopeIdMap(ctx, d, meta)
	if err != nil {
		return diag.Errorf("failed to get application scope consent grant: %v", err)
	}

	revokeListIds := make([]string, 0)
	for _, scope := range revokeList {
		revokeListIds = append(revokeListIds, scopeMap[scope])
	}
	err = revokeOAuthApiScope(ctx, d, meta, revokeListIds)
	if err != nil {
		return diag.Errorf("failed to revoke application scope consent grant: %v", err)
	}

	return resourceAppOAuthAPIScopeRead(ctx, d, meta)
}

func resourceAppOAuthAPIScopeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	scopeMap, err := getOAuthApiScopeIdMap(ctx, d, meta)
	if err != nil {
		return diag.Errorf("failed to get application scope consent grant: %v", err)
	}
	revokeListIds := make([]string, 0)
	scopes := utils.ConvertInterfaceToStringSetNullable(d.Get("scopes"))
	for _, scope := range scopes {
		revokeListIds = append(revokeListIds, scopeMap[scope])
	}
	err = revokeOAuthApiScope(ctx, d, meta, revokeListIds)
	if err != nil {
		return diag.Errorf("failed to revoke application scope consent grant: %v", err)
	}

	return nil
}

// Resource Helpers
// Creates a new OAuth2ScopeConsentGrant struct
func newOAuthApiScope(scopeId, issuer string) *sdk.OAuth2ScopeConsentGrant {
	return &sdk.OAuth2ScopeConsentGrant{
		Issuer:  issuer,
		ScopeId: scopeId,
	}
}

// Creates a list of OAuth2ScopeConsentGrant structs from a string list with scope names
func getOAuthApiScopeList(scopeIds []string, issuer string) []*sdk.OAuth2ScopeConsentGrant {
	result := make([]*sdk.OAuth2ScopeConsentGrant, len(scopeIds))
	for i, scopeId := range scopeIds {
		result[i] = newOAuthApiScope(scopeId, issuer)
	}
	return result
}

// Fetches current granted application scopes and returns a map with names and IDs.
func getOAuthApiScopeIdMap(ctx context.Context, d *schema.ResourceData, meta interface{}) (map[string]string, error) {
	result := make(map[string]string)
	currentScopes, resp, err := GetOktaClientFromMetadata(meta).Application.ListScopeConsentGrants(ctx, d.Get("app_id").(string), nil)
	if err := utils.SuppressErrorOn404(resp, err); err != nil {
		return nil, fmt.Errorf("failed to get application scope consent grants: %v", err)
	}
	for _, currentScope := range currentScopes {
		result[currentScope.ScopeId] = currentScope.Id
	}
	return result, nil
}

// set resource schema from a list scopes
func setOAuthApiScopes(d *schema.ResourceData, to []*sdk.OAuth2ScopeConsentGrant) error {
	scopes := make([]string, len(to))
	for i, scope := range to {
		scopes[i] = scope.ScopeId
	}
	d.SetId(d.Get("app_id").(string))
	_ = d.Set("scopes", utils.ConvertStringSliceToSet(scopes))
	return nil
}

// Grant a list of scopes to an OAuth application. For convenience this function takes a list of OAuth2ScopeConsentGrant structs.
func grantOAuthApiScopes(ctx context.Context, d *schema.ResourceData, meta interface{}, scopeGrants []*sdk.OAuth2ScopeConsentGrant) error {
	for _, scopeGrant := range scopeGrants {
		_, _, err := GetOktaClientFromMetadata(meta).Application.GrantConsentToScope(ctx, d.Get("app_id").(string), *scopeGrant)
		if err != nil {
			return fmt.Errorf("failed to grant application api scope: %v", err)
		}
	}
	return nil
}

// Revoke a list of scopes from an OAuth application. The scope ID is needed for a revoke.
func revokeOAuthApiScope(ctx context.Context, d *schema.ResourceData, meta interface{}, ids []string) error {
	for _, id := range ids {
		resp, err := GetOktaClientFromMetadata(meta).Application.RevokeScopeConsentGrant(ctx, d.Get("app_id").(string), id)
		if err := utils.SuppressErrorOn404(resp, err); err != nil {
			return fmt.Errorf("failed to revoke application api scope: %v", err)
		}
	}
	return nil
}

// Diff function to identify which scope needs to be added or removed to the application
func getOAuthApiScopeUpdateLists(d *schema.ResourceData, from []*sdk.OAuth2ScopeConsentGrant) (grantList, revokeList []string) {
	desiredScopes := make([]string, 0)
	currentScopes := make([]string, 0)

	scopes := utils.ConvertInterfaceToStringSetNullable(d.Get("scopes"))
	desiredScopes = append(desiredScopes, scopes...)

	// extract scope list form []okta.OAuth2ScopeConsentGrant
	for _, currentScope := range from {
		currentScopes = append(currentScopes, currentScope.ScopeId)
	}

	// return scopes that should be added or removed
	return splitTargets(desiredScopes, currentScopes)
}

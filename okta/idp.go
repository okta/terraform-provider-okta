package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

const (
	postBindingAlias     = "HTTP-POST"
	redirectBindingAlias = "HTTP-REDIRECT"
)

var (
	baseIdpSchema = map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Name of the IdP",
		},
		"status": statusSchema,
		"account_link_action": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "AUTO",
		},
		"account_link_group_include": {
			Type:     schema.TypeSet,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
		},
		"provisioning_action": {
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: elemInSlice([]string{"AUTO", "DISABLED", ""}),
			Default:          "AUTO",
		},
		"deprovisioned_action": actionSchema,
		"suspended_action":     actionSchema,
		"groups_action": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          "NONE",
			ValidateDiagFunc: elemInSlice([]string{"NONE", "SYNC", "APPEND", "ASSIGN"}),
		},
		"groups_attribute": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"groups_assignment": {
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeSet,
		},
		"groups_filter": {
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeSet,
		},
		"username_template": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "idpuser.email",
		},
		"subject_match_type": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          "USERNAME",
			ValidateDiagFunc: elemInSlice([]string{"USERNAME", "EMAIL", "USERNAME_OR_EMAIL", "CUSTOM_ATTRIBUTE"}),
		},
		"subject_match_attribute": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"profile_master": {
			Type:     schema.TypeBool,
			Optional: true,
		},
	}

	actionSchema = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Default:  "NONE",
	}

	samlRequestSignatureAlgorithmSchema = &schema.Schema{
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "The XML digital Signature Algorithm used when signing an <AuthnRequest> message",
		ValidateDiagFunc: elemInSlice([]string{"SHA-256", "SHA-1"}),
		Default:          "SHA-256",
	}
	samlRequestSignatureScopeSchema = &schema.Schema{
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "Specifies whether to digitally sign <AuthnRequest> messages to the IdP",
		ValidateDiagFunc: elemInSlice([]string{"REQUEST", "NONE"}),
		Default:          "REQUEST",
	}

	samlResponseSignatureAlgorithmSchema = &schema.Schema{
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "The minimum XML digital Signature Algorithm allowed when verifying a <SAMLResponse> message or <Assertion> element",
		ValidateDiagFunc: elemInSlice([]string{"SHA-256", "SHA-1"}),
		Default:          "SHA-256",
	}
	samlResponseSignatureScopeSchema = &schema.Schema{
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "Specifies whether to verify a <SAMLResponse> message or <Assertion> element XML digital signature",
		ValidateDiagFunc: elemInSlice([]string{"RESPONSE", "ASSERTION", "ANY"}),
		Default:          "ANY",
	}

	oidcRequestSignatureAlgorithmSchema = &schema.Schema{
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "The HMAC Signature Algorithm used when signing an authorization request",
		ValidateDiagFunc: elemInSlice([]string{"HS256", "HS384", "HS512", "SHA-256", "RS256", "RS384", "RS512"}),
		Default:          "HS256",
	}

	oidcRequestSignatureScopeSchema = &schema.Schema{
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "Specifies whether to digitally sign an authorization request to the IdP",
		ValidateDiagFunc: elemInSlice([]string{"REQUEST", "NONE"}),
		Default:          "REQUEST",
	}

	optBindingSchema = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}

	optURLSchema = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}

	optionalURLSchema = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}

	bindingSchema = &schema.Schema{
		Type:             schema.TypeString,
		Required:         true,
		ValidateDiagFunc: elemInSlice([]string{"HTTP-POST", "HTTP-REDIRECT"}),
	}

	optionalBindingSchema = &schema.Schema{
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: elemInSlice([]string{"HTTP-POST", "HTTP-REDIRECT"}),
	}

	issuerMode = &schema.Schema{
		Type:             schema.TypeString,
		Description:      "Indicates whether Okta uses the original Okta org domain URL, or a custom domain URL",
		ValidateDiagFunc: elemInSlice([]string{"ORG_URL", "CUSTOM_URL_DOMAIN"}),
		Default:          "ORG_URL",
		Optional:         true,
	}

	urlSchema = &schema.Schema{
		Type:             schema.TypeString,
		Required:         true,
		ValidateDiagFunc: stringIsURL(validURLSchemes...),
	}
)

func buildIdpSchema(idpSchema map[string]*schema.Schema) map[string]*schema.Schema {
	return buildSchema(baseIdpSchema, idpSchema)
}

func getIdentityProviderByID(ctx context.Context, m interface{}, id, providerType string) (*okta.IdentityProvider, error) {
	idp, _, err := getOktaClientFromMetadata(m).IdentityProvider.GetIdentityProvider(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get identity provider with id '%s': %v", id, err)
	}
	if idp.Type != providerType {
		return nil, fmt.Errorf("identity provider of type '%s' with id '%s' does not exist", providerType, id)
	}
	return idp, nil
}

func getIdpByNameAndType(ctx context.Context, m interface{}, name, providerType string) (*okta.IdentityProvider, error) {
	queryParams := &query.Params{Limit: 1, Q: name, Type: providerType}
	idps, _, err := getOktaClientFromMetadata(m).IdentityProvider.ListIdentityProviders(ctx, queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to find identity provider '%s': %v", name, err)
	}
	if len(idps) < 1 || idps[0].Name != name {
		return nil, fmt.Errorf("identity provider with name '%s' and type '%s' does not exist: %v", name, providerType, err)
	}
	return idps[0], nil
}

func resourceIdpDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(m)
	_, resp, err := client.IdentityProvider.DeactivateIdentityProvider(ctx, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to deactivate identity provider: %v", err)
	}
	resp, err = client.IdentityProvider.DeleteIdentityProvider(ctx, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to delete identity provider: %v", err)
	}
	return nil
}

func setIdpStatus(ctx context.Context, d *schema.ResourceData, client *okta.Client, status string) error {
	desiredStatus := d.Get("status").(string)
	if status == desiredStatus {
		return nil
	}
	var err error
	if desiredStatus == statusInactive {
		_, _, err = client.IdentityProvider.DeactivateIdentityProvider(ctx, d.Id())
	} else {
		_, _, err = client.IdentityProvider.ActivateIdentityProvider(ctx, d.Id())
	}
	return err
}

func syncEndpoint(key string, e *okta.ProtocolEndpoint, d *schema.ResourceData) {
	if e != nil {
		_ = d.Set(key+"_binding", e.Binding)
		_ = d.Set(key+"_url", e.Url)
	}
}

func syncGroupActions(d *schema.ResourceData, groups *okta.ProvisioningGroups) error {
	if groups == nil {
		return nil
	}
	_ = d.Set("groups_action", groups.Action)
	_ = d.Set("groups_attribute", groups.SourceAttributeName)
	return setNonPrimitives(d, map[string]interface{}{
		"groups_assignment": groups.Assignments,
		"groups_filter":     groups.Filter,
	})
}

func buildPolicyAccountLink(d *schema.ResourceData) *okta.PolicyAccountLink {
	link := convertInterfaceToStringSet(d.Get("account_link_group_include"))
	var filter *okta.PolicyAccountLinkFilter

	if len(link) > 0 {
		filter = &okta.PolicyAccountLinkFilter{
			Groups: &okta.PolicyAccountLinkFilterGroups{
				Include: link,
			},
		}
	}
	return &okta.PolicyAccountLink{
		Action: d.Get("account_link_action").(string),
		Filter: filter,
	}
}

func buildIdPProvisioning(d *schema.ResourceData) *okta.Provisioning {
	return &okta.Provisioning{
		Action:        d.Get("provisioning_action").(string),
		ProfileMaster: boolPtr(d.Get("profile_master").(bool)),
		Conditions: &okta.ProvisioningConditions{
			Deprovisioned: &okta.ProvisioningDeprovisionedCondition{
				Action: d.Get("deprovisioned_action").(string),
			},
			Suspended: &okta.ProvisioningSuspendedCondition{
				Action: d.Get("suspended_action").(string),
			},
		},
		Groups: &okta.ProvisioningGroups{
			Action:              d.Get("groups_action").(string),
			Assignments:         convertInterfaceToStringSetNullable(d.Get("groups_assignment")),
			Filter:              convertInterfaceToStringSetNullable(d.Get("groups_filter")),
			SourceAttributeName: d.Get("groups_attribute").(string),
		},
	}
}

func buildAlgorithms(d *schema.ResourceData) *okta.ProtocolAlgorithms {
	return &okta.ProtocolAlgorithms{
		Request:  buildProtocolAlgorithmType(d, "request"),
		Response: buildProtocolAlgorithmType(d, "response"),
	}
}

func buildProtocolAlgorithmType(d *schema.ResourceData, key string) *okta.ProtocolAlgorithmType {
	scopeKey := fmt.Sprintf("%s_signature_scope", key)
	scope, ok := d.GetOk(scopeKey)
	if !ok || scope.(string) == "" {
		return nil
	}
	return &okta.ProtocolAlgorithmType{
		Signature: &okta.ProtocolAlgorithmTypeSignature{
			Algorithm: d.Get(fmt.Sprintf("%s_signature_algorithm", key)).(string),
			Scope:     scope.(string),
		},
	}
}

func buildProtocolEndpoints(d *schema.ResourceData) *okta.ProtocolEndpoints {
	return &okta.ProtocolEndpoints{
		Authorization: buildProtocolEndpoint(d, "authorization"),
		Token:         buildProtocolEndpoint(d, "token"),
		UserInfo:      buildProtocolEndpoint(d, "user_info"),
		Jwks:          buildProtocolEndpoint(d, "jwks"),
	}
}

func buildProtocolEndpoint(d *schema.ResourceData, key string) *okta.ProtocolEndpoint {
	binding := d.Get(fmt.Sprintf("%s_binding", key)).(string)
	url := d.Get(fmt.Sprintf("%s_url", key)).(string)
	if binding != "" && url != "" {
		return &okta.ProtocolEndpoint{
			Binding: binding,
			Url:     url,
		}
	}
	return nil
}

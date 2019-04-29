package okta

import (
	"fmt"
	articulateOkta "github.com/articulate/oktasdk-go/okta"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"net/http"
)

// DEPRECATED - see okta_idp and okta_saml_idp
func resourceIdentityProvider() *schema.Resource {
	return &schema.Resource{
		Create: resourceIdentityProviderCreate,
		Read:   resourceIdentityProviderRead,
		Update: resourceIdentityProviderUpdate,
		Delete: resourceIdentityProviderDelete,
		Exists: idpExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		DeprecationMessage: "This resource is being deprecated in favor of okta_oidc_idp & okta_saml_idp",
		Schema: map[string]*schema.Schema{
			"active": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the IDP is active or not - can only be issued post-creation",
			},
			"authorization_url": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"authorization_url_binding": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"client_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "OAUTH2 client ID",
			},
			"client_secret": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "OAUTH2 client secret",
				Sensitive:   true,
			},
			"issuer_mode": &schema.Schema{
				Type:         schema.TypeString,
				Description:  "Indicates whether Okta uses the original Okta org domain URL, or a custom domain URL",
				ValidateFunc: validation.StringInSlice([]string{"ORG_URL", "CUSTOM_URL_DOMAIN"}, false),
				Default:      "ORG_URL",
				Optional:     true,
			},
			"links_authorized_hints_allow": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"links_authorized_href": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"links_authorized_templated": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"links_client_redirect_uri_hints_allow": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"links_client_redirect_uri_href": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Identity Provider Resource",
			},
			"policy_account_link_action": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "AUTO",
				Description: "Policy Account Link Action",
			},
			"policy_account_link_filter": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Policy Account Link Filter",
			},
			"policy_max_clock_skew": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Policy Max Clock Skew",
			},
			"policy_provisioning_action": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "AUTO",
				Description: "Policy Provisioning Action",
			},
			"policy_provisioning_conditions_deprovisioned_action": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "NONE",
				Description: "Policy Provisioning Conditions Deprovisioned Action",
			},
			"policy_provisioning_conditions_suspended_action": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "NONE",
				Description: "Policy Provisioning Conditions Suspended Action",
			},
			"policy_provisioning_group_assignments": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Policy Provisioning Groups Assignment",
			},
			"policy_provisioning_groups_action": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Policy Provisioning Groups Action",
			},
			"policy_provisioning_profile_master": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Policy Provisioning Profile Master",
			},
			"policy_subject_filter": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Policy Subject Filter",
			},
			"policy_subject_match_type": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "USERNAME",
				Description: "Policy Subject Match Type",
			},
			"policy_subject_username_template": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "idpuser.firstName",
				Description: "Policy Subject Username Template",
			},
			"protocol_scopes": &schema.Schema{
				Type:        schema.TypeList,
				MinItems:    1,
				Required:    true,
				Description: "Scopes provided to the Idp, e.g. 'openid', 'email', 'profile'",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"protocol_type": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "OIDC",
				Description: "IDP Protocol type to use - ie. OAUTH2",
			},
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"GOOGLE"}, false),
				Description:  "Identity Provider Type: GOOGLE",
			},
		},
	}
}

func activationStatus(active bool) string {
	if active {
		return "ACTIVE"
	}

	return "INACTIVE"
}

func assembleIdentityProvider() *articulateOkta.IdentityProvider {
	idp := &articulateOkta.IdentityProvider{}

	client := &articulateOkta.IdpClient{}
	credentials := &articulateOkta.Credentials{Client: client}
	authorization := &articulateOkta.Authorization{}
	endpoints := &articulateOkta.Endpoints{Authorization: authorization}
	protocol := &articulateOkta.Protocol{
		Credentials: credentials,
		Endpoints:   endpoints,
	}
	idpGroups := &articulateOkta.IdpGroups{}
	deprovisioned := &articulateOkta.Deprovisioned{}
	suspended := &articulateOkta.Suspended{}

	conditions := &articulateOkta.Conditions{
		Deprovisioned: deprovisioned,
		Suspended:     suspended,
	}

	provisioning := &articulateOkta.Provisioning{
		ProfileMaster: true,
		Groups:        idpGroups,
		Conditions:    conditions,
	}

	accountLink := &articulateOkta.AccountLink{}

	userNameTemplate := &articulateOkta.UserNameTemplate{}

	subject := &articulateOkta.Subject{
		UserNameTemplate: userNameTemplate,
	}

	policy := &articulateOkta.IdpPolicy{
		Provisioning: provisioning,
		AccountLink:  accountLink,
		Subject:      subject,
		MaxClockSkew: 0,
	}

	authorize := &articulateOkta.Authorize{
		Hints: &articulateOkta.Hints{},
	}

	clientRedirectUri := &articulateOkta.ClientRedirectUri{
		Hints: &articulateOkta.Hints{},
	}

	idpLinks := &articulateOkta.IdpLinks{
		Authorize:         authorize,
		ClientRedirectUri: clientRedirectUri,
	}

	idp.Protocol = protocol
	idp.Policy = policy
	idp.Links = idpLinks

	return idp
}

// Populates the IdentityProvider struct (used by the Okta SDK for API operaations) with the state provided by TF
func populateIdentityProvider(idp *articulateOkta.IdentityProvider, d *schema.ResourceData) *articulateOkta.IdentityProvider {

	idp.Type = d.Get("type").(string)
	idp.Name = d.Get("name").(string)
	idp.Protocol.Endpoints.Authorization.Url = d.Get("authorization_url").(string)
	idp.Protocol.Endpoints.Authorization.Binding = d.Get("authorization_url_binding").(string)
	idp.Protocol.Type = d.Get("protocol_type").(string)

	scopes := make([]string, 0)
	for _, vals := range d.Get("protocol_scopes").([]interface{}) {
		scopes = append(scopes, vals.(string))
	}
	idp.Protocol.Scopes = scopes
	idp.Policy.Provisioning.Action = d.Get("policy_provisioning_action").(string)
	idp.Policy.Provisioning.ProfileMaster = d.Get("policy_provisioning_profile_master").(bool)

	if vals, ok := d.GetOk("policy_provisioning_group_assignments"); ok {
		groupAssignments := make([]string, 0)
		for _, val := range vals.([]interface{}) {
			groupAssignments = append(groupAssignments, val.(string))
		}
		idp.Policy.Provisioning.Groups.Action = "ASSIGN"
		idp.Policy.Provisioning.Groups.Assignments = groupAssignments
	} else {
		idp.Policy.Provisioning.Groups.Action = "NONE"
	}

	idp.Policy.Provisioning.Conditions.Deprovisioned.Action = d.Get("policy_provisioning_conditions_deprovisioned_action").(string)
	idp.Policy.Provisioning.Conditions.Suspended.Action = d.Get("policy_provisioning_conditions_suspended_action").(string)
	idp.Policy.AccountLink.Filter = d.Get("policy_account_link_filter").(string)
	idp.Policy.AccountLink.Action = d.Get("policy_account_link_action").(string)
	idp.Policy.Subject.UserNameTemplate.Template = d.Get("policy_subject_username_template").(string)
	idp.Policy.Subject.Filter = d.Get("policy_subject_filter").(string)
	idp.Policy.Subject.MatchType = d.Get("policy_subject_match_type").(string)
	idp.Policy.MaxClockSkew = d.Get("policy_max_clock_skew").(int)
	idp.Links.Authorize.Href = d.Get("links_authorized_href").(string)
	idp.Links.Authorize.Templated = d.Get("links_authorized_templated").(bool)
	idp.Links.ClientRedirectUri.Href = d.Get("links_client_redirect_uri_href").(string)
	idp.Protocol.Credentials.Client.ClientID = d.Get("client_id").(string)
	idp.Protocol.Credentials.Client.ClientSecret = d.Get("client_secret").(string)
	idp.IssuerMode = d.Get("issuer_mode").(string)

	return idp
}

func resourceIdentityProviderCreate(d *schema.ResourceData, m interface{}) error {
	client := getClientFromMetadata(m)
	idp := assembleIdentityProvider()
	populateIdentityProvider(idp, d)

	returnedIdp, _, err := client.IdentityProviders.CreateIdentityProvider(idp)
	if err != nil {
		return err
	}

	d.SetId(returnedIdp.ID)

	return resourceIdentityProviderRead(d, m)
}

func resourceIdentityProviderRead(d *schema.ResourceData, m interface{}) error {
	var idp *articulateOkta.IdentityProvider
	client := getClientFromMetadata(m)
	idp, _, err := client.IdentityProviders.GetIdentityProvider(d.Id())
	if err != nil {
		return err
	}

	d.Set("type", idp.Type)
	d.Set("name", idp.Name)
	d.Set("active", idp.Status == "ACTIVE")
	d.Set("authorization_url", idp.Protocol.Endpoints.Authorization.Url)
	d.Set("authorization_url_binding", idp.Protocol.Endpoints.Authorization.Binding)
	d.Set("protocol_type", idp.Protocol.Type)
	d.Set("policy_provisioning_action", idp.Policy.Provisioning.Action)
	d.Set("policy_provisioning_profile_master", idp.Policy.Provisioning.ProfileMaster)
	d.Set("policy_provisioning_groups_action", idp.Policy.Provisioning.Groups.Action)
	d.Set("policy_provisioning_conditions_deprovisioned_action", idp.Policy.Provisioning.Conditions.Deprovisioned.Action)
	d.Set("policy_provisioning_conditions_suspended_action", idp.Policy.Provisioning.Conditions.Suspended.Action)
	d.Set("policy_account_link_filter", idp.Policy.AccountLink.Filter)
	d.Set("policy_account_link_action", idp.Policy.AccountLink.Action)
	d.Set("policy_subject_username_template", idp.Policy.Subject.UserNameTemplate.Template)
	d.Set("policy_subject_filter", idp.Policy.Subject.Filter)
	d.Set("policy_subject_match_type", idp.Policy.Subject.MatchType)
	d.Set("policy_max_clock_skew", idp.Policy.MaxClockSkew)
	d.Set("links_authorized_href", idp.Links.Authorize.Href)
	d.Set("links_authorized_templated", idp.Links.Authorize.Templated)
	d.Set("links_client_redirect_uri_href", idp.Links.ClientRedirectUri.Href)
	d.Set("client_id", idp.Protocol.Credentials.Client.ClientID)
	d.Set("client_secret", idp.Protocol.Credentials.Client.ClientSecret)

	if idp.IssuerMode != "" {
		d.Set("issuer_mode", idp.IssuerMode)
	}

	agTypeMap := map[string]interface{}{
		"protocol_scopes":                       idp.Protocol.Scopes,
		"links_authorized_hints_allow":          idp.Links.Authorize.Hints.Allow,
		"links_client_redirect_uri_hints_allow": idp.Links.ClientRedirectUri.Hints.Allow,
	}
	assignmentList := idp.Policy.Provisioning.Groups.Assignments

	if len(assignmentList) > 0 {
		agTypeMap["policy_provisioning_group_assignments"] = assignmentList
	}

	return setNonPrimitives(d, agTypeMap)
}

func resourceIdentityProviderUpdate(d *schema.ResourceData, m interface{}) error {
	var idp = assembleIdentityProvider()
	populateIdentityProvider(idp, d)

	// can only update IDP status in Update operation
	status := activationStatus(d.Get("active").(bool))
	idp.Status = status
	client := getClientFromMetadata(m)
	idp, _, err := client.IdentityProviders.UpdateIdentityProvider(d.Id(), idp)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Updating Identity Provider with Okta: %v", err)
	}

	return resourceIdentityProviderRead(d, m)
}

// check if IDP exists in Okta
func idpExists(d *schema.ResourceData, m interface{}) (bool, error) {
	_, resp, err := getClientFromMetadata(m).IdentityProviders.GetIdentityProvider(d.Id())
	return resp.StatusCode == http.StatusOK, err
}

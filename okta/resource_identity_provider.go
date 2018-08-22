package okta

import (
	"fmt"
	"github.com/articulate/oktasdk-go/okta"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"log"
)

func resourceIdentityProviders() *schema.Resource {
	return &schema.Resource{
		Create: resourceIdentityProviderCreate,
		Read:   resourceIdentityProviderRead,
		Update: resourceIdentityProviderUpdate,
		Delete: resourceIdentityProviderDelete,
		Exists: idpExists,

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
				Sensitive: true,
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
			"policy_provisioning_groups_action": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "NONE",
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

func assembleIdentityProvider() *okta.IdentityProvider {
	idp := &okta.IdentityProvider{}

	client := &okta.IdpClient{}
	credentials := &okta.Credentials{Client: client}
	authorization := &okta.Authorization{}
	endpoints := &okta.Endpoints{Authorization: authorization}
	protocol := &okta.Protocol{
		Credentials: credentials,
		Endpoints:   endpoints,
	}
	idpGroups := &okta.IdpGroups{}
	deprovisioned := &okta.Deprovisioned{}
	suspended := &okta.Suspended{}

	conditions := &okta.Conditions{
		Deprovisioned: deprovisioned,
		Suspended:     suspended,
	}

	provisioning := &okta.Provisioning{
		ProfileMaster: true,
		Groups:        idpGroups,
		Conditions:    conditions,
	}

	accountLink := &okta.AccountLink{}

	userNameTemplate := &okta.UserNameTemplate{}

	subject := &okta.Subject{
		UserNameTemplate: userNameTemplate,
	}

	policy := &okta.IdpPolicy{
		Provisioning: provisioning,
		AccountLink:  accountLink,
		Subject:      subject,
		MaxClockSkew: 0,
	}

	authorize := &okta.Authorize{
		Hints: &okta.Hints{},
	}

	clientRedirectUri := &okta.ClientRedirectUri{
		Hints: &okta.Hints{},
	}

	idpLinks := &okta.IdpLinks{
		Authorize:         authorize,
		ClientRedirectUri: clientRedirectUri,
	}

	idp.Protocol = protocol
	idp.Policy = policy
	idp.Links = idpLinks

	return idp
}

func populateIdentityProvider(idp *okta.IdentityProvider, d *schema.ResourceData) *okta.IdentityProvider {

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
	idp.Policy.Provisioning.Groups.Action = d.Get("policy_provisioning_groups_action").(string)
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

	return idp
}

func resourceIdentityProviderCreate(d *schema.ResourceData, m interface{}) error {
	if !d.Get("active").(bool) {
		return fmt.Errorf("[ERROR] Okta will not allow an IDP to be created as INACTIVE. Can set to false for existing IDPs only.")
	}

	client := m.(*Config).oktaClient
	idp := assembleIdentityProvider()

	idp.Type = d.Get("type").(string)
	idp.Name = d.Get("name").(string)

	idp.Protocol.Type = d.Get("protocol_type").(string)

	if len(d.Get("protocol_scopes").([]interface{})) > 0 {
		scopes := make([]string, 0)
		for _, vals := range d.Get("protocol_scopes").([]interface{}) {
			scopes = append(scopes, vals.(string))
		}
		idp.Protocol.Scopes = scopes
	}

	idp.Protocol.Credentials.Client.ClientID = d.Get("client_id").(string)
	idp.Protocol.Credentials.Client.ClientSecret = d.Get("client_secret").(string)

	// Hardcode required values
	idp.Policy.Provisioning.Action = d.Get("policy_provisioning_action").(string)
	idp.Policy.Provisioning.Groups.Action = d.Get("policy_provisioning_groups_action").(string)
	idp.Policy.Provisioning.Conditions.Deprovisioned.Action = d.Get("policy_provisioning_conditions_deprovisioned_action").(string)
	idp.Policy.Provisioning.Conditions.Suspended.Action = d.Get("policy_provisioning_conditions_suspended_action").(string)
	idp.Policy.AccountLink.Action = d.Get("policy_account_link_action").(string)
	idp.Policy.Subject.UserNameTemplate.Template = d.Get("policy_subject_username_template").(string)
	idp.Policy.Subject.MatchType = d.Get("policy_subject_match_type").(string)

	returnedIdp, _, err := client.IdentityProviders.CreateIdentityProvider(*idp)

	d.SetId(returnedIdp.ID)
	if err != nil {
		return err
	}
	return resourceIdentityProviderRead(d, m)
}

func resourceIdentityProviderRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] List Identity Provider %v", d.Get("name").(string))

	var idp *okta.IdentityProvider
	client := m.(*Config).oktaClient
	exists, err := idpExists(d, m)
	if err != nil {
		return err
	}

	if exists == true {
		idp, _, err = client.IdentityProviders.GetIdentityProvider(d.Id())
	} else {
		d.SetId("")
		return nil
	}

	d.Set("type", idp.Type)
	d.Set("name", idp.Name)
	d.Set("active", idp.Status == "ACTIVE")
	d.Set("authorization_url", idp.Protocol.Endpoints.Authorization.Url)
	d.Set("authorization_url_binding", idp.Protocol.Endpoints.Authorization.Binding)
	d.Set("protocol_type", idp.Protocol.Type)
	d.Set("protocol_scopes", idp.Protocol.Scopes)
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
	d.Set("links_authorized_hints_allow", idp.Links.Authorize.Hints.Allow)
	d.Set("links_client_redirect_uri_href", idp.Links.ClientRedirectUri.Href)
	d.Set("links_client_redirect_uri_hints_allow", idp.Links.ClientRedirectUri.Hints.Allow)
	d.Set("client_id", idp.Protocol.Credentials.Client.ClientID)
	d.Set("client_secret", idp.Protocol.Credentials.Client.ClientSecret)
	return nil
}

func resourceIdentityProviderUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] List Identity Provider %v", d.Get("name").(string))

	var idp = assembleIdentityProvider()
	populateIdentityProvider(idp, d)

	// can only update IDP status in Update operation
	status := activationStatus(d.Get("active").(bool))
	idp.Status = status

	client := m.(*Config).oktaClient
	exists, err := idpExists(d, m)
	if err != nil {
		return err
	}

	if exists == true {
		idp, _, err = client.IdentityProviders.UpdateIdentityProvider(d.Id(), idp)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Updating Identity Provider with Okta: %v", err)
		}
	} else {
		d.SetId("")
		return nil
	}

	return resourceIdentityProviderRead(d, m)
}

func resourceIdentityProviderDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Delete Identity Provider %v", d.Get("name").(string))
	client := m.(*Config).oktaClient
	exists, err := idpExists(d, m)
	if err != nil {
		return err
	}
	if exists == true {
		_, err = client.IdentityProviders.DeleteIdentityProvider(d.Id())
		if err != nil {
			return fmt.Errorf("[ERROR] Error Deleting Identity Providers from Okta: %v", err)
		}
	}

	updatedExists, err2 := idpExists(d, m)
	if err2 != nil {
		return err2
	}

	if updatedExists == true {
		return fmt.Errorf("%v", "Resource Still Exits after Destroy called")
	}
	return nil
}

// check if IDP exists in Okta
func idpExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*Config).oktaClient
	_, _, err := client.IdentityProviders.GetIdentityProvider(d.Id())

	if client.OktaErrorCode == "E0000007" {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("[ERROR] Error Listing Identity Provider in Okta: %v", err)
	}
	return true, nil
}

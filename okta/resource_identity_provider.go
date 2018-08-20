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
		Create:        resourceIdentityProviderCreate,
		Read:          resourceIdentityProviderRead,
		Update:        resourceIdentityProviderUpdate,
		Delete:        resourceIdentityProviderDelete,
		CustomizeDiff: func(d *schema.ResourceDiff, v interface{}) error { return nil },

		Schema: map[string]*schema.Schema{
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"GOOGLE"}, false),
				Description:  "Identity Provider Type: GOOGLE",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Identity Provider Resource",
			},
			"protocol_type": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "OAUTH2",
				Description: "IDP Protocol type to use - ie. OAUTH2",
			},
			"protocol_scopes": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Scopes provided to the Idp, e.g. 'openid', 'email', 'profile'",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"client_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "OAUTH2 client ID",
			},
			"client_secret": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "OAUTH2 client secret",
			},
			"status": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
			},
			"authorization_url": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
			},
			"authorization_url_binding": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIdentityProviderCreate(d *schema.ResourceData, m interface{}) error {

	client := m.(*Config).oktaClient
	idp := client.IdentityProviders.IdentityProvider()

	idpClient := &okta.IdpClient{}

	credentials := &okta.Credentials{Client: idpClient}
	protocol := &okta.Protocol{Credentials: credentials}

	idpGroups := &okta.IdpGroups{Action:"NONE"}
	deprovisioned := &okta.Deprovisioned{Action:"NONE"}
	suspended := &okta.Suspended{Action:"NONE"}	
	
	conditions := &okta.Conditions{
		Deprovisioned: deprovisioned,
		Suspended: suspended,
	}

	provisioning := &okta.Provisioning{
		Action: "AUTO",
		ProfileMaster: true,
		Groups: idpGroups,
		Conditions: conditions,
	}

	accountLink := &okta.AccountLink{
		Action: "AUTO",
	}

	userNameTemplate := &okta.UserNameTemplate{
		Template: "idpuser.firstName",
	}
	
	subject := &okta.Subject{
		UserNameTemplate: userNameTemplate,
		MatchType: "USERNAME",
	}
	
	policy := &okta.IdpPolicy{
		Provisioning: provisioning,
		AccountLink: accountLink,
		Subject: subject,
		MaxClockSkew: 0,
	}

	idp.Type = d.Get("type").(string)
	idp.Name = d.Get("name").(string)

	protocol.Type = d.Get("protocol_type").(string)

	if len(d.Get("protocol_scopes").([]interface{})) > 0 {
		scopes := make([]string, 0)
		for _, vals := range d.Get("protocol_scopes").([]interface{}) {
			scopes = append(scopes, vals.(string))
		}
		protocol.Scopes = scopes
	}

	protocol.Credentials.Client.ClientID     = d.Get("client_id").(string)
	protocol.Credentials.Client.ClientSecret = d.Get("client_secret").(string)

	idp.Protocol = protocol
	idp.Policy   = policy

	returnedIdp, _, err := client.IdentityProviders.CreateIdentityProvider(idp)

	d.SetId(returnedIdp.ID);
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
		return fmt.Errorf("[ERROR] Error Identity Provider not found in Okta: %v", err)
	}

	d.Set("type", idp.Type)
	d.Set("name", idp.Name)
	d.Set("status", idp.Status)
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
	return nil
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
	} else {
		return fmt.Errorf("[ERROR] Error Identity Provider not found in Okta: %v", err)
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

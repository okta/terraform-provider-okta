package okta

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
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
			"protocol": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Conditions that must be met during Policy Evaluation",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "OAUTH2",
							Description: "IDP Protocol type to use - ie. OAUTH2",
						},
						"scopes": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Scopes provided to the Idp, e.g. 'openid', 'email', 'profile'",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"credentials": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Credentials",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"client": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"client_id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "OAUTH2 client ID",
												},
												"client_secret": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "OAUTH2 client secret",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"policy": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"provisioning": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "",
										Default:     "AUTO",
									},
									"profile_master": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "",
										Default:     true,
									},
									"groups": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"action": {
													Type:        schema.TypeString,
													Optional:    true,
													Default:     "NONE",
													Description: "",
												},
											},
										},
									},
									"conditions": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"deprovisioned": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"action": {
																Type:        schema.TypeString,
																Optional:    true,
																Default:     "NONE",
																Description: "",
															},
														},
													},
												},
												"suspended": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"action": {
																Type:        schema.TypeString,
																Optional:    true,
																Default:     "NONE",
																Description: "",
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceIdentityProviderCreate(d *schema.ResourceData, m interface{}) error {

	client := m.(*Config).oktaClient
	idp := client.IdentityProviders.IdentityProvider()

	
	idp.Type = d.Get("type").(string)
	idp.Name = d.Get("name").(string)
	idp.Protocol.Type = d.Get("protocol.type").(string)
	idp.Protocol.Scopes = d.Get("protocol.scopes").([]string)
	idp.Protocol.Credentials.Client.ClientID = d.Get("protocol.credentials.client.client_id").(string)
	idp.Protocol.Credentials.Client.ClientSecret = d.Get("protocol.credentials.client.client_secret").(string)
	idp.Policy.Provisioning.Action = d.Get("policy.provisioning.action").(string)
	idp.Policy.Provisioning.ProfileMaster = d.Get("policy.provisioning.profile_master").(bool)
	idp.Policy.Provisioning.Groups.Action = d.Get("policy.provisioning.groups.action").(string)
	idp.Policy.Provisioning.Conditions.Deprovisioned.Action = d.Get("policy.provisioning.conditions.deprovisioned.action").(string)
	idp.Policy.Provisioning.Conditions.Suspended.Action = d.Get("policy.provisioning.conditions.suspended.action").(string)
	idp.Policy.AccountLink.Filter = ""
	idp.Policy.AccountLink.Action = "AUTO"
	idp.Policy.Subject.UserNameTemplate.Template = "idpuser.userPrincipalName"
	idp.Policy.Subject.Filter = ""
	idp.Policy.Subject.MatchType = "USERNAME"
	idp.Policy.MaxClockSkew = 0

	_, _, err := client.IdentityProviders.CreateIdentityProvider(idp)
	if err != nil {
		fmt.Println("ERRORE OMG PROTECC ME!!!")
		fmt.Println(err)
		return err
	}
	return nil
}

func resourceIdentityProviderRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceIdentityProviderUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceIdentityProviderDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

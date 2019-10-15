package okta

import (
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-okta/sdk"
)

func resourceAuthServerPolicy() *schema.Resource {
	return &schema.Resource{
		Create:   resourceAuthServerPolicyCreate,
		Exists:   resourceAuthServerPolicyExists,
		Read:     resourceAuthServerPolicyRead,
		Update:   resourceAuthServerPolicyUpdate,
		Delete:   resourceAuthServerPolicyDelete,
		Importer: createNestedResourceImporter([]string{"auth_server_id", "id"}),

		Schema: map[string]*schema.Schema{
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "OAUTH_AUTHORIZATION_POLICY",
				Description: "Auth server policy type, unlikely this will be anything other then the default",
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"auth_server_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"status": statusSchema,
			"priority": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Priority of the auth server policy",
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"client_whitelist": &schema.Schema{
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Use [\"ALL_CLIENTS\"] when unsure.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func buildAuthServerPolicy(d *schema.ResourceData) *sdk.AuthorizationServerPolicy {
	return &sdk.AuthorizationServerPolicy{
		Name:        d.Get("name").(string),
		Type:        d.Get("type").(string),
		Status:      d.Get("status").(string),
		Priority:    d.Get("priority").(int),
		Description: d.Get("description").(string),
		Conditions: &sdk.PolicyConditions{
			Clients: &sdk.Whitelist{
				Include: convertInterfaceToStringSet(d.Get("client_whitelist")),
			},
		},
	}
}

func resourceAuthServerPolicyCreate(d *schema.ResourceData, m interface{}) error {
	authServerPolicy := buildAuthServerPolicy(d)
	c := getSupplementFromMetadata(m)
	responseAuthServerPolicy, _, err := c.CreateAuthorizationServerPolicy(d.Get("auth_server_id").(string), *authServerPolicy, nil)
	if err != nil {
		return err
	}

	d.SetId(responseAuthServerPolicy.Id)

	return resourceAuthServerPolicyRead(d, m)
}

func resourceAuthServerPolicyExists(d *schema.ResourceData, m interface{}) (bool, error) {
	g, err := fetchAuthServerPolicy(d, m)

	return err == nil && g != nil, err
}

func resourceAuthServerPolicyRead(d *schema.ResourceData, m interface{}) error {
	authServerPolicy, err := fetchAuthServerPolicy(d, m)

	if authServerPolicy == nil {
		d.SetId("")
		return nil
	}

	if err != nil {
		return err
	}

	d.Set("name", authServerPolicy.Name)
	d.Set("description", authServerPolicy.Description)
	d.Set("status", authServerPolicy.Status)
	d.Set("priority", authServerPolicy.Priority)
	d.Set("client_whitelist", convertStringSetToInterface(authServerPolicy.Conditions.Clients.Include))

	return nil
}

func resourceAuthServerPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	authServerPolicy := buildAuthServerPolicy(d)
	c := getSupplementFromMetadata(m)
	_, _, err := c.UpdateAuthorizationServerPolicy(d.Get("auth_server_id").(string), d.Id(), *authServerPolicy, nil)
	if err != nil {
		return err
	}

	return resourceAuthServerPolicyRead(d, m)
}

func resourceAuthServerPolicyDelete(d *schema.ResourceData, m interface{}) error {
	_, err := getSupplementFromMetadata(m).DeleteAuthorizationServerPolicy(d.Get("auth_server_id").(string), d.Id())

	return err
}

func fetchAuthServerPolicy(d *schema.ResourceData, m interface{}) (*sdk.AuthorizationServerPolicy, error) {
	c := getSupplementFromMetadata(m)
	auth, resp, err := c.GetAuthorizationServerPolicy(d.Get("auth_server_id").(string), d.Id(), sdk.AuthorizationServerPolicy{})

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	return auth, err
}

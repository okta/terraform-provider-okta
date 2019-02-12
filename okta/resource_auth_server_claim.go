package okta

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform/helper/validation"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAuthServerClaim() *schema.Resource {
	return &schema.Resource{
		Create: resourceAuthServerClaimCreate,
		Exists: resourceAuthServerClaimExists,
		Read:   resourceAuthServerClaimRead,
		Update: resourceAuthServerClaimUpdate,
		Delete: resourceAuthServerClaimDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				parts := strings.Split(d.Id(), "/")
				if len(parts) != 2 {
					return nil, fmt.Errorf("Invalid policy rule specifier. Expecting {auth_server_id}/{id}")
				}
				d.Set("auth_server_id", parts[0])
				d.SetId(parts[1])
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Auth server claim name",
			},
			"auth_server_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Auth server ID",
			},
			"scopes": &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Auth server claim list of scopes",
			},
			"status": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ACTIVE",
				ValidateFunc: validation.StringInSlice([]string{"ACTIVE", "INACTIVE"}, false),
			},
			"value": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"value_type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"EXPRESSION"}, false),
			},
			"claim_type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"RESOURCE"}, false),
			},
			"always_include_in_token": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func buildAuthServerClaim(d *schema.ResourceData) *AuthorizationServerClaim {
	return &AuthorizationServerClaim{
		Status:               d.Get("status").(string),
		ClaimType:            d.Get("claim_type").(string),
		ValueType:            d.Get("value_type").(string),
		Value:                d.Get("value").(string),
		AlwaysIncludeInToken: d.Get("always_include_in_token").(bool),
		Name:                 d.Get("name").(string),
		Conditions:           &Conditions{Scopes: convertInterfaceToStringArr(d.Get("scopes"))},
	}
}

func resourceAuthServerClaimCreate(d *schema.ResourceData, m interface{}) error {
	authServerClaim := buildAuthServerClaim(d)
	c := getSupplementFromMetadata(m)
	responseAuthServerClaim, _, err := c.CreateAuthorizationServerClaim(d.Get("auth_server_id").(string), *authServerClaim, nil)
	if err != nil {
		return err
	}

	d.SetId(responseAuthServerClaim.Id)

	return resourceAuthServerClaimRead(d, m)
}

func resourceAuthServerClaimExists(d *schema.ResourceData, m interface{}) (bool, error) {
	g, err := fetchAuthServerClaim(d, m)

	return err == nil && g != nil, err
}

func resourceAuthServerClaimRead(d *schema.ResourceData, m interface{}) error {
	authServerClaim, err := fetchAuthServerClaim(d, m)
	if err != nil {
		return err
	}

	if authServerClaim.Conditions != nil && len(authServerClaim.Conditions.Scopes) > 0 {
		d.Set("scopes", authServerClaim.Conditions.Scopes)
	}

	d.Set("name", authServerClaim.Name)
	d.Set("status", authServerClaim.Status)
	d.Set("value", authServerClaim.Value)
	d.Set("value_type", authServerClaim.ValueType)
	d.Set("claim_type", authServerClaim.ClaimType)
	d.Set("always_include_in_token", authServerClaim.AlwaysIncludeInToken)

	return nil
}

func resourceAuthServerClaimUpdate(d *schema.ResourceData, m interface{}) error {
	authServerClaim := buildAuthServerClaim(d)
	c := getSupplementFromMetadata(m)
	_, _, err := c.UpdateAuthorizationServerClaim(d.Get("auth_server_id").(string), d.Id(), *authServerClaim, nil)
	if err != nil {
		return err
	}

	return resourceAuthServerClaimRead(d, m)
}

func resourceAuthServerClaimDelete(d *schema.ResourceData, m interface{}) error {
	_, err := getSupplementFromMetadata(m).DeleteAuthorizationServerClaim(d.Get("auth_server_id").(string), d.Id())

	return err
}

func fetchAuthServerClaim(d *schema.ResourceData, m interface{}) (*AuthorizationServerClaim, error) {
	c := getSupplementFromMetadata(m)
	auth, resp, err := c.GetAuthorizationServerClaim(d.Get("auth_server_id").(string), d.Id(), AuthorizationServerClaim{})

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	return auth, err
}

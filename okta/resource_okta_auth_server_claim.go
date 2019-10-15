package okta

import (
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-okta/sdk"
)

func resourceAuthServerClaim() *schema.Resource {
	return &schema.Resource{
		Create:   resourceAuthServerClaimCreate,
		Exists:   resourceAuthServerClaimExists,
		Read:     resourceAuthServerClaimRead,
		Update:   resourceAuthServerClaimUpdate,
		Delete:   resourceAuthServerClaimDelete,
		Importer: createNestedResourceImporter([]string{"auth_server_id", "id"}),

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
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"status": statusSchema,
			"value": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"value_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"EXPRESSION", "GROUPS"}, false),
				Default:      "EXPRESSION",
			},
			"claim_type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"RESOURCE", "IDENTITY"}, false),
			},
			"always_include_in_token": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"group_filter_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"STARTS_WITH", "EQUALS", "CONTAINS", "REGEX"}, false),
				Description:  "Required when value_type is GROUPS",
			},
		},
	}
}

func buildAuthServerClaim(d *schema.ResourceData) *sdk.AuthorizationServerClaim {
	return &sdk.AuthorizationServerClaim{
		Status:               d.Get("status").(string),
		ClaimType:            d.Get("claim_type").(string),
		ValueType:            d.Get("value_type").(string),
		Value:                d.Get("value").(string),
		AlwaysIncludeInToken: d.Get("always_include_in_token").(bool),
		Name:                 d.Get("name").(string),
		Conditions:           &sdk.ClaimConditions{Scopes: convertInterfaceToStringSetNullable(d.Get("scopes"))},
		GroupFilterType:      d.Get("group_filter_type").(string),
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

	if authServerClaim == nil {
		d.SetId("")
		return nil
	}

	if err != nil {
		return err
	}

	if authServerClaim.Conditions != nil && len(authServerClaim.Conditions.Scopes) > 0 {
		d.Set("scopes", convertStringSetToInterface(authServerClaim.Conditions.Scopes))
	}

	d.Set("name", authServerClaim.Name)
	d.Set("status", authServerClaim.Status)
	d.Set("value", authServerClaim.Value)
	d.Set("value_type", authServerClaim.ValueType)
	d.Set("claim_type", authServerClaim.ClaimType)
	d.Set("always_include_in_token", authServerClaim.AlwaysIncludeInToken)
	d.Set("group_filter_type", authServerClaim.GroupFilterType)

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

func fetchAuthServerClaim(d *schema.ResourceData, m interface{}) (*sdk.AuthorizationServerClaim, error) {
	c := getSupplementFromMetadata(m)
	auth, resp, err := c.GetAuthorizationServerClaim(d.Get("auth_server_id").(string), d.Id(), sdk.AuthorizationServerClaim{})

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	return auth, err
}

package okta

import (
	"fmt"
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
		CustomizeDiff: func(d *schema.ResourceDiff, v interface{}) error {},
		Schema: map[string]*schema.Schema{}
	}
}

func resourceIdentityProviderCreate(d *schema.ResourceData, m interface{}) error {}
func resourceIdentityProviderRead(d *schema.ResourceData, m interface{}) error {}
func resourceIdentityProviderUpdate(d *schema.ResourceData, m interface{}) error {}
func resourceIdentityProviderDelete(d *schema.ResourceData, m interface{}) error {}

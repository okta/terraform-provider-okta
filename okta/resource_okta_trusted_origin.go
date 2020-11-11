package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func resourceTrustedOrigin() *schema.Resource {
	return &schema.Resource{
		Create: resourceTrustedOriginCreate,
		Read:   resourceTrustedOriginRead,
		Update: resourceTrustedOriginUpdate,
		Delete: resourceTrustedOriginDelete,
		Exists: resourceTrustedOriginExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the Trusted Origin is active or not - can only be issued post-creation",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique name for this trusted origin",
			},
			"origin": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique origin URL for this trusted origin",
			},
			"scopes": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Scopes of the Trusted Origin - can either be CORS or REDIRECT only",
			},
		},
	}
}

func resourceTrustedOriginCreate(d *schema.ResourceData, m interface{}) error {
	if !d.Get("active").(bool) {
		return fmt.Errorf("can not create inactive trusted origin, only existing trusted origins can be deactivated")
	}

	client := getOktaClientFromMetadata(m)

	returnedTrustedOrigin, _, err := client.TrustedOrigin.CreateOrigin(context.Background(), buildTrustedOrigin(d))
	if err != nil {
		return fmt.Errorf("failed to create trusted origin: %v", err)
	}
	return setTrustedOrigin(d, returnedTrustedOrigin)
}

func resourceTrustedOriginRead(d *schema.ResourceData, m interface{}) error {
	to, resp, err := getOktaClientFromMetadata(m).TrustedOrigin.GetOrigin(context.Background(), d.Id())
	if resp != nil && is404(resp.StatusCode) {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to get trusted origin: %v", err)
	}
	return setTrustedOrigin(d, to)
}

func resourceTrustedOriginUpdate(d *schema.ResourceData, m interface{}) error {
	client1 := getOktaClientFromMetadata(m)

	if d.HasChange("active") {
		var err error
		if d.Get("active").(bool) {
			_, _, err = client1.TrustedOrigin.ActivateOrigin(context.Background(), d.Id())
		} else {
			_, _, err = client1.TrustedOrigin.DeactivateOrigin(context.Background(), d.Id())
		}
		if err != nil {
			return fmt.Errorf("failed to change trusted origin status: %v", err)
		}
	}
	to, _, err := client1.TrustedOrigin.UpdateOrigin(context.Background(), d.Id(), buildTrustedOrigin(d))
	if err != nil {
		return fmt.Errorf("failed to update trusted origin: %v", err)
	}

	return setTrustedOrigin(d, to)
}

func resourceTrustedOriginDelete(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	_, err := client.TrustedOrigin.DeleteOrigin(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("failed to delete trusted origin: %v", err)
	}
	return nil
}

// check if Trusted Origin exists in Okta
func resourceTrustedOriginExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := getOktaClientFromMetadata(m)
	_, resp, err := client.TrustedOrigin.GetOrigin(context.Background(), d.Id())
	if resp != nil && is404(resp.StatusCode) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check trusted origin existence: %v", err)
	}
	return true, nil
}

// Creates Trusted Origin struct with the data resource provided by TF
func buildTrustedOrigin(d *schema.ResourceData) okta.TrustedOrigin {
	trustedOrigin := okta.TrustedOrigin{
		Name:   d.Get("name").(string),
		Origin: d.Get("origin").(string),
	}
	if d.Get("active").(bool) {
		trustedOrigin.Status = statusActive
	} else {
		trustedOrigin.Status = statusInactive
	}

	resScopes := d.Get("scopes").([]interface{})

	trustedOrigin.Scopes = make([]*okta.Scope, len(resScopes))

	for i, vals := range resScopes {
		trustedOrigin.Scopes[i] = &okta.Scope{
			Type: vals.(string),
		}
	}

	return trustedOrigin
}

func setTrustedOrigin(d *schema.ResourceData, to *okta.TrustedOrigin) error {
	d.SetId(to.Id)
	scopes := make([]string, len(to.Scopes))
	for i, scope := range to.Scopes {
		scopes[i] = scope.Type
	}

	_ = d.Set("active", to.Status == statusActive)
	_ = d.Set("origin", to.Origin)
	_ = d.Set("name", to.Name)

	return setNonPrimitives(d, map[string]interface{}{
		"scopes": scopes,
	})
}

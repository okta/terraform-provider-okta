package okta

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func resourceTrustedOrigin() *schema.Resource {
	return &schema.Resource{
		Create: resourceTrustedOriginCreate,
		Read:   resourceTrustedOriginRead,
		Update: resourceTrustedOriginUpdate,
		Delete: resourceTrustedOriginDelete,
		Exists: trustedOriginExists,
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
				Description: "Name of the Trusted Origin Resource",
			},
			"origin": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The origin to trust",
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

// Populates the Trusted Origin struct (used by the Okta SDK for API operaations) with the data resource provided by TF
func populateTrustedOrigin(trustedOrigin *okta.TrustedOrigin, d *schema.ResourceData) *okta.TrustedOrigin {
	trustedOrigin.Name = d.Get("name").(string)
	trustedOrigin.Origin = d.Get("origin").(string)

	var scopes []*okta.Scope

	for _, vals := range d.Get("scopes").([]interface{}) {
		scopes = append(scopes, &okta.Scope{Type: vals.(string)})
	}

	trustedOrigin.Scopes = scopes

	return trustedOrigin
}

func resourceTrustedOriginCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Create Trusted Origin %v", d.Get("name").(string))

	if !d.Get("active").(bool) {
		return fmt.Errorf("[ERROR] Okta will not allow a Trusted Origin to be created as INACTIVE. Can set to false for existing Trusted Origins only.")
	}

	ctx, client := getOktaClientFromMetadata(m)
	trustedOrigin := &okta.TrustedOrigin{}
	trustedOrigin = populateTrustedOrigin(trustedOrigin, d)

	returnedTrustedOrigin, _, err := client.TrustedOrigin.CreateOrigin(ctx, *trustedOrigin)

	if err != nil {
		return fmt.Errorf("[ERROR] %v.", err)
	}

	d.SetId(returnedTrustedOrigin.Id)

	return resourceTrustedOriginRead(d, m)
}

func resourceTrustedOriginRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Read Trusted Origin %v", d.Get("name").(string))

	var trustedOrigin *okta.TrustedOrigin

	ctx, client := getOktaClientFromMetadata(m)

	exists, err := trustedOriginExists(d, m)
	if err != nil {
		return err
	}

	if exists == true {
		trustedOrigin, _, err = client.TrustedOrigin.GetOrigin(ctx, d.Id())
	} else {
		d.SetId("")
		return nil
	}

	scopes := make([]string, 0)
	for _, scope := range trustedOrigin.Scopes {
		scopes = append(scopes, scope.Type)
	}

	d.Set("active", trustedOrigin.Status == "ACTIVE")
	d.Set("origin", trustedOrigin.Origin)
	d.Set("name", trustedOrigin.Name)

	return setNonPrimitives(d, map[string]interface{}{
		"scopes": scopes,
	})
}

func resourceTrustedOriginUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Update Trusted Origin %v", d.Get("name").(string))

	var trustedOrigin = &okta.TrustedOrigin{}
	trustedOrigin = populateTrustedOrigin(trustedOrigin, d)

	ctx, client := getOktaClientFromMetadata(m)

	exists, err := trustedOriginExists(d, m)
	if err != nil {
		return err
	}

	if exists == true {
		if d.HasChange("active") {
			_, _, err = client.TrustedOrigin.ActivateOrigin(ctx, d.Id())
			if err != nil {
				return fmt.Errorf("[ERROR] Error Updating Trusted Origin with Okta: %v", err)
			}
		}

		_, _, err = client.TrustedOrigin.UpdateOrigin(ctx, d.Id(), *trustedOrigin)

		if err != nil {
			return fmt.Errorf("[ERROR] Error Updating Trusted Origin with Okta: %v", err)
		}
	} else {
		d.SetId("")
		return nil
	}

	return resourceTrustedOriginRead(d, m)
}

func resourceTrustedOriginDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Delete Trusted Origin %v", d.Get("name").(string))

	ctx, client := getOktaClientFromMetadata(m)

	exists, err := trustedOriginExists(d, m)
	if err != nil {
		return err
	}
	if exists == true {
		_, err = client.TrustedOrigin.DeleteOrigin(ctx, d.Id())
		if err != nil {
			return fmt.Errorf("[ERROR] Error Deleting Trusted Origin from Okta: %v", err)
		}
	}

	return nil
}

// check if Trusted Origin exists in Okta
func trustedOriginExists(d *schema.ResourceData, m interface{}) (bool, error) {
	ctx, client := getOktaClientFromMetadata(m)
	_, resp, err := client.TrustedOrigin.GetOrigin(ctx, d.Id())

	if resp.Status == "404 Not Found" {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("[ERROR] Error Getting Trusted Origin in Okta: %v", err)
	}
	return true, nil
}

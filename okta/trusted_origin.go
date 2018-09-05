package okta

import (
  "github.com/articulate/oktasdk-go/okta"
  "github.com/hashicorp/terraform/helper/schema"
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
      "active": &schema.Schema{
        Type:        schema.TypeBool,
        Optional:    true,
        Default:     true,
        Description: "Whether the Trusted Origin is active or not - can only be issued post-creation",
      },
      "name": &schema.Schema{
        Type:        schema.TypeString,
        Required:    true,
        Description: "Name of the Trusted Origin Resource",
      },
      "origin": &schema.Schema{
        Type:        schema.TypeString,
        Required:    true,
        Description: "The origin to trust",
      },
      "scopes": &schema.Schema{
        Type:        schema.TypeList,
        Optional:    true,
        Elem:        &schema.Schema{Type: schema.TypeString},
        Description: "Scopes of the Trusted Origin - can either be CORS or Redirect only",
      },
    },
  }
}

func assembleTrustedOrigin() *okta.TrustedOrigin {
  hints := &okta.TrustedOriginHints{}

  deactivate := &okta.TrustedOriginDeactive{
    Hints: hints,
  }

  self := &okta.TrustedOriginSelf{
    Hints: hints,
  }

  links := &okta.TrustedOriginLinks{
    Self: self,
    Deactivate: deactivate,
  }

  trustedOrigin := &okta.TrustedOrigin{
    Links: links,
  }

  return trustedOrigin
}

// Populates the Trusted Origin struct (used by the Okta SDK for API operaations) with the data resource provided by TF
func populateTrustedOrigin(trustedOrigin *okta.TrustedOrigin, d *schema.ResourceData) *okta.TrustedOrigin {
  trustedOrigin.ID = d.Get("id").(string)
  trustedOrigin.Name = d.Get("name").(string)
  trustedOrigin.Origin = d.Get("origin").(string)

  var scopes []map[string]string

  for _, vals := range d.Get("scopes").([]interface{}) {
    scopes = append(scopes, map[string]string{"Type": vals.(string)})
  }

  trustedOrigin.Scopes = scopes

  return trustedOrigin
}

func resourceTrustedOriginCreate(d *schema.ResourceData, m interface{}) error {
  return nil
}

func resourceTrustedOriginRead(d *schema.ResourceData, m interface{}) error {
  return nil
}

func resourceTrustedOriginUpdate(d *schema.ResourceData, m interface{}) error {
  return nil
}

func resourceTrustedOriginDelete(d *schema.ResourceData, m interface{}) error {
  return nil
}

// check if Trusted Origin exists in Okta
func trustedOriginExists(d *schema.ResourceData, m interface{}) (bool, error) {
  return true, nil
}

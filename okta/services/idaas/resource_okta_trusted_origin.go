package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
)

func resourceTrustedOrigin() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTrustedOriginCreate,
		ReadContext:   resourceTrustedOriginRead,
		UpdateContext: resourceTrustedOriginUpdate,
		DeleteContext: resourceTrustedOriginDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "Creates a Trusted Origin. This resource allows you to create and configure a Trusted Origin.",
		Schema: map[string]*schema.Schema{
			"active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the Trusted Origin is active or not - can only be issued post-creation. By default, it is `true`.",
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
				Description: "Scopes of the Trusted Origin - can either be `CORS` and/or `REDIRECT`",
			},
		},
	}
}

func resourceTrustedOriginCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if !d.Get("active").(bool) {
		return diag.Errorf("can not create inactive trusted origin, only existing trusted origins can be deactivated")
	}
	trustedOrigin, _, err := getOktaClientFromMetadata(meta).TrustedOrigin.CreateOrigin(ctx, buildTrustedOrigin(d))
	if err != nil {
		return diag.Errorf("failed to create trusted origin: %v", err)
	}
	d.SetId(trustedOrigin.Id)
	err = setTrustedOrigin(d, trustedOrigin)
	if err != nil {
		return diag.Errorf("failed to set trusted origin's properties: %v", err)
	}
	return nil
}

func resourceTrustedOriginRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	trustedOrigin, resp, err := getOktaClientFromMetadata(meta).TrustedOrigin.GetOrigin(ctx, d.Id())
	if err := utils.SuppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get trusted origin: %v", err)
	}
	if trustedOrigin == nil {
		d.SetId("")
		return nil
	}
	err = setTrustedOrigin(d, trustedOrigin)
	if err != nil {
		return diag.Errorf("failed to set trusted origin's properties: %v", err)
	}
	return nil
}

func resourceTrustedOriginUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(meta)
	if d.HasChange("active") {
		var err error
		if d.Get("active").(bool) {
			_, _, err = client.TrustedOrigin.ActivateOrigin(ctx, d.Id())
		} else {
			_, _, err = client.TrustedOrigin.DeactivateOrigin(ctx, d.Id())
		}
		if err != nil {
			return diag.Errorf("failed to change trusted origin's status: %v", err)
		}
	}
	trustedOrigin, _, err := client.TrustedOrigin.UpdateOrigin(ctx, d.Id(), buildTrustedOrigin(d))
	if err != nil {
		return diag.Errorf("failed to update trusted origin: %v", err)
	}
	err = setTrustedOrigin(d, trustedOrigin)
	if err != nil {
		return diag.Errorf("failed to set trusted origin's properties: %v", err)
	}
	return nil
}

func resourceTrustedOriginDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_, err := getOktaClientFromMetadata(meta).TrustedOrigin.DeleteOrigin(ctx, d.Id())
	if err != nil {
		return diag.Errorf("failed to delete trusted origin: %v", err)
	}
	return nil
}

// Creates Trusted Origin struct with the data resource provided by TF
func buildTrustedOrigin(d *schema.ResourceData) sdk.TrustedOrigin {
	trustedOrigin := sdk.TrustedOrigin{
		Name:   d.Get("name").(string),
		Origin: d.Get("origin").(string),
	}
	if d.Get("active").(bool) {
		trustedOrigin.Status = StatusActive
	} else {
		trustedOrigin.Status = StatusInactive
	}
	resScopes := d.Get("scopes").([]interface{})
	trustedOrigin.Scopes = make([]*sdk.Scope, len(resScopes))
	for i := range resScopes {
		trustedOrigin.Scopes[i] = &sdk.Scope{
			Type: resScopes[i].(string),
		}
	}
	return trustedOrigin
}

func setTrustedOrigin(d *schema.ResourceData, to *sdk.TrustedOrigin) error {
	scopes := make([]string, len(to.Scopes))
	for i, scope := range to.Scopes {
		scopes[i] = scope.Type
	}
	_ = d.Set("active", to.Status == StatusActive)
	_ = d.Set("origin", to.Origin)
	_ = d.Set("name", to.Name)
	return utils.SetNonPrimitives(d, map[string]interface{}{
		"scopes": scopes,
	})
}

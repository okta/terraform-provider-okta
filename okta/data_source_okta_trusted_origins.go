package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

func dataSourceTrustedOrigin() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTrustedOriginsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "Trusted origin Id",
				Required:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Unique name for this trusted origin",
				Required:    true,
			},
			"origin": {
				Type:        schema.TypeString,
				Description: "Unique origin URL for this trusted origin",
				Required:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Trusted origin's status whether it is active or not",
				Required:    true,
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

func dataSourceTrustedOriginsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	query := d.Get("query").(*query.Params)
	trustedOrigins, resp, err := getOktaClientFromMetadata(m).TrustedOrigin.ListOrigins(ctx, query)

	if err != nil || resp.StatusCode != 200 {
		return diag.Errorf("failed to list trusted origins: %v", err)
	}

	if trustedOrigins == nil {
		d.SetId("")
		return nil
	}

	arr := make([]map[string]interface{}, len(trustedOrigins))
	for i := range trustedOrigins {
		arr[i] = flattenTrustedOrigins(trustedOrigins[i])
	}
	_ = d.Set("trustedOrigins", arr)
	return nil
}

func flattenTrustedOrigins(c *okta.TrustedOrigin) map[string]interface{} {
	m := map[string]interface{}{
		"id":     c.Id,
		"name":   c.Name,
		"origin": c.Origin,
		"status": c.Status,
		"scopes": c.Scopes,
	}
	return m
}

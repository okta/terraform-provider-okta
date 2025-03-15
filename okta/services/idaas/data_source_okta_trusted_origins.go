package idaas

import (
	"context"
	"fmt"
	"hash/crc32"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
	"github.com/okta/terraform-provider-okta/sdk/query"
)

func dataSourceTrustedOrigins() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTrustedOriginsRead,
		Schema: map[string]*schema.Schema{
			"filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter criteria. Filter value will be URL-encoded by the provider",
			},
			"trusted_origins": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique identifier",
						},
						"active": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the Trusted Origin is active or not - can only be issued post-creation",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique name for this trusted origin",
						},
						"origin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique origin URL for this trusted origin",
						},
						"scopes": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Scopes of the Trusted Origin - can either be CORS or REDIRECT only",
						},
					},
				},
			},
		},
		Description: "Get List of Trusted Origins using filters.",
	}
}

func dataSourceTrustedOriginsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	qp := &query.Params{Limit: utils.DefaultPaginationLimit}
	filter, ok := d.GetOk("filter")
	if ok {
		qp.Filter = filter.(string)
	}
	trustedOrigins, err := collectTrustedOrigins(ctx, getOktaClientFromMetadata(meta), qp)
	if err != nil {
		return diag.Errorf("failed to trusted origins: %v", err)
	}
	d.SetId(fmt.Sprintf("%d", crc32.ChecksumIEEE([]byte(qp.String()))))
	arr := make([]map[string]interface{}, len(trustedOrigins))
	for i := range trustedOrigins {
		scopes := make([]string, len(trustedOrigins[i].Scopes))
		for j := range trustedOrigins[i].Scopes {
			scopes[j] = trustedOrigins[i].Scopes[j].Type
		}
		arr[i] = map[string]interface{}{
			"id":     trustedOrigins[i].Id,
			"active": trustedOrigins[i].Status == StatusActive,
			"name":   trustedOrigins[i].Name,
			"origin": trustedOrigins[i].Origin,
			"scopes": scopes,
		}
	}
	_ = d.Set("trusted_origins", arr)
	return nil
}

func collectTrustedOrigins(ctx context.Context, client *sdk.Client, qp *query.Params) ([]*sdk.TrustedOrigin, error) {
	trustedOrigins, resp, err := client.TrustedOrigin.ListOrigins(ctx, qp)
	if err != nil {
		return nil, err
	}
	for resp.HasNextPage() {
		var nextTrustedOrigins []*sdk.TrustedOrigin
		resp, err = resp.Next(ctx, &nextTrustedOrigins)
		if err != nil {
			return nil, err
		}
		trustedOrigins = append(trustedOrigins, nextTrustedOrigins...)
	}
	return trustedOrigins, nil
}

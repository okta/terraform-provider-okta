package okta

import (
	"context"
	"fmt"
	"hash/crc32"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOrgSupport() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOrgSupportCreate,
		ReadContext:   resourceOrgSupportRead,
		DeleteContext: resourceOrgSupportDelete,
		Importer:      nil,
		Schema: map[string]*schema.Schema{
			"extend_by": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Number of days the support should be extended by",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of Okta Support",
			},
			"expiration": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Expiration of Okta Support",
			},
		},
	}
}

func resourceOrgSupportCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	support, _, err := getOktaClientFromMetadata(m).OrgSetting.GrantOktaSupport(ctx)
	if err != nil {
		return diag.Errorf("failed to grant org support: %v", err)
	}
	eb, ok := d.GetOk("extend_by")
	if ok && eb.(int) > 0 {
		for i := 0; i < eb.(int); i++ {
			support, _, err = getOktaClientFromMetadata(m).OrgSetting.ExtendOktaSupport(ctx)
			if err != nil {
				return diag.Errorf("failed to extend org support: %v", err)
			}
		}
	}
	d.SetId(fmt.Sprintf("%d", crc32.ChecksumIEEE([]byte(support.Expiration.String()))))
	_ = d.Set("expiration", support.Expiration.String())
	_ = d.Set("status", support.Support)
	return nil
}

func resourceOrgSupportRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	support, _, err := getOktaClientFromMetadata(m).OrgSetting.GetOrgOktaSupportSettings(ctx)
	if err != nil {
		return diag.Errorf("failed to get org support settings: %v", err)
	}
	if support.Expiration != nil {
		_ = d.Set("expiration", support.Expiration.String())
	}
	_ = d.Set("status", support.Support)
	return nil
}

func resourceOrgSupportDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	_, _, err := getOktaClientFromMetadata(m).OrgSetting.RevokeOktaSupport(ctx)
	if err != nil {
		return diag.Errorf("failed to revoke okta support: %v", err)
	}
	return nil
}

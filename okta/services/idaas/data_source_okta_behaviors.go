package idaas

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"io"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk/query"
)

func dataSourceBehaviors() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBehaviorsRead,
		Schema: map[string]*schema.Schema{
			"q": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Searches the name property of behaviors for matching value",
			},
			"behaviors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Behavior ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Behavior name.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Behavior status.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Behavior type.",
						},
						"settings": {
							Type:        schema.TypeMap,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "Map of behavior settings.",
						},
					},
				},
			},
		},
		Description: "Get a behaviors by search criteria.",
	}
}

func dataSourceBehaviorsRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	behaviorRulesFromRawResp, err := getBehaviorRules(ctx, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	qp := &query.Params{Limit: utils.DefaultPaginationLimit}
	q, ok := d.GetOk("q")
	if ok {
		qp.Q = q.(string)
	}
	behaviorRules := []map[string]any{}
	for _, behaviorRule := range behaviorRulesFromRawResp {
		if !strings.Contains(behaviorRule["name"].(string), q.(string)) {
			continue // since name does not contain our "q" substring partial match, continue
		}
		settingsMap := make(map[string]string)
		for k, v := range behaviorRule["settings"].(map[string]any) {
			settingsMap[k] = fmt.Sprint(v)
		}
		behaviorRules = append(behaviorRules, map[string]any{
			"id":       behaviorRule["id"],
			"name":     behaviorRule["name"],
			"type":     behaviorRule["type"],
			"status":   behaviorRule["status"],
			"settings": settingsMap,
		})
	}
	d.SetId(fmt.Sprintf("%d", crc32.ChecksumIEEE([]byte(qp.String()))))
	err = d.Set("behaviors", behaviorRules)
	return diag.FromErr(err)
}

func getBehaviorRules(ctx context.Context, meta any) ([]map[string]any, error) {
	var behaviorRulesFromRawResp []map[string]any
	_, rawResp, err := getOktaV5ClientFromMetadata(meta).BehaviorAPI.ListBehaviorDetectionRules(ctx).Execute()
	if err != nil {
		if 200 <= rawResp.StatusCode && rawResp.StatusCode <= 299 {
			if strings.HasPrefix(err.Error(), "parsing time") {
				logger(meta).Info("error when parsing time, will process raw HTTP response")
			}
			if strings.Contains(err.Error(), "cannot unmarshal number") {
				logger(meta).Info("error when parsing number, will process raw HTTP response")
			}
		} else {
			return behaviorRulesFromRawResp, fmt.Errorf("failed to query for behaviors: %v", err)
		}
	}
	rawRespBody, err := io.ReadAll(rawResp.Body)
	if err != nil {
		return behaviorRulesFromRawResp, err
	}
	err = json.Unmarshal(rawRespBody, &behaviorRulesFromRawResp)
	return behaviorRulesFromRawResp, err
}

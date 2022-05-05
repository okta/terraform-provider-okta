package okta

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

const (
	push     = "PUSH"
	dontPush = "DONT_PUSH"
)

func resourceProfileMapping() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProfileMappingCreate,
		ReadContext:   resourceProfileMappingRead,
		UpdateContext: resourceProfileMappingUpdate,
		DeleteContext: resourceProfileMappingDelete,
		Schema: map[string]*schema.Schema{
			"source_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The source id of the mapping to manage.",
			},
			"delete_when_absent": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "When turned on this flag will trigger the provider to delete mapping properties that are not defined in config. By default, we do not delete missing properties.",
			},
			"source_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The target id of the mapping to manage.",
			},
			"target_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mappings": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     mappingResource,
			},
			"always_apply": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether apply the changes to all users with this profile after updating or creating the these mappings.",
				Default:     false,
			},
		},
	}
}

var mappingResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The mapping property key.",
		},
		"expression": {
			Type:     schema.TypeString,
			Required: true,
		},
		"push_status": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          dontPush,
			ValidateDiagFunc: elemInSlice([]string{push, dontPush}),
		},
	},
}

func resourceProfileMappingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sourceID := d.Get("source_id").(string)
	targetID := d.Get("target_id").(string)
	mapping, resp, err := getProfileMappingBySourceID(ctx, sourceID, targetID, m)
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get profile mapping: %v", err)
	}
	if mapping == nil {
		return diag.Errorf("no profile mappings found for source ID '%s' and target ID '%s'", sourceID, targetID)
	}
	d.SetId(mapping.Id)
	newMapping := buildMapping(d)
	if d.Get("delete_when_absent").(bool) {
		newMapping.Properties = mergeProperties(newMapping.Properties, getDeleteProperties(d, mapping.Properties))
	}
	_, _, err = getOktaClientFromMetadata(m).ProfileMapping.UpdateProfileMapping(ctx, mapping.Id, newMapping)
	if err != nil {
		return diag.Errorf("failed to create profile mapping: %v", err)
	}
	err = applyMapping(ctx, d, m, mapping)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceProfileMappingRead(ctx, d, m)
}

func resourceProfileMappingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	mapping, resp, err := getOktaClientFromMetadata(m).ProfileMapping.GetProfileMapping(ctx, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get profile mapping: %v", err)
	}
	if mapping == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("source_type", mapping.Source.Type)
	_ = d.Set("source_name", mapping.Source.Name)
	_ = d.Set("target_type", mapping.Target.Type)
	_ = d.Set("target_id", mapping.Target.Id)
	_ = d.Set("target_name", mapping.Target.Name)
	if !d.Get("delete_when_absent").(bool) {
		current := buildMappingProperties(d.Get("mappings").(*schema.Set))
		for k := range mapping.Properties {
			if _, ok := current[k]; !ok {
				delete(mapping.Properties, k)
			}
		}
	}
	_ = d.Set("mappings", flattenMappingProperties(mapping.Properties))
	return nil
}

func resourceProfileMappingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sourceID := d.Get("source_id").(string)
	targetID := d.Get("target_id").(string)
	mapping, resp, err := getOktaClientFromMetadata(m).ProfileMapping.GetProfileMapping(ctx, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get profile mapping: %v", err)
	}
	if mapping == nil {
		return diag.Errorf("no profile mappings found for source ID '%s' and target ID '%s'", sourceID, targetID)
	}
	newMapping := buildMapping(d)
	if d.Get("delete_when_absent").(bool) {
		newMapping.Properties = mergeProperties(newMapping.Properties, getDeleteProperties(d, mapping.Properties))
	}
	_, _, err = getOktaClientFromMetadata(m).ProfileMapping.UpdateProfileMapping(ctx, mapping.Id, newMapping)
	if err != nil {
		return diag.Errorf("failed to update profile mapping: %v", err)
	}
	err = applyMapping(ctx, d, m, mapping)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceProfileMappingRead(ctx, d, m)
}

func resourceProfileMappingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sourceID := d.Get("source_id").(string)
	targetID := d.Get("target_id").(string)
	mapping, resp, err := getOktaClientFromMetadata(m).ProfileMapping.GetProfileMapping(ctx, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get profile mapping: %v", err)
	}
	if mapping == nil {
		return diag.Errorf("no profile mappings found for source ID '%s' and target ID '%s'", sourceID, targetID)
	}
	for k := range mapping.Properties {
		if k == "login" {
			continue
		}
		mapping.Properties[k] = nil
	}
	_, _, err = getOktaClientFromMetadata(m).ProfileMapping.UpdateProfileMapping(ctx, mapping.Id, *mapping)
	if err != nil {
		return diag.Errorf("failed to delete profile mapping: %v", err)
	}
	return nil
}

func getDeleteProperties(d *schema.ResourceData, actual map[string]*okta.ProfileMappingProperty) map[string]*okta.ProfileMappingProperty {
	toDelete := map[string]*okta.ProfileMappingProperty{}
	config := buildMappingProperties(d.Get("mappings").(*schema.Set))
	for key := range actual {
		if _, ok := config[key]; !ok {
			toDelete[key] = nil
		}
	}
	return toDelete
}

func mergeProperties(target, b map[string]*okta.ProfileMappingProperty) map[string]*okta.ProfileMappingProperty {
	for k, v := range b {
		target[k] = v
	}
	return target
}

func flattenMappingProperties(src map[string]*okta.ProfileMappingProperty) *schema.Set {
	var arr []interface{}
	for k, v := range src {
		arr = append(arr, map[string]interface{}{
			"id":          k,
			"push_status": v.PushStatus,
			"expression":  v.Expression,
		})
	}
	return schema.NewSet(schema.HashResource(mappingResource), arr)
}

func buildMappingProperties(set *schema.Set) map[string]*okta.ProfileMappingProperty {
	res := map[string]*okta.ProfileMappingProperty{}
	for _, rawMap := range set.List() {
		if m, ok := rawMap.(map[string]interface{}); ok {
			k := m["id"].(string)

			res[k] = &okta.ProfileMappingProperty{
				Expression: m["expression"].(string),
				PushStatus: m["push_status"].(string),
			}
		}
	}
	return res
}

func buildMapping(d *schema.ResourceData) okta.ProfileMapping {
	return okta.ProfileMapping{
		Id:         d.Id(),
		Properties: buildMappingProperties(d.Get("mappings").(*schema.Set)),
	}
}

func applyMapping(ctx context.Context, d *schema.ResourceData, m interface{}, mapping *okta.ProfileMapping) error {
	if !d.Get("always_apply").(bool) {
		return nil
	}
	source := d.Get("source_id").(string)
	target := d.Get("target_id").(string)
	var appID string
	if mapping.Source.Type == "appuser" {
		appID = mapping.Source.Id
	}
	if mapping.Target.Type == "appuser" {
		appID = mapping.Target.Id
	}
	appUserTypes, _, err := getSupplementFromMetadata(m).GetAppUserTypes(ctx, appID)
	if err != nil {
		return fmt.Errorf("failed to list app user types: %v", err)
	}
	if len(appUserTypes) == 0 || len(appUserTypes) > 2 {
		log.Println("[WARN] mappings were not applied")
		return nil
	}
	if mapping.Source.Type == "appuser" {
		source = appUserTypes[0].Id
	} else {
		target = appUserTypes[0].Id
	}
	// FIXME uses internal api
	_, err = getSupplementFromMetadata(m).ApplyMappings(ctx, source, target)
	if err != nil {
		return fmt.Errorf("failed to apply mappings for source '%s' and target '%s': %v", source, target, err)
	}
	return nil
}

func getProfileMappingBySourceID(ctx context.Context, sourceId, targetId string, m interface{}) (*okta.ProfileMapping, *okta.Response, error) {
	qp := query.Params{
		Limit: defaultPaginationLimit,
	}
	if sourceId != "" {
		qp.SourceId = sourceId
	}
	if targetId != "" {
		qp.TargetId = targetId
	}

	client := getOktaClientFromMetadata(m)
	mappings, resp, err := client.ProfileMapping.ListProfileMappings(ctx, &qp)
	if err != nil {
		return nil, resp, err
	}

	for _, mapping := range mappings {
		if mapping.Source.Id == sourceId {
			return client.ProfileMapping.GetProfileMapping(ctx, mapping.Id)
		}
	}

	return nil, resp, err
}

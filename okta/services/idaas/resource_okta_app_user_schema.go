package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
)

func resourceAppUserSchema() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppUserSchemaCreate,
		ReadContext:   resourceAppUserSchemaRead,
		UpdateContext: resourceAppUserSchemaUpdate,
		DeleteContext: resourceAppUserSchemaDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: `Manages the entire app user schema for an application.

This resource manages all custom properties in an application's user schema as a single object. This approach aligns with how the Okta API actually works (single mutable schema object) and provides better visibility into auto-created properties when provisioning is enabled.

**Advantages over okta_app_user_schema_property:**
- Manages all properties in one place
- Auto-created properties (from provisioning) are visible in plan/state
- Single import operation captures entire schema
- More efficient (fewer API calls)

**IMPORTANT:** With 'enum', list its values as strings even though the 'type' may be something other than string. This is a limitation of the schema definition in the Terraform Plugin SDK runtime and we juggle the type correctly when making Okta API calls. Same holds for the 'const' value of 'one_of' as well as the 'array_*' variation of 'enum' and 'one_of'`,
		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The Application's ID",
			},
			"custom_property": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Custom properties in the schema",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"index": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The property name/index",
						},
						"title": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Display name for the property",
						},
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The type of the schema property. It can be `string`, `boolean`, `number`, `integer`, `array`, or `object`",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The description of the property",
						},
						"required": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether the property is required",
						},
						"scope": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "NONE",
							Description: "Determines whether an app user attribute can be set at the Personal `SELF` or Group `NONE` level. Default value is `NONE`.",
						},
						"min_length": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The minimum length of the property value. Only applies to type `string`",
						},
						"max_length": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The maximum length of the property value. Only applies to type `string`",
						},
						"enum": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Array of values a primitive property can be set to",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"external_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "External name of the property",
						},
						"external_namespace": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "External namespace of the property",
						},
						"master": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Master priority for the property. It can be set to `PROFILE_MASTER` or `OKTA`",
						},
						"permissions": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "READ_ONLY",
							Description: "Access control permissions for the property. It can be set to `READ_WRITE`, `READ_ONLY`, `HIDE`. Default: `READ_ONLY`",
						},
						"union": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "If `type` is set to `array`, used to set whether attribute value is determined by group priority `false`, or combine values across groups `true`. Can not be set to `true` if `scope` is set to `SELF`.",
						},
					},
				},
			},
		},
	}
}

func resourceAppUserSchemaCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appId := d.Get("app_id").(string)
	d.SetId(appId)
	return resourceAppUserSchemaUpdate(ctx, d, meta)
}

func resourceAppUserSchemaRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appId := d.Id()
	us, resp, err := getOktaClientFromMetadata(meta).UserSchema.GetApplicationUserSchema(ctx, appId)
	if err != nil {
		if utils.SuppressErrorOn404(resp, err) == nil {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to get application user schema: %v", err)
	}

	_ = d.Set("app_id", appId)

	// Read custom properties
	if us.Definitions != nil && us.Definitions.Custom != nil && us.Definitions.Custom.Properties != nil {
		customProps := make([]interface{}, 0)
		for index, attr := range us.Definitions.Custom.Properties {
			propMap := flattenSchemaProperty(index, attr)
			customProps = append(customProps, propMap)
		}
		_ = d.Set("custom_property", customProps)
	}

	return nil
}

func resourceAppUserSchemaUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appId := d.Get("app_id").(string)

	// Build the full schema update
	schema := &sdk.UserSchema{
		Definitions: &sdk.UserSchemaDefinitions{
			Custom: &sdk.UserSchemaPublic{
				Id:         "#custom",
				Type:       "object",
				Properties: make(map[string]*sdk.UserSchemaAttribute),
			},
		},
	}

	// Add custom properties from config
	if v, ok := d.GetOk("custom_property"); ok {
		// TypeSet is stored as a Set interface with a List() method
		if customPropsSet, ok := v.(interface{ List() []interface{} }); ok {
			customProps := customPropsSet.List()
			for _, prop := range customProps {
				propMap := prop.(map[string]interface{})
				index := propMap["index"].(string)
				attr := expandSchemaProperty(propMap)
				schema.Definitions.Custom.Properties[index] = attr
			}
		}
	}

	// Call API
	_, resp, err := getOktaClientFromMetadata(meta).UserSchema.UpdateApplicationUserProfile(ctx, appId, *schema)
	if err != nil {
		if utils.SuppressErrorOn404(resp, err) == nil {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to update application user schema: %v", err)
	}

	return resourceAppUserSchemaRead(ctx, d, meta)
}

func resourceAppUserSchemaDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appId := d.Get("app_id").(string)

	// To delete all custom properties, we send an empty custom definitions object
	schema := &sdk.UserSchema{
		Definitions: &sdk.UserSchemaDefinitions{
			Custom: &sdk.UserSchemaPublic{
				Id:         "#custom",
				Type:       "object",
				Properties: make(map[string]*sdk.UserSchemaAttribute),
			},
		},
	}

	_, resp, err := getOktaClientFromMetadata(meta).UserSchema.UpdateApplicationUserProfile(ctx, appId, *schema)
	if err != nil {
		if utils.SuppressErrorOn404(resp, err) == nil {
			return nil
		}
		return diag.Errorf("failed to delete application user schema: %v", err)
	}

	return nil
}

// flattenSchemaProperty converts an SDK UserSchemaAttribute to a map for Terraform state
func flattenSchemaProperty(index string, attr *sdk.UserSchemaAttribute) map[string]interface{} {
	m := make(map[string]interface{})

	m["index"] = index

	if attr.Title != "" {
		m["title"] = attr.Title
	}
	if attr.Type != "" {
		m["type"] = attr.Type
	}
	if attr.Description != "" {
		m["description"] = attr.Description
	}
	if attr.Required != nil {
		m["required"] = *attr.Required
	}
	if attr.Scope != "" {
		m["scope"] = attr.Scope
	}
	if attr.MinLengthPtr != nil {
		m["min_length"] = *attr.MinLengthPtr
	}
	if attr.MaxLengthPtr != nil {
		m["max_length"] = *attr.MaxLengthPtr
	}
	if attr.Enum != nil {
		m["enum"] = attr.Enum
	}
	if attr.ExternalName != "" {
		m["external_name"] = attr.ExternalName
	}
	if attr.ExternalNamespace != "" {
		m["external_namespace"] = attr.ExternalNamespace
	}
	if attr.Master != nil {
		if attr.Master.Type != "" {
			m["master"] = attr.Master.Type
		}
	}
	if attr.Permissions != nil {
		for _, perm := range attr.Permissions {
			if perm.Action != "" {
				m["permissions"] = perm.Action
				break
			}
		}
	}
	if attr.Union != "" {
		m["union"] = attr.Union == "ENABLE"
	}

	return m
}

// expandSchemaProperty converts a Terraform config map to an SDK UserSchemaAttribute
func expandSchemaProperty(m map[string]interface{}) *sdk.UserSchemaAttribute {
	attr := &sdk.UserSchemaAttribute{}

	if v, ok := m["title"]; ok {
		attr.Title = v.(string)
	}
	if v, ok := m["type"]; ok {
		attr.Type = v.(string)
	}
	if v, ok := m["description"]; ok {
		attr.Description = v.(string)
	}
	if v, ok := m["required"]; ok {
		required := v.(bool)
		attr.Required = &required
	}
	if v, ok := m["scope"]; ok {
		attr.Scope = v.(string)
	}
	if v, ok := m["min_length"]; ok {
		minLen := int64(v.(int))
		attr.MinLengthPtr = &minLen
	}
	if v, ok := m["max_length"]; ok {
		maxLen := int64(v.(int))
		attr.MaxLengthPtr = &maxLen
	}
	if v, ok := m["enum"]; ok {
		attr.Enum = v.([]interface{})
	}
	if v, ok := m["external_name"]; ok {
		attr.ExternalName = v.(string)
	}
	if v, ok := m["external_namespace"]; ok {
		attr.ExternalNamespace = v.(string)
	}
	if v, ok := m["master"]; ok {
		attr.Master = &sdk.UserSchemaAttributeMaster{
			Type: v.(string),
		}
	}
	if v, ok := m["permissions"]; ok {
		attr.Permissions = []*sdk.UserSchemaAttributePermission{
			{
				Action: v.(string),
			},
		}
	}
	if v, ok := m["union"]; ok {
		if v.(bool) {
			attr.Union = "ENABLE"
		} else {
			attr.Union = "DISABLE"
		}
	}

	return attr
}

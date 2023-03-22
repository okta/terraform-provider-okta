package okta

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
)

var emailCustomizationsDataSourceSchema = map[string]*schema.Schema{
	"email_customizations": {
		Type:        schema.TypeSet,
		Computed:    true,
		Description: "List of `okta_email_customization` belonging to the named email template of the brand in the organization",
		Elem: &schema.Resource{
			Schema: emailCustomizationDataSourceSchema,
		},
		Set: hashEmailCustomization,
	},
}

var emailCustomizationDataSourceSchema = map[string]*schema.Schema{
	"id": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The ID of the customization",
	},
	"links": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Link relations for this object - JSON HAL - Discoverable resources related to the email template",
	},
	"language": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The language supported by the customization",
	},
	"is_default": {
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Whether the customization is the default",
	},
	"subject": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The subject of the customization",
	},
	"body": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The body of the customization",
	},
}

var emailCustomizationResourceSchema = map[string]*schema.Schema{
	"id": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The ID of the customization",
	},
	"brand_id": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Brand ID",
	},
	"template_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Template Name",
	},
	"links": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Link relations for this object - JSON HAL - Discoverable resources related to the email template",
	},
	"language": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The language supported by the customization",
	},
	"is_default": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Whether the customization is the default",
	},
	"subject": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The subject of the customization",
	},
	"body": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The body of the customization",
	},
	"force_is_default": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Force is_default on the create and delete by deleting all email customizations. Comma separated string with values of 'create' or 'destroy' or both `create,destroy'.",
	},
}

func flattenEmailCustomization(emailCustomization *sdk.EmailTemplateCustomization) map[string]interface{} {
	attrs := map[string]interface{}{}
	attrs["id"] = emailCustomization.Id
	attrs["language"] = emailCustomization.Language
	attrs["is_default"] = false
	if emailCustomization.IsDefault != nil {
		attrs["is_default"] = emailCustomization.IsDefault
	}
	attrs["subject"] = emailCustomization.Subject
	attrs["body"] = emailCustomization.Body
	links, _ := json.Marshal(emailCustomization.Links)
	attrs["links"] = string(links)

	return attrs
}

func hashEmailCustomization(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf(
		"%s-%s-%s-",
		m["id"].(string),
		m["language"].(string),
		m["subject"].(string),
	))
	return schema.HashString(buf.String())
}

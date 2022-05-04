package okta

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

var emailCustomizationsDataSourceSchema = map[string]*schema.Schema{
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

func flattenEmailCustomization(brandId, templateName string, emailCustomization *okta.EmailTemplateCustomization) map[string]interface{} {
	attrs := map[string]interface{}{}
	attrs["id"] = emailCustomization.Id
	attrs["brand_id"] = brandId
	attrs["template_name"] = templateName
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
		"%s-%s-%s-%s-",
		m["brand_id"].(string),
		m["template_name"].(string),
		m["language"].(string),
		m["subject"].(string),
	))
	return schema.HashString(buf.String())
}

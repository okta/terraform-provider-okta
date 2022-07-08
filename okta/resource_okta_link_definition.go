package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func resourceLinkDefinition() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLinkDefinitionCreate,
		ReadContext:   resourceLinkDefinitionRead,
		DeleteContext: resourceLinkDefinitionDelete,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
		Schema: map[string]*schema.Schema{
			"primary_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "API name of the primary link.",
				ForceNew:    true,
			},
			"primary_title": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name of the primary link.",
				ForceNew:    true,
			},
			"primary_description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Description of the primary relationship.",
				ForceNew:    true,
			},
			"associated_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "API name of the associated link.",
				ForceNew:    true,
			},
			"associated_title": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name of the associated link.",
				ForceNew:    true,
			},
			"associated_description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Description of the associated relationship.",
				ForceNew:    true,
			},
		},
	}
}

func resourceLinkDefinitionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	linkedObject := okta.LinkedObject{
		Primary: &okta.LinkedObjectDetails{
			Name:        d.Get("primary_name").(string),
			Title:       d.Get("primary_title").(string),
			Description: d.Get("primary_description").(string),
			Type:        "USER",
		},
		Associated: &okta.LinkedObjectDetails{
			Name:        d.Get("associated_name").(string),
			Title:       d.Get("associated_title").(string),
			Description: d.Get("associated_description").(string),
			Type:        "USER",
		},
	}
	_, _, err := getOktaClientFromMetadata(m).LinkedObject.AddLinkedObjectDefinition(ctx, linkedObject)
	if err != nil {
		return diag.Errorf("failed to create linked object: %v", err)
	}
	d.SetId(d.Get("primary_name").(string))
	return nil
}

func resourceLinkDefinitionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	linkedObject, resp, err := getOktaClientFromMetadata(m).LinkedObject.GetLinkedObjectDefinition(ctx, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get linked object: %v", err)
	}
	if linkedObject == nil {
		d.SetId("")
		return nil
	}
	if d.Id() != linkedObject.Primary.Name {
		d.SetId(linkedObject.Primary.Name)
	}
	_ = d.Set("primary_name", linkedObject.Primary.Name)
	_ = d.Set("primary_title", linkedObject.Primary.Title)
	_ = d.Set("primary_description", linkedObject.Primary.Description)
	_ = d.Set("associated_name", linkedObject.Associated.Name)
	_ = d.Set("associated_title", linkedObject.Associated.Title)
	_ = d.Set("associated_description", linkedObject.Associated.Description)
	return nil
}

func resourceLinkDefinitionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resp, err := getOktaClientFromMetadata(m).LinkedObject.DeleteLinkedObjectDefinition(ctx, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to remove linked object: %v", err)
	}
	return nil
}

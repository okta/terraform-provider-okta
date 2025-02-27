package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
)

func resourceLinkDefinition() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLinkDefinitionCreate,
		ReadContext:   resourceLinkDefinitionRead,
		DeleteContext: resourceLinkDefinitionDelete,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
		Description: `Manages the creation and removal of the link definitions.
		
Link definition operations allow you to manage the creation and removal of the link definitions. If you remove a link 
definition, links based on that definition are unavailable. Note that this resource is immutable, thus can not be modified.
~> **NOTE:** Links reappear if you recreate the definition. However, Okta is likely to change this behavior so that links don't reappear. Don't rely on this behavior in production environments.`,
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

func resourceLinkDefinitionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// NOTE: Okta API will ignore parallel calls to `POST
	// /api/v1/meta/schemas/user/linkedObjects` so a mutex to affect TF
	// `-parallelism=1` behavior is needed here.
	oktaMutexKV.Lock(resources.OktaIDaaSLinkDefinition)
	defer oktaMutexKV.Unlock(resources.OktaIDaaSLinkDefinition)

	linkedObject := sdk.LinkedObject{
		Primary: &sdk.LinkedObjectDetails{
			Name:        d.Get("primary_name").(string),
			Title:       d.Get("primary_title").(string),
			Description: d.Get("primary_description").(string),
			Type:        "USER",
		},
		Associated: &sdk.LinkedObjectDetails{
			Name:        d.Get("associated_name").(string),
			Title:       d.Get("associated_title").(string),
			Description: d.Get("associated_description").(string),
			Type:        "USER",
		},
	}
	_, _, err := getOktaClientFromMetadata(meta).LinkedObject.AddLinkedObjectDefinition(ctx, linkedObject)
	if err != nil {
		return diag.Errorf("failed to create linked object: %v", err)
	}
	d.SetId(d.Get("primary_name").(string))
	return nil
}

func resourceLinkDefinitionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	linkedObject, resp, err := getOktaClientFromMetadata(meta).LinkedObject.GetLinkedObjectDefinition(ctx, d.Id())
	if err := utils.SuppressErrorOn404(resp, err); err != nil {
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

func resourceLinkDefinitionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// NOTE: Okta API will ignore parallel calls to `DELETE
	// /api/v1/meta/schemas/user/linkedObjects` so a mutex to affect TF
	// `-parallelism=1` behavior is needed here.
	oktaMutexKV.Lock(resources.OktaIDaaSLinkDefinition)
	defer oktaMutexKV.Unlock(resources.OktaIDaaSLinkDefinition)

	resp, err := getOktaClientFromMetadata(meta).LinkedObject.DeleteLinkedObjectDefinition(ctx, d.Id())
	if err := utils.SuppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to remove linked object: %v", err)
	}
	return nil
}

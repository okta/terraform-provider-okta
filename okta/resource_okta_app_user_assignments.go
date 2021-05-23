package okta

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func resourceOktaAppUserAssignments() *schema.Resource {
	return &schema.Resource{
		CreateContext: nil,
		ReadContext:   nil,
		UpdateContext: nil,
		DeleteContext: nil,
		Importer:      nil,
		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "App to associate user with",
			},
			"users": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Set of users to associate with the app",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "User associated with the application",
						},
						"username": {
							Type:     schema.TypeString,
							Required: true,
						},
						"password": {
							Type:      schema.TypeString,
							Sensitive: true,
							Optional:  true,
						},
						"profile": {
							Type:             schema.TypeString,
							ValidateDiagFunc: stringIsJSON,
							StateFunc:        normalizeDataJSON,
							Optional:         true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return new == ""
							},
						},
						"retain_assignment": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Retain the user assignment on destroy. If set to true, the resource will be removed from state but not from the Okta app.",
						},
					},
				},
			},
		},
	}
}

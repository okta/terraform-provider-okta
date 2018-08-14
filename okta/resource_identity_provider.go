package okta

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"log"
)

func resourceIdentityProviders() *schema.Resource {
	return &schema.Resource{
		Create: resourceIdentityProviderCreate,
		Read:   resourceIdentityProviderRead,
		Update: resourceIdentityProviderUpdate,
		Delete: resourceIdentityProviderDelete,
		CustomizeDiff: func(d *schema.ResourceDiff, v interface{}) error {},

		Schema: map[string]*schema.Schema{
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"GOOGLE"}, false),
				Description:  "Identity Provider Type: GOOGLE",
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the Identity Provider Resource",
			},
			"protocol": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Conditions that must be met during Policy Evaluation",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "IDP Protocol type to use - ie. OAUTH2",
						},
						"scopes": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of Group IDs to Include",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"credentials": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"client": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"client_id": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "",
												},
												"client_secret": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"policy": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"provisioning": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "",
									},
				          "profileMaster": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "",
									},
				          "groups": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
						            "action": {
													Type:        schema.TypeString,
													Optional:    true,
													Default:     "NONE"
													Description: "",
												},
											},
										},
				          },
				          "conditions": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"deprovisioned": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"action": {
																Type:        schema.TypeString,
																Optional:    true,
																Default:     "NONE"
																Description: "",
															},
														},
													},
													"suspended": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"action": {
																Type:        schema.TypeString,
																Optional:    true,
																Default:     "NONE"
																Description: "",
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceIdentityProviderCreate(d *schema.ResourceData, m interface{}) error {}

func resourceIdentityProviderRead(d *schema.ResourceData, m interface{}) error {}
func resourceIdentityProviderUpdate(d *schema.ResourceData, m interface{}) error {}
func resourceIdentityProviderDelete(d *schema.ResourceData, m interface{}) error {}

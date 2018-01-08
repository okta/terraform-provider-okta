package main

import (
    "net/http"
    "time"

    okta "github.com/curtisallen/go-okta"
    "github.com/hashicorp/terraform/terraform"
    "github.com/hashicorp/terraform/helper/schema"
)

// Provider entry point for okta provider
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"organization": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Okta orgizantion id e.g. dev-1234",
			},

			"token": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OKTA_TOKEN", nil),
				Description: "Okta API token",
			},

			"preview": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
				DefaultFunc: func() (interface{}, error) {
					return false, nil
				},
				Description: "Okta API token preview",
			},
		},

//		ResourcesMap: map[string]*schema.Resource{
//			"okta_group":      resourceGroup(),
//			"okta_membership": resourceMembership(),
//		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	// configure the okta client
	client := okta.NewClient(
		d.Get("token").(string),
		d.Get("organization").(string),
		d.Get("preview").(bool),
		&http.Client{Timeout: 60 * time.Second},
	)

	return client, nil
}

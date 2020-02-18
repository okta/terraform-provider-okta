package okta

import (
	"strings"

	"github.com/crewjam/saml"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func syncSamlIndexEndpointBinding(d *schema.ResourceData, services []saml.IndexedEndpoint) {
	// Always grab the last one just for simplicity. Should never have duplicates.
	for _, service := range services {
		switch service.Binding {
		case postBinding:
			d.Set("http_post_binding", service.Location)
		case redirectBinding:
			d.Set("http_redirect_binding", service.Location)
		}
	}
}

func syncSamlEndpointBinding(d *schema.ResourceData, services []saml.Endpoint) {
	// Always grab the last one just for simplicity. Should never have duplicates.
	for _, service := range services {
		switch service.Binding {
		case postBinding:
			d.Set("http_post_binding", service.Location)
		case redirectBinding:
			d.Set("http_redirect_binding", service.Location)
		}
	}
}

func getExternalID(url string, pattern string) string {
	// Default idp issuer is such that I can extract the ID. If someone enters a custom value
	// this will result in "" most likely, which seems fine
	pur := strings.Replace(pattern, "${org.externalKey}", "", -1)
	return strings.Replace(url, pur, "", -1)
}

func syncSamlCertificates(d *schema.ResourceData, descriptors []saml.KeyDescriptor) {
	for _, desc := range descriptors {
		switch desc.Use {
		case "encryption":
			d.Set("encryption_certificate", desc.KeyInfo.Certificate)
		case "signing":
			d.Set("signing_certificate", desc.KeyInfo.Certificate)
		}
	}
}

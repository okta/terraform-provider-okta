package okta

import (
	"strings"

	"github.com/crewjam/saml"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func syncSamlIndexEndpointBinding(d *schema.ResourceData, services []saml.IndexedEndpoint) {
	// Always grab the last one just for simplicity. Should never have duplicates.
	for _, service := range services {
		switch service.Binding {
		case postBinding:
			_ = d.Set("http_post_binding", service.Location)
		case redirectBinding:
			_ = d.Set("http_redirect_binding", service.Location)
		}
	}
}

func syncSamlEndpointBinding(d *schema.ResourceData, services []saml.Endpoint) {
	// Always grab the last one just for simplicity. Should never have duplicates.
	for _, service := range services {
		switch service.Binding {
		case postBinding:
			_ = d.Set("http_post_binding", service.Location)
		case redirectBinding:
			_ = d.Set("http_redirect_binding", service.Location)
		}
	}
}

func getExternalID(url, pattern string) string {
	// Default idp issuer is such that I can extract the ID. If someone enters a custom value
	// this will result in "" most likely, which seems fine
	pur := strings.ReplaceAll(pattern, "${org.externalKey}", "")
	return strings.ReplaceAll(url, pur, "")
}

func syncSamlCertificates(d *schema.ResourceData, descriptors []saml.KeyDescriptor) {
	for _, desc := range descriptors {
		switch desc.Use {
		case "encryption":
			_ = d.Set("encryption_certificate", desc.KeyInfo.Certificate)
		case "signing":
			_ = d.Set("signing_certificate", desc.KeyInfo.Certificate)
		}
	}
}

package idaas

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
)

func TestBuildSamlAppAssertionLifetime(t *testing.T) {
	raw := map[string]interface{}{
		"label":                           "test",
		"sso_url":                         "https://example.com/sso",
		"recipient":                       "https://example.com/sso",
		"destination":                     "https://example.com/sso",
		"audience":                        "https://example.com/audience",
		"subject_name_id_template":        "${user.userName}",
		"subject_name_id_format":          "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified",
		"signature_algorithm":             "RSA_SHA256",
		"digest_algorithm":                "SHA256",
		"authn_context_class_ref":         "urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport",
		"response_signed":                 true,
		"saml_assertion_lifetime_seconds": 7200,
	}
	d := schema.TestResourceDataRaw(t, resourceAppSaml().Schema, raw)
	app, err := buildSamlApp(d)
	if err != nil {
		t.Fatal(err)
	}
	got := app.Settings.SignOn.SamlAssertionLifetimeSeconds
	if got == nil || *got != 7200 {
		t.Fatalf("expected samlAssertionLifetimeSeconds 7200, got %v", got)
	}

	delete(raw, "saml_assertion_lifetime_seconds")
	d = schema.TestResourceDataRaw(t, resourceAppSaml().Schema, raw)
	app, err = buildSamlApp(d)
	if err != nil {
		t.Fatal(err)
	}
	if app.Settings.SignOn.SamlAssertionLifetimeSeconds != nil {
		t.Fatalf("expected samlAssertionLifetimeSeconds to be omitted when unset, got %d", *app.Settings.SignOn.SamlAssertionLifetimeSeconds)
	}
}

func TestSetSamlSettingsAssertionLifetime(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceAppSaml().Schema, map[string]interface{}{})
	lifetime := int64(1800)
	if err := setSamlSettings(d, &sdk.SamlApplicationSettingsSignOn{SamlAssertionLifetimeSeconds: &lifetime}); err != nil {
		t.Fatal(err)
	}
	if got := d.Get("saml_assertion_lifetime_seconds").(int); got != 1800 {
		t.Fatalf("expected saml_assertion_lifetime_seconds 1800, got %d", got)
	}
}

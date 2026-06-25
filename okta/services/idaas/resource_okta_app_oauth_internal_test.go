package idaas

import (
	"testing"

	v6okta "github.com/okta/okta-sdk-golang/v6/okta"
)

// TestSetOAuthAppType verifies that the "type" attribute is populated correctly
// on read/import, including the preconfigured OIN OIDC app cases that regressed
// in GH-2868 (nil oauthClient or empty application_type). The key invariant is
// that "type" is Required + ForceNew, so it must never be clobbered with an
// empty value, otherwise a post-import plan destroys and recreates the live app.
func TestSetOAuthAppType(t *testing.T) {
	newClient := func(appType *string) *v6okta.OpenIdConnectApplicationSettingsClient {
		c := v6okta.NewOpenIdConnectApplicationSettingsClientWithDefaults()
		c.ApplicationType = appType
		return c
	}
	strPtr := func(s string) *string { return &s }

	tests := []struct {
		name string
		// existing is the value already in state/config before read (mimics an
		// import where the user has declared type, or a prior read).
		existing    string
		oauthClient *v6okta.OpenIdConnectApplicationSettingsClient
		want        string
	}{
		{
			name:        "concrete application_type is set",
			existing:    "",
			oauthClient: newClient(strPtr("web")),
			want:        "web",
		},
		{
			name:        "concrete application_type overwrites existing",
			existing:    "native",
			oauthClient: newClient(strPtr("browser")),
			want:        "browser",
		},
		{
			// Preconfigured OIN app: API omits application_type. We must keep the
			// declared/existing value rather than clobber it to "" (GH-2868).
			name:        "empty application_type preserves existing",
			existing:    "web",
			oauthClient: newClient(nil),
			want:        "web",
		},
		{
			// Preconfigured OIN app: API omits oauthClient entirely.
			name:        "nil oauthClient preserves existing",
			existing:    "web",
			oauthClient: nil,
			want:        "web",
		},
		{
			name:        "empty application_type with no existing stays empty",
			existing:    "",
			oauthClient: newClient(strPtr("")),
			want:        "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := resourceAppOAuth().TestResourceData()
			if tt.existing != "" {
				if err := d.Set("type", tt.existing); err != nil {
					t.Fatalf("failed to seed existing type: %v", err)
				}
			}

			setOAuthAppType(d, tt.oauthClient)

			if got := d.Get("type").(string); got != tt.want {
				t.Errorf("type = %q, want %q", got, tt.want)
			}
		})
	}
}

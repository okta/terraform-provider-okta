package idaas

import (
	"context"
	"encoding/json"
	"testing"

	ctyjson "github.com/hashicorp/go-cty/cty/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func oauthUpdateData(t *testing.T, stateVersion string, cfg map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := resourceAppOAuth()
	var state *terraform.InstanceState
	if stateVersion != "" {
		state = &terraform.InstanceState{
			ID: "0oa1",
			Attributes: map[string]string{
				"id":                             "0oa1",
				"omit_secret":                    "true",
				"client_basic_secret_wo_version": stateVersion,
			},
		}
	}
	// The raw config has to carry every attribute, as it does in a real apply,
	// so build it against the resource's schema.
	block := r.CoreConfigSchema()
	raw, err := json.Marshal(cfg)
	if err != nil {
		t.Fatal(err)
	}
	val, err := ctyjson.Unmarshal(raw, block.ImpliedType())
	if err != nil {
		t.Fatal(err)
	}
	diff, err := r.Diff(context.Background(), state, terraform.NewResourceConfigShimmed(val, block), nil)
	if err != nil {
		t.Fatal(err)
	}
	// The gRPC server attaches the raw config to the diff the same way.
	diff.RawConfig = val
	d, err := schema.InternalMap(r.Schema).Data(state, diff)
	if err != nil {
		t.Fatal(err)
	}
	return d
}

func oauthConfig(version int) map[string]interface{} {
	return map[string]interface{}{
		"label":                          "test",
		"type":                           "web",
		"grant_types":                    []interface{}{"authorization_code"},
		"response_types":                 []interface{}{"code"},
		"redirect_uris":                  []interface{}{"https://example.com/callback"},
		"omit_secret":                    true,
		"client_basic_secret_wo":         "new-secret",
		"client_basic_secret_wo_version": version,
	}
}

func clientSecretAfterOmit(t *testing.T, d *schema.ResourceData) string {
	t.Helper()
	app, err := buildAppOAuthV6(d, false)
	if err != nil {
		t.Fatal(err)
	}
	omitOAuthClientSecretOnUpdate(d, &app)
	return app.OpenIdConnectApplication.GetCredentials().OauthClient.GetClientSecret()
}

func TestOmitOAuthClientSecretOnUpdateKeepsRotatedWriteOnlySecret(t *testing.T) {
	d := oauthUpdateData(t, "1", oauthConfig(2))
	if got := clientSecretAfterOmit(t, d); got != "new-secret" {
		t.Fatalf("the write-only secret should be sent when its version changes, got %q", got)
	}
}

func TestOmitOAuthClientSecretOnUpdateStillOmitsOtherwise(t *testing.T) {
	// version unchanged: another attribute changed, do not re-send the secret
	cfg := oauthConfig(1)
	cfg["label"] = "renamed"
	d := oauthUpdateData(t, "1", cfg)
	if got := clientSecretAfterOmit(t, d); got != "" {
		t.Fatalf("the secret should be omitted when its version did not change, got %q", got)
	}
}

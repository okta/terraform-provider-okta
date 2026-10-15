package idaas

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	v6okta "github.com/okta/okta-sdk-golang/v6/okta"
	"github.com/stretchr/testify/require"
)

func TestBuildAppOAuthV6DefaultsToSigningUseWhenUseOmitted(t *testing.T) {
	jwks, err := buildOAuthAppJwksForTest(t, []interface{}{
		map[string]interface{}{
			"kty": "RSA",
			"kid": "SIGNING_KEY_RSA",
			"e":   "AQAB",
			"n":   "modulus",
		},
	})

	require.NoError(t, err)
	keys := jwks.GetKeys()
	require.Len(t, keys, 1)
	require.NotNil(t, keys[0].OAuth2ClientJsonSigningKeyResponse)
	require.NotNil(t, keys[0].OAuth2ClientJsonSigningKeyResponse.OAuth2ClientJsonWebKeyRsaResponse)
	require.Nil(t, keys[0].OAuth2ClientJsonEncryptionKeyResponse)
	require.Equal(t, "sig", keys[0].OAuth2ClientJsonSigningKeyResponse.OAuth2ClientJsonWebKeyRsaResponse.AdditionalProperties["use"])

	body, err := json.Marshal(jwks)
	require.NoError(t, err)
	require.JSONEq(t, `{"keys":[{"e":"AQAB","kid":"SIGNING_KEY_RSA","kty":"RSA","n":"modulus","use":"sig"}]}`, string(body))
}

func TestBuildAppOAuthV6SerializesSigningUse(t *testing.T) {
	jwks, err := buildOAuthAppJwksForTest(t, []interface{}{
		map[string]interface{}{
			"kty": "RSA",
			"kid": "SIGNING_KEY_RSA",
			"e":   "AQAB",
			"n":   "modulus",
			"use": "sig",
		},
	})

	require.NoError(t, err)
	keys := jwks.GetKeys()
	require.Len(t, keys, 1)
	signingKey := keys[0].OAuth2ClientJsonSigningKeyResponse
	require.NotNil(t, signingKey)
	require.Equal(t, "sig", signingKey.OAuth2ClientJsonWebKeyRsaResponse.AdditionalProperties["use"])

	body, err := json.Marshal(jwks)
	require.NoError(t, err)
	require.JSONEq(t, `{"keys":[{"e":"AQAB","kid":"SIGNING_KEY_RSA","kty":"RSA","n":"modulus","use":"sig"}]}`, string(body))
}

func TestBuildAppOAuthV6SerializesEcSigningUse(t *testing.T) {
	jwks, err := buildOAuthAppJwksForTest(t, []interface{}{
		map[string]interface{}{
			"kty": "EC",
			"kid": "SIGNING_KEY_EC",
			"x":   "x-coordinate",
			"y":   "y-coordinate",
			"use": "sig",
		},
	})

	require.NoError(t, err)
	body, err := json.Marshal(jwks)
	require.NoError(t, err)
	require.JSONEq(t, `{"keys":[{"kid":"SIGNING_KEY_EC","kty":"EC","use":"sig","x":"x-coordinate","y":"y-coordinate"}]}`, string(body))
}

func TestBuildAppOAuthV6SerializesEncryptionUse(t *testing.T) {
	jwks, err := buildOAuthAppJwksForTest(t, []interface{}{
		map[string]interface{}{
			"kty": "RSA",
			"kid": "ENCRYPTION_KEY_RSA",
			"e":   "AQAB",
			"n":   "modulus",
			"use": "enc",
		},
	})

	require.NoError(t, err)
	keys := jwks.GetKeys()
	require.Len(t, keys, 1)
	require.NotNil(t, keys[0].OAuth2ClientJsonEncryptionKeyResponse)
	require.Nil(t, keys[0].OAuth2ClientJsonSigningKeyResponse)
	require.Equal(t, "enc", keys[0].OAuth2ClientJsonEncryptionKeyResponse.GetUse())

	body, err := json.Marshal(jwks)
	require.NoError(t, err)
	require.JSONEq(t, `{"keys":[{"e":"AQAB","kid":"ENCRYPTION_KEY_RSA","kty":"RSA","n":"modulus","use":"enc"}]}`, string(body))
}

func TestBuildAppOAuthV6RejectsEcEncryptionKey(t *testing.T) {
	_, err := buildOAuthAppJwksForTest(t, []interface{}{
		map[string]interface{}{
			"kty": "EC",
			"kid": "ENCRYPTION_KEY_EC",
			"x":   "x",
			"y":   "y",
			"use": "enc",
		},
	})

	require.EqualError(t, err, `OAuth application JWKS encryption keys must use kty "RSA", got "EC"`)
}

func TestBuildAppOAuthV6RejectsUnsupportedSigningKeyType(t *testing.T) {
	_, err := buildOAuthAppJwksForTest(t, []interface{}{
		map[string]interface{}{
			"kty": "oct",
			"kid": "SIGNING_KEY_OCT",
			"use": "sig",
		},
	})

	require.EqualError(t, err, `OAuth application JWKS signing keys must use kty "RSA" or "EC", got "oct"`)
}

func buildOAuthAppJwksForTest(t *testing.T, jwksData []interface{}) (*v6okta.OpenIdConnectApplicationSettingsClientKeys, error) {
	t.Helper()

	d := resourceDataWithRawConfigForTest(t, map[string]interface{}{
		"label":                      "test",
		"type":                       "service",
		"response_types":             []interface{}{"token"},
		"grant_types":                []interface{}{"client_credentials"},
		"token_endpoint_auth_method": "private_key_jwt",
		"jwks":                       jwksData,
	})
	app, err := buildAppOAuthV6(d, true)
	if err != nil {
		return nil, err
	}

	oidc, err := verifyOidcAppTypeV6(app)
	if err != nil {
		return nil, err
	}
	if oidc.Settings.OauthClient == nil || oidc.Settings.OauthClient.Jwks == nil {
		return nil, errors.New("OAuth application JWKS was not set")
	}
	return oidc.Settings.OauthClient.Jwks, nil
}

func resourceDataWithRawConfigForTest(t *testing.T, raw map[string]interface{}) *schema.ResourceData {
	t.Helper()

	sm := schema.InternalMap(resourceAppOAuth().Schema)
	diff, err := sm.Diff(context.Background(), nil, terraform.NewResourceConfigRaw(raw), nil, nil, true)
	require.NoError(t, err)
	diff.RawConfig = cty.ObjectVal(map[string]cty.Value{
		"auto_key_rotation":        cty.NullVal(cty.Bool),
		"pkce_required":            cty.NullVal(cty.Bool),
		"dpop_bound_access_tokens": cty.NullVal(cty.Bool),
		"client_basic_secret_wo":   cty.NullVal(cty.String),
	})

	d, err := sm.Data(nil, diff)
	require.NoError(t, err)
	return d
}

func TestSetOAuthClientSettingsV6SetsSigningUse(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceAppOAuth().Schema, nil)

	diags := setOAuthClientSettingsV6(d, testOAuthClientSettingsWithSigningJwks())
	require.False(t, diags.HasError(), diags)

	jwks := d.Get("jwks").([]interface{})
	require.Len(t, jwks, 1)
	require.Equal(t, "sig", jwks[0].(map[string]interface{})["use"])
}

func TestSetOAuthClientSettingsV6NormalizesNullSigningUse(t *testing.T) {
	var jwks v6okta.OpenIdConnectApplicationSettingsClientKeys
	err := json.Unmarshal([]byte(`{"keys":[{"e":"AQAB","kid":"remote-key","kty":"RSA","n":"modulus","use":null}]}`), &jwks)
	require.NoError(t, err)

	oauthClient := v6okta.NewOpenIdConnectApplicationSettingsClientWithDefaults()
	oauthClient.SetJwks(jwks)

	d := schema.TestResourceDataRaw(t, resourceAppOAuth().Schema, nil)
	diags := setOAuthClientSettingsV6(d, oauthClient)
	require.False(t, diags.HasError(), diags)

	stateJwks := d.Get("jwks").([]interface{})
	require.Len(t, stateJwks, 1)
	require.Equal(t, "sig", stateJwks[0].(map[string]interface{})["use"])
}

func TestSetOAuthClientSettingsV6ReordersJwksToMatchConfig(t *testing.T) {
	raw := map[string]interface{}{
		"jwks": []interface{}{
			map[string]interface{}{
				"kid": "SIGNING_KEY_RSA",
				"kty": "RSA",
				"e":   "AQAB",
				"n":   "modulus",
				"use": "sig",
			},
			map[string]interface{}{
				"kid": "SIGNING_KEY_EC",
				"kty": "EC",
				"x":   "x-coordinate",
				"y":   "y-coordinate",
				"use": "sig",
			},
		},
	}
	d := schema.TestResourceDataRaw(t, resourceAppOAuth().Schema, raw)

	ecKey := v6okta.NewOAuth2ClientJsonWebKeyECResponseWithDefaults()
	ecKey.SetKid("SIGNING_KEY_EC")
	ecKey.SetKty("EC")
	ecKey.SetX("x-coordinate")
	ecKey.SetY("y-coordinate")
	ecSigning := v6okta.OAuth2ClientJsonWebKeyECResponseAsOAuth2ClientJsonSigningKeyResponse(ecKey)

	rsaKey := v6okta.NewOAuth2ClientJsonWebKeyRsaResponseWithDefaults()
	rsaKey.SetKid("SIGNING_KEY_RSA")
	rsaKey.SetKty("RSA")
	rsaKey.SetE("AQAB")
	rsaKey.SetN("modulus")
	rsaSigning := v6okta.OAuth2ClientJsonWebKeyRsaResponseAsOAuth2ClientJsonSigningKeyResponse(rsaKey)

	jwks := v6okta.NewOpenIdConnectApplicationSettingsClientKeysWithDefaults()
	jwks.SetKeys([]v6okta.OpenIdConnectApplicationSettingsClientKeysKeysInner{
		v6okta.OAuth2ClientJsonSigningKeyResponseAsOpenIdConnectApplicationSettingsClientKeysKeysInner(&ecSigning),
		v6okta.OAuth2ClientJsonSigningKeyResponseAsOpenIdConnectApplicationSettingsClientKeysKeysInner(&rsaSigning),
	})

	oauthClient := v6okta.NewOpenIdConnectApplicationSettingsClientWithDefaults()
	oauthClient.SetJwks(*jwks)

	diags := setOAuthClientSettingsV6(d, oauthClient)
	require.False(t, diags.HasError(), diags)

	stateJwks := d.Get("jwks").([]interface{})
	require.Len(t, stateJwks, 2)
	require.Equal(t, "SIGNING_KEY_RSA", stateJwks[0].(map[string]interface{})["kid"])
	require.Equal(t, "SIGNING_KEY_EC", stateJwks[1].(map[string]interface{})["kid"])
}

func TestSetOAuthClientSettingsV6KeepsApiOrderWhenJwksKidCannotBeUniquelyMatched(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceAppOAuth().Schema, nil)

	rsaKey := v6okta.NewOAuth2ClientJsonWebKeyRsaResponseWithDefaults()
	rsaKey.SetKid("DUPLICATE_KEY")
	rsaKey.SetKty("RSA")
	rsaKey.SetE("AQAB")
	rsaKey.SetN("modulus")
	rsaSigning := v6okta.OAuth2ClientJsonWebKeyRsaResponseAsOAuth2ClientJsonSigningKeyResponse(rsaKey)

	ecKey := v6okta.NewOAuth2ClientJsonWebKeyECResponseWithDefaults()
	ecKey.SetKid("DUPLICATE_KEY")
	ecKey.SetKty("EC")
	ecKey.SetX("x-coordinate")
	ecKey.SetY("y-coordinate")
	ecSigning := v6okta.OAuth2ClientJsonWebKeyECResponseAsOAuth2ClientJsonSigningKeyResponse(ecKey)

	jwks := v6okta.NewOpenIdConnectApplicationSettingsClientKeysWithDefaults()
	jwks.SetKeys([]v6okta.OpenIdConnectApplicationSettingsClientKeysKeysInner{
		v6okta.OAuth2ClientJsonSigningKeyResponseAsOpenIdConnectApplicationSettingsClientKeysKeysInner(&rsaSigning),
		v6okta.OAuth2ClientJsonSigningKeyResponseAsOpenIdConnectApplicationSettingsClientKeysKeysInner(&ecSigning),
	})

	oauthClient := v6okta.NewOpenIdConnectApplicationSettingsClientWithDefaults()
	oauthClient.SetJwks(*jwks)

	diags := setOAuthClientSettingsV6(d, oauthClient)
	require.False(t, diags.HasError(), diags)

	stateJwks := d.Get("jwks").([]interface{})
	require.Len(t, stateJwks, 2)
	require.Equal(t, "RSA", stateJwks[0].(map[string]interface{})["kty"])
	require.Equal(t, "EC", stateJwks[1].(map[string]interface{})["kty"])
}

func testOAuthClientSettingsWithSigningJwks() *v6okta.OpenIdConnectApplicationSettingsClient {
	rsaKey := v6okta.NewOAuth2ClientJsonWebKeyRsaResponseWithDefaults()
	rsaKey.SetKid("remote-key")
	rsaKey.SetKty("RSA")
	rsaKey.SetE("AQAB")
	rsaKey.SetN("modulus")

	signingKey := v6okta.OAuth2ClientJsonWebKeyRsaResponseAsOAuth2ClientJsonSigningKeyResponse(rsaKey)
	key := v6okta.OAuth2ClientJsonSigningKeyResponseAsOpenIdConnectApplicationSettingsClientKeysKeysInner(&signingKey)

	jwks := v6okta.NewOpenIdConnectApplicationSettingsClientKeysWithDefaults()
	jwks.SetKeys([]v6okta.OpenIdConnectApplicationSettingsClientKeysKeysInner{key})

	oauthClient := v6okta.NewOpenIdConnectApplicationSettingsClientWithDefaults()
	oauthClient.SetJwks(*jwks)
	return oauthClient
}

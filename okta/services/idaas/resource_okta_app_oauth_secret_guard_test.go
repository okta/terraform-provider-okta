package idaas

import (
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

type fakeAppOAuthSecretGetter struct {
	values map[string]interface{}
	woSet  bool
	woVal  string
}

func (f fakeAppOAuthSecretGetter) Get(key string) interface{} { return f.values[key] }

func (f fakeAppOAuthSecretGetter) GetOk(key string) (interface{}, bool) {
	v, ok := f.values[key]
	if !ok || v == "" || v == nil {
		return v, false
	}
	return v, true
}

func (f fakeAppOAuthSecretGetter) GetRawConfigAt(cty.Path) (cty.Value, diag.Diagnostics) {
	if !f.woSet {
		return cty.NullVal(cty.String), nil
	}
	return cty.StringVal(f.woVal), nil
}

func TestAppOAuthWouldSilentlyRotateSecret(t *testing.T) {
	tests := []struct {
		name string
		f    fakeAppOAuthSecretGetter
		want bool
	}{
		{
			name: "no secret anywhere, auth method requires one -> risky",
			f: fakeAppOAuthSecretGetter{values: map[string]interface{}{
				"omit_secret":                false,
				"token_endpoint_auth_method": "client_secret_basic",
				"client_secret":               "",
			}},
			want: true,
		},
		{
			name: "state has client_secret -> safe",
			f: fakeAppOAuthSecretGetter{values: map[string]interface{}{
				"omit_secret":                false,
				"token_endpoint_auth_method": "client_secret_basic",
				"client_secret":               "existing-secret",
			}},
			want: false,
		},
		{
			name: "omit_secret true -> safe regardless",
			f: fakeAppOAuthSecretGetter{values: map[string]interface{}{
				"omit_secret":                true,
				"token_endpoint_auth_method": "client_secret_basic",
				"client_secret":               "",
			}},
			want: false,
		},
		{
			name: "private_key_jwt does not need a secret -> safe",
			f: fakeAppOAuthSecretGetter{values: map[string]interface{}{
				"omit_secret":                false,
				"token_endpoint_auth_method": "private_key_jwt",
				"client_secret":               "",
			}},
			want: false,
		},
		{
			name: "none does not need a secret -> safe",
			f: fakeAppOAuthSecretGetter{values: map[string]interface{}{
				"omit_secret":                false,
				"token_endpoint_auth_method": "none",
				"client_secret":               "",
			}},
			want: false,
		},
		{
			name: "client_basic_secret set -> safe",
			f: fakeAppOAuthSecretGetter{values: map[string]interface{}{
				"omit_secret":                false,
				"token_endpoint_auth_method": "client_secret_post",
				"client_secret":               "",
				"client_basic_secret":         "user-provided",
			}},
			want: false,
		},
		{
			name: "client_basic_secret_wo set via raw config -> safe",
			f: fakeAppOAuthSecretGetter{
				values: map[string]interface{}{
					"omit_secret":                false,
					"token_endpoint_auth_method": "client_secret_jwt",
					"client_secret":               "",
				},
				woSet: true,
				woVal: "wo-secret",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := appOAuthWouldSilentlyRotateSecret(tt.f); got != tt.want {
				t.Errorf("appOAuthWouldSilentlyRotateSecret() = %v, want %v", got, tt.want)
			}
		})
	}
}

package config

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-hclog"
)

func TestConfigLoadAndValidate(t *testing.T) {
	tests := []struct {
		name         string
		accessToken  string
		apiToken     string
		clientID     string
		privateKey   string
		privateKeyID string
		scopes       []string
		expectError  bool
	}{
		{"access_token = pass", "accessToken", "", "", "", "", nil, false},
		// NOTE: don't test apiToken, it causes a hit to the wire with a "GET
		//       /api/v1/users/me" and the test tokens are scrubbed for this test
		// {"api_token = pass", "", "apiToken", "", "", "", nil, false},
		{"client_id, private_key, scopes = pass", "", "", "clientID", "privateKey", "", []string{"scope1", "scope2"}, false},
		{"client_id, private_key, private_key_id, scopes = pass", "", "", "clientID", "privateKey", "privateKeyID", []string{"scope1", "scope2"}, false},
	}

	for _, test := range tests {
		config := Config{
			OrgName:      "test",
			Domain:       "okta.com",
			AccessToken:  test.accessToken,
			ApiToken:     test.apiToken,
			ClientID:     test.clientID,
			PrivateKey:   test.privateKey,
			PrivateKeyId: test.privateKeyID,
			Scopes:       test.scopes,
			LogLevel:     int(hclog.Warn),
		}
		config.Logger = hclog.New(hclog.DefaultOptions)

		err := config.VerifyCredentials(context.TODO())
		if test.expectError && err == nil {
			t.Errorf("test %q: expected error but received none", test.name)
		}
		if !test.expectError && err != nil {
			t.Errorf("test %q: did not expect error but received error: %+v", test.name, err)
			fmt.Println()
		}
	}
}

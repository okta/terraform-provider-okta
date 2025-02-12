package okta

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
			orgName:      "test",
			domain:       "okta.com",
			accessToken:  test.accessToken,
			apiToken:     test.apiToken,
			clientID:     test.clientID,
			privateKey:   test.privateKey,
			privateKeyId: test.privateKeyID,
			scopes:       test.scopes,
			logLevel:     int(hclog.Warn),
		}
		config.logger = hclog.New(hclog.DefaultOptions)

		err := config.loadClients()
		if err == nil {
			err = config.verifyCredentials(context.TODO())
		}

		if test.expectError && err == nil {
			t.Errorf("test %q: expected error but received none", test.name)
		}
		if !test.expectError && err != nil {
			t.Errorf("test %q: did not expect error but received error: %+v", test.name, err)
			fmt.Println()
		}
	}
}

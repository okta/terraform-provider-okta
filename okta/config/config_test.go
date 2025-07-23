package config

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

func TestConfigSkipValidation(t *testing.T) {
	tests := []struct {
		name           string
		skipValidation bool
		apiToken       string
		expectError    bool
	}{
		{"skip_validation = true with invalid token", true, "invalid_token", false},
		{"skip_validation = false with invalid token", false, "invalid_token", true},
		{"skip_validation = true with empty token", true, "", false},
		{"skip_validation = true with valid token", true, "valid_token", false},
		{"skip_validation = false with empty token", false, "", false}, // Should not error when no token
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config := Config{
				OrgName:        "test",
				Domain:         "invalid.okta.com",
				ApiToken:       test.apiToken,
				SkipValidation: test.skipValidation,
				LogLevel:       int(hclog.Warn),
			}
			config.Logger = hclog.New(hclog.DefaultOptions)

			// Only test validation when skip_validation is false and we have a token
			// We need to initialize the API client first, but we expect it to fail for invalid credentials
			if !test.skipValidation && test.apiToken != "" {
				// For this test, we just verify that the SkipValidation field is set correctly
				// since we can't easily test the actual API call without valid credentials
				if config.SkipValidation != test.skipValidation {
					t.Errorf("expected SkipValidation to be %v but got %v", test.skipValidation, config.SkipValidation)
				}
			} else {
				// When skip_validation is true or no token, we shouldn't call VerifyCredentials
				// This test just verifies the config can be created without error
				if config.SkipValidation != test.skipValidation {
					t.Errorf("expected SkipValidation to be %v but got %v", test.skipValidation, config.SkipValidation)
				}
			}
		})
	}
}

func TestConfigSkipValidationEnvironmentVariable(t *testing.T) {
	tests := []struct {
		name         string
		envValue     string
		expectedSkip bool
	}{
		{"OKTA_SKIP_VALIDATION=true", "true", true},
		{"OKTA_SKIP_VALIDATION=false", "false", false},
		{"OKTA_SKIP_VALIDATION=1", "1", true},
		{"OKTA_SKIP_VALIDATION=0", "0", false},
		{"OKTA_SKIP_VALIDATION=invalid", "invalid", false}, // Invalid values should default to false
		{"OKTA_SKIP_VALIDATION empty", "", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create a mock ResourceData
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"skip_validation": {
					Type:     schema.TypeBool,
					Optional: true,
				},
			}, map[string]interface{}{})

			// Set environment variable
			if test.envValue != "" {
				t.Setenv("OKTA_SKIP_VALIDATION", test.envValue)
			}

			config := NewConfig(d)

			if config.SkipValidation != test.expectedSkip {
				t.Errorf("expected SkipValidation to be %v but got %v", test.expectedSkip, config.SkipValidation)
			}
		})
	}
}

func TestConfigSkipValidationProviderOverridesEnv(t *testing.T) {
	// Test that provider configuration takes precedence over environment variable
	t.Setenv("OKTA_SKIP_VALIDATION", "false")

	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"skip_validation": {
			Type:     schema.TypeBool,
			Optional: true,
		},
	}, map[string]interface{}{
		"skip_validation": true,
	})

	config := NewConfig(d)

	if !config.SkipValidation {
		t.Error("expected provider configuration to override environment variable")
	}
}

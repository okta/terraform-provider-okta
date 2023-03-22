package sdk

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/stretchr/testify/require"
)

func TestPolicyMarshal(t *testing.T) {
	example := SdkPolicy{}
	_json, err := json.Marshal(&example)
	require.NoError(t, err)
	require.Equal(t, `{}`, string(_json))

	example = SdkPolicy{
		Policy: okta.Policy{Id: "1"},
	}
	_json, err = json.Marshal(&example)
	require.NoError(t, err)
	require.Equal(t, `{"id":"1"}`, string(_json))

	example = SdkPolicy{
		Policy:   okta.Policy{Id: "1"},
		Settings: &SdkPolicySettings{},
	}
	_json, err = json.Marshal(&example)
	require.NoError(t, err)
	require.Equal(t, `{"id":"1","settings":{}}`, string(_json))

	example = SdkPolicy{
		Policy: okta.Policy{Id: "1"},
		Settings: &SdkPolicySettings{
			Type: "test",
		},
	}
	_json, err = json.Marshal(&example)
	require.NoError(t, err)
	require.Equal(t, `{"id":"1","settings":{"type":"test"}}`, string(_json))
}

func TestPolicyUnmarshal(t *testing.T) {
	_json := `{"id":"1","priority":7,"settings":{"type":"test"}}`

	var policy SdkPolicy
	err := json.Unmarshal([]byte(_json), &policy)
	require.NoError(t, err)
	require.NotNil(t, policy.Settings)
	require.Equal(t, policy.Settings.Type, "test")
}

// TestPolicyUnmarshalAdvanced sdk.Policy has okta.Policy embedded within it.
// Here, we are making sure the marshaling of the complex objects as presented
// by API responses is unmarshaled correctly
func TestPolicyUnmarshalAdvanced(t *testing.T) {
	pwd, err := os.Getwd()
	require.NoError(t, err)
	policyPath := filepath.Join(pwd, "../test/fixtures/password_policy_default.json")
	policyJSON, err := os.ReadFile(policyPath)
	require.NoError(t, err)

	var policy SdkPolicy
	err = json.Unmarshal(policyJSON, &policy)

	// make sure the marshaling from okta SDK is correct
	require.NoError(t, err)
	require.Equal(t, policy.Status, "ACTIVE")
	require.Equal(t, policy.Type, "PASSWORD")
	require.NotNil(t, policy.PriorityPtr)
	require.Equal(t, *policy.PriorityPtr, int64(2))

	// make sure marshaling of the local sdk's policy settings is correct
	require.NotNil(t, policy.Settings)
	require.NotNil(t, policy.Settings.Password)
	require.NotNil(t, policy.Settings.Password.Complexity)
	require.NotNil(t, policy.Settings.Password.Complexity.MinLengthPtr)
	require.Equal(t, *policy.Settings.Password.Complexity.MinLengthPtr, int64(8))
	require.NotNil(t, policy.Settings.Password.Age)
	require.NotNil(t, policy.Settings.Password.Age.HistoryCountPtr)
	require.Equal(t, *policy.Settings.Password.Age.HistoryCountPtr, int64(5))
	require.NotNil(t, policy.Settings.Password.Lockout)
	require.NotNil(t, policy.Settings.Password.Lockout.MaxAttemptsPtr)
	require.Equal(t, *policy.Settings.Password.Lockout.MaxAttemptsPtr, int64(10))
	require.NotNil(t, policy.Settings.Recovery)
	require.NotNil(t, policy.Settings.Recovery.Factors)
	require.NotNil(t, policy.Settings.Recovery.Factors.RecoveryQuestion)
	require.NotNil(t, policy.Settings.Recovery.Factors.RecoveryQuestion.Properties)
	require.NotNil(t, policy.Settings.Recovery.Factors.RecoveryQuestion.Properties.Complexity)
	require.NotNil(t, policy.Settings.Recovery.Factors.RecoveryQuestion.Properties.Complexity.MinLengthPtr)
	require.Equal(t, *policy.Settings.Recovery.Factors.RecoveryQuestion.Properties.Complexity.MinLengthPtr, int64(4))
	require.NotNil(t, policy.Settings.Delegation)
	require.NotNil(t, policy.Settings.Delegation.Options)
	require.NotNil(t, policy.Settings.Delegation.Options.SkipUnlock)
	require.Equal(t, *policy.Settings.Delegation.Options.SkipUnlock, false)
}

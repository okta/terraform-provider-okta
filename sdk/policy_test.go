package sdk

import (
	"encoding/json"
	"testing"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/stretchr/testify/require"
)

func TestPolicyMarshal(t *testing.T) {
	example := Policy{}
	_json, err := json.Marshal(&example)
	require.NoError(t, err)
	require.Equal(t, `{}`, string(_json))

	example = Policy{
		Policy: okta.Policy{Id: "1"},
	}
	_json, err = json.Marshal(&example)
	require.NoError(t, err)
	require.Equal(t, `{"id":"1"}`, string(_json))

	example = Policy{
		Policy:   okta.Policy{Id: "1"},
		Settings: &PolicySettings{},
	}
	_json, err = json.Marshal(&example)
	require.NoError(t, err)
	require.Equal(t, `{"id":"1","settings":{}}`, string(_json))

	example = Policy{
		Policy: okta.Policy{Id: "1"},
		Settings: &PolicySettings{
			Type: "test",
		},
	}
	_json, err = json.Marshal(&example)
	require.NoError(t, err)
	require.Equal(t, `{"id":"1","settings":{"type":"test"}}`, string(_json))
}

func TestPolicyUnmarshal(t *testing.T) {
	_json := `{"id":"1","priority":7,"settings":{"type":"test"}}`

	var policy Policy
	err := json.Unmarshal([]byte(_json), &policy)
	require.NoError(t, err)
	require.NotNil(t, policy.Settings)
	require.Equal(t, policy.Settings.Type, "test")
}

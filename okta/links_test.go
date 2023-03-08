package okta

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLinksValue(t *testing.T) {
	type links struct {
		Links interface{} `json:"_links,omitempty"`
	}
	tests := []struct {
		name      string
		linksJSON string
		keys      []string
		expected  string
	}{
		{
			name:      "missing links",
			linksJSON: `{}`,
			keys:      []string{"appLinks", "href"},
			expected:  "",
		},
		{
			name:      "nil links object",
			linksJSON: `{"_links":null}`,
			keys:      []string{"appLinks", "href"},
			expected:  "",
		},
		{
			name:      "empty links object",
			linksJSON: `{"_links":{}}`,
			keys:      []string{"appLinks", "href"},
			expected:  "",
		},
		{
			name:      "links object with nil appLinks",
			linksJSON: `{"_links":{"appLinks":null}}`,
			keys:      []string{"appLinks", "href"},
			expected:  "",
		},
		{
			// before fixing linksValue function this test case would cause the panic
			// "panic: runtime error: index out of range [0] with length 0 [recovered]""
			// seen in issue 1480
			// https://github.com/okta/terraform-provider-okta/issues/1480
			name:      "links object with empty appLinks array",
			linksJSON: `{"_links":{"appLinks":[]}}`,
			keys:      []string{"appLinks", "href"},
			expected:  "",
		},
		{
			name: "links object with an appLinks array",
			linksJSON: `{
				            "_links": {
				            	"appLinks":[
				            			{
				            				"name": "something",
				            				"href": "https://org/home/something",
				            				"type": "text/html"
				            			}						
				            		]
				                }
					}`,
			keys:     []string{"appLinks", "href"},
			expected: "https://org/home/something",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var _links links
			err := json.Unmarshal([]byte(test.linksJSON), &_links)
			require.NoError(t, err)
			result := linksValue(_links.Links, test.keys...)
			require.Equal(t, test.expected, result)
		})
	}
}

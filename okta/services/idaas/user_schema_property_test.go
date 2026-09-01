package idaas

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
	"github.com/stretchr/testify/require"
)

func TestBuildUserCustomSchemaAttributeDefaultStringEnum(t *testing.T) {
	d := userCustomSchemaResourceData(t, map[string]interface{}{
		"index":   "customSource",
		"title":   "Custom Source",
		"type":    "string",
		"default": "PRIMARY",
		"enum":    []interface{}{"PRIMARY", "SECONDARY"},
		"one_of": []interface{}{
			map[string]interface{}{"const": "PRIMARY", "title": "Primary"},
			map[string]interface{}{"const": "SECONDARY", "title": "Secondary"},
		},
	})

	attribute, err := buildUserCustomSchemaAttribute(d)
	require.NoError(t, err)
	require.Equal(t, "PRIMARY", attribute.Default)

	payload, err := json.Marshal(BuildCustomUserSchema(d.Get("index").(string), attribute))
	require.NoError(t, err)
	require.Contains(t, string(payload), `"default":"PRIMARY"`)
}

func TestBuildUserCustomSchemaAttributeDefaultCoercesPrimitiveTypes(t *testing.T) {
	tests := []struct {
		name          string
		attributeType string
		defaultValue  string
		expected      interface{}
	}{
		{
			name:          "string",
			attributeType: "string",
			defaultValue:  "PRIMARY",
			expected:      "PRIMARY",
		},
		{
			name:          "number",
			attributeType: "number",
			defaultValue:  "12.5",
			expected:      12.5,
		},
		{
			name:          "integer",
			attributeType: "integer",
			defaultValue:  "12",
			expected:      12,
		},
		{
			name:          "integer zero",
			attributeType: "integer",
			defaultValue:  "0",
			expected:      0,
		},
		{
			name:          "boolean",
			attributeType: "boolean",
			defaultValue:  "true",
			expected:      true,
		},
		{
			name:          "boolean false",
			attributeType: "boolean",
			defaultValue:  "false",
			expected:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := userCustomSchemaResourceData(t, map[string]interface{}{
				"index":   "customSource",
				"title":   "Custom Source",
				"type":    tt.attributeType,
				"default": tt.defaultValue,
			})

			attribute, err := buildUserCustomSchemaAttribute(d)
			require.NoError(t, err)
			require.Equal(t, tt.expected, attribute.Default)
		})
	}
}

func TestBuildUserCustomSchemaAttributeDefaultMarshalsZeroValues(t *testing.T) {
	tests := []struct {
		name          string
		attributeType string
		defaultValue  string
		expectedJSON  string
	}{
		{
			name:          "integer zero",
			attributeType: "integer",
			defaultValue:  "0",
			expectedJSON:  `"default":0`,
		},
		{
			name:          "boolean false",
			attributeType: "boolean",
			defaultValue:  "false",
			expectedJSON:  `"default":false`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := userCustomSchemaResourceData(t, map[string]interface{}{
				"index":   "customSource",
				"title":   "Custom Source",
				"type":    tt.attributeType,
				"default": tt.defaultValue,
			})

			attribute, err := buildUserCustomSchemaAttribute(d)
			require.NoError(t, err)

			payload, err := json.Marshal(BuildCustomUserSchema(d.Get("index").(string), attribute))
			require.NoError(t, err)
			require.Contains(t, string(payload), tt.expectedJSON)
		})
	}
}

func TestBuildUserCustomSchemaAttributeDefaultJSON(t *testing.T) {
	tests := []struct {
		name          string
		attributeType string
		defaultValue  string
		expected      interface{}
	}{
		{
			name:          "array",
			attributeType: "array",
			defaultValue:  `["PRIMARY","SECONDARY"]`,
			expected:      []interface{}{"PRIMARY", "SECONDARY"},
		},
		{
			name:          "object",
			attributeType: "object",
			defaultValue:  `{"source":"PRIMARY"}`,
			expected:      map[string]interface{}{"source": "PRIMARY"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := userCustomSchemaResourceData(t, map[string]interface{}{
				"index":   "customSource",
				"title":   "Custom Source",
				"type":    tt.attributeType,
				"default": tt.defaultValue,
			})

			attribute, err := buildUserCustomSchemaAttribute(d)
			require.NoError(t, err)
			require.Equal(t, tt.expected, attribute.Default)
		})
	}
}

func TestSyncCustomUserSchemaDefault(t *testing.T) {
	tests := []struct {
		name          string
		attributeType string
		defaultValue  interface{}
		expected      string
	}{
		{
			name:          "string",
			attributeType: "string",
			defaultValue:  "PRIMARY",
			expected:      "PRIMARY",
		},
		{
			name:          "number",
			attributeType: "number",
			defaultValue:  12.5,
			expected:      "12.5",
		},
		{
			name:          "integer",
			attributeType: "integer",
			defaultValue:  float64(12),
			expected:      "12",
		},
		{
			name:          "boolean",
			attributeType: "boolean",
			defaultValue:  true,
			expected:      "true",
		},
		{
			name:          "array",
			attributeType: "array",
			defaultValue:  []interface{}{"PRIMARY", "SECONDARY"},
			expected:      `["PRIMARY","SECONDARY"]`,
		},
		{
			name:          "object",
			attributeType: "object",
			defaultValue:  map[string]interface{}{"source": "PRIMARY"},
			expected:      `{"source":"PRIMARY"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := userCustomSchemaResourceData(t, map[string]interface{}{
				"index": "customSource",
				"title": "Custom Source",
				"type":  tt.attributeType,
			})

			err := syncCustomUserSchema(d, &sdk.UserSchemaAttribute{
				Title:       "Custom Source",
				Type:        tt.attributeType,
				Default:     tt.defaultValue,
				Required:    boolPtr(false),
				Permissions: []*sdk.UserSchemaAttributePermission{{Action: "READ_ONLY", Principal: "SELF"}},
			})
			require.NoError(t, err)
			require.Equal(t, tt.expected, d.Get("default"))
		})
	}
}

func userCustomSchemaResourceData(t *testing.T, raw map[string]interface{}) *schema.ResourceData {
	t.Helper()
	return schema.TestResourceDataRaw(t, resourceUserCustomSchemaProperty().Schema, raw)
}

func boolPtr(value bool) *bool {
	return &value
}

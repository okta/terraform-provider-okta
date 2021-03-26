package sdk

import (
	"context"
	"encoding/json"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

type (
	UserSchema struct {
		Schema      string                 `json:"$schema,omitempty"`
		Created     string                 `json:"created,omitempty"`
		Definitions *UserSchemaDefinitions `json:"definitions,omitempty"`
		ID          string                 `json:"id,omitempty"`
		LastUpdated string                 `json:"lastUpdated,omitempty"`
		Name        string                 `json:"name,omitempty"`
		Properties  *UserSchemaProperties  `json:"properties,omitempty"`
		Title       string                 `json:"title,omitempty"`
		Type        string                 `json:"type,omitempty"`
	}

	UserSchemaPermission struct {
		Action    string `json:"action,omitempty"`
		Principal string `json:"principal,omitempty"`
	}

	UserSchemaPropertyProfile struct {
		AllOf []*UserSchemaRef `json:"allOf,omitempty"`
	}

	UserSchemaDefinitions struct {
		Base   *UserSubSchemaProperties `json:"base,omitempty"`
		Custom *UserSubSchemaProperties `json:"custom,omitempty"`
	}

	UserSchemaItem struct {
		Enum  []string          `json:"enum,omitempty"`
		OneOf []*UserSchemaEnum `json:"oneOf,omitempty"`
		Type  string            `json:"type,omitempty"`
	}

	UserSchemaMaster struct {
		Type     string                     `json:"type,omitempty"`
		Priority []UserSchemaMasterPriority `json:"priority,omitempty"`
	}

	UserSchemaMasterPriority struct {
		Type  string `json:"type,omitempty"`
		Value string `json:"value,omitempty"`
	}

	UserSchemaEnum struct {
		Const string `json:"const,omitempty"`
		Title string `json:"title,omitempty"`
	}

	UserSubSchema struct {
		Description       string                  `json:"description,omitempty"`
		Enum              []string                `json:"enum,omitempty"`
		Format            string                  `json:"format,omitempty"`
		Items             *UserSchemaItem         `json:"items,omitempty"`
		Master            *UserSchemaMaster       `json:"master,omitempty"`
		MaxLength         *int                    `json:"maxLength,omitempty"`
		MinLength         *int                    `json:"minLength,omitempty"`
		Mutability        string                  `json:"mutability,omitempty"`
		OneOf             []*UserSchemaEnum       `json:"oneOf,omitempty"`
		Pattern           *string                 `json:"pattern,omitempty"`
		Permissions       []*UserSchemaPermission `json:"permissions,omitempty"`
		Required          *bool                   `json:"required,omitempty"`
		Scope             string                  `json:"scope,omitempty"`
		Title             string                  `json:"title,omitempty"`
		Type              string                  `json:"type,omitempty"`
		Union             string                  `json:"union,omitempty"`
		Unique            string                  `json:"unique,omitempty"`
		ExternalName      string                  `json:"externalName,omitempty"`
		ExternalNamespace string                  `json:"externalNamespace,omitempty"`
		IsLogin           bool                    `json:"-"`
	}

	UserSubSchemaProperties struct {
		ID         string                    `json:"id,omitempty"`
		Properties map[string]*UserSubSchema `json:"properties,omitempty"`
		Required   []interface{}             `json:"required,omitempty"`
		Type       string                    `json:"type,omitempty"`
	}

	UserSchemaProperties struct {
		Profile *UserSchemaPropertyProfile `json:"profile,omitempty"`
	}

	UserSchemaRef struct {
		Ref string `json:"$ref,omitempty"`
	}
)

// This is workaround for issue, when we should set 'pattern' to 'null' explicitly to use default Email Format for login
// but for other cases `null` causes 500 error
func (u *UserSubSchema) MarshalJSON() ([]byte, error) {
	type localIDX UserSubSchema
	m, err := json.Marshal((*localIDX)(u))
	if !u.IsLogin {
		return m, err
	}
	if err != nil {
		return nil, err
	}
	var a interface{}
	_ = json.Unmarshal(m, &a)
	b := a.(map[string]interface{})
	p := b["pattern"]
	if p == nil || p.(string) == "" {
		b["pattern"] = nil
	}
	ret, err := json.Marshal(b)
	return ret, err
}

func (m *ApiSupplement) DeleteUserSchemaProperty(ctx context.Context, schemaURL string, id string) (*okta.Response, error) {
	req, err := m.RequestExecutor.NewRequest("POST", schemaURL, getCustomUserSchema(id, nil))
	if err != nil {
		return nil, err
	}

	return m.RequestExecutor.Do(ctx, req, nil)
}

func (m *ApiSupplement) UpdateCustomUserSchemaProperty(ctx context.Context, schemaURL string, id string, schema *UserSubSchema) (*UserSchema, *okta.Response, error) {
	return m.UpdateUserSchemaProperty(ctx, schemaURL, getCustomUserSchema(id, schema))
}

func (m *ApiSupplement) UpdateBaseUserSchemaProperty(ctx context.Context, schemaURL string, id string, schema *UserSubSchema) (*UserSchema, *okta.Response, error) {
	return m.UpdateUserSchemaProperty(ctx, schemaURL, getBaseUserSchema(id, schema))
}

func (m *ApiSupplement) UpdateUserSchemaProperty(ctx context.Context, schemaURL string, schema *UserSchema) (*UserSchema, *okta.Response, error) {
	req, err := m.RequestExecutor.NewRequest("POST", schemaURL, schema)
	if err != nil {
		return nil, nil, err
	}

	fullSchema := &UserSchema{}
	resp, err := m.RequestExecutor.Do(ctx, req, fullSchema)
	return fullSchema, resp, err
}

func (m *ApiSupplement) GetUserSchema(ctx context.Context, schemaURL string) (*UserSchema, *okta.Response, error) {
	req, err := m.RequestExecutor.NewRequest("GET", schemaURL, nil)
	if err != nil {
		return nil, nil, err
	}
	schema := &UserSchema{}
	resp, err := m.RequestExecutor.Do(ctx, req, schema)
	return schema, resp, err
}

func getBaseUserSchema(index string, schema *UserSubSchema) *UserSchema {
	return &UserSchema{
		Definitions: &UserSchemaDefinitions{
			Base: GetUserSchemaProp("#base", index, schema),
		},
	}
}

func getCustomUserSchema(index string, schema *UserSubSchema) *UserSchema {
	return &UserSchema{
		Definitions: &UserSchemaDefinitions{
			Custom: GetUserSchemaProp("#custom", index, schema),
		},
	}
}

func GetUserSchemaProp(id, index string, schema *UserSubSchema) *UserSubSchemaProperties {
	return &UserSubSchemaProperties{
		ID:   id,
		Type: "object",
		Properties: map[string]*UserSubSchema{
			index: schema,
		},
	}
}

package sdk

import (
	"github.com/okta/okta-sdk-golang/okta"
)

var (
	schemaUrl = "/api/v1/meta/schemas/user/default"
)

func (m *ApiSupplement) RemoveUserSchemaProperty(id string) (*okta.Response, error) {
	req, err := m.RequestExecutor.NewRequest("DELETE", schemaUrl, nil)
	if err != nil {
		return nil, err
	}

	return m.RequestExecutor.Do(req, nil)
}

func (m *ApiSupplement) AddCustomUserSchemaProperty(schema *UserSubSchema) (*UserSchema, *okta.Response, error) {
	req, err := m.RequestExecutor.NewRequest("POST", schemaUrl, schema)
	if err != nil {
		return nil, nil, err
	}

	fullSchema := &UserSchema{}
	resp, err := m.RequestExecutor.Do(req, fullSchema)
	return fullSchema, resp, err
}

func (m *ApiSupplement) UpdateCustomUserSchemaProperty(id string, schema *UserSubSchema) (*UserSchema, *okta.Response, error) {
	return m.UpdateUserSchemaProperty(getUserSchema("#custom", id, schema))
}

func (m *ApiSupplement) UpdateBaseUserSchemaProperty(id string, schema *UserSubSchema) (*UserSchema, *okta.Response, error) {
	return m.UpdateUserSchemaProperty(getUserSchema("#base", id, schema))
}

func (m *ApiSupplement) UpdateUserSchemaProperty(schema *UserSchema) (*UserSchema, *okta.Response, error) {
	req, err := m.RequestExecutor.NewRequest("POST", schemaUrl, schema)
	if err != nil {
		return nil, nil, err
	}

	fullSchema := &UserSchema{}
	resp, err := m.RequestExecutor.Do(req, fullSchema)
	return fullSchema, resp, err
}

func (m *ApiSupplement) GetUserSchema() (*UserSchema, *okta.Response, error) {
	req, err := m.RequestExecutor.NewRequest("GET", schemaUrl, nil)
	if err != nil {
		return nil, nil, err
	}
	schema := &UserSchema{}
	resp, err := m.RequestExecutor.Do(req, schema)
	return schema, resp, err
}

func getUserSchema(id, index string, schema *UserSubSchema) *UserSchema {
	return &UserSchema{
		Definitions: &UserSchemaDefinitions{
			Custom: &UserSubSchemaProperties{
				ID:   id,
				Type: "object",
				Properties: map[string]*UserSubSchema{
					index: schema,
				},
			},
		},
	}
}

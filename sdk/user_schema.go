package sdk

import (
	"context"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

var (
	schemaUrl = "/api/v1/meta/schemas/user/default"
)

func (m *ApiSupplement) DeleteUserSchemaProperty(id string) (*okta.Response, error) {
	req, err := m.RequestExecutor.NewRequest("POST", schemaUrl, getCustomUserSchema(id, nil))
	if err != nil {
		return nil, err
	}

	return m.RequestExecutor.Do(context.Background(), req, nil)
}

func (m *ApiSupplement) AddCustomUserSchemaProperty(schema *UserSubSchema) (*UserSchema, *okta.Response, error) {
	req, err := m.RequestExecutor.NewRequest("POST", schemaUrl, schema)
	if err != nil {
		return nil, nil, err
	}

	fullSchema := &UserSchema{}
	resp, err := m.RequestExecutor.Do(context.Background(), req, fullSchema)
	return fullSchema, resp, err
}

func (m *ApiSupplement) UpdateCustomUserSchemaProperty(id string, schema *UserSubSchema) (*UserSchema, *okta.Response, error) {
	return m.UpdateUserSchemaProperty(getCustomUserSchema(id, schema))
}

func (m *ApiSupplement) UpdateBaseUserSchemaProperty(id string, schema *UserSubSchema) (*UserSchema, *okta.Response, error) {
	return m.UpdateUserSchemaProperty(getBaseUserSchema(id, schema))
}

func (m *ApiSupplement) UpdateUserSchemaProperty(schema *UserSchema) (*UserSchema, *okta.Response, error) {
	req, err := m.RequestExecutor.NewRequest("POST", schemaUrl, schema)
	if err != nil {
		return nil, nil, err
	}

	fullSchema := &UserSchema{}
	resp, err := m.RequestExecutor.Do(context.Background(), req, fullSchema)
	return fullSchema, resp, err
}

func (m *ApiSupplement) GetUserSchema() (*UserSchema, *okta.Response, error) {
	req, err := m.RequestExecutor.NewRequest("GET", schemaUrl, nil)
	if err != nil {
		return nil, nil, err
	}
	schema := &UserSchema{}
	resp, err := m.RequestExecutor.Do(context.Background(), req, schema)
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

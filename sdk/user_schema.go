package sdk

import (
	"context"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

func (m *ApiSupplement) DeleteUserSchemaProperty(schemaURL string, id string) (*okta.Response, error) {
	req, err := m.RequestExecutor.NewRequest("POST", schemaURL, getCustomUserSchema(id, nil))
	if err != nil {
		return nil, err
	}

	return m.RequestExecutor.Do(context.Background(), req, nil)
}

func (m *ApiSupplement) AddCustomUserSchemaProperty(schemaURL string, schema *UserSubSchema) (*UserSchema, *okta.Response, error) {
	req, err := m.RequestExecutor.NewRequest("POST", schemaURL, schema)
	if err != nil {
		return nil, nil, err
	}

	fullSchema := &UserSchema{}
	resp, err := m.RequestExecutor.Do(context.Background(), req, fullSchema)
	return fullSchema, resp, err
}

func (m *ApiSupplement) UpdateCustomUserSchemaProperty(schemaURL string, id string, schema *UserSubSchema) (*UserSchema, *okta.Response, error) {
	return m.UpdateUserSchemaProperty(schemaURL, getCustomUserSchema(id, schema))
}

func (m *ApiSupplement) UpdateBaseUserSchemaProperty(schemaURL string, id string, schema *UserSubSchema) (*UserSchema, *okta.Response, error) {
	return m.UpdateUserSchemaProperty(schemaURL, getBaseUserSchema(id, schema))
}

func (m *ApiSupplement) UpdateUserSchemaProperty(schemaURL string, schema *UserSchema) (*UserSchema, *okta.Response, error) {
	req, err := m.RequestExecutor.NewRequest("POST", schemaURL, schema)
	if err != nil {
		return nil, nil, err
	}

	fullSchema := &UserSchema{}
	resp, err := m.RequestExecutor.Do(context.Background(), req, fullSchema)
	return fullSchema, resp, err
}

func (m *ApiSupplement) GetUserSchema(schemaURL string) (*UserSchema, *okta.Response, error) {
	req, err := m.RequestExecutor.NewRequest("GET", schemaURL, nil)
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

package sdk

import (
	"github.com/okta/okta-sdk-golang/okta"
)

func (m *ApiSupplement) DeleteUserSchemaProperty(schemaUrl string, id string) (*okta.Response, error) {
	req, err := m.RequestExecutor.NewRequest("POST", schemaUrl, getCustomUserSchema(id, nil))
	if err != nil {
		return nil, err
	}

	return m.RequestExecutor.Do(req, nil)
}

func (m *ApiSupplement) AddCustomUserSchemaProperty(schemaUrl string, schema *UserSubSchema) (*UserSchema, *okta.Response, error) {
	req, err := m.RequestExecutor.NewRequest("POST", schemaUrl, schema)
	if err != nil {
		return nil, nil, err
	}

	fullSchema := &UserSchema{}
	resp, err := m.RequestExecutor.Do(req, fullSchema)
	return fullSchema, resp, err
}

func (m *ApiSupplement) UpdateCustomUserSchemaProperty(schemaUrl string, id string, schema *UserSubSchema) (*UserSchema, *okta.Response, error) {
	return m.UpdateUserSchemaProperty(schemaUrl, getCustomUserSchema(id, schema))
}

func (m *ApiSupplement) UpdateBaseUserSchemaProperty(schemaUrl string, id string, schema *UserSubSchema) (*UserSchema, *okta.Response, error) {
	return m.UpdateUserSchemaProperty(schemaUrl, getBaseUserSchema(id, schema))
}

func (m *ApiSupplement) UpdateUserSchemaProperty(schemaUrl string, schema *UserSchema) (*UserSchema, *okta.Response, error) {
	req, err := m.RequestExecutor.NewRequest("POST", schemaUrl, schema)
	if err != nil {
		return nil, nil, err
	}

	fullSchema := &UserSchema{}
	resp, err := m.RequestExecutor.Do(req, fullSchema)
	return fullSchema, resp, err
}

func (m *ApiSupplement) GetUserSchema(schemaUrl string) (*UserSchema, *okta.Response, error) {
	req, err := m.RequestExecutor.NewRequest("GET", schemaUrl, nil)
	if err != nil {
		return nil, nil, err
	}
	schema := &UserSchema{}
	resp, err := m.RequestExecutor.Do(req, schema)
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

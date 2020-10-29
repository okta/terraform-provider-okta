package sdk

import (
	"context"
	"fmt"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

var (
	appUserSchemaUrl = "/api/v1/meta/schemas/apps/%s/default"
)

func (m *ApiSupplement) UpdateAppUserSchema(appId string, schema *UserSchema) (*UserSchema, *okta.Response, error) {
	req, err := m.RequestExecutor.NewRequest("POST", fmt.Sprintf(appUserSchemaUrl, appId), schema)
	if err != nil {
		return nil, nil, err
	}

	fullSchema := &UserSchema{}
	resp, err := m.RequestExecutor.Do(context.Background(), req, fullSchema)
	return fullSchema, resp, err
}

func (m *ApiSupplement) GetAppUserSchema(appId string) (*UserSchema, *okta.Response, error) {
	req, err := m.RequestExecutor.NewRequest("GET", fmt.Sprintf(appUserSchemaUrl, appId), nil)
	if err != nil {
		return nil, nil, err
	}
	schema := &UserSchema{}
	resp, err := m.RequestExecutor.Do(context.Background(), req, schema)
	return schema, resp, err
}

func (m *ApiSupplement) DeleteAppUserSchemaProperty(id, appId string) (*okta.Response, error) {
	req, err := m.RequestExecutor.NewRequest("POST", fmt.Sprintf(appUserSchemaUrl, appId), getCustomUserSchema(id, nil))
	if err != nil {
		return nil, err
	}

	return m.RequestExecutor.Do(context.Background(), req, nil)
}

func (m *ApiSupplement) UpdateCustomAppUserSchemaProperty(id, appId string, schema *UserSubSchema) (*UserSchema, *okta.Response, error) {
	return m.UpdateAppUserSchema(appId, getCustomUserSchema(id, schema))
}

func (m *ApiSupplement) UpdateBaseAppUserSchemaProperty(id, appId string, schema *UserSubSchema) (*UserSchema, *okta.Response, error) {
	return m.UpdateAppUserSchema(appId, getBaseUserSchema(id, schema))
}

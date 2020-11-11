package sdk

import (
	"context"
	"fmt"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

var (
	appUserSchemaURL = "/api/v1/meta/schemas/apps/%s/default"
)

func (m *ApiSupplement) UpdateAppUserSchema(appID string, schema *UserSchema) (*UserSchema, *okta.Response, error) {
	req, err := m.RequestExecutor.NewRequest("POST", fmt.Sprintf(appUserSchemaURL, appID), schema)
	if err != nil {
		return nil, nil, err
	}

	fullSchema := &UserSchema{}
	resp, err := m.RequestExecutor.Do(context.Background(), req, fullSchema)
	return fullSchema, resp, err
}

func (m *ApiSupplement) GetAppUserSchema(appID string) (*UserSchema, *okta.Response, error) {
	req, err := m.RequestExecutor.NewRequest("GET", fmt.Sprintf(appUserSchemaURL, appID), nil)
	if err != nil {
		return nil, nil, err
	}
	schema := &UserSchema{}
	resp, err := m.RequestExecutor.Do(context.Background(), req, schema)
	return schema, resp, err
}

func (m *ApiSupplement) DeleteAppUserSchemaProperty(id, appID string) (*okta.Response, error) {
	req, err := m.RequestExecutor.NewRequest("POST", fmt.Sprintf(appUserSchemaURL, appID), getCustomUserSchema(id, nil))
	if err != nil {
		return nil, err
	}

	return m.RequestExecutor.Do(context.Background(), req, nil)
}

func (m *ApiSupplement) UpdateCustomAppUserSchemaProperty(id, appID string, schema *UserSubSchema) (*UserSchema, *okta.Response, error) {
	return m.UpdateAppUserSchema(appID, getCustomUserSchema(id, schema))
}

func (m *ApiSupplement) UpdateBaseAppUserSchemaProperty(id, appID string, schema *UserSubSchema) (*UserSchema, *okta.Response, error) {
	return m.UpdateAppUserSchema(appID, getBaseUserSchema(id, schema))
}

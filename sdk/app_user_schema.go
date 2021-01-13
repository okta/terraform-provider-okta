package sdk

import (
	"context"
	"fmt"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

var appUserSchemaURL = "/api/v1/meta/schemas/apps/%s/default"

func (m *ApiSupplement) UpdateAppUserSchema(ctx context.Context, appID string, schema *UserSchema) (*UserSchema, *okta.Response, error) {
	req, err := m.RequestExecutor.NewRequest("POST", fmt.Sprintf(appUserSchemaURL, appID), schema)
	if err != nil {
		return nil, nil, err
	}

	fullSchema := &UserSchema{}
	resp, err := m.RequestExecutor.Do(ctx, req, fullSchema)
	return fullSchema, resp, err
}

func (m *ApiSupplement) GetAppUserSchema(ctx context.Context, appID string) (*UserSchema, *okta.Response, error) {
	req, err := m.RequestExecutor.NewRequest("GET", fmt.Sprintf(appUserSchemaURL, appID), nil)
	if err != nil {
		return nil, nil, err
	}
	schema := &UserSchema{}
	resp, err := m.RequestExecutor.Do(ctx, req, schema)
	return schema, resp, err
}

func (m *ApiSupplement) DeleteAppUserSchemaProperty(ctx context.Context, id, appID string) (*okta.Response, error) {
	req, err := m.RequestExecutor.NewRequest("POST", fmt.Sprintf(appUserSchemaURL, appID), getCustomUserSchema(id, nil))
	if err != nil {
		return nil, err
	}

	return m.RequestExecutor.Do(ctx, req, nil)
}

func (m *ApiSupplement) UpdateCustomAppUserSchemaProperty(ctx context.Context, id, appID string, schema *UserSubSchema) (*UserSchema, *okta.Response, error) {
	return m.UpdateAppUserSchema(ctx, appID, getCustomUserSchema(id, schema))
}

func (m *ApiSupplement) UpdateBaseAppUserSchemaProperty(ctx context.Context, id, appID string, schema *UserSubSchema) (*UserSchema, *okta.Response, error) {
	return m.UpdateAppUserSchema(ctx, appID, getBaseUserSchema(id, schema))
}

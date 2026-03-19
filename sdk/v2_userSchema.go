// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

import (
	"context"
	"fmt"
)

type UserSchemaResource resource

type UserSchema struct {
	Schema      string                 `json:"$schema,omitempty"`
	Links       interface{}            `json:"_links,omitempty"`
	Created     string                 `json:"created,omitempty"`
	Definitions *UserSchemaDefinitions `json:"definitions,omitempty"`
	Id          string                 `json:"id,omitempty"`
	LastUpdated string                 `json:"lastUpdated,omitempty"`
	Name        string                 `json:"name,omitempty"`
	Properties  *UserSchemaProperties  `json:"properties,omitempty"`
	Title       string                 `json:"title,omitempty"`
	Type        string                 `json:"type,omitempty"`
}

// Fetches the Schema for an App User
func (m *UserSchemaResource) GetApplicationUserSchema(ctx context.Context, appInstanceId string) (*UserSchema, *Response, error) {
	url := fmt.Sprintf("/api/v1/meta/schemas/apps/%v/default", appInstanceId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var userSchema *UserSchema

	resp, err := rq.Do(ctx, req, &userSchema)
	if err != nil {
		return nil, resp, err
	}

	return userSchema, resp, nil
}

// Partial updates on the User Profile properties of the Application User Schema.
func (m *UserSchemaResource) UpdateApplicationUserProfile(ctx context.Context, appInstanceId string, body UserSchema) (*UserSchema, *Response, error) {
	url := fmt.Sprintf("/api/v1/meta/schemas/apps/%v/default", appInstanceId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	var userSchema *UserSchema

	resp, err := rq.Do(ctx, req, &userSchema)
	if err != nil {
		return nil, resp, err
	}

	return userSchema, resp, nil
}

// Fetches the schema for a Schema Id.
func (m *UserSchemaResource) GetUserSchema(ctx context.Context, schemaId string) (*UserSchema, *Response, error) {
	url := fmt.Sprintf("/api/v1/meta/schemas/user/%v", schemaId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var userSchema *UserSchema

	resp, err := rq.Do(ctx, req, &userSchema)
	if err != nil {
		return nil, resp, err
	}

	return userSchema, resp, nil
}

// Partial updates on the User Profile properties of the user schema.
func (m *UserSchemaResource) UpdateUserProfile(ctx context.Context, schemaId string, body UserSchema) (*UserSchema, *Response, error) {
	url := fmt.Sprintf("/api/v1/meta/schemas/user/%v", schemaId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	var userSchema *UserSchema

	resp, err := rq.Do(ctx, req, &userSchema)
	if err != nil {
		return nil, resp, err
	}

	return userSchema, resp, nil
}

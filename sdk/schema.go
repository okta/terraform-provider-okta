package sdk

import (
	"time"
)

type (
	SchemaDefinitions struct {
		Base   *SubSchema `json:"base"`
		Custom *SubSchema `json:"custom"`
	}

	SubSchema struct {
		ID         string           `json:"id"`
		Type       string           `json:"type"`
		Properties []*BaseSubSchema `json:"properties"`
		Required   []string         `json:"required"`
	}

	Items struct {
		Type string `json:"type,omitempty"`
	}

	Schema struct {
		Created     time.Time          `json:"created"`
		Definitions *SchemaDefinitions `json:"definitions"`
		ID          string             `json:"id"`
		LastUpdated time.Time          `json:"lastUpdated"`
		Name        string             `json:"name"`
		Schema      string             `json:"$schema"`
		Title       string             `json:"title"`
		Type        string             `json:"type"`
	}

	// User Profiles Base SubSchema
	BaseSubSchema struct {
		Format      string         `json:"format,omitempty"`
		Index       string         `json:"-"`
		Master      *Master        `json:"master,omitempty"`
		MaxLength   int            `json:"maxLength,omitempty"`
		MinLength   int            `json:"minLength,omitempty"`
		Mutability  string         `json:"mutablity,omitempty"`
		Permissions []*Permissions `json:"permissions"`
		Required    bool           `json:"required,omitempty"`
		Scope       string         `json:"scope,omitempty"`
		Title       string         `json:"title"`
		Type        string         `json:"type"`
	}

	// User Profiles Custom SubSchema
	CustomSubSchema struct {
		Description string         `json:"description,omitempty"`
		Enum        []string       `json:"enum,omitempty"`
		Format      string         `json:"format,omitempty"`
		Index       string         `json:"-"`
		Items       *Items         `json:"items,omitempty"`
		Master      *Master        `json:"master,omitempty"`
		MaxLength   int            `json:"maxLength,omitempty"`
		MinLength   int            `json:"minLength,omitempty"`
		Mutability  string         `json:"mutablity,omitempty"`
		OneOf       []*OneOf       `json:"oneOf,omitempty"`
		Permissions []*Permissions `json:"permissions"`
		Required    bool           `json:"required,omitempty"`
		Scope       string         `json:"scope,omitempty"`
		Title       string         `json:"title"`
		Type        string         `json:"type"`
		Union       string         `json:"union,omitempty"`
	}

	Master struct {
		Type string `json:"type,omitempty"`
	}

	// Permissions obj for User Profiles SubSchemas
	Permissions struct {
		Action    string `json:"action"`
		Principal string `json:"principal"`
	}

	// OneOf obj for User Profiles Custom SubSchema
	OneOf struct {
		Const string `json:"const"`
		Title string `json:"title"`
	}
)

package okta

import (
	"encoding/json"
	"fmt"
	"time"
)

// SchemasService handles communication with the Schema data related
// methods of the OKTA API.
type SchemasService service

// Return the BaseSubSchema object. Used to update the Base User SubSchema
func (p *SchemasService) BaseSubSchema() BaseSubSchema {
	return BaseSubSchema{}
}

// Return the CustomSubSchema object. Used to create & update Custom User SubSchema
func (p *SchemasService) CustomSubSchema() CustomSubSchema {
	return CustomSubSchema{}
}

// Return the Permissions object. Used to create & update User SubSchemas Permissions
func (p *SchemasService) Permissions() Permissions {
	return Permissions{}
}

// Return the OneOf object. Used to create & update Custom User SubSchema OneOf
func (p *SchemasService) OneOf() OneOf {
	return OneOf{}
}

// User Profiles Schema obj
type Schema struct {
	ID          string    `json:"id"`
	Schema      string    `json:"$schema"`
	Name        string    `json:"name"`
	Title       string    `json:"title"`
	Created     time.Time `json:"created"`
	LastUpdated time.Time `json:"lastUpdated"`
	Definitions struct {
		Base struct {
			ID         string          `json:"id"`
			Type       string          `json:"type"`
			Properties []BaseSubSchema `json:"properties"`
			Required   []string        `json:"required"`
		}
		Custom struct {
			ID         string            `json:"id"`
			Type       string            `json:"type"`
			Properties []CustomSubSchema `json:"properties"`
			Required   []string          `json:"required"`
		} `json:"custom"`
	} `json:"definitions"`
	Type string `json:"type"`
}

// User Profiles Base SubSchema
type BaseSubSchema struct {
	Index       string        `json:"-"`
	Title       string        `json:"title"`
	Type        string        `json:"type"`
	Format      string        `json:"format,omitempty"`
	Required    bool          `json:"required,omitempty"`
	Mutability  string        `json:"mutablity,omitempty"`
	Scope       string        `json:"scope,omitempty"`
	MinLength   int           `json:"minLength,omitempty"`
	MaxLength   int           `json:"maxLength,omitempty"`
	Permissions []Permissions `json:"permissions"`
	Master      *Master       `json:"master,omitempty"`
}

// User Profiles Custom SubSchema
type CustomSubSchema struct {
	Index       string `json:"-"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
	Format      string `json:"format,omitempty"`
	Required    bool   `json:"required,omitempty"`
	Mutability  string `json:"mutablity,omitempty"`
	Scope       string `json:"scope,omitempty"`
	MinLength   int    `json:"minLength,omitempty"`
	MaxLength   int    `json:"maxLength,omitempty"`
	Items       struct {
		Type string `json:"type,omitempty"`
	} `json:"items,omitempty"`
	Union       string        `json:"union,omitempty"`
	Enum        []string      `json:"enum,omitempty"`
	OneOf       []OneOf       `json:"oneOf,omitempty"`
	Permissions []Permissions `json:"permissions"`
	Master      *Master       `json:"master,omitempty"`
}

type Master struct {
	Type string `json:"type,omitempty"`
}

// Permissions obj for User Profiles SubSchemas
type Permissions struct {
	Principal string `json:"principal"`
	Action    string `json:"action"`
}

// OneOf obj for User Profiles Custom SubSchema
type OneOf struct {
	Const string `json:"const"`
	Title string `json:"title"`
}

// GetRawUserSchema returns the User Profile Schema as a map[string]interface{}
func (s *SchemasService) GetRawUserSchema() (map[string]interface{}, *Response, error) {
	u := fmt.Sprintf("meta/schemas/user/default")
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	var obj map[string]interface{}
	resp, err := s.client.Do(req, &obj)
	if err != nil {
		return nil, resp, err
	}
	return obj, resp, err
}

// GetUserSchema returns the User Profile Schema as a Schema struct
func (s *SchemasService) GetUserSchema() (*Schema, *Response, error) {
	obj, resp, err := s.client.Schemas.GetRawUserSchema()
	if err != nil {
		return nil, resp, err
	}
	schema, err := s.client.Schemas.userSchema(obj)
	if err != nil {
		return nil, resp, err
	}
	return schema, resp, nil
}

// userSchema (unexported) used to populate the Schema struct from a map[string]interface{}
// input is a map[string]interface{} of the User Profile Schema, such as from GetRawUserSchema()
func (s *SchemasService) userSchema(obj map[string]interface{}) (*Schema, error) {
	layout := "2006-01-02T15:04:05.000Z"
	create, err := time.Parse(layout, obj["created"].(string))
	if err != nil {
		return nil, err
	}
	update, _ := time.Parse(layout, obj["lastUpdated"].(string))
	if err != nil {
		return nil, err
	}
	schema := new(Schema)
	schema.ID = obj["id"].(string)
	schema.Schema = obj["$schema"].(string)
	schema.Name = obj["name"].(string)
	schema.Title = obj["title"].(string)
	schema.Created = create
	schema.LastUpdated = update
	schema.Type = obj["type"].(string)
	for k, v := range obj["definitions"].(map[string]interface{}) {
		switch k {
		case "base":
			for k2, v2 := range v.(map[string]interface{}) {
				switch k2 {
				case "id":
					schema.Definitions.Base.ID = v2.(string)
				case "type":
					schema.Definitions.Base.Type = v2.(string)
				case "required":
					var req []string
					for _, v3 := range v2.([]interface{}) {
						req = append(req, v3.(string))
					}
					schema.Definitions.Base.Required = req
				case "properties":
					for k3, v3 := range v2.(map[string]interface{}) {
						sub, err := s.client.Schemas.GetUserBaseSubSchema(k3, v3.(map[string]interface{}))
						if err != nil {
							return nil, err
						}
						schema.Definitions.Base.Properties = append(schema.Definitions.Base.Properties, *sub)
					}
				}
			}
		case "custom":
			for k2, v2 := range v.(map[string]interface{}) {
				switch k2 {
				case "id":
					schema.Definitions.Custom.ID = v2.(string)
				case "type":
					schema.Definitions.Custom.Type = v2.(string)
				case "required":
					var req []string
					for _, v3 := range v2.([]interface{}) {
						req = append(req, v3.(string))
					}
					schema.Definitions.Custom.Required = req
				case "properties":
					for k3, v3 := range v2.(map[string]interface{}) {
						sub, err := s.client.Schemas.GetUserCustomSubSchema(k3, v3.(map[string]interface{}))
						if err != nil {
							return nil, err
						}
						schema.Definitions.Custom.Properties = append(schema.Definitions.Custom.Properties, *sub)
					}
				}
			}
		}
	}
	return schema, err
}

// userSubSchemaPropMap (unexported) returns the User Profile Schema Properties as a map[string]interface{}
// input is a string subschema scope "base" or "custom"
func (s *SchemasService) userSubSchemaPropMap(scope string) (map[string]interface{}, *Response, error) {
	if scope != "base" && scope != "custom" {
		return nil, nil, fmt.Errorf("[ERROR] SubSchema Properties Map scope input supports values \"base\" or \"custom\"")
	}
	obj, resp, err := s.client.Schemas.GetRawUserSchema()
	if err != nil {
		return nil, resp, err
	}
	for k, v := range obj["definitions"].(map[string]interface{}) {
		if k == scope {
			for k2, v2 := range v.(map[string]interface{}) {
				if k2 == "properties" {
					return v2.(map[string]interface{}), nil, nil
				}
			}
		}
	}
	return nil, nil, nil
}

// GetUserSubSchemaPropMap returns the User Profile SubSchema as a map[string]interface{}
// inputs are a string subschema scope "base" or "custom" & the index key for the User Profile SubSchema
func (s *SchemasService) GetUserSubSchemaPropMap(scope string, index string) (map[string]interface{}, *Response, error) {
	prop, resp, err := s.client.Schemas.userSubSchemaPropMap(scope)
	if err != nil {
		return nil, resp, err
	}
	if v, ok := prop[index]; ok {
		return v.(map[string]interface{}), resp, err
	} else {
		return nil, resp, fmt.Errorf("[ERROR] GetUserSubSchemaPropMap subschema %v not found in Okta", index)
	}

	return nil, resp, err
}

// GetUserSubSchemaIndex returns an array of User Profile SubSchema index keys
// input is a string subschema scope "base" or "custom"
func (s *SchemasService) GetUserSubSchemaIndex(scope string) ([]string, *Response, error) {
	var index []string
	prop, resp, err := s.client.Schemas.userSubSchemaPropMap(scope)
	if err != nil {
		return nil, resp, err
	}
	for key := range prop {
		index = append(index, key)
	}
	return index, resp, err
}

// GetUserBaseSubSchema returns the User Base Profile SubSchema as a BaseSubSchema struct
// inputs are a string index key for the SubSchema & a map[string]interface{} for the
// User Profile SubSchema, such as from GetUserSubSchemaPropMap()
func (s *SchemasService) GetUserBaseSubSchema(index string, obj map[string]interface{}) (*BaseSubSchema, error) {
	subSchema := new(BaseSubSchema)
	subSchema.Index = index
	if v, ok := obj["title"]; ok {
		subSchema.Title = v.(string)
	} else {
		// if we cant find a title field, we'll assume this obj map is not correct
		return nil, fmt.Errorf("[ERROR] GetUserBaseSubSchema interface map parsing error")
	}
	if v, ok := obj["type"]; ok {
		subSchema.Type = v.(string)
	}
	if v, ok := obj["format"]; ok {
		subSchema.Format = v.(string)
	}
	if v, ok := obj["required"]; ok {
		subSchema.Required = v.(bool)
	}
	if v, ok := obj["mutability"]; ok {
		subSchema.Mutability = v.(string)
	}
	if v, ok := obj["scope"]; ok {
		subSchema.Scope = v.(string)
	}
	if v, ok := obj["minLength"]; ok {
		subSchema.MinLength = int(v.(float64))
	}
	if v, ok := obj["maxLength"]; ok {
		subSchema.MaxLength = int(v.(float64))
	}
	if v, ok := obj["permissions"]; ok {
		perms := make([]Permissions, len(v.([]interface{})))
		for k2, v2 := range v.([]interface{}) {
			for k3, v3 := range v2.(map[string]interface{}) {
				switch k3 {
				case "principal":
					perms[k2].Principal = v3.(string)
				case "action":
					perms[k2].Action = v3.(string)
				}
			}
		}
		subSchema.Permissions = perms
	}
	if v, ok := obj["master"]; ok {
		for k2, v2 := range v.(map[string]interface{}) {
			switch k2 {
			case "type":
				subSchema.Master = &Master{Type: v2.(string)}
			}
		}
	}
	return subSchema, nil
}

// GetUserCustomSubSchema returns the User Custom Profile SubSchema as a CustomSubSchema struct
// inputs are a string index key for the SubSchema & a map[string]interface{} for the
// User Profile SubSchema, such as from GetUserSubSchemaPropMap()
func (s *SchemasService) GetUserCustomSubSchema(index string, obj map[string]interface{}) (*CustomSubSchema, error) {
	subSchema := new(CustomSubSchema)
	subSchema.Index = index
	if v, ok := obj["title"]; ok {
		subSchema.Title = v.(string)
	} else {
		// if we cant find a title field, we'll assume this obj map is not correct
		return nil, fmt.Errorf("[ERROR] GetUserCustomSubSchema interface map parsing error")
	}
	if v, ok := obj["type"]; ok {
		subSchema.Type = v.(string)
	}
	if v, ok := obj["description"]; ok {
		subSchema.Description = v.(string)
	}
	if v, ok := obj["format"]; ok {
		subSchema.Format = v.(string)
	}
	if v, ok := obj["required"]; ok {
		subSchema.Required = v.(bool)
	}
	if v, ok := obj["mutability"]; ok {
		subSchema.Mutability = v.(string)
	}
	if v, ok := obj["scope"]; ok {
		subSchema.Scope = v.(string)
	}
	if v, ok := obj["minLength"]; ok {
		subSchema.MinLength = int(v.(float64))
	}
	if v, ok := obj["maxLength"]; ok {
		subSchema.MaxLength = int(v.(float64))
	}
	if v, ok := obj["items"]; ok {
		for k2, v2 := range v.(map[string]interface{}) {
			switch k2 {
			case "type":
				subSchema.Items.Type = v2.(string)
			}
		}
	}
	if v, ok := obj["union"]; ok {
		subSchema.Union = v.(string)
	}
	if v, ok := obj["enum"]; ok {
		// assuming here all enum values are strings, I hope i'm right
		enum := make([]string, 0)
		for _, v2 := range v.([]interface{}) {
			enum = append(enum, v2.(string))
		}
		subSchema.Enum = enum
	}
	if v, ok := obj["oneOf"]; ok {
		oneof := make([]OneOf, len(v.([]interface{})))
		for k2, v2 := range v.([]interface{}) {
			for k3, v3 := range v2.(map[string]interface{}) {
				switch k3 {
				case "const":
					oneof[k2].Const = v3.(string)
				case "title":
					oneof[k2].Title = v3.(string)
				}
			}
		}
		subSchema.OneOf = oneof
	}
	if v, ok := obj["permissions"]; ok {
		perms := make([]Permissions, len(v.([]interface{})))
		for k2, v2 := range v.([]interface{}) {
			for k3, v3 := range v2.(map[string]interface{}) {
				switch k3 {
				case "principal":
					perms[k2].Principal = v3.(string)
				case "action":
					perms[k2].Action = v3.(string)
				}
			}
		}
		subSchema.Permissions = perms
	}
	if v, ok := obj["master"]; ok {
		for k2, v2 := range v.(map[string]interface{}) {
			switch k2 {
			case "type":
				subSchema.Master = &Master{Type: v2.(string)}
			}
		}
	}
	return subSchema, nil
}

// UpdateUserCustomSubSchema Adds or Updates a Custom SubSchema
// input is a CustomSubSchema struct
func (s *SchemasService) UpdateUserCustomSubSchema(update CustomSubSchema) (*Schema, *Response, error) {
	index := update.Index
	subschema, err := json.Marshal(update)
	if err != nil {
		return nil, nil, err
	}
	raw := fmt.Sprintf(`{ "definitions": { "custom": { "id": "#custom", "type": "object", "properties": { "%s": %s }, "required": [] } } }`, index, string(subschema))
	// remove the escaped double quotes during NewRequest Marshal serialization
	ser := json.RawMessage(raw)
	u := fmt.Sprintf("meta/schemas/user/default")
	req, err := s.client.NewRequest("POST", u, ser)
	if err != nil {
		return nil, nil, err
	}
	var obj map[string]interface{}
	resp, err := s.client.Do(req, &obj)
	if err != nil {
		return nil, resp, err
	}
	schema, err := s.client.Schemas.userSchema(obj)
	return schema, resp, err
}

// DeleteUserCustomSubSchema deletes a Custom SubSchema
// input is a string of the custom subschema index key
func (s *SchemasService) DeleteUserCustomSubSchema(index string) (*Schema, *Response, error) {
	raw := fmt.Sprintf(`{ "definitions": { "custom": { "id": "#custom", "type": "object", "properties": { "%s": null }, "required": [] } } }`, index)
	// remove the escaped double quotes during NewRequest Marshal serialization
	ser := json.RawMessage(raw)
	u := fmt.Sprintf("meta/schemas/user/default")
	req, err := s.client.NewRequest("POST", u, ser)
	if err != nil {
		return nil, nil, err
	}
	var obj map[string]interface{}
	resp, err := s.client.Do(req, &obj)
	if err != nil {
		return nil, resp, err
	}
	schema, err := s.client.Schemas.userSchema(obj)
	return schema, resp, err
}

// UpdateUserBaseSubSchema Updates a Base SubSchema
// can only update subschema permissions & the nullability of the firstName and lastName subschemas
// input is a BaseSubSchema struct
func (s *SchemasService) UpdateUserBaseSubSchema(update BaseSubSchema) (*Schema, *Response, error) {
	index := update.Index
	subschema, err := json.Marshal(update)
	if err != nil {
		return nil, nil, err
	}
	raw := fmt.Sprintf(`{ "definitions": { "base": { "id": "#base", "type": "object", "properties": { "%s": %s }, "required": [] } } }`, index, string(subschema))
	// remove the escaped double quotes during NewRequest Marshal serialization
	ser := json.RawMessage(raw)
	u := fmt.Sprintf("meta/schemas/user/default")
	req, err := s.client.NewRequest("POST", u, ser)
	if err != nil {
		return nil, nil, err
	}
	var obj map[string]interface{}
	resp, err := s.client.Do(req, &obj)
	if err != nil {
		return nil, resp, err
	}
	schema, err := s.client.Schemas.userSchema(obj)
	return schema, resp, err
}

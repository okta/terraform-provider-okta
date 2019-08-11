package sdk

type (
	UserSchema struct {
		Schema      string                 `json:"$schema,omitempty"`
		Created     string                 `json:"created,omitempty"`
		Definitions *UserSchemaDefinitions `json:"definitions,omitempty"`
		ID          string                 `json:"id,omitempty"`
		LastUpdated string                 `json:"lastUpdated,omitempty"`
		Name        string                 `json:"name,omitempty"`
		Properties  *UserSchemaProperties  `json:"properties,omitempty"`
		Title       string                 `json:"title,omitempty"`
		Type        string                 `json:"type,omitempty"`
	}

	UserSchemaPermission struct {
		Action    string `json:"action,omitempty"`
		Principal string `json:"principal,omitempty"`
	}

	UserSchemaPropertyProfile struct {
		AllOf []*UserSchemaRef `json:"allOf,omitempty"`
	}

	UserSchemaDefinitions struct {
		Base   *BaseUserSchema   `json:"base,omitempty"`
		Custom *CustomUserSchema `json:"custom,omitempty"`
	}

	UserSubSchema struct {
		Description string                  `json:"description,omitempty"`
		Format      string                  `json:"format,omitempty"`
		MaxLength   int64                   `json:"maxLength,omitempty"`
		MinLength   int64                   `json:"minLength,omitempty"`
		Permissions []*UserSchemaPermission `json:"permissions,omitempty"`
		Required    bool                    `json:"required,omitempty"`
		Title       string                  `json:"title,omitempty"`
		Type        string                  `json:"type,omitempty"`
	}

	BaseSchemaProperties struct {
		Email     *UserSubSchema `json:"email,omitempty"`
		FirstName *UserSubSchema `json:"firstName,omitempty"`
		LastName  *UserSubSchema `json:"lastName,omitempty"`
		Login     *UserSubSchema `json:"login,omitempty"`
	}

	BaseUserSchema struct {
		ID         string                `json:"id,omitempty"`
		Properties *BaseSchemaProperties `json:"properties,omitempty"`
		Required   []string              `json:"required,omitempty"`
		Type       string                `json:"type,omitempty"`
	}

	CustomUserSchema struct {
		ID         string                    `json:"id,omitempty"`
		Properties map[string]*UserSubSchema `json:"properties,omitempty"`
		Required   []interface{}             `json:"required,omitempty"`
		Type       string                    `json:"type,omitempty"`
	}

	UserSchemaProperties struct {
		Profile *UserSchemaPropertyProfile `json:"profile,omitempty"`
	}

	UserSchemaRef struct {
		Ref string `json:"$ref,omitempty"`
	}
)

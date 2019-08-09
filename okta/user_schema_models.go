package okta

type UserSchema struct {
	Schema      string    `json:"$schema,omitempty"`
	Created     string    `json:"created,omitempty"`
	Definitions *UserSchemaDefinitions  `json:"definitions,omitempty"`
	ID          string    `json:"id,omitempty"`
	LastUpdated string    `json:"lastUpdated,omitempty"`
	Name        string    `json:"name,omitempty"`
	Properties  *UserSchemaProperties `json:"properties,omitempty"`
	Title       string    `json:"title,omitempty"`
	Type        string    `json:"type,omitempty"`
}

type UserSchemaPermission struct {
	Action    string `json:"action,omitempty"`
	Principal string `json:"principal,omitempty"`
}

type UserSchemaPropertyProfile struct {
	AllOf []*UserSchemaRef `json:"allOf,omitempty"`
}

type UserSchemaDefinitions struct {
	Base   *BaseUserSchema `json:"base,omitempty"`
	Custom *CustomUserSchema `json:"custom,omitempty"`
}

type UserSubSchema struct {
	Description string     `json:"description,omitempty"`
	Format      string     `json:"format,omitempty"`
	MaxLength   int64      `json:"maxLength,omitempty"`
	MinLength   int64      `json:"minLength,omitempty"`
	Permissions []*UserSchemaPermission `json:"permissions,omitempty"`
	Required    bool       `json:"required,omitempty"`
	Title       string     `json:"title,omitempty"`
	Type        string     `json:"type,omitempty"`
}

type BaseSchemaProperties struct {
	Email     *UserSubSchema `json:"email,omitempty"`
	FirstName *UserSubSchema  `json:"firstName,omitempty"`
	LastName  *UserSubSchema  `json:"lastName,omitempty"`
	Login     *UserSubSchema  `json:"login,omitempty"`
}

type BaseUserSchema struct {
	ID         string   `json:"id,omitempty"`
	Properties *BaseSchemaProperties `json:"properties,omitempty"`
	Required   []string `json:"required,omitempty"`
	Type       string   `json:"type,omitempty"`
}

type CustomUserSchema struct {
	ID         string        `json:"id,omitempty"`
	Properties map[string]*UserSubSchema      `json:"properties,omitempty"`
	Required   []interface{} `json:"required,omitempty"`
	Type       string        `json:"type,omitempty"`
}

type UserSchemaProperties struct {
	Profile *UserSchemaPropertyProfile `json:"profile,omitempty"`
}

type *UserSchemaRef struct {
	Ref string `json:"$ref,omitempty"`
}

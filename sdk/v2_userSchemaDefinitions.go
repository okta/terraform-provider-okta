package sdk

type UserSchemaDefinitions struct {
	Base   *UserSchemaBase   `json:"base,omitempty"`
	Custom *UserSchemaPublic `json:"custom,omitempty"`
}

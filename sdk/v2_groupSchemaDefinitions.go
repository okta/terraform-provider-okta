package sdk

type GroupSchemaDefinitions struct {
	Base   *GroupSchemaBase   `json:"base,omitempty"`
	Custom *GroupSchemaCustom `json:"custom,omitempty"`
}

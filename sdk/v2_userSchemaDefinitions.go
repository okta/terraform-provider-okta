// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type UserSchemaDefinitions struct {
	Base   *UserSchemaBase   `json:"base,omitempty"`
	Custom *UserSchemaPublic `json:"custom,omitempty"`
}

// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type OrgOktaCommunicationSettingResource resource

type OrgOktaCommunicationSetting struct {
	Links            interface{} `json:"_links,omitempty"`
	OptOutEmailUsers *bool       `json:"optOutEmailUsers,omitempty"`
}

package sdk

type OrgOktaCommunicationSettingResource resource

type OrgOktaCommunicationSetting struct {
	Links            interface{} `json:"_links,omitempty"`
	OptOutEmailUsers *bool       `json:"optOutEmailUsers,omitempty"`
}

// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

import (
	"time"
)

type OrgOktaSupportSettingsObjResource resource

type OrgOktaSupportSettingsObj struct {
	Links      interface{} `json:"_links,omitempty"`
	Expiration *time.Time  `json:"expiration,omitempty"`
	Support    string      `json:"support,omitempty"`
}

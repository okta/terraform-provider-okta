// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type ApplicationVisibility struct {
	AppLinks          map[string]bool            `json:"appLinks,omitempty"`
	AutoLaunch        *bool                      `json:"autoLaunch,omitempty"`
	AutoSubmitToolbar *bool                      `json:"autoSubmitToolbar,omitempty"`
	Hide              *ApplicationVisibilityHide `json:"hide,omitempty"`
}

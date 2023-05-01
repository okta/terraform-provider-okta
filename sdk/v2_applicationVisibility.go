package sdk

type ApplicationVisibility struct {
	AppLinks          map[string]bool            `json:"appLinks,omitempty"`
	AutoLaunch        *bool                      `json:"autoLaunch,omitempty"`
	AutoSubmitToolbar *bool                      `json:"autoSubmitToolbar,omitempty"`
	Hide              *ApplicationVisibilityHide `json:"hide,omitempty"`
}

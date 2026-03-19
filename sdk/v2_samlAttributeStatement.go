// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type SamlAttributeStatement struct {
	FilterType  string   `json:"filterType,omitempty"`
	FilterValue string   `json:"filterValue,omitempty"`
	Name        string   `json:"name,omitempty"`
	Namespace   string   `json:"namespace,omitempty"`
	Type        string   `json:"type,omitempty"`
	Values      []string `json:"values,omitempty"`
}

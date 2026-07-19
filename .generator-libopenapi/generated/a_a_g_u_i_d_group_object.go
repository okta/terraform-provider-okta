// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// AAGUIDGroupObject represents the AAGUIDGroupObject schema
type AAGUIDGroupObject struct {
	// A list of YubiKey hardware FIDO2 AAGUIDs. The available [AAGUIDs](https://support.yubico.com/hc/en-us/articles/360016648959-YubiKey-Hardware-FIDO2-AAGUIDs) are provided by the FIDO Alliance Metadat...
	Aaguids []string `json:"aaguids,omitempty"`
	// A name to identify the group of YubiKey hardware FIDO2 AAGUIDs
	Name string `json:"name,omitempty"`
}

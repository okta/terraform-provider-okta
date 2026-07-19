// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SecurePasswordStoreApplicationSettingsApplication represents the SecurePasswordStoreApplicationSettingsApplication schema
type SecurePasswordStoreApplicationSettingsApplication struct {
	// Name of the optional value in the sign-in form
	OptionalField2Value string `json:"optionalField2Value,omitempty"`
	// Name of the optional parameter in the sign-in form
	OptionalField3 string `json:"optionalField3,omitempty"`
	// CSS selector for the **Password** field in the sign-in form
	PasswordField string `json:"passwordField"`
	// CSS selector for the **Username** field in the sign-in form
	UsernameField string `json:"usernameField"`
	// Name of the optional parameter in the sign-in form
	OptionalField2 string `json:"optionalField2,omitempty"`
	// Name of the optional value in the sign-in form
	OptionalField3Value string `json:"optionalField3Value,omitempty"`
	// The URL of the sign-in page for this app
	Url string `json:"url"`
	// Name of the optional parameter in the sign-in form
	OptionalField1 string `json:"optionalField1,omitempty"`
	// Name of the optional value in the sign-in form
	OptionalField1Value string `json:"optionalField1Value,omitempty"`
}

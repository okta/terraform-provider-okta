// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// PasswordPolicyPasswordSettingsComplexity represents the PasswordPolicyPasswordSettingsComplexity schema
// Complexity settings
type PasswordPolicyPasswordSettingsComplexity struct {
	// Indicates if a password must contain at least one symbol (For example: !@#$%^&*): `0` indicates no, `1` indicates yes
	MinSymbol int `json:"minSymbol,omitempty"`
	// <x-lifecycle-container><x-lifecycle class="oie"></x-lifecycle></x-lifecycle-container>Use an [Expression Language](https://developer.okta.com/docs/reference/okta-expression-language-in-identity-eng...
	OelStatement string `json:"oelStatement,omitempty"`
	Dictionary interface{} `json:"dictionary,omitempty"`
	// Indicates if a password must contain at least one number: `0` indicates no, `1` indicates yes
	MinNumber int `json:"minNumber,omitempty"`
	// Indicates if a password must contain at least one upper case letter: `0` indicates no, `1` indicates yes
	MinUpperCase int `json:"minUpperCase,omitempty"`
	// The User profile attributes whose values must be excluded from the password: currently only supports `firstName` and `lastName`
	ExcludeAttributes []string `json:"excludeAttributes,omitempty"`
	// Indicates if the Username must be excluded from the password
	ExcludeUsername bool `json:"excludeUsername,omitempty"`
	// <x-lifecycle-container><x-lifecycle class="oie"></x-lifecycle></x-lifecycle-container>Specifies the maximum number of consecutive repeating characters that can be used in a password
	MaxConsecutiveCharacters int `json:"maxConsecutiveCharacters,omitempty"`
	// Minimum password length
	MinLength int `json:"minLength,omitempty"`
	// Indicates if a password must contain at least one lower case letter: `0` indicates no, `1` indicates yes
	MinLowerCase int `json:"minLowerCase,omitempty"`
}

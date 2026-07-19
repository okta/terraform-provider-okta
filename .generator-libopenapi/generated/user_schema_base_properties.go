// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// UserSchemaBaseProperties represents the UserSchemaBaseProperties schema
type UserSchemaBaseProperties struct {
	// State or region component of the user's address (`region`)
	State interface{} `json:"state,omitempty"`
	// Name of a cost center assigned to the user
	CostCenter interface{} `json:"costCenter,omitempty"`
	// Country name component of the user's address (`country`.) This property uses [ISO 3166-1 alpha 2 "short" code format](https://tools.ietf.org/html/draft-ietf-scim-core-schema-22#ref-ISO3166).
	CountryCode interface{} `json:"countryCode,omitempty"`
	// Name of the user's division
	Division interface{} `json:"division,omitempty"`
	// Primary email address of the user. This property is formatted according to [RFC 5322 Section 3.2.3](https://datatracker.ietf.org/doc/html/rfc5322#section-3.2.3).
	Email interface{} `json:"email,omitempty"`
	// Mobile phone number of the user
	MobilePhone interface{} `json:"mobilePhone,omitempty"`
	// Honorific prefix(es) of the user or title in most Western languages
	HonorificPrefix interface{} `json:"honorificPrefix,omitempty"`
	// Honorific suffix(es) of the user
	HonorificSuffix interface{} `json:"honorificSuffix,omitempty"`
	// Middle name(s) of the user
	MiddleName interface{} `json:"middleName,omitempty"`
	// Unique identifier for the user (`userName`)  The login property is validated according to its pattern attribute, which is a string. By default, the attribute is null. When the attribute is null, th...
	Login interface{} `json:"login,omitempty"`
	// Given name of the user (`givenName`)
	FirstName interface{} `json:"firstName,omitempty"`
	// Name of the user's organization
	Organization interface{} `json:"organization,omitempty"`
	// User's time zone. This property is formatted according to the [IANA Time Zone database format](https://tools.ietf.org/html/rfc6557).
	Timezone interface{} `json:"timezone,omitempty"`
	// Family name of the user (`familyName`)
	LastName interface{} `json:"lastName,omitempty"`
	// URL of the user's online profile (for example, a web page.) This property is formatted according to the [Relative Uniform Resource Locators specification](https://tools.ietf.org/html/draft-ietf-sci...
	ProfileUrl interface{} `json:"profileUrl,omitempty"`
	// Secondary email address of the user typically used for account recovery. This property is formatted according to [RFC 5322 Section 3.2.3](https://datatracker.ietf.org/doc/html/rfc5322#section-3.2.3).
	SecondEmail interface{} `json:"secondEmail,omitempty"`
	// Organization or company assigned unique identifier for the user
	EmployeeNumber interface{} `json:"employeeNumber,omitempty"`
	// The `displayName` of the user's manager
	Manager interface{} `json:"manager,omitempty"`
	// Mailing address component of the user's address
	PostalAddress interface{} `json:"postalAddress,omitempty"`
	// User's title, such as "Vice President"
	Title interface{} `json:"title,omitempty"`
	// Used to describe the organization to the user relationship such as "Employee" or "Contractor".  **Note:** The `userType` field is an arbitrary string value and isn't related to the newer [User Type...
	UserType interface{} `json:"userType,omitempty"`
	// The `id` of the user's manager
	ManagerId interface{} `json:"managerId,omitempty"`
	// Casual way to address the user in real life
	NickName interface{} `json:"nickName,omitempty"`
	// Full street address component of the user's address
	StreetAddress interface{} `json:"streetAddress,omitempty"`
	// User's default location for purposes of localizing items such as currency, date time format, numerical representations, and so on.  A locale value is a concatenation of the ISO 639-1 two-letter lan...
	Locale interface{} `json:"locale,omitempty"`
	// Name of the user's department
	Department interface{} `json:"department,omitempty"`
	// Name of the user, suitable for display to end users
	DisplayName interface{} `json:"displayName,omitempty"`
	// ZIP code or postal code component of the user's address (`postalCode`)
	ZipCode interface{} `json:"zipCode,omitempty"`
	// City or locality component of the user's address (`locality`)
	City interface{} `json:"city,omitempty"`
	// User's preferred written or spoken languages. This property is formatted according to [RFC 7231 Section 5.3.5](https://tools.ietf.org/html/rfc7231#section-5.3.5).
	PreferredLanguage interface{} `json:"preferredLanguage,omitempty"`
	// Primary phone number of the user, such as home number
	PrimaryPhone interface{} `json:"primaryPhone,omitempty"`
}

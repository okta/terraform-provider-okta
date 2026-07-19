// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// UserProfile represents the UserProfile schema
// Specifies the default and custom profile properties for a user.  The default user profile is based on the [System for Cross-domain Identity Management: Core Schema](https://datatracker.ietf.org/doc...
type UserProfile struct {
	// Name of the user suitable for display to end users
	DisplayName string `json:"displayName,omitempty"`
	// Honorific prefix(es) of the user, or title in most Western languages
	HonorificPrefix string `json:"honorificPrefix,omitempty"`
	// The family name of the user (`familyName`)
	LastName string `json:"lastName,omitempty"`
	// The mobile phone number of the user
	MobilePhone string `json:"mobilePhone,omitempty"`
	// The state or region component of the user's address (`region`)
	State string `json:"state,omitempty"`
	// The full street address component of the user's address
	StreetAddress string `json:"streetAddress,omitempty"`
	// The ZIP code or postal code component of the user's address (`postalCode`)
	ZipCode string `json:"zipCode,omitempty"`
	// Name of the cost center assigned to a user
	CostCenter string `json:"costCenter,omitempty"`
	// The primary email address of the user. For validation, see [RFC 5322 Section 3.2.3](https://datatracker.ietf.org/doc/html/rfc5322#section-3.2.3).
	Email string `json:"email,omitempty"`
	// The primary phone number of the user such as a home number
	PrimaryPhone string `json:"primaryPhone,omitempty"`
	// The secondary email address of the user typically used for account recovery. For validation, see [RFC 5322 Section 3.2.3](https://datatracker.ietf.org/doc/html/rfc5322#section-3.2.3).
	SecondEmail string `json:"secondEmail,omitempty"`
	// The city or locality of the user's address (`locality`)
	City string `json:"city,omitempty"`
	// The `id` of the user's manager
	ManagerId string `json:"managerId,omitempty"`
	// The casual way to address the user in real life
	NickName string `json:"nickName,omitempty"`
	// Mailing address component of the user's address
	PostalAddress string `json:"postalAddress,omitempty"`
	// The user's time zone
	Timezone string `json:"timezone,omitempty"`
	// The property used to describe the organization-to-user relationship, such as employee or contractor  > **Note:** The `userType` property is a standard string attribute and should be treated as a de...
	UserType string `json:"userType,omitempty"`
	// The organization or company assigned unique identifier for the user
	EmployeeNumber string `json:"employeeNumber,omitempty"`
	// Honorific suffix(es) of the user
	HonorificSuffix string `json:"honorificSuffix,omitempty"`
	// The user's default location for purposes of localizing items such as currency, date time format, numerical representations, and so on. A locale value is a concatenation of the ISO 639-1 two-letter ...
	Locale string `json:"locale,omitempty"`
	// The unique identifier for the user (`username`). For validation, see [Login pattern validation](https://developer.okta.com/docs/reference/api/schemas/#login-pattern-validation).  Every user within ...
	Login string `json:"login,omitempty"`
	// The `displayName` of the user's manager
	Manager string `json:"manager,omitempty"`
	// Name of the the user's organization
	Organization string `json:"organization,omitempty"`
	// The user's preferred written or spoken language. For validation, see [RFC 7231 Section 5.3.5](https://datatracker.ietf.org/doc/html/rfc7231#section-5.3.5).
	PreferredLanguage string `json:"preferredLanguage,omitempty"`
	// Name of the user's department
	Department string `json:"department,omitempty"`
	// The middle name of the user
	MiddleName string `json:"middleName,omitempty"`
	// The user's title, such as Vice President
	Title string `json:"title,omitempty"`
	// The country name component of the user's address (`country`). For validation, see [ISO 3166-1 alpha 2 "short" code format](https://datatracker.ietf.org/doc/html/draft-ietf-scim-core-schema-22#ref-I...
	CountryCode string `json:"countryCode,omitempty"`
	// Name of the user's division
	Division string `json:"division,omitempty"`
	// Given name of the user (`givenName`)
	FirstName string `json:"firstName,omitempty"`
	// The URL of the user's online profile. For example, a web page. See [URL](https://datatracker.ietf.org/doc/html/rfc1808).
	ProfileUrl string `json:"profileUrl,omitempty"`
}

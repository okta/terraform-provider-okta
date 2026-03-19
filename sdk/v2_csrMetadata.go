// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type CsrMetadata struct {
	Subject         *CsrMetadataSubject         `json:"subject,omitempty"`
	SubjectAltNames *CsrMetadataSubjectAltNames `json:"subjectAltNames,omitempty"`
}

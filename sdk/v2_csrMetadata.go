package sdk

type CsrMetadata struct {
	Subject         *CsrMetadataSubject         `json:"subject,omitempty"`
	SubjectAltNames *CsrMetadataSubjectAltNames `json:"subjectAltNames,omitempty"`
}

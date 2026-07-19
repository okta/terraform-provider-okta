// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// FCMConfiguration represents the FCMConfiguration schema
type FCMConfiguration struct {
	// (Optional) File name for Admin Console display
	FileName string `json:"fileName,omitempty"`
	// Project ID of FCM configuration
	ProjectId string `json:"projectId,omitempty"`
	// JSON containing the private service account key and service account details. See [Creating and managing service account keys](https://cloud.google.com/iam/docs/creating-managing-service-account-key...
	ServiceAccountJson map[string]interface{} `json:"serviceAccountJson,omitempty"`
}

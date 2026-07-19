// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// AppServiceAccount represents the AppServiceAccount schema
type AppServiceAccount struct {
	// Timestamp when the app service account was last updated
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	// The user-defined name for the app service account
	Name string `json:"name"`
	// The app service account password. Required for apps that don't have provisioning enabled or don't support password synchronization.
	Password string `json:"password,omitempty"`
	Status interface{} `json:"status,omitempty"`
	// The username that serves as the direct link to your managed app account. Ensure that this value precisely matches the identifier of the target app account.
	Username string `json:"username"`
	// Timestamp when the app service account was created
	Created *time.Time `json:"created,omitempty"`
	// The description of the app service account
	Description string `json:"description,omitempty"`
	// The UUID of the app service account
	ID string `json:"id,omitempty"`
	// A list of IDs of the Okta groups who own the app service account
	OwnerGroupIds []string `json:"ownerGroupIds,omitempty"`
	// A list of IDs of the Okta users who own the app service account
	OwnerUserIds []string `json:"ownerUserIds,omitempty"`
	StatusDetail interface{} `json:"statusDetail,omitempty"`
	// The key name of the app in the Okta Integration Network (OIN)
	ContainerGlobalName string `json:"containerGlobalName,omitempty"`
	// The app instance label
	ContainerInstanceName string `json:"containerInstanceName,omitempty"`
	// The [ORN](/openapi/okta-management/guides/roles/#okta-resource-name-orn) of the relevant resource.  Use the specific app ORN format (`orn:{partition}:idp:{yourOrgId}:apps:{appType}:{appId}`) to ide...
	ContainerOrn string `json:"containerOrn"`
}

// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ImportScheduleSettings represents the ImportScheduleSettings schema
type ImportScheduleSettings struct {
	// The import schedule in UNIX cron format
	Expression string `json:"expression"`
	// The import schedule time zone in Internet Assigned Numbers Authority (IANA) time zone name format
	Timezone string `json:"timezone,omitempty"`
}

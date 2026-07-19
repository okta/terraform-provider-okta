// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// AutoUpdateSchedule represents the AutoUpdateSchedule schema
// The schedule of auto-update configured by the admin
type AutoUpdateSchedule struct {
	// Delay in days
	Delay int `json:"delay,omitempty"`
	// Duration in minutes
	Duration int `json:"duration,omitempty"`
	// Timestamp when the update finished (only for a successful or failed update, not for a cancelled update). Null is returned if the job hasn't finished once yet.
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	// Timezone of where the scheduled job takes place
	Timezone string `json:"timezone,omitempty"`
	// The schedule of the update in cron format. The cron settings are limited to only the day of the month or the nth-day-of-the-week configurations. For example, `0 8 ? * 6#3` indicates every third Sat...
	Cron string `json:"cron,omitempty"`
}

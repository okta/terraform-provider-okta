// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// LogBotProtection represents the LogBotProtection schema
// <x-lifecycle-container><x-lifecycle class="ea"></x-lifecycle></x-lifecycle-container>The result of the bot protection detection associated with the event
type LogBotProtection struct {
	// The bot detected level associated with the bot protection configuration target
	Level string `json:"level,omitempty"`
}

// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// LogTransaction represents the LogTransaction schema
// A `transaction` object comprises contextual information associated with its respective event. This information is useful for understanding sequences of correlated events. For example, a `transactio...
type LogTransaction struct {
	// Details for this transaction.
	Detail map[string]interface{} `json:"detail,omitempty"`
	// Unique identifier for this transaction.
	ID string `json:"id,omitempty"`
	// Describes the kind of transaction. `WEB` indicates a web request. `JOB` indicates an asynchronous task.
	Type string `json:"type,omitempty"`
}

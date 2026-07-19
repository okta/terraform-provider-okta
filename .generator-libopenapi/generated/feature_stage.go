// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// FeatureStage represents the FeatureStage schema
// Current release cycle stage of a feature  If a feature's stage value is `EA`, the state is `null` and not returned. If the value is `BETA`, the state is `OPEN` or `CLOSED` depending on whether the ...
type FeatureStage struct {
	Value interface{} `json:"value,omitempty"`
	State interface{} `json:"state,omitempty"`
}

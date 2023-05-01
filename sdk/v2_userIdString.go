package sdk

type UserIdString struct {
	Links  interface{} `json:"_links,omitempty"`
	UserId string      `json:"userId,omitempty"`
}

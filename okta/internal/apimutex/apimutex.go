package apimutex

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

const (
	APPS_KEY      = "apps"
	APPID_KEY     = "app-id"
	CAS_KEY       = "cas-id"
	CLIENTS_KEY   = "clients"
	DEVICES_KEY   = "devices"
	EVENTS_KEY    = "events"
	GROUPS_KEY    = "groups"
	GROUPID_KEY   = "group-id"
	LOGS_KEY      = "logs"
	USERS_KEY     = "users"
	USERID_KEY    = "user-id"
	USERIDGET_KEY = "user-id-get"
	OTHER_KEY     = "other"
)

// APIMutex synchronizes keeping account of current known rate limit values
// from Okta management endpoints. Specifically apps, users, and other, see:
// https://developer.okta.com/docs/reference/rl-global-mgmt/ The Okta Terraform
// Provider can not account for all other kinds of clients utilization of API
// limits but it can account for its own usage and attempt to preemptively
// react appropriately.
type APIMutex struct {
	lock     sync.Mutex
	status   map[string]*APIStatus
	capacity int
}

// APIStatus is used to hold rate limit information from Okta's API, see:
// https://developer.okta.com/docs/reference/rl-best-practices/
type APIStatus struct {
	limit     int
	remaining int
	reset     int64 // UTC epoch time in seconds
	class     string
}

// NewAPIMutex returns a new api mutex object that represents untilized
// capacity under the specified capacity percentage.
func NewAPIMutex(capacity int) (*APIMutex, error) {
	if capacity < 1 || capacity > 100 {
		return nil, fmt.Errorf("expecting capacity as whole number > 0 and <= 100, was %d", capacity)
	}
	status := map[string]*APIStatus{
		APPS_KEY:      {class: APPS_KEY},
		APPID_KEY:     {class: APPID_KEY},
		CAS_KEY:       {class: CAS_KEY},
		CLIENTS_KEY:   {class: CLIENTS_KEY},
		DEVICES_KEY:   {class: DEVICES_KEY},
		LOGS_KEY:      {class: LOGS_KEY},
		EVENTS_KEY:    {class: EVENTS_KEY},
		GROUPS_KEY:    {class: GROUPS_KEY},
		GROUPID_KEY:   {class: GROUPID_KEY},
		OTHER_KEY:     {class: OTHER_KEY},
		USERS_KEY:     {class: USERS_KEY},
		USERID_KEY:    {class: USERID_KEY},
		USERIDGET_KEY: {class: USERIDGET_KEY},
	}
	return &APIMutex{
		capacity: capacity,
		status:   status,
	}, nil
}

// HasCapacity approximates if there is capacity below the api mutex's maximum
// capacity threshold.
func (m *APIMutex) HasCapacity(method, endPoint string) bool {
	status := m.get(method, endPoint)

	// if the status hasn't been updated recently assume there is capacity
	if status.reset+60 < time.Now().Unix() {
		return true
	}

	// calculate utilization
	utilization := 100.0 * (float32(status.limit-status.remaining) / float32(status.limit))

	return utilization <= float32(m.capacity)
}

// Update updates the known status for the given API endpoint. It is synchronous
// and intelligently accounts for new values regardless of parallelism.
func (m *APIMutex) Update(method, endPoint string, limit, remaining int, reset int64) {
	m.lock.Lock()
	defer m.lock.Unlock()

	key := m.normalizeKey(method, endPoint)
	status := m.status[key]

	if reset > status.reset {
		// reset value greater than current reset implies we are in a new Okta API
		// one minute window. set/reset values.
		status.reset = reset
		status.remaining = remaining
		status.limit = limit
		return
	}

	if reset <= (status.reset - 60) {
		// these values are from the previous one minute window, ignore
		return
	}

	if remaining < status.remaining {
		status.remaining = remaining
	}
}

// Status return the APIStatus for the given class of endpoint.
func (m *APIMutex) Status(method, endPoint string) *APIStatus {
	return m.get(method, endPoint)
}

var reAppId = regexp.MustCompile("/api/v1/apps/[^/]+$")
var reGroupId = regexp.MustCompile("/api/v1/groups/[^/]+$")

func (m *APIMutex) normalizeKey(method, endPoint string) string {
	var result string
	switch {
	case reAppId.MatchString(endPoint):
		result = APPID_KEY
	case strings.Contains(endPoint, "/api/v1/apps"):
		result = APPS_KEY
	case endPoint == "/api/v1/certificateAuthorities":
		result = CAS_KEY
	case endPoint == "/oauth2/v1/clients":
		result = CLIENTS_KEY
	case endPoint == "/api/v1/devices":
		result = DEVICES_KEY
	case endPoint == "/api/v1/events":
		result = EVENTS_KEY
	case reGroupId.MatchString(endPoint):
		result = GROUPID_KEY
	case strings.Contains(endPoint, "/api/v1/groups"):
		result = GROUPS_KEY
	case endPoint == "/api/v1/logs":
		result = LOGS_KEY
	case endPoint == "/api/v1/users":
		result = USERS_KEY
	case method == http.MethodGet && strings.HasPrefix(endPoint, "/api/v1/users/"):
		result = USERIDGET_KEY
	case strings.HasPrefix(endPoint, "/api/v1/users/"):
		result = USERID_KEY
	default:
		result = "other"
	}
	return result
}

// Reset returns the current reset value of the api status object.
func (s *APIStatus) Reset() int64 {
	return s.reset
}

// Limit returns the current limit value of the api status object.
func (s *APIStatus) Limit() int {
	return s.limit
}

// Remaining returns the current remaining value of the api status object.
func (s *APIStatus) Remaining() int {
	return s.remaining
}

// Class returns the api endpoint class for this status.
func (s *APIStatus) Class() string {
	return s.class
}

func (m *APIMutex) get(method, endPoint string) *APIStatus {
	key := m.normalizeKey(method, endPoint)
	return m.status[key]
}

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
// from Okta management endpoints. See:
// https://developer.okta.com/docs/reference/rl-global-mgmt/
//
// The Okta Terraform Provider can not account for other clients consumption of
// API limits but it can account for its own usage and attempt to preemptively
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

var (
	reAppId   = regexp.MustCompile("/api/v1/apps/[^/]+$")
	reGroupId = regexp.MustCompile("/api/v1/groups/[^/]+$")
	reUserId  = regexp.MustCompile("/api/v1/users/[^/]+$")
)

func (m *APIMutex) normalizeKey(method, endPoint string) string {
	// Okta internal: see rate-limit-mappings-CLASSIC-DEFAULT.txt file in core
	// repo.  It corresponds to:
	// https://developer.okta.com/docs/reference/rl-best-practices/
	//
	// TODO: API rate limits can be overwritten by the org admin, we should come
	//       up with a way to accommodate for that.  Perhaps caching an APIStatus
	//       struct on the SHA of http method + URI path.

	getPutDelete := (http.MethodGet == method) ||
		(http.MethodPut == method) ||
		(http.MethodDelete == method)
	postPutDelete := (http.MethodPost == method) ||
		(http.MethodPut == method) ||
		(http.MethodDelete == method)
	var result string

	switch {
	//  1. [GET|PUT|DELETE] /api/v1/apps/${id}
	case reAppId.MatchString(endPoint) && getPutDelete:
		result = APPID_KEY

	//  2. starts with /api/v1/apps
	case strings.HasPrefix(endPoint, "/api/v1/apps"):
		result = APPS_KEY

	//  3. [GET|PUT|DELETE] /api/v1/groups/${id}
	case reGroupId.MatchString(endPoint) && getPutDelete:
		result = GROUPID_KEY

	//  4. starts with /api/v1/groups
	case strings.HasPrefix(endPoint, "/api/v1/groups"):
		result = GROUPS_KEY

	//  5. [POST|PUT|DELETE] /api/v1/users/${id}
	case reUserId.MatchString(endPoint) && postPutDelete:
		result = USERID_KEY

	//  6. [GET] /api/v1/users/${idOrLogin}
	case reUserId.MatchString(endPoint) && method == http.MethodGet:
		result = USERIDGET_KEY

	//  7. starts with /api/v1/users
	case strings.HasPrefix(endPoint, "/api/v1/users"):
		result = USERS_KEY

	//  8. GET /api/v1/logs
	case endPoint == "/api/v1/logs" && method == http.MethodGet:
		result = LOGS_KEY

	//  9. GET /api/v1/events
	case endPoint == "/api/v1/events" && method == http.MethodGet:
		result = EVENTS_KEY

	// 10. GET /oauth2/v1/clients
	case endPoint == "/oauth2/v1/clients" && method == http.MethodGet:
		result = CLIENTS_KEY

	// 11. GET /api/v1/certificateAuthorities
	case endPoint == "/api/v1/certificateAuthorities" && method == http.MethodGet:
		result = CAS_KEY

	// 12. GET /api/v1/devices
	case endPoint == "/api/v1/devices" && method == http.MethodGet:
		result = DEVICES_KEY

	// 13. GET /api/v1
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

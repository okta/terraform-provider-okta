package apimutex

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"
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
	capacity int
	status   map[string]*APIStatus
	buckets  map[string]string
}

// APIStatus is used to hold rate limit information from Okta's API, see:
// https://developer.okta.com/docs/reference/rl-best-practices/
type APIStatus struct {
	limit     int
	remaining int
	reset     int64 // UTC epoch time in seconds
}

// NewAPIMutex returns a new api mutex object that represents untilized
// capacity under the specified capacity percentage.
func NewAPIMutex(capacity int) (*APIMutex, error) {
	rootStatus := &APIStatus{}
	mutex := &APIMutex{
		capacity: capacity,
		status: map[string]*APIStatus{
			"/": rootStatus,
		},
		buckets: map[string]string{},
	}
	mutex.initRateLimitLookup()

	return mutex, nil
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

	status := m.get(method, endPoint)
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

// Status Returns the APIStatus for the given method + endpoint combination.
func (m *APIMutex) Status(method, endPoint string) *APIStatus {
	return m.get(method, endPoint)
}

// Class Returns the api endpoint class.
func (m *APIMutex) Class(method, endPoint string) string {
	return m.normalizedKey(method, endPoint)
}

func (m *APIMutex) normalizedKey(method, endPoint string) string {
	return fmt.Sprintf("%s %s", method, endPoint)
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

var (
	reOktaID = regexp.MustCompile(`\A[\w]{20}\z`)
)

func (m *APIMutex) get(method, endPoint string) *APIStatus {
	// the important point here is the replace all is performing this
	// transformation for the bucket lookup /api/v1/users/abcdefghij0123456789
	// to /api/v1/users/ID
	path := reOktaID.ReplaceAllString(endPoint, "ID")
	key := m.normalizedKey(method, path)
	bucket, ok := m.buckets[key]
	if !ok {
		return m.status["/"]
	}
	return m.status[bucket]
}

func (m *APIMutex) initRateLimitLookup() {
	for _, line := range rateLimitLines {
		vals := strings.Split(line, " ")
		path := vals[0]
		method := vals[1]
		bucket := vals[2]

		key := m.normalizedKey(method, path)
		m.buckets[key] = bucket

		if _, ok := m.status[bucket]; !ok {
			m.status[bucket] = &APIStatus{}
		}
	}
}

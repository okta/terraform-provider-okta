package idaas

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/okta/terraform-provider-okta/okta/api"
	"github.com/okta/terraform-provider-okta/okta/config"
	"github.com/okta/terraform-provider-okta/okta/internal/apprefresh"
	"github.com/stretchr/testify/require"
)

type batchTransport struct {
	mu          sync.Mutex
	calls       map[string]int
	failPage    bool
	unavailable int
}

func (b *batchTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	key := r.Method + " " + r.URL.Path
	b.calls[key]++
	body := ""
	code := 200
	header := make(http.Header)
	switch {
	case r.Method == http.MethodPut:
		body = `{"id":"g1","priority":1,"profile":{}}`
	case strings.HasSuffix(r.URL.Path, "/groups"):
		if r.URL.Query().Get("after") == "" {
			body = `[{"id":"g1","priority":1,"profile":{"value":"kept"}}]`
			header.Set("Link", `<https://example.okta.com/api/v1/apps/app/groups?after=g1>; rel="next"`)
		} else {
			body = assignmentPage()
			if b.failPage {
				code = 500
				body = `{"errorCode":"E0000009","errorSummary":"second page failed"}`
			}
		}
	case strings.HasSuffix(r.URL.Path, "/groups/g1"):
		body = `{"id":"g1","priority":1,"profile":{}}`
	case strings.HasSuffix(r.URL.Path, "/groups/missing"):
		code = 404
		body = `{"errorCode":"E0000007","errorSummary":"not found"}`
	case strings.HasSuffix(r.URL.Path, "/group-push/mappings"):
		if r.URL.Query().Get("after") == "" {
			body = `[{"id":"p1","sourceGroupId":"g1","targetGroupId":"t1","status":"ACTIVE"}]`
			header.Set("Link", `<https://example.okta.com/api/v1/apps/app/group-push/mappings?after=p1>; rel="next"`)
		} else {
			body = mappingPage()
			if b.failPage {
				code = 500
				body = `{"errorCode":"E0000009","errorSummary":"second page failed"}`
			}
		}
	case strings.HasSuffix(r.URL.Path, "/group-push/mappings/missing"):
		body = `{"id":"missing","sourceGroupId":"g3","targetGroupId":"t3","status":"ACTIVE"}`
	default:
		return nil, fmt.Errorf("unexpected request: %s", key)
	}
	if b.unavailable != 0 && strings.HasSuffix(r.URL.Path, "/groups") {
		code = b.unavailable
		body = `{"errorCode":"E0000007","errorSummary":"collection unavailable"}`
	}
	header.Set("Content-Type", "application/json")
	return &http.Response{StatusCode: code, Status: fmt.Sprintf("%d %s", code, http.StatusText(code)), Header: header, Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
}

func batchConfig(t *testing.T, b *batchTransport) *config.Config {
	t.Helper()
	// Explicit fake credentials; no environment-derived provider configuration.
	c := &config.Config{AppReadBatching: true}
	client, err := api.NewOktaIDaaSAPIClient(&api.OktaAPIConfig{OrgName: "example", Domain: "okta.com", ApiToken: "test-only", Logger: hclog.NewNullLogger(), AppReads: &c.AppReads})
	require.NoError(t, err)
	c.SetIdaasAPIClient(client)
	client.OktaSDKClientV2().GetConfig().HttpClient.Transport = apprefresh.InvalidateWrites(b, &c.AppReads)
	client.OktaSDKClientV6().GetConfig().HTTPClient.Transport = apprefresh.InvalidateWrites(b, &c.AppReads)
	return c
}

func TestBatchAppAssignmentPaginationAndWrites(t *testing.T) {
	b := &batchTransport{calls: map[string]int{}}
	c := batchConfig(t, b)
	g, _, err := readAppAssignment(t.Context(), c, "app", "g1")
	require.NoError(t, err)
	require.Equal(t, "g1", g.Id)
	g, _, err = readAppAssignment(t.Context(), c, "app", "g2")
	require.NoError(t, err)
	require.EqualValues(t, 2, *g.PriorityPtr)
	for i := 3; i <= 100; i++ {
		g, _, err = readAppAssignment(t.Context(), c, "app", fmt.Sprintf("g%d", i))
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf("g%d", i), g.Id)
	}
	require.Equal(t, 2, b.calls["GET /api/v1/apps/app/groups"])
	// A write through the other SDK must also invalidate the assignment snapshot.
	req, err := http.NewRequest(http.MethodPut, "https://example.okta.com/api/v1/apps/app/groups/g1", nil)
	require.NoError(t, err)
	res, err := c.OktaIDaaSClient.OktaSDKClientV6().GetConfig().HTTPClient.Do(req)
	require.NoError(t, err)
	res.Body.Close()
	_, _, err = readAppAssignment(t.Context(), c, "app", "g1")
	require.NoError(t, err)
	require.Equal(t, 2, b.calls["GET /api/v1/apps/app/groups"])
	require.Equal(t, 1, b.calls["GET /api/v1/apps/app/groups/g1"])
	_, resp, err := readAppAssignment(t.Context(), c, "app", "missing")
	require.Error(t, err)
	require.Equal(t, 404, resp.StatusCode)
	require.Equal(t, 1, b.calls["GET /api/v1/apps/app/groups/missing"])
}

func TestBatchPushMappingPaginationAndMissingFallback(t *testing.T) {
	b := &batchTransport{calls: map[string]int{}}
	c := batchConfig(t, b)
	p, err := readPushMapping(t.Context(), c, "app", "p1")
	require.NoError(t, err)
	require.Equal(t, "t1", p.GetTargetGroupId())
	p, err = readPushMapping(t.Context(), c, "app", "p2")
	require.NoError(t, err)
	require.Equal(t, "INACTIVE", p.GetStatus())
	for i := 3; i <= 100; i++ {
		p, err = readPushMapping(t.Context(), c, "app", fmt.Sprintf("p%d", i))
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf("p%d", i), p.GetId())
	}
	require.Equal(t, 2, b.calls["GET /api/v1/apps/app/group-push/mappings"])
	p, err = readPushMapping(t.Context(), c, "app", "missing")
	require.NoError(t, err)
	require.Equal(t, "missing", p.GetId())
	require.Equal(t, 1, b.calls["GET /api/v1/apps/app/group-push/mappings/missing"])
}

func TestBatchPartialCollectionFailsWithoutRemovingState(t *testing.T) {
	b := &batchTransport{calls: map[string]int{}, failPage: true}
	c := batchConfig(t, b)
	_, resp, err := readAppAssignment(t.Context(), c, "app", "g1")
	require.Error(t, err)
	require.Nil(t, resp)
	_, err = readPushMapping(t.Context(), c, "app", "p1")
	require.Error(t, err)
	b.failPage = false
	_, _, err = readAppAssignment(t.Context(), c, "app", "g1")
	require.NoError(t, err)
	_, err = readPushMapping(t.Context(), c, "app", "p1")
	require.NoError(t, err)
	require.Equal(t, 4, b.calls["GET /api/v1/apps/app/groups"])
	require.Equal(t, 4, b.calls["GET /api/v1/apps/app/group-push/mappings"])
}

func assignmentPage() string {
	rows := make([]string, 0, 99)
	for i := 2; i <= 100; i++ {
		rows = append(rows, fmt.Sprintf(`{"id":"g%d","priority":%d,"profile":{}}`, i, i))
	}
	return "[" + strings.Join(rows, ",") + "]"
}
func mappingPage() string {
	rows := make([]string, 0, 99)
	for i := 2; i <= 100; i++ {
		rows = append(rows, fmt.Sprintf(`{"id":"p%d","sourceGroupId":"g%d","targetGroupId":"t%d","status":"INACTIVE"}`, i, i, i))
	}
	return "[" + strings.Join(rows, ",") + "]"
}

func TestBatchDisabledUsesIndividualRead(t *testing.T) {
	b := &batchTransport{calls: map[string]int{}}
	c := batchConfig(t, b)
	c.AppReadBatching = false
	_, _, err := readAppAssignment(t.Context(), c, "app", "g1")
	require.NoError(t, err)
	_, err = readPushMapping(t.Context(), c, "app", "missing")
	require.NoError(t, err)
	require.Zero(t, b.calls["GET /api/v1/apps/app/groups"])
	require.Zero(t, b.calls["GET /api/v1/apps/app/group-push/mappings"])
}

func TestBatchUnavailableAssignmentCollectionUsesIndividualRead(t *testing.T) {
	for _, status := range []int{http.StatusForbidden, http.StatusNotFound} {
		b := &batchTransport{calls: map[string]int{}, unavailable: status}
		c := batchConfig(t, b)
		g, _, err := readAppAssignment(t.Context(), c, "app", "g1")
		require.NoError(t, err)
		require.Equal(t, "g1", g.Id)
		_, resp, err := readAppAssignment(t.Context(), c, "app", "missing")
		require.Error(t, err)
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	}
}

package idaas

import (
	"testing"

	"github.com/okta/terraform-provider-okta/sdk"
	"github.com/stretchr/testify/require"
)

func tfGroup(id string, priority int, profile string) interface{} {
	return map[string]interface{}{"id": id, "priority": priority, "profile": profile}
}

func assignmentIDs(assignments []*sdk.ApplicationGroupAssignment) []string {
	ids := make([]string, 0, len(assignments))
	for _, a := range assignments {
		ids = append(ids, a.Id)
	}
	return ids
}

func assignmentPriorities(assignments []*sdk.ApplicationGroupAssignment) []int64 {
	priorities := make([]int64, 0, len(assignments))
	for _, a := range assignments {
		priorities = append(priorities, *a.PriorityPtr)
	}
	return priorities
}

func TestPlanGroupAssignments_SendsReorderInAscendingPriority(t *testing.T) {
	old := []interface{}{
		tfGroup("a", 1, `{"role":"user"}`),
		tfGroup("b", 2, `{"role":"user"}`),
		tfGroup("c", 3, `{"role":"user"}`),
		tfGroup("d", 4, `{"role":"user"}`),
	}
	planned := []interface{}{
		tfGroup("d", 4, `{"role":"user"}`),
		tfGroup("c", 1, `{"role":"user"}`),
		tfGroup("a", 3, `{"role":"user"}`),
		tfGroup("b", 2, `{"role":"user"}`),
	}

	for i := 0; i < 50; i++ {
		toAssign, toRemove := planGroupAssignments(old, planned)
		require.Empty(t, toRemove)
		require.Equal(t, []string{"c", "b", "a", "d"}, assignmentIDs(toAssign))
		require.Equal(t, []int64{1, 2, 3, 4}, assignmentPriorities(toAssign))
	}
}

func TestPlanGroupAssignments_RemovalResendsAllInOrder(t *testing.T) {
	old := []interface{}{
		tfGroup("a", 1, `{}`),
		tfGroup("x", 2, `{}`),
		tfGroup("b", 3, `{}`),
		tfGroup("y", 4, `{}`),
	}
	planned := []interface{}{
		tfGroup("a", 1, `{}`),
		tfGroup("b", 2, `{}`),
	}

	toAssign, toRemove := planGroupAssignments(old, planned)
	require.Equal(t, []string{"x", "y"}, assignmentIDs(toRemove))
	require.Equal(t, []string{"a", "b"}, assignmentIDs(toAssign))
	require.Equal(t, []int64{1, 2}, assignmentPriorities(toAssign))
}

func TestPlanGroupAssignments_ProfileOnlyChangeSendsOnlyThatGroup(t *testing.T) {
	old := []interface{}{
		tfGroup("a", 1, `{"groups":null,"role":"coadmin"}`),
		tfGroup("b", 2, `{"groups":null,"role":"user"}`),
		tfGroup("c", 3, `{"groups":null,"role":"user"}`),
	}
	planned := []interface{}{
		tfGroup("a", 1, `{"role":"coadmin","groups":null}`),
		tfGroup("b", 2, `{"groups":null,"role":"coadmin"}`),
		tfGroup("c", 3, `{"groups":null,"role":"user"}`),
	}

	toAssign, toRemove := planGroupAssignments(old, planned)
	require.Empty(t, toRemove)
	require.Equal(t, []string{"b"}, assignmentIDs(toAssign))
}

func TestPlanGroupAssignments_AppendSendsOnlyNewGroup(t *testing.T) {
	old := []interface{}{
		tfGroup("a", 1, `{}`),
		tfGroup("b", 2, `{}`),
	}
	planned := []interface{}{
		tfGroup("a", 1, `{}`),
		tfGroup("b", 2, `{}`),
		tfGroup("c", 3, `{}`),
	}

	toAssign, toRemove := planGroupAssignments(old, planned)
	require.Empty(t, toRemove)
	require.Equal(t, []string{"c"}, assignmentIDs(toAssign))
}

func TestPlanGroupAssignments_InsertResendsAllInOrder(t *testing.T) {
	old := []interface{}{
		tfGroup("a", 1, `{}`),
		tfGroup("b", 2, `{}`),
	}
	planned := []interface{}{
		tfGroup("a", 1, `{}`),
		tfGroup("c", 2, `{}`),
		tfGroup("b", 3, `{}`),
	}

	toAssign, toRemove := planGroupAssignments(old, planned)
	require.Empty(t, toRemove)
	require.Equal(t, []string{"a", "c", "b"}, assignmentIDs(toAssign))
	require.Equal(t, []int64{1, 2, 3}, assignmentPriorities(toAssign))
}

func TestPlanGroupAssignments_NoChangeSendsNothing(t *testing.T) {
	groups := []interface{}{
		tfGroup("a", 1, `{"role":"user"}`),
		tfGroup("b", 2, `{"role":"user"}`),
	}

	toAssign, toRemove := planGroupAssignments(groups, groups)
	require.Empty(t, toRemove)
	require.Empty(t, toAssign)
}

func TestPlanGroupAssignments_CreateFromEmptyState(t *testing.T) {
	planned := []interface{}{
		tfGroup("b", 2, `{}`),
		tfGroup("a", 1, `{}`),
	}

	toAssign, toRemove := planGroupAssignments(nil, planned)
	require.Empty(t, toRemove)
	require.Equal(t, []string{"a", "b"}, assignmentIDs(toAssign))
}

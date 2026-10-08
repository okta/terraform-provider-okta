package governance

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func idModels(ids ...string) []IdModel {
	out := make([]IdModel, 0, len(ids))
	for _, id := range ids {
		out = append(out, IdModel{Id: types.StringValue(id)})
	}
	return out
}

func TestKeepOrderIfSameIds(t *testing.T) {
	tests := []struct {
		name      string
		api       []IdModel
		reference []IdModel
		want      []IdModel
	}{
		{
			name:      "same order",
			api:       idModels("a", "b", "c"),
			reference: idModels("a", "b", "c"),
			want:      idModels("a", "b", "c"),
		},
		{
			name:      "same ids reordered keeps reference order",
			api:       idModels("a", "b", "c"),
			reference: idModels("c", "a", "b"),
			want:      idModels("c", "a", "b"),
		},
		{
			name:      "id added by api returns api",
			api:       idModels("a", "b", "c"),
			reference: idModels("a", "b"),
			want:      idModels("a", "b", "c"),
		},
		{
			name:      "id missing from api returns api",
			api:       idModels("a"),
			reference: idModels("b", "a"),
			want:      idModels("a"),
		},
		{
			name:      "id replaced returns api",
			api:       idModels("a", "x"),
			reference: idModels("b", "a"),
			want:      idModels("a", "x"),
		},
		{
			name:      "duplicates counted",
			api:       idModels("a", "a", "b"),
			reference: idModels("b", "a", "a"),
			want:      idModels("b", "a", "a"),
		},
		{
			name:      "duplicate count differs returns api",
			api:       idModels("a", "b", "b"),
			reference: idModels("b", "a", "a"),
			want:      idModels("a", "b", "b"),
		},
		{
			name:      "empty reference (import) returns api",
			api:       idModels("b", "a"),
			reference: nil,
			want:      idModels("b", "a"),
		},
		{
			name:      "empty api returns api",
			api:       nil,
			reference: idModels("a"),
			want:      nil,
		},
		{
			name:      "both empty",
			api:       nil,
			reference: nil,
			want:      nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := keepOrderIfSameIds(tt.api, tt.reference)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("keepOrderIfSameIds() = %v, want %v", got, tt.want)
			}
		})
	}
}

package impose

import (
	"reflect"
	"testing"
)

func TestCalculatePageOrder(t *testing.T) {
	tests := []struct {
		name       string
		totalPages int
		want       []int
	}{
		{
			name:       "4 pages",
			totalPages: 4,
			want:       []int{3, 0, 1, 2},
		},
		{
			name:       "8 pages",
			totalPages: 8,
			want:       []int{7, 0, 1, 6, 5, 2, 3, 4},
		},
		{
			name:       "6 pages (padded to 8)",
			totalPages: 6,
			want:       []int{7, 0, 1, 6, 5, 2, 3, 4},
		},
		{
			name:       "12 pages",
			totalPages: 12,
			want:       []int{11, 0, 1, 10, 9, 2, 3, 8, 7, 4, 5, 6},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculatePageOrder(tt.totalPages)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CalculatePageOrder(%d) = %v, want %v",
					tt.totalPages, got, tt.want)
			}
		})
	}
}

package drift

import "testing"

func TestGrade_Boundaries(t *testing.T) {
	cases := []struct {
		score float64
		want  string
	}{
		{100.0, "A"},
		{90.0, "A"},
		{89.9, "B"},
		{75.0, "B"},
		{74.9, "C"},
		{50.0, "C"},
		{49.9, "D"},
		{25.0, "D"},
		{24.9, "F"},
		{0.0, "F"},
	}
	for _, tc := range cases {
		got := grade(tc.score)
		if got != tc.want {
			t.Errorf("grade(%.1f) = %s, want %s", tc.score, got, tc.want)
		}
	}
}

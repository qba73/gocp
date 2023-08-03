package geometry_test

import (
	"math"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/qba73/gocp/geometry"
)

var floatComparer = cmp.Comparer(func(x, y float64) bool {
	delta := math.Abs(x - y)
	mean := math.Abs(x+y) / 2.0
	return delta/mean < 0.00001
})

func TestDistance_CalculatesDistanceBetweenTwoPoints(t *testing.T) {
	t.Parallel()

	want := 1.41422

	p := geometry.Point{1, 1}
	q := geometry.Point{2, 2}
	got := geometry.Distance(p, q)

	if !cmp.Equal(got, want, floatComparer) {
		t.Error(cmp.Diff(got, want))
	}
}

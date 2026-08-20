package chatserver

import (
	"math"
	"testing"
)

func TestHaversineKmSamePoint(t *testing.T) {
	t.Parallel()

	got := haversineKm(52.52, 13.405, 52.52, 13.405)
	if got != 0 {
		t.Fatalf("expected 0 for identical points, got %v", got)
	}
}

// One degree of latitude change (same longitude) is a great-circle arc, so
// its length is independently derivable from earthRadiusKm without going
// through haversineKm's own trig — a real check on the formula, not a
// tautology.
func TestHaversineKmOneDegreeLatitude(t *testing.T) {
	t.Parallel()

	got := haversineKm(0, 0, 1, 0)
	want := earthRadiusKm * (math.Pi / 180)
	if diff := math.Abs(got - want); diff > 0.01 {
		t.Fatalf("expected ~%.4fkm for 1 degree of latitude, got %.4fkm (diff %.4f)", want, got, diff)
	}
}

// Antipodal points are exactly half the Earth's circumference apart —
// another independently-derivable reference value.
func TestHaversineKmAntipodal(t *testing.T) {
	t.Parallel()

	got := haversineKm(10, 20, -10, -160)
	want := math.Pi * earthRadiusKm
	if diff := math.Abs(got - want); diff > 0.01 {
		t.Fatalf("expected ~%.4fkm for antipodal points, got %.4fkm (diff %.4f)", want, got, diff)
	}
}

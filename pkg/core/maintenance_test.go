package core

import "testing"

func TestMatchMaintenance(t *testing.T) {
	windows := []MaintenanceWindow{{
		Project:  "cupress",
		Reason:   "deploy",
		StartsAt: "2026-10-21 10:00:00",
		EndsAt:   "2026-10-21 12:00:00",
	}}
	if _, ok := MatchMaintenance("2026-10-21 10:30:00", windows); !ok {
		t.Fatal("a timestamp inside the window should match")
	}
	if _, ok := MatchMaintenance("2026-10-21 09:59:00", windows); ok {
		t.Fatal("before the window")
	}
	if _, ok := MatchMaintenance("2026-10-21 12:00:00", windows); ok {
		t.Fatal("the end instant is outside the window")
	}
}

func TestValidateMaintenanceHours(t *testing.T) {
	if err := ValidateMaintenanceHours(4); err != nil {
		t.Fatal(err)
	}
	if err := ValidateMaintenanceHours(0); err == nil {
		t.Fatal("zero hours must be rejected")
	}
	if err := ValidateMaintenanceHours(48); err == nil {
		t.Fatal("48 hours must be rejected")
	}
}

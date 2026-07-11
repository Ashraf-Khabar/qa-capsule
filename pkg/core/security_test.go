package core

import "testing"

func TestEmailMatchesAllowedDomain_EmptyPolicy(t *testing.T) {
	if !EmailMatchesAllowedDomain("any@where.com", "") {
		t.Error("empty policy should allow all")
	}
}

func TestEmailMatchesAllowedDomain_Normalized(t *testing.T) {
	if !EmailMatchesAllowedDomain("User@Company.COM", "company.com") {
		t.Error("expected case-insensitive suffix match")
	}
}

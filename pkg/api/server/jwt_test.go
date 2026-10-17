package server

import (
	"testing"

	"github.com/SRE-Capsule/sre-capsule-community/pkg/core"
)

func TestIsDevelopmentEnv(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	if isDevelopmentEnv() {
		t.Fatal("production should not be development")
	}
	t.Setenv("APP_ENV", "development")
	if !isDevelopmentEnv() {
		t.Fatal("expected development")
	}
}

func TestInitJWT_allowsDevelopmentFallback(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	t.Setenv("SRECAPSULE_JWT_SECRET", "")
	InitJWT(&core.Config{})
	if len(jwtKey) == 0 {
		t.Fatal("expected dev fallback JWT key")
	}
}

package server

import (
	"log"
	"os"
	"strings"

	"github.com/SRE-Capsule/sre-capsule-community/pkg/core"
)

// jwtKey is the secret used to sign JWT tokens (set via InitJWT at startup).
var jwtKey []byte

const defaultJWTSecret = "sre-super-secret-jwt-key"

func isDevelopmentEnv() bool {
	v := strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
	return v == "development" || v == "dev" || v == "local"
}

// InitJWT configures the signing key from config or SRECAPSULE_JWT_SECRET.
// Production requires an explicit secret; APP_ENV=development allows the dev fallback.
func InitJWT(config *core.Config) {
	secret := ""
	if config != nil {
		secret = strings.TrimSpace(config.Security.JWTSecret)
	}
	if secret == "" {
		secret = strings.TrimSpace(os.Getenv("SRECAPSULE_JWT_SECRET"))
	}
	if secret == "" {
		if isDevelopmentEnv() {
			secret = defaultJWTSecret
			log.Println("[SECURITY] APP_ENV=development: using built-in JWT secret. Set security.jwt_secret or SRECAPSULE_JWT_SECRET for production-like setups.")
		} else {
			log.Fatal("[FATAL] JWT secret required: set security.jwt_secret in config.yaml or SRECAPSULE_JWT_SECRET (or APP_ENV=development for local dev fallback only)")
		}
	}
	jwtKey = []byte(secret)
}

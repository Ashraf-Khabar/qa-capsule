package main

import (
	"log/slog"
	"os"
	"strings"

	"github.com/SRE-Capsule/sre-capsule-community/pkg/api/server"
	"github.com/SRE-Capsule/sre-capsule-community/pkg/core"
)

func ensureLocalDevJWTEnv(config core.Config) {
	if strings.TrimSpace(os.Getenv("APP_ENV")) != "" {
		return
	}
	if strings.TrimSpace(os.Getenv("SRECAPSULE_JWT_SECRET")) != "" {
		return
	}
	if strings.TrimSpace(config.Security.JWTSecret) != "" {
		return
	}
	_ = os.Setenv("APP_ENV", "development")
	slog.Warn("no JWT secret configured; APP_ENV=development enabled for this local run — set security.jwt_secret or SRECAPSULE_JWT_SECRET in production")
}

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	core.EnsureProjectRoot()
	config := core.LoadConfig()
	ensureLocalDevJWTEnv(config)
	if err := core.InitRemediationEngine(config.Plugins.Directory); err != nil {
		slog.Error("remediation engine init failed", "error", err)
		os.Exit(1)
	}
	core.InitDB("sre-capsule.db")
	core.InitSuperApp()
	server.Start(config)
}

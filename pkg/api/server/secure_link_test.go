package server

import (
	"strings"
	"testing"
	"time"
)

func TestCommentLink_RoundTripAndRejectsTampering(t *testing.T) {
	jwtKey = []byte("dev-compose-change-me")
	now := time.Unix(1_700_000_000, 0)
	token, err := mintCommentLink(2, 10, now)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(token, "2") && strings.Contains(token, "10") && len(token) < 24 {
		t.Fatalf("token looks like a raw id: %s", token)
	}
	if strings.ContainsAny(token, "/+=") {
		t.Fatalf("token is not a path segment: %s", token)
	}
	incidentID, commentID, err := openCommentLink(token, now.Add(24*time.Hour))
	if err != nil || incidentID != 2 || commentID != 10 {
		t.Fatalf("open: id=%d comment=%d err=%v", incidentID, commentID, err)
	}
	raw := []byte(token)
	raw[8] ^= 0x2a
	if _, _, err := openCommentLink(string(raw), now); err == nil {
		t.Fatal("tampered token was accepted")
	}
	if _, _, err := openCommentLink(token, now.Add(linkTTL+time.Hour)); err == nil {
		t.Fatal("expired token was accepted")
	}
}

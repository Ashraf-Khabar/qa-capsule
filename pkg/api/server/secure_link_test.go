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

func TestLinkPath_UsesTwoShardsPlusTheHash(t *testing.T) {
	jwtKey = []byte("dev-compose-change-me")
	now := time.Unix(1_700_000_000, 0)
	token, err := mintCommentLink(2, 10, now)
	if err != nil {
		t.Fatal(err)
	}
	path := linkPath(token)
	a, b, got, ok := splitDynamicLink(path)
	if !ok || got != token {
		t.Fatalf("path %s did not split", path)
	}
	wantA, wantB := linkShards(token)
	if a != wantA || b != wantB {
		t.Fatalf("shards = %s/%s, want %s/%s", a, b, wantA, wantB)
	}
	if strings.Contains(path, "/l/") || strings.Count(path, "/") != 3 {
		t.Fatalf("path is not /XXXX/YYYY/hash: %s", path)
	}
	other := linkPath(token + "x")
	if other == path {
		t.Fatal("different tokens produced the same path")
	}
}

func TestViewLink_HidesThePageName(t *testing.T) {
	jwtKey = []byte("dev-compose-change-me")
	now := time.Unix(1_700_000_000, 0)
	token, err := sealLink(secureLink{K: "view", V: "plugins", T: "", E: now.Add(time.Hour).Unix()})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(token, "plugins") {
		t.Fatalf("page name leaked into the token: %s", token)
	}
	link, err := openLink(token, now)
	if err != nil || link.K != "view" || link.V != "plugins" {
		t.Fatalf("open view: %#v %v", link, err)
	}
}

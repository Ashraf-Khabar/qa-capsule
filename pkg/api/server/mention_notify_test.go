package server

import (
	"strings"
	"testing"
)

func TestMentionEmail_NamesTheAuthorAndEscapes(t *testing.T) {
	subject, body := mentionEmail("Achraf KHABAR", "https://capsule.example/l/abc")
	if subject != "Achraf KHABAR mentioned you" {
		t.Fatalf("subject: %s", subject)
	}
	if !strings.Contains(body, "Achraf KHABAR") || !strings.Contains(body, "https://capsule.example/l/abc") {
		t.Fatalf("body missing link: %s", body)
	}
	_, hostile := mentionEmail("Ann<script>", `https://capsule.example/l/a"b`)
	if strings.Contains(hostile, "<script>") || strings.Contains(hostile, `a"b`) {
		t.Fatalf("email was not escaped: %s", hostile)
	}
}

func TestMentionsSelf(t *testing.T) {
	blocks := []CommentBlock{{Mention: "Ada@Example.com"}}
	if !mentionsSelf(blocks, "ada@example.com") {
		t.Fatal("author should not be able to mention themselves")
	}
	if mentionsSelf(blocks, "sam@example.com") {
		t.Fatal("another person is a valid mention")
	}
	if got := mentionRejection("ada@example.com", blocks); got != "You cannot mention yourself" {
		t.Fatalf("rejection: %q", got)
	}
}

func TestMentionsPlatformAdmin(t *testing.T) {
	users := map[string]directoryUser{
		"admin":           {Username: "admin", Fullname: "Platform Administrator", Role: "admin"},
		"ada@example.com": {Username: "ada@example.com", Fullname: "Ada", Role: "observer"},
	}
	if !mentionsPlatformAdmin([]CommentBlock{{Mention: "admin"}}, users) {
		t.Fatal("admin mention should be rejected")
	}
	if mentionsPlatformAdmin([]CommentBlock{{Mention: "ada@example.com"}}, users) {
		t.Fatal("a teammate mention should be allowed")
	}
}

func TestMentionedUsers_SkipsPeopleAlreadyTagged(t *testing.T) {
	blocks := []CommentBlock{{Mention: "Ada@Example.com"}, {Mention: "sam@example.com"}, {T: " hi"}}
	got := mentionedUsers(blocks)
	if len(got) != 2 {
		t.Fatalf("users: %#v", got)
	}
	previous := mentionedInHTML(`<span class="mention" data-user="ada@example.com">@Ada</span>`)
	fresh := 0
	for _, user := range got {
		if !previous[strings.ToLower(user)] {
			fresh++
		}
	}
	if fresh != 1 {
		t.Fatalf("expected only the new mention, got %d", fresh)
	}
}

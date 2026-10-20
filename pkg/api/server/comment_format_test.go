package server

import (
	"strings"
	"testing"
)

func TestRenderCommentBlocks_FormatsAndEscapes(t *testing.T) {
	html, ok := RenderCommentBlocks([]CommentBlock{
		{T: "Hello ", B: true},
		{Mention: "ada"},
		{T: "<script>", Code: true, Color: "#ff00aa", Size: "18px"},
		{Pre: "x := 1"},
	})
	if !ok {
		t.Fatal("expected content")
	}
	if strings.Contains(html, "<script>") {
		t.Fatalf("script tag was not escaped: %s", html)
	}
	if !strings.Contains(html, "<strong>Hello </strong>") {
		t.Fatalf("missing bold: %s", html)
	}
	if !strings.Contains(html, `data-user="ada"`) || !strings.Contains(html, ">@ada<") {
		t.Fatalf("missing mention: %s", html)
	}
	if !strings.Contains(html, "color:#ff00aa") || !strings.Contains(html, "font-size:18px") {
		t.Fatalf("missing style: %s", html)
	}
	if !strings.Contains(html, "<pre><code>x := 1</code></pre>") {
		t.Fatalf("missing code block: %s", html)
	}
}

func TestRenderCommentBlocks_Lists(t *testing.T) {
	html, ok := RenderCommentBlocks([]CommentBlock{
		{T: "one", Li: "ol", Item: true},
		{T: "<two>", B: true, Li: "ol", Item: true, Color: "#3b82f6"},
		{T: "dot", Li: "ul", Item: true},
		{T: " after"},
	})
	if !ok {
		t.Fatal("expected content")
	}
	want := `<ol><li>one</li><li><span style="color:#3b82f6;"><strong>&lt;two&gt;</strong></span></li></ol><ul><li>dot</li></ul> after`
	if html != want {
		t.Fatalf("lists:\n got %s\nwant %s", html, want)
	}
	bad, ok := RenderCommentBlocks([]CommentBlock{{T: "x", Li: "script", Item: true}})
	if !ok || strings.Contains(bad, "<script") || strings.Contains(bad, "<li>") {
		t.Fatalf("list kind must stay ul or ol: %s", bad)
	}
}

func TestRenderCommentBlocks_Empty(t *testing.T) {
	if _, ok := RenderCommentBlocks(nil); ok {
		t.Fatal("nil blocks should be empty")
	}
	if _, ok := RenderCommentBlocks([]CommentBlock{{Br: true}, {T: "   "}}); ok {
		t.Fatal("whitespace should be empty")
	}
	if _, ok := RenderCommentBlocks([]CommentBlock{{Mention: "bad name"}}); ok {
		t.Fatal("invalid mention should be dropped")
	}
}

func TestRenderCommentBlocks_MentionShowsName(t *testing.T) {
	blocks := []CommentBlock{{Mention: "ada@example.com"}, {T: " please look"}}
	applyMentionLabels(blocks, map[string]string{"ada@example.com": "Ada Lovelace"})
	html, ok := RenderCommentBlocks(blocks)
	if !ok {
		t.Fatal("expected content")
	}
	if !strings.Contains(html, `data-user="ada@example.com"`) || !strings.Contains(html, ">@Ada Lovelace<") {
		t.Fatalf("mention should show the name: %s", html)
	}
	if strings.Contains(html, ">@ada@example.com<") {
		t.Fatalf("mention should not show the email: %s", html)
	}
	if strings.Contains(html, "<script>") {
		t.Fatal("label was not escaped")
	}
	blocks[0].MentionLabel = `<script>`
	html, _ = RenderCommentBlocks(blocks)
	if strings.Contains(html, "<script>") {
		t.Fatalf("script label was not escaped: %s", html)
	}
}

func TestRewriteMentionHTML_UsesCurrentName(t *testing.T) {
	in := `<span class="mention" data-user="ada@example.com">@ada@example.com</span> hi`
	out := rewriteMentionHTML(in, map[string]string{"ada@example.com": "Ada Lovelace"})
	if !strings.Contains(out, ">@Ada Lovelace<") || !strings.Contains(out, `data-user="ada@example.com"`) {
		t.Fatalf("rewrite failed: %s", out)
	}
	if !mentionedInHTML(out)["ada@example.com"] {
		t.Fatal("username should stay readable for later edits")
	}
}

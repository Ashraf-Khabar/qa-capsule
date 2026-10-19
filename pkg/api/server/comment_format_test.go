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
	if !strings.Contains(html, `data-user="ada"`) {
		t.Fatalf("missing mention: %s", html)
	}
	if !strings.Contains(html, "color:#ff00aa") || !strings.Contains(html, "font-size:18px") {
		t.Fatalf("missing style: %s", html)
	}
	if !strings.Contains(html, "<pre><code>x := 1</code></pre>") {
		t.Fatalf("missing code block: %s", html)
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

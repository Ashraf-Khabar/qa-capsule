package server

import (
	"html"
	"regexp"
	"strings"
	"unicode/utf8"
)

// CommentBlock is one formatted span produced by the comment editor.
// The server renders these blocks itself so stored HTML cannot carry scripts.
type CommentBlock struct {
	T       string `json:"t,omitempty"`
	B       bool   `json:"b,omitempty"`
	I       bool   `json:"i,omitempty"`
	U       bool   `json:"u,omitempty"`
	Code    bool   `json:"code,omitempty"`
	Pre     string `json:"pre,omitempty"`
	Mention string `json:"mention,omitempty"`
	Color   string `json:"color,omitempty"`
	Size    string `json:"size,omitempty"`
	Br      bool   `json:"br,omitempty"`
}

var (
	commentColorRe = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)
	commentSizeRe  = regexp.MustCompile(`^(12|14|16|18|20|24|28)px$`)
	commentUserRe  = regexp.MustCompile(`^[a-zA-Z0-9._@+\-]{1,80}$`)
)

// RenderCommentBlocks turns editor blocks into a safe HTML fragment.
// The boolean is false when the comment has no visible text.
func RenderCommentBlocks(blocks []CommentBlock) (string, bool) {
	if len(blocks) == 0 {
		return "", false
	}
	if len(blocks) > 400 {
		blocks = blocks[:400]
	}
	var b strings.Builder
	textRunes := 0
	for _, block := range blocks {
		if textRunes > 8000 {
			break
		}
		if block.Br {
			b.WriteString("<br>")
			continue
		}
		if strings.TrimSpace(block.Pre) != "" {
			pre := truncateRunes(block.Pre, 4000)
			textRunes += utf8.RuneCountInString(pre)
			b.WriteString("<pre><code>")
			b.WriteString(html.EscapeString(pre))
			b.WriteString("</code></pre>")
			continue
		}
		if strings.TrimSpace(block.Mention) != "" {
			user := strings.TrimSpace(block.Mention)
			if !commentUserRe.MatchString(user) {
				continue
			}
			esc := html.EscapeString(user)
			b.WriteString(`<span class="mention" data-user="`)
			b.WriteString(esc)
			b.WriteString(`">@`)
			b.WriteString(esc)
			b.WriteString(`</span>`)
			textRunes += utf8.RuneCountInString(user) + 1
			continue
		}
		text := block.T
		if strings.TrimSpace(text) == "" {
			if text != "" && textRunes > 0 {
				b.WriteString(html.EscapeString(text))
			}
			continue
		}
		if utf8.RuneCountInString(text) > 2000 {
			text = truncateRunes(text, 2000)
		}
		textRunes += utf8.RuneCountInString(text)
		inner := html.EscapeString(text)
		if block.Code {
			inner = "<code>" + inner + "</code>"
		}
		if block.B {
			inner = "<strong>" + inner + "</strong>"
		}
		if block.I {
			inner = "<em>" + inner + "</em>"
		}
		if block.U {
			inner = "<u>" + inner + "</u>"
		}
		style := ""
		if commentColorRe.MatchString(block.Color) {
			style += "color:" + block.Color + ";"
		}
		if commentSizeRe.MatchString(block.Size) {
			style += "font-size:" + block.Size + ";"
		}
		if style != "" {
			inner = `<span style="` + style + `">` + inner + `</span>`
		}
		b.WriteString(inner)
	}
	out := strings.TrimSpace(b.String())
	if out == "" || textRunes == 0 {
		return "", false
	}
	return out, true
}

func truncateRunes(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n])
}

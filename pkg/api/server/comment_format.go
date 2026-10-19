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
	T            string `json:"t,omitempty"`
	B            bool   `json:"b,omitempty"`
	I            bool   `json:"i,omitempty"`
	U            bool   `json:"u,omitempty"`
	Code         bool   `json:"code,omitempty"`
	Pre          string `json:"pre,omitempty"`
	Mention      string `json:"mention,omitempty"`
	MentionLabel string `json:"-"`
	Color        string `json:"color,omitempty"`
	Size         string `json:"size,omitempty"`
	Br           bool   `json:"br,omitempty"`
}

var (
	commentColorRe = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)
	commentSizeRe  = regexp.MustCompile(`^(12|14|16|18|20|24|28)px$`)
	commentUserRe  = regexp.MustCompile(`^[a-zA-Z0-9._@+\-]{1,80}$`)
	mentionSpanRe  = regexp.MustCompile(`<span class="mention" data-user="([^"]*)">@[^<]*</span>`)
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
			label := strings.TrimSpace(block.MentionLabel)
			if label == "" {
				label = user
			}
			b.WriteString(`<span class="mention" data-user="`)
			b.WriteString(html.EscapeString(user))
			b.WriteString(`">@`)
			b.WriteString(html.EscapeString(label))
			b.WriteString(`</span>`)
			textRunes += utf8.RuneCountInString(label) + 1
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

// applyMentionLabels sets the visible name for each @mention.
// names is keyed by lowercase username. The stored identity stays the username.
func applyMentionLabels(blocks []CommentBlock, names map[string]string) {
	for i := range blocks {
		user := strings.TrimSpace(blocks[i].Mention)
		if user == "" {
			continue
		}
		if name := strings.TrimSpace(names[strings.ToLower(user)]); name != "" {
			blocks[i].MentionLabel = name
		}
	}
}

// rewriteMentionHTML replaces the visible @mention text with the current display name.
func rewriteMentionHTML(body string, names map[string]string) string {
	if body == "" || len(names) == 0 {
		return body
	}
	return mentionSpanRe.ReplaceAllStringFunc(body, func(span string) string {
		parts := mentionSpanRe.FindStringSubmatch(span)
		if len(parts) < 2 {
			return span
		}
		user := html.UnescapeString(parts[1])
		label := strings.TrimSpace(names[strings.ToLower(user)])
		if label == "" {
			return span
		}
		return `<span class="mention" data-user="` + html.EscapeString(user) + `">@` + html.EscapeString(label) + `</span>`
	})
}

// mentionedInHTML returns the lowercase usernames already present in rendered comment HTML.
func mentionedInHTML(body string) map[string]bool {
	out := map[string]bool{}
	for _, parts := range mentionSpanRe.FindAllStringSubmatch(body, -1) {
		if len(parts) < 2 {
			continue
		}
		user := strings.ToLower(strings.TrimSpace(html.UnescapeString(parts[1])))
		if user != "" {
			out[user] = true
		}
	}
	return out
}

func mentionedUsers(blocks []CommentBlock) []string {
	seen := map[string]bool{}
	out := make([]string, 0)
	for _, block := range blocks {
		user := strings.TrimSpace(block.Mention)
		if user == "" || !commentUserRe.MatchString(user) {
			continue
		}
		key := strings.ToLower(user)
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, user)
	}
	return out
}

func truncateRunes(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n])
}

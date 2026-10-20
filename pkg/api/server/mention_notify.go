package server

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/SRE-Capsule/sre-capsule-community/pkg/core"
)

type directoryUser struct {
	Username string
	Fullname string
	Role     string
}

func activeUsersByName() map[string]directoryUser {
	out := map[string]directoryUser{}
	if core.DB == nil {
		return out
	}
	rows, err := core.DB.Query(`SELECT username, COALESCE(fullname, ''), COALESCE(role, '') FROM users WHERE is_active = 1`)
	if err != nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		var user directoryUser
		if err := rows.Scan(&user.Username, &user.Fullname, &user.Role); err != nil {
			continue
		}
		user.Username = strings.TrimSpace(user.Username)
		user.Fullname = strings.TrimSpace(user.Fullname)
		if user.Username == "" {
			continue
		}
		out[strings.ToLower(user.Username)] = user
	}
	return out
}

func fullnameMap(users map[string]directoryUser) map[string]string {
	names := make(map[string]string, len(users))
	for key, user := range users {
		if user.Fullname != "" {
			names[key] = user.Fullname
		}
	}
	return names
}

func mentionsPlatformAdmin(blocks []CommentBlock, users map[string]directoryUser) bool {
	for _, mentioned := range mentionedUsers(blocks) {
		user, ok := users[strings.ToLower(mentioned)]
		if ok && core.IsAdmin(user.Role) {
			return true
		}
	}
	return false
}

func mentionsSelf(blocks []CommentBlock, author string) bool {
	author = strings.ToLower(strings.TrimSpace(author))
	if author == "" {
		return false
	}
	for _, mentioned := range mentionedUsers(blocks) {
		if strings.EqualFold(mentioned, author) {
			return true
		}
	}
	return false
}

// mentionRejection blocks tags the author is not allowed to send.
func mentionRejection(author string, blocks []CommentBlock) string {
	if mentionsSelf(blocks, author) {
		return "You cannot mention yourself"
	}
	if mentionsPlatformAdmin(blocks, activeUsersByName()) {
		return "You cannot mention a platform administrator"
	}
	return ""
}

func stampMentionLabels(blocks []CommentBlock) {
	applyMentionLabels(blocks, fullnameMap(activeUsersByName()))
}

func mentionEmail(authorName, link string) (string, string) {
	name := strings.TrimSpace(authorName)
	if name == "" {
		name = "A teammate"
	}
	name = strings.NewReplacer("\r", " ", "\n", " ").Replace(name)
	subject := name + " mentioned you"
	body := fmt.Sprintf(`
<div style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif; background-color: #0d1117; color: #c9d1d9; padding: 30px; border-radius: 8px; border: 1px solid #30363d; max-width: 600px; margin: 0 auto;">
	<h2 style="color: #58a6ff; margin-top: 0;">SRE Capsule</h2>
	<p><strong>%s</strong> mentioned you in a comment.</p>
	<p style="margin: 28px 0;">
		<a href="%s" style="display: inline-block; background: #238636; color: #ffffff; text-decoration: none; padding: 10px 16px; border-radius: 6px; font-weight: 600;">Open the comment</a>
	</p>
	<p style="color: #8b949e; font-size: 13px;">Sign in to open it. The link stays on this comment.</p>
</div>`, html.EscapeString(name), html.EscapeString(link))
	return subject, body
}

// notifyNewMentions emails people who were just tagged. Editing a note does not email someone already tagged.
func notifyNewMentions(config *core.Config, r *http.Request, author string, incidentID, commentID int64, blocks []CommentBlock, previousHTML string) {
	if config == nil || incidentID <= 0 || commentID <= 0 {
		return
	}
	users := activeUsersByName()
	authorName := author
	if user, ok := users[strings.ToLower(author)]; ok && user.Fullname != "" {
		authorName = user.Fullname
	}
	previous := mentionedInHTML(previousHTML)
	recipients := make([]string, 0)
	seen := map[string]bool{}
	for _, mentioned := range mentionedUsers(blocks) {
		key := strings.ToLower(mentioned)
		if previous[key] || strings.EqualFold(mentioned, author) || seen[key] {
			continue
		}
		user, ok := users[key]
		if !ok || core.IsAdmin(user.Role) || !strings.Contains(user.Username, "@") {
			continue
		}
		seen[key] = true
		recipients = append(recipients, user.Username)
	}
	if len(recipients) == 0 {
		return
	}
	token, err := mintCommentLink(incidentID, commentID, time.Now())
	if err != nil {
		log.Printf("[MENTION] could not build a direct link: %v", err)
		return
	}
	link := requestPublicBase(r) + linkPath(secureLink{K: "comment"}, token)
	subject, body := mentionEmail(authorName, link)
	cfg := config
	for _, to := range recipients {
		to := to
		go sendEmail(cfg, to, subject, body)
	}
}

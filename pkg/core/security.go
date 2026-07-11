package core

import "strings"

// EmailMatchesAllowedDomain returns true when allowedDomain is empty or email ends with the policy suffix.
// Policy may be stored as "@company.com" or "company.com" (case-insensitive).
func EmailMatchesAllowedDomain(email, allowedDomain string) bool {
	email = strings.TrimSpace(strings.ToLower(email))
	allowedDomain = strings.TrimSpace(strings.ToLower(allowedDomain))
	if allowedDomain == "" {
		return true
	}
	if !strings.Contains(email, "@") {
		return false
	}
	if !strings.HasPrefix(allowedDomain, "@") {
		allowedDomain = "@" + allowedDomain
	}
	return strings.HasSuffix(email, allowedDomain)
}

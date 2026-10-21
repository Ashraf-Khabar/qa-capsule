package server

import (
	"net/http"
	"strings"
	"time"

	"github.com/SRE-Capsule/sre-capsule-community/pkg/core"
)

func appOrStaticHandler() http.Handler {
	static := staticWebHandler()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			static.ServeHTTP(w, r)
			return
		}
		if token, ok := namedLinkToken(r.URL.Path); ok {
			link, err := openLink(token, time.Now())
			if err != nil || !pathMatchesLink(r.URL.Path, link, token) {
				http.NotFound(w, r)
				return
			}
			serveAppIndex(w, r)
			return
		}
		a, b, token, ok := splitDynamicLink(r.URL.Path)
		if !ok {
			static.ServeHTTP(w, r)
			return
		}
		wantA, wantB := linkShards(token)
		if a != wantA || b != wantB {
			static.ServeHTTP(w, r)
			return
		}
		if _, err := openLink(token, time.Now()); err != nil {
			http.NotFound(w, r)
			return
		}
		serveAppIndex(w, r)
	})
}

func staticWebHandler() http.Handler {
	fs := http.FileServer(http.Dir("./web"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/" || path == "/index.html" || strings.HasSuffix(path, ".html") || strings.HasSuffix(path, ".js") || strings.HasSuffix(path, ".css") {
			w.Header().Set("Cache-Control", "no-store")
		}
		w.Header().Set("X-SRE-Capsule-Edition", core.EditionID())
		w.Header().Set("X-SRE-Capsule-Version", core.Version)
		fs.ServeHTTP(w, r)
	})
}

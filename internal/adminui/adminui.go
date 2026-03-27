// Package adminui embeds the compiled admin dashboard SPA and serves it
// for all /admin/* paths. The frontend is built from admin-ui/ using Vite.
package adminui

import (
	"embed"
	"io"
	"io/fs"
	"net/http"
	"strings"
)

//go:embed all:dist
var dist embed.FS

// Handler returns an http.Handler that serves the admin SPA.
// Mount it at /admin and /admin/* in the router.
//
// Static assets are served directly. All other paths fall through to
// index.html so the React router handles client-side navigation.
func Handler() http.Handler {
	sub, err := fs.Sub(dist, "dist")
	if err != nil {
		panic("adminui: " + err.Error())
	}
	fileServer := http.FileServer(http.FS(sub))

	return http.StripPrefix("/admin", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			// Root /admin/ — serve index.html directly (bypass FileServer redirect).
			serveIndex(w, r, sub)
			return
		}

		// Try to open the file in the embedded FS.
		f, err := sub.Open(path)
		if err != nil {
			// Unknown path: serve index.html for React Router client-side navigation.
			serveIndex(w, r, sub)
			return
		}

		stat, err := f.Stat()
		f.Close()
		if err != nil || stat.IsDir() {
			serveIndex(w, r, sub)
			return
		}

		fileServer.ServeHTTP(w, r)
	}))
}

// serveIndex writes index.html directly to the response, bypassing the
// http.FileServer redirect that turns /index.html back into ./.
func serveIndex(w http.ResponseWriter, r *http.Request, sub fs.FS) {
	f, err := sub.Open("index.html")
	if err != nil {
		http.Error(w, "admin not found", http.StatusNotFound)
		return
	}
	defer f.Close()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	io.Copy(w, f) //nolint:errcheck
}

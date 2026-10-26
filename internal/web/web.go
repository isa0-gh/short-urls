package web

import (
	"embed"
	"net/http"
)

//go:embed *.html *.css *.js
var files embed.FS

func RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(files))))
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/login", http.StatusFound)
	})

	for route, filename := range map[string]string{
		"/login":          "login.html",
		"/login.html":     "login.html",
		"/register":       "register.html",
		"/register.html":  "register.html",
		"/dashboard":      "dashboard.html",
		"/dashboard.html": "dashboard.html",
	} {
		mux.HandleFunc("GET "+route, page(filename))
	}
}

func page(filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contents, err := files.ReadFile(filename)
		if err != nil {
			http.Error(w, "Page not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(contents)
	}
}

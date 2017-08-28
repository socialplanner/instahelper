package handlers

import (
	"mime"
	"net/http"
	"path/filepath"
)

// MimeTypeMiddleware will set the mimetype of the content being requested
func MimeTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := mime.TypeByExtension(filepath.Ext(r.URL.String()))

		w.Header().Set("Content-Type", contentType)

		next.ServeHTTP(w, r)
	})
}

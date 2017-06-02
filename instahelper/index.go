package instahelper

import "net/http"

// Index of instahelper "/"
func Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

// Favicon .ico of instahelper "/favicon.ico"
func Favicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/favicon.ico")
}

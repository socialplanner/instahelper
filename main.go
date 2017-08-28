package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/socialplanner/instahelper/app/handlers"
)

func main() {
	// To be removed on working prototype :)
	fmt.Println("Rome wasn't built in a day.")

	r := chi.NewRouter()

	// MIDDLEWARE
	// gzip compress
	r.Use(middleware.DefaultCompress)
	// do not cache
	r.Use(middleware.NoCache)
	// recovers from panic gracefully
	r.Use(middleware.Recoverer)
	// timeout requests after 30 seconds
	r.Use(middleware.Timeout(time.Second * 30))
	// redirect "/url/" to "/url"
	r.Use(middleware.RedirectSlashes)

	// Pages
	for _, p := range handlers.Pages {
		r.Get(p.Link, p.Handler())
	}

	// Assets
	r.With(handlers.MimeTypeMiddleware).Get(
		"/assets/*",
		handlers.AssetsHandler,
	)

	http.ListenAndServe(":3000", r)
}

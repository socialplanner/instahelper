package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/skratchdot/open-golang/open"

	"github.com/sirupsen/logrus"

	"github.com/socialplanner/instahelper/app/auth"
	"github.com/socialplanner/instahelper/app/config"
	"github.com/socialplanner/instahelper/app/handlers"
	"github.com/socialplanner/instahelper/app/insta"
	"github.com/socialplanner/instahelper/app/notifications"
)

func main() {
	port := flag.Int("port", 3333, "Port to run Instahelper on")
	debug := flag.Bool("debug", false, "Run in debug mode")
	username := flag.String("user", "", "Username for instahelper")
	password := flag.String("pass", "", "Password for instahelper")
	useAuth := flag.Bool("auth", false, "Run using http basic auth")

	flag.Parse()

	// To be removed on working prototype :)
	fmt.Println("Rome wasn't built in a day.")

	err := config.Open()

	if err != nil {
		logrus.Fatal("While opening db ", err)
	}

	defer config.Close()

	c, err := config.Config()

	if err != nil {
		logrus.Fatal(err)
	}

	r := chi.NewRouter()

	// They passed in a command line argument to use basic auth
	if *useAuth || *password != "" || *username != "" {
		if *password == "" {
			pass, err := insta.Decrypt(c.Password)
			if err != nil {
				logrus.Fatal(err)
			}

			password = &pass
		}

		if *username == "" {
			*username = c.Username
		}
		logrus.Info("Using authentication")
		r.Use(auth.SimpleBasicAuth(*username, *password))
	}
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

	// ROUTES
	// Pages

	if *debug {
		for _, p := range handlers.Pages {
			r.With(middleware.Logger).Get(p.Link, p.Handler)
		}
	} else {
		for _, p := range handlers.Pages {
			r.Get(p.Link, p.Handler)
		}
	}

	// Assets
	r.With(handlers.MimeTypeMiddleware).Get(
		"/assets/*",
		handlers.AssetsHandler,
	)

	// API
	r.Route("/api", func(r chi.Router) {

		// Accounts
		r.Route("/accounts", func(r chi.Router) {
			r.Get("/", handlers.APIAccountsHandler)
			r.Post("/create", handlers.APICreateAccountHandler)
			r.Delete("/{username}", handlers.APIDeleteAccountHandler)
			r.Post("/update/{username}", handlers.APIUpdateAccountHandler)
		})

		r.Route("/notifications", func(r chi.Router) {
			r.Delete("/", handlers.APIDeleteNotificationsHandler)
		})

		r.Route("/update", func(r chi.Router) {
			r.Post("/to/{version}", handlers.APIUpdateToHandler)
			r.Post("/", handlers.APIUpdateHandler)
		})

		r.Route("/settings", func(r chi.Router) {
			r.Post("/edit", handlers.APISettingsEditHandler)
		})
	})

	// Websocket handler
	go notifications.Hub.Start()
	r.Get("/ws", notifications.WSHandler)

	// Use config if no port passed
	if *port == 3333 {
		*port = c.Port
	}

	go func() {
		time.Sleep(time.Second)
		ip, err := localIP()

		if err != nil {
			ip = "localhost"
		}

		logrus.Infof("Up and running at http://%s:%d !", ip, *port)
		err = open.Run(fmt.Sprintf("http://%s:%d", ip, *port))

		if err != nil {
			logrus.Error(err)
		}

	}()

	logrus.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), r))
}

// localIP returns the local ip of the current device
func localIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}

	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String(), nil
}

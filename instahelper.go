package main

import (
	"net/http"

	"github.com/skratchdot/open-golang/open"

	"time"

	app "github.com/socialplanner/instahelper/instahelper"
)

func main() {

	log := app.Log

	// index
	http.HandleFunc("/", app.Index)

	// favicon.ico
	http.HandleFunc("/favicon.ico", app.Favicon)

	log.Info("Running on port 8080. Visit http://localhost:8080")

	go func() {
		time.Sleep(1 * time.Second)
		err := open.Run("http://localhost:8080")

		if err != nil {
			log.Error(err)
		}
	}()
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal(err)
	}

}

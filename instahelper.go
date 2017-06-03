package main

import (
	"net/http"

	"github.com/skratchdot/open-golang/open"

	"time"

	"fmt"

	app "github.com/socialplanner/instahelper/instahelper"
)

const logLevel = "debug"
const port = "8080"

func main() {

	log := app.Log

	app.SetLoggingLevel(logLevel)

	// index
	http.HandleFunc("/", app.Index)
	log.Debug("Set index handler.")

	// favicon.ico
	http.HandleFunc("/favicon.ico", app.Favicon)
	log.Debug("Set favicon.ico handler")

	ip, err := GetIP()

	if err != nil {
		log.Info("Couldn't fetch your IP. To visit instahelper on other devices within your network visit [DEVICEIP]:8080.")
		log.Info("If this is running on a VPS, visit [DOMAINNAME]:8080")
		ip = "localhost"
	}

	url := fmt.Sprintf("http://%s:%s", ip, port)

	log.Infof("Running on port %s. Visit %s", port, url)

	// Opens the url in the default web browser
	if logLevel != "debug" {
		go func() {
			time.Sleep(1 * time.Second)
			err = open.Run(url)

			if err != nil {
				log.Error(err)
			}
		}()
	}

	err = http.ListenAndServe(":"+port, nil)

	if err != nil {
		log.Fatal(err)
	}

}

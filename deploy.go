// +build ignore

package main

import (
	"io/ioutil"
	"fmt"
	"net/http"
	"os"

	"github.com/socialplanner/instahelper/app/update"
)

var client = &http.Client{}

func main() {
	// DELETE NIGHTLY BUILD IN PREPARATION FOR NIGHTLY FROM TRAVIS
	key := os.Getenv("GITHUB_KEY")
	releases, err := update.ListReleases()

	for _, r := range releases {
		if r.Version == "nightly" {
			err := delete(r.ID)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}
}

// delete will delete the release based off of ID
func delete(releaseid int) error {
	req, _ := http.NewRequest(
		"DELETE", 
		fmt.Sprintf(
			"https://api.github.com/repos/socialplanner/instahelper/releases/%d",
			releaseid),
		nil
	)
	req.BasicAuth("jaynagpaul", GITHUB_KEY, true)

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("Expected 200 got %d. %s", resp.StatusCode, ioutil.ReadAll(resp.Body))
	}
}

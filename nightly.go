// +build ignore

package main

// Used to autogenerate nightly builds on commit.

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/socialplanner/instahelper/app/update"
)

var client = &http.Client{}

func main() {
	// DELETE NIGHTLY BUILD IN PREPARATION FOR NIGHTLY FROM TRAVIS
	key := os.Getenv("GITHUB_KEY")

	releases, err := listReleases("jaynagpaul", key)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, r := range releases {
		if r.Name == "Nightly" {
			err := delete("jaynagpaul", key, r.ID)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}
}

// delete will delete the release based off of ID
func delete(username, password string, releaseid int) error {

	fmt.Println("Deleting release with ID", releaseid)

	req, _ := http.NewRequest(
		"DELETE",
		fmt.Sprintf(
			"https://api.github.com/repos/socialplanner/instahelper/releases/%d",
			releaseid),
		nil,
	)

	req.SetBasicAuth(username, password)

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		b, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("Expected 200 got %d. %s", resp.StatusCode, string(b))
	}
	return nil
}

// Special version of listReleases which takes a username/password to avoid hitting the ratelimit on Travis' IP.
// ListReleases collects information about github releases
func listReleases(username, password string) ([]update.Release, error) {
	req, _ := http.NewRequest("GET", "https://api.github.com/repos/socialplanner/instahelper/releases", nil)

	req.SetBasicAuth(username, password)

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		time.Sleep(3 * time.Second)
		return listReleases(username, password)
	}

	b, _ := ioutil.ReadAll(resp.Body)

	var releases []update.Release

	err = json.Unmarshal(b, &releases)

	if err != nil {
		return nil, err
	}

	return releases, err
}

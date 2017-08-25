// +build ignore

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/socialplanner/instahelper/app/update"
)

var client = &http.Client{}

func main() {
	// DELETE NIGHTLY BUILD IN PREPARATION FOR NIGHTLY FROM TRAVIS
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
	key := os.Getenv("GITHUB_KEY")
	req, _ := http.NewRequest(
		"DELETE",
		fmt.Sprintf(
			"https://api.github.com/repos/socialplanner/instahelper/releases/%d",
			releaseid),
		nil,
	)

	req.SetBasicAuth("jaynagpaul", key)

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("Expected 200 got %d. %s", resp.StatusCode, string(b))
	}
	return nil
}

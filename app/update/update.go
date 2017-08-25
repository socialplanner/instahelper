// Package update provides methods to find the version of the package and update it.
package update

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/coreos/go-semver/semver"
	. "github.com/socialplanner/instahelper/app/log"
)

const (
	baseURL = "https://api.github.com/repos/zyedidia/micro/"
)

// Update replaces the binary runnning this command with a newer one fetched from github releases
func Update(version string) (*Asset, error) {
	releases, err := ListReleases()

	if err != nil {
		Log.Error(err)
		return nil, err
	}

	currentVer, err := semver.NewVersion(version)

	if err != nil {
		Log.Error(err)
		return nil, err
	}

	for _, r := range releases {
		ver, err := semver.NewVersion(strings.Replace(
			r.Version,
			"v",
			"",
			-1,
		))

		if err != nil {
			continue
		}

		if !ver.LessThan(*currentVer) {
			if asset := pickAsset(r.Assets); asset != nil {
				err := download(asset.DownloadURL)
				return asset, err
			}

		}
	}
	return nil, errors.New("No available download")
}

// ListReleases collects information about github releases
func ListReleases() ([]Release, error) {
	resp, err := http.Get(baseURL + "releases")

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)

	var releases []Release

	err = json.Unmarshal(b, &releases)

	if err != nil {
		return nil, err
	}

	return releases, err
}

// streams the url and copys it to os.Executable()
func download(url string) error {
	resp, err := http.Get(url)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	path, err := os.Executable()

	if err != nil {
		return err
	}

	path, err = filepath.EvalSymlinks(path)

	if err != nil {
		return err
	}

	f, err := os.Create(path)

	if err != nil {
		return err
	}

	_, err = io.Copy(f, resp.Body)

	return err
}

func pickAsset(assets []Asset) *Asset {
	for _, asset := range assets {
		name := asset.Name
		if ext := filepath.Ext(name); ext != "" {
			name = strings.Replace(asset.Name, filepath.Ext(name), "", -1)
		}

		info := strings.Split(name, "-")

		goos, goarch := info[2], info[3]

		if goos == "macos" {
			goos = "darwin"
		}

		if goos == runtime.GOOS {
			switch goarch {
			case "32", "arm":
				goarch = "386"
			case "64", "arm64":
				goarch = "amd64"
			}

			if goarch == runtime.GOARCH {
				return &asset
			}
		}
	}

	return nil
}

// Release is a github release
type Release struct {
	Name        string `json:"name"`
	Description string `json:"body"`
	Version     string `json:"tag_name"`
	URL         string `json:"url"`

	Assets []Asset `json:"assets"`
}

// Asset is a single download and it's various info for a github release
type Asset struct {
	Name        string `json:"name"`
	DownloadURL string `json:"browser_download_url"`
	Size        int    `json:"size"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

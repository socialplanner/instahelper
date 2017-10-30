// Package update provides methods to find the version of the package and update it.
package update

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/socialplanner/instahelper/app/notifications"

	"github.com/sirupsen/logrus"

	"github.com/Masterminds/semver"
	"github.com/socialplanner/instahelper/app/config"
)

// DIFFERENTVERSION tells instahelper if the current version differs from the downloaded version
var DIFFERENTVERSION = false

func init() {
	go func() {
		for {
			time.Sleep(time.Minute * 30)

			// Check config every 30 minutes instead of once, so the settings are allowed to change without a restart.
			c, err := config.Config()
			if err != nil {
				logrus.Error(err)
				continue
			}

			// Automatic Update Enabled
			if c.AutomaticUpdates {
				if asset, err := ToLatest(VERSION); asset != nil && err == nil {
					logrus.Info("Updated to the latest version!")
					if err := notifications.NewNotification("Updated to the latest version! Make sure to restart the app to experience the latest changes!", "/update"); err != nil {
						logrus.Error(err)
					}
					// Stop checking for updates.
					break
				} else {
					if err.Error() != "No available download" {
						logrus.Error(err)
						continue
					}
					logrus.Debug("No available update")
				}
			}
		}
	}()
}

const (
	baseURL = "https://api.github.com/repos/socialplanner/instahelper/"
)

// ToLatest replaces the binary runnning this command with the latest binary available on github releases
// If version is an empty string, it will fetch and replace the binary no matter what.
// Else it will only replace it if the version is greater than the current version
func ToLatest(version string) (*Asset, error) {
	// HACKY
	if version == "" {
		version = "0.0.1"
	}

	asset, err := HigherVersion(version)

	if err != nil {
		return nil, err
	}
	return asset, download(asset.DownloadURL)
}

// To will replace the current binary with the binary with version ver.
// Will return error if ver not found
// Update is allowed to update to a lower version
func To(ver string) (*Asset, error) {
	releases, err := ListReleases()

	if err != nil {
		return nil, err
	}

	for _, r := range releases {
		if ver == r.Version {
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
	if r := config.Get("releases"); r != nil {
		if rel, ok := r.([]Release); ok {
			return rel, nil
		}
	}

	resp, err := http.Get(baseURL + "releases")

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		time.Sleep(3 * time.Second)
		return ListReleases()
	}

	b, _ := ioutil.ReadAll(resp.Body)

	var releases []Release

	err = json.Unmarshal(b, &releases)

	if err != nil {
		return nil, err
	}

	config.Set("releases", releases, 30*time.Minute)

	return releases, nil
}

// HigherVersion returns the asset if there is a higher version available to download.
func HigherVersion(version string) (*Asset, error) {
	releases, err := ListReleases()

	if err != nil {
		return nil, err
	}

	currentVer, err := semver.NewVersion(version)

	if err != nil {
		return nil, err
	}

	for _, r := range releases {
		// Skip release if it is a prerelease
		if r.PreRelease {
			//continue
		}

		ver, err := semver.NewVersion(strings.Replace(
			r.Version,
			"v",
			"",
			-1,
		))

		if err != nil {
			continue
		}

		if !ver.LessThan(currentVer) {
			if asset := pickAsset(r.Assets); asset != nil {
				return asset, nil
			}
		}
	}

	return nil, errors.New("No available download")
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

	// fix symbolic links
	path, err = filepath.EvalSymlinks(path)

	if err != nil {
		return err
	}

	dir := filepath.Dir(path)

	// Downloads file to dir/tmp-instahelper/temp.zip
	zipFile := filepath.Join(dir, "temp.zip")

	f, err := os.Create(zipFile)

	if err != nil {
		return err
	}

	defer f.Close()
	defer os.RemoveAll(filepath.Join(dir, "tmp-instahelper"))

	_, err = io.Copy(f, resp.Body)

	if err != nil {
		return err
	}

	// unzip and deletes the file
	err = unzip(zipFile, filepath.Join(dir, "tmp-instahelper"))

	if err != nil {
		return err
	}

	// Open file to the executable
	exe, err := os.Create(path)

	if err != nil {
		return err
	}
	defer f.Close()

	var name string

	if runtime.GOOS == "windows" {
		name = "instahelper.exe"
	} else {
		name = "instahelper"
	}

	// Open file to the new executable
	// Uses name to account for .exe vs no suffix
	f, err = os.Open(filepath.Join(dir, "tmp-instahelper", name))

	if err != nil {
		return err
	}
	defer f.Close()

	// Copies it over
	_, err = io.Copy(exe, f)

	if err != nil {
		return err
	}

	DIFFERENTVERSION = true
	return nil
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

// unzip a .zip file located at src to dest
func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	// Delete src
	return os.Remove(src)
}

// Release is a github release
type Release struct {
	Name        string  `json:"name,omitempty"`
	Description string  `json:"body,omitempty"`
	Version     string  `json:"tag_name,omitempty"`
	URL         string  `json:"url,omitempty"`
	ID          int     `json:"id,omitempty"`
	PreRelease  bool    `json:"prerelease,omitempty"`
	PublishedAt string  `json:"published_at,omitempty"`
	Assets      []Asset `json:"assets,omitempty"`
	InfoURL     string  `json:"html_url,omitempty"`
}

// Semver will return a semver.Version from the release version
func (r *Release) Semver() (*semver.Version, error) {
	return semver.NewVersion(r.Version)
}

// Asset is a single download and it's various info for a github release
type Asset struct {
	Name        string `json:"name"`
	DownloadURL string `json:"browser_download_url"`
	Size        int    `json:"size"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

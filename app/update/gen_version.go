// +build ignore

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/coreos/go-semver/semver"
)

var tmpl = `package update

// VERSION of instahelper
var VERSION = "%s"`

func main() {
	path, _ := filepath.Abs(".")
	path = filepath.Join(path, "app", "update", "version.go")

	out := fmt.Sprintf(tmpl, os.Args[1])

	// Check if version is semver
	semver.New(os.Args[1])

	err := ioutil.WriteFile(path, []byte(out), 0644)
	if err != nil {
		fmt.Errorf("", err)
	}

}

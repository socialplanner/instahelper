// +build windows

package config

import (
	"os"
	"path/filepath"
)

// ConfigDir is the folder where we should store our config
var ConfigDir = filepath.Join(os.Getenv("APPDATA"), "instahelper")

// Package config implements several helper functions to abstract away getting information stored.
// Database used is BoltDB for fast concurrently safe actions.
// To further abstract away BoltDB which is a K/V store we use github.com/asdine/storm
// You may read more about these at https://github.com/asdine/storm & https://github.com/boltdb/bolt
package config

import (
	"crypto/rand"
	"os"
	"path/filepath"

	"github.com/asdine/storm"
)

// DB is the wrapper around BoltDB. It contains an instance of BoltDB and uses it to perform all the
// needed operations
var DB *storm.DB

// Open opens a database at the set location
func Open() error {
	var existed bool

	if _, err := os.Stat(filepath.Join(ConfigDir, "instahelper.db")); err == nil {
		existed = true
	}

	db, err := storm.Open(filepath.Join(ConfigDir, "instahelper.db"))

	if err != nil {
		return err
	}

	if !existed {
		Init()
	}

	Migrate()
	DB = db

	return nil
}

// Close the database. Should be deferred after opening the database
func Close() error {
	return DB.Close()
}

// Init will create the InstahelperConfig
func Init() error {
	key := make([]byte, 64)
	rand.Read(key)

	return DB.Save(&InstahelperConfig{
		ID:     1,
		AESKey: key,
	})
}

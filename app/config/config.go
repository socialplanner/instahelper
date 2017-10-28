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
	"github.com/asdine/storm/codec/gob"
)

// DB is the wrapper around BoltDB. It contains an instance of BoltDB and uses it to perform all the
// needed operations
var DB *storm.DB

// Open opens a database at the set location
func Open() error {
	// Creates the directory, ignores err.
	os.Mkdir(ConfigDir, 0777)

	var existed bool

	if _, err := os.Stat(filepath.Join(ConfigDir, "instahelper.db")); err == nil {
		existed = true
	}

	// Uses gob codec
	db, err := storm.Open(filepath.Join(ConfigDir, "instahelper.db"), storm.Codec(gob.Codec))

	if err != nil {
		return err
	}

	DB = db

	err = Migrate()

	if err != nil {
		return err
	}

	if !existed {
		createConfig()
	}

	// Refresh the cache
	DB.Drop("cache")

	return nil
}

// Close the database. Should be deferred after opening the database
func Close() error {
	return DB.Close()
}

// createConfig will create the InstahelperConfig
func createConfig() error {
	key := make([]byte, 32)
	rand.Read(key)

	c := &[]InstahelperConfig{}

	if DB.All(c); len(*c) > 0 {
		return nil
	}
	// Defaults
	return DB.Save(&InstahelperConfig{
		ID:               1,
		AESKey:           key,
		Analytics:        true,
		AutomaticUpdates: true,
	})
}

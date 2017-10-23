package config

import (
	"time"
)

// Set is a quick interface to a key value cache.
// dur is the amount of time to cache the object before deleting
// if 0 then will cache infinitely
func Set(key string, value interface{}, dur time.Duration) error {

	if err := DB.Set("cache", key, value); err != nil {
		return err
	}

	if dur.Seconds() > 0 {
		go func() {
			time.Sleep(dur)
			Delete(key)
		}()
	}

	return nil
}

// Get a value from the key value cache.
// Returns nil if not found.
func Get(key string) interface{} {
	var i interface{}

	DB.Get("cache", key, &i)

	return i
}

// Delete a value from the key value cache.
func Delete(key string) error {
	return DB.Delete("cache", key)
}

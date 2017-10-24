package handlers

import (
	"github.com/sirupsen/logrus"
	"github.com/socialplanner/instahelper/app/config"
)

// If the user has google analytics enabled
func analyticsEnabled() bool {
	c, err := config.Config()

	if err != nil {
		// Default
		logrus.Error(err)
		return true
	}

	return c.Analytics
}

package config

import (
	"time"
)

// InstahelperConfig is information about the Instahelper package as a whole
type InstahelperConfig struct {
	// AESKey used to encrypt password and Account.CachedInsta
	AESKey []byte

	ID int `storm:"id"`

	Port int
}

// Account is an Instagram Account
type Account struct {
	ID int `storm:"id,increment"`

	Username string `storm:"unique"`
	Password string

	// Cached GoInsta object
	CachedInsta []byte

	// Settings is not inline to be able to copy over settings
	Settings Settings

	// AddedAt is when the user added this account
	AddedAt time.Time `storm:"index"`
}

// Settings for a given account
type Settings struct {
	FollowsPerDay   int
	CommentsPerDay  int
	LikesPerDay     int
	UnfollowsPerDay int

	// Proxy to make requests with
	Proxy string

	// UnfollowAt is the number of follows when the bot should start unfollowing
	UnfollowAt int
	// UnfollowNonFollowers will decide if we unfollow those who do not follow after one day
	UnfollowNonFollowers bool

	// Tags to follow, comment, or like
	Tags []string
	// CommentList is the list of comments to choose from when commenting
	CommentList []string

	// Blacklist is a list of accounts to avoid following, commenting, or liking
	Blacklist []string
	// Whitelist is the list of users to only follow, comment, and like on
	Whitelist []string

	// FollowPrivate will decide if we follow private accounts
	FollowPrivate bool
}

var models = []interface{}{
	&Account{}, &InstahelperConfig{},
}

// Migrate will reindex all fields
func Migrate() {

	for _, m := range models {
		DB.ReIndex(m)
	}
}

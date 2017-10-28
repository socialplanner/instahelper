// Package insta provides a wrapper over goinsta
package insta

import (
	"time"

	"github.com/ahmdrz/goinsta"
	"github.com/ahmdrz/goinsta/store"
	"github.com/socialplanner/instahelper/app/config"
)

// Accounts is a map of all Instagram accounts currently in use.
// username => Instagram
// This is used to reuse instagram connections so an account is only logged in once.
// TODO implement wrapper funcs and replace Login, etc...
var Accounts = map[string]*goinsta.Instagram{}

// Acc is a wrapper function to return a Instagram Account by name.
// Automatically logs in the account if not logged and updates respective fields in the Account model.
func Acc(name string) (*goinsta.Instagram, error) {
	if ig, ok := Accounts["name"]; ok {
		return ig, nil
	}

	var acc = config.Account{}

	if err := config.DB.One("Username", name, &acc); err != nil {
		return nil, err
	}

	ig, err := Import(acc.CachedInsta)

	if err != nil {
		return nil, err
	}

	pass, err := Decrypt(acc.Password)
	if err != nil {
		return nil, err
	}

	// Password has been updated
	if p := ig.Informations.Password; p != pass {
		ig.Informations.Password = p
		ig.IsLoggedIn = false

		b, err := ExportCached(ig)

		if err != nil {
			return nil, err
		}

		acc.CachedInsta = b
	}

	// Save it to the map
	Accounts[name] = ig

	if !ig.IsLoggedIn {
		if err := ig.Login(); err != nil {
			return nil, err
		}
	}

	// Only call at the most every 25 minutes.
	if time.Now().Sub(acc.LastAccess).Minutes() > 25 {
		// We are using a request that should pass under most circumstances, to check
		// if the account is still logged in. The password could've changed between now and then, or any number of things.
		p, err := ig.GetProfileData()

		if err != nil {
			if err := ig.Login(); err != nil {
				return nil, err
			}
		}

		user := p.User

		// Update vanity profile data
		acc.Bio = user.Biography
		acc.Private = user.IsPrivate
		acc.ProfilePic = user.HDProfilePicURLInfo.URL

		following, err := ig.SelfTotalUserFollowing()
		if err != nil {
			return nil, err
		}

		followers, err := ig.SelfTotalUserFollowers()
		if err != nil {
			return nil, err
		}

		acc.Following = len(following.Users)
		acc.Followers = len(followers.Users)
		acc.LastUpdate = time.Now()
	}

	acc.LastAccess = time.Now()

	err = acc.Update()

	if err != nil {
		return nil, err
	}
	return ig, nil
}

// ExportCached will export the cached instagram object using the configs AESKey.
func ExportCached(ig *goinsta.Instagram) ([]byte, error) {
	c, err := config.Config()

	if err != nil {
		return []byte{}, err
	}
	return store.Export(ig, c.AESKey)
}

// Login will connect to Instagram through a proxy if one is passed
func Login(username, password, proxy string) (*goinsta.Instagram, error) {
	var ig *goinsta.Instagram

	// If proxy passed create a goinsta connection using that proxy
	if proxy == "" {
		ig = goinsta.New(username, password)
	} else {
		ig = goinsta.NewViaProxy(username, password, password)
	}

	err := ig.Login()

	if err != nil {
		return nil, err
	}

	return ig, nil
}

// Import an account from it's cached bytes
func Import(b []byte) (*goinsta.Instagram, error) {
	c, err := config.Config()

	if err != nil {
		return nil, err
	}

	return store.Import(b, c.AESKey)

}

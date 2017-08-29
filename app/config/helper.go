package config

// Helpers to get all of a specific model

// Accounts will return all accounts in the database
func Accounts() (*[]Account, error) {

	accs := &[]Account{}

	err := DB.All(accs)

	return accs, err
}

// Config will return the instahelper config
func Config() (*InstahelperConfig, error) {

	c := &InstahelperConfig{}

	err := DB.One("ID", 1, c)

	return c, err
}

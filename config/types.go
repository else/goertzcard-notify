package config

type Config struct {
	Config   struct{} `yaml:"config"`
	Accounts []struct {
		Cards []struct {
			Ean           string `yaml:"ean" validate:"required"`
			MinimumAmount string `yaml:"minimumAmount" validate:"required"`
			Notify        string `yaml:"notify"`
		} `yaml:"cards" validate:"required"`
		Credentials struct {
			Password string `yaml:"password" validate:"required"`
			Username string `yaml:"username" validate:"required"`
		} `yaml:"credentials" validate:"required"`
		Notifier struct {
			Pushover struct {
				Type  string `yaml:"type" validate:"required"`
				Token string `yaml:"token" validate:"required"`
				User  string `yaml:"user" validate:"required"`
			} `yaml:"pushover" validate:"required"`
		} `yaml:"notifier" validate:"required"`
		Owner string `yaml:"owner" validate:"required"`
	} `yaml:"accounts" validate:"required"`
}

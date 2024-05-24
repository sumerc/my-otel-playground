package myauthext

import "fmt"

// Config represents the receiver config settings within the collector's config.yaml
type Config struct {
	Dummy string `mapstructure:"dummy"`
}

// Validate checks if the receiver configuration is valid
func (cfg *Config) Validate() error {
	fmt.Println("validate->", cfg.Dummy)
	return nil
}

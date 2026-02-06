package config

import "fmt"

func ValidateConfig(c *Litmos) error {
	if c.ApiKey == "" {
		return fmt.Errorf("api-key is required")
	}
	if c.Source == "" {
		return fmt.Errorf("source is required")
	}
	return nil
}

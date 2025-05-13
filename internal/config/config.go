package config

import ierr "github.com/borisdvlpr/gotail/internal/error"

type Config struct {
	ExitNode     string `yaml:"exit_node"`
	SubnetRouter string `yaml:"subnet_router"`
	Subnets      string `yaml:"subnets"`
	Hostname     string `yaml:"hostname"`
	AuthKey      string `yaml:"auth_key"`
}

func (c *Config) Validate() error {
	if c.AuthKey == "" {
		return ierr.StatusError{Status: "auth key is required", StatusCode: 1}
	}

	if c.SubnetRouter == "y" && c.Subnets == "" {
		return ierr.StatusError{Status: "subnets are required when subnet router is enabled", StatusCode: 1}
	}

	if c.Hostname == "" {
		return ierr.StatusError{Status: "hostname is required", StatusCode: 1}
	}

	return nil
}

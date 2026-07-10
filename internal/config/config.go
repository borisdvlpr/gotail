package config

import ierror "github.com/borisdvlpr/gotail/internal/error"

type Config struct {
	ExitNode     string `yaml:"exit_node"`
	SubnetRouter string `yaml:"subnet_router"`
	Subnets      string `yaml:"subnets"`
	Hostname     string `yaml:"hostname"`
	Tags         string `yaml:"tags"`
	AuthKey      string `yaml:"auth_key"`
}

func (c *Config) Validate() error {
	if c.AuthKey == "" {
		return ierror.StatusError{Status: "auth key is required", StatusCode: 1}
	}

	if c.SubnetRouter == "y" && c.Subnets == "" {
		return ierror.StatusError{Status: "subnets are required when subnet router is enabled", StatusCode: 1}
	}

	if c.Hostname == "" {
		return ierror.StatusError{Status: "hostname is required", StatusCode: 1}
	}

	return nil
}

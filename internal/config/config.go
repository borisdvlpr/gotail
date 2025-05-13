package config

import "errors"

type Config struct {
	ExitNode     string `yaml:"exit_node"`
	SubnetRouter string `yaml:"subnet_router"`
	Subnets      string `yaml:"subnets"`
	Hostname     string `yaml:"hostname"`
	AuthKey      string `yaml:"auth_key"`
}

func (c *Config) Validate() error {
	if c.AuthKey == "" {
		return errors.New("auth_key is required")
	}

	if c.SubnetRouter == "y" && c.Subnets == "" {
		return errors.New("subnets are required when subnet_router is enabled")
	}

	if c.Hostname == "" {
		return errors.New("hostname is required")
	}

	return nil
}

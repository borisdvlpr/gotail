package config

import (
	"fmt"

	ierror "github.com/borisdvlpr/gotail/internal/error"
)

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

	if c.ExitNode != "" && c.ExitNode != "y" && c.ExitNode != "n" {
		return ierror.StatusError{Status: fmt.Sprintf(`exit_node must be "y" or "n", got %q`, c.ExitNode), StatusCode: 1}
	}

	if c.SubnetRouter != "" && c.SubnetRouter != "y" && c.SubnetRouter != "n" {
		return ierror.StatusError{Status: fmt.Sprintf(`subnet_router must be "y" or "n", got %q`, c.SubnetRouter), StatusCode: 1}
	}

	return nil
}

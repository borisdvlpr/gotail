package config

type Config struct {
	ExitNode     bool   `yaml:"exit_node"`
	SubnetRouter bool   `yaml:"subnet_router"`
	Subnets      string `yaml:"subnets"`
	Hostname     string `yaml:"hostname"`
	AuthKey      string `yaml:"auth_key"`
}

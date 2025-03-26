package config

type Config struct {
	ExitNode     bool   `mapstructure:"exit_node"`
	SubnetRouter bool   `mapstructure:"subnet_router"`
	Subnets      string `mapstructure:"subnets"`
	Hostname     string `mapstructure:"hostname"`
	AuthKey      string `mapstructure:"auth_key"`
}

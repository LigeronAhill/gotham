package config

import "github.com/spf13/viper"

func Init() *viper.Viper {
	c := viper.New()
	c.SetEnvPrefix("{{ .EnvPrefixLower }}")
	c.SetDefault("host", "localhost")
	c.SetDefault("port", 3000)
	c.SetDefault("admin_password", "{{ .AdminPwd }}")
	c.AutomaticEnv()
	return c
}

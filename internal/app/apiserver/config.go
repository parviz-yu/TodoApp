package apiserver

type Config struct {
	BindAddr    string `toml:"bind_addr"`
	LogLevel    string `toml:"log_level"`
	DatabaseURL string `toml:"database_url"`
	SessioKey   string `toml:"session_key"`
}

// NewConfig return new Config instance
func NewConfig() *Config {
	return &Config{}
}

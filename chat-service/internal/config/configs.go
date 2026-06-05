package config

type Config struct {
	Port         string   `mapstructure:"PORT"`
	AppName      string   `mapstructure:"APP_NAME"`
	DBSchemaName string   `mapstructure:"DB_SCHEMA_NAME"`
	PG           PGConfig `mapstructure:"DB"`
	//JWTConfig    JWTConfig `mapstructure:"JWT"`
	AuthServiceURL string `mapstructure:"AUTH_SERVICE_URL"`
	Env            string `mapstructure:"ENV"`
	LogLevel       string `mapstructure:"LOG_LEVEL"`
}

type PGConfig struct {
	URL string `mapstructure:"URL"`
}

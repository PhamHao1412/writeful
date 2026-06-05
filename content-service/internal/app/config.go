package app

type Config struct {
	Port           string   `mapstructure:"PORT"`
	Env            string   `mapstructure:"ENV"`
	AppName        string   `mapstructure:"APP_NAME"`
	DBSchemaName   string   `mapstructure:"DB_SCHEMA_NAME"`
	PG             PGConfig `mapstructure:"DB"`
	AuthServiceURL string   `mapstructure:"AUTH_SERVICE_URL"`
	LogLevel       string   `mapstructure:"LOG_LEVEL"`
}

type PGConfig struct {
	URL string `mapstructure:"URL"`
}

package app

type Config struct {
	Port         string   `mapstructure:"PORT"`
	AppName      string   `mapstructure:"APP_NAME"`
	DBSchemaName string   `mapstructure:"DB_SCHEMA_NAME"`
	PG           PGConfig `mapstructure:"DB"`
	//JWTConfig    JWTConfig `mapstructure:"JWT"`
}

type PGConfig struct {
	URL string `mapstructure:"URL"`
}

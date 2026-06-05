package app

type Config struct {
	Port         string    `mapstructure:"PORT"`
	AppName      string    `mapstructure:"APP_NAME"`
	DBSchemaName string    `mapstructure:"DB_SCHEMA_NAME"`
	PG           PGConfig  `mapstructure:"DB"`
	JWTConfig    JWTConfig `mapstructure:"JWT"`
}

type PGConfig struct {
	URL string `mapstructure:"URL"`
}

type JWTConfig struct {
	AdminKey                         string `mapstructure:"ADMIN_KEY"`
	Secret                           string `mapstructure:"SECRET" json:"Secret"`
	Issuer                           string `mapstructure:"ISSUER" json:"Issuer"`
	Audience                         string `mapstructure:"AUDIENCE" json:"Audience"`
	DurationInMinutes                int    `mapstructure:"DURATION_IN_MINUTES" json:"DurationInMinutes"`
	DurationInMinutesForRefreshToken int    `mapstructure:"DURATION_IN_MINUTES_FOR_REFRESH_TOKEN" json:"DurationInMinutesForRefreshToken"`
	Alg                              string `mapstructure:"ALG" json:"alg"`
	Kid                              string `mapstructure:"KID" json:"kid"`
}

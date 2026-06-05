package app

type Config struct {
	Port         string           `mapstructure:"PORT"`
	AppName      string           `mapstructure:"APP_NAME"`
	DBSchemaName string           `mapstructure:"DB_SCHEMA_NAME"`
	PG           PGConfig         `mapstructure:"DB"`
	Cloudinary   CloudinaryConfig `mapstructure:"CLOUDINARY"`
}

type PGConfig struct {
	URL string `mapstructure:"URL"`
}

type CloudinaryConfig struct {
	ApiKey    string `mapstructure:"API_KEY"`
	Secret    string `mapstructure:"SECRET"`
	CloudName string `mapstructure:"CLOUD_NAME"`
}

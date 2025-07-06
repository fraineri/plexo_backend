package settings

type Settings struct {
	Name    string `env:"APP_INFO_NAME,required"`
	Version string `env:"APP_INFO_VERSION,required"`
}

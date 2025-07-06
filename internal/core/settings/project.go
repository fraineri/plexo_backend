package settings

type ProjectSettings struct {
	Environment string `env:"PROJECT_ENVIRONMENT,required"`
	Debug       bool   `env:"PROJECT_DEBUG,required"`
}

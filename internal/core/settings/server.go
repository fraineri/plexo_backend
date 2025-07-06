package settings

type ServerSettings struct {
	Port string `env:"SERVER_PORT,required" envDefault:"8080"`
}

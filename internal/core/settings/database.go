package settings

import "fmt"

type DatabaseSettings struct {
	PostgresUserWrite     string `env:"POSTGRES_USER_WRITE,required"`
	PostgresPasswordWrite string `env:"POSTGRES_PASSWORD_WRITE,required"`
	PostgresHostWrite     string `env:"POSTGRES_HOST_WRITE,required"`
	PostgresPortWrite     int    `env:"POSTGRES_PORT_WRITE,required"`
	PostgresDBWrite       string `env:"POSTGRES_DB_WRITE,required"`
	PostgresSSLModeWrite  string `env:"POSTGRES_SSL_MODE_WRITE,required"`

	PostgresUserRead     string `env:"POSTGRES_USER_READ,required"`
	PostgresPasswordRead string `env:"POSTGRES_PASSWORD_READ,required"`
	PostgresHostRead     string `env:"POSTGRES_HOST_READ,required"`
	PostgresPortRead     int    `env:"POSTGRES_PORT_READ,required"`
	PostgresDBRead       string `env:"POSTGRES_DB_READ,required"`
	PostgresSSLModeRead  string `env:"POSTGRES_SSL_MODE_READ,required"`
}

func (d *DatabaseSettings) GetPostgresWriteDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		d.PostgresHostWrite,
		fmt.Sprintf("%d", d.PostgresPortWrite),
		d.PostgresUserWrite,
		d.PostgresPasswordWrite,
		d.PostgresDBWrite,
		d.PostgresSSLModeWrite,
	)
}
func (d *DatabaseSettings) GetPostgresReadDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		d.PostgresHostRead,
		fmt.Sprintf("%d", d.PostgresPortRead),
		d.PostgresUserRead,
		d.PostgresPasswordRead,
		d.PostgresDBRead,
		d.PostgresSSLModeRead,
	)
}

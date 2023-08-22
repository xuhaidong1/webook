package config

type WebookConfig struct {
	DB    DBConfig
	Redis RedisConfig
	Port  string
}

type DBConfig struct {
	DSN string
}

type RedisConfig struct {
	Addr string
}

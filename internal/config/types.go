package config

type Config struct {
	Postgres   PostgresCfg   `mapstructure:"postgres"`
	HTTPServer HTTPServerCfg `mapstructure:"http_server"`
	Redis      RedisCfg      `mapstructure:"redis"`
	RabbitMq   RabbitMqCfg   `mapstructure:"rabbit_mq"`
}

type PostgresCfg struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Name     string `mapstructure:"name"`
	Password string `mapstructure:"password"`
}

type HTTPServerCfg struct {
	Address     string `mapstructure:"address"`
	Timeout     int    `mapstructure:"timeout"`
	IdleTimeout int    `mapstructure:"idle_timeout"`
}

type RedisCfg struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
}

type RabbitMqCfg struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

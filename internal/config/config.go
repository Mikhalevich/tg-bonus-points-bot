package config

import (
	"time"
)

type ConsumerBot struct {
	LogLevel    string      `yaml:"log_level" required:"true"`
	Tracing     Tracing     `yaml:"tracing" required:"true"`
	Bot         Bot         `yaml:"bot" required:"true"`
	Postgres    Postgres    `yaml:"postgres" required:"true"`
	ButtonRedis ButtonRedis `yaml:"button_redis" required:"true"`
	CartRedis   CartRedis   `yaml:"cart_redis" required:"true"`
}

type ManagerHTTPService struct {
	LogLevel string   `yaml:"log_level" required:"true"`
	Tracing  Tracing  `yaml:"tracing" required:"true"`
	Bot      Bot      `yaml:"bot" required:"true"`
	Postgres Postgres `yaml:"postgres" required:"true"`
	HTTPPort int      `yaml:"http_port" required:"true"`
}

type Bot struct {
	Token string `yaml:"token" required:"true"`
}

type Tracing struct {
	Endpoint    string `yaml:"endpoint" required:"true"`
	ServiceName string `yaml:"service_name" required:"true"`
}

type Postgres struct {
	Connection string `yaml:"connection" required:"true"`
}

type ButtonRedis struct {
	Addr string        `yaml:"addr" required:"true"`
	Pwd  string        `yaml:"pwd" required:"true"`
	DB   int           `yaml:"db" required:"true"`
	TTL  time.Duration `yaml:"ttl" required:"true"`
}

type CartRedis struct {
	Addr string        `yaml:"addr" required:"true"`
	Pwd  string        `yaml:"pwd" required:"true"`
	DB   int           `yaml:"db" required:"true"`
	TTL  time.Duration `yaml:"ttl" required:"true"`
}

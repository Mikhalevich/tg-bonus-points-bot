package config

type ScheduleBot struct {
	LogLevel string   `yaml:"log_level" required:"true"`
	Tracing  Tracing  `yaml:"tracing" required:"true"`
	Bot      Bot      `yaml:"bot" required:"true"`
	Postgres Postgres `yaml:"postgres" required:"true"`
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

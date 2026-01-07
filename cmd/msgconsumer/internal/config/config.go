package config

type Config struct {
	LogLevel string  `yaml:"log_level" required:"true"`
	Tracing  Tracing `yaml:"tracing" required:"true"`
	Bot      Bot     `yaml:"bot" required:"true"`
	Kafka    Kafka   `yaml:"kafka" required:"true"`
}

type Tracing struct {
	Endpoint    string `yaml:"endpoint" required:"true"`
	ServiceName string `yaml:"service_name" required:"true"`
}

type Bot struct {
	Token        string `yaml:"token" required:"true"`
	PaymentToken string `yaml:"payment_token" required:"true"`
}

type Kafka struct {
	Brokers       []string `yaml:"brokers" required:"true"`
	Topic         string   `yaml:"topic" required:"true"`
	ConsumerGroup string   `yaml:"consumer_group" required:"true"`
}

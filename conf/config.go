package conf

type Config struct {
	App struct {
		IsPrd bool `yaml:"isPrd"`
	} `yaml:"app"`
	Binance struct {
		ApiKey    string `yaml:"ApiKey"`
		SecretKey string `yaml:"SecretKey"`
	} `yaml:"binance"`
	Wx struct {
		MessagePushUrl string `yaml:"messagePushUrl"`
	} `yaml:"wx"`
}

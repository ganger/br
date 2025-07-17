package conf

type Config struct {
	Binance struct {
		ApiKey    string `yaml:"ApiKey"`
		SecretKey string `yaml:"SecretKey"`
	} `yaml:"binance"`
	Wx struct {
		MessagePushUrl string `yaml:"messagePushUrl"`
	} `yaml:"wx"`
}

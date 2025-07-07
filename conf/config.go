package conf

type Config struct {
	Binance struct {
		ReadApiKey    string `yaml:"readApiKey"`
		ReadSecretKey string `yaml:"readSecretKey"`
	} `yaml:"binance"`
	Wx struct {
		MessagePushUrl string `yaml:"messagePushUrl"`
	} `yaml:"wx"`
}

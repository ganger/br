package conf

type Config struct {
	App struct {
		IsPrd bool `yaml:"isPrd"`
	} `yaml:"app"`
	Binance struct {
		ApiKey     string `yaml:"ApiKey"`
		SecretKey  string `yaml:"SecretKey"`
		ApiKey2    string `yaml:"ApiKey2"`
		SecretKey2 string `yaml:"SecretKey2"`
	} `yaml:"binance"`
	Wx struct {
		MessagePushUrl string `yaml:"messagePushUrl"`
	} `yaml:"wx"`
}

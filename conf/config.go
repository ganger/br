package conf

type Config struct {
	App struct {
		MessagePushUrl string  `yaml:"messagePushUrl"`
		WalletAddress  string  `yaml:"walletAddress"`
		ReadApiKey     string  `yaml:"readApiKey"`
		ReadSecretKey  string  `yaml:"readSecretKey"`
		JlpBalance     float64 `yaml:"jlpBalance"`
		ShortBTC       float64 `yaml:"shortBTC"`
		ShortUSDC      float64 `yaml:"shortUSDC"`
		IsPrd          bool    `yaml:"isPrd"`
		TradeApiKey    string  `yaml:"tradeApiKey"`
		TradeSecretKey string  `yaml:"tradeSecretKey"`
		IgnoreWX       bool    `yaml:"ignoreWX"`
		InitTime       string  `yaml:"initTime"`
	} `yaml:"app"`
}

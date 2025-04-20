package config

type Config struct {
	Credentials Credentials `koanf:"credentials"`
}

type Credentials struct {
	APIKey   string `json:"apikey"`
	Username string `json:"username"`
	Password string `json:"password"`
}

package config

type Config struct {
	Credentials Credentials `koanf:"credentials" json:"credentials"`
	Download    Download    `koanf:"download" json:"download"`
	Network     Network     `koanf:"network" json:"network"`
}

type Credentials struct {
	Username string `koanf:"username" json:"username"`
	Password string `koanf:"password" json:"password"`
}

type Download struct {
	StorageDirectory string `koanf:"storage_directory" json:"storage_directory"`
}

type Network struct {
	XAppVar   string `koanf:"x_app_var" json:"x_app_var"`
	XApiKey   string `koanf:"x_api_key" json:"x_api_key"`
	UserAgent string `koanf:"user_agent" json:"user_agent"`
}

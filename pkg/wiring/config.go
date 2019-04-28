package wiring

// Config wraps app configs.
type Config struct {
	HTTP struct {
		Port int `envconfig:"default=3001"`
	}
	LOG struct {
		Level string
	}
	BITLY struct {
		Token string
	}
	TRANSMIT struct {
		Apikey string
		Secret string
	}
}

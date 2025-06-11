package config

type Config struct {
	Next        string
	Previous    string
	CurrentArgs []string
}

func New() *Config {
	return &Config{
		Next:        "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20",
		Previous:    "",
		CurrentArgs: []string{},
	}
}

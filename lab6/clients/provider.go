package clients

import (
	"fmt"
	"strings"
)

type Provider string

const (
	ProviderOpenWeather Provider = "openweather"
	ProviderGoogle      Provider = "google"
)

func ParseProvider(value string) (Provider, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", "openweather":
		return ProviderOpenWeather, nil
	case "google":
		return ProviderGoogle, nil
	default:
		return "", fmt.Errorf("unknown provider: %s", value)
	}
}

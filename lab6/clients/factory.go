package clients

import "fmt"

type WeatherClientFactory struct {
	clients map[Provider]WeatherDataClient
}

func NewWeatherClientFactory(openWeather WeatherDataClient, google WeatherDataClient) *WeatherClientFactory {
	return &WeatherClientFactory{
		clients: map[Provider]WeatherDataClient{
			ProviderOpenWeather: openWeather,
			ProviderGoogle:      google,
		},
	}
}

func (f *WeatherClientFactory) Get(provider Provider) (WeatherDataClient, error) {
	client, ok := f.clients[provider]
	if !ok || client == nil {
		return nil, fmt.Errorf("provider not configured: %s", provider)
	}

	return client, nil
}

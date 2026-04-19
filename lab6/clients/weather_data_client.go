package clients

import (
	"time"

	"github.com/shopspring/decimal"
)

type ForecastEntry struct {
	Time        time.Time
	Temperature decimal.Decimal
}

type WeatherDataClient interface {
	LocationCurrentTemperature(lat decimal.Decimal, lon decimal.Decimal) (temperature decimal.Decimal, err error)
	LocationForecast(lat decimal.Decimal, lon decimal.Decimal, hours int) ([]ForecastEntry, error)
}

package locations

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

type Location struct {
	Name string          `json:"name,omitempty"`
	Lat  decimal.Decimal `json:"lat"`
	Lon  decimal.Decimal `json:"lon"`
}

var knownCities = map[string]Location{
	"минск":    {Name: "Minsk", Lat: decimal.NewFromFloat(53.9006), Lon: decimal.NewFromFloat(27.5590)},
	"minsk":    {Name: "Minsk", Lat: decimal.NewFromFloat(53.9006), Lon: decimal.NewFromFloat(27.5590)},
	"лондон":   {Name: "London", Lat: decimal.NewFromFloat(51.5074), Lon: decimal.NewFromFloat(-0.1278)},
	"london":   {Name: "London", Lat: decimal.NewFromFloat(51.5074), Lon: decimal.NewFromFloat(-0.1278)},
	"токио":    {Name: "Tokyo", Lat: decimal.NewFromFloat(35.6762), Lon: decimal.NewFromFloat(139.6503)},
	"tokyo":    {Name: "Tokyo", Lat: decimal.NewFromFloat(35.6762), Lon: decimal.NewFromFloat(139.6503)},
	"шанхай":   {Name: "Shanghai", Lat: decimal.NewFromFloat(31.2304), Lon: decimal.NewFromFloat(121.4737)},
	"shanghai": {Name: "Shanghai", Lat: decimal.NewFromFloat(31.2304), Lon: decimal.NewFromFloat(121.4737)},
	"варшава":  {Name: "Warsaw", Lat: decimal.NewFromFloat(52.2297), Lon: decimal.NewFromFloat(21.0122)},
	"warsaw":   {Name: "Warsaw", Lat: decimal.NewFromFloat(52.2297), Lon: decimal.NewFromFloat(21.0122)},
}

func ResolveCity(name string) (Location, error) {
	key := strings.ToLower(strings.TrimSpace(name))
	if key == "" {
		return Location{}, fmt.Errorf("city name is empty")
	}

	location, ok := knownCities[key]
	if !ok {
		return Location{}, fmt.Errorf("unknown city: %s", name)
	}

	return location, nil
}

func ResolveLocation(city, latRaw, lonRaw string) (Location, error) {
	city = strings.TrimSpace(city)
	latRaw = strings.TrimSpace(latRaw)
	lonRaw = strings.TrimSpace(lonRaw)

	hasCity := city != ""
	hasCoords := latRaw != "" || lonRaw != ""

	if hasCity && hasCoords {
		return Location{}, fmt.Errorf("provide either city or coordinates, not both")
	}

	if hasCity {
		return ResolveCity(city)
	}

	if latRaw == "" || lonRaw == "" {
		return Location{}, fmt.Errorf("city or both lat and lon are required")
	}

	lat, errLat := decimal.NewFromString(latRaw)
	lon, errLon := decimal.NewFromString(lonRaw)
	if errLat != nil || errLon != nil {
		return Location{}, fmt.Errorf("invalid coordinates")
	}

	return Location{Lat: lat, Lon: lon}, nil
}

func ResolveCities(names []string) ([]Location, error) {
	if len(names) == 0 {
		return nil, fmt.Errorf("at least one city is required")
	}

	locations := make([]Location, 0, len(names))
	for _, name := range names {
		location, err := ResolveCity(name)
		if err != nil {
			return nil, err
		}
		locations = append(locations, location)
	}

	return locations, nil
}

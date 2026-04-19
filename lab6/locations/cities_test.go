package locations

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestResolveCityKnownCities(t *testing.T) {
	tests := []struct {
		input    string
		wantName string
	}{
		{input: "Минск", wantName: "Minsk"},
		{input: "minsk", wantName: "Minsk"},
		{input: "London", wantName: "London"},
		{input: "Токио", wantName: "Tokyo"},
		{input: "shanghai", wantName: "Shanghai"},
		{input: "Варшава", wantName: "Warsaw"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			location, err := ResolveCity(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if location.Name != tt.wantName {
				t.Fatalf("got name %q, want %q", location.Name, tt.wantName)
			}
		})
	}
}

func TestResolveCityUnknown(t *testing.T) {
	_, err := ResolveCity("Paris")
	if err == nil {
		t.Fatalf("expected error for unknown city")
	}
}

func TestResolveCityEmpty(t *testing.T) {
	_, err := ResolveCity("   ")
	if err == nil {
		t.Fatalf("expected error for empty city")
	}
}

func TestResolveLocationByCity(t *testing.T) {
	location, err := ResolveLocation("London", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if location.Name != "London" {
		t.Fatalf("got %q, want London", location.Name)
	}
}

func TestResolveLocationByCoordinates(t *testing.T) {
	location, err := ResolveLocation("", "10.5", "20.25")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !location.Lat.Equal(decimal.NewFromFloat(10.5)) || !location.Lon.Equal(decimal.NewFromFloat(20.25)) {
		t.Fatalf("unexpected coordinates: %s,%s", location.Lat, location.Lon)
	}
}

func TestResolveLocationBothCityAndCoordinates(t *testing.T) {
	_, err := ResolveLocation("London", "1", "2")
	if err == nil {
		t.Fatalf("expected error when both city and coordinates provided")
	}
}

func TestResolveLocationMissingInput(t *testing.T) {
	_, err := ResolveLocation("", "", "")
	if err == nil {
		t.Fatalf("expected error when location is missing")
	}
}

func TestResolveLocationInvalidCoordinates(t *testing.T) {
	_, err := ResolveLocation("", "abc", "20")
	if err == nil {
		t.Fatalf("expected invalid coordinates error")
	}
}

func TestResolveCities(t *testing.T) {
	locations, err := ResolveCities([]string{"Minsk", "Tokyo"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(locations) != 2 {
		t.Fatalf("got %d locations, want 2", len(locations))
	}
}

func TestResolveCitiesEmpty(t *testing.T) {
	_, err := ResolveCities(nil)
	if err == nil {
		t.Fatalf("expected error for empty list")
	}
}

func TestResolveCitiesWithUnknownCity(t *testing.T) {
	_, err := ResolveCities([]string{"Minsk", "Berlin"})
	if err == nil {
		t.Fatalf("expected error for unknown city in batch")
	}
}

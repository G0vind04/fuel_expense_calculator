// main.go
package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// OnStartup is called when the app starts, before the frontend has loaded
func (a *App) OnStartup(ctx context.Context) {
	a.ctx = ctx
}

// Route represents a route response from Google Maps
type Route struct {
	Distance string  `json:"distance"`
	Duration string  `json:"duration"`
	DistanceKm float64 `json:"distanceKm"`
}

// FuelPrice represents fuel price data
type FuelPrice struct {
	State   string  `json:"state"`
	Petrol  float64 `json:"petrol"`
	Diesel  float64 `json:"diesel"`
}

// TripCalculation represents the final calculation result
type TripCalculation struct {
	Route     Route     `json:"route"`
	FuelPrice FuelPrice `json:"fuelPrice"`
	FuelCost  float64   `json:"fuelCost"`
	FuelType  string    `json:"fuelType"`
	Mileage   float64   `json:"mileage"`
}

// GetRoute calculates route using OpenRouteService API (free, no key required)
func (a *App) GetRoute(origin, destination string) (Route, error) {
	// Using OpenRouteService which provides free routing without API key
	// Alternative: Use Nominatim for geocoding + distance calculation
	
	// For demo purposes, we'll use a geocoding service to get coordinates
	// then calculate approximate distance
	originCoords, err := a.geocodeLocation(origin)
	if err != nil {
		return Route{}, fmt.Errorf("failed to geocode origin: %v", err)
	}
	
	destCoords, err := a.geocodeLocation(destination)
	if err != nil {
		return Route{}, fmt.Errorf("failed to geocode destination: %v", err)
	}
	
	// Calculate haversine distance
	distance := haversineDistance(originCoords.Lat, originCoords.Lon, destCoords.Lat, destCoords.Lon)
	
	// Estimate road distance (typically 1.2-1.4x of straight line distance)
	roadDistance := distance * 1.3
	
	// Estimate duration (assuming average speed of 50 km/h)
	durationHours := roadDistance / 50
	
	return Route{
		Distance:   fmt.Sprintf("%.1f km", roadDistance),
		Duration:   formatDuration(durationHours),
		DistanceKm: roadDistance,
	}, nil
}

// Coordinates represents lat/lon coordinates
type Coordinates struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// geocodeLocation uses Nominatim (OpenStreetMap) to get coordinates
func (a *App) geocodeLocation(location string) (Coordinates, error) {
	baseURL := "https://nominatim.openstreetmap.org/search"
	
	params := url.Values{}
	params.Add("q", location+", India") // Append India for better results
	params.Add("format", "json")
	params.Add("limit", "1")
	params.Add("countrycodes", "in") // Restrict to India
	
	// Add user agent as required by Nominatim
	req, err := http.NewRequest("GET", baseURL+"?"+params.Encode(), nil)
	if err != nil {
		return Coordinates{}, err
	}
	req.Header.Set("User-Agent", "FuelCalculator/1.0")
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Coordinates{}, fmt.Errorf("failed to geocode: %v", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Coordinates{}, fmt.Errorf("failed to read response: %v", err)
	}
	
	var results []struct {
		Lat string `json:"lat"`
		Lon string `json:"lon"`
	}
	
	if err := json.Unmarshal(body, &results); err != nil {
		return Coordinates{}, fmt.Errorf("failed to parse response: %v", err)
	}
	
	if len(results) == 0 {
		return Coordinates{}, fmt.Errorf("location not found")
	}
	
	lat, err := strconv.ParseFloat(results[0].Lat, 64)
	if err != nil {
		return Coordinates{}, fmt.Errorf("invalid latitude")
	}
	
	lon, err := strconv.ParseFloat(results[0].Lon, 64)
	if err != nil {
		return Coordinates{}, fmt.Errorf("invalid longitude")
	}
	
	return Coordinates{Lat: lat, Lon: lon}, nil
}

// haversineDistance calculates the distance between two points on Earth
func haversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Earth's radius in kilometers
	
	dLat := (lat2 - lat1) * (3.14159265359 / 180)
	dLon := (lon2 - lon1) * (3.14159265359 / 180)
	
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*(3.14159265359/180))*math.Cos(lat2*(3.14159265359/180))*
		math.Sin(dLon/2)*math.Sin(dLon/2)
	
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	
	return R * c
}

// formatDuration formats duration in hours to readable format
func formatDuration(hours float64) string {
	if hours < 1 {
		minutes := int(hours * 60)
		return fmt.Sprintf("%d mins", minutes)
	}
	
	h := int(hours)
	m := int((hours - float64(h)) * 60)
	
	if m == 0 {
		return fmt.Sprintf("%d hour", h)
	}
	
	return fmt.Sprintf("%d hours %d mins", h, m)
}

// GetFuelPrices fetches current fuel prices for Indian states
func (a *App) GetFuelPrices(state string) (FuelPrice, error) {
	// This is a mock implementation. In reality, you'd integrate with:
	// - Indian Oil Corporation API
	// - Bharat Petroleum API  
	// - Or scrape from reliable fuel price websites
	
	// Mock data for demonstration
	fuelPrices := map[string]FuelPrice{
		"kerala": {State: "Kerala", Petrol: 102.85, Diesel: 89.15},
		"tamil nadu": {State: "Tamil Nadu", Petrol: 101.50, Diesel: 88.75},
		"karnataka": {State: "Karnataka", Petrol: 102.86, Diesel: 88.94},
		"maharashtra": {State: "Maharashtra", Petrol: 106.31, Diesel: 94.27},
		"delhi": {State: "Delhi", Petrol: 96.72, Diesel: 89.62},
		"gujarat": {State: "Gujarat", Petrol: 96.77, Diesel: 92.91},
		"rajasthan": {State: "Rajasthan", Petrol: 107.49, Diesel: 92.91},
		"uttar pradesh": {State: "Uttar Pradesh", Petrol: 96.57, Diesel: 89.76},
		"west bengal": {State: "West Bengal", Petrol: 106.03, Diesel: 92.76},
		"punjab": {State: "Punjab", Petrol: 108.53, Diesel: 94.61},
	}
	
	stateKey := strings.ToLower(strings.TrimSpace(state))
	if price, exists := fuelPrices[stateKey]; exists {
		return price, nil
	}
	
	// Default to national average if state not found
	return FuelPrice{
		State:  "India (Average)",
		Petrol: 102.50,
		Diesel: 90.25,
	}, nil
}

// CalculateTrip performs the complete trip calculation
func (a *App) CalculateTrip(origin, destination, state, fuelType string, mileage float64) (TripCalculation, error) {
	// Get route information
	route, err := a.GetRoute(origin, destination)
	if err != nil {
		return TripCalculation{}, fmt.Errorf("route calculation failed: %v", err)
	}
	
	// Get fuel prices
	fuelPrice, err := a.GetFuelPrices(state)
	if err != nil {
		return TripCalculation{}, fmt.Errorf("fuel price fetch failed: %v", err)
	}
	
	// Calculate fuel cost
	var pricePerLiter float64
	if strings.ToLower(fuelType) == "petrol" {
		pricePerLiter = fuelPrice.Petrol
	} else {
		pricePerLiter = fuelPrice.Diesel
	}
	
	fuelNeeded := route.DistanceKm / mileage
	totalCost := fuelNeeded * pricePerLiter
	
	return TripCalculation{
		Route:     route,
		FuelPrice: fuelPrice,
		FuelCost:  totalCost,
		FuelType:  fuelType,
		Mileage:   mileage,
	}, nil
}

// GetIndianStates returns list of Indian states
func (a *App) GetIndianStates() []string {
	return []string{
		"Andhra Pradesh", "Arunachal Pradesh", "Assam", "Bihar", "Chhattisgarh",
		"Goa", "Gujarat", "Haryana", "Himachal Pradesh", "Jharkhand", "Karnataka",
		"Kerala", "Madhya Pradesh", "Maharashtra", "Manipur", "Meghalaya", "Mizoram",
		"Nagaland", "Odisha", "Punjab", "Rajasthan", "Sikkim", "Tamil Nadu",
		"Telangana", "Tripura", "Uttar Pradesh", "Uttarakhand", "West Bengal",
		"Delhi", "Jammu and Kashmir", "Ladakh", "Puducherry", "Chandigarh",
		"Dadra and Nagar Haveli and Daman and Diu", "Lakshadweep",
		"Andaman and Nicobar Islands",
	}
}

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "Fuel Cost Calculator",
		Width:  1200,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup:        app.OnStartup,
		Colour:           options.RGBA{R: 27, G: 38, B: 54, A: 1},
		MinWidth:         800,
		MinHeight:        600,
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

// wails.json configuration file
/*
{
  "name": "fuel-calculator",
  "version": "1.0.0",
  "description": "Fuel Cost Calculator for Indian Roads",
  "author": {
    "name": "Your Name",
    "email": "your.email@example.com"
  },
  "wails": {
    "build": {
      "frontend": {
        "dir": "./frontend",
        "install": "npm install",
        "build": "npm run build"
      }
    }
  }
}
*/

// go.mod
/*
module fuel-calculator

go 1.21

require github.com/wailsapp/wails/v2 v2.6.0

require (
	github.com/bep/debounce v1.2.1 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/jchv/go-winloader v0.0.0-20210711035445-715c2860da7e // indirect
	github.com/labstack/echo/v4 v4.10.2 // indirect
	github.com/labstack/gommon v0.4.0 // indirect
	github.com/leaanthony/go-ansi-parser v1.6.0 // indirect
	github.com/leaanthony/gosod v1.0.3 // indirect
	github.com/leaanthony/slicer v1.6.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/pkg/browser v0.0.0-20210911075715-681adbf594b8 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rivo/uniseg v0.4.4 // indirect
	github.com/samber/lo v1.38.1 // indirect
	github.com/tkrajina/go-reflector v0.5.6 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/wailsapp/go-webview2 v1.0.1 // indirect
	github.com/wailsapp/mimetype v1.4.1 // indirect
	golang.org/x/crypto v0.6.0 // indirect
	golang.org/x/exp v0.0.0-20230522175609-2e198f4a06a1 // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/sys v0.7.0 // indirect
	golang.org/x/text v0.7.0 // indirect
)
*/
package main

import (
	"context"
	"fmt"
	"math"
)

// App struct
type App struct {
	ctx context.Context
}

// FuelCalculation represents the result of a fuel cost calculation
type FuelCalculation struct {
	Distance        float64 `json:"distance"`
	FuelEfficiency  float64 `json:"fuelEfficiency"`
	FuelPrice       float64 `json:"fuelPrice"`
	FuelNeeded      float64 `json:"fuelNeeded"`
	TotalCost       float64 `json:"totalCost"`
	CostPerKm       float64 `json:"costPerKm"`
}

// TripComparison represents a comparison between different vehicles/routes
type TripComparison struct {
	Vehicle1 FuelCalculation `json:"vehicle1"`
	Vehicle2 FuelCalculation `json:"vehicle2"`
	Savings  float64         `json:"savings"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// OnStartup is called when the app starts, before the frontend has loaded
func (a *App) OnStartup(ctx context.Context) {
	a.ctx = ctx
	fmt.Println("Fuel Calculator App Started!")
}

// Greetings returns a greeting for the given name
func (a *App) Greetings(name string) string {
	return fmt.Sprintf("Hello %s, Welcome to Fuel Calculator!", name)
}

// CalculateFuelCost calculates the fuel cost for a trip
func (a *App) CalculateFuelCost(distance, efficiency, pricePerLiter float64) FuelCalculation {
	// Input validation
	if distance <= 0 || efficiency <= 0 || pricePerLiter <= 0 {
		return FuelCalculation{}
	}

	// Calculate fuel needed (in liters)
	fuelNeeded := distance / efficiency

	// Calculate total cost
	totalCost := fuelNeeded * pricePerLiter

	// Calculate cost per kilometer
	costPerKm := totalCost / distance

	return FuelCalculation{
		Distance:       distance,
		FuelEfficiency: efficiency,
		FuelPrice:      pricePerLiter,
		FuelNeeded:     math.Round(fuelNeeded*100) / 100,     // Round to 2 decimal places
		TotalCost:      math.Round(totalCost*100) / 100,      // Round to 2 decimal places
		CostPerKm:      math.Round(costPerKm*100) / 100,      // Round to 2 decimal places
	}
}

// CalculateRoundTrip calculates fuel cost for a round trip
func (a *App) CalculateRoundTrip(distance, efficiency, pricePerLiter float64) FuelCalculation {
	return a.CalculateFuelCost(distance*2, efficiency, pricePerLiter)
}

// CompareFuelCosts compares fuel costs between two different scenarios
func (a *App) CompareFuelCosts(
	distance1, efficiency1, price1 float64,
	distance2, efficiency2, price2 float64,
) TripComparison {
	
	vehicle1 := a.CalculateFuelCost(distance1, efficiency1, price1)
	vehicle2 := a.CalculateFuelCost(distance2, efficiency2, price2)
	
	savings := math.Abs(vehicle1.TotalCost - vehicle2.TotalCost)
	savings = math.Round(savings*100) / 100

	return TripComparison{
		Vehicle1: vehicle1,
		Vehicle2: vehicle2,
		Savings:  savings,
	}
}

// CalculateFuelNeeded calculates how much fuel is needed for a given distance
func (a *App) CalculateFuelNeeded(distance, efficiency float64) float64 {
	if distance <= 0 || efficiency <= 0 {
		return 0
	}
	
	fuelNeeded := distance / efficiency
	return math.Round(fuelNeeded*100) / 100
}

// CalculateMaxDistance calculates maximum distance possible with given fuel amount
func (a *App) CalculateMaxDistance(fuelAmount, efficiency float64) float64 {
	if fuelAmount <= 0 || efficiency <= 0 {
		return 0
	}
	
	maxDistance := fuelAmount * efficiency
	return math.Round(maxDistance*100) / 100
}

// ConvertMPGToKmpl converts Miles Per Gallon to Kilometers Per Liter
func (a *App) ConvertMPGToKmpl(mpg float64) float64 {
	if mpg <= 0 {
		return 0
	}
	
	// 1 MPG = 0.425144 km/l
	kmpl := mpg * 0.425144
	return math.Round(kmpl*100) / 100
}

// ConvertKmplToMPG converts Kilometers Per Liter to Miles Per Gallon
func (a *App) ConvertKmplToMPG(kmpl float64) float64 {
	if kmpl <= 0 {
		return 0
	}
	
	// 1 km/l = 2.35214 MPG
	mpg := kmpl * 2.35214
	return math.Round(mpg*100) / 100
}

// GetFuelEfficiencyCategory returns a category description based on fuel efficiency
func (a *App) GetFuelEfficiencyCategory(efficiency float64) string {
	if efficiency <= 0 {
		return "Invalid efficiency"
	}
	
	if efficiency >= 20 {
		return "Excellent (20+ km/l)"
	} else if efficiency >= 15 {
		return "Good (15-20 km/l)"
	} else if efficiency >= 10 {
		return "Average (10-15 km/l)"
	} else if efficiency >= 5 {
		return "Poor (5-10 km/l)"
	} else {
		return "Very Poor (< 5 km/l)"
	}
}

// CalculateAnnualFuelCost calculates annual fuel cost based on monthly usage
func (a *App) CalculateAnnualFuelCost(monthlyDistance, efficiency, pricePerLiter float64) map[string]float64 {
	if monthlyDistance <= 0 || efficiency <= 0 || pricePerLiter <= 0 {
		return map[string]float64{
			"monthlyDistance": 0,
			"monthlyCost":     0,
			"annualDistance":  0,
			"annualCost":      0,
		}
	}

	monthlyCost := a.CalculateFuelCost(monthlyDistance, efficiency, pricePerLiter).TotalCost
	annualDistance := monthlyDistance * 12
	annualCost := monthlyCost * 12

	return map[string]float64{
		"monthlyDistance": monthlyDistance,
		"monthlyCost":     math.Round(monthlyCost*100) / 100,
		"annualDistance":  annualDistance,
		"annualCost":      math.Round(annualCost*100) / 100,
	}
}
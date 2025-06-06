package main

import (
	"context"
	"fmt"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts, before the frontend has loaded
func (a *App) OnStartup(ctx context.Context) {
	a.ctx = ctx
	fmt.Println("Fuel Calculator App Started!")
}

// Greet returns a greeting for the given name
func (a *App) Greetings(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
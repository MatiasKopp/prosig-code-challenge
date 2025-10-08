package main

import "github.com/MatiasKopp/prosig-code-challenge/internal/app"

func main() {
	api := app.New()
	api.Start()
}

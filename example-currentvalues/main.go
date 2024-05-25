package main

import (
	"fmt"
	"os"

	aseko "github.com/danielswe88/aseko-pool-live-api"
)

func main() {
	// SET UP CLIENT
	fmt.Println("Setting up client...")
	asekoClient, err := aseko.NewClient(
		os.Getenv("username"),
		os.Getenv("password"),
		os.Getenv("deviceid"),
		true,
	)
	if err != nil {
		panic(err)
	}

	// GETTING CURRENT VALUES
	fmt.Println("Getting current values...")
	currentValues, err := asekoClient.CurrentValues()
	if err != nil {
		panic(err)
	}

	fmt.Println("Current values:")
	fmt.Printf("Current errors: %+v\n", currentValues.Errors)
	for _, v := range currentValues.Variables {
		fmt.Println(v.Name, v.CurrentValue)
	}
	fmt.Println("")

}

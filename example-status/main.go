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

	// GET GENERAL STATUS OF SYSTEM
	fmt.Println("Getting status...")
	status, err := asekoClient.Status()
	if err != nil {
		panic(err)
	}

	fmt.Println("Status:")
	fmt.Println("SerialNumber: ", status.SerialNumber)
	fmt.Println("Type: ", status.Type)
	fmt.Println("Name: ", status.Name)
	fmt.Println("Timezone: ", status.Timezone)
	fmt.Println("IsOnline: ", status.IsOnline)
	fmt.Println("DateLastData: ", status.DateLastData)
	fmt.Println("HasError: ", status.HasError)
}

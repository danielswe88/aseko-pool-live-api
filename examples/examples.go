package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	aseko "github.com/danielswe88/aseko-pool-live-api"
)

func main() {
	// SET UP CLIENT
	fmt.Println("Setting up client...")
	debugModeOn := false
	asekoClient := aseko.NewClient(debugModeOn)

	fmt.Println("Logging in...")

	deviceID, err := strconv.Atoi(os.Getenv("deviceid"))
	if err != nil {
		fmt.Println("Could not convert device id to integer")
		panic(err)
	}

	// LOGIN
	err = asekoClient.Login(
		os.Getenv("username"),
		os.Getenv("password"),
		int32(deviceID),
	)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Successfully logged in")
	}

	// GET GENERAL STATUS OF SYSTEM
	fmt.Println("Getting status...")
	status, err := asekoClient.Status()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Status: %+v\n\n", status)

	// GETTING CURRENT VALUES
	fmt.Println("Getting current values...")
	currentValues, err := asekoClient.CurrentValues()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Current values: %+v\n\n", currentValues)

	// GETTING GRAPH DATA
	fmt.Println("Getting graph data...")
	startDate, _ := time.Parse(time.RFC3339, "2022-08-10T10:00:00Z")
	endDate := startDate.Add(24 * time.Hour)
	graphData, err := asekoClient.GetGraphData(startDate, endDate)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Graph Data: %+v\n\n", graphData)
}

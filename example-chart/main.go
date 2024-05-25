package main

import (
	"fmt"
	"os"
	"time"

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

	// GETTING GRAPH DATA
	fmt.Println("Getting graph data...")
	endDate := time.Now()
	startDate := endDate.Add(-15 * time.Minute)
	graphData, err := asekoClient.Chart(startDate, endDate)
	if err != nil {
		panic(err)
	}

	fmt.Println("Graph Data:")
	for _, item := range graphData.Items {
		fmt.Printf("Timestamp: %+v. ", item.Timestamp)
		fmt.Printf("PH: %v. ", item.PH)
		fmt.Printf("Redox: %v. ", item.Redox)
		fmt.Printf("Temperature: %v", item.WaterTemp)
		fmt.Println("")
	}
}

package main

import (
	"log"
	"os"

	"github.com/blewb/bubblebeam/stream"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	apiURL := os.Getenv("API_URL")
	apiToken := os.Getenv("API_TOKEN")

	if apiURL == "" || apiToken == "" {
		log.Fatal("Missing required API URL or token")
	}

	api := stream.NewAPI(apiURL, apiToken)

	api.GetBranches()

}

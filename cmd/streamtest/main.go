package main

import (
	"log"
	"os"
	"strconv"
	"fmt"

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
		log.Fatal("missing required API URL or token")
	}

	apiUser := os.Getenv("USER_ID")
	userid, err := strconv.ParseInt(apiUser, 10, 64)
	if err != nil {
		log.Fatal("invalid user ID")
	}

	api := stream.NewAPI(apiURL, apiToken, userid)
	err = api.LoadJobs()

	fmt.Println(err)

	// res, err := api.GetJobItems(12, 34)

}

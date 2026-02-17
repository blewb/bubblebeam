package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/blewb/bubblebeam/span"
	"github.com/blewb/bubblebeam/stream"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
)

var flagSample int

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

	flag.Parse()

	if flagSample <= 0 {
		flagSample = 1
	}

	sp := span.NewSpan()
	if err := sp.Read(flagSample); err != nil {
		log.Fatal(err)
	}

	api := stream.NewAPI(apiURL, apiToken, userid)
	if err := api.LoadJobs(); err != nil {
		log.Fatal(err)
	}

	p := tea.NewProgram(initialModel(sp, api), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

}

func init() {
	flag.IntVar(&flagSample, "sample", 0, "Use sample data file #")
	flag.IntVar(&flagSample, "s", 0, "Use sample data file #")
}

package main

import (
	_ "embed"
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

const (
	SELECTION_DAYS = 21
)

//go:embed title.txt
var appTitle string
var flagSample, flagState int

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
		flagSample = 1 // Development default
		// log.Fatalf("invalid sample: %d", flagSample)
	}

	if flagState < 0 || flagState > int(StateConfirm) {
		flagState = 0 // Development default
		// log.Fatalf("invalid sample: %d", flagSample)
	}

	sp := span.NewSpan()
	sp.Read(flagSample)

	api := stream.NewAPI(apiURL, apiToken, userid)
	// api.LoadJobs()

	p := tea.NewProgram(initialModel(sp, api, modelState(flagState)), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

}

func init() {

	flag.IntVar(&flagSample, "sample", 0, "Use sample data file #")
	flag.IntVar(&flagSample, "s", 0, "Use sample data file #")

	flag.IntVar(&flagState, "state", 0, "Launch into a given app state")
	flag.IntVar(&flagState, "t", 0, "Launch into a given app state")

}

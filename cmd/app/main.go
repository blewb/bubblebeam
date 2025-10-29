package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"

	"github.com/blewb/bubblebeam/span"
	tea "github.com/charmbracelet/bubbletea"
)

//go:embed title.txt
var appTitle string
var flagSample, flagState int

func main() {

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

	p := tea.NewProgram(initialModel(sp, modelState(flagState)), tea.WithAltScreen())
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

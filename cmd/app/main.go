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
var flagSample int

func main() {

	flag.Parse()

	if flagSample <= 0 {
		flagSample = 1 // Development default
		// log.Fatalf("invalid sample: %d", flagSample)
	}

	sp := span.NewSpan()
	sp.Read(flagSample)

	p := tea.NewProgram(initialModel(sp), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

}

func init() {

	flag.IntVar(&flagSample, "sample", 0, "Use sample data file #")
	flag.IntVar(&flagSample, "s", 0, "Use sample data file #")

}

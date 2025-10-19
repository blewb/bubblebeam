package main

// These imports will be used later on the tutorial. If you save the file
// now, Go might complain they are unused, but that's fine.
// You may also need to run `go mod tidy` to download bubbletea and its
// dependencies.
import (
	"flag"
	"fmt"
	"os"

	"github.com/blewb/bubblebeam/span"
	tea "github.com/charmbracelet/bubbletea"
)

var flagSample int

func main() {

	flag.Parse()

	if flagSample <= 0 {
		flagSample = 3 // Development default
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

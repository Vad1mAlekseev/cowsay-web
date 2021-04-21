package cowsay

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

var errorRequestingFigs = errors.New("error requesting figures")

type Cowsay struct {
	listCache []string
}

// List returns supported figures.
func (s *Cowsay) List() ([]string, error) {
	if s.listCache == nil {
		figures, err := allFigures()
		if err != nil {
			return nil, err
		}

		s.listCache = figures
	}

	return s.listCache, nil
}

// Make a new figure by name.
func (s *Cowsay) Make(figName string, text string) ([]byte, error) {
	if text == "random" {
		fortCmd := exec.Command("fortune")

		fort, err := fortCmd.Output()
		if err != nil {
			return nil, fmt.Errorf("error running fortune: %w", err)
		}

		text = string(fort)
	}

	cmd := exec.Command("cowsay", "-f", figName, text)

	figure, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error creating the figure %s: %w", figName, err)
	}

	return figure, nil
}

func allFigures() ([]string, error) {
	figure, err := exec.Command("cowsay", "-l").Output()
	if err != nil {
		return nil, fmt.Errorf("error while check available figures: %w", err)
	}

	output := string(figure)
	lines := strings.Split(output, "\n")
	// Skip help info
	lines = lines[1:]
	output = strings.Join(lines, " ")
	figures := strings.Fields(output)

	if len(figures) == 0 {
		return nil, errorRequestingFigs
	}

	return figures, nil
}

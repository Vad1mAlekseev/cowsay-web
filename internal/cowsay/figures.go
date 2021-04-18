package cowsay

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

type Cowsay struct {
	listCache []string
}

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

func (s *Cowsay) Make(figName string, text string) ([]byte, error) {
	if text == "random" {
		fortCmd := exec.Command("fortune")

		fort, err := fortCmd.Output()
		if err != nil {
			return nil, err
		}
		text = string(fort)
	}

	cmd := exec.Command("cowsay", "-f", figName, text)

	figure, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return figure, nil
}

func allFigures() ([]string, error) {
	figure, err := exec.Command("cowsay", "-l").Output()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error while check available figures: %v\n", err))
	}

	output := string(figure)
	lines := strings.Split(output, "\n")
	// Skip help info
	lines = lines[1:]
	output = strings.Join(lines, " ")
	figures := strings.Fields(output)
	if len(figures) <= 0 {
		return nil, errors.New("error requesting figures")
	}

	return figures, nil
}

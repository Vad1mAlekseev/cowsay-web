package figures

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

var availableCache []string

func Available() ([]string, error) {
	if availableCache == nil {
		figures, err := getInstalledFigures()
		if err != nil {
			return nil, err
		}
		availableCache = figures
	}

	return availableCache, nil
}

func Create(phrase string) {
	return
}

func callCowsay(flags []string, text string) (string, error) {
	args := make([]string, 0, len(flags) + 1)
	for _, flag := range flags {
		args = append(args, "-" + flag)
	}

	if text == "fortune" {
		fortCmd := exec.Command("fortune")

		fort, err := fortCmd.Output()
		if err != nil {
			return "", err
		}
		args[len(args) - 1] = string(fort)
	}

	figure, err := exec.Command("cowsay", args...).Output()
	if err != nil {
		return "", err
	}

	return string(figure), nil
}

func getInstalledFigures() ([]string, error)  {
	output, err := callCowsay([]string{"l"}, "")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error while check available figures: %v\n", err))
	}

	lines := strings.Split(output, "\n")
	// Skip verbose info
	lines = lines[1:]
	output = strings.Join(lines, " ")
	figures := strings.Fields(output)
	if len(figures) <= 0 {
		return nil, errors.New("error while requesting figures")
	}

	return figures, nil
}

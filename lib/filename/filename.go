package filename

import (
	"errors"
	"strings"
)

func GetFromArgs(args []string) (filename string, err error) {
	if len(args) > 0 {
		filename = strings.TrimSpace(args[0])
	}

	if filename == "" {
		return "", errors.New("Please provide a valid filename")
	}

	return filename, nil
}

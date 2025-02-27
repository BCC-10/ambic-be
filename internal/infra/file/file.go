package file

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ReadHTML(path, filename string) (string, error) {
	fullPath := fmt.Sprintf("%s/%s.html", path, filename)
	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		return "", err
	}

	file, err := os.Open(absPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var builder strings.Builder
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		builder.WriteString(scanner.Text() + "\n")
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return builder.String(), nil
}

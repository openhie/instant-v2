package file

import (
	"bufio"
	"os"
	"strings"
)

func Sed(fileName, from, to string) error {
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0777)
	if err != nil {
		return err
	}
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	var commands []string
	for fileScanner.Scan() {
		text := fileScanner.Text()
		trimmedText := strings.TrimSpace(text)

		match := strings.Contains(trimmedText, from)

		if match && to == "" {
			continue
		} else if match {
			commands = append(commands, strings.ReplaceAll(text, from, to))
			continue
		}

		commands = append(commands, text)
	}
	file.Close()

	file, err = os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(file)

	for _, command := range commands {
		b := []byte(command)

		_, err = writer.Write(append(b, '\n'))
		if err != nil {
			return err
		}
	}

	err = writer.Flush()
	if err != nil {
		return err
	}

	return file.Close()
}

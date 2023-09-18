package main

import (
	"bufio"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var path, newText, textToReplace string
	_, err := fmt.Fscan(in, &path, &textToReplace, &newText)
	if err != nil {
		fmt.Printf("Reading data error %s", err.Error())
		return
	}
	if len(path) == 0 {
		log.Fatal("Path is not defined")
		return
	} else if len(newText) == 0 {
		log.Fatal("New text is nil")
		return
	} else if len(textToReplace) == 0 {
		log.Fatal("Text to replace is nil")
		return
	}
	e := run(path, newText, textToReplace)
	if e != nil {
		log.Fatal(e)
	}
	slog.Info("All done")
}

func run(path string, newText string, textToReplace string) error {
	logFile, err := os.CreateTemp(path, "log.*.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := logFile.Close()
		slog.Warn(err.Error())
	}()
	logger := log.New(logFile, "", log.LstdFlags)
	files, e := os.ReadDir(path)
	if e != nil {
		_, errLog := logFile.WriteString("Something went wrong, check path")
		return errLog
	}
	for _, file := range files {
		if !file.IsDir() {
			pathStr := []string{path, file.Name()}
			name, errorOpFile := os.OpenFile(strings.Join(pathStr, "\\"), os.O_RDWR, 0644)
			if errorOpFile != nil {
				logger.Println("File reading error, path: %s, file name: %s", path, name)
			}
			scanner := bufio.NewScanner(name)
			var lines []string
			lineNumber := 1
			for scanner.Scan() {
				line := scanner.Text()
				if strings.Contains(line, textToReplace) {
					var newStr string
					ctr := strings.Count(line, textToReplace)
					for i := ctr; i > 0; i-- {
						index := strings.Index(line, textToReplace)
						newStr = strings.Replace(line, textToReplace, newText, 1)
						logger.Printf("file: %s: line %d, pos: %d \n %s -> %s \n", name.Name(), lineNumber, index, getSnippet(newStr, textToReplace, 5), getSnippet(newStr, newText, 5))
					}
					lines = append(lines, newStr)
				}
				lineNumber++
			}
			err = name.Truncate(0)
			_, err = name.Seek(0, 0)
			writer := bufio.NewWriter(name)
			for _, line := range lines {
				writer.WriteString(line + "\n")
			}
			err = writer.Flush()
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	return err
}

func getSnippet(oldText string, newText string, snippetLength int) string {
	index := strings.Index(oldText, newText)
	if index == -1 {
		return ""
	}
	start := index - snippetLength
	if start < 0 {
		start = 0
	}
	end := index + len(newText) + snippetLength
	if end > len(oldText) {
		end = len(oldText)
	}
	return oldText[start:end]
}

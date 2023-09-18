package main

import (
	"bufio"
	"fmt"
	"log"
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
		fmt.Println("Path is not defined")
		return
	} else if len(newText) == 0 {
		fmt.Println("New text is nil")
		return
	} else if len(textToReplace) == 0 {
		fmt.Println("Text to replace is nil")
		return
	}
	e := run(path, newText, textToReplace)
	fmt.Println(e)
}

func run(path string, newText string, textToReplace string) error {
	logFile, err := os.Create("log.txt") //создали лог файл
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	logger := log.New(logFile, "", log.LstdFlags)

	files, e := os.ReadDir(path)
	if e != nil {
		logFile.WriteString("Something went wrong, check path")
		log.Fatal(e)
	}
	for _, file := range files {
		if !file.IsDir() {
			pathStr := []string{path, file.Name()}
			name, e := os.OpenFile(strings.Join(pathStr, "\\"), os.O_RDWR, 0644)
			if e != nil {
				logger.Println("File reading error, path: %s, file name: %s", path, name)
			}
			scanner := bufio.NewScanner(name)
			var lines []string
			lineNumber := 1
			for scanner.Scan() {
				line := scanner.Text()
				if strings.Contains(line, textToReplace) {
					newStr := strings.ReplaceAll(line, textToReplace, newText)
					//делать запись о кажждой замене в тексте в лог файле
					logger.Printf("%s: line %d \n %s -> %s \n", name.Name(), lineNumber, getSnippet(line, textToReplace, 5), getSnippet(newStr, newText, 5))
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
	//snippetLength - длина фрагмента текста перед и после заменяемого текста
	index := strings.Index(oldText, newText)
	if index == -1 {
		return ""
	}
	start := index - snippetLength
	if start < 0 {
		start = 0
	}
	end := index + snippetLength
	if end > len(oldText) {
		end = len(oldText)
	}
	return oldText[start:end]
}

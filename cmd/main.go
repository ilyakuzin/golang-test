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
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	//как лучше будет обработать аргументы?
	var path, newText, textToReplace string
	fmt.Fscan(in, path, newText, textToReplace) //обрабатывать ли ошибку?

	if len(path) == 0 || len(newText) == 0 || len(textToReplace) == 0 {
		fmt.Fprintln(out, "One of args not defined")
		//выход?
	}
	run(path, newText, textToReplace)
}

func run(path string, newText string, textToReplace string) {
	//year, month, day := time.Now().Date() //а как пребразовать дату в строчку? :)
	//чтобы в названии лог-файла указать дату, что-то типо logFile, err := os.Create("log%с%d%c.txt", day, month, year)
	logFile, err := os.Create("log.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	logger := log.New(logFile, "", log.LstdFlags)

	files, err := os.ReadDir(path)
	if err != nil {
		logFile.WriteString("Something went wrong, check path")
		log.Fatal(err)
		//defer files.Close() ?
	}
	for _, file := range files {
		if !file.IsDir() {
			name, err := os.Open(file.Name())
			if err != nil {
				logger.Println("File reading error, path: %s, file name: %c", path, name)
			}
			scanner := bufio.NewScanner(name)
			var lines []string
			lineNumber := 1
			for scanner.Scan() {
				line := scanner.Text()
				if strings.Contains(line, textToReplace) {
					newStr := strings.ReplaceAll(line, textToReplace, newText)
					logger.Println("%s: line %d, %s -> &s", path, lineNumber,
						getSnippet(line, textToReplace, 5),
						getSnippet(newText, newStr, 5))
					lines = append(lines, newStr)
				}
				lineNumber++
			}
			name.Truncate(0) //очищаем файл
			name.Seek(0, 0)  //возвращаемся в начало
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
}

// регуляркой не придумал как, поэтому отдельный метод
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

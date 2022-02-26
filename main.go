package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(output io.Writer, currDir string, printFiles bool) error {
	readAndPrint("", output, currDir, printFiles)
	return nil
}

func readAndPrint(prependingString string, output io.Writer, currDir string, printFiles bool) {
	// Считать строку
	fileObj, err := os.Open(currDir)
	//Закрыть файл
	defer fileObj.Close()
	// если есть ошибки
	if err != nil {
		log.Fatalf("Could not open %s: %s", currDir, err.Error())
	}
	// Считать имя файла
	fileName := fileObj.Name()
	//Считать файл
	files, err := ioutil.ReadDir(fileName)
	// если есть ошибки
	if err != nil {
		log.Fatalf("Could not read dir names in %s: %s", currDir, err.Error())
	}
	// Применить функцию сортировки к считанному файлу
	files = sortFiles(files)
	// создание нового списка
	var newFileList []os.FileInfo = []os.FileInfo{}

	var length int
	if !printFiles {
		for _, file := range files {
			// Проверка на существование директории, если существует, добавить к списку
			if file.IsDir() {
				newFileList = append(newFileList, file)
			}
		}
		files = newFileList
	}
	length = len(files)
	for i, file := range files {
		if file.IsDir() {
			var stringPrepender string
			if length > i+1 {
				fmt.Fprintf(output, prependingString+"├───"+"%s\n", file.Name())
				stringPrepender = prependingString + "│\t"
			} else {
				fmt.Fprintf(output, prependingString+"└───"+"%s\n", file.Name())
				stringPrepender = prependingString + "\t"
			}
			newDir := filepath.Join(currDir, file.Name())
			readAndPrint(stringPrepender, output, newDir, printFiles)
		} else if printFiles {
			if file.Size() > 0 {
				if length > i+1 {
					fmt.Fprintf(output, prependingString+"├───%s (%vb)\n", file.Name(), file.Size())
				} else {
					fmt.Fprintf(output, prependingString+"└───%s (%vb)\n", file.Name(), file.Size())
				}
			} else {
				if length > i+1 {
					fmt.Fprintf(output, prependingString+"├───%s (empty)\n", file.Name())
				} else {
					fmt.Fprintf(output, prependingString+"└───%s (empty)\n", file.Name())
				}
			}
		}
	}
}

//Функция сортировки файлов
func sortFiles(files []os.FileInfo) (sortedFilesArr []os.FileInfo) {
	var filesMap map[string]os.FileInfo = map[string]os.FileInfo{}
	var unSortedFilesNameArr []string = []string{}
	for _, file := range files {
		unSortedFilesNameArr = append(unSortedFilesNameArr, file.Name())
		filesMap[file.Name()] = file
	}
	sort.Strings(unSortedFilesNameArr)
	for _, stringName := range unSortedFilesNameArr {
		sortedFilesArr = append(sortedFilesArr, filesMap[stringName])
	}
	return sortedFilesArr
}

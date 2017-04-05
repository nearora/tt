package main

import (
	"fmt"
	"path/filepath"
	"os"
)

func addFilenameToList(fileList *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, _ error) error {
		if info.IsDir() {
			fmt.Println("Descending in ", path)
		} else {
			*fileList = append(*fileList, path)
		}

		return nil
	}
}

func main() {
	searchPaths := os.Args[1:]

	if len(searchPaths) < 1 {
		fmt.Printf("Usage: %s path1 [path2] [path3] ...\n", os.Args[0])
		os.Exit(1)
	}

	var fileList []string

	for _, s := range searchPaths {
		err := filepath.Walk(s, addFilenameToList(&fileList))
		if err != nil {
			fmt.Println("Error when walking ", s, ": ", err)
		}

	}

	for _, l := range fileList {
		fmt.Println(l)
	}
}


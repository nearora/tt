package main

import (
	"fmt"
	"path/filepath"
	"os"
)

func main() {
	searchPaths := os.Args[1:]

	if len(searchPaths) < 1 {
		fmt.Printf("Usage: %s path1 [path2] [path3] ...\n", os.Args[0])
		os.Exit(1)
	}

	var fileList []string

	for _, s := range searchPaths {
		err :=
			filepath.Walk(s, func(path string, info os.FileInfo, _ error) error {
				if info.IsDir() {
					fmt.Println("Descending in ", path)
				} else {
					fileList = append(fileList, path)
				}

				return nil
			})

		if err != nil {
			fmt.Println("Error when walking ", s, ": ", err)
		}

	}

	for _, l := range fileList {
		fmt.Println(l)
	}
}


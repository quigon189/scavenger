package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func isTextByExt(path string) bool {
	extensions := map[string]bool{
		".txt": true, ".md": true, ".go": true, 
		".js": true, ".html": true, ".css": true, ".json": true,
		".py": true, ".mod": true, ".yaml": true, ".yml": true,
	}
	ext := filepath.Ext(path)
	return extensions[strings.ToLower(ext)]
}

func printTree(path string, prefix string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for i, entry := range entries {
		name := entry.Name()

		if len(name) > 0 && name[0] == '.' {
			continue
		}

		isLast := i == len(entries)-1

		connector := "├── "
		if isLast {
			connector = "└── "
		}

		fmt.Println(prefix + connector + name)

		if entry.IsDir() {
			newPrefix := prefix
			if isLast {
				newPrefix += "    "
			} else {
				newPrefix += "|   "
			}
			printTree(filepath.Join(path, name), newPrefix)
		}
	}

	return nil
}

func procesDir(root string) error {
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && isTextByExt(path) {
			fmt.Println("--------------------------")
			fmt.Println("Файл:", path)

			content, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("Ошибка чтения файла: %v", err)
			}

			fmt.Println(string(content))
		}

		return nil
	})

	return err
}

func main() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Ошибка получения текущего каталога:", err)
		return
	}
	fmt.Println("Текущий каталог:", dir)
	printTree(dir, "")
	procesDir(dir)
}

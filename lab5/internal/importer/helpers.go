package importer

import (
	"fmt"
	"os"
)

func readFile(filePath string) ([]byte, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла %s: %w", filePath, err)
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("файл пуст: %s", filePath)
	}
	return data, nil
}

package exporter

import (
	"archive/zip"
	"bytes"
	"fmt"
	"lab5/internal/domain"
)

// ZipDecorator — декоратор, упаковывающий результат в ZIP-архив (паттерн Декоратор)
type ZipDecorator struct {
	inner ExportWriter
}

func NewZipDecorator(inner ExportWriter) *ZipDecorator {
	return &ZipDecorator{inner: inner}
}

func (d *ZipDecorator) Write(results []domain.DeliveryResult) ([]byte, string, error) {
	data, ext, err := d.inner.Write(results)
	if err != nil {
		return nil, "", err
	}

	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	fileWriter, err := zipWriter.Create("result" + ext)
	if err != nil {
		return nil, "", fmt.Errorf("ошибка создания файла в архиве: %w", err)
	}

	if _, err := fileWriter.Write(data); err != nil {
		return nil, "", fmt.Errorf("ошибка записи в архив: %w", err)
	}

	if err := zipWriter.Close(); err != nil {
		return nil, "", fmt.Errorf("ошибка закрытия архива: %w", err)
	}

	return buf.Bytes(), ".zip", nil
}

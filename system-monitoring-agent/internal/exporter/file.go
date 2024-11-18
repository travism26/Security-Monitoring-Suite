// internal/exporter/file.go
package exporter

import (
	"encoding/json"
	"log"
	"os"
)

type FileExporter struct {
	outputFilePath string
}

func NewFileExporter(outputFilePath string) *FileExporter {
	return &FileExporter{outputFilePath: outputFilePath}
}

func (e *FileExporter) Export(data map[string]interface{}) error {
	file, err := os.OpenFile(e.outputFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = file.WriteString(string(jsonData) + "\n")
	if err != nil {
		return err
	}

	log.Println("Metrics exported successfully.")
	return nil
}

func (e *FileExporter) Close() error {
	return nil
}

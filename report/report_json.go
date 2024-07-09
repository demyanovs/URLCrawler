package report

import (
	"encoding/json"
	"github.com/demyanovs/urlcrawler/parser"
	"os"
)

// JSONReport represents a JSON report.
type JSONReport struct {
	filePath    string
	firstInsert bool
}

// NewJSONReport creates a new JSONReport.
func NewJSONReport(filePath string) *JSONReport {
	return &JSONReport{
		filePath:    filePath,
		firstInsert: true,
	}
}

// SaveBulk saves multiple records to the file.
func (r *JSONReport) SaveBulk(records []parser.PageData) error {
	file, err := os.OpenFile(r.filePath, os.O_CREATE|os.O_RDWR, 0644)
	defer file.Close()

	if err != nil {
		return err
	}

	// Truncate the file if it's the first insert
	if r.firstInsert {
		err = r.truncateFile(file)
		if err != nil {
			return err
		}

		r.firstInsert = false
	}

	decoder := json.NewDecoder(file)

	var data []parser.PageData
	for decoder.More() {
		err = decoder.Decode(&data)
		if err != nil {
			return err
		}
	}

	data = append(data, records...)

	err = r.truncateFile(file)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	err = encoder.Encode(data)
	if err != nil {
		return err
	}

	return nil
}

func (r *JSONReport) truncateFile(file *os.File) error {
	err := file.Truncate(0)
	if err != nil {
		return err
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	return nil
}

package utils

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

var header = []string{"URL", "StatusCode", "Title", "Description", "Keywords"}

// CSVReport represents a CSV report.
type CSVReport struct {
	filePath    string
	firstInsert bool
}

// NewCSVReport creates a new CSVReport.
func NewCSVReport(filePath string) *CSVReport {
	return &CSVReport{
		filePath:    filePath,
		firstInsert: true,
	}
}

// SaveBulk saves multiple records to the file.
func (r *CSVReport) SaveBulk(records []PageData) error {
	if r.firstInsert == true {
		err := r.addHeader()
		if err != nil {
			return err
		}

		r.firstInsert = false
	}

	file, err := os.OpenFile(r.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	defer file.Close()
	if err != nil {
		log.Fatalln("failed to open file", err)
	}

	w := csv.NewWriter(file)
	defer w.Flush()

	var data [][]string
	for _, record := range records {
		row := []string{record.URL, strconv.Itoa(record.StatusCode), record.Title, record.Desc, record.Keywords}
		data = append(data, row)
	}

	err = w.WriteAll(data)
	if err != nil {
		return err
	}

	return nil
}

func (r *CSVReport) addHeader() error {
	if _, err := os.Stat(r.filePath); err == nil {
		err = os.Truncate(r.filePath, 0)
		if err != nil {
			return err
		}
	}

	file, err := os.OpenFile(r.filePath, os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()

	if err != nil {
		return err
	}

	w := csv.NewWriter(file)
	defer w.Flush()

	err = w.Write(header)
	if err != nil {
		return err
	}

	return nil
}

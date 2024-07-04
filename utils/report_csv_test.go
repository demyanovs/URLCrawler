package utils

import (
	"encoding/csv"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var records = PagesData{
	{
		URL:        "https://en.wikipedia.org/wiki/Yuri_Gagarin",
		StatusCode: 200,
		Title:      "Yuri Gagarin - Wikipedia",
		Desc:       "Soviet pilot and cosmonaut who, aboard the first successful crewed spaceflight, became the first human to journey into outer space.",
		Keywords:   "Pilot, cosmonaut",
	},
	{
		URL:        "https:://en.wikipedia.org/wiki/Fyodor_Dostoevsky",
		StatusCode: 200,
		Title:      "Fyodor Dostoevsky - Wikipedia",
		Desc:       "Russian novelist, short story writer, essayist and journalist",
		Keywords:   "Fyodor Dostoevsky, novelist, essayist, journalist",
	},
	{
		URL:        "https://en.wikipedia.org/wiki/Ilya_Repin",
		StatusCode: 200,
		Title:      "Ilya Repin - Wikipedia",
		Desc:       "Russian painter",
		Keywords:   "",
	},
}

func TestSaveBulkCSV_WithHeaderSuccess(t *testing.T) {
	filePath := "result_test.scv"
	reporter := NewCSVReport(filePath)
	reporter.firstInsert = true

	err := reporter.SaveBulk(records)
	require.NoError(t, err)

	defer os.Remove(filePath)

	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	rows, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	require.Equal(t, 4, len(rows))

}

func TestSaveBulkCSV_WithoutHeaderSuccess(t *testing.T) {
	filePath := "result_test.scv"
	reporter := NewCSVReport(filePath)
	reporter.firstInsert = false

	err := reporter.SaveBulk(records)
	require.NoError(t, err)

	defer os.Remove(filePath)

	f, err := os.Open(filePath)
	if err != nil {
		log.Println("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	rows, err := csvReader.ReadAll()
	if err != nil {
		log.Println(err)
	}

	require.Equal(t, 3, len(rows))
}

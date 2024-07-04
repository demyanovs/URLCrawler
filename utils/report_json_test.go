package utils

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var data = PagesData{
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

func TestSaveBulkJSON_WithHeaderSuccess(t *testing.T) {
	filePath := "result_test.json"
	reporter := NewJSONReport(filePath)
	reporter.firstInsert = true

	err := reporter.SaveBulk(records)
	require.NoError(t, err)

	defer os.Remove(filePath)

	var parsedData PagesData
	f, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
	}

	defer f.Close()

	jsonData, err := io.ReadAll(f)
	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal(jsonData, &parsedData)
	if err != nil {
		log.Println(err)
	}
	require.Equal(t, data, parsedData)

}

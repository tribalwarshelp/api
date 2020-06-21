package repository

import (
	"encoding/csv"
	"net/http"
	"time"
)

var client = &http.Client{
	Timeout: 20 * time.Second,
}

func getCSVData(url string) ([][]string, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return csv.NewReader(resp.Body).ReadAll()
}

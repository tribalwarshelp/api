package repository

import (
	"encoding/csv"
	"net/http"
)

func getCSVData(url string) ([][]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return csv.NewReader(resp.Body).ReadAll()
}

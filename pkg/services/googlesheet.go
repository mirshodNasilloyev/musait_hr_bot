package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type SheetData struct {
	Range          string     `json:"range"`
	MajorDimension string     `json:"majorDimension"`
	Values         [][]string `json:"values"`
}

func GetSheetData(url string) ([]map[string]string, error) {
	//HTTP Get Data request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error on the http request: %v\n", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error on read body: %v\n", err)
	}

	var sheetData SheetData
	err = json.Unmarshal(body, &sheetData)
	if err != nil {
		return nil, fmt.Errorf("error on unmarshal: %v\n", err)
	}
	headers := sheetData.Values[0]
	allRows := []map[string]string{}

	for i := 1; i < len(sheetData.Values); i++ {
		rowMap := make(map[string]string)
		row := sheetData.Values[i]
		for j, header := range headers {
			if j < len(row) {
				rowMap[header] = row[j]
			} else {
				rowMap[header] = ""
			}
		}
		allRows = append(allRows, rowMap)
	}
	return allRows, nil
}
func NewURL(spreadsheetID string, sheetName string, row string, column string, apiKey string) string {
	return fmt.Sprintf("https://sheets.googleapis.com/v4/spreadsheets/%s/values/%s!%s:%s?key=%s", spreadsheetID, sheetName, row, column, apiKey)
}

package companies

import (
	"encoding/json"
	"os"
)

type Company struct {
	CompanyName string `json:"companyName"`
	CVRNumber   string `json:"cvrNumber"`
	Title       string `json:"title"`
	Link        string `json:"link"`
}

const COMPANY_DATA_FILE_PATH = "company-data.json"

// Parse reads data from the filePath provided and attempts to return a slice of companies if they exist.
// if none exist, an empty slice of companies is returned instead.
func Parse(filePath string) ([]Company, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return []Company{}, err
	}
	companies := []Company{}
	if err := json.Unmarshal(data, &companies); err != nil {
		return []Company{}, err
	}
	return companies, nil
}

func GetCompanies(filePath string) ([]Company, error) {
	companies, err := Parse(filePath)
	return companies, err
}

package main

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/gocolly/colly/v2"
)

func scrapeAndSaveToCSV(url, outputFile string) {
	// Create a new collector
	c := colly.NewCollector()

	// Slice to store the scraped data
	var data [][]string

	// On every <tr> element which has height attribute
	c.OnHTML("tr[height]", func(e *colly.HTMLElement) {
		companyName := e.ChildText("td:nth-child(1) p.bodycopy-small")
		cvr := e.ChildText("td:nth-child(2) p.bodycopy-small")

		// Append the scraped data to the data slice
		if companyName != "" && cvr != "" {
			data = append(data, []string{companyName, cvr})
		}
	})

	// Start the scraping
	err := c.Visit(url)
	if err != nil {
		log.Fatalf("Failed to visit %s: %v", url, err)
	}

	// Save the data to CSV
	file, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Cannot create file %s: %s", outputFile, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header to CSV
	writer.Write([]string{"company name", "CVR"})

	// Write the data
	for _, row := range data {
		writer.Write(row)
	}
}

func main() {
	scrapeAndSaveToCSV("https://nyidanmark.dk/en-GB/Words%20and%20Concepts%20Front%20Page/SIRI/List%20certified%20companies", "output.csv")
}

// Netherlands: https://ind.nl/en/public-register-recognised-sponsors/public-register-regular-labour-and-highly-skilled-migrants
// UK: https://assets.publishing.service.gov.uk/government/uploads/system/uploads/attachment_data/file/1185202/2023-09-15_-_Worker_and_Temporary_Worker.csv/preview

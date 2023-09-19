package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Docs - https://programmablesearchengine.google.com/controlpanel/overview?cx=27ab32d7e2b9c46bb / https://developers.google.com/custom-search/v1/using_rest

const (
	// Example URL - "https://www.googleapis.com/customsearch/v1?key=INSERT_YOUR_API_KEY&cx=017576662512468239146&q=lectures"
	DefaultEndpoint = "https://www.googleapis.com/customsearch/v1"
)

type Client struct {
	HTTPClient *http.Client
	Endpoint   string
}

// Do performs the *http.Request and decodes the http.Response.Body into v and return the *http.Response. If v is an io.Writer it will copy the body to the writer.
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	req.RequestURI = ""

	if c.HTTPClient == nil {
		c.HTTPClient = http.DefaultClient
	}

	if c.Endpoint == "" {
		c.Endpoint = DefaultEndpoint
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		_, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		return res, fmt.Errorf("http error code: %d", res.StatusCode)
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, res.Body)
		} else {
			decErr := json.NewDecoder(res.Body).Decode(v)
			if decErr == io.EOF {
				decErr = nil // ignore EOF errors caused by empty response body
			}
			if decErr != nil {
				return nil, decErr
			}
		}
	}
	return res, nil
}

type SearchResponse struct {
	Kind  string `json:"kind"`
	Items []Item `json:"items"`
}

type Item struct {
	Kind             string `json:"kind"`
	Title            string `json:"title"`
	HTMLTitle        string `json:"htmlTitle"`
	Link             string `json:"link"`
	DisplayLink      string `json:"displayLink"`
	Snippet          string `json:"snippet"`
	HTMLSnippet      string `json:"htmlSnippet"`
	CacheID          string `json:"cacheId"`
	FormattedURL     string `json:"formattedUrl"`
	HTMLFormattedURL string `json:"htmlFormattedUrl"`
}

func (c *Client) Search(query string) (SearchResponse, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	searchEngineApiKey := os.Getenv("PROGRAMMABLE_SEARCH_ENGINE_API_KEY")
	searchEngineId := os.Getenv("PROGRAMMABLE_SEARCH_ENGINE_ID")
	fmt.Println("query", fmt.Sprintf("%s?key=%s&cx=%s&q=%s", c.Endpoint, searchEngineApiKey, searchEngineId, query))
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?key=%s&cx=%s&q=%s", c.Endpoint, searchEngineApiKey, searchEngineId, url.QueryEscape(query)), nil)
	if err != nil {
		return SearchResponse{}, fmt.Errorf("error creating new request: %w", err)
	}

	var searchResponse SearchResponse
	_, err = c.Do(req, &searchResponse)
	if err != nil {
		return SearchResponse{}, fmt.Errorf("error searching: %w", err)
	}
	return searchResponse, nil
}

func main() {
	fmt.Println("Welcome to the searcher!")

	searchClient := Client{
		Endpoint: DefaultEndpoint,
	}

	// Open the CSV file
	file, err := os.Open("companies.csv")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// Create a CSV reader from the opened file
	r := csv.NewReader(file)

	// Read the CSV records
	records, err := r.ReadAll()
	if err != nil {
		log.Fatalf("Failed to read records: %v", err)
	}

	// Create (or open) the TypeScript file
	tsFile, err := os.Create("companiesList.ts")
	if err != nil {
		panic(err)
	}
	defer tsFile.Close()

	// Write the header to the TypeScript file
	tsFile.WriteString("export const testList: Company[] = [\n")

	// Loop over the records and print the company names
	for i, record := range records {
		if i == 0 { // Assuming the first row is the header and skipping it
			continue
		}
		searchQuery := record[0]
		cvrNumber := record[1]
		time.Sleep(1000 * time.Millisecond)
		res, err := searchClient.Search(searchQuery)
		if err != nil {
			fmt.Println(err)
		}

		firstResult := res.Items[0]
		companyDetails := struct {
			Title   string
			Link    string
			Snippet string
		}{
			Title:   firstResult.Title,
			Link:    firstResult.Link,
			Snippet: firstResult.Snippet,
		}

		// Write each company's data to the TypeScript file
		tsFile.WriteString(fmt.Sprintf("  { companyName: \"%s\", cvrNumber: \"%s\", title: \"%s\", link: \"%s\", snippet: \"%s\" },\n", record[0], cvrNumber, companyDetails.Title, companyDetails.Link, companyDetails.Snippet))
	}
	tsFile.WriteString("]")

}

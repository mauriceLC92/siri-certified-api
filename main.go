package main

import (
	"fmt"
	"log"
	"siri-certified-api/companies"
)

func main() {
	fmt.Println("hey")

	c, err := companies.GetCompanies(companies.COMPANY_DATA_FILE_PATH)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("c: %v\n", c)
}

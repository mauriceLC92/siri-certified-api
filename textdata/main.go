package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Company struct {
	CompanyName string
	CvrNumber   string
}

func main() {
	file, err := os.Open("data.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var companies []Company
	scanner := bufio.NewScanner(file)
	var currentCompany Company

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "companyName:") {
			currentCompany.CompanyName = strings.TrimPrefix(line, "companyName:")
			currentCompany.CompanyName = strings.TrimSpace(currentCompany.CompanyName)
		} else if strings.HasPrefix(line, "cvrNumber:") {
			currentCompany.CvrNumber = strings.TrimPrefix(line, "cvrNumber:")
			currentCompany.CvrNumber = strings.TrimSpace(currentCompany.CvrNumber)
			companies = append(companies, currentCompany)
			currentCompany = Company{} // Reset currentCompany
		}
	}

	if scanner.Err() != nil {
		panic(err)
	}

	// Now let's write to the data.ts file
	tsFile, err := os.Create("data.ts")
	if err != nil {
		panic(err)
	}
	defer tsFile.Close()

	writer := bufio.NewWriter(tsFile)

	// Writing the array of objects
	writer.WriteString("export const companies = [\n")
	for _, company := range companies {
		writer.WriteString(fmt.Sprintf("  { companyName: \"%s\", cvrNumber: \"%s\" },\n", company.CompanyName, company.CvrNumber))
	}
	writer.WriteString("];\n")
	writer.Flush()
}

package companies_test

import (
	"siri-certified-api/companies"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseReadsJSONFileOfCompaniesAndReturnsASliceOfCompanies(t *testing.T) {
	t.Parallel()

	wantedCompanies := []companies.Company{
		{CompanyName: "&TRADITION A/S", CVRNumber: "18169304", Title: "&Tradition", Link: "https://www.andtradition.com/"},
		{CompanyName: "3Shape A/S", CVRNumber: "25553489", Title: "Dental 3D Scanners & Software for CAD/CAM Dentistry | 3shape", Link: "https://www.3shape.com/"},
	}

	got, err := companies.Parse("../testdata/test-company-data.json")
	if err != nil {
		t.Fatal("error parsing file", err)
	}

	if !cmp.Equal(got, wantedCompanies) {
		t.Error(cmp.Diff(wantedCompanies, got))
	}
}

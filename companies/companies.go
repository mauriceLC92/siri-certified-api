package companies

type Company struct {
	CompanyName string `json:"companyName"`
	CVRNumber   string `json:"cvrNumber"`
	Title       string `json:"title"`
	Link        string `json:"link"`
}

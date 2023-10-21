package main

type Job struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Company     string `json:"company"`
	Location    string `json:"location"`
	Posted      string `json:"posted"`
	Link        string `json:"link"`
	Description string `json:"description"`
	CompanyLogo string `json:"company_logo"`
}

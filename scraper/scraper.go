package scraper

// Item stores details of a single product
type Item struct {
	Title       string `json:"title"`
	Size        string `json:"size"`
	UnitPrice   string `json:"unitPrice"`
	Description string `json:"description"`
}

// Result stores details of the scraped products
type Result struct {
	Items []Item  `json:"items"`
	Total float32 `json:"total"`
}

// Scrape function parses provided URL for product links
func Scrape(urls []string) Result {
	return Result{}
}

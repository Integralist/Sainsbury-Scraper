package scraper

import (
	"fmt"
	"regexp"
	"strconv"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

// Item stores details of a single product
type Item struct {
	Title       string `json:"title"`
	Size        string `json:"size"`
	UnitPrice   string `json:"unitPrice"`
	Description string `json:"description"`
}

// Result stores details of the scraped products
type Result struct {
	Items []Item `json:"items"`
	Total string `json:"total"`
}

func (r *Result) calculate() {
	var total float64

	for _, v := range r.Items {
		if s, err := strconv.ParseFloat(v.UnitPrice, 32); err == nil {
			total = total + s
		}
	}

	r.Total = fmt.Sprintf("%.2f", total)
}

var ch chan Item
var wg sync.WaitGroup

// Scrape function parses provided URL for product links
func Scrape(urls []string) Result {
	ch = make(chan Item, len(urls))

	result := Result{}

	for _, url := range urls {
		wg.Add(1)
		go getItem(url)
		result.Items = append(result.Items, <-ch)
	}
	wg.Wait()

	result.calculate()

	return result
}

func getItem(url string) {
	defer wg.Done()

	doc, err := goquery.NewDocument(url)
	if err != nil {
		fmt.Println(
			fmt.Errorf("Unable to create a new document: %s", err.Error()),
		)
	}

	item := Item{}

	item.Title = doc.Find("h1").Text()
	item.UnitPrice = extractPrice(doc.Find(".pricePerUnit").Text())
	item.Size = doc.Find(".productText").Eq(3).Text()
	item.Description = doc.Find(".productText").First().Text()

	ch <- item
}

func extractPrice(text string) string {
	re := regexp.MustCompile("\\d+\\.\\d+")
	match := re.FindString(text)

	if p, err := strconv.ParseFloat(match, 32); err == nil {
		return fmt.Sprintf("%.2f", p)
	}

	return "0.00"
}

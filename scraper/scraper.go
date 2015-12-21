package scraper

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

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
	Items []Item `json:"results"`
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

type extendedDocument struct {
	Size     string
	Document *goquery.Document
}

var ch chan Item

// Scrape function parses provided URL for product links
func Scrape(urls []string) Result {
	ch = make(chan Item, len(urls))

	result := Result{}

	for _, url := range urls {
		go getItem(url)
		result.Items = append(result.Items, <-ch)
	}

	result.calculate()

	return result
}

func extendDocument(url string) (extendedDocument, error) {
	res, err := http.Get(url)
	if err != nil {
		return extendedDocument{}, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return extendedDocument{}, err
	}
	size := strconv.Itoa(len(body)/1000) + "kb"

	// Rewind response body so it can be re-read by goquery
	res.Body = ioutil.NopCloser(bytes.NewReader(body))

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return extendedDocument{}, err
	}

	return extendedDocument{size, doc}, nil
}

var getItem = func(url string) {
	d, err := extendDocument(url)
	if err != nil {
		fmt.Println(
			fmt.Errorf("Unable to create a new document: %s", err.Error()),
		)
	}

	item := Item{}

	item.Title = d.Document.Find("h1").Text()
	item.UnitPrice = extractPrice(d.Document.Find(".pricePerUnit").Text())
	item.Description = strings.TrimSpace(d.Document.Find(".productText").First().Text())
	item.Size = d.Size

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

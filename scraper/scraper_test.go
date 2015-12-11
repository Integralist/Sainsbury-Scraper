package scraper

import "testing"

func TestScrape(t *testing.T) {
	getItem = func(url string) {
		defer wg.Done()

		ch <- Item{
			"FooTitle",
			"FooSize",
			"10.00",
			"FooDescription",
		}
	}

	urls := []string{
		"http://foo.com/",
		"http://bar.com/",
		"http://baz.com/",
	}

	result := Scrape(urls)
	response := result.Total
	expected := "30.00"

	if response != expected {
		t.Errorf("The response:\n '%s'\ndidn't match the expectation:\n '%s'", response, expected)
	}
}

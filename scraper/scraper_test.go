package scraper

import "testing"

func TestScrapeResults(t *testing.T) {
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
	first := result.Items[0]

	if expected := "FooTitle"; first.Title != expected {
		err(first.Title, expected, t)
	}

	if expected := "FooSize"; first.Size != expected {
		err(first.Size, expected, t)
	}

	if expected := "10.00"; first.UnitPrice != expected {
		err(first.UnitPrice, expected, t)
	}

	if expected := "FooDescription"; first.Description != expected {
		err(first.Description, expected, t)
	}

	if expected := "30.00"; result.Total != expected {
		err(result.Total, expected, t)
	}
}

func err(response, expected string, t *testing.T) {
	t.Errorf("The response:\n '%s'\ndidn't match the expectation:\n '%s'", response, expected)
}

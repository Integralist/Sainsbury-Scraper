package scraper

import "testing"

func TestScrapeResults(t *testing.T) {
	getItem = func(url string) {
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

	var suite = []struct {
		response string
		expected string
	}{
		{first.Title, "FooTitle"},
		{first.Size, "FooSize"},
		{first.UnitPrice, "10.00"},
		{first.Description, "FooDescription"},
		{result.Total, "30.00"},
	}

	for _, v := range suite {
		if v.response != v.expected {
			err(v.response, v.expected, t)
		}
	}
}

func err(response, expected string, t *testing.T) {
	t.Errorf("The response:\n '%s'\ndidn't match the expectation:\n '%s'", response, expected)
}

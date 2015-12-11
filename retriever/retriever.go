package retriever

import "github.com/PuerkitoBio/goquery"

// Collection holds all parsed links from supplied URL
type Collection []string

// Retrieve function parses provided URL for product links
func Retrieve(url string) (Collection, error) {
	coll := Collection{}

	doc, err := goquery.NewDocument(url)
	if err != nil {
		return Collection{}, err
	}

	doc.Find(".productInfo a").Each(func(i int, s *goquery.Selection) {
		if v, exists := s.Attr("href"); exists {
			coll = append(coll, v)
		}
	})

	return coll, nil
}

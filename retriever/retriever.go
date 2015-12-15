package retriever

import "github.com/PuerkitoBio/goquery"

// Collection holds all parsed links from supplied URL
type Collection []string

// DocumentBuilder is a type abstraction over our injected dependency
type DocumentBuilder func(url string) (*goquery.Document, error)

// Retrieve function parses provided URL for product links
func Retrieve(url string, newDoc DocumentBuilder) (Collection, error) {
	coll := Collection{}

	doc, err := newDoc(url)
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

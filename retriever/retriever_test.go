package retriever

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

const href = "http://bar.com/"

func fakeNewDocument(url string) (*goquery.Document, error) {
	body := `
		<html>
			<body>
				<div class="productInfo">
					<a href="{}">Bar</a>
				</div>
			</body>
		<html>
	`
	body = strings.Replace(body, "{}", href, 1)

	resp := &http.Response{
		Status:        "200 OK",
		StatusCode:    200,
		Proto:         "HTTP/1.0",
		ProtoMajor:    1,
		ProtoMinor:    0,
		Body:          ioutil.NopCloser(bytes.NewBufferString(body)),
		ContentLength: int64(len(body)),
		Request:       &http.Request{},
	}

	return goquery.NewDocumentFromResponse(resp)
}

func TestRetrieveReturnValue(t *testing.T) {
	url := "http://foo.com/"

	coll, _ := Retrieve(url, fakeNewDocument)

	if response := coll[0]; response != href {
		t.Errorf("The response:\n '%s'\ndidn't match the expectation:\n '%s'", response, href)
	}
}

## How to run?

```bash
ss
```

> Note: there was no requirement at this stage to define any flag options

## How to build?

```bash
go build -o ss
```

You could use [Gox](http://github.com/mitchellh/gox) to more easily build the binary for multiple systems

```bash
gox -osarch="linux/amd64" -osarch="darwin/amd64" -osarch="windows/amd64" -output="ss.{{.OS}}"
```

## How to run the tests?

```bash
go test -v ./...
```

## Architecture

<small><u>Time spent</u>: 20 minutes</small>

&lt;insert diagram&gt;

> This includes considering my solution as well as writing up this README

## Development

<small><u>Time spent</u>: ? minutes</small>

1. Read-Me Driven Development
2. Create CLI structure
3. Define entry command
4. Define 'retriever' package
5. Define 'scraper' package
6. Define 'utilities' package

### Retriever

The retriever should be handed a URL and return an Array of sub page resource URLs, like so:

```json
[
  "http://hiring-tests.s3-website-eu-west-1.amazonaws.com/2015_Developer_Scrape/sainsburys-apricot-ripe---ready-320g.html",
  "http://hiring-tests.s3-website-eu-west-1.amazonaws.com/2015_Developer_Scrape/sainsburys-avocado-xl-pinkerton-loose-300g.html",
  "http://hiring-tests.s3-website-eu-west-1.amazonaws.com/2015_Developer_Scrape/sainsburys-avocado--ripe---ready-x2.html",
  "http://hiring-tests.s3-website-eu-west-1.amazonaws.com/2015_Developer_Scrape/sainsburys-avocados--ripe---ready-x4.html",
  "http://hiring-tests.s3-website-eu-west-1.amazonaws.com/2015_Developer_Scrape/sainsburys-conference-pears--ripe---ready-x4-%28minimum%29.html",
  "http://hiring-tests.s3-website-eu-west-1.amazonaws.com/2015_Developer_Scrape/sainsburys-golden-kiwi--taste-the-difference-x4-685641-p-44.html",
  "http://hiring-tests.s3-website-eu-west-1.amazonaws.com/2015_Developer_Scrape/sainsburys-kiwi-fruit--ripe---ready-x4.html"
]
```

> Note: I use `.productInfo a` as my filter

### Scraper

The scraper should be passed an Array of URLs (see above for example) so it can concurrently request each resource and parse it for the relevant information:

- Resource size
- Product title
- Product unit size
- Product description
- Product size

The scraper should return a Struct with a key of `results` which is assigned an Array of the details collated details and a key of `total` which details the total cost, like so:


```json
{
  "results": [
    {
      "title":"Sainsbury's Avocado, Ripe & Ready x2",
      "size": "90.6kb",
      "unit_price": 1.80,
      "description": "Great to eat now - refrigerate at home 1 of 5 a day 1 avocado counts as 1 of your 5..."
    }, 
    {
      "title":"Sainsbury's Avocado, Ripe & Ready x4",
      "size": "87kb",
      "unit_price": 2.00,
      "description": "Great to eat now - refrigerate at home 1 of 5 a day 1 avocado counts as 1 of your 5..."
    }
  ],
  "total": 3.80
}
```

If the code needs to be made more *reusable*, then we could also look to inject the Array of 'filters' rather than hardcode them. This would allow the package to be reused on different page types.

> Note:
> 
> I use a multitude of filters such as `h1`, `.pricePerUnit` and `productText`. The last selector isn't very flexible though. I would've used a selector such as `nth-child` but it made the code harder to reason about and so I opted against it.
> 
> Also, if the order of the page changes, then a more robust solution would be to inspect the text at each `productDataItemHeader` element and if a match is made we know we can extract the content that follows it

### Utilities

The utilities package is consumed by the Scraper and is expected to provide utility functions to any other packages

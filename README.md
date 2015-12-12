## How to build?

```bash
go build -o ss
```

You could use [Gox](http://github.com/mitchellh/gox) to more easily build the binary for multiple systems

```bash
gox -osarch="linux/amd64" -osarch="darwin/amd64" -osarch="windows/amd64" -output="ss.{{.OS}}"
```

## How to run compiled binary?

```bash
ss
```

> Note: there was no requirement at this stage to define any flag options

## How to run the tests?

```bash
go test -v ./...
```

## Architecture

![Architecture](https://cloud.githubusercontent.com/assets/180050/11756388/72c1d13a-a051-11e5-860c-7a30bf3e3b49.png)

## Dependencies

The main dependency is [goquery](https://github.com/PuerkitoBio/goquery/) which abstracts away a lot of the complexity of having to manually parse HTML content

The other dependency is [codegangsta/cli](https://github.com/codegangsta/cli) which abstracts away a lot of the boilerplate required for creating a console based application

> Note: I'm a big fan of Dave Cheney's [gb](https://getgb.io/) for managing vendored dependencies. Although the BBC prefers to use [Godep](https://godoc.org/github.com/tools/godep). I opted for neither as there were only two dependencies, and so it felt a little overkill for this small project. Once Go 1.6 is released hopefully we'll see an official/native implementation for vendored dependencies

## Development

1. Read-Me Driven Development
2. Create CLI structure
3. Define entry command
4. Define 'retriever' package
5. Define 'scraper' package

### Retriever

The retriever should be handed a URL and return a Slice of sub page resource URLs, like so:

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

The scraper should return a Struct with a field of `Items` which is assigned an Array of collated details and a field of `Total` which details the total cost. Once it's converted to JSON it'll look something like:


```json
{
    "results": [
        {
            "title": "Sainsbury's Apricot Ripe \u0026 Ready x5",
            "size": "39kb",
            "unitPrice": "3.50",
            "description": "Apricots"
        },
        {
            "title": "Sainsbury's Avocado Ripe \u0026 Ready XL Loose 300g",
            "size": "39kb",
            "unitPrice": "1.50",
            "description": "Avocados"
        },
        {
            "title": "Sainsbury's Avocado, Ripe \u0026 Ready x2",
            "size": "44kb",
            "unitPrice": "1.80",
            "description": "Avocados"
        },
        {
            "title": "Sainsbury's Avocados, Ripe \u0026 Ready x4",
            "size": "39kb",
            "unitPrice": "3.20",
            "description": "Avocados"
        },
        {
            "title": "Sainsbury's Conference Pears, Ripe \u0026 Ready x4 (minimum)",
            "size": "39kb",
            "unitPrice": "1.50",
            "description": "Conference"
        },
        {
            "title": "Sainsbury's Golden Kiwi x4",
            "size": "39kb",
            "unitPrice": "1.80",
            "description": "Gold Kiwi"
        },
        {
            "title": "Sainsbury's Kiwi Fruit, Ripe \u0026 Ready x4",
            "size": "39kb",
            "unitPrice": "1.80",
            "description": "Kiwi"
        }
    ],
    "total": "15.10"
}
```

If the code needs to be made more *reusable*, then we could also look to inject the 'filters' rather than hardcode them. This would allow the package to be reused on different page types.

> Note:
> I use a multitude of filters such as `h1`, `.pricePerUnit`, `productText` and `productDataItemHeader`.

## Commit History

For the purposes of this quick test project I was committing straight to master (which in the real-world is a big no-no). At the BBC we have a specific git workflow for how we merge our PRs. Effectively we squash/rebase before cherry picking, while referencing issues/PRs allows us to close them dynamically upon push to master). [I've documented the workflow here](http://www.integralist.co.uk/posts/github-workflow.html))

## Additional comments

- In order to fulfil the requirement for displaying the size of the HTML page being linked to, I needed to build decorator that mimic'ed an internal function from the goquery dependency. This allowed me to adapt an interface that was effectively the same but incorporated the additional requirement of calculating the response body size. This introduced an issue where by the response body needed to be read twice (once by my implementation and once again when passed to the goquery dependency) 
- I wrote tests for the Scraper package and tried (with what time I had left) to write a test for the Retriever package, but due to how the goquery dependency library had been developed it made testing the package quite difficult (most of its method signatures utilised pointers to actual objects so I wasn't always able to swap out a concrete implementation with an interface). On top of that I realised I might have ended up trying to unit test the dependency rather than verifying expected behaviour of the Retriever package itself (which doesn't really do much other than utilise the goquery API).
- I ended up spending a bit too much time trying to produce the price in the JSON object response as a float rather than a string. The issue I was having was with regards to floats rounding off the last zero (e.g. converting the string into a float would result in something like `15.10` being translated into `15.1`) which was misleading output I felt and so after trying quite a few work arounds, I had to settle on implementing it as a string type instead
- Spent a bit of time investigating the Unicode code points being placed into the JSON output instead of the actual rune character being rendered (e.g. the Struct would show `&` but when marshaled into JSON it would be transformed into the code point `\u0026`). It seems that this is expected behaviour according to the Go documentation. If you paste the JSON output into a browser console then you'll find the code point is translated back to the actual rune character

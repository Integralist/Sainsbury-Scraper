package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	r "github.com/integralist/sainsbury-scraper/retriever"
	s "github.com/integralist/sainsbury-scraper/scraper"
)

const url = "http://hiring-tests.s3-website-eu-west-1.amazonaws.com/2015_Developer_Scrape/5_products.html"

var authorList []cli.Author

func author(name, email string) cli.Author {
	return cli.Author{
		Name:  name,
		Email: email,
	}
}

func authors() []cli.Author {
	authorList = append(
		authorList,
		author("Mark McDonnell", "mark.mcdx@gmail.com"),
	)

	return authorList
}

func commandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}

func main() {
	app := cli.NewApp()
	app.Name = "Sainsbury Scraper"
	app.Version = "0.0.1"
	app.Authors = authors()
	app.Usage = "CLI tool for scraping contents from Sainsbury website"

	app.Action = func(c *cli.Context) {
		coll, err := r.Retrieve(url)
		if err != nil {
			fmt.Printf("There was an issue retrieving links from the page: %s", err.Error())
			os.Exit(1)
		}

		b, err := json.MarshalIndent(s.Scrape(coll), "", "    ")
		if err != nil {
			fmt.Printf("There was an issue converting our data into JSON: %s", err.Error())
			os.Exit(1)
		}

		fmt.Println(string(b))
	}

	app.CommandNotFound = commandNotFound
	app.Run(os.Args)
}

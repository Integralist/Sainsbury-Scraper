package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

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
		fmt.Println("do something")
	}

	app.CommandNotFound = commandNotFound
	app.Run(os.Args)
}

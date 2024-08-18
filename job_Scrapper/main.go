package main

import (
	"jobScrapper/scrapper"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo"
)

const fileName string = "jobs.csv"

func main() {
	// Echo instance
	e := echo.New()

	// Routes
	e.GET("/", hello)

	e.POST("/scrape", handlerScrape)
	// Start server
	e.Logger.Fatal(e.Start(":1223"))
}

// Handler
func hello(c echo.Context) error {
	HTMLcode := `
	<!DOCTYPE html>
	<html lang="en">
  	<head>
    	<meta charset="UTF-8" />
    	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
    	<meta http-equiv="X-UA-Compatible" content="ie=edge" />
    	<title>Go Jobs</title>
  	</head>
  	<body>
    	<h1>Go Jobs</h1>
    	<form method="POST" action="/scrape">
      		<input placeholder="what job do you want" name="term" />
      	<button>Search</button>
    	</form>
  	</body>
	</html>
	`
	return c.HTML(http.StatusOK, HTMLcode)
}

func handlerScrape(c echo.Context) error {
	defer os.Remove(fileName)
	term := strings.ToLower(scrapper.CleanString(c.FormValue("term")))
	scrapper.Scrape(term)
	return c.Attachment(fileName, fileName)
}

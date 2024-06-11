package main

import (
	"log"
	"net/http"
)

var baseurl string = "https://www.youtube.com/?app=desktop&hl=ko&gl=KR"

func main() {
	getPages()
}

func getPages() int {
	res, err := http.Get(baseurl)
	if err != nil {
		log.Fatalln(err)
	}
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status:", res.StatusCode)
	}
	return 0
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

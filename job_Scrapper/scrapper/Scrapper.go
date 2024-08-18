package scrapper

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type extractedJob struct {
	link     string
	title    string
	location string
	summary  string
	company  string
}

// Scrape saramin by a term
func Scrape(term string) {
	var baseURL string = "https://www.saramin.co.kr/zf_user/search/recruit?&searchword=" + term + "&recruitCount=50"
	total := getPages(baseURL)

	var jobs []extractedJob
	c := make(chan []extractedJob)
	for i := 1; i <= total; i++ {
		go getPage(i, baseURL, c)
	}

	for i := 1; i <= total; i++ {
		job := <-c
		jobs = append(jobs, job...)
		// same
		// jobs = append(jobs, <- c...)
	}
	writeJobs(jobs)
}

func writeJobs(jobs []extractedJob) {
	file, err := os.Create("jobs.csv")
	checkErr(err)

	w := csv.NewWriter(file)
	defer w.Flush() // must

	headers := []string{"Link", "Title", "Location", "Summary", "Company"}

	Werr := w.Write(headers)
	checkErr(Werr)

	for _, job := range jobs {
		jobSlice := []string{"https://www.saramin.co.kr/zf_user/jobs/relay/view?rec_idx=" + job.link, job.title, job.location, job.summary, job.company}
		jobErr := w.Write(jobSlice)
		checkErr(jobErr)
	}
}

func getPage(page int, baseURL string, mainC chan<- []extractedJob) {
	pageURL := baseURL + "&recruitPage=" + strconv.Itoa(page)
	fmt.Println("Requesting :", pageURL)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkStatusCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	var jobs []extractedJob
	c := make(chan extractedJob)

	cards := doc.Find(".item_recruit")
	cards.Each(func(i int, card *goquery.Selection) {
		go extractJob(card, baseURL, c)
	})

	for i := 0; i < cards.Length(); i++ {
		job := <-c
		jobs = append(jobs, job)
	}
	mainC <- jobs
}

func CleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), "")
}

func extractJob(card *goquery.Selection, baseURL string, c chan<- extractedJob) {
	link, _ := card.Attr("value")
	title := CleanString(card.Find(".job_tit>a").Text())
	location := CleanString(card.Find(".job_condition>span>a").Text())
	summary := CleanString(card.Find(".job_sector").Clone().ChildrenFiltered(".job_day").Remove().End().Text())
	company := CleanString(card.Find(".area_corp>strong>a").Text())
	c <- extractedJob{
		link:     link,
		title:    title,
		location: location,
		summary:  summary,
		company:  company,
	}
}

func getPages(baseURL string) int {
	pages := 0
	res, err := http.Get(baseURL)
	checkErr(err)
	checkStatusCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length()
	})

	return pages
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkStatusCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status:", res.StatusCode, res.Status)
	}
}

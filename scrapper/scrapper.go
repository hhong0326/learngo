package scrapper

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type requestResult struct {
	url    string
	status string
}

type extractedJob struct {
	id       string
	title    string
	location string
	salary   string
	summary  string
}

// Scrape
func Scrape(term string) {

	var baseURL string = "https://kr.indeed.com/jobs?q=" + term
	var jobs []extractedJob
	c := make(chan []extractedJob)
	tp := getPages(baseURL)
	fmt.Println(tp)

	for i := 0; i < tp; i++ {
		go getPage(i, baseURL, c)
	}

	for i := 0; i < 5; i++ {
		extractedJobs := <-c
		jobs = append(jobs, extractedJobs...) // [] + [] + [] = []
	}

	writeJobs(jobs)
	fmt.Println("Done, ", len(jobs))
}

func getPage(page int, url string, mainC chan<- []extractedJob) {
	var jobs []extractedJob
	c := make(chan extractedJob)
	pageURL := url + "&start=" + strconv.Itoa(page*10)
	fmt.Println(pageURL)

	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCard := doc.Find(".jobsearch-SerpJobCard")

	searchCard.Each(func(i int, card *goquery.Selection) {
		go extractJob(card, c)
		// jobs = append(jobs, job)
	})

	for i := 0; i < searchCard.Length(); i++ {
		job := <-c
		jobs = append(jobs, job)
	}

	mainC <- jobs

}

func writeJobs(jobs []extractedJob) {

	c := make(chan []string)
	file, err := os.Create("jobs.csv")
	checkErr(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"Link", "Title", "Location", "Salary", "Summary"}

	wErr := w.Write(headers)
	checkErr(wErr)

	for _, job := range jobs {
		jobSlice := []string{"https://kr.indeed.com/viewjob?jk=" + job.id, job.title, job.location, job.salary, job.summary}
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)
	}

	for _, job := range jobs {
		go writeToCsv(job, c)
	}

	jobSlice := <-c

	jwErr := w.Write(jobSlice)
	checkErr(jwErr)
}

func writeToCsv(job extractedJob, wC chan<- []string) {
	wC <- []string{"https://kr.indeed.com/viewjob?jk=" + job.id, job.title, job.location, job.salary, job.summary}
}

func extractJob(card *goquery.Selection, c chan<- extractedJob) {
	id, _ := card.Attr("data-jk") // s == div
	title := CleanString(card.Find("a").Text())
	location := CleanString(card.Find(".sjcl>span").Text())
	summary := CleanString(card.Find(".summary").Text())
	// fmt.Println(id)

	c <- extractedJob{
		id:       id,
		title:    title,
		location: location,
		salary:   "4000",
		summary:  summary,
	}
}

// CleanString
func CleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

func getPages(url string) int {
	var pages int = 0
	res, err := http.Get(url)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find(".pagination .pagination-list").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length()
		fmt.Println(s.Find("li").Text())

	})
	return pages
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed")
	}
}

var errRequestFailed = errors.New("Request Failed")

func hitURL(url string, c chan<- requestResult) {

	status := "OK"
	res, err := http.Get(url)
	if err != nil || res.StatusCode >= 400 {
		status = "FAILED"
	}

	c <- requestResult{url: url, status: status}

}

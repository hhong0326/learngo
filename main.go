package main

import (
	"os"
	"strings"

	"github.com/labstack/echo"
	"github.com/learngo/algorithms"
	"github.com/learngo/scrapper"
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

var fileName string = "jobs.csv"

func main() {

	// account := account.NewAccount("sunil")
	// account.Deposit(100)

	// err := account.Withdraw(1000)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Println(account.Balance(), account.Owner())

	// account.ChangeOwner("hong")
	// fmt.Println(account)

	// dictionary := mydict.Dictionary{"first": "First Word"}

	// definition, err := dictionary.Search("first")
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(definition)
	// }

	// word := "second"
	// err = dictionary.Add(word, "ff Word")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// err = dictionary.Add(word, "sd Word")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// definition, err = dictionary.Search(word)
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(definition)
	// }

	// err = dictionary.Update(word, "Second Word")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// word, _ = dictionary.Search(word)
	// fmt.Println(word)

	// err = dictionary.Delete("dfd")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// results := make(map[string]string)
	// c := make(chan requestResult)
	// urls := []string{
	// 	"https://www.airbnb.com/",
	// 	"https://www.google.com/",
	// 	"https://www.amazon.com/",
	// 	"https://www.google.com/",
	// 	"https://soundcloud.com/",
	// 	"https://www.naver.com/",
	// 	"https://www.instagram.com/",
	// }

	// for _, url := range urls {
	// 	go hitURL(url, c)

	// }

	// for i := 0; i < len(urls); i++ {
	// 	result := <-c
	// 	results[result.url] = result.status
	// }

	// for url, status := range results {
	// 	fmt.Println(url, status)
	// }

	// scrapper
	// e := echo.New()

	// e.GET("/", handleHome)
	// e.POST("/scrape", handleScrape)

	// e.Logger.Fatal(e.Start(":1323"))

	//algorithms
	// algorithms.StartQuick()
	// algorithms.StartBST()
	// algorithms.StartSelectSort()
	// algorithms.StartBFS()
	// algorithms.StartDFS()
	// algorithms.StartLN()
	// algorithms.Queens(0)
	algorithms.StartTN()
}

func handleHome(c echo.Context) error {
	return c.File("home.html")
}

func handleScrape(c echo.Context) error {

	defer os.Remove(fileName)

	t := strings.ToLower(scrapper.CleanString(c.FormValue("term")))
	scrapper.Scrape(t)
	// Attachment = 첨부파일 리턴
	return c.Attachment(fileName, fileName) // a를 b 이름으로 다운로드
}

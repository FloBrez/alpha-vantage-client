package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	var ticker string
	var dataset string
	var apikey string
	var tickerfile string
	var outputDir string

	flag.StringVar(&ticker, "t", "IBM", "Specify ticker symbol. Default is IBM.")
	flag.StringVar(&tickerfile, "f", "", "Specify ticker symbol. Default is an empty string.")
	flag.StringVar(&outputDir, "o", "", "Specify output directory.")
	flag.StringVar(&dataset, "d", "TIME_SERIES_DAILY", "Specify dataset type.")
	flag.StringVar(&apikey, "k", "demo", "Specify apikey")
	flag.Parse() // after declaring flags we need to call it

	var tickerList []string
	var err error
	if tickerfile != "" {
		tickerList, err = readTickerFile(tickerfile)
		if err != nil {
			fmt.Println("Error reading ticker file.")
			tickerList = []string{ticker}
		}
	} else {
		tickerList = []string{ticker}
	}

	fmt.Println("Will attempt to retrieve data for ticker symbols", strings.Join(tickerList, ", "))

	fmt.Println("Requesting data from Alpha Vantage...")

	for _, symbol := range tickerList {
		fmt.Println("for", symbol)
		res := request(dataset, symbol, apikey)
		write(fmt.Sprintf("%s%s.json", outputDir, symbol), res)
		time.Sleep(15 * time.Second) // free API limited to 5 calls per minute and 500 per day
	}

}

func readTickerFile(filename string) ([]string, error) {
	contents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("File reading error", err)
		return []string{}, err
	}
	v := strings.Split(string(contents), "\n")
	return v, nil
}

func request(dataset string, symbol string, apikey string) string {

	url := fmt.Sprintf("https://www.alphavantage.co/query?function=%s&symbol=%s&outputsize=full&apikey=%s&datatype=csv", dataset, symbol, apikey)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return "error"
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "error"
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "error"
	}
	//fmt.Println(string(body))
	return string(body)
}

func write(filename string, content string) {
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	l, err := f.WriteString(content)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	fmt.Println(l, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

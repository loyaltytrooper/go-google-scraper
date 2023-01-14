package utils

import (
	"fmt"
	"go-google-scraper/models"
	"net/http"
	"strings"
	"time"
)

var googleDomains = map[string]string{
	"com": "https://www.google.com/search?q=",
	"in":  "https://www.google.co.in/search?q=",
	"fr":  "https://www.google.fr/search?q=",
}

func buildGoogleURLs(searchTerm, countryCode, languageCode string, pages, count int) ([]string, error) {
	toScrape := []string{}
	searchTerm = strings.Trim(searchTerm, " ")
	searchTerm = strings.Replace(searchTerm, " ", "+", -1)
	baseUrl, ok := googleDomains[countryCode]
	if ok == false {
		fmt.Printf("No Country Code found")
		return nil, fmt.Errorf("wrong country code provided (%s)", countryCode)
	}

	for i := 0; i < pages; i++ {
		start := i * count
		scrapeUrl := fmt.Sprint(baseUrl, searchTerm, "&num=", count, "&hl=", languageCode, "&start=", start, "&filter=0")
		toScrape = append(toScrape, scrapeUrl)
	}

	return toScrape, nil
}

func scrapeClientRequest(url string, proxyString interface{}) (*http.Response, error) {
	baseClient := getScrapeClient(proxyString)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", randomUserAgents())

	res, err := baseClient.Do(req)
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("scraper recieved a non-200 status Code, suggesting a ban")
	}

	if err != nil {
		return nil, err
	}

	return res, nil
}

func GoogleScrape(searchTerm, countryCode, languageCode string, proxyString interface{}, pages, count int, sleepTime time.Duration) ([]models.SearchResults, error) {
	var results []models.SearchResults
	resultCounter := 0

	googlePages, err := buildGoogleURLs(searchTerm, countryCode, languageCode, pages, count)
	if err != nil {
		return nil, err
	}

	for _, url := range googlePages {
		res, err := scrapeClientRequest(url, proxyString)
		if err != nil {
			return nil, err
		}

		data, err := googleResultParsing(res, resultCounter)
		if err != nil {
			return nil, err
		}

		resultCounter += len(data)
		for _, result := range data {
			results = append(results, result)
		}
		time.Sleep(sleepTime)
	}
	return results, nil
}

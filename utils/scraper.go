package utils

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"go-google-scraper/models"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func getScrapeClient(proxyString interface{}) *http.Client {
	switch v := proxyString.(type) {
	case string:
		proxyUrl, _ := url.Parse(v)
		return &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			},
		}

	default:
		return &http.Client{}
	}
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

func googleResultParsing(response *http.Response, rank int) ([]models.SearchResults, error) {
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	results := []models.SearchResults{}
	sel := doc.Find("div.g")
	rank++

	for i := range sel.Nodes {
		item := sel.Eq(i)
		linkTag := item.Find("a")
		link, _ := linkTag.Attr("href")
		titleTag := item.Find("h3.r")
		descTag := item.Find("span.st")
		desc := descTag.Text()
		title := titleTag.Text()
		link = strings.Trim(link, " ")

		if link != "" && link != "#" && strings.HasPrefix(link, "/") {
			result := models.SearchResults{
				ResultRank:  rank,
				ResultURL:   link,
				ResultTitle: title,
				ResultDesc:  desc,
			}
			results = append(results, result)
		}
	}
	return results, nil
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

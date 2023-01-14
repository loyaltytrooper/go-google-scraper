package main

import (
	"fmt"
	"go-google-scraper/utils"
	"time"
)

func main() {
	responses, err := utils.GoogleScrape("Rajat Kumar Ventura", "in", "en", nil, 1, 10, time.Duration(time.Millisecond*10))
	if err == nil {
		for _, res := range responses {
			fmt.Printf("Result: %s\n", res)
		}
	} else {
		fmt.Printf("error %s", err.Error())
	}
}

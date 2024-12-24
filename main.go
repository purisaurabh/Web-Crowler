package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("not enough arguments provided")
		fmt.Println("usage: crawler <baseURL> <maxConcurrency> <maxPages>")
		return
	}
	if len(os.Args) > 4 {
		fmt.Println("too many arguments provided")
		return
	}
	rawBaseURL := os.Args[1]
	maxConcurrencyString := os.Args[2]
	maxPagesString := os.Args[3]

	maxConcurrency, err := strconv.Atoi(maxConcurrencyString)
	if err != nil {
		fmt.Printf("Error - maxConcurrency: %v", err)
		return
	}
	maxPages, err := strconv.Atoi(maxPagesString)
	if err != nil {
		fmt.Printf("Error - maxPages: %v", err)
		return
	}

	cfg, err := configure(rawBaseURL, maxConcurrency, maxPages)
	if err != nil {
		fmt.Printf("Error - configure: %v", err)
		return
	}

	fmt.Printf("starting crawl of: %s...\n", rawBaseURL)

	cfg.wg.Add(1)
	go cfg.crawlPage(rawBaseURL)
	cfg.wg.Wait()

	printReport(cfg.pages, rawBaseURL)
}

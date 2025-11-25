package crawler

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

type Page struct {
	URL          string          `json:"url"`
	Count        int             `json:"count"`
	Title        string          `json:"title,omitempty"`
	Description  string          `json:"description,omitempty"`
	Keywords     string          `json:"keywords,omitempty"`
	Author       string          `json:"author,omitempty"`
	Canonical    string          `json:"canonical,omitempty"`
	Language     string          `json:"language,omitempty"`
	Charset      string          `json:"charset,omitempty"`
	OGImage      string          `json:"og_image,omitempty"`
	OGType       string          `json:"og_type,omitempty"`
	OGURL        string          `json:"og_url,omitempty"`
	OGSiteName   string          `json:"og_site_name,omitempty"`
	TwitterCard  string          `json:"twitter_card,omitempty"`
	TwitterSite  string          `json:"twitter_site,omitempty"`
	TwitterImage string          `json:"twitter_image,omitempty"`
	Suggestions  *AnalysisResult `json:"suggestions,omitempty"`
}

func PrintReport(pages map[string]*PageData, baseURL string, jsonOutput bool, outputFile string) {
	if outputFile != "" {
		saveReportToFile(pages, outputFile, jsonOutput)
		return
	}

	if jsonOutput {
		printJSONReport(pages)
		return
	}

	fmt.Printf(`
=============================
  REPORT for %s
=============================
`, baseURL)

	sortedPages := sortPages(pages)
	for _, page := range sortedPages {
		url := page.URL
		count := page.Count
		fmt.Printf("Found %d internal links to %s\n", count, url)
	}
}

func saveReportToFile(pages map[string]*PageData, outputFile string, jsonOutput bool) {
	f, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer f.Close()

	sortedPages := sortPages(pages)

	if jsonOutput {
		jsonData, err := json.MarshalIndent(sortedPages, "", "  ")
		if err != nil {
			fmt.Printf("Error marshalling JSON: %v\n", err)
			return
		}
		f.Write(jsonData)
	} else {
		for _, page := range sortedPages {
			line := fmt.Sprintf("Found %d internal links to %s\n", page.Count, page.URL)
			f.WriteString(line)
		}
	}
	fmt.Printf("Report saved to %s\n", outputFile)
}

func printJSONReport(pages map[string]*PageData) {
	sortedPages := sortPages(pages)
	jsonData, err := json.MarshalIndent(sortedPages, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling JSON: %v\n", err)
		return
	}
	fmt.Println(string(jsonData))
}

func sortPages(pages map[string]*PageData) []Page {
	pagesSlice := []Page{}
	for url, data := range pages {
		pagesSlice = append(pagesSlice, Page{
			URL:          url,
			Count:        data.LinkCount,
			Title:        data.Title,
			Description:  data.Description,
			Keywords:     data.Keywords,
			Author:       data.Author,
			Canonical:    data.Canonical,
			Language:     data.Language,
			Charset:      data.Charset,
			OGImage:      data.OGImage,
			OGType:       data.OGType,
			OGURL:        data.OGURL,
			OGSiteName:   data.OGSiteName,
			TwitterCard:  data.TwitterCard,
			TwitterSite:  data.TwitterSite,
			TwitterImage: data.TwitterImage,
			Suggestions:  data.Suggestions,
		})
	}
	sort.Slice(pagesSlice, func(i, j int) bool {
		if pagesSlice[i].Count == pagesSlice[j].Count {
			return pagesSlice[i].URL < pagesSlice[j].URL
		}
		return pagesSlice[i].Count > pagesSlice[j].Count
	})
	return pagesSlice
}

package crawler

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

func (cfg *Config) CrawlPage(rawCurrentURL string) {
	cfg.ConcurrencyControl <- struct{}{}
	defer func() {
		<-cfg.ConcurrencyControl
		cfg.WG.Done()
	}()

	if cfg.PagesLen() >= cfg.MaxPages {
		return
	}

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error - crawlPage: couldn't parse URL '%s': %v\n", rawCurrentURL, err)
		return
	}

	// skip other websites
	if currentURL.Hostname() != cfg.BaseURL.Hostname() {
		return
	}

	// Check robots.txt
	if !cfg.Robots.IsAllowed(rawCurrentURL) {
		return
	}

	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error - normalizedURL: %v", err)
		return
	}

	isFirst := cfg.addPageVisit(normalizedURL)
	if !isFirst {
		return
	}

	fmt.Printf("crawling %s\n", rawCurrentURL)

	// Rate limiting
	time.Sleep(cfg.RateLimit)

	htmlBody, err := cfg.getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error - getHTML: %v", err)
		return
	}

	// Extract metadata
	title, description, keywords, author, canonical, language, charset, ogImage, ogType, ogURL, ogSiteName, twitterCard, twitterSite, twitterImage := extractMetadata(htmlBody)
	cfg.Mu.Lock()
	if data, ok := cfg.Pages[normalizedURL]; ok {
		data.Title = title
		data.Description = description
		data.Keywords = keywords
		data.Author = author
		data.Canonical = canonical
		data.Language = language
		data.Charset = charset
		data.OGImage = ogImage
		data.OGType = ogType
		data.OGURL = ogURL
		data.OGSiteName = ogSiteName
		data.TwitterCard = twitterCard
		data.TwitterSite = twitterSite
		data.TwitterImage = twitterImage

		// AI Analysis if enabled
		if cfg.Analyzer != nil {
			analysis, err := cfg.Analyzer.AnalyzePage(rawCurrentURL, title, description)
			if err != nil {
				fmt.Printf("Warning - AI analysis failed: %v\n", err)
			} else {
				data.Suggestions = analysis
			}
		}
	}
	cfg.Mu.Unlock()

	nextURLs, err := getURLsFromHTML(htmlBody, cfg.BaseURL)
	if err != nil {
		fmt.Printf("Error - getURLsFromHTML: %v", err)
		return
	}

	for _, nextURL := range nextURLs {
		cfg.WG.Add(1)
		go cfg.CrawlPage(nextURL)
	}
}

func extractMetadata(htmlBody string) (title, description, keywords, author, canonical, language, charset, ogImage, ogType, ogURL, ogSiteName, twitterCard, twitterSite, twitterImage string) {
	doc, err := html.Parse(strings.NewReader(htmlBody))
	if err != nil {
		return
	}

	var (
		ogDescription      string
		twitterDescription string
		jsonLDDescription  string
	)

	var f func(*html.Node)
	f = func(n *html.Node) {
		// Extract language from html tag
		if n.Type == html.ElementNode && n.Data == "html" {
			for _, a := range n.Attr {
				if a.Key == "lang" {
					language = a.Val
				}
			}
		}

		if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
			title = n.FirstChild.Data
		}

		if n.Type == html.ElementNode && n.Data == "link" {
			var rel, href string
			for _, a := range n.Attr {
				if a.Key == "rel" {
					rel = a.Val
				}
				if a.Key == "href" {
					href = a.Val
				}
			}
			if rel == "canonical" {
				canonical = href
			}
		}

		if n.Type == html.ElementNode && n.Data == "meta" {
			var name, property, content, charsetAttr string
			for _, a := range n.Attr {
				if a.Key == "name" {
					name = a.Val
				}
				if a.Key == "property" {
					property = a.Val
				}
				if a.Key == "content" {
					content = a.Val
				}
				if a.Key == "charset" {
					charsetAttr = a.Val
				}
			}

			// Charset
			if charsetAttr != "" {
				charset = charsetAttr
			}

			// Standard meta tags
			if name == "description" {
				description = content
			} else if name == "keywords" {
				keywords = content
			} else if name == "author" {
				author = content
			}

			// Open Graph
			if property == "og:description" {
				ogDescription = content
			} else if property == "og:image" {
				ogImage = content
			} else if property == "og:type" {
				ogType = content
			} else if property == "og:url" {
				ogURL = content
			} else if property == "og:site_name" {
				ogSiteName = content
			}

			// Twitter Cards
			if name == "twitter:description" {
				twitterDescription = content
			} else if name == "twitter:card" {
				twitterCard = content
			} else if name == "twitter:site" {
				twitterSite = content
			} else if name == "twitter:image" {
				twitterImage = content
			}
		}

		if n.Type == html.ElementNode && n.Data == "script" {
			isJSONLD := false
			for _, a := range n.Attr {
				if a.Key == "type" && a.Val == "application/ld+json" {
					isJSONLD = true
					break
				}
			}
			if isJSONLD && n.FirstChild != nil {
				jsonContent := n.FirstChild.Data
				// Basic string extraction to avoid full JSON parsing overhead/complexity for now
				// Looking for "description": "..."
				if idx := strings.Index(jsonContent, `"description":`); idx != -1 {
					start := idx + 14
					rest := jsonContent[start:]
					rest = strings.TrimSpace(rest)
					if strings.HasPrefix(rest, `"`) {
						rest = rest[1:]
						if end := strings.Index(rest, `",`); end != -1 {
							jsonLDDescription = rest[:end]
						} else if end := strings.Index(rest, `"`+"\n"); end != -1 { // Handle end of line
							jsonLDDescription = rest[:end]
						}
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	if description == "" {
		if ogDescription != "" {
			description = ogDescription
		} else if twitterDescription != "" {
			description = twitterDescription
		} else if jsonLDDescription != "" {
			description = jsonLDDescription
		}
	}

	return
}

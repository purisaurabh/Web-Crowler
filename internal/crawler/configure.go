package crawler

import (
	"fmt"
	"net/url"
	"sync"
	"time"
)

type PageData struct {
	LinkCount    int
	Title        string
	Description  string
	Keywords     string
	Author       string
	Canonical    string
	Language     string
	Charset      string
	OGImage      string
	OGType       string
	OGURL        string
	OGSiteName   string
	TwitterCard  string
	TwitterSite  string
	TwitterImage string
	Suggestions  *AnalysisResult
}

type Config struct {
	Pages              map[string]*PageData
	BaseURL            *url.URL
	Mu                 *sync.Mutex
	ConcurrencyControl chan struct{}
	WG                 *sync.WaitGroup
	MaxPages           int
	Robots             *RobotsChecker
	RateLimit          time.Duration
	UserAgent          string
	JSONOutput         bool
	Analyzer           *AIAnalyzer
}

func (cfg *Config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.Mu.Lock()
	defer cfg.Mu.Unlock()

	if pageData, visited := cfg.Pages[normalizedURL]; visited {
		pageData.LinkCount++
		return false
	}

	cfg.Pages[normalizedURL] = &PageData{LinkCount: 1}
	return true
}

func (cfg *Config) PagesLen() int {
	cfg.Mu.Lock()
	defer cfg.Mu.Unlock()
	return len(cfg.Pages)
}

func Configure(rawBaseURL string, maxConcurrency int, maxPages int, rateLimit time.Duration, userAgent string, jsonOutput bool, apiKey string, aiProvider string) (*Config, error) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse base URL: %v", err)
	}

	var analyzer *AIAnalyzer
	if apiKey != "" {
		analyzer = NewAIAnalyzer(apiKey, aiProvider)
	}

	return &Config{
		Pages:              make(map[string]*PageData),
		BaseURL:            baseURL,
		Mu:                 &sync.Mutex{},
		ConcurrencyControl: make(chan struct{}, maxConcurrency),
		WG:                 &sync.WaitGroup{},
		MaxPages:           maxPages,
		Robots:             NewRobotsChecker(baseURL, userAgent),
		RateLimit:          rateLimit,
		UserAgent:          userAgent,
		JSONOutput:         jsonOutput,
		Analyzer:           analyzer,
	}, nil
}

package crawler

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type RobotsChecker struct {
	disallowed []string
	mu         sync.Mutex
	userAgent  string
}

func NewRobotsChecker(baseURL *url.URL, userAgent string) *RobotsChecker {
	rc := &RobotsChecker{
		userAgent: userAgent,
	}
	rc.fetchRobotsTxt(baseURL)
	return rc
}

func (rc *RobotsChecker) fetchRobotsTxt(baseURL *url.URL) {
	robotsURL := baseURL.Scheme + "://" + baseURL.Host + "/robots.txt"
	req, err := http.NewRequest("GET", robotsURL, nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", rc.userAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	rc.parse(string(body))
}

func (rc *RobotsChecker) parse(body string) {
	lines := strings.Split(body, "\n")
	userAgentMatches := false
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Remove comments
		if idx := strings.Index(line, "#"); idx != -1 {
			line = line[:idx]
		}
		if line == "" {
			continue
		}

		lowerLine := strings.ToLower(line)
		if strings.HasPrefix(lowerLine, "user-agent:") {
			agent := strings.TrimSpace(line[11:])
			if agent == "*" || strings.EqualFold(agent, rc.userAgent) {
				userAgentMatches = true
			} else {
				userAgentMatches = false
			}
		} else if userAgentMatches && strings.HasPrefix(lowerLine, "disallow:") {
			path := strings.TrimSpace(line[9:])
			if path != "" {
				rc.disallowed = append(rc.disallowed, path)
			}
		}
	}
}

func (rc *RobotsChecker) IsAllowed(u string) bool {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return false
	}
	path := parsedURL.Path
	if path == "" {
		path = "/"
	}
	for _, disallowed := range rc.disallowed {
		if strings.HasPrefix(path, disallowed) {
			return false
		}
	}
	return true
}

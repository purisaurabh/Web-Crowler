package crawler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type AIAnalyzer struct {
	apiKey   string
	provider string // "openai", "gemini", "anthropic"
}

type AnalysisResult struct {
	SEO            []string `json:"seo"`
	ContentQuality []string `json:"content_quality"`
	Accessibility  []string `json:"accessibility"`
	Performance    []string `json:"performance"`
}

func NewAIAnalyzer(apiKey, provider string) *AIAnalyzer {
	if provider == "" {
		provider = "openai"
	}
	return &AIAnalyzer{
		apiKey:   apiKey,
		provider: provider,
	}
}

func (a *AIAnalyzer) AnalyzePage(url, title, description string) (*AnalysisResult, error) {
	if a.apiKey == "" {
		return nil, fmt.Errorf("API key not provided")
	}

	switch a.provider {
	case "openai":
		return a.analyzeWithOpenAI(url, title, description)
	case "gemini":
		return a.analyzeWithGemini(url, title, description)
	case "anthropic":
		return a.analyzeWithAnthropic(url, title, description)
	default:
		return nil, fmt.Errorf("unsupported AI provider: %s", a.provider)
	}
}

func (a *AIAnalyzer) analyzeWithOpenAI(url, title, description string) (*AnalysisResult, error) {
	prompt := fmt.Sprintf(`Analyze this webpage and provide improvement suggestions in JSON format:

URL: %s
Title: %s
Description: %s

Provide suggestions in these categories:
1. SEO (meta tags, title optimization, structured data)
2. Content Quality (title/description effectiveness)
3. Accessibility (common issues)
4. Performance (potential bottlenecks)

Return ONLY a JSON object with this structure:
{
  "seo": ["suggestion 1", "suggestion 2"],
  "content_quality": ["suggestion 1"],
  "accessibility": ["suggestion 1"],
  "performance": ["suggestion 1"]
}`, url, title, description)

	requestBody := map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You are a web development expert specializing in SEO, accessibility, and performance optimization. Provide concise, actionable suggestions.",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"temperature": 0.7,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+a.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("OpenAI API error: %s - %s", resp.Status, string(body))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	// Parse the JSON response from the AI
	var analysis AnalysisResult
	content := result.Choices[0].Message.Content

	// Try to extract JSON from markdown code blocks if present
	if idx := bytes.Index([]byte(content), []byte("```json")); idx != -1 {
		content = content[idx+7:]
		if endIdx := bytes.Index([]byte(content), []byte("```")); endIdx != -1 {
			content = content[:endIdx]
		}
	} else if idx := bytes.Index([]byte(content), []byte("```")); idx != -1 {
		content = content[idx+3:]
		if endIdx := bytes.Index([]byte(content), []byte("```")); endIdx != -1 {
			content = content[:endIdx]
		}
	}

	if err := json.Unmarshal([]byte(content), &analysis); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %v - Content: %s", err, content)
	}

	return &analysis, nil
}

func (a *AIAnalyzer) analyzeWithGemini(url, title, description string) (*AnalysisResult, error) {
	prompt := fmt.Sprintf(`Analyze this webpage and provide improvement suggestions in JSON format:

URL: %s
Title: %s
Description: %s

Provide suggestions in these categories:
1. SEO (meta tags, title optimization, structured data)
2. Content Quality (title/description effectiveness)
3. Accessibility (common issues)
4. Performance (potential bottlenecks)

Return ONLY a JSON object with this structure:
{
  "seo": ["suggestion 1", "suggestion 2"],
  "content_quality": ["suggestion 1"],
  "accessibility": ["suggestion 1"],
  "performance": ["suggestion 1"]
}`, url, title, description)

	requestBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{
						"text": prompt,
					},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"temperature": 0.7,
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	apiURL := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key=%s", a.apiKey)
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Gemini API error: %s - %s", resp.Status, string(body))
	}

	var result struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no response from Gemini")
	}

	// Parse the JSON response from the AI
	var analysis AnalysisResult
	content := result.Candidates[0].Content.Parts[0].Text

	// Try to extract JSON from markdown code blocks if present
	if idx := bytes.Index([]byte(content), []byte("```json")); idx != -1 {
		content = content[idx+7:]
		if endIdx := bytes.Index([]byte(content), []byte("```")); endIdx != -1 {
			content = content[:endIdx]
		}
	} else if idx := bytes.Index([]byte(content), []byte("```")); idx != -1 {
		content = content[idx+3:]
		if endIdx := bytes.Index([]byte(content), []byte("```")); endIdx != -1 {
			content = content[:endIdx]
		}
	}

	if err := json.Unmarshal([]byte(content), &analysis); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %v - Content: %s", err, content)
	}

	return &analysis, nil
}

func (a *AIAnalyzer) analyzeWithAnthropic(url, title, description string) (*AnalysisResult, error) {
	// Placeholder for Anthropic implementation
	return nil, fmt.Errorf("Anthropic integration not yet implemented")
}

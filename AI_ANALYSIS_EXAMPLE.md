# AI Analysis Example

This example demonstrates how to use the AI-powered analysis feature.

## Prerequisites
- OpenAI API key (get one at https://platform.openai.com/api-keys)

## Basic Usage

```bash
# Crawl a website with AI analysis
go run cmd/crawler/main.go \
  -url https://cadicient.com \
  -pages 5 \
  -json \
  -analyze \
  -api-key YOUR_OPENAI_API_KEY \
  -out analysis.json
```

## Example Output

The output will include AI-generated suggestions for each page:

```json
[
  {
    "URL": "cadicient.com",
    "Count": 1,
    "Title": "Cadicient",
    "Description": "Cadicient Engineering Consultants provides...",
    "suggestions": {
      "seo": [
        "Add structured data (Schema.org) for better search visibility",
        "Consider adding more descriptive keywords to the title"
      ],
      "content_quality": [
        "The description is comprehensive but could be more concise"
      ],
      "accessibility": [
        "Ensure all images have alt text",
        "Check color contrast ratios for WCAG compliance"
      ],
      "performance": [
        "Consider lazy loading images",
        "Minify CSS and JavaScript files"
      ]
    }
  }
]
```

## Supported AI Providers

Currently supported:
- **OpenAI** (default): Uses GPT-3.5-turbo
- **Gemini**: Coming soon
- **Anthropic Claude**: Coming soon

## Notes

- AI analysis adds ~2-3 seconds per page
- Requires an active internet connection
- API costs apply based on your provider's pricing

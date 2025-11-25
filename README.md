# Web Crawler

A Golang CLI application that crawls a website and generates an internal links report.

## Features
-   **Concurrent Crawling**: Uses Goroutines to crawl pages in parallel.
-   **Robots.txt Support**: Respects `robots.txt` rules (User-agent: *).
-   **Rate Limiting**: Configurable delay between requests to avoid overwhelming servers.
-   **JSON Output**: Option to export report in JSON format.
-   **File Output**: Save report to a specific file.
-   **Rich Page Data**: Extracts comprehensive metadata including:
    -   Title, Description, Keywords, Author
    -   Canonical URL, Language, Charset
    -   Open Graph data (image, type, URL, site name)
    -   Twitter Card data (card type, site, image)
-   **Configurable User-Agent**: Set custom User-Agent string.
-   **AI-Powered Analysis**: Get AI-generated suggestions for SEO, content quality, accessibility, and performance improvements.
-   **Environment Variables**: Load API keys from `.env` file for security.

## Setup

1.  **Clone the repository**
2.  **Create a `.env` file** (optional, for AI analysis):
    ```bash
    cp .env.example .env
    # Edit .env and add your API key
    ```

## Usage

```bash
go run cmd/crawler/main.go -url <baseURL> [flags]
```

### Flags
-   `-url`: Base URL to crawl (required).
-   `-concurrency`: Maximum number of concurrent requests (default 10).
-   `-pages`: Maximum number of pages to crawl (default 100).
-   `-json`: Output report in JSON format (default false).
-   `-out`: Output file path (optional).
-   `-user-agent`: User-Agent string to use (default "Crawler").
-   `-delay`: Delay between requests (default 500ms).
-   `-analyze`: Enable AI-powered analysis (requires API key in .env file).
-   `-ai-provider`: AI provider to use: openai/gemini/anthropic (default "openai").

### Examples

**Basic Crawl:**
```bash
go run cmd/crawler/main.go -url https://wagslane.dev
```

**Save to JSON File:**
```bash
go run cmd/crawler/main.go -url https://wagslane.dev -json -out report.json
```

**AI-Powered Analysis:**
```bash
go run cmd/crawler/main.go -url https://cadicient.com -json -analyze -api-key YOUR_OPENAI_API_KEY -out analysis.json
```

**Custom Configuration:**
```bash
go run cmd/crawler/main.go -url https://wagslane.dev -concurrency 20 -pages 50 -delay 100ms -user-agent "MyBot"
```

## Project Structure
-   `cmd/crawler`: Entry point of the application.
-   `internal/crawler`: Core logic (crawler, configuration, robots.txt, etc.).
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/purisaurabh/web-crowler/internal/crawler"
)

func main() {
	// Load .env file (ignore error if file doesn't exist)
	_ = godotenv.Load()

	urlFlag := flag.String("url", "", "Base URL to crawl")
	concurrencyFlag := flag.Int("concurrency", 10, "Maximum number of concurrent requests")
	pagesFlag := flag.Int("pages", 100, "Maximum number of pages to crawl")
	jsonFlag := flag.Bool("json", false, "Output report in JSON format")
	outFlag := flag.String("out", "", "Output file path (optional)")
	userAgentFlag := flag.String("user-agent", "Crawler", "User-Agent string to use")
	delayFlag := flag.Duration("delay", 500*time.Millisecond, "Delay between requests")
	analyzeFlag := flag.Bool("analyze", false, "Enable AI-powered analysis (requires API key in .env file)")
	aiProviderFlag := flag.String("ai-provider", "openai", "AI provider (openai/gemini/anthropic)")

	flag.Parse()

	if *urlFlag == "" {
		// Fallback to positional arguments for backward compatibility or ease of use
		args := flag.Args()
		if len(args) >= 1 {
			*urlFlag = args[0]
		}
		if len(args) >= 2 {
			fmt.Sscanf(args[1], "%d", concurrencyFlag)
		}
		if len(args) >= 3 {
			fmt.Sscanf(args[2], "%d", pagesFlag)
		}
	}

	if *urlFlag == "" {
		fmt.Println("usage: crawler -url <baseURL> [-concurrency <n>] [-pages <n>] [-json] [-out <file>] [-user-agent <s>] [-delay <d>] [-analyze] [-ai-provider <provider>]")
		fmt.Println("\nFor AI analysis, set API key in .env file:")
		fmt.Println("  OPENAI_API_KEY=your-key-here")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Handle AI analysis - only from environment variables
	apiKey := ""
	if *analyzeFlag {
		// Check which API key is available
		openaiKey := os.Getenv("OPENAI_API_KEY")
		geminiKey := os.Getenv("GEMINI_API_KEY")
		anthropicKey := os.Getenv("ANTHROPIC_API_KEY")

		// Auto-detect provider based on available API key
		if *aiProviderFlag == "openai" && openaiKey != "" {
			apiKey = openaiKey
		} else if *aiProviderFlag == "gemini" && geminiKey != "" {
			apiKey = geminiKey
		} else if *aiProviderFlag == "anthropic" && anthropicKey != "" {
			apiKey = anthropicKey
		} else if geminiKey != "" {
			// Auto-select Gemini if available
			apiKey = geminiKey
			*aiProviderFlag = "gemini"
		} else if openaiKey != "" {
			// Auto-select OpenAI if available
			apiKey = openaiKey
			*aiProviderFlag = "openai"
		} else if anthropicKey != "" {
			// Auto-select Anthropic if available
			apiKey = anthropicKey
			*aiProviderFlag = "anthropic"
		}

		if apiKey == "" {
			fmt.Println("Error: AI analysis requires an API key in .env file")
			fmt.Println("Add one of the following to your .env file:")
			fmt.Println("  OPENAI_API_KEY=your-openai-key")
			fmt.Println("  GEMINI_API_KEY=your-gemini-key")
			fmt.Println("  ANTHROPIC_API_KEY=your-anthropic-key")
			os.Exit(1)
		}
	}

	cfg, err := crawler.Configure(*urlFlag, *concurrencyFlag, *pagesFlag, *delayFlag, *userAgentFlag, *jsonFlag, apiKey, *aiProviderFlag)
	if err != nil {
		fmt.Printf("Error - configure: %v", err)
		return
	}

	if !*jsonFlag && *outFlag == "" {
		fmt.Printf("starting crawl of: %s...\n", *urlFlag)
		if *analyzeFlag {
			fmt.Printf("AI analysis enabled using %s\n", *aiProviderFlag)
		}
	}

	cfg.WG.Add(1)
	go cfg.CrawlPage(*urlFlag)
	cfg.WG.Wait()

	crawler.PrintReport(cfg.Pages, *urlFlag, *jsonFlag, *outFlag)
}

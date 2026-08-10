package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

// Crawler stores everything our crawler needs.
type Crawler struct {
	client      *http.Client
	visited     map[string]bool
	statusCodes map[int]int
	errors      map[string]int // Track error types

	mu sync.Mutex

	maxPages  int
	workers   int
	baseHost  string
	delay     time.Duration // Politeness delay between requests
	lastFetch time.Time     // Track last request time
}

// NewCrawler creates a new crawler.
func NewCrawler(startURL string, workers, maxPages int) (*Crawler, error) {
	parsedURL, err := url.Parse(startURL)
	if err != nil {
		return nil, err
	}

	return &Crawler{
		client: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 100,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		visited:     make(map[string]bool),
		statusCodes: make(map[int]int),
		errors:      make(map[string]int),
		maxPages:    maxPages,
		workers:     workers,
		baseHost:    parsedURL.Host,
		delay:       time.Millisecond * 100, // Default 100ms politeness delay
	}, nil
}

// SetDelay sets the politeness delay between requests.
func (c *Crawler) SetDelay(delay time.Duration) {
	c.delay = delay
}

// normalizeURL cleans a URL.
func (c *Crawler) normalizeURL(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	parsed.Fragment = ""
	parsed.Scheme = strings.ToLower(parsed.Scheme)
	parsed.Host = strings.ToLower(parsed.Host)
	return parsed.String()
}

// isAllowed checks whether we are allowed to crawl this URL.
func (c *Crawler) isAllowed(rawURL string) bool {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return false
	}
	return parsed.Host == c.baseHost
}

// markVisited checks and records a URL.
// Returns true if this is a new URL and we are under the page limit.
func (c *Crawler) markVisited(pageURL string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.maxPages > 0 && len(c.visited) >= c.maxPages {
		return false
	}
	if c.visited[pageURL] {
		return false
	}
	c.visited[pageURL] = true
	return true
}

// fetch downloads a webpage and extracts its raw links.
func (c *Crawler) fetch(pageURL string) ([]string, int, error) {
	c.mu.Lock()
	timeSinceLastFetch := time.Since(c.lastFetch)
	if timeSinceLastFetch < c.delay {
		c.mu.Unlock()
		time.Sleep(c.delay - timeSinceLastFetch)
		c.mu.Lock()
	}
	c.lastFetch = time.Now()
	c.mu.Unlock()

	req, err := http.NewRequest(http.MethodGet, pageURL, nil)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("User-Agent", "GentleMan")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	status := resp.StatusCode

	// Read body ONCE
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("failed to read body: %w", err)
	}
	bodyString := string(bodyBytes)

	// Bot detection check
	if strings.Contains(bodyString, "verify you are human") ||
		strings.Contains(bodyString, "captcha") ||
		strings.Contains(bodyString, "unusual traffic") {
		return nil, 200, fmt.Errorf("bot detection page, not real content")
	}

	if status < 200 || status >= 300 {
		return nil, status, fmt.Errorf("HTTP %d: %s", status, resp.Status)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		return nil, status, fmt.Errorf("not HTML: %s", contentType)
	}

	// Parse from string — NOT the drained body
	doc, err := html.Parse(strings.NewReader(bodyString))
	if err != nil {
		return nil, status, err
	}

	var links []string
	var walk func(*html.Node)
	walk = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "a" {
			for _, attr := range node.Attr {
				if attr.Key == "href" {
					links = append(links, attr.Val)
				}
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}
	walk(doc)

	return links, status, nil
}

// extractLinks converts raw links into valid absolute URLs.
func (c *Crawler) extractLinks(pageURL string, rawLinks []string) []string {
	base, err := url.Parse(pageURL)
	if err != nil {
		return nil
	}

	var links []string
	seen := make(map[string]bool) // Deduplicate links from same page

	for _, rawLink := range rawLinks {
		rawLink = strings.TrimSpace(rawLink)
		if rawLink == "" {
			continue
		}

		parsedLink, err := url.Parse(rawLink)
		if err != nil {
			continue
		}

		fullURL := base.ResolveReference(parsedLink)
		fullURL.Fragment = ""
		normalized := c.normalizeURL(fullURL.String())

		if normalized == "" || !c.isAllowed(normalized) {
			continue
		}

		// Avoid returning duplicate links from the same page
		if !seen[normalized] {
			seen[normalized] = true
			links = append(links, normalized)
		}
	}

	return links
}

// crawlWorker is one crawler worker.
func (c *Crawler) crawlWorker(id int, jobs <-chan string, results chan<- []string, wg *sync.WaitGroup) {
	defer wg.Done()

	for pageURL := range jobs {
		fmt.Printf("🕷️  [Worker %3d] Crawling: %s\n", id, pageURL)

		rawLinks, status, err := c.fetch(pageURL)

		// Track status codes and errors
		c.mu.Lock()
		c.statusCodes[status]++
		if err != nil {
			// Categorize errors
			errType := "unknown"
			switch {
			case status == 404:
				errType = "404 Not Found"
			case status == 403:
				errType = "403 Forbidden"
			case status == 429:
				errType = "429 Rate Limited"
			case status >= 500:
				errType = fmt.Sprintf("5xx Server Error (%d)", status)
			case status >= 400:
				errType = fmt.Sprintf("4xx Client Error (%d)", status)
			case status >= 300:
				errType = fmt.Sprintf("3xx Redirect (%d)", status)
			case status == 0:
				errType = "Network Error"
			}
			c.errors[errType]++
		}
		c.mu.Unlock()

		if err != nil {
			fmt.Printf("   ❌ [Worker %3d] Error: %v\n", id, err)
			results <- nil // Still signal completion
			continue
		}

		links := c.extractLinks(pageURL, rawLinks)
		fmt.Printf("   ✅ [Worker %3d] Found %d unique links\n", id, len(links))
		results <- links
	}
}

// Start begins crawling using a proper worker pool.
func (c *Crawler) Start(startURL string) {
	startURL = c.normalizeURL(startURL)
	if startURL == "" {
		fmt.Println("❌ Invalid start URL")
		return
	}

	if !c.markVisited(startURL) {
		fmt.Println("❌ Max pages reached or start URL not allowed, exiting.")
		return
	}

	fmt.Println("🚀 Starting crawl...")
	fmt.Printf("📋 Config: %d workers | %d max pages | %v delay\n",
		c.workers, c.maxPages, c.delay)
	fmt.Println(strings.Repeat("─", 60))

	jobs := make(chan string, c.workers)      // Buffered channel
	results := make(chan []string, c.workers) // Buffered channel

	var workerWG sync.WaitGroup
	for i := 1; i <= c.workers; i++ {
		workerWG.Add(1)
		go c.crawlWorker(i, jobs, results, &workerWG)
	}

	// Track statistics
	startTime := time.Now()
	queue := []string{startURL}
	active := 0
	processed := 0

	// Progress reporting ticker
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// Coordinator loop
	for {
		// Check if done
		if active == 0 && len(queue) == 0 {
			break
		}

		// Prepare a send only if there's work and an idle worker
		var sendCh chan string
		var sendURL string
		if len(queue) > 0 && active < c.workers {
			sendCh = jobs
			sendURL = queue[0]
		}

		select {
		case sendCh <- sendURL:
			// Successfully dispatched a job
			queue = queue[1:]
			active++

		case links := <-results:
			active--
			processed++

			// Add discovered links to queue
			for _, link := range links {
				if c.markVisited(link) {
					queue = append(queue, link)
				}
			}

		case <-ticker.C:
			// Progress update every second
			elapsed := time.Since(startTime)
			c.mu.Lock()
			visited := len(c.visited)
			c.mu.Unlock()

			fmt.Printf("📊 Progress: %d/%d pages | %d active | %d queued | %v elapsed\n",
				visited, c.maxPages, active, len(queue), elapsed.Round(time.Millisecond))
		}
	}

	elapsed := time.Since(startTime)
	fmt.Println(strings.Repeat("─", 60))
	fmt.Printf("✅ Crawl complete in %v\n", elapsed.Round(time.Millisecond))

	// No more jobs to send; close channel so workers exit
	close(jobs)
	workerWG.Wait()
}

// PrintSummary prints final crawl statistics.
func (c *Crawler) PrintSummary() {
	c.mu.Lock()
	defer c.mu.Unlock()

	fmt.Println()
	fmt.Println("╔════════════════════════════════════╗")
	fmt.Println("║        CRAWL SUMMARY              ║")
	fmt.Println("╚════════════════════════════════════╝")
	fmt.Println()

	// Success rate
	totalRequests := 0
	for _, count := range c.statusCodes {
		totalRequests += count
	}

	successfulRequests := c.statusCodes[200]
	successRate := 0.0
	if totalRequests > 0 {
		successRate = float64(successfulRequests) / float64(totalRequests) * 100
	}

	fmt.Printf("📈 Total Requests:    %d\n", totalRequests)
	fmt.Printf("✅ Successful (200):  %d (%.1f%%)\n", successfulRequests, successRate)
	fmt.Printf("📄 Pages Visited:     %d\n", len(c.visited))
	fmt.Println()

	// HTTP Status Codes breakdown
	fmt.Println("📊 HTTP Status Codes:")
	fmt.Println("   ┌──────────┬──────────┬─────────────────────┐")
	fmt.Println("   │  Status  │  Count   │  Description        │")
	fmt.Println("   ├──────────┼──────────┼─────────────────────┤")

	statusDescriptions := map[int]string{
		200: "OK",
		301: "Moved Permanently",
		302: "Found",
		304: "Not Modified",
		400: "Bad Request",
		401: "Unauthorized",
		403: "Forbidden",
		404: "Not Found",
		429: "Too Many Requests",
		500: "Internal Server Error",
		502: "Bad Gateway",
		503: "Service Unavailable",
		0:   "Network Error",
	}

	// Sort status codes for consistent output
	for code := 200; code <= 600; code++ {
		if count, exists := c.statusCodes[code]; exists && count > 0 {
			desc, ok := statusDescriptions[code]
			if !ok {
				desc = "Unknown"
			}
			fmt.Printf("   │  %3d     │  %4d    │ %-20s│\n", code, count, desc)
		}
	}

	// Network errors (status 0)
	if count, exists := c.statusCodes[0]; exists && count > 0 {
		fmt.Printf("   │   0     │  %4d    │ %-20s│\n", count, "Network Error")
	}

	fmt.Println("   └──────────┴──────────┴─────────────────────┘")
	fmt.Println()

	// Error breakdown
	if len(c.errors) > 0 {
		fmt.Println("🔍 Error Breakdown:")
		for errType, count := range c.errors {
			fmt.Printf("   • %s: %d\n", errType, count)
		}
		fmt.Println()
	}

	fmt.Println(strings.Repeat("═", 60))
}

func main() {
	fmt.Println("╔════════════════════════════════════════════════╗")
	fmt.Println("║         GO WEB CRAWLER - ENHANCED              ║")
	fmt.Println("╚════════════════════════════════════════════════╝")
	fmt.Println()

	startURL := "https://www.amazon.in/s?i=electronics&rh=n%3A1805560031%2Cp_36%3A2460000-4510000"
	workers := 50
	maxPages := 100

	crawler, err := NewCrawler(startURL, workers, maxPages)
	if err != nil {
		fmt.Println("❌ Crawler error:", err)
		return
	}

	// Configure politeness
	crawler.SetDelay(200 * time.Millisecond) // 200ms between requests

	fmt.Printf("🎯 Target: %s\n", startURL)
	fmt.Println()

	start := time.Now()
	crawler.Start(startURL)
	crawler.PrintSummary()

	fmt.Printf("\n⏱️  Total Time: %v\n", time.Since(start).Round(time.Millisecond))
}

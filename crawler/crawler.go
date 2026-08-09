package main

import (
	"fmt"
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

	mu sync.Mutex

	maxPages int
	workers  int
	baseHost string
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
		},
		visited:     make(map[string]bool),
		statusCodes: make(map[int]int),
		maxPages:    maxPages,
		workers:     workers,
		baseHost:    parsedURL.Host,
	}, nil
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
	req, err := http.NewRequest(http.MethodGet, pageURL, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("User-Agent", "GoCrawler/1.0")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	status := resp.StatusCode
	if status < 200 || status >= 300 {
		return nil, status, fmt.Errorf("HTTP status: %s", resp.Status)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		return nil, status, fmt.Errorf("not an HTML page")
	}

	doc, err := html.Parse(resp.Body)
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
		links = append(links, normalized)
	}
	return links
}

// crawlWorker is one crawler worker.
func (c *Crawler) crawlWorker(jobs <-chan string, results chan<- []string, wg *sync.WaitGroup) {
	defer wg.Done()
	for pageURL := range jobs {
		fmt.Println("🕷️ Crawling:", pageURL)

		rawLinks, status, err := c.fetch(pageURL)

		c.mu.Lock()
		c.statusCodes[status]++
		c.mu.Unlock()

		if err != nil {
			fmt.Println("   Error:", err)
			results <- nil // still signal completion
			continue
		}

		links := c.extractLinks(pageURL, rawLinks)
		fmt.Printf("   Found %d links\n", len(links))
		results <- links
	}
}

// Start begins crawling using a proper worker pool.
func (c *Crawler) Start(startURL string) {
	startURL = c.normalizeURL(startURL)
	if !c.markVisited(startURL) {
		fmt.Println("Max pages reached or start URL not allowed, exiting.")
		return
	}

	jobs := make(chan string)
	results := make(chan []string)

	var workerWG sync.WaitGroup
	for i := 0; i < c.workers; i++ {
		workerWG.Add(1)
		go c.crawlWorker(jobs, results, &workerWG)
	}

	// queue holds URLs ready to be crawled.
	queue := []string{startURL}
	active := 0 // number of jobs currently being processed by workers

	// Coordinator loop: dispatch jobs when workers are idle and process results.
	for {
		// If there's nothing left to do, stop.
		if active == 0 && len(queue) == 0 {
			break
		}

		// Prepare a send only if there's work and an idle worker.
		var sendCh chan string
		var sendURL string
		if len(queue) > 0 && active < c.workers {
			sendCh = jobs
			sendURL = queue[0]
		}

		select {
		case sendCh <- sendURL:
			// Successfully dispatched a job.
			queue = queue[1:]
			active++

		case links := <-results:
			active--
			for _, link := range links {
				if c.markVisited(link) {
					queue = append(queue, link)
				}
			}
		}
	}

	// No more jobs to send; close channel so workers exit their range loop.
	close(jobs)
	workerWG.Wait()
}

// PrintSummary prints final crawl statistics.
func (c *Crawler) PrintSummary() {
	c.mu.Lock()
	defer c.mu.Unlock()

	fmt.Println()
	fmt.Println("================================")
	fmt.Println("          CRAWL SUMMARY")
	fmt.Println("================================")
	fmt.Println("Pages visited:", len(c.visited))
	fmt.Println()
	fmt.Println("HTTP status codes:")
	for code, count := range c.statusCodes {
		fmt.Printf("  %d → %d pages\n", code, count)
	}
	fmt.Println("================================")
}

func main() {
	startURL := "https://en.wikipedia.org/wiki/Web_scraping"
	workers := 100
	maxPages := 100

	crawler, err := NewCrawler(startURL, workers, maxPages)
	if err != nil {
		fmt.Println("Crawler error:", err)
		return
	}

	start := time.Now()
	crawler.Start(startURL)
	crawler.PrintSummary()
	fmt.Println("Time:", time.Since(start).Round(time.Millisecond))
}

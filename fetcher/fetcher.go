package fetcher

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
)

// Fetch navigates to the URL, waits for it to render, and returns the outer HTML.
func Fetch(url, waitFor string, timeout time.Duration) (string, error) {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Initialize headless browser options
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("disable-images", true),       // Optimization: don't load images
		chromedp.Flag("disable-animations", true),   // Optimization: skip animations
		chromedp.Flag("disable-popup-blocking", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	// Create a new browser context
	taskCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	var html string
	
	// Define the tasks to run in the browser
	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
	}

	// Wait for a specific selector if provided, otherwise wait for network idle (simulated by sleep for now)
	if waitFor != "" {
		tasks = append(tasks, chromedp.WaitVisible(waitFor, chromedp.ByQuery))
	} else {
		// A simple heuristic for waiting for SPAs to load if no specific selector is provided
		tasks = append(tasks, chromedp.Sleep(2*time.Second))
	}

	// Extract the outer HTML of the root element
	tasks = append(tasks, chromedp.OuterHTML("html", &html, chromedp.ByQuery))

	// Execute the tasks
	if err := chromedp.Run(taskCtx, tasks); err != nil {
		return "", err
	}

	return html, nil
}
package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/go-rod/rod"
)

// LoadPageWithTimeout loads a page with a specified timeout
func LoadPageWithTimeout(browser *rod.Browser, pageURL string, timeout time.Duration) (*rod.Page, error) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Create a page
	page := browser.MustPage(pageURL)

	// Wait for the page to load with context
	err := page.Timeout(timeout).MustWaitLoad()

	// Check if context was canceled due to timeout
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("timeout while loading page: %s", ctx.Err())
	default:
		// Continue normally if no timeout occurred
	}

	if err != nil {
		return nil, fmt.Errorf("error loading page: %w", err)
	}

	return page, nil
}

// LoadElementWithTimeout loads the page and attempts to find an element with a specified tag.
// It cancels the operation if the timeout is reached.
func LoadElementWithTimeout(page *rod.Page, tag string, timeout time.Duration) (*rod.Element, error) {

	// Attempt to find the element within the given timeout
	element, err := page.Element(tag)
	if err != nil {
		return nil, err
	}

	// Return the found element
	return element, nil
}

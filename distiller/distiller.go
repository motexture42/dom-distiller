package distiller

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

// Distill parses raw HTML and extracts semantic/actionable nodes into a tree structure.
func Distill(rawHTML string) (*Node, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(rawHTML))
	if err != nil {
		return nil, err
	}

	// 1. Cleanup Phase: remove non-essential and hidden tags
	doc.Find("script, style, noscript, svg, iframe, [style*='display: none'], [style*='display:none']").Remove()

	// Initialize ID counters
	var linkCounter, btnCounter, inputCounter int

	// The root document
	root := &Node{
		Type: "document",
	}

	// Traverse the body
	doc.Find("body").Each(func(i int, s *goquery.Selection) {
		for _, n := range s.Nodes {
			processNode(n, root, &linkCounter, &btnCounter, &inputCounter, "/html/body")
		}
	})

	return root, nil
}

// processNode recursively distills an HTML node into an internal Node
func processNode(n *html.Node, parent *Node, linkCount, btnCount, inputCount *int, currentXPath string) {
	if n.Type == html.TextNode {
		content := strings.TrimSpace(n.Data)
		if content != "" {
			parent.Children = append(parent.Children, &Node{
				Type:    "text",
				Content: content,
			})
		}
		return
	}

	if n.Type != html.ElementNode {
		return
	}

	tagName := strings.ToLower(n.Data)

	// Calculate XPath for this specific element
	elementIndex := 1
	for prev := n.PrevSibling; prev != nil; prev = prev.PrevSibling {
		if prev.Type == html.ElementNode && strings.ToLower(prev.Data) == tagName {
			elementIndex++
		}
	}
	myXPath := fmt.Sprintf("%s/%s[%d]", currentXPath, tagName, elementIndex)

	newNode := &Node{
		Type:       tagName,
		XPath:      myXPath,
		Attributes: make(map[string]string),
	}

	isActionable := false

	switch tagName {
	case "a":
		isActionable = true
		*linkCount++
		newNode.ActionID = fmt.Sprintf("LINK_%d", *linkCount)
		newNode.Type = "link"
		if href := getAttr(n, "href"); href != "" {
			newNode.Attributes["href"] = href
		}
	case "button":
		isActionable = true
		*btnCount++
		newNode.ActionID = fmt.Sprintf("BUTTON_%d", *btnCount)
		newNode.Type = "button"
	case "input":
		isActionable = true
		*inputCount++
		newNode.ActionID = fmt.Sprintf("INPUT_%d", *inputCount)
		newNode.Type = "input"
		if typ := getAttr(n, "type"); typ != "" {
			newNode.Attributes["type"] = typ
		}
		if ph := getAttr(n, "placeholder"); ph != "" {
			newNode.Attributes["placeholder"] = ph
		}
		if val := getAttr(n, "value"); val != "" {
			newNode.Attributes["value"] = val
		}
	case "h1", "h2", "h3", "h4", "h5", "h6":
		newNode.Type = "heading"
		newNode.Attributes["level"] = strings.TrimPrefix(tagName, "h")
	case "p", "div", "span", "ul", "ol", "li", "main", "article", "section":
		// These are structural or text containers, keep them if they have meaningful text or children
		newNode.Type = tagName
	default:
		// Ignore unknown tags but process children
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			processNode(c, parent, linkCount, btnCount, inputCount, myXPath)
		}
		return
	}

	// Process children for the current element
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		processNode(c, newNode, linkCount, btnCount, inputCount, myXPath)
	}

	// Post-processing: compress nodes if they are just wrappers
	compressNode(newNode)

	// Only append if it's actionable, has content, or has children
	if isActionable || len(newNode.Children) > 0 || strings.TrimSpace(newNode.Content) != "" {
		parent.Children = append(parent.Children, newNode)
	}
}

// compressNode simplifies the tree by flattening nested text or empty containers
func compressNode(n *Node) {
	if n.Type == "div" || n.Type == "span" || n.Type == "p" || n.Type == "li" {
		if len(n.Children) == 1 && n.Children[0].Type == "text" {
			n.Content = n.Children[0].Content
			n.Children = nil
		}
	}
	// Extract content for heading, link, button if it's purely text
	if n.Type == "heading" || n.Type == "link" || n.Type == "button" {
		var contentBuilder strings.Builder
		for _, c := range n.Children {
			if c.Type == "text" {
				contentBuilder.WriteString(c.Content + " ")
			}
		}
		if contentBuilder.Len() > 0 {
			n.Content = strings.TrimSpace(contentBuilder.String())
		}
		// Try to keep actionable children inside, but text is lifted to content.
	}
}

// getAttr gets the value of an attribute by key
func getAttr(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}
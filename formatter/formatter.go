package formatter

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/motextur3/dom-distiller/distiller"
)

// ToMarkdown converts the distilled Node tree into a token-optimized Markdown string.
func ToMarkdown(root *distiller.Node) string {
	var sb strings.Builder
	writeMarkdown(root, &sb, 0)
	return strings.TrimSpace(sb.String())
}

func writeMarkdown(n *distiller.Node, sb *strings.Builder, depth int) {
	indent := strings.Repeat("  ", depth)

	switch n.Type {
	case "document":
		for _, c := range n.Children {
			writeMarkdown(c, sb, depth)
		}
	case "heading":
		level := "#"
		if l, ok := n.Attributes["level"]; ok {
			switch l {
			case "1": level = "#"
			case "2": level = "##"
			case "3": level = "###"
			case "4": level = "####"
			case "5": level = "#####"
			case "6": level = "######"
			}
		}
		sb.WriteString(fmt.Sprintf("\n%s %s\n", level, n.Content))
	case "text":
		if n.Content != "" {
			sb.WriteString(fmt.Sprintf("%s%s\n", indent, n.Content))
		}
	case "link":
		href := n.Attributes["href"]
		// Condense to [LINK_ID: Content] (href)
		content := n.Content
		if content == "" {
			content = "Link"
		}
		sb.WriteString(fmt.Sprintf("%s* [%s: %s] (%s)\n", indent, n.ActionID, content, href))
	case "button":
		content := n.Content
		if content == "" {
			content = "Button"
		}
		sb.WriteString(fmt.Sprintf("%s* [%s: %s]\n", indent, n.ActionID, content))
	case "input":
		typ := n.Attributes["type"]
		if typ == "" {
			typ = "text"
		}
		ph := n.Attributes["placeholder"]
		val := n.Attributes["value"]
		sb.WriteString(fmt.Sprintf("%s* [%s: Input Type=%s Placeholder='%s' Value='%s']\n", indent, n.ActionID, typ, ph, val))
	default:
		// For containers (div, section, etc)
		if n.Content != "" {
			sb.WriteString(fmt.Sprintf("%s%s\n", indent, n.Content))
		}
		for _, c := range n.Children {
			writeMarkdown(c, sb, depth)
		}
	}
}

// ToJSON converts the distilled Node tree into compressed JSON.
func ToJSON(root *distiller.Node) string {
	b, err := json.Marshal(root)
	if err != nil {
		return "{}"
	}
	return string(b)
}
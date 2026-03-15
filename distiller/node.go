package distiller

// Node represents a distilled, semantic DOM element.
type Node struct {
	Type       string            `json:"type"`       // e.g., "text", "heading", "link", "button", "input"
	Content    string            `json:"content,omitempty"`    // Visible text content
	ActionID   string            `json:"action_id,omitempty"`   // Unique ID for interactive elements (e.g., "BTN_12")
	Attributes map[string]string `json:"attributes,omitempty"` // Relevant attributes like "href" or "placeholder"
	Children   []*Node           `json:"children,omitempty"`   // Nested elements
}
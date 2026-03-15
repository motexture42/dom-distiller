# dom-distiller

A blazing fast, stateless CLI utility optimized for allowing AI Agents and LLMs to interact with the web.

Traditional scraping tools (like returning raw HTML to an LLM) bloat the agent's context window with `<div>` soup, inline styles, and tracking scripts, burying the actionable elements.

`dom-distiller` solves this by decoupling the **"Reading"** phase from the **"Acting"** phase. It renders a page headlessly, gives the LLM a highly compressed "Agent View" in clean Markdown, and secretly saves a mapping of all the complex XPaths locally. When the agent is ready to click, it simply resolves its decision against the local map.

## Features

- **Token Efficiency:** Drops non-essential tokens (scripts, styles, SVGs) reducing massive webpages into lightweight, semantic maps.
- **Stateless Action Mapping (The Sidecar):** Automatically assigns unique IDs (e.g., `[BUTTON_4]`) to interactive elements. It saves the actual, heavy `XPath` mappings to a local JSON file so your LLM never has to look at them.
- **Headless Rendering:** Uses `chromedp` natively in Go. It executes JavaScript and waits for network idle, fully supporting complex SPAs (React, Vue, etc.).
- **Language Agnostic:** It's just a CLI. It acts as a perfect sidecar utility for agents written in Python, Node.js, or standard bash scripts.

## Installation

Ensure you have Go 1.21+ installed, then run:

```bash
go install github.com/motexture42/dom-distiller@latest
```

*(Or clone the repository and run `go build -o dom-distiller`)*

---

## The Agent Workflow (End-to-End Example)

This utility is designed to be called by your agent's code. Here is how an agent interacts with a webpage using `dom-distiller`:

### 1. The "Read" Phase

The Agent needs to see what is on the page. It calls `dom-distiller fetch`, requesting Markdown output and telling the utility to save the XPath map locally.

```bash
dom-distiller fetch https://news.ycombinator.com --format=markdown --save-map=/tmp/hn_map.json
```

**What the LLM sees (The Output):**
```markdown
# Hacker News

* [LINK_1: login] (https://news.ycombinator.com/login)
* [LINK_2: Show HN: My New App] (https://example.com)
* [BUTTON_1: upvote]
```
*Notice how clean the context window is. The LLM only sees actionable IDs and text.*

### 2. The "Decision" Phase

You prompt your LLM: *"Based on the page, you want to upvote the post. Which ID do you interact with?"*
The LLM replies: `BUTTON_1`.

### 3. The "Resolve" Phase

The Agent takes the LLM's decision (`BUTTON_1`) and asks `dom-distiller` for the actual DOM identifier using the map it saved earlier.

```bash
dom-distiller resolve /tmp/hn_map.json BUTTON_1
```

**Output:**
```text
/html/body/center/table/tbody/tr[3]/td/table/tbody/tr[1]/td[2]/center/a[1]
```

### 4. The "Act" Phase

The Agent immediately passes that exact XPath into its native browser automation tool (like Playwright or Puppeteer) to perform the action.

**Example Python Playwright integration:**
```python
import subprocess
from playwright.sync_api import sync_playwright

action_id = "BUTTON_1" # The ID the LLM chose

# 1. Resolve the XPath using the CLI
xpath = subprocess.check_output([
    "dom-distiller", "resolve", "/tmp/hn_map.json", action_id
]).decode().strip()

# 2. Click it in the actual browser session
with sync_playwright() as p:
    browser = p.chromium.launch(headless=False)
    page = browser.new_page()
    page.goto("https://news.ycombinator.com")
    
    # Execute the click flawlessly!
    page.locator(f"xpath={xpath}").click()
```

---

## Additional Usage

### JSON Output (Ideal for strict APIs)
```bash
dom-distiller fetch https://github.com/trending --format=json
```

### Wait for specific elements (SPAs)
If you are scraping a React/Vue app, you can tell the headless browser to wait for a specific CSS selector before distilling.
```bash
dom-distiller fetch https://my-react-app.com --wait-for=".dashboard-loaded"
```
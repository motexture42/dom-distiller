# dom-distiller

A blazing fast, CLI-based web scraper optimized specifically for AI Agents and LLMs. 

Unlike standard scraping tools that return bloated raw HTML (`<div>` soup, inline styles, massive SVG strings), `dom-distiller` renders a web page using a headless Chromium browser, strips out the noise, and returns a highly compressed, token-optimized **Agent View** in Markdown or JSON.

## Features

- **Built for Agents:** Automatically maps interactive elements (`<a>`, `<button>`, `<input>`) to unique semantic IDs (e.g., `[BUTTON_4]`) so your LLM can reliably reference them for interaction later.
- **Token Efficiency:** Drops tokens used by tracking scripts, `<style>` tags, base64 images, and deep `<div>` wrappers, easily reducing a 150,000-token web page down to a 3,000-token semantic map.
- **Action Resolution (No Token Bloat):** Calculates the exact XPath for every element under the hood, allowing you to save a "Sidecar Map" locally. Your agent can resolve `BUTTON_4` to an exact XPath instantly without cluttering its context window.
- **Headless Rendering:** Uses `chromedp` natively in Go. It executes JavaScript and waits for network idle, meaning it can successfully scrape complex Single Page Applications (React, Vue, etc.) without writing complex custom wait logic.
- **Language Agnostic:** It's just a CLI. You can call it from Python, Node.js, Rust, or standard bash scripts.

## Installation

Ensure you have Go 1.21+ installed, then run:

```bash
go install github.com/motexture42/dom-distiller@latest
```

## Usage

### 1. Fetching & Mapping (The Agent Reads)

Fetch a page, get the clean markdown, and save the hidden XPath mappings to a local file.

```bash
dom-distiller fetch https://news.ycombinator.com --format=markdown --save-map=/tmp/page_map.json
```

**Example Markdown Output:**
```markdown
# Hacker News

* [LINK_1: login] (https://news.ycombinator.com/login)
* [LINK_2: Show HN: My New App] (https://example.com)
* [BUTTON_1: upvote]
```

### 2. Resolving Actions (The Agent Interacts)

Once the LLM decides it wants to click `BUTTON_1`, it doesn't need to guess the XPath. It simply calls the `resolve` command using the map we just saved.

```bash
dom-distiller resolve /tmp/page_map.json BUTTON_1
```

**Output:**
```
/html/body/center/table/tbody/tr[3]/td/table/tbody/tr[1]/td[2]/center/a[1]
```

Your agent can now pass this exact string directly into Playwright or Puppeteer to perform the click!

```python
# Example Agent Logic (Python/Playwright)
action_id = "BUTTON_1" # Provided by LLM
xpath = subprocess.check_output(["dom-distiller", "resolve", "/tmp/page_map.json", action_id]).decode().strip()

page.locator(f"xpath={xpath}").click()
```

### JSON Output (Ideal for function calling/APIs)
```bash
dom-distiller fetch https://github.com/trending --format=json
```

### Wait for specific elements (SPAs)
```bash
dom-distiller fetch https://my-react-app.com --wait-for=".dashboard-loaded"
```
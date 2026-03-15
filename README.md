# dom-distiller

A blazing fast, CLI-based web scraper optimized specifically for AI Agents and LLMs. 

Unlike standard scraping tools that return bloated raw HTML (`<div>` soup, inline styles, massive SVG strings), `dom-distiller` renders a web page using a headless Chromium browser, strips out the noise, and returns a highly compressed, token-optimized **Agent View** in Markdown or JSON.

## Features

- **Built for Agents:** Automatically maps interactive elements (`<a>`, `<button>`, `<input>`) to unique semantic IDs (e.g., `[BUTTON_4]`) so your LLM can reliably reference them for interaction later.
- **Token Efficiency:** Drops tokens used by tracking scripts, `<style>` tags, base64 images, and deep `<div>` wrappers, easily reducing a 150,000-token web page down to a 3,000-token semantic map.
- **Headless Rendering:** Uses `chromedp` natively in Go. It executes JavaScript and waits for network idle, meaning it can successfully scrape complex Single Page Applications (React, Vue, etc.) without writing complex custom wait logic.
- **Language Agnostic:** It's just a CLI. You can call it from Python, Node.js, Rust, or standard bash scripts.

## Installation

Ensure you have Go 1.21+ installed, then run:

```bash
go install github.com/username/dom-distiller@latest
```

## Usage

### Basic Markdown Distillation
```bash
dom-distiller fetch https://news.ycombinator.com
```

### JSON Output (Ideal for function calling/APIs)
```bash
dom-distiller fetch https://github.com/trending --format=json
```

### Wait for specific elements (SPAs)
```bash
dom-distiller fetch https://my-react-app.com --wait-for=".dashboard-loaded"
```

## Example Output (Markdown Mode)

**Input:** A complex webpage with hundreds of nested divs, ads, and scripts.
**Output:**
```markdown
# Page Title

Main Content Article Text Here...

* [LINK_1: Login] (https://example.com/login)
* [LINK_2: Sign Up] (https://example.com/signup)
* [INPUT_1: Input Type=text Placeholder='Search...' Value='']
* [BUTTON_1: Submit]
```

## Using it with an Agent

When writing an agent prompt, simply provide the output of `dom-distiller` and instruct your agent to output the `ActionID` it wants to interact with.

*Prompt Example:*
> Based on the provided distilled DOM, which button should I click to submit the form? Reply only with the ActionID.

*Agent Output:*
> BUTTON_1
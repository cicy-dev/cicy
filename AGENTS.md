# AGENTS.md - CICY MCP Message Communication System

This document provides guidelines for AI agents working on the CICY codebase.

## Build, Lint, and Test Commands

### Installation
```bash
npm install
```

### Running the Server
```bash
npm start                          # Start server on port 13001
PORT=3000 npm start               # Custom port
npm run dev:server                 # Auto-restart with nodemon
```

### Running Clients
```bash
npm run client:cli                 # CLI client (recommended, tmux-compatible)
npm run dev:tui                    # TUI client with graphics
npm run dev                       # Run server and CLI client concurrently
```

### Testing
```bash
./test.sh                         # Full integration test suite
node test-image.js               # Image display test (TUI)
```

## Code Style Guidelines

### General Principles
- Keep functions focused and single-purpose
- Maximum function length: ~50 lines recommended
- Maximum file length: ~400 lines recommended
- Remove debug logs before committing
- No TODO comments left in code

### Naming Conventions
- Variables and functions: `camelCase` (e.g., `sendMessage`, `totalMessages`)
- Constants: `UPPER_SNAKE_CASE` for config values (e.g., `API_URL`, `PORT`)
- File names: `kebab-case` (e.g., `client-tui.js`, `test-image.js`)
- Descriptive names preferred over abbreviations (use `messageBox` not `msgBox`)

### Import Order
```javascript
// 1. Node.js core modules
const express = require('express');
const fs = require('fs');
const path = require('path');

// 2. External dependencies
const blessed = require('blessed');
const axios = require('axios');

// 3. Utilities (promisify, etc.)
const { exec } = require('child_process');
const { promisify } = require('util');
```

### Formatting
- Indentation: 4 spaces (not tabs)
- Semicolons: Always use semicolons
- Line length: ~100 characters max
- Blank lines: One blank line between function definitions
- Curly braces: Same-line style (`function() {`)

### Error Handling
- Always wrap async operations in try-catch blocks
- Provide user-friendly error messages
- Distinguish error types (e.g., `ECONNREFUSED` vs generic errors)
- Log errors with timestamps for debugging

```javascript
try {
    const response = await axios.post(url, data, { timeout: 5000 });
} catch (error) {
    if (error.code === 'ECONNREFUSED') {
        addMessage(`{red-fg} Error: Server not running{/red-fg}`);
    } else {
        addMessage(`{red-fg} Error: ${error.message}{/red-fg}`);
    }
}
```

### Async/Await Patterns
- Always use async/await over raw promises
- Set reasonable timeouts on HTTP requests
- Minimum loading animation duration: 500ms for UX
- Handle connection refused errors explicitly

### Comments
- Use Chinese comments for user-facing text and UX messages
- Use English for code logic comments
- Document function purpose at top of function
- No commented-out code left in production

### TUI/UI Development
- Always set `LANG=en_US.UTF-8` for blessed compatibility
- Use `{tags}` for colored output in blessed
- Handle `Ctrl+C` gracefully with double-press confirmation
- Provide visual feedback (loading animations, status updates)

### Git Workflow
- Commit messages in English
- Follow conventional commits format
- No force pushes to shared branches
- Always run `./test.sh` before committing

### Project Structure
```
/client-tui.js        # Main TUI client
/server.js            # Express MCP server
/test.sh              # Integration tests
/docs/                # Documentation
/example/             # Example code
/todos/               # Task tracking
```

### Protocol-Specific Guidelines
- MCP Protocol: Use version "2024-11-05"
- JSON-RPC 2.0 responses must include `jsonrpc`, `id`, and `result/error`
- Tools must have `name`, `description`, and `inputSchema`
- Error codes: -32600 (Invalid Request), -32601 (Method Not Found), -32602 (Invalid Params)

## Quick Reference

| Task | Command |
|------|---------|
| Start server | `npm start` |
| Dev mode | `npm run dev:server` |
| Run tests | `./test.sh` |
| TUI client | `npm run dev:tui` |
| CLI client | `npm run client:cli` |

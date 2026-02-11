# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-02-11

### Added
- MCP (Model Context Protocol) server implementation
- JSON-RPC 2.0 support
- 5 MCP tools: send_message, send_image, get_messages, get_images, clear_messages
- Node.js TUI client with blessed
- Remote client for connecting to remote servers
- Go TUI client with Bubble Tea framework
- Command support: /help, /quit, /clear, /list
- Loading animations and status indicators
- Message history display
- Tokyo Night color scheme
- Complete documentation (README, AGENTS.md)
- Integration tests

### Features
- **Server**
  - Express.js based MCP server
  - Port 13001 (configurable)
  - Text message support
  - Base64 image transfer
  - In-memory message storage

- **Node.js Client**
  - TUI interface with blessed
  - Masu style design
  - Loading animations
  - Completion time display
  - Error handling

- **Remote Client**
  - Connect to any CICY server
  - Real-time connection status
  - Auto-reconnect
  - Command support

- **Go Client**
  - High performance (~10MB memory)
  - Fast startup (<10ms)
  - Single binary deployment
  - Full command support
  - Message history (last 5)
  - Help interface

### Documentation
- Complete README with usage examples
- AGENTS.md for AI development
- Architecture documentation
- Development guidelines
- Testing guide

### Performance
- Go client: 10x faster startup than Node.js
- Go client: 5x less memory usage
- Async message handling
- Efficient rendering

## [Unreleased]

### Planned
- Image display in Go client
- Configuration file support
- More themes
- Plugin system
- WebSocket support
- Multi-language support

---

## Version History

- **1.0.0** (2026-02-11) - Initial release

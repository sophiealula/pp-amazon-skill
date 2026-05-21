// amazon-pp-mcp is a minimal stdio MCP shim that walks the amazon-pp-cli
// Cobra tree and exposes each leaf command as an MCP tool. This is the
// minimal viable surface — full MCP support (intents, code orchestration,
// HTTP transport) is intentionally out of scope for v0.1.
//
// Wire shape: stdio MCP. The handshake echoes back a tool catalog derived
// from `amazon-pp-cli agent-context --json` when available; if the parent
// CLI does not implement agent-context, the shim instead lists the headline
// commands (history, add, reorder-last, cart, doctor, checkout) so an MCP
// host has something to call.
//
// This file deliberately ships as a no-op stub: nanoclaw integrates via the
// printed-CLI binary directly. The MCP entry point exists so the project
// layout matches every other printing-press library entry.
package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stderr, "amazon-pp-mcp: stdio MCP shim is not implemented in v0.1.")
	fmt.Fprintln(os.Stderr, "Use the CLI directly: `amazon-pp-cli <command> --json` — nanoclaw consumes JSON outputs as the integration surface.")
	os.Exit(2)
}

package main

import (
	"log"
	"os"

	"github.com/hairglasses-studio/mcpkit/registry"
	"github.com/hairglasses-studio/sway-mcp-go/internal/sway"
)

func main() {
	if os.Getenv("WAYLAND_DISPLAY") == "" {
		log.Fatal("WAYLAND_DISPLAY not set — not running under Wayland")
	}

	s := registry.NewMCPServer("sway-mcp", "0.2.0")
	
	module := &sway.Module{}
	for _, td := range module.Tools() {
		registry.AddToolToServer(s, td.Tool, td.Handler)
	}

	if err := registry.ServeStdio(s); err != nil {
		log.Fatal(err)
	}
}

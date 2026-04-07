package main

import (
	"context"
	"log"
	"os"

	"github.com/hairglasses-studio/mcpkit/bridge/a2a"
	"github.com/hairglasses-studio/mcpkit/registry"
	"github.com/hairglasses-studio/sway-mcp-go/internal/sway"
)

func main() {
	if os.Getenv("WAYLAND_DISPLAY") == "" {
		log.Fatal("WAYLAND_DISPLAY not set — not running under Wayland")
	}

	// 1. Register tools
	reg := registry.NewToolRegistry()
	module := &sway.Module{}
	reg.RegisterModule(module)

	// 2. Configure and create the bridge
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	url := "http://localhost" + addr

	b, err := a2a.NewBridge(reg, a2a.BridgeConfig{
		Name:        "sway-agent",
		Description: "Sway/Wayland desktop management agent",
		URL:         url,
		Addr:        addr,
	})
	if err != nil {
		log.Fatal(err)
	}

	// 3. Start the bridge
	log.Printf("Sway A2A agent listening on %s", addr)
	if err := b.Start(context.Background()); err != nil {
		log.Fatal(err)
	}
}

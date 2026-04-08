package main

import (
	"log"

	"github.com/hairglasses-studio/mcpkit/prompts"
	"github.com/hairglasses-studio/mcpkit/registry"
	"github.com/hairglasses-studio/mcpkit/resources"
	"github.com/hairglasses-studio/sway-mcp-go/internal/sway"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	const version = "0.2.0"

	toolReg := registry.NewToolRegistry()
	resReg := resources.NewResourceRegistry()
	promptReg := prompts.NewPromptRegistry()

	toolReg.RegisterModule(&sway.Module{})
	toolReg.RegisterModule(&sway.ContractToolModule{
		ToolRegistry:     toolReg,
		ResourceRegistry: resReg,
		PromptRegistry:   promptReg,
		Version:          version,
	})
	resReg.RegisterModule(&sway.ContractResourceModule{
		ToolRegistry:   toolReg,
		PromptRegistry: promptReg,
		Version:        version,
	})
	promptReg.RegisterModule(&sway.ContractPromptModule{})

	s := registry.NewMCPServer(
		"sway-mcp",
		version,
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(false, true),
		server.WithPromptCapabilities(true),
	)
	toolReg.RegisterWithServer(s)
	resReg.RegisterWithServer(s)
	promptReg.RegisterWithServer(s)

	if err := registry.ServeStdio(s); err != nil {
		log.Fatal(err)
	}
}

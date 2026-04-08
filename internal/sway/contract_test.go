package sway

import (
	"testing"

	"github.com/hairglasses-studio/mcpkit/mcptest"
	"github.com/hairglasses-studio/mcpkit/prompts"
	"github.com/hairglasses-studio/mcpkit/registry"
	"github.com/hairglasses-studio/mcpkit/resources"
)

func newContractClient(t *testing.T) *mcptest.Client {
	t.Helper()

	toolReg := registry.NewToolRegistry()
	resReg := resources.NewResourceRegistry()
	promptReg := prompts.NewPromptRegistry()

	toolReg.RegisterModule(&Module{})
	toolReg.RegisterModule(&ContractToolModule{
		ToolRegistry:     toolReg,
		ResourceRegistry: resReg,
		PromptRegistry:   promptReg,
		Version:          "test",
	})
	resReg.RegisterModule(&ContractResourceModule{
		ToolRegistry:   toolReg,
		PromptRegistry: promptReg,
		Version:        "test",
	})
	promptReg.RegisterModule(&ContractPromptModule{})

	srv := mcptest.NewServer(t, toolReg)
	resReg.RegisterWithServer(srv.MCP)
	promptReg.RegisterWithServer(srv.MCP)

	return mcptest.NewClient(t, srv)
}

func TestContractCatalogAndHealth(t *testing.T) {
	client := newContractClient(t)

	catalog := client.CallTool("sway_tool_catalog", map[string]any{})
	mcptest.AssertNotError(t, catalog)

	var out CatalogOutput
	mcptest.AssertStructured(t, catalog, &out)
	if out.Count != 17 {
		t.Fatalf("catalog count = %d, want 17", out.Count)
	}

	health := client.CallTool("sway_server_health", map[string]any{})
	mcptest.AssertNotError(t, health)

	var status struct {
		Status        string `json:"status"`
		ToolCount     int    `json:"tool_count"`
		ResourceCount int    `json:"resource_count"`
		PromptCount   int    `json:"prompt_count"`
	}
	mcptest.AssertStructured(t, health, &status)
	if status.Status != "ok" {
		t.Fatalf("status = %q, want ok", status.Status)
	}
	if status.ToolCount != 17 {
		t.Fatalf("tool_count = %d, want 17", status.ToolCount)
	}
	if status.ResourceCount != 2 {
		t.Fatalf("resource_count = %d, want 2", status.ResourceCount)
	}
	if status.PromptCount != 1 {
		t.Fatalf("prompt_count = %d, want 1", status.PromptCount)
	}
}

func TestContractResourcesAndPrompt(t *testing.T) {
	client := newContractClient(t)

	overview := client.ReadResource("sway://server/overview")
	mcptest.AssertResourceContains(t, overview, "sway-mcp-go")
	mcptest.AssertResourceContains(t, overview, "sway_server_health")

	requirements := client.ReadResource("sway://runtime/requirements")
	mcptest.AssertResourceContains(t, requirements, "WAYLAND_DISPLAY")
	mcptest.AssertResourceContains(t, requirements, "ydotool")

	prompt := client.GetPrompt("sway_start_triage", map[string]string{"goal": "focus the browser"})
	mcptest.AssertPromptMessages(t, prompt, 1)
	mcptest.AssertPromptContains(t, prompt, "sway_server_health")
	mcptest.AssertPromptContains(t, prompt, "focus the browser")
}

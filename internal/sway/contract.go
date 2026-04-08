package sway

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"

	"github.com/hairglasses-studio/mcpkit/handler"
	"github.com/hairglasses-studio/mcpkit/prompts"
	"github.com/hairglasses-studio/mcpkit/registry"
	"github.com/hairglasses-studio/mcpkit/resources"
	"github.com/mark3labs/mcp-go/mcp"
)

type ContractToolModule struct {
	ToolRegistry     *registry.ToolRegistry
	ResourceRegistry *resources.ResourceRegistry
	PromptRegistry   *prompts.PromptRegistry
	Version          string
}

type ContractResourceModule struct {
	ToolRegistry   *registry.ToolRegistry
	PromptRegistry *prompts.PromptRegistry
	Version        string
}

type ContractPromptModule struct{}

type ToolCatalogInput struct {
	Category string `json:"category,omitempty" jsonschema:"description=Optional category filter: capture\\,input\\,window\\,clipboard\\,display\\,discovery"`
}

type ToolSearchInput struct {
	Query    string `json:"query" jsonschema:"required,description=Keyword query for tool names\\, descriptions\\, categories\\, or search terms"`
	Category string `json:"category,omitempty" jsonschema:"description=Optional category filter: capture\\,input\\,window\\,clipboard\\,display\\,discovery"`
}

type ToolSchemaInput struct {
	Name string `json:"name" jsonschema:"required,description=Exact sway tool name to inspect"`
}

type ToolCatalogEntry struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	IsWrite     bool     `json:"is_write"`
	SearchTerms []string `json:"search_terms,omitempty"`
}

type CatalogOutput struct {
	Tools []ToolCatalogEntry `json:"tools"`
	Count int                `json:"count"`
}

func (m *ContractToolModule) Name() string { return "contract" }

func (m *ContractToolModule) Description() string {
	return "Discovery and server health tools for sway-mcp-go"
}

func (m *ContractToolModule) Tools() []registry.ToolDefinition {
	catalog := handler.TypedHandler[ToolCatalogInput, CatalogOutput](
		"sway_tool_catalog",
		"List sway-mcp-go tools with optional category filtering. Prefer this when you need the complete discovery surface before choosing a specific action.",
		func(_ context.Context, input ToolCatalogInput) (CatalogOutput, error) {
			entries := m.catalogEntries(input.Category)
			return CatalogOutput{Tools: entries, Count: len(entries)}, nil
		},
	)
	catalog.Category = "discovery"
	catalog.SearchTerms = []string{"tool catalog", "tool list", "discovery", "wayland", "sway"}

	search := handler.TypedHandler[ToolSearchInput, CatalogOutput](
		"sway_tool_search",
		"Search sway-mcp-go tools by keyword. Use this when you know the task shape but not the exact tool name.",
		func(_ context.Context, input ToolSearchInput) (CatalogOutput, error) {
			query := strings.ToLower(strings.TrimSpace(input.Query))
			if query == "" {
				return CatalogOutput{}, fmt.Errorf("query is required")
			}

			entries := m.catalogEntries(input.Category)
			filtered := make([]ToolCatalogEntry, 0, len(entries))
			for _, entry := range entries {
				haystack := strings.ToLower(strings.Join([]string{
					entry.Name,
					entry.Description,
					entry.Category,
					strings.Join(entry.SearchTerms, " "),
				}, " "))
				if strings.Contains(haystack, query) {
					filtered = append(filtered, entry)
				}
			}
			return CatalogOutput{Tools: filtered, Count: len(filtered)}, nil
		},
	)
	search.Category = "discovery"
	search.SearchTerms = []string{"tool search", "find tool", "keyword search", "discovery"}

	schema := handler.TypedHandler[ToolSchemaInput, map[string]any](
		"sway_tool_schema",
		"Inspect one sway-mcp-go tool descriptor including schema, category, and write-safety hints.",
		func(_ context.Context, input ToolSchemaInput) (map[string]any, error) {
			td, ok := m.ToolRegistry.GetTool(input.Name)
			if !ok {
				return nil, fmt.Errorf("tool not found: %s", input.Name)
			}
			annotated := registry.ApplyToolMetadata(td, "", false)
			return map[string]any{
				"name":          td.Tool.Name,
				"description":   td.Tool.Description,
				"category":      td.Category,
				"is_write":      td.IsWrite,
				"input_schema":  td.Tool.InputSchema,
				"output_schema": annotated.Tool.OutputSchema,
				"search_terms":  td.SearchTerms,
			}, nil
		},
	)
	schema.Category = "discovery"
	schema.SearchTerms = []string{"tool schema", "tool descriptor", "input schema", "output schema"}

	statsTool := handler.TypedHandler[struct{}, map[string]any](
		"sway_tool_stats",
		"Show sway-mcp-go tool counts by category plus resource and prompt coverage.",
		func(_ context.Context, _ struct{}) (map[string]any, error) {
			stats := m.ToolRegistry.GetToolStats()
			resourceCount := 0
			promptCount := 0
			if m.ResourceRegistry != nil {
				resourceCount = m.ResourceRegistry.ResourceCount() + m.ResourceRegistry.TemplateCount()
			}
			if m.PromptRegistry != nil {
				promptCount = m.PromptRegistry.PromptCount()
			}
			return map[string]any{
				"tool_count":       stats.TotalTools,
				"module_count":     stats.ModuleCount,
				"resource_count":   resourceCount,
				"prompt_count":     promptCount,
				"by_category":      stats.ByCategory,
				"by_runtime_group": stats.ByRuntimeGroup,
				"write_tools":      stats.WriteToolsCount,
				"read_only_tools":  stats.ReadOnlyCount,
			}, nil
		},
	)
	statsTool.Category = "discovery"
	statsTool.SearchTerms = []string{"tool stats", "catalog stats", "tool counts", "contract stats"}

	health := handler.TypedHandler[struct{}, map[string]any](
		"sway_server_health",
		"Show sway-mcp-go runtime health, protocol coverage, and required desktop dependency status.",
		func(_ context.Context, _ struct{}) (map[string]any, error) {
			stats := m.ToolRegistry.GetToolStats()
			resourceCount := 0
			promptCount := 0
			if m.ResourceRegistry != nil {
				resourceCount = m.ResourceRegistry.ResourceCount() + m.ResourceRegistry.TemplateCount()
			}
			if m.PromptRegistry != nil {
				promptCount = m.PromptRegistry.PromptCount()
			}
			return map[string]any{
				"server":              "sway-mcp",
				"version":             m.Version,
				"status":              "ok",
				"go_version":          runtime.Version(),
				"wayland_display_set": os.Getenv("WAYLAND_DISPLAY") != "",
				"tool_count":          stats.TotalTools,
				"resource_count":      resourceCount,
				"prompt_count":        promptCount,
				"write_tools":         stats.WriteToolsCount,
				"read_only_tools":     stats.ReadOnlyCount,
				"required_binaries": map[string]any{
					"swaymsg":  binaryStatus("swaymsg"),
					"grim":     binaryStatus("grim"),
					"magick":   binaryStatus("magick"),
					"wtype":    binaryStatus("wtype"),
					"ydotool":  binaryStatus("ydotool"),
					"wl-copy":  binaryStatus("wl-copy"),
					"wl-paste": binaryStatus("wl-paste"),
				},
				"discovery_tools": []string{
					"sway_tool_catalog",
					"sway_tool_search",
					"sway_tool_schema",
					"sway_tool_stats",
					"sway_server_health",
				},
			}, nil
		},
	)
	health.Category = "discovery"
	health.SearchTerms = []string{"server health", "contract", "status", "dependencies", "wayland"}

	return []registry.ToolDefinition{catalog, search, schema, statsTool, health}
}

func (m *ContractToolModule) catalogEntries(category string) []ToolCatalogEntry {
	definitions := m.ToolRegistry.GetAllToolDefinitions()
	entries := make([]ToolCatalogEntry, 0, len(definitions))
	for _, td := range definitions {
		if category != "" && td.Category != category {
			continue
		}
		entries = append(entries, ToolCatalogEntry{
			Name:        td.Tool.Name,
			Description: td.Tool.Description,
			Category:    td.Category,
			IsWrite:     td.IsWrite,
			SearchTerms: td.SearchTerms,
		})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})
	return entries
}

func (m *ContractResourceModule) Name() string { return "server_context" }

func (m *ContractResourceModule) Description() string {
	return "Overview and runtime requirement resources for sway-mcp-go"
}

func (m *ContractResourceModule) Resources() []resources.ResourceDefinition {
	return []resources.ResourceDefinition{
		{
			Resource: mcp.NewResource(
				"sway://server/overview",
				"sway-mcp-go Overview",
				mcp.WithResourceDescription("Compact server card for discovery-first Sway and Wayland automation"),
				mcp.WithMIMEType("text/markdown"),
			),
			Handler: func(_ context.Context, _ mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
				toolCount := 0
				promptCount := 0
				if m.ToolRegistry != nil {
					toolCount = m.ToolRegistry.ToolCount()
				}
				if m.PromptRegistry != nil {
					promptCount = m.PromptRegistry.PromptCount()
				}
				return []mcp.ResourceContents{
					mcp.TextResourceContents{
						URI:      "sway://server/overview",
						MIMEType: "text/markdown",
						Text: strings.Join([]string{
							"# sway-mcp-go",
							"",
							"Sway and Wayland desktop automation retained as a repo-local surface even though the main workstation now targets Hyprland.",
							"",
							fmt.Sprintf("- Version: `%s`", m.Version),
							fmt.Sprintf("- Registered tools: `%d`", toolCount),
							fmt.Sprintf("- Registered prompt workflows: `%d`", promptCount),
							"",
							"1. Start with `sway_server_health` to confirm `WAYLAND_DISPLAY` and runtime binaries.",
							"2. Use `sway_tool_search` or `sway_tool_catalog` to choose the smallest matching tool.",
							"3. Prefer `sway_list_windows`, `sway_get_outputs`, or screenshots before write actions like click, type, key, move, or clipboard write.",
						}, "\n"),
					},
				}, nil
			},
			Category: "overview",
			Tags:     []string{"sway", "overview", "workflow"},
		},
		{
			Resource: mcp.NewResource(
				"sway://runtime/requirements",
				"sway-mcp-go Runtime Requirements",
				mcp.WithResourceDescription("Required binaries and environment expectations for sway-mcp-go"),
				mcp.WithMIMEType("text/markdown"),
			),
			Handler: func(_ context.Context, _ mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
				return []mcp.ResourceContents{
					mcp.TextResourceContents{
						URI:      "sway://runtime/requirements",
						MIMEType: "text/markdown",
						Text: strings.Join([]string{
							"# Runtime Requirements",
							"",
							"- `WAYLAND_DISPLAY` must be set for live desktop interaction.",
							"- Required binaries: `swaymsg`, `grim`, `magick`, `wtype`, `ydotool`, `wl-copy`, `wl-paste`.",
							"- Screenshot tools return base64-encoded PNG payloads sized for agent consumption.",
							"- Write actions are intentionally separate from discovery tools; inspect first whenever state is unclear.",
						}, "\n"),
					},
				}, nil
			},
			Category: "runtime",
			Tags:     []string{"runtime", "dependencies", "sway", "wayland"},
		},
	}
}

func (m *ContractResourceModule) Templates() []resources.TemplateDefinition { return nil }

func (m *ContractPromptModule) Name() string { return "server_prompts" }

func (m *ContractPromptModule) Description() string {
	return "Prompt workflows for discovery-first sway-mcp-go usage"
}

func (m *ContractPromptModule) Prompts() []prompts.PromptDefinition {
	return []prompts.PromptDefinition{
		{
			Prompt: mcp.NewPrompt(
				"sway_start_triage",
				mcp.WithPromptDescription("Start a Sway automation task with health and discovery checks before desktop mutation"),
				mcp.WithArgument("goal", mcp.ArgumentDescription("Optional short description of the Sway task to accomplish")),
			),
			Handler: func(_ context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
				goal := strings.TrimSpace(req.Params.Arguments["goal"])
				if goal == "" {
					goal = "the current Sway desktop task"
				}
				return mcp.NewGetPromptResult("Start sway-mcp-go triage", []mcp.PromptMessage{
					mcp.NewPromptMessage(mcp.RoleUser, mcp.NewTextContent(
						fmt.Sprintf("Triage %s using sway-mcp-go in a discovery-first way. Start with `sway_server_health`, then `sway_tool_search` or `sway_tool_catalog`, then use read-only inspection tools like `sway_list_windows`, `sway_get_outputs`, or screenshots before taking any write action.", goal),
					)),
				}), nil
			},
			Category: "workflow",
			Tags:     []string{"triage", "workflow", "discovery", "sway"},
		},
	}
}

func binaryStatus(name string) map[string]any {
	path, err := exec.LookPath(name)
	if err != nil {
		return map[string]any{
			"available": false,
			"path":      "",
		}
	}
	return map[string]any{
		"available": true,
		"path":      path,
	}
}

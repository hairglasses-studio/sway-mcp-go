# sway-mcp-go Roadmap

Last updated: 2026-04-08.

## Current State

Standalone Sway-focused MCP server retained separately while the main desktop control plane has moved to Hyprland.

- Tier: `standalone`
- Lifecycle: `maintenance-only`
- Language profile: `Go`
- Visibility / sensitivity: `PUBLIC` / `public`
<!-- whiteclaw-rollout:start -->
## Whiteclaw-Derived Overhaul (2026-04-08)

This tranche applies the highest-value whiteclaw findings that fit this repo's real surface: engineer briefs, bounded skills/runbooks, searchable provenance, scoped MCP packaging, and explicit verification ladders.

### Strategic Focus
- This repo is retained as a Sway-specific maintenance surface after the Hyprland migration.
- Use whiteclaw patterns to clarify when this repo is still the right target, what live-session requirements remain, and how to verify it without a running compositor.
- Keep the roadmap centered on contract trust and compatibility boundaries, not broad new autonomy work.

### Recommended Work
- [ ] [Boundary docs] Document the Hyprland-vs-Sway boundary and when this repo is still the correct control surface.
- [ ] [Smoke tests] Add smoke tests for no-WAYLAND discovery mode and for live-session helper requirements when Sway is present.
- [ ] [Contract snapshots] Snapshot the tool/resource/prompt schemas before further maintenance backports change them implicitly.
- [ ] [Compatibility] Record any session/runtime assumptions that future maintainers need to reproduce local behavior safely.

### Rationale Snapshot
- Tier / lifecycle: `standalone` / `maintenance-only`
- Language profile: `Go`
- Visibility / sensitivity: `PUBLIC` / `public`
- Surface baseline: AGENTS=yes, skills=yes, codex=yes, mcp_manifest=configured, ralph=no, roadmap=yes
- Whiteclaw transfers in scope: Sway-vs-Hyprland boundary, no-WAYLAND smoke tests, contract snapshots, runtime compatibility docs
- Live repo notes: AGENTS, skills, Codex config, configured .mcp.json, 9 workflow(s)

<!-- whiteclaw-rollout:end -->

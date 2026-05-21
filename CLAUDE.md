## graphify

This project has a knowledge graph at graphify-out/ with god nodes, community structure, and cross-file relationships. Current scope: `core/`, `api/`, `server/` (control plane).

Rules:
- ALWAYS read graphify-out/GRAPH_REPORT.md before reading any source files, running grep/glob searches, or answering codebase questions. The graph is your primary map of the codebase.
- IF graphify-out/wiki/index.md EXISTS, navigate it instead of reading raw files
- For cross-module "how does X relate to Y" questions, prefer `graphify query "<question>"`, `graphify path "<A>" "<B>"`, or `graphify explain "<concept>"` over grep — these traverse the graph's EXTRACTED + INFERRED edges instead of scanning files
- After modifying code, run `graphify update .` to keep the graph current (AST-only, no API cost).
- To extend scope (e.g. add `vehicle/` or `charger/` directories), re-run `/graphify` on the project — current index intentionally skips those to save tokens.

## gitnexus

This project is also indexed in **GitNexus** (broader symbol/relationship graph covering the whole repo, not just core/api/server). GitNexus is available via MCP tools (`mcp__gitnexus__*`) and is faster for raw symbol lookups across the full Go codebase.

Rules:
- Use `mcp__gitnexus__impact` before refactoring or renaming any non-trivial symbol — especially in `core/`, `api/`, `server/` — to see the blast radius (direct callers + indirect dependencies).
- Use `mcp__gitnexus__context` instead of grep when you need 360° view of one symbol (callers, callees, processes it participates in).
- Use `mcp__gitnexus__query` for execution-flow questions ("how does charging start", "what handles the OCPP message"); it returns ranked call chains, not file matches.
- Use `mcp__gitnexus__detect_changes` before committing — maps your uncommitted diff to affected processes.
- This repo is one of many indexed; ALWAYS pass `repo: "evcc-ng"` (or full path) to GitNexus calls.
- Reindex after large refactors: `npx gitnexus analyze --skip-agents-md` (the `--skip-agents-md` flag prevents auto-updates to AGENTS.md/CLAUDE.md that would otherwise pollute commits).

When graphify and gitnexus disagree, gitnexus is authoritative for symbol relationships in code (broader scope, AST-only); graphify is authoritative for concept-level relationships pulled from `core/optimizer.md`, `core/planner/planner.md`, `core/soc/README.md` and the OpenAPI spec.

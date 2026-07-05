# evcc-ng — community fork

Fork of `evcc-io/evcc` focused on fixing long-standing backlog issues and
adding top-voted features. Origin: `danusha2345/evcc-ng`. Sync `master` with
upstream periodically by **merge** (not rebase — branch protection blocks
force-push). Releases tagged `0.307.0-ngN`; latest is `0.307.0-ng10`.
Donations: a Boosty link in the README, nothing more.

Releases: pushing a tag runs `.github/workflows/release.yml`, which reuses the
fork's own `default.yml` (all `ubuntu-latest`) then publishes a **GitHub Release
only** via goreleaser with `.goreleaser-ng.yml` (cross-platform tar.gz/zip +
.deb + checksums, same asset set as upstream). Upstream's Docker/Cloudsmith-APT/
Homebrew/Hassio/fly.io steps are intentionally dropped — they need Depot
self-hosted runners and org secrets the fork lacks (that was why every earlier
`-ngN` release run sat queued for 24h and got auto-cancelled).

Upstream PRs are NOT welcome here — the maintainer rejected our submissions as
"AI cruft" and closed them all. Work only in this fork; act on any technical
feedback they give, but don't re-open the PR route.

## Changes on top of upstream

Bug fixes:
- #29922 race on `POST /api/vehicles/{id}/plan/strategy` — per-vehicle mutex (`server/http_vehicle_handler.go`)
- #29864 Peugeot/PSA `evcc token` bootstrap — skip validating instantiation (`cmd/token.go`)
- #30006 Hyundai 12V drain — passive vehicle test polling (`server/http_config_helper.go`)
- #28652 Tesla vehicle-api wakeup on `ErrAsleep` — only wakes a connected vehicle (`lp.connected()` guard) so a car that drove off isn't woken repeatedly / drained (`core/loadpoint.go`)
- #29682 solar forecast interval safety floors (`tariff/tariff.go`)
- #14418 EEBus/Elli loadpoint minCurrent respected (`core/loadpoint_effective.go`)
- #30541 battery boost starved below prioritySoc — `Site.sitePower()` zeroed the home-battery charge power below prioritySoc, stripping the boost's headroom so the loadpoint under-charged / stayed 1p. `sitePower()` now takes a per-loadpoint `boost` flag and skips the priority zeroing only when that loadpoint is boosting (`core/site.go`; `TestSitePowerBatteryBoost`). Non-boost / above-prioritySoc behaviour unchanged.

Features (backend + UI unless noted):
- #6144 vehicle SoC start/end in charge log (`core/session/`, Sessions UI)
- #14661 per-phase 1p/3p current limits (`core/loadpoint*`, settings modal)
- #21747 zero feed-in / PV curtailment on negative prices (`core/site*`, forecast view). Adapted to the upstream tristate curtail API (#30116): `shouldFeedInCurtail()` returns `*bool` (nil = not managed), merged with `circuitCurtailed()` in `core/site.go`, fed to `curtailPV(*bool)`.
- #14496 graceful startup when a charger/meter fails to init — opt-in via `--graceful-start` (default off keeps upstream failsafe + fatal banner so the config-fatals e2e test passes). When on: failing devices wrap as offline (`charger/wrapper.go`, `meter/wrapper.go`, `cmd/setup.go`) and the loadpoint card shows an offline badge.
- #19649 current forecasted solar power ("Now") in forecast view (`assets/js/components/Forecast/SolarDetails.vue`)
- #30068 don't erase external PV limits when not managing curtailment — now **subsumed by upstream tristate #30116**: a nil curtail value means "leave the inverter's registers alone" (e.g. a static Huawei 70% limit set in FusionSolar), so our earlier explicit gating is no longer needed.
- #21144 per-device enable/disable toggle in config UI — a disabled device stays in the config but is instantiated as a quiet offline stub (no I/O, no log noise, references still resolve so the site starts), independent of `--graceful-start`. Migration-safe `Disabled` Property (zero value = active, `util/config/config.go`); `NewDisabledWrapper` stubs in `charger/meter/vehicle wrapper.go` advertise `api.Offline` but **not** `api.Retryable` (manual re-enable only); skip wired in `cmd/setup.go`. Also suppresses the per-cycle ERROR spam from any offline charger in the loadpoint `Update` loop (`core/loadpoint.go`), which also quiets `--graceful-start` failure wrappers. UI: "Gerät aktiv" toggle in the device modal + "Inactive" badge on meter/vehicle cards (`assets/js/components/Config/`).
- #19650 loadpoint priority on the main-UI settings modal — frontend-only: the `POST /loadpoints/{id}/priority/{value}` route + `priority` state already existed, so it just adds a 0–10 priority select (shown when >1 loadpoint) to `assets/js/components/Loadpoints/SettingsModal.vue`, changeable on the fly without the config-UI restart prompt.
- #30705 per-vehicle fixed phase override — a vehicle that cannot do dynamic 1p/3p switching (e.g. MG4) can pin the loadpoint to a fixed phase count while it is the active vehicle; 0 = auto/dynamic (today's behaviour). New optional `api.PhaseConfigurer` interface (`api/api.go`) + `fixedphases` field on the vehicle `embed` (`vehicle/embed.go`). `lp.effectivePhasesConfigured()` (`core/loadpoint_phases.go`) lets the active vehicle's override win over `lp.phasesConfigured` on a phase-switching charger (read-time only, never mutates the loadpoint setting); the phase decision-path reads (`scalePhasesRequired`, `scalePhasesIfAvailable`, `pvScalePhases` scale-down gate, `maxActivePhases`, Update dispatch) route through it. Template-driven `fixedphases` choice (0/1/3, advanced) in the vehicle modal via `util/templates/defaults.yaml` (param names are lowercase — defaults merge via `ToLower`). `TestVehiclePhasesConfigured`. Override is ignored on non-switching chargers and must be ≤ the descriptive `phases` (max) field.

Conventions: Go tests via `make test` (or `CGO_ENABLED=0 go test ./...`),
UI via `npm run lint` (eslint+vue-tsc) + `npm test`. New i18n strings go in
both `i18n/en.json` and `i18n/de.json`. Frontend is Vue 3 Options API.
Runtime site/loadpoint settings persist via `settings.Set*` (not YAML) and
publish to the UI through the WebSocket state snapshot.

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

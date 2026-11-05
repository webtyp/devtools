---
PLAN: "feat: devtools — database diagram and API explorer for an app's dev build; nothing without -tags dev"
EXECUTOR: jules
REVIEWER: none
---

> Read AGENTS.md at the repo root before any change.
> **Finish by yourself:** when the delivery check passes, commit, push and open the pull request. Never ask whether to
> proceed or finalize — nobody answers, and a session that waits ends without your work.

# Plan — development tools live in one piece and never ship to production

## Tools

```bash
go install webtyp.com/devflow/cmd/gotest@latest
go get webtyp.com/appkit@v0.2.8 webtyp.com/migrate@v0.0.2 webtyp.com/erd@v0.0.16 webtyp.com/router@v0.4.17 \
       webtyp.com/layout@v0.3.99 webtyp.com/model@v0.2.13
```

Environment: `GOPRIVATE=webtyp.com/*,github.com/veltylabs/*`, `GOFLAGS=-mod=mod`, `go mod tidy` after `go get`.

## Why

Today `appkit` mounts the database diagram (`/_schema`) and the API route table (`/_routes`) on every server and
always shows «API explorer» and «Database» in the menu. Those are development tools; `appkit` is production only. This
repo owns them, behind the `dev` build tag (the tag `webtyp dev` already passes to the server and the browser
builds). The generated mains of every app call it unconditionally; in production the calls return nothing:

```go
// web/server/main.go
func main() { if err := appkit.Server(config.App, devtools.Server(config.App)...); err != nil { log.Fatal(err) } }
// web/browser/main.go
func main() { appkit.Browser(config.App, devtools.Menu()) }
```

`appkit` (a parallel plan) changes to `Server(app App, mounts ...func(r router.Router, policy model.PolicyDescriber) error) error`
and `Browser(app App, tools appmenu.Tools)`. This repo does **not** call those; it only returns values of those types.

## Ecosystem search

What exists and is reused (no copy): `erd.ReadSteps`, `erd.Mount`, `erd.SchemaPath`, `erd.ResourceSchema`;
`router.MountIntrospection`, `router.IntrospectionPath`, `router.ResourceIntrospection`; `migrate.Steps`;
`apiexplorer.FromServer`, `schemaview.FromServer`, `appmenu.Tools`. The current wiring to move is
`appkit/server.go:53-61` and `appkit/run_browser.go:139` (`_temp/appkit/tests/server_test.go` shows how those are
tested today with `webtyp.com/router/mock`).

## Design gate

1. Prior art (unverified): Rails `web-console` only in the development group; Django `debug_toolbar` only in dev
   settings; Spring Boot DevTools excluded from the packaged jar — the tool is left out of the deliverable.
2. Names: `devtools.Server(app)`, `devtools.Menu()` — what each side gets.
3. Ledger: +2 exported functions; ways to add a dev tool: 1 (this repo).
4. Belongs here: one concern, development tools. The tag sits inside the library (same pattern as
   `webtyp.com/auth/authority` `devlogin.go` / `devlogin_off.go`), so apps have no tagged files.
5. Deletes: the gonew scaffold (`devtools.go`: `type Devtools`, `New`). `appkit`'s copy is deleted by its own plan.

## APIs used

| Symbol | Signature | Tag |
|---|---|---|
| `erd.ReadSteps` | `func ReadSteps(steps ...ddl.Step) (erd.Diagram, error)` | erd v0.0.16 |
| `erd.Mount` | `func Mount(r router.Router, path string, d erd.Diagram) router.Route` | erd v0.0.16 |
| `erd.SchemaPath`, `erd.ResourceSchema` | `const = "/_schema"`, `const model.Resource = "db_schema"` | erd v0.0.16 |
| `router.MountIntrospection` | `func MountIntrospection(r router.Router, path string, policy model.PolicyDescriber) router.Route` | router v0.4.17 |
| `router.IntrospectionPath`, `router.ResourceIntrospection` | `"/_routes"`, `model.Resource = "api_explorer"` | router v0.4.17 |
| `Route.Requires` | `.Requires(res model.Resource, a model.Action)` as used in `appkit/server.go:53` | router v0.4.17 |
| `migrate.Steps` | `func Steps(app appkit.App) []ddl.Step` | migrate v0.0.2 |
| `appmenu.Tools` (`webtyp.com/layout/app`) | `struct{ API *apiexplorer.Config; Database *schemaview.Source }` | layout v0.3.99 |
| `apiexplorer.FromServer` / `schemaview.FromServer` | `func FromServer() *apiexplorer.Config` / `*schemaview.Source` | layout v0.3.99 |

If a signature in the published tag differs from this table, stop and name the symbol in the PR; never adapt it by guessing.

## Stage 1 — server side

Delete `devtools.go` (scaffold). `server.go`:
```go
//go:build dev && !wasm

// Package devtools gives an app's development build its tools: the database diagram and the API route table on the
// server, and their screens in the menu. Built without -tags dev, every function returns nothing.
package devtools

import (
	"webtyp.com/appkit"
	"webtyp.com/erd"
	"webtyp.com/migrate"
	"webtyp.com/model"
	"webtyp.com/router"
)

// Server is what the development server mounts: the route table at router.IntrospectionPath (api_explorer:read)
// and the database diagram of the app's migrations at erd.SchemaPath (db_schema:read).
//
// Example: appkit.Server(config.App, devtools.Server(config.App)...)
func Server(app appkit.App) []func(r router.Router, policy model.PolicyDescriber) error {
	return []func(router.Router, model.PolicyDescriber) error{
		func(r router.Router, policy model.PolicyDescriber) error {
			router.MountIntrospection(r, router.IntrospectionPath, policy).Requires(router.ResourceIntrospection, model.Read)
			return nil
		},
		func(r router.Router, _ model.PolicyDescriber) error {
			steps := migrate.Steps(app)
			if len(steps) == 0 {
				return nil
			}
			diagram, err := erd.ReadSteps(steps...)
			if err != nil {
				return err
			}
			erd.Mount(r, erd.SchemaPath, diagram).Requires(erd.ResourceSchema, model.Read)
			return nil
		},
	}
}
```
`server_off.go`:
```go
//go:build !dev && !wasm

package devtools

import (
	"webtyp.com/appkit"
	"webtyp.com/model"
	"webtyp.com/router"
)

// Server returns nothing in a build without -tags dev: production mounts no development tool.
func Server(appkit.App) []func(r router.Router, policy model.PolicyDescriber) error { return nil }
```

## Stage 2 — browser side

`menu.go` (`//go:build dev`): `func Menu() appmenu.Tools { return appmenu.Tools{API: apiexplorer.FromServer(), Database: schemaview.FromServer()} }`
with doc «Menu is the «Developer tools» group of the menu: API explorer and Database. Example:
appkit.Browser(config.App, devtools.Menu())».
`menu_off.go` (`//go:build !dev`): `func Menu() appmenu.Tools { return appmenu.Tools{} }` — «no tools: the menu leaves
the group out». If `webtyp.com/layout/app` does not compile without the wasm tag, put `&& wasm` / `&& wasm` on both
menu files and say so in the PR.

## Stage 3 — tests (`tests/`, package `tests`)

Use `webtyp.com/router/mock` as `_temp/appkit/tests/server_test.go` does (`mock.Router`, `Configure(mock.Config{Authn,
Authorize})`, `Invoke`, `mock.Context`, the `get` helper, `shopStep`, `allowSchema`). Mount with
`for _, m := range devtools.Server(app) { m(r, nil) }`, then `r.Verify()`.

`tests/dev_test.go` (`//go:build dev && !wasm`), `app := appkit.App{Modules: []appkit.Module{{Name: "shop", Server: appkit.ModuleServer{Migrate: shopStep.Run}}}}`:

| Test | Setup → call | Expected |
|---|---|---|
| `TestServerGivesTwoMounts` | `devtools.Server(app)` | `len == 2` |
| `TestRoutesNeedTheirPermission` | authorize only `fake` → `get(r, router.IntrospectionPath)` | 403 |
| `TestRoutesOpenWithPermission` | authorize `router.ResourceIntrospection` read → same | 200 |
| `TestSchemaServedWithPermission` | authorize `erd.ResourceSchema` read → `get(r, erd.SchemaPath)` | 200; body decodes to `erd.Report` whose SVG contains `id="erd-table-invoice"` |
| `TestNoSchemaWithoutSteps` | `appkit.App{}` → mount → `get(r, erd.SchemaPath)` | 404 |
| `TestMenuHasTools` | `devtools.Menu()` | `API != nil && Database != nil` |

`tests/off_test.go` (`//go:build !dev && !wasm`):

| Test | Setup → call | Expected |
|---|---|---|
| `TestServerEmptyInProduction` | `devtools.Server(app)` | `nil` |
| `TestMenuEmptyInProduction` | `devtools.Menu()` | `API == nil && Database == nil` |

## Stage 4 — docs and cleanup

`README.md` (≤150 lines, 4 sections; Without/With: today `appkit` ships the tools to production / with devtools only
the dev build has them), Step-by-Step = the two generated mains above, Documentation Map table. `docs/ARCHITECTURE.md`
(4 sections): the tag pairs, why no `init`, what each mount guards. Delete `_temp/` (`test ! -e _temp`).

## Delivery check

1. `gotest` passes (it also runs the `-tags dev` pass). A PR that does not compile, or whose tests could not run, is a
   failed delivery.
2. `grep -rn "func init\|type Devtools" --include=*.go .` → empty; `test ! -e _temp`.

Commit after every stage. Never run `git reset`, `git restore`, `git checkout -- <file>` or `git stash` on your own
work. Stuck 30 minutes: open the PR anyway listing what fails.

## Stages

| Stage | Files | Done when |
|---|---|---|
| 1 | `devtools.go` (deleted), `server.go`, `server_off.go`, `go.mod` | compiles with and without `-tags dev` |
| 2 | `menu.go`, `menu_off.go` | compiles |
| 3 | `tests/dev_test.go`, `tests/off_test.go` | tests pass |
| 4 | `README.md`, `docs/ARCHITECTURE.md`, `_temp/` (deleted) | done |

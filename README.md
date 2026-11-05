# devtools
<img src="docs/img/badges.svg">

Development tools of a webtyp app (database diagram, API explorer).

## Without devtools

Today, `appkit` mounts the database diagram and API route table on every server, and always shows "API explorer" and "Database" in the menu. Since `appkit` is production-only, those screens end up shipping to production.

## With devtools

These tools now live in this repository behind the `dev` build tag. When built with `-tags dev`, they are mounted to the app. When built without `-tags dev` (like in production), they return nothing and are omitted entirely from the binary.

## Step-by-Step

Generated mains of every app call the tools unconditionally, but in production they return nothing:

```go
// web/server/main.go
func main() {
    if err := appkit.Server(config.App, devtools.Server(config.Migrations...)...); err != nil {
        log.Fatal(err)
    }
}
```

```go
// web/browser/main.go
func main() {
    appkit.Browser(config.App, devtools.Menu())
}
```

## Documentation Map

| Resource | Description |
|---|---|
| `docs/ARCHITECTURE.md` | Architecture decisions, tag pairs, missing `init`, and security of mounts |

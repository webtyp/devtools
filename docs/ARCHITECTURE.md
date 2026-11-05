# Architecture

## Tag Pairs

Every tool in this library is behind the `dev` build tag. This ensures that only development builds get development tools, and production builds carry none of them. The files come in pairs:
- `server.go` (`dev && !wasm`) and `server_off.go` (`!dev && !wasm`)
- `menu.go` (`dev`) and `menu_off.go` (`!dev`)

The same exported names exist on both sides of each pair.

## Why No Init

We do not use `init()` functions for side effects or automatic registration. Everything is mounted explicitly via function arguments. This explicit passing makes sure that development tools aren't automatically injected unless strictly configured in the caller.

## Mount Guards

The mounts for the development tools guard access with explicit requirements:
- The API route table (`router.IntrospectionPath`) requires `api_explorer:read`.
- The database diagram (`erd.SchemaPath`) requires `db_schema:read`.

These permission checks ensure that even in development mode, only users explicitly granted those permissions can access the introspective screens.

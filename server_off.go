//go:build !dev && !wasm

package devtools

import (
	"webtyp.com/ddl"
	"webtyp.com/model"
	"webtyp.com/router"
)

// Server returns nothing in a build without -tags dev: production mounts no development tool.
func Server(ddl.Schema) []func(r router.Router, policy model.PolicyDescriber) error { return nil }

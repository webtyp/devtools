//go:build !dev && !wasm

package devtools

import (
	"webtyp.com/appkit"
	"webtyp.com/model"
	"webtyp.com/router"
)

// Server returns nothing in a build without -tags dev: production mounts no development tool.
func Server(appkit.App) []func(r router.Router, policy model.PolicyDescriber) error { return nil }

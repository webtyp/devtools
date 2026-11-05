//go:build dev && !wasm

// Package devtools gives an app's development build its tools: the database diagram and the API route table on the
// server, and their screens in the menu. Built without -tags dev, every function returns nothing.
package devtools

import (
	"webtyp.com/ddl"
	"webtyp.com/erd"
	"webtyp.com/model"
	"webtyp.com/router"
)

// Server is what the development server mounts: the route table at router.IntrospectionPath (api_explorer:read)
// and the database diagram of the given migrations at erd.SchemaPath (db_schema:read).
//
// Example: appkit.Server(config.App, devtools.Server(config.Migrations...)...)
func Server(steps ...ddl.Step) []func(r router.Router, policy model.PolicyDescriber) error {
	return []func(router.Router, model.PolicyDescriber) error{
		func(r router.Router, policy model.PolicyDescriber) error {
			router.MountIntrospection(r, router.IntrospectionPath, policy).Requires(router.ResourceIntrospection, model.Read)
			return nil
		},
		func(r router.Router, _ model.PolicyDescriber) error {
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

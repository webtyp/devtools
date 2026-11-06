//go:build !dev && !wasm

package tests

import (
	"testing"

	"webtyp.com/ddl"
	"webtyp.com/devtools"
)

func TestServerEmptyInProduction(t *testing.T) {
	mounts := devtools.Server(ddl.Schema{})
	if mounts != nil {
		t.Fatalf("devtools.Server(ddl.Schema{}) in production gave %v, want nil", mounts)
	}
}

func TestMenuEmptyInProduction(t *testing.T) {
	tools := devtools.Menu()
	if tools.API != nil || tools.Database != nil {
		t.Fatalf("devtools.Menu() in production gave API=%v, Database=%v, want nil for both", tools.API, tools.Database)
	}
}

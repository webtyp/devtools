//go:build !dev && !wasm

package tests

import (
	"testing"
	"webtyp.com/appkit"
	"webtyp.com/devtools"
)

func TestServerEmptyInProduction(t *testing.T) {
	mounts := devtools.Server(appkit.App{})
	if mounts != nil {
		t.Fatalf("devtools.Server(app) in production gave %v, want nil", mounts)
	}
}

func TestMenuEmptyInProduction(t *testing.T) {
	tools := devtools.Menu()
	if tools.API != nil || tools.Database != nil {
		t.Fatalf("devtools.Menu() in production gave API=%v, Database=%v, want nil for both", tools.API, tools.Database)
	}
}

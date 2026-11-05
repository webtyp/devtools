//go:build dev && !wasm

package tests

import (
	"strings"
	"testing"
	"webtyp.com/ddl"
	"webtyp.com/devtools"
	"webtyp.com/erd"
	"webtyp.com/json"
	"webtyp.com/model"
	"webtyp.com/router"
	"webtyp.com/router/mock"
)

type table struct {
	name string
	ref  string
}

func (t *table) ModelName() string { return t.name }
func (t *table) Schema() []model.Field {
	f := []model.Field{{Name: "id", Type: model.Text(), DB: &model.FieldDB{PK: true}}}
	if t.ref != "" {
		f = append(f, model.Field{Name: "ref_id", Type: model.Text()})
	}
	return f
}
func (t *table) SchemaExt() []model.FieldExt {
	if t.ref == "" {
		return nil
	}
	return []model.FieldExt{{Field: t.Schema()[1], Ref: t.ref}}
}
func (t *table) Pointers() []any                  { return nil }
func (t *table) IsNil() bool                      { return t == nil }
func (t *table) EncodeFields(_ model.FieldWriter) {}
func (t *table) DecodeFields(_ model.FieldReader) {}

var shopStep = ddl.Step{Name: "shop", Run: func(conn ddl.Execer, c ddl.Compiler) error {
	d := ddl.New(conn, c)
	if err := d.Sync(&table{name: "customer"}); err != nil {
		return err
	}
	return d.Sync(&table{name: "invoice", ref: "tenant"})
}}

func fakeAuthn(userID string) router.Middleware {
	return func(next router.HandlerFunc) router.HandlerFunc {
		return func(ctx router.Context) {
			ctx.SetUserID(userID)
			next(ctx)
		}
	}
}

func get(r *mock.Router, path string) (int, string) {
	ctx := &mock.Context{
		InMethod: "GET",
		InPath:   path,
	}
	r.Invoke("GET", path, ctx)
	status := ctx.Status
	if status == 0 {
		status = 200
	}
	return status, string(ctx.ResponseBody())
}


func TestServerGivesTwoMounts(t *testing.T) {
	mounts := devtools.Server(shopStep)
	if len(mounts) != 2 {
		t.Fatalf("devtools.Server(app) gave %d mounts, want 2", len(mounts))
	}
}

func TestRoutesNeedTheirPermission(t *testing.T) {
	mounts := devtools.Server(shopStep)
	r := &mock.Router{}
	r.Configure(mock.Config{
		Authn: fakeAuthn("test-user"),
		Authorize: func(userID string, res model.Resource, a model.Action) bool {
			return res == "fake"
		},
	})

	for _, m := range mounts {
		m(r, nil)
	}

	if err := r.Verify(); err != nil {
		t.Fatalf("r.Verify: %v", err)
	}

	status, _ := get(r, router.IntrospectionPath)
	if status != 403 {
		t.Fatalf("GET %s without permission: want 403, got %d", router.IntrospectionPath, status)
	}
}

func TestRoutesOpenWithPermission(t *testing.T) {
	mounts := devtools.Server(shopStep)
	r := &mock.Router{}
	r.Configure(mock.Config{
		Authn: fakeAuthn("test-user"),
		Authorize: func(userID string, res model.Resource, a model.Action) bool {
			return res == router.ResourceIntrospection && a == model.Read
		},
	})

	for _, m := range mounts {
		m(r, nil)
	}
	if err := r.Verify(); err != nil {
		t.Fatalf("r.Verify: %v", err)
	}

	status, _ := get(r, router.IntrospectionPath)
	if status != 200 {
		t.Fatalf("GET %s with permission: want 200, got %d", router.IntrospectionPath, status)
	}
}

func TestSchemaServedWithPermission(t *testing.T) {
	mounts := devtools.Server(shopStep)
	r := &mock.Router{}
	r.Configure(mock.Config{
		Authn: fakeAuthn("test-user"),
		Authorize: func(userID string, res model.Resource, a model.Action) bool {
			return res == erd.ResourceSchema && a == model.Read
		},
	})

	for _, m := range mounts {
		m(r, nil)
	}
	if err := r.Verify(); err != nil {
		t.Fatalf("r.Verify: %v", err)
	}

	status, body := get(r, erd.SchemaPath)
	if status != 200 {
		t.Fatalf("GET %s with permission: want 200, got %d", erd.SchemaPath, status)
	}

	var rep erd.Report
	if err := json.Decode([]byte(body), &rep); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(rep.SVG, `id="erd-table-invoice"`) {
		t.Fatalf("diagram SVG missing invoice table, got: %s", rep.SVG)
	}
}

func TestNoSchemaWithoutSteps(t *testing.T) {
	mounts := devtools.Server()
	r := &mock.Router{}
	r.Configure(mock.Config{
		Authn: fakeAuthn("test-user"),
		Authorize: func(userID string, res model.Resource, a model.Action) bool {
			return true
		},
	})

	for _, m := range mounts {
		m(r, nil)
	}
	if err := r.Verify(); err != nil {
		t.Fatalf("r.Verify: %v", err)
	}

	status, _ := get(r, erd.SchemaPath)
	if status != 404 {
		t.Fatalf("GET %s on app without migrations: want 404, got %d", erd.SchemaPath, status)
	}
}

func TestMenuHasTools(t *testing.T) {
	tools := devtools.Menu()
	if tools.API == nil {
		t.Errorf("API tool is nil")
	}
	if tools.Database == nil {
		t.Errorf("Database tool is nil")
	}
}

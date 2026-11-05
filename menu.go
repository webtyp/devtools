//go:build dev

package devtools

import (
	"webtyp.com/layout/apiexplorer"
	appmenu "webtyp.com/layout/app"
	"webtyp.com/layout/schemaview"
)

// Menu is the "Developer tools" group of the menu: API explorer and Database.
//
// Example: appkit.Browser(config.App, devtools.Menu())
func Menu() appmenu.Tools {
	return appmenu.Tools{
		API:      apiexplorer.FromServer(),
		Database: schemaview.FromServer(),
	}
}

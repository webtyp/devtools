//go:build dev

package devtools

import (
	"webtyp.com/layout/apiexplorer"
	appmenu "webtyp.com/layout/app"
	"webtyp.com/layout/schemaview"
)

// Menu is the "Developer tools" group of the menu: API explorer and Database.
//
// Example: browser.Run(config.App, screens.Browser, layoutapp.New(devtools.Menu()))
func Menu() appmenu.Tools {
	return appmenu.Tools{
		API:      apiexplorer.FromServer(),
		Database: schemaview.FromServer(),
		ViewAs:   &appmenu.ViewAsConfig{},
	}
}

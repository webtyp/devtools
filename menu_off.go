//go:build !dev

package devtools

import appmenu "webtyp.com/layout/app"

// Menu returns no tools: the menu leaves the group out in a build without -tags dev.
func Menu() appmenu.Tools { return appmenu.Tools{} }

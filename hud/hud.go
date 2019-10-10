package hud

import "cavernal.com/lib/g3n/engine/gui"

type IHUD interface {
	Panel() *gui.Panel
	Dispose()
	PreRender()
}
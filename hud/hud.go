package hud

import "github.com/lquesada/cavernal/lib/g3n/engine/gui"

type IHUD interface {
	Panel() *gui.Panel
	Dispose()
	PreRender()
}
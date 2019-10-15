package model

import (
	"github.com/lquesada/cavernal/lib/g3n/engine/geometry"
	"github.com/lquesada/cavernal/lib/g3n/engine/material"
	"github.com/lquesada/cavernal/lib/g3n/engine/graphic"
	"github.com/lquesada/cavernal/lib/g3n/engine/math32"
	"github.com/lquesada/cavernal/lib/g3n/engine/core"
)

func NewShadow() INode {
	shadowGeom := geometry.NewCircle(1, 15)
    shadowMat := material.NewPhong(&math32.Color{0, 0, 0})
    shadow := graphic.NewMesh(shadowGeom, shadowMat)
    node := core.NewNode()
    node.Add(shadow)
    return &modelNode{node: node}
}
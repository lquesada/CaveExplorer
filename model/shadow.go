package model

import (
	"cavernal.com/lib/g3n/engine/geometry"
	"cavernal.com/lib/g3n/engine/material"
	"cavernal.com/lib/g3n/engine/graphic"
	"cavernal.com/lib/g3n/engine/math32"
	"cavernal.com/lib/g3n/engine/core"
)

func NewShadow() INode {
	shadowGeom := geometry.NewCircle(1, 15)
    shadowMat := material.NewPhong(&math32.Color{0, 0, 0})
    shadow := graphic.NewMesh(shadowGeom, shadowMat)
    node := core.NewNode()
    node.Add(shadow)
    return &modelNode{node: node}
}
package model

import "github.com/lquesada/cavernal/lib/g3n/engine/graphic"
import "github.com/lquesada/cavernal/lib/g3n/engine/material"
import        "github.com/lquesada/cavernal/lib/g3n/engine/core"
import         "github.com/lquesada/cavernal/lib/g3n/engine/math32"

type INode interface {
	G3NNode() *core.Node
	G3NNodeMask(r, g, b, a float32) *core.Node
	Transform(ts ...*Transform) INode
	RGBA(color *math32.Color4) INode
}

type node struct {
	nodes []INode
	transform *Transform
	rgba *math32.Color4
}

func NewNode(ns ...INode) INode {
	return &node{nodes: ns}
}

func (n *node) G3NNode() *core.Node {
	return n.G3NNodeMask(1, 1, 1, 1)
}

func (n *node) G3NNodeMask(r, g, b, a float32) *core.Node {
	node := core.NewNode()
	for _, i := range n.nodes {
		if n.rgba != nil {
			node.Add(i.G3NNodeMask(r*n.rgba.R, g*n.rgba.G, b*n.rgba.B, a*n.rgba.A))
		} else {
			node.Add(i.G3NNodeMask(r, g, b, a))
		}			
	}
	if n.transform != nil {
		node = n.transform.Apply(node)
	}
	return node
}

func (n *node) Add(is ...INode) {
	for _, i := range is {
		n.nodes = append(n.nodes, i)
	}
}

func (n *node) Transform(ts ...*Transform) INode {
	return applyTransforms(n, ts)
}

func (n *node) RGBA(color *math32.Color4) INode {
	return applyRGBA(n, color)
}

// --

type modelNode struct {
	node *core.Node
}

func (n *modelNode) G3NNode() *core.Node {
	return n.G3NNodeMask(1, 1, 1, 1)
}

func (n *modelNode) G3NNodeMask(r, g, b, a float32) *core.Node {
	for _, c := range n.node.Children() {
		mesh, ok := c.(*graphic.Mesh)
		if !ok {
			continue
		}
		for _, mat := range mesh.Materials() {
			m, ok := mat.IMaterial().(*material.Phong)
			if !ok {
				continue
			}
			m.SetAmbientColor(&math32.Color{r, g, b})
			m.SetOpacity(a)
			if a < 1 {
				m.SetTransparent(true)
			} else {
				m.SetTransparent(false)
			}

		}
	}
	return n.node
}

func (n *modelNode) Transform(ts ...*Transform) INode {
	return applyTransforms(n, ts)
}

func (n *modelNode) RGBA(color *math32.Color4) INode {
	return applyRGBA(n, color)
}

// --

func applyTransforms(n INode, ts []*Transform) INode {
	for _, t := range ts {
		n = &node{
			nodes: []INode{n},
			transform: t,
		}
	}
	return n
}

func applyRGBA(n INode, color *math32.Color4) INode {
	return &node{
		nodes: []INode{n},
		rgba: color,
	}
}
package world

import "github.com/lquesada/cavernal/model"

type ITile interface {
	model.IDrawable

	SetName(string)
	Name() string

	BaseNode() model.INode

	SetWalkThrough(bool)
	WalkThrough() bool
	SetSeeThrough(bool)
	SeeThrough() bool
	SetFallThrough(bool)
	FallThrough() bool
	SetConnect(bool)
	Connect() bool
	SetCleanable(bool)
	Cleanable() bool

	SetY(v float32)
	Y() float32
}

type tile struct{
	name string

	walkThrough bool
	seeThrough bool
	fallThrough bool

	connect bool
	cleanable bool

	y float32

	node model.INode
}

func NewTile(name string, node model.INode) ITile {
		return &tile{
				name: name,
				node: node,
		}
}

func (t *tile) BaseNode() model.INode {
	return t.node
}

func (t *tile) Node() model.INode {
	if t.node == nil {
		return nil
	}
	return t.node
}

func (t *tile) SetName(v string) {
	t.name = v
}
func (t *tile) Name() string {
	return t.name
}

func (t *tile) SetWalkThrough(v bool) {
	t.walkThrough = v
}

func (t *tile) WalkThrough() bool {
	return t.walkThrough
}

func (t *tile) SetSeeThrough(v bool) {
	t.seeThrough = v
}

func (t *tile) SeeThrough() bool {
	return t.seeThrough
}

func (t *tile) SetFallThrough(v bool) {
	t.fallThrough = v
}

func (t *tile) FallThrough() bool {
	return t.fallThrough
}

func (t *tile) SetConnect(v bool) {
	t.connect = v
}

func (t *tile) Connect() bool {
	return t.connect
}

func (t *tile) SetCleanable(v bool) {
	t.cleanable = v
}

func (t *tile) Cleanable() bool {
	return t.cleanable
}

func (t *tile) SetY(v float32) {
	t.y = v
}

func (t *tile) Y() float32 {
	return t.y
}
package helpers

import "github.com/lquesada/cavernal/model"
import "github.com/lquesada/cavernal/entity"

type BasicEntity struct {
	entity.Entity

	node model.INode
}


func (e *BasicEntity) Node() model.INode {
	return e.node
}

type BasicEntitySpecification struct {
	Name string
	Radius float32
	ShadowOpacity float32
	Height float32
	Node model.INode
}

func (s *BasicEntitySpecification) New() *BasicEntity {
	p := entity.NewEntity(s.Name)
	p.SetRadius(s.Radius)
	p.SetHeight(s.Height)
	p.SetShadowOpacity(s.ShadowOpacity)

	return &BasicEntity{
		Entity: *p,
		node: s.Node,
	}
}

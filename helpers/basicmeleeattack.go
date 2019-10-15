package helpers


import "github.com/lquesada/cavernal/entity"
import "github.com/lquesada/cavernal/model"

type BasicMeleeAttack struct {
	entity.Attack		
}


type BasicMeleeAttackSpecification struct {
	Name string
	Radius float32
	ShadowOpacity float32
	MaxWidth float32
	Reach float32
	Height float32
	Y float32
	Filter entity.Filter

	node model.INode
}

func (s *BasicMeleeAttackSpecification) New() *BasicMeleeAttack {
	p := entity.NewEntity(s.Name)
	p.SetRadius(s.Radius)
	p.SetHeight(s.Height)
	p.SetShadowOpacity(s.ShadowOpacity)

	a := entity.NewAttack(p, s.Filter, entity.GenerateCylinders(s.MaxWidth, s.Reach, 1.6, 0, 0), 1)
	a.SetReach(s.Reach)

	return &BasicMeleeAttack{
		Attack: *a,
	}
}
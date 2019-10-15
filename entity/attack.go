package entity

import	"github.com/lquesada/cavernal/model"

type IAttack interface {
	IEntity

    InnerEntity() IEntity

	Reach() float32

	TargetFilter() Filter
	Hit(c ICharacter, delta float32)
}

type Attack struct {
	Entity

	attackColision []*RelativeCylinder
	targetFilter	Filter

	reach float32
	
	damage float32

	node model.INode
}

func NewAttack(p *Entity, targetFilter Filter, attackColision []*RelativeCylinder, damage float32) *Attack {
	return &Attack{
		Entity: *p,
		targetFilter: targetFilter,
		attackColision: attackColision,
		damage: damage,
	}
}

func (e *Attack) Colision() []*Cylinder {
	return RelativeToAbsoluteList(e.attackColision, e.lookAngle)
}

func (e *Attack) InnerEntity() IEntity {
	return &e.Entity
}

func (e *Attack) PreTick() {
	e.Entity.PreTick()
}

func (e *Attack) Tick(delta float32) {
	e.Entity.Tick(delta)
}

func (e *Attack) PostTick(delta float32) {
	e.Entity.PostTick(delta)
	e.Destroy()
}

func (e *Attack) Node() model.INode {
	return e.node
}

func (e *Attack) TargetFilter() Filter {
	return e.targetFilter
}

func (e *Attack) SetTargetFilter(v Filter) {
	e.targetFilter = v
}

func (e *Attack) Reach() float32 {
	return e.reach
}

func (e *Attack) SetReach(v float32) {
	e.reach = v
}

func (e *Attack) Hit(c ICharacter, delta float32) {
	c.Damaged(e.damage)
	e.Destroy()
}


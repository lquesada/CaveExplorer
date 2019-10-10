package entity

import (
	"cavernal.com/model"
	"cavernal.com/lib/g3n/engine/math32"
)

type IEquipable interface {
	IItem

	EquippedNode() model.INode
	
	AllowRepeat() bool
	Attack(delta float32, damage float32, position *math32.Vector3, attackAngle, legHeight, handHeight, totalHeight float32, filter Filter) []IAttack
	CoverDrawSlots() []DrawSlotId

	AttackValue() int
	DefenseValue() int

	BodyType() BodyType

	IsTwoHanded() bool
}

type Equipable struct {
	Item

	// Config
	bodyType BodyType
	allowRepeat bool
	isTwoHanded bool
	attackValue int
	defenseValue int
	coverDrawSlots []DrawSlotId

	// Runtime, need to be exported
	attackIntention bool
}


func NewEquipable(i *Item, bodyType BodyType) *Equipable {
	return &Equipable{
		Item: *i,
		bodyType: bodyType,
		coverDrawSlots: []DrawSlotId{},
	}
}

func (e *Equipable) AttackValue() int {
	return e.attackValue
}

func (e *Equipable) SetAttackValue(v int) {
	e.attackValue = v
}

func (e *Equipable) DefenseValue() int {
	return e.defenseValue
}

func (e *Equipable) SetDefenseValue(v int) {
	e.defenseValue = v
}

func (e *Equipable)	CoverDrawSlots() []DrawSlotId {
	return e.coverDrawSlots
}

func (e *Equipable) SetCoverDrawSlots(v []DrawSlotId) {
	e.coverDrawSlots = v
}

func (e *Equipable) BodyType() BodyType {
 	return e.bodyType
}

func (e *Equipable) SetAllowRepeat(v bool) {
 	e.allowRepeat = v
}

func (e *Equipable) AllowRepeat() bool {
 	return e.allowRepeat
}

func (e *Equipable) SetTwoHanded(v bool) {
 	e.isTwoHanded = v
}

func (e *Equipable) IsTwoHanded() bool {
 	return e.isTwoHanded
}

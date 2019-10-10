package entity

import (
	"cavernal.com/lib/g3n/engine/math32"
	"math/rand"
)

type ICharacter interface {
    IEntity
    InnerEntity() IEntity

	MoveIntention() *math32.Vector3
	SetLookIntention(float32)
	LookIntention() float32

    Inventory() IInventory
    
	PushForce() float32

	SetDontPickUpUntilFarDistance(float32)
	DontPickUpUntilFarDistance() float32
	AddDontPickUpUntilFar(item IItem) 

	PickUp(IItem)
	WantToPickUp(IItem) bool

	StyleEquipables() map[SlotId]IEquipable

 	EquippedWeapon() IEquipable
	AttackIntention() bool
	TargetFilter() Filter

	HandHeight() float32

	WantToAttack() bool
	WantToRepeatAttack() bool
	IsAttacking() bool
	SetAttacking(bool)
	Attack(delta float32) []IAttack

	Jump()
	RepeatJump()
	IsJumping() bool
	SetJumping(bool)

    OuterDistanceTo(ICharacter) float32

	Think()
}

type Character struct {
	Entity
	
	// Config
	walkAnimationSpeed float32
	walkAnimationWhenStopping float32
	attackAnimationSpeed float32
	pushForce float32
	pushFactor float32
	handHeight float32
	targetFilter	Filter
	styleEquipables map[SlotId]IEquipable
	dontPickUpUntilFarDistance float32
	jumpMoveFactor float32

	// RunTime
	dontPickUpUntilFar map[IItem]bool
	moveIntention *math32.Vector3
	lookIntention float32
	attacking bool
	jumping bool
}

func NewCharacter(p *Entity) *Character {
	return &Character{
		Entity: *p,
		styleEquipables: map[SlotId]IEquipable{},
		dontPickUpUntilFar: map[IItem]bool{},
		moveIntention: &math32.Vector3{},
	}
}

func (e *Character) InnerEntity() IEntity {
	return &e.Entity
}

func (e *Character) WalkAnimationSpeed() float32 {
	return math32.Max(e.walkAnimationSpeed, 0.001)
}

func (e *Character) SetWalkAnimationSpeed(v float32) {
	e.walkAnimationSpeed = v
}

func (e *Character) WalkAnimationWhenStopping() float32 {
	return math32.Max(e.walkAnimationWhenStopping, 0.001)
}

func (e *Character) SetWalkAnimationWhenStopping(v float32) {
	e.walkAnimationWhenStopping = v
}

func (e *Character) AttackAnimationSpeed() float32 {
	return math32.Max(e.attackAnimationSpeed, 0.001)
}

func (e *Character) SetHandHeight(v float32) {
	e.handHeight = v
}

func (e *Character) HandHeight() float32 {
	return e.handHeight
}

func (e *Character) SetAttackAnimationSpeed(v float32) {
	e.attackAnimationSpeed = v
}

func (e *Character) PushForce() float32 {
	return e.pushForce
}

func (e *Character) SetPushForce(v float32) {
	e.pushForce = v
}

func (e *Character) PushFactor() float32 {
	return math32.Max(e.pushFactor, 0.001)
}

func (e *Character) SetPushFactor(v float32) {
	e.pushFactor = v
}

func (e *Character) TargetFilter() Filter {
	return e.targetFilter
}

func (e *Character) SetTargetFilter(v Filter) {
	e.targetFilter = v
}

func (e *Character) OuterDistanceTo(c ICharacter) float32 {
	return e.CenterDistanceTo(c)-e.OuterRadius()-c.OuterRadius()
}

func (e *Character) StyleEquipables() map[SlotId]IEquipable {
	return e.styleEquipables
}

func (e *Character) SetDontPickUpUntilFarDistance(v float32) {
	e.dontPickUpUntilFarDistance = v
}

func (e *Character) DontPickUpUntilFarDistance() float32 {
	return math32.Max(e.dontPickUpUntilFarDistance, 0.001)
}

func (e *Character) SetJumpMoveFactor(v float32) {
	e.jumpMoveFactor = v
}

func (e *Character) JumpMoveFactor() float32 {
	return e.jumpMoveFactor
}

func (e *Character) AddDontPickUpUntilFar(item IItem) {
	e.dontPickUpUntilFar[item] = true
}

func (e *Character) MoveIntention() *math32.Vector3 {
	return e.moveIntention
}

func (e *Character) SetLookIntention(v float32) {
	e.lookIntention = v
}

func (e *Character) LookIntention() float32 {
	return e.lookIntention
}

func (e *Character) PreTick() {
	e.Entity.PreTick()
	for i := range e.dontPickUpUntilFar {
		if e.Entity.BorderDistanceTo(i) > e.dontPickUpUntilFarDistance {
			delete(e.dontPickUpUntilFar, i)
		}
	}
	e.moveIntention.Zero()
	e.lookIntention = e.lookAngle
}

func (e *Character) Think() {
}

func (e *Character) WantToAttack() bool {
	return false
}

func (e *Character) WantToRepeatAttack() bool {
	return false
}

func (e *Character) Attack(delta float32) {
}

func (e *Character) Jump() {
}

func (e *Character) RepeatJump() {
}

func (e *Character) IsJumping() bool {
	return e.jumping
}

func (e *Character) SetJumping(v bool) {
	e.jumping = v
}

func (e *Character) IsAttacking() bool {
	return e.attacking
}

func (e *Character) SetAttacking(v bool) {
	e.attacking = v
}

func (e *Character) PickUp(IItem) {
}


func (e *Character) WantToPickUp(item IItem) bool {
	if e.dontPickUpUntilFar[item] {
		return false
	}
	return true
}

func (e *Character) Tick(delta float32) {
    // Acceleration
    if e.FallingToVoidPosition() == nil {
    	var jumpFactor float32 = 1
    	if e.IsJumping() {
    		jumpFactor *= e.jumpMoveFactor
    	}
    	e.speed.X += e.moveIntention.X * e.acceleration * delta * jumpFactor
    	e.speed.Z += e.moveIntention.Z * e.acceleration * delta * jumpFactor
    }

	e.Entity.Tick(delta)

	// Update look angle
	if e.lookIntention != e.lookAngle {
	    var epsilon float32 = 0.00000001
		var deltaAngle float32 = NormalizeAngle(e.lookIntention - e.lookAngle + epsilon*100 * rand.Float32() - epsilon*100/2)
		var applyDelta float32
		if deltaAngle < 0 {
			applyDelta = -(e.rotationSpeed * delta + epsilon)
			if applyDelta < deltaAngle {
				applyDelta = deltaAngle
			}
		} else {
			applyDelta = e.rotationSpeed * delta + epsilon
			if applyDelta > deltaAngle {
				applyDelta = deltaAngle
			}
		}
		e.lookAngle = NormalizeAngle(e.lookAngle+applyDelta)
	}
}

func (e *Character) PostTick(delta float32) {
	e.Entity.PostTick(delta)
}

func (e *Character) Damaged(damage float32) {
	e.Entity.Damaged(damage)
}

package entity

import (
	"cavernal.com/model"
	"cavernal.com/lib/g3n/engine/math32"
)

type IEntity interface {
	model.IDrawable
	ITickable
	Name() string
	ShadowNode() model.INode
	
	Colision() []*Cylinder
    Radius() float32
	OuterRadius() float32
	ClimbRadius() float32
	ClimbReach() float32
	Height() float32

	Position() *math32.Vector3
	FormerPosition() *math32.Vector3
	Speed() *math32.Vector3
	MaxSpeed() float32
	SetLookAngle(float32)
	LookAngle() float32
	Destroyed() bool
	Health() float32
	SetHealth(v float32)
	MaxHealth() float32
	SetMaxHealth(v float32)

	Gravity(v float32)
    SetOnGround(bool)
    OnGround() bool

    BorderDistanceTo(IEntity) float32
    CenterDistanceTo(IEntity) float32

	FallingToVoidPosition() (*math32.Vector3)
	SetFallingToVoidPosition(*math32.Vector3)

    Generate() []IEntity
	Healed(float32)
	Damaged(float32)
	Destroy()
}

type Entity struct {
	// Config
	name string
	maxSpeed float32
	friction float32
	acceleration float32
	minSpeed float32
	rotationSpeed float32
	colision []*SimpleCylinder
	radius float32
	outerRadius float32
	climbRadius float32
	climbReach float32
	height float32
	maxHealth float32
	damageable bool
	shadowOpacity float32

	// Runtime
	health	float32
	gravity float32
	speed *math32.Vector3
	shadow model.INode
	onGround bool
	formerPosition *math32.Vector3
        position *math32.Vector3
	lookAngle float32
	fallingToVoidPosition *math32.Vector3

	destroyed bool
}

func NewEntity(name string) *Entity {

	return &Entity{
		name: name,
		position: &math32.Vector3{},
		formerPosition: &math32.Vector3{},
		speed: &math32.Vector3{},
		maxHealth: 1,
		health: 1,
		colision: []*SimpleCylinder{},
		shadow: model.NewShadow(),
		shadowOpacity: 0.2,
	}
}

func (e *Entity) Name() string {
	return e.name
}

func (e *Entity) MaxHorizontalSpeed() float32 {
	return e.maxSpeed
}

func (e *Entity) SetMaxSpeed(v float32) {
	e.maxSpeed = v
}

func (e *Entity) MaxSpeed() float32 {
	return e.maxSpeed
}

func (e *Entity) Friction() float32 {
	return e.friction
}

func (e *Entity) SetFriction(v float32) {
	e.friction = v
}

func (e *Entity) Acceleration() float32 {
	return e.acceleration
}

func (e *Entity) SetAcceleration(v float32) {
	e.acceleration = v
}

func (e *Entity) ShadowOpacity() float32 {
	return e.shadowOpacity
}

func (e *Entity) SetShadowOpacity(v float32) {
	e.shadowOpacity = v
}

func (e *Entity) MinHorizontalSpeed() float32 {
	return e.minSpeed
}

func (e *Entity) SetMinSpeed(v float32) {
	e.minSpeed = v
}

func (e *Entity) RotationHorizontalSpeed() float32 {
	return e.rotationSpeed
}

func (e *Entity) SetRotationSpeed(v float32) {
	e.rotationSpeed = v
}

func (e *Entity) Colision() []*Cylinder {
	return []*Cylinder{
		SimpleToAbsolute(&SimpleCylinder{Radius: e.radius, Height: e.height}),
		}
}

func (e *Entity) Radius() float32 {
	return e.radius
}

func (e *Entity) SetRadius(v float32) {
	e.radius = v
}

func (e *Entity) OuterRadius() float32 {
	return math32.Max(0.001, math32.Max(e.outerRadius, e.radius))
}

func (e *Entity) SetOuterRadius(v float32) {
	e.outerRadius = v
}

func (e *Entity) ClimbRadius() float32 {
	return math32.Max(0.001, e.climbRadius)
}

func (e *Entity) SetClimbRadius(v float32) {
	e.climbRadius = v
}

func (e *Entity) ClimbReach() float32 {
	return math32.Max(0.001, e.climbReach)
}

func (e *Entity) SetClimbReach(v float32) {
	e.climbReach = v
}
func (e *Entity) Height() float32 {
	return e.height
}

func (e *Entity) SetHeight(v float32) {
	e.height = v
}

func (e *Entity) MaxHealth() float32 {
	return e.maxHealth
}

func (e *Entity) SetMaxHealth(v float32) {
	e.maxHealth = v
}

func (e *Entity) Damageable() bool {
	return e.damageable
}

func (e *Entity) SetDamageable(v bool) {
	e.damageable = v
}

func (e *Entity) Health() float32 {
	return e.health
}

func (e *Entity) SetHealth(v float32) {
	e.health = v
}

func (e *Entity) PreTick() {
	e.FormerPosition().Copy(e.Position())
}

func (e *Entity) Tick(delta float32) {
    // Friction
    if e.OnGround() {
	    frictionX := e.friction * e.speed.X
	    if e.speed.X < 0 {
	            e.speed.X -= frictionX * delta
	            if e.speed.X > 0 {
	                    e.speed.X = 0
	            }
	    } else if e.speed.X > 0 {
	            e.speed.X -= frictionX * delta
	            if e.speed.X < 0 {
	                    e.speed.X = 0
	            }
	    }
	    frictionZ := e.friction * e.speed.Z
	    if e.speed.Z < 0 {
	            e.speed.Z -= frictionZ * delta
	            if e.speed.Z > 0 {
	                    e.speed.Z = 0
	            }
	    } else if e.speed.Z > 0 {
	            e.speed.Z -= frictionZ * delta
	            if e.speed.Z < 0 {
	                    e.speed.Z = 0
	            }
	    }
	}

 	// Min-Max speed
	if e.HorizontalSpeed() > e.maxSpeed {
		e.SetSpeed(e.maxSpeed)
	}
	if e.HorizontalSpeed() < e.minSpeed {
		e.SetSpeed(0)
	}

    e.speed.Y += e.gravity * delta

    // Update position
    e.position.X += e.speed.X * delta
    e.position.Y += e.speed.Y * delta
    e.position.Z += e.speed.Z * delta
}

func (e *Entity) PostTick(delta float32) {
}

func (e *Entity) Position() *math32.Vector3 {
    return e.position
}

func (e *Entity) FormerPosition() *math32.Vector3 {
        return e.formerPosition
}

func (e *Entity) Speed() *math32.Vector3 {
        return e.speed
}

func (e *Entity) SetLookAngle(v float32) {
        e.lookAngle = v
}

func (e *Entity) LookAngle() float32 {
        return e.lookAngle
}

func (e *Entity) BorderDistanceTo(c IEntity) float32 {
	return e.CenterDistanceTo(c)-e.Radius()-c.Radius()
}

func (e *Entity) CenterDistanceTo(c IEntity) float32 {
	return Distance2D(e.position, c.Position())
}

func (e *Entity) HorizontalSpeed() float32 {
        return math32.Sqrt(e.speed.X * e.speed.X + e.speed.Z * e.speed.Z)
}

func (e *Entity) SetSpeed(v float32) {
	speed := e.HorizontalSpeed()
	var delta float32 = 1
	if speed > 0 {
		delta = v / e.HorizontalSpeed()
	}
	e.speed.X *= delta
	e.speed.Z *= delta
}

func (e *Entity) Destroy() {
	e.destroyed = true
}

func (e *Entity) Destroyed() bool {
	return e.destroyed
}

func (e *Entity) Healed(points float32) {
	e.health += points
	if e.health > e.maxHealth {
		e.health = e.maxHealth
	}
}

func (e *Entity) Damaged(damage float32) {
	if !e.damageable {
		return
	}
	e.health -= damage
	if e.health <= 0 {
		e.health = 0
		e.Destroy()
	}
}

func (e *Entity) Node() model.INode {
	return nil
}

func (e *Entity) ShadowNode() model.INode {
  	if e.radius > 0 && e.shadowOpacity > 0 {
		return e.shadow.
			RGBA(&math32.Color4{0, 0, 0, e.shadowOpacity}).
			Transform(
				&model.Transform{
					Rotation: &math32.Vector3{-math32.Pi/2, 0, 0},
					},
				&model.Transform{
					Scale: &math32.Vector3{e.radius, 1, e.radius},
					Position: &math32.Vector3{e.position.X, 0.01, e.position.Z},
				},
				)
	}
	return nil
}

func (e *Entity) Generate() []IEntity {
	return nil
}

func (e *Entity) Gravity(v float32) {
	e.gravity = v
}

func (e *Entity) SetOnGround(v bool) {
	e.onGround = v
}

func (e *Entity) OnGround() bool {
	return e.onGround
}

func (e *Entity) FallingToVoidPosition() (v *math32.Vector3) {
	return e.fallingToVoidPosition
}

func (e *Entity) SetFallingToVoidPosition(v *math32.Vector3) {
	e.fallingToVoidPosition = v
}

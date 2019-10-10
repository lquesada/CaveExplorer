package helpers

import "cavernal.com/model"
import "cavernal.com/entity"
import    "hash/fnv"
import    "math"
import "cavernal.com/lib/g3n/engine/math32"

type BasicEquipable struct {
	entity.Equipable

	cloneSpecification *BasicEquipableSpecification

 	attackGenerator func() *BasicMeleeAttack

	equippedNode model.INode
	droppedNode model.INode
	inventoryNode model.INode
}


func (e *BasicEquipable) EquippedNode() model.INode {
	return e.equippedNode
}


func (e *BasicEquipable) DroppedNode() model.INode {
	return e.droppedNode
}

func (e *BasicEquipable) InventoryNode() model.INode {
	return e.inventoryNode
}

func (e *BasicEquipable) Node() model.INode {
	return e.DroppedNode()
}

func (e *BasicEquipable) Attack(delta float32, damage float32, position *math32.Vector3, attackAngle, legHeight, handHeight, totalHeight float32, filter entity.Filter) []entity.IAttack {
	a := e.attackGenerator()
	a.SetTargetFilter(filter)
	a.Position().X = position.X + math32.Sin(attackAngle) * a.Reach()/2
	a.Position().Y = legHeight
	a.Position().Z = position.Z + math32.Cos(attackAngle) * a.Reach()/2
	return []entity.IAttack{a}
}

func (e *BasicEquipable) SimpleValue() float32 {
     h := fnv.New32a()
     h.Write([]byte(e.Equipable.InnerEntity().Name()))
    hash := h.Sum32()
    return float32(e.Equipable.AttackValue() + e.Equipable.DefenseValue()) + float32(hash)/float32(math.MaxUint32)
}

func (e *BasicEquipable) Clone() entity.IItem {
  return e.cloneSpecification.New()
}

type BasicEquipableSpecification struct {
	Name string
	ItemType entity.ItemType
	BodyType entity.BodyType
	Radius float32
	ShadowOpacity float32
	Height float32
	AllowRepeat bool
	IsTwoHanded bool
	AttackValue int
	DefenseValue int
	AttackGenerator func() *BasicMeleeAttack
	EquippedNode model.INode
	DroppedNode model.INode
	InventoryNode model.INode
	CountableId entity.CountableId
	ProvidesAmmo entity.AmmoType
	RequiresAmmo entity.AmmoType
	Category entity.ItemCategory
}

func (s *BasicEquipableSpecification) New() *BasicEquipable {
	p := entity.NewEntity(s.Name)
	p.SetRadius(s.Radius)
	p.SetHeight(s.Height)
	p.SetShadowOpacity(s.ShadowOpacity)

	i := entity.NewItem(p, s.ItemType)
	i.SetCountableId(s.CountableId)
	i.SetProvidesAmmo(s.ProvidesAmmo)
	i.SetRequiresAmmo(s.RequiresAmmo)
	i.SetCategory(s.Category)

	e := entity.NewEquipable(i, s.BodyType)
	e.SetAllowRepeat(s.AllowRepeat)
	e.SetTwoHanded(s.IsTwoHanded)
	e.SetAttackValue(s.AttackValue)
	e.SetDefenseValue(s.DefenseValue)

	return &BasicEquipable{
		Equipable: *e,
		equippedNode: s.EquippedNode,
		droppedNode: s.DroppedNode,
		inventoryNode: s.InventoryNode,
		attackGenerator: s.AttackGenerator,
		cloneSpecification: s,
	}
}

package entity

import	"github.com/lquesada/cavernal/model"

type IItem interface {
	IEntity
    InnerEntity() IEntity

	DroppedNode() model.INode
	InventoryNode() model.INode

   	WaitBeforePickupable(seconds float32)
	IsPickupable() bool

	SimpleValue() float32

	ItemType() ItemType

	ProvidesAmmo() AmmoType
	SetProvidesAmmo(v AmmoType)
	RequiresAmmo() AmmoType
	SetRequiresAmmo(v AmmoType)

	SetCountableId(countableId CountableId)
	IsCountable() bool
	CountableId() CountableId
	SetCount(v int)
	Count() int

	SetCategory(ItemCategory)
	Category() ItemCategory

	Clone() IItem
}

var (
	OneHandedWeapon ItemType = "OneHandedWeapon"
	TwoHandedWeapon ItemType = "TwoHandedWeapon"
	Shield ItemType = "Shield"
	Other ItemType = "Other"

	WeaponItemType = map[ItemType]bool{
		OneHandedWeapon: true,
		TwoHandedWeapon: true,
	}

	HandItemType  = map[ItemType]bool{
		OneHandedWeapon: true,
		Shield: true,
		TwoHandedWeapon: true,
		Other: true,
	}
)

type Item struct {
	Entity // When dropped

	// Config
	itemType ItemType
	countableId CountableId
	providesAmmo AmmoType
	requiresAmmo AmmoType
	
	// Runtime, need to be exported
	pickupableSequence *model.Sequence
	count int
	
	category ItemCategory
}

func NewItem(p *Entity, itemType ItemType) *Item {
	return &Item{
		Entity: *p,
		category: Common,
		itemType: itemType,
		pickupableSequence: model.NewTimer(1),
		count: 1,
	}
}

func (e *Item) InnerEntity() IEntity {
	return &e.Entity
}

func (e *Item) PreTick() {
	e.Entity.PreTick()
}

func (e *Item) Tick(delta float32) {
	e.Entity.Tick(delta)
	e.pickupableSequence.Add(delta)
}

func (e *Item) PostTick(delta float32) {
	e.Entity.PostTick(delta)
}

func (e *Item) WaitBeforePickupable(seconds float32) {
	e.pickupableSequence = model.NewTimer(seconds)
	e.pickupableSequence.Start()
}

func (e *Item) IsPickupable() bool {
	return !e.pickupableSequence.Running()
}

func (e *Item) ItemType() ItemType {
	return e.itemType
}

func (e *Item) IsCountable() bool {
	if e.countableId != "" {
		return true
	}
	return false
}

func (e *Item) SetCountableId(countableId CountableId) {
	e.countableId = countableId
}

func (e *Item) CountableId() CountableId {
	return e.countableId
}

func (e *Item) Count() int {
	return e.count
}

func (e *Item) SetCount(v int) {
	e.count = v
}

func (e *Item) ProvidesAmmo() AmmoType {
	return e.providesAmmo
}

func (e *Item) SetProvidesAmmo(v AmmoType) {
	e.providesAmmo = v
}

func (e *Item) RequiresAmmo() AmmoType {
	return e.requiresAmmo
}

func (e *Item) SetRequiresAmmo(v AmmoType) {
	e.requiresAmmo = v
}

func CanStack(source, destination IItem) bool {
	if source.IsCountable() && destination.IsCountable() && source.CountableId() == destination.CountableId() {
		return true
	}
	return false
}

func Stack(source, destination IItem) {
	if source.IsCountable() && destination.IsCountable() && source.CountableId() == destination.CountableId() {
		destination.SetCount(destination.Count() + source.Count())
		source.SetCount(0)
	}
}

func Split(source IItem, amount int) IItem {
  if !source.IsCountable() || source.Count() <= amount {
    return nil
  }
  destination := source.Clone()
  destination.SetCount(amount)
  source.SetCount(source.Count() - amount)
  return destination
}

func (e *Item) Category() ItemCategory {
	return e.category
}

func (e *Item) SetCategory(v ItemCategory) {
	e.category = v
}
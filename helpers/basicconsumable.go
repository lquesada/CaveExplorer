package helpers

import "cavernal.com/model"
import "cavernal.com/entity"

type BasicConsumable struct {
	entity.Consumable
	
	cloneSpecification *BasicConsumableSpecification
	
	droppedNode model.INode
	inventoryNode model.INode
}

func (e *BasicConsumable) DroppedNode() model.INode {
	return e.droppedNode
}

func (e *BasicConsumable) InventoryNode() model.INode {
	return e.inventoryNode
}

func (e *BasicConsumable) Node() model.INode {
	return e.DroppedNode()
}

func (e *BasicConsumable) SimpleValue() float32 {
	return 1 // TODO
}

func (e *BasicConsumable) Clone() entity.IItem {
  return e.cloneSpecification.New()
}

type BasicConsumableSpecification struct {
	Name string
	ItemType entity.ItemType
	Radius float32
	ShadowOpacity float32
	Height float32
	DroppedNode model.INode
	InventoryNode model.INode
	CountableId entity.CountableId
	Category entity.ItemCategory
}

func (s *BasicConsumableSpecification) New() *BasicConsumable {
	p := entity.NewEntity(s.Name)
	p.SetRadius(s.Radius)
	p.SetHeight(s.Height)
	p.SetShadowOpacity(s.ShadowOpacity)

	i := entity.NewItem(p, s.ItemType)
	i.SetCountableId(s.CountableId)
	i.SetCategory(s.Category)

	c := entity.NewConsumable(i)


	return &BasicConsumable{
		Consumable: *c,
		droppedNode: s.DroppedNode,
		inventoryNode: s.InventoryNode,
		cloneSpecification: s,
	}
}
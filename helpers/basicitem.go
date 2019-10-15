package helpers

import "github.com/lquesada/cavernal/model"
import "github.com/lquesada/cavernal/entity"

type BasicItem struct {
	entity.Item
	
	cloneSpecification *BasicItemSpecification

	droppedNode model.INode
	inventoryNode model.INode
}


func (e *BasicItem) DroppedNode() model.INode {
	return e.droppedNode
}

func (e *BasicItem) InventoryNode() model.INode {
	return e.inventoryNode
}


func (e *BasicItem) Node() model.INode {
	return e.DroppedNode()
}

func (e *BasicItem) SimpleValue() float32 {
	return 1 // TODO
}

func (e *BasicItem) Clone() entity.IItem {
  return e.cloneSpecification.New()
}


type BasicItemSpecification struct {
	Name string
	ItemType entity.ItemType
	Radius float32
	Height float32
	DroppedNode model.INode
	InventoryNode model.INode
	CountableId entity.CountableId
	ShadowOpacity float32
	Category entity.ItemCategory
}

func (s *BasicItemSpecification) New() *BasicItem {
	p := entity.NewEntity(s.Name)
	p.SetRadius(s.Radius)
	p.SetHeight(s.Height)
	p.SetShadowOpacity(s.ShadowOpacity)

	i := entity.NewItem(p, s.ItemType)
	i.SetCountableId(s.CountableId)
	i.SetCategory(s.Category)

	return &BasicItem{
		Item: *i,
		droppedNode: s.DroppedNode,
		inventoryNode: s.InventoryNode,
		cloneSpecification: s,
	}
}

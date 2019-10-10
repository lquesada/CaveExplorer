package entity

type IConsumable interface {
	IItem

	Consume()
}

type Consumable struct {
	Item
}

func NewConsumable(i *Item) *Consumable {
	return &Consumable{
		Item: *i,
	}
}

func (c *Consumable) Consume() {
	
}
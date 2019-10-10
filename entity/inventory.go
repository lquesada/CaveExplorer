package entity

type IInventory interface{
	ITickable
	
	Slots() map[SlotId]IItem
	EquipmentSlots() map[SlotId]IEquipable
	CanEquip(itemSlot int, slot SlotId) bool
	CanEquipItem(item IItem, slot SlotId) bool
	AutoEquip(itemSlot int)
	Equip(itemSlot int, slot SlotId)
	AutoUnequip(slot SlotId)
	Unequip(slot SlotId, itemSlot int)
	ShuffleEquipment(slot1, slot2 SlotId) bool

	Items() []IItem

	EquippedWeapon() IEquipable

	WhereHas(IItem) int
	WhereEquipped(item IItem) SlotId

	CanPickUp(item IItem) bool
	PickUp(item IItem) int
	ShuffleItem(itemSlot1, itemSlot2 int)	

	SetHoldAmount(v int)
	HoldAmount() int
	SetHoldItem(item IItem)
	HoldItem() IItem

	RemoveItemSlot(itemSlot int)
	RemoveEquipped(slot SlotId)

	DropItemSlot(itemSlot int)
	DropEquipped(slot SlotId)
	DropAllLoot()
	Generate() []IItem

	Flip()

	Attack(bool) interface{}
}